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

	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/configgtm"

	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/models"
	"github.com/sapcc/andromeda/internal/utils"
)

func (s *AkamaiAgent) SyncProperty(domain *models.Domain) error {
	var members []*models.Member
	var monitors []*models.Monitor

	pools := domain.GetPools()
	if len(pools) > 0 {
		// flatten Members and Monitors
		for _, pool := range pools {
			members = append(members, pool.GetMembers()...)
			monitors = append(monitors, pool.GetMonitors()...)
		}
	}

	var datacenters []*gtm.Datacenter
	datacenters, err := s.gtm.ListDatacenters(context.Background(), config.Global.AkamaiConfig.Domain)
	if err != nil {
		return err
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
	}

	// Add new Members
	for _, member := range members {
		trafficTarget := gtm.TrafficTarget{
			Name:    member.GetId(),
			Enabled: member.GetAdminStateUp(),
			Servers: []string{utils.InetNtoa(member.Address).String()},
			Weight:  50,
		}
		datacenterID := member.GetDatacenter()

		if len(datacenterID) > 0 {
			found := false
			for _, datacenter := range datacenters {
				if datacenter.Nickname == datacenterID {
					trafficTarget.DatacenterId = datacenter.DatacenterId
					found = true
				}
			}

			if !found {
				// Create Datacenter
				aDatacenter, err := s.GetDatacenter(datacenterID)
				if err != nil {
					return err
				}
				datacenter, err := s.SyncDatacenter(aDatacenter)
				if err != nil {
					return err
				}
				trafficTarget.DatacenterId = datacenter.DatacenterId
				datacenters = append(datacenters, datacenter)
			}
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
		case models.Monitor_HTTP:
			if monitor.GetSend() == "" {
				livenessTest.TestObject = "/"
			} else {
				livenessTest.TestObject = monitor.GetSend()
			}
		case models.Monitor_TCP:
			livenessTest.RequestString = monitor.GetSend()
			livenessTest.ResponseString = monitor.GetReceive()
		}
		property.LivenessTests = append(property.LivenessTests, &livenessTest)
	}

	if _, err := s.gtm.UpdateProperty(context.Background(), &property, config.Global.AkamaiConfig.Domain); err != nil {
		return err
	}

	return nil
}
