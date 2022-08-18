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
	"go-micro.dev/v4/logger"

	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/models"
	"github.com/sapcc/andromeda/internal/rpc/server"
	"github.com/sapcc/andromeda/internal/utils"
)

func (s *AkamaiAgent) UpdateStatus(domain *models.Domain) error {
	property, err := s.gtm.GetProperty(context.Background(), domain.GetFqdn(), config.Global.AkamaiConfig.Domain)
	if err != nil {
		return err
	}

	logger.Debug(property)

	if domain.ProvisioningStatus != "ACTIVE" {
		s.UpdateStatusToActive("DOMAIN", domain.Id)
	}

	// flatten Members and Monitors
	for _, pool := range domain.Pools {
		// No pool representation in Akamai
		if pool.ProvisioningStatus != "ACTIVE" {
			s.UpdateStatusToActive("POOL", pool.Id)
		}
		for _, monitor := range pool.Monitors {
			s.UpdateMonitor(monitor, property.LivenessTests)
		}
		for _, member := range pool.Members {
			s.UpdateMember(member, property.TrafficTargets)
		}
	}

	return nil
}

func (s *AkamaiAgent) UpdateMonitor(monitor *models.Monitor, livenessTests []*gtm.LivenessTest) {
	// Todo: Update all kinds of statuses
	if monitor.ProvisioningStatus == "ACTIVE" {
		return
	}

	for _, livenessTest := range livenessTests {
		if livenessTest.Name != monitor.Id {
			continue
		}

		s.UpdateStatusToActive("MONITOR", monitor.Id)
		return
	}
}

func (s *AkamaiAgent) UpdateMember(member *models.Member, trafficTargets []*gtm.TrafficTarget) {
	// Todo: Update all kinds of statuses
	if member.ProvisioningStatus == "ACTIVE" {
		return
	}

	for _, trafficTarget := range trafficTargets {
		// Akamai bug, name is not persisted
		// if trafficTarget.Name != member.Id {
		if len(trafficTarget.Servers) != 1 || trafficTarget.Servers[0] != utils.InetNtoa(member.Address).String() {
			continue
		}

		s.UpdateStatusToActive("MEMBER", member.Id)
		return
	}
}

func (s *AkamaiAgent) UpdateStatusToActive(modelType string, id string) {
	model := server.ProvisioningStatusRequest_ProvisioningStatus_Model(
		server.ProvisioningStatusRequest_ProvisioningStatus_Model_value[modelType])
	logger.Infof("Status change of %s '%s' -> ACTIVE", modelType, id)

	if _, err := s.rpc.UpdateProvisioningStatus(context.Background(),
		&server.ProvisioningStatusRequest{ProvisioningStatus: []*server.ProvisioningStatusRequest_ProvisioningStatus{
			{
				Id:     id,
				Model:  model,
				Status: server.ProvisioningStatusRequest_ProvisioningStatus_ACTIVE,
			},
		}}); err != nil {
		logger.Warn(err)
	}
}
