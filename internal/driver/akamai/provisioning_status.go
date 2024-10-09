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

	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/driver"
	"github.com/sapcc/andromeda/internal/rpc/server"
	"github.com/sapcc/andromeda/internal/rpcmodels"
	"github.com/sapcc/andromeda/models"

	"github.com/apex/log"
)

type ProvRequests []*server.ProvisioningStatusRequest_ProvisioningStatus

func (s *AkamaiAgent) CascadeUpdateDomainProvisioningStatus(domain *rpcmodels.Domain, value string) ProvRequests {
	var provisioningStatusRequests ProvRequests

	provisioningStatusRequests = append(provisioningStatusRequests,
		driver.GetProvisioningStatusRequest(domain.Id, "DOMAIN", value))
	if value == "DELETED" {
		// Only delete domain, set related objects to active
		value = models.DomainProvisioningStatusACTIVE
	}

	for _, datacenter := range domain.Datacenters {
		if datacenter.ProvisioningStatus != models.DatacenterProvisioningStatusACTIVE {
			provisioningStatusRequests = append(provisioningStatusRequests,
				driver.GetProvisioningStatusRequest(datacenter.Id, "DATACENTER", value))
		}
	}

	for _, pool := range domain.Pools {
		// No pool representation in Akamai
		if pool.ProvisioningStatus != models.PoolProvisioningStatusACTIVE {
			provisioningStatusRequests = append(provisioningStatusRequests,
				driver.GetProvisioningStatusRequest(pool.Id, "POOL", value))
		}
		for _, monitor := range pool.Monitors {
			if monitor.ProvisioningStatus != models.MonitorProvisioningStatusACTIVE {
				provisioningStatusRequests = append(provisioningStatusRequests,
					driver.GetProvisioningStatusRequest(monitor.Id, "MONITOR", value))
			}
		}
		for _, member := range pool.Members {
			if member.ProvisioningStatus != models.MemberProvisioningStatusACTIVE {
				provisioningStatusRequests = append(provisioningStatusRequests,
					driver.GetProvisioningStatusRequest(member.Id, "MEMBER", value))
			}
		}
	}

	return provisioningStatusRequests
}

func (s *AkamaiAgent) syncProvisioningStatus(domain *rpcmodels.Domain) (string, error) {
	// Check for running domain's propagation state
	status, err := s.gtm.GetDomainStatus(context.Background(), config.Global.AkamaiConfig.Domain)
	if err != nil {
		return "UNKNOWN", err
	}

	// Tracks the status of the domain's propagation state. Either PENDING, COMPLETE, or DENIED.
	// A DENIED value indicates that the domain configuration is invalid,
	// and doesn't propagate until the validation errors are resolved.
	switch status.PropagationStatus {
	case "PENDING":
		log.Debug("Akamai Backend: pending configuration change")
	case "DENIED":
		if domain == nil {
			log.Error("Akamai Backend: configuration change failed")
			return status.PropagationStatus, nil
		}

		return status.PropagationStatus, fmt.Errorf("Domain %s failed syncing: %s", domain.Id, status.Message)
	case "COMPLETE":
		if domain == nil {
			log.Info("Akamai Backend: configuration change completed")
			return status.PropagationStatus, nil
		}
		log.Infof("Domain %s has been propagated", domain.Id)
		if err = s.syncMemberStatus(domain); err != nil {
			log.Warnf("syncProvisioningStatus: %s", err.Error())
		}
	}
	return status.PropagationStatus, nil
}
