// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package akamai

import (
	"context"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/gtm"
	"github.com/apex/log"

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
	request := gtm.GetDomainRequest{DomainName: config.Global.AkamaiConfig.Domain}
	if _, err := s.gtm.GetDomain(context.Background(), request); err != nil {
		log.Warnf("Akamai Domain %s doesn't exist, creating...", config.Global.AkamaiConfig.Domain)
		domain := gtm.Domain{
			Name: config.Global.AkamaiConfig.Domain,
			Type: domainType,
		}
		createRequest := gtm.CreateDomainRequest{Domain: &domain}
		if _, err = s.gtm.CreateDomain(context.Background(), createRequest); err != nil {
			return err
		}
	}
	return nil
}

func (s *AkamaiAgent) FetchAndSyncDomains(domains []string, force bool) error {
	if s.executing {
		return nil
	}

	s.executing = true
	defer func() { s.executing = false }()

	log.Debugf("Running FetchAndSyncDomains(domains=%+v, force=%t)", domains, force)
	response, err := s.rpc.GetDomains(context.Background(), &server.SearchRequest{
		Provider:       "akamai",
		PageNumber:     0,
		ResultPerPage:  1,
		FullyPopulated: true,
		Pending:        domains == nil && !force,
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
			if datacenter.ProvisioningStatus != models.DatacenterProvisioningStatusACTIVE &&
				datacenter.ProvisioningStatus != models.DatacenterProvisioningStatusPENDINGDELETE {
				datacenters = append(datacenters, datacenter.Id)
			}
		}
	}
	if len(datacenters) > 0 {
		if err := s.FetchAndSyncDatacenters(datacenters, false); err != nil {
			log.Warnf("FetchAndSyncDomains(%+v) prior FetchAndSyncDatacenters(%+v): %v", domains, datacenters, err)
		}
	}

	for _, domain := range res {
		if err = func() error {
			var provRequests ProvRequests
			log.Infof("domainSync(%s) running...", domain.Id)
			defer s.gtmLock.Unlock()
			s.gtmLock.Lock()

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
					log.Error(err.Error())
				}
			}
			driver.UpdateProvisioningStatus(s.rpc, provRequests)
			log.Infof("FetchAndSyncDomains(%s) finished with '%s'", domain.Id, status)
			return nil
		}(); err != nil {
			return err
		}
	}
	return nil
}
