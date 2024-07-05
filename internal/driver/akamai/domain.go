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
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/gtm"
	"go-micro.dev/v4/logger"

	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/driver"
	"github.com/sapcc/andromeda/internal/rpc/server"
	"github.com/sapcc/andromeda/internal/rpcmodels"
	"github.com/sapcc/andromeda/models"
)

var MONITOR_LIVENESS_TYPE_MAP = map[rpcmodels.Monitor_MonitorType]string{
	rpcmodels.Monitor_HTTP:  "HTTP",
	rpcmodels.Monitor_HTTPS: "HTTPS",
	rpcmodels.Monitor_TCP:   "TCP",
}

func (s *AkamaiAgent) EnsureDomain(domainType string) error {
	if _, err := s.gtm.GetDomain(context.Background(), config.Global.AkamaiConfig.Domain); err != nil {
		logger.Warnf("Akamai Domain %s doesn't exist, creating...", config.Global.AkamaiConfig.Domain)
		domain := gtm.Domain{
			Name: config.Global.AkamaiConfig.Domain,
			Type: domainType,
		}
		if _, err := s.gtm.CreateDomain(context.Background(), &domain, nil); err != nil {
			return err
		}
	}
	return nil
}

func (s *AkamaiAgent) FetchAndSyncDomains(domains []string) error {
	if s.executing {
		return nil
	}

	s.executing = true
	defer func() { s.executing = false }()

	logger.Debugf("Running FetchAndSyncDomains(domains=%+v)", domains)
	response, err := s.rpc.GetDomains(context.Background(), &server.SearchRequest{
		Provider:       "akamai",
		PageNumber:     0,
		ResultPerPage:  1,
		FullyPopulated: true,
		Pending:        domains == nil,
		Ids:            domains,
	})
	if err != nil {
		return err
	}

	res := response.GetResponse()
	if len(res) == 0 {
		return nil
	}

	// TODO: support multiple trafficManagementDomains due to limit of 100 properties
	trafficManagementDomain := config.Global.AkamaiConfig.Domain

	// Sync all required datacenters first
	var datacenters []string
	for _, domain := range res {
		for _, datacenter := range domain.Datacenters {
			if datacenter.ProvisioningStatus != models.DatacenterProvisioningStatusACTIVE {
				datacenters = append(datacenters, datacenter.Id)
			}
		}
	}
	if len(datacenters) > 0 {
		if err := s.FetchAndSyncDatacenters(datacenters, false); err != nil {
			return err
		}
	}

	for _, domain := range res {
		var provRequests ProvRequests
		logger.Infof("domainSync(%s) running...", domain.Id)
		if err := s.gtmLock.Lock(trafficManagementDomain); err != nil {
			return err
		}

		if domain.ProvisioningStatus == models.DomainProvisioningStatusPENDINGDELETE {
			// Run Delete
			if err := s.DeleteProperty(domain, trafficManagementDomain); err != nil {
				return err
			}
			provRequests = s.CascadeUpdateDomainProvisioningStatus(domain, "DELETED")
		} else {
			// Run Sync
			if provRequests, err = s.SyncProperty(domain, trafficManagementDomain); err != nil {
				return err
			}
		}

		// Wait for status propagation
		var status string
		for ok := true; ok; ok = status == "PENDING" {
			time.Sleep(5 * time.Second)
			status, err = s.syncProvisioningStatus(domain)
			if err != nil {
				logger.Error(err)
			}
		}
		driver.UpdateProvisioningStatus(s.rpc, provRequests)

		logger.Infof("FetchAndSyncDomains(%s) finished with '%s'", domain.Id, status)
		if err := s.gtmLock.Unlock(trafficManagementDomain); err != nil {
			return err
		}
	}
	return nil
}
