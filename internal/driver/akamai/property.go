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
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_4"

	"github.com/sapcc/andromeda/internal/models"
	"github.com/sapcc/andromeda/internal/utils"
)

func (s *AkamaiAgent) SyncProperty(domain *models.Domain, members []*models.Member, monitors []*models.Monitor) (*configgtm.Property, error) {
	var datacenters []*configgtm.Datacenter
	datacenters, err := configgtm.ListDatacenters(domain.GetFqdn())
	if err != nil {
		return nil, err
	}

	// Create property
	property := configgtm.NewProperty("origin")
	property.Type = PROPERTY_TYPE_MAP[domain.GetMode()]
	property.Comments = domain.Id
	property.ScoreAggregationType = "mean"
	property.HandoutMode = "normal"
	property.FailbackDelay = 0
	property.FailoverDelay = 0

	// Add new Members
	for _, member := range members {
		trafficTarget := property.NewTrafficTarget()
		trafficTarget.Name = member.GetId()
		trafficTarget.Enabled = member.GetAdminStateUp()
		trafficTarget.Servers = []string{utils.InetNtoa(member.Address).String()}
		trafficTarget.Weight = 50

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
					return nil, err
				}
				datacenter, err := s.SyncDatacenter(domain, aDatacenter)
				if err != nil {
					return nil, err
				}
				trafficTarget.DatacenterId = datacenter.DatacenterId
				datacenters = append(datacenters, datacenter)
			}
		}
		property.TrafficTargets = append(property.TrafficTargets, trafficTarget)
	}

	// Add new Monitors
	for _, monitor := range monitors {
		livenessTest := property.NewLivenessTest(
			monitor.GetId(),
			MONITOR_LIVENESS_TYPE_MAP[monitor.GetType()],
			int(monitor.GetInterval()),
			float32(monitor.GetTimeout()))
		livenessTest.Disabled = !monitor.GetAdminStateUp()

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
		property.LivenessTests = append(property.LivenessTests, livenessTest)
	}

	/*if _, err := property.Update(domain.GetFqdn()); err != nil {
		return err
	}*/

	return property, nil
}
