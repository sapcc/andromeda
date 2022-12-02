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

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/configgtm"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"

	"github.com/sapcc/andromeda/internal/config"
	_ "github.com/sapcc/andromeda/internal/plugins"
	"github.com/sapcc/andromeda/internal/rpc/server"
	"github.com/sapcc/andromeda/internal/utils"
)

func ExecuteAkamaiStatusAgent() error {
	meta := map[string]string{
		"type":    "Akamai",
		"version": "2.0",
	}
	service := micro.NewService(
		micro.Name("andromeda.agent.akamai-status"),
		micro.Metadata(meta),
		micro.RegisterTTL(time.Second*60),
		micro.RegisterInterval(time.Second*30),
		utils.ConfigureTransport(),
	)
	service.Init()

	// Create F5 worker instance with Server RPC interface
	s, domainType := NewAkamaiSession()
	akamai := AkamaiAgent{
		gtm.Client(*s),
		domainType,
		server.NewRPCServerService("andromeda.server", service.Client()),
	}

	go akamai.periodicStatusSync()
	return service.Run()
}

func (s *AkamaiAgent) periodicStatusSync() {
	t := time.NewTicker(30 * time.Second)
	defer t.Stop()
	for {
		<-t.C // Activate periodically
		if err := s.syncStatus(); err != nil {
			logger.Error(err)
		}
	}
}

func (s *AkamaiAgent) syncStatus() error {
	// Akamai backend can only process one change at a time
	response, err := s.rpc.GetDomains(context.Background(), &server.SearchRequest{
		Provider:       "akamai",
		PageNumber:     0,
		ResultPerPage:  1,
		FullyPopulated: false,
		Pending:        true,
	})
	if err != nil {
		return err
	}

	// Check for running domain's propagation state
	status, err := s.gtm.GetDomainStatus(context.Background(), config.Global.AkamaiConfig.Domain)
	if err != nil {
		return err
	}

	if status.PropagationStatus == "PENDING" {
		return nil
	}

	for _, domain := range response.GetResponse() {
		// Tracks the status of the domain's propagation state. Either PENDING, COMPLETE, or DENIED.
		// A DENIED value indicates that the domain configuration is invalid,
		// and doesn't propagate until the validation errors are resolved.
		switch status.PropagationStatus {
		case "PENDING":
			logger.Debug("Backend has pending configuration change")
		case "DENIED":
			logger.Errorf("Domain %s failed syncing: %s", domain.Id, status.Message)
			if err := s.UpdateDomainProvisioningStatus(domain, "ERROR"); err != nil {
				return err
			}
		case "COMPLETE":
			logger.Infof("Domain %s is in sync", domain.Id)
			if err := s.UpdateDomainProvisioningStatus(domain, "ACTIVE"); err != nil {
				return err
			}
			// Ready for new Changes
		}
	}
	return nil
}
