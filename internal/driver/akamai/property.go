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
	"go-micro.dev/v4/logger"

	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/configgtm"

	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/rpcmodels"
	"github.com/sapcc/andromeda/internal/utils"
)

func addTrafficTarget(property *gtm.Property, member *rpcmodels.Member) bool {
	for _, trafficTarget := range property.TrafficTargets {
		if trafficTarget.Name == member.GetDatacenter() {
			trafficTarget.Servers = append(trafficTarget.Servers, utils.InetNtoa(member.Address).String())
			return true
		}
	}
	return false
}

func (s *AkamaiAgent) SyncProperty(domain *rpcmodels.Domain) error {
	var members []*rpcmodels.Member
	var monitors []*rpcmodels.Monitor

	pools := domain.GetPools()
	if len(pools) > 0 {
		// flatten Members and Monitors
		for _, pool := range pools {
			members = append(members, pool.GetMembers()...)
			monitors = append(monitors, pool.GetMonitors()...)
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
		// Add member to existing traffic target within the same Datacenter
		if addTrafficTarget(&property, member) {
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
			var aDatacenter *rpcmodels.Datacenter
			for _, datacenter := range domain.Datacenters {
				if datacenter.GetId() == datacenterUUID {
					aDatacenter = datacenter
					break
				}
			}

			// Sync datacenter first
			var err error
			aDatacenter, err = s.SyncDatacenter(aDatacenter, false)
			if err != nil {
				return err
			}

			// DatacenterId is a unique number for an akamai datacenter
			trafficTarget.DatacenterId = int(aDatacenter.GetMeta())
		}
		property.TrafficTargets = append(property.TrafficTargets, &trafficTarget)
	}

	// Add new Monitors
	for _, monitor := range monitors {
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
		case rpcmodels.Monitor_TCP:
			livenessTest.RequestString = monitor.GetSend()
			livenessTest.ResponseString = monitor.GetReceive()
		}
		property.LivenessTests = append(property.LivenessTests, &livenessTest)
	}

	// Pre-Validation
	if len(property.TrafficTargets) == 0 {
		// Need traffictargets with datacenters
		return nil
	}

	existingProperty, err := s.gtm.GetProperty(context.Background(), property.Name, config.Global.AkamaiConfig.Domain)
	if err != nil {
		logger.Warnf("Property %s doesn't exist, creating...", property.Name)
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
		"TrafficTargets.Name",
		"LivenessTests",
		"LivenessTests.Name",
		"LivenessTests.TestObject",
		"LivenessTests.TestInterval",
		"LivenessTests.TestTimeout",
		"LivenessTests.RequestString",
		"LivenessTests.ResponseString",
		"LivenessTests.TestObjectProtocol",
	}
	if !utils.DeepEqualFields(&property, existingProperty, fieldsToCompare) {
		// Update
		logger.Infof("UpdateProperty(%s) of domain %s", property.Name, config.Global.AkamaiConfig.Domain)
		_, err := s.gtm.UpdateProperty(context.Background(), &property, config.Global.AkamaiConfig.Domain)
		if err != nil {
			return err
		}

	}
	return nil
}
