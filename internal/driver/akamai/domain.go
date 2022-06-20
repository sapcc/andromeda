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
	"hash/fnv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_4"
	"go-micro.dev/v4/logger"

	"github.com/sapcc/andromeda/internal/models"
	"github.com/sapcc/andromeda/internal/rpc/server"
	"github.com/sapcc/andromeda/internal/utils"
)

var MONITOR_LIVENESS_TYPE_MAP = map[models.Monitor_MonitorType]string{
	models.Monitor_HTTP:  "HTTP",
	models.Monitor_HTTPS: "HTTPS",
	models.Monitor_TCP:   "TCP",
}

func hash(s string) int {
	h := fnv.New32a()
	if _, err := h.Write([]byte(s)); err != nil {
		panic(err)
	}
	return int(h.Sum32())
}

func (s *AkamaiAgent) SyncDomain(domain *models.Domain) error {
	logger.Debugf("SyncDomain('%s')", domain.Id)

	// Populate Domain
	expectedDomain := configgtm.NewDomain(domain.GetFqdn(), DOMAIN_MODE_MAP[domain.GetMode()])
	pools := domain.GetPools()
	if len(pools) > 0 {
		var members []*models.Member
		var monitors []*models.Monitor

		// flatten Members and Monitors
		for _, pool := range pools {
			members = append(members, pool.GetMembers()...)
			monitors = append(monitors, pool.GetMonitors()...)
		}

		property, err := s.SyncProperty(domain, members, monitors)
		if err != nil {
			return err
		}
		expectedDomain.Properties = []*configgtm.Property{property}
	}

	logger.Infof("processing domain '%s' %s -> ACTIVE", domain.GetFqdn(), domain.ProvisioningStatus)
	akamaiDomain, err := configgtm.GetDomain(domain.GetFqdn())
	if err != nil && (domain.ProvisioningStatus == "PENDING_CREATE" || domain.ProvisioningStatus == "PENDING_UPDATE") {
		// Create
		if _, err := expectedDomain.Create(nil); err != nil {
			return err
		}
		if _, err := s.rpc.UpdateProvisioningStatus(context.Background(),
			&server.ProvisioningStatusRequest{ProvisioningStatus: []*server.ProvisioningStatusRequest_ProvisioningStatus{
				{
					Id:     domain.Id,
					Model:  server.ProvisioningStatusRequest_ProvisioningStatus_DOMAIN,
					Status: server.ProvisioningStatusRequest_ProvisioningStatus_ACTIVE,
				},
			}}); err != nil {
			return err
		}
		return nil
	} else {
		// Delete properties/resources/datacenter since domains are not deletable
		if domain.ProvisioningStatus == "PENDING_DELETE" {
			for _, property := range akamaiDomain.Properties {
				if _, err := property.Delete(domain.GetFqdn()); err != nil {
					return err
				}
			}
			for _, resource := range akamaiDomain.Resources {
				if _, err := resource.Delete(domain.GetFqdn()); err != nil {
					return err
				}
			}
			for _, datacenter := range akamaiDomain.Datacenters {
				if _, err := datacenter.Delete(domain.GetFqdn()); err != nil {
					return err
				}
			}
			if _, err := s.rpc.UpdateProvisioningStatus(context.Background(),
				&server.ProvisioningStatusRequest{ProvisioningStatus: []*server.ProvisioningStatusRequest_ProvisioningStatus{
					{
						Id:     domain.Id,
						Model:  server.ProvisioningStatusRequest_ProvisioningStatus_DOMAIN,
						Status: server.ProvisioningStatusRequest_ProvisioningStatus_ACTIVE,
					},
				}}); err != nil {
				return err
			}
		}
	}

	fieldsToCompare := []string{
		"Type",
		"Datacenters",
		"Datacenters.City",
		"Datacenters.Continent",
		"Datacenters.Country"}
	if !utils.DeepEqualFields(expectedDomain, akamaiDomain, fieldsToCompare) {
		// Update
		if _, err := expectedDomain.Update(nil); err != nil {
			return err
		}
	}

	return nil
}
