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

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/gtm"
	"github.com/apex/log"
	"github.com/go-openapi/swag"

	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/driver"
	"github.com/sapcc/andromeda/internal/rpcmodels"
	"github.com/sapcc/andromeda/internal/utils"
	"github.com/sapcc/andromeda/models"
)

func (s *AkamaiAgent) DeleteProperty(domain *rpcmodels.Domain, trafficManagementDomain string) error {
	// Delete
	log.Infof("DeleteProperty(domain=%s, property=%s)", trafficManagementDomain, domain.GetFqdn())

	property, err := s.gtm.GetProperty(context.Background(), domain.GetFqdn(), config.Global.AkamaiConfig.Domain)
	if err != nil {
		log.Warnf("Property '%s' doesn't exist...", domain.GetFqdn())
		return nil
	}

	ret, err := s.gtm.DeleteProperty(context.Background(), property, trafficManagementDomain)
	if err != nil {
		return fmt.Errorf("request %s: %w", PrettyJson(property), err)
	}

	log.Debugf("Request: %s\nResponse: %s",
		PrettyJson(property),
		PrettyJson(ret))
	return nil
}

func (s *AkamaiAgent) SyncProperty(domain *rpcmodels.Domain, trafficManagementDomain string) (ProvRequests, error) {
	var provRequests ProvRequests
	var members []*rpcmodels.Member
	var monitors []*rpcmodels.Monitor

	pools := domain.GetPools()
	if len(pools) > 0 {
		// flatten Members and Monitors
		for _, pool := range pools {
			if pool.ProvisioningStatus == models.PoolProvisioningStatusPENDINGDELETE {
				provRequests = append(provRequests,
					driver.GetProvisioningStatusRequest(pool.Id, "POOL", "DELETED"))
				continue
			}
			members = append(members, pool.GetMembers()...)
			monitors = append(monitors, pool.GetMonitors()...)
			provRequests = append(provRequests,
				driver.GetProvisioningStatusRequest(pool.Id, "POOL", models.PoolProvisioningStatusACTIVE))
		}
	}

	// Create property
	property := gtm.Property{
		Name:                 domain.GetFqdn(),
		Type:                 PROPERTY_TYPE_MAP[domain.GetMode()],
		Comments:             domain.Id,
		ScoreAggregationType: "best",
		HandoutMode:          "all-live-ips",
		FailbackDelay:        0,
		FailoverDelay:        0,
		TrafficTargets:       []*gtm.TrafficTarget{},
		LivenessTests:        []*gtm.LivenessTest{},
	}

	// Process Members
MEMBERLOOP:
	for _, member := range members {
		if member.ProvisioningStatus == models.MemberProvisioningStatusPENDINGDELETE {
			provRequests = append(provRequests,
				driver.GetProvisioningStatusRequest(member.Id, "MEMBER", "DELETED"))
			continue
		}

		datacenterUUID := member.GetDatacenter()
		var datacenterID int
		if len(datacenterUUID) > 0 {
			var err error
			datacenterID, err = s.GetDatacenterMeta(datacenterUUID, domain.Datacenters)
			if err != nil {
				log.WithField("datacenterUUID", datacenterUUID).
					WithError(err).Error("Failed to get datacenter meta")
				continue
			}

			// check if we have already a traffic target for this datacenter
			for _, target := range property.TrafficTargets {
				if target.DatacenterID == datacenterID {
					// just add the server to the existing traffic target
					target.Servers = append(target.Servers, utils.InetNtoa(member.Address).String())
					provRequests = append(provRequests,
						driver.GetProvisioningStatusRequest(member.Id, "MEMBER", "ACTIVE"))
					continue MEMBERLOOP
				}
			}
		}

		// Add new traffic target
		trafficTarget := gtm.TrafficTarget{
			Name:         member.GetDatacenter(),
			Enabled:      member.GetAdminStateUp(),
			Servers:      []string{utils.InetNtoa(member.Address).String()},
			Weight:       50,
			DatacenterID: datacenterID,
		}
		property.TrafficTargets = append(property.TrafficTargets, &trafficTarget)
		provRequests = append(provRequests,
			driver.GetProvisioningStatusRequest(member.Id, "MEMBER", "ACTIVE"))
	}

	// due shortcoming of individual liveness test per member, we have to replicate a monitor per unique member port
	// collect unique member ports
	uniquePorts := make(map[uint32]interface{})
	for _, member := range members {
		uniquePorts[member.GetPort()] = nil
	}

	// Add new Monitors
	for _, monitor := range monitors {
		if monitor.ProvisioningStatus == models.MonitorProvisioningStatusPENDINGDELETE {
			provRequests = append(provRequests,
				driver.GetProvisioningStatusRequest(monitor.Id, "MONITOR", "DELETED"))
			continue
		}

	monitorLoop:
		for testPort := range uniquePorts {
			livenessTest := gtm.LivenessTest{
				Name:               fmt.Sprintf("%s-%d", monitor.GetId(), testPort),
				TestObjectPort:     int(testPort),
				TestObjectProtocol: MONITOR_LIVENESS_TYPE_MAP[monitor.GetType()],
				TestInterval:       int(monitor.GetInterval()),
				TestTimeout:        float32(monitor.GetTimeout()),
				Disabled:           !monitor.GetAdminStateUp(),
			}

			switch monitor.GetType() {
			case rpcmodels.Monitor_HTTPS:
				fallthrough
			case rpcmodels.Monitor_HTTP:
				if monitor.GetSend() == "" {
					livenessTest.TestObject = "/"
				} else {
					livenessTest.TestObject = monitor.GetSend()
				}
				livenessTest.HTTPHeaders = []*gtm.HTTPHeader{{Name: "Host", Value: domain.GetFqdn()}}
				if domainName := monitor.GetDomainName(); domainName != "" {
					livenessTest.HTTPHeaders = []*gtm.HTTPHeader{{Name: "Host", Value: domainName}}
				}
				livenessTest.HTTPMethod = swag.String(monitor.GetMethod().String())
			case rpcmodels.Monitor_TCP:
				livenessTest.RequestString = monitor.GetSend()
				livenessTest.ResponseString = monitor.GetReceive()
			default:
				// unsupported type
				log.Warnf("Unsupported monitor type: %s", monitor.GetType())
				provRequests = append(provRequests,
					driver.GetProvisioningStatusRequest(monitor.Id, "MONITOR", models.MonitorProvisioningStatusERROR))
				continue monitorLoop
			}
			property.LivenessTests = append(property.LivenessTests, &livenessTest)
		}
		provRequests = append(provRequests,
			driver.GetProvisioningStatusRequest(monitor.Id, "MONITOR", models.MonitorProvisioningStatusACTIVE))
	}

	provRequests = append(provRequests,
		driver.GetProvisioningStatusRequest(domain.Id, "DOMAIN", models.DomainProvisioningStatusACTIVE))

	// Pre-Validation
	if len(property.TrafficTargets) == 0 {
		// Need traffictargets with datacenters before posting
		log.Debugf("Skipping Property '%s': No traffic targets", property.Name)
		return provRequests, nil
	}

	existingProperty, err2 := s.gtm.GetProperty(context.Background(), property.Name, config.Global.AkamaiConfig.Domain)
	if err2 != nil {
		log.Debugf("Property '%s' doesn't exist, creating...", property.Name)
	}

	fieldsToCompare := []string{
		"Name",
		"Type",
		"Comments",
		"HandoutMode",
		"TrafficTargets",
		"ScoreAggregationType",
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
		return provRequests, nil
	}

	// Update
	log.Infof("UpdateProperty(domain=%s, property=%s)", trafficManagementDomain, property.Name)
	ret, err3 := s.gtm.UpdateProperty(context.Background(), &property, trafficManagementDomain)
	if err3 != nil {
		return nil, fmt.Errorf("request %s: %w", PrettyJson(property), err3)
	}

	log.Debugf("Request: %s\nResponse: %s",
		PrettyJson(property),
		PrettyJson(ret))
	return provRequests, nil
}
