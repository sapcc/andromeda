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
	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/driver"
	"github.com/sapcc/andromeda/internal/rpc/server"
	"github.com/sapcc/andromeda/internal/rpcmodels"
	"go-micro.dev/v4/logger"
)

func (s *AkamaiAgent) UpdateDomainProvisioningStatus(domain *rpcmodels.Domain, value string) error {
	var provisioningStatusRequests []*server.ProvisioningStatusRequest_ProvisioningStatus

	provisioningStatusRequests = append(provisioningStatusRequests,
		driver.GetProvisioningStatusRequest(domain.Id, "DOMAIN", value))
	if value == "DELETED" {
		// Only delete domain, set related objects to active
		value = "ACTIVE"
	}

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
			}
		}
	}

	driver.UpdateProvisioningStatus(s.rpc, provisioningStatusRequests)
	return nil
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
		logger.Debug("Akamai Backend: pending configuration change")
	case "DENIED":
		if domain == nil {
			logger.Error("Akamai Backend: configuration change failed")
			return status.PropagationStatus, nil
		}

		logger.Errorf("Domain %s failed syncing: %s", domain.Id, status.Message)
		if err := s.UpdateDomainProvisioningStatus(domain, "ERROR"); err != nil {
			return "UNKNOWN", err
		}
	case "COMPLETE":
		if domain == nil {
			logger.Info("Akamai Backend: configuration change completed")
			return status.PropagationStatus, nil
		}
		logger.Infof("Domain %s has been propagated", domain.Id)
		provStatus := "ACTIVE"
		if domain.ProvisioningStatus == "PENDING_DELETE" {
			provStatus = "DELETED"
		} else if err := s.syncMemberStatus(domain); err != nil {
			logger.Warn(err)
		}

		if err := s.UpdateDomainProvisioningStatus(domain, provStatus); err != nil {
			return "UNKNOWN", err
		}

	}
	return status.PropagationStatus, nil
}
