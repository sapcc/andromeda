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
	"github.com/sapcc/andromeda/internal/driver"
	"github.com/sapcc/andromeda/internal/models"
	"github.com/sapcc/andromeda/internal/rpc/server"
)

func (s *AkamaiAgent) UpdateDomainProvisioningStatus(domain *models.Domain, value string) error {
	provisioningStatusRequests := []*server.ProvisioningStatusRequest_ProvisioningStatus{
		driver.GetProvisioningStatusRequest(domain.Id, "DOMAIN", value),
	}
	memberStatusRequests := []*server.MemberStatusRequest_MemberStatus{}

	for _, datacenter := range domain.Datacenters {
		if datacenter.ProvisioningStatus != "ACTIVE" {
			provisioningStatusRequests = append(provisioningStatusRequests,
				driver.GetProvisioningStatusRequest(datacenter.Id, "DATACENTER", value))
		}
	}

	for _, pool := range domain.Pools {
		// No pool representation in Akamai
		if pool.ProvisioningStatus != "ACTIVE" {
			provisioningStatusRequests = append(provisioningStatusRequests,
				driver.GetProvisioningStatusRequest(pool.Id, "POOL", value))
		}
		for _, monitor := range pool.Monitors {
			if monitor.ProvisioningStatus != "ACTIVE" {
				provisioningStatusRequests = append(provisioningStatusRequests,
					driver.GetProvisioningStatusRequest(monitor.Id, "MONITOR", value))
			}
		}
		for _, member := range pool.Members {
			if member.ProvisioningStatus != "ACTIVE" {
				provisioningStatusRequests = append(provisioningStatusRequests,
					driver.GetProvisioningStatusRequest(member.Id, "MEMBER", value))
				memberStatusRequests = append(memberStatusRequests,
					driver.GetMemberStatusRequest(member.Id, "ONLINE"))
			}
		}
	}

	driver.UpdateProvisioningStatus(s.rpc, provisioningStatusRequests)
	driver.UpdateMemberStatus(s.rpc, memberStatusRequests)
	return nil
}
