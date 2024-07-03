/*
 *   Copyright 2022 SAP SE
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package akamai

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/gtm"
	"github.com/go-openapi/swag"
	"go-micro.dev/v4/logger"

	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/driver"
	"github.com/sapcc/andromeda/internal/rpcmodels"
	"github.com/sapcc/andromeda/internal/utils"
)

func (s *AkamaiAgent) DeleteProperty(domain *rpcmodels.Domain, trafficManagementDomain string) error {
	// Delete
	logger.Infof("DeleteProperty(domain=%s, property=%s)", trafficManagementDomain, domain.GetFqdn())

	property, err := s.gtm.GetProperty(context.Background(), domain.GetFqdn(), config.Global.AkamaiConfig.Domain)
	if err != nil {
		logger.Warnf("Property '%s' doesn't exist...", domain.GetFqdn())
		return nil
	}

	ret, err := s.gtm.DeleteProperty(context.Background(), property, trafficManagementDomain)
	if err != nil {
		return fmt.Errorf("Request %s: %w", PrettyJson(property), err)
	}

	if !logger.V(logger.DebugLevel, nil) {
		logger.Debugf("Request: %s\nResponse: %s",
			PrettyJson(property),
			PrettyJson(ret))
	}
	return nil
}

func (s *AkamaiAgent) SyncProperty(domain *rpcmodels.Domain, trafficManagementDomain string) (provRequests ProvRequests, err error) {
	var members []*rpcmodels.Member
	var monitors []*rpcmodels.Monitor

	pools := domain.GetPools()
	if len(pools) > 0 {
		// flatten Members and Monitors
		for _, pool := range pools {
			if pool.ProvisioningStatus == "PENDING_DELETE" {
				provRequests = append(provRequests,
					driver.GetProvisioningStatusRequest(pool.Id, "POOL", "DELETED"))
				continue
			}
			members = append(members, pool.GetMembers()...)
			monitors = append(monitors, pool.GetMonitors()...)
			provRequests = append(provRequests,
				driver.GetProvisioningStatusRequest(pool.Id, "POOL", "ACTIVE"))
		}
	}

	// Create property
	property := gtm.Property{
		Name:                 domain.GetFqdn(),
		Type:                 PROPERTY_TYPE_MAP[domain.GetMode()],
		Comments:             domain.Id,
		ScoreAggregationType: "mean",
		HandoutMode:          "normal",
		FailbackDelay:        0,
		FailoverDelay:        0,
		TrafficTargets:       []*gtm.TrafficTarget{},
		LivenessTests:        []*gtm.LivenessTest{},
	}

	// Add new Members
	for _, member := range members {
		if member.ProvisioningStatus == "PENDING_DELETE" {
			provRequests = append(provRequests,
				driver.GetProvisioningStatusRequest(member.Id, "MEMBER", "DELETE"))
			continue
		}

		// Add new traffic target
		trafficTarget := gtm.TrafficTarget{
			Name:    member.GetDatacenter(),
			Enabled: member.GetAdminStateUp(),
			Servers: []string{utils.InetNtoa(member.Address).String()},
			Weight:  50,
		}
		datacenterUUID := member.GetDatacenter()

		if len(datacenterUUID) > 0 {
			if trafficTarget.DatacenterID, err = s.GetDatacenterMeta(datacenterUUID, domain.Datacenters); err != nil {
				logger.Error(err)
				continue
			}
		}
		property.TrafficTargets = append(property.TrafficTargets, &trafficTarget)
		provRequests = append(provRequests,
			driver.GetProvisioningStatusRequest(member.Id, "MEMBER", "ACTIVE"))
	}

	// Add new Monitors
	for _, monitor := range monitors {
		if monitor.ProvisioningStatus == "PENDING_DELETE" {
			provRequests = append(provRequests,
				driver.GetProvisioningStatusRequest(monitor.Id, "MONITOR", "DELETED"))
			continue
		}

		livenessTest := gtm.LivenessTest{
			Name:               monitor.GetId(),
			TestObjectProtocol: MONITOR_LIVENESS_TYPE_MAP[monitor.GetType()],
			TestInterval:       int(monitor.GetInterval()),
			TestTimeout:        float32(monitor.GetTimeout()),
			Disabled:           !monitor.GetAdminStateUp(),
		}

		switch monitor.GetType() {
		case rpcmodels.Monitor_HTTP:
			if monitor.GetSend() == "" {
				livenessTest.TestObject = "/"
			} else {
				livenessTest.TestObject = monitor.GetSend()
			}
			var testPort uint32 = 80
			for _, member := range members {
				testPort = member.GetPort()
				break
			}
			livenessTest.TestObjectPort = int(testPort)
			livenessTest.HTTPHeaders = []*gtm.HTTPHeader{{Name: "Host", Value: domain.GetFqdn()}}
			livenessTest.HTTPMethod = swag.String(monitor.GetMethod().String())
		case rpcmodels.Monitor_TCP:
			livenessTest.RequestString = monitor.GetSend()
			livenessTest.ResponseString = monitor.GetReceive()
		default:
			// unsupported type
			provRequests = append(provRequests,
				driver.GetProvisioningStatusRequest(monitor.Id, "MONITOR", "ERROR"))
			continue
		}
		property.LivenessTests = append(property.LivenessTests, &livenessTest)
		provRequests = append(provRequests,
			driver.GetProvisioningStatusRequest(monitor.Id, "MONITOR", "ACTIVE"))
	}

	provRequests = append(provRequests,
		driver.GetProvisioningStatusRequest(domain.Id, "DOMAIN", "ACTIVE"))

	// Pre-Validation
	if len(property.TrafficTargets) == 0 {
		// Need traffictargets with datacenters before posting
		logger.Debugf("Skipping Property '%s': No traffic targets", property.Name)
		return
	}

	existingProperty, err2 := s.gtm.GetProperty(context.Background(), property.Name, config.Global.AkamaiConfig.Domain)
	if err2 != nil {
		logger.Debugf("Property '%s' doesn't exist, creating...", property.Name)
	}

	fieldsToCompare := []string{
		"Name",
		"Type",
		"Comments",
		"TrafficTargets",
		"TrafficTargets.DatacenterId",
		"TrafficTargets.Enabled",
		"TrafficTargets.Weight",
		"TrafficTargets.Servers",
		//"TrafficTargets.Name", # bug in Akamai API
		"LivenessTests",
		"LivenessTests.Name",
		"LivenessTests.TestObject",
		"LivenessTests.TestObjectPort",
		"LivenessTests.TestInterval",
		"LivenessTests.TestTimeout",
		"LivenessTests.RequestString",
		"LivenessTests.ResponseString",
		"LivenessTests.TestObjectProtocol",
		"LivenessTests.HTTPMethod",
	}
	if utils.DeepEqualFields(&property, existingProperty, fieldsToCompare) {
		return
	}

	// Update
	logger.Infof("UpdateProperty(domain=%s, property=%s)", trafficManagementDomain, property.Name)
	ret, err3 := s.gtm.UpdateProperty(context.Background(), &property, trafficManagementDomain)
	if err3 != nil {
		err = fmt.Errorf("request %s: %w", PrettyJson(property), err3)
		return
	}

	if !logger.V(logger.DebugLevel, nil) {
		logger.Debugf("Request: %s\nResponse: %s",
			PrettyJson(property),
			PrettyJson(ret))
	}
	return
}
