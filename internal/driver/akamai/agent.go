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
	"net/http"
	"os"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/configgtm"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/sapcc/andromeda/internal/config"
	_ "github.com/sapcc/andromeda/internal/plugins"
	"github.com/sapcc/andromeda/internal/rpc/server"
	"github.com/sapcc/andromeda/internal/rpc/worker"
	"github.com/sapcc/andromeda/internal/rpcmodels"
	"github.com/sapcc/andromeda/internal/utils"
	"github.com/sapcc/andromeda/models"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/metadata"
)

var PROPERTY_TYPE_MAP = map[string]string{
	models.DomainModeAVAILABILITY: "failover",
	models.DomainModeGEOGRAPHIC:   "geographic",
	models.DomainModeWEIGHTED:     "weighted-round-robin",
	models.DomainModeROUNDROBIN:   "weighted-round-robin",
}

type AkamaiAgent struct {
	gtm        gtm.GTM
	domainType string
	rpc        server.RPCServerService
}

// All methods of Sub will be executed when
// a message is received
type Sub struct {
	akamai *AkamaiAgent
}

// Method can be of any name
func (s *Sub) Process(ctx context.Context, request *worker.SyncRequest) error {
	md, _ := metadata.FromContext(ctx)
	logger.Infof("[pubsub.1] Received event %+v with metadata %+v\n", request, md)
	// do something with event
	return nil
}

func ExecuteAkamaiAgent() error {
	meta := map[string]string{
		"type":    "Akamai",
		"version": "2.0",
	}
	service := micro.NewService(
		micro.Name("andromeda.agent.akamai"),
		micro.Metadata(meta),
		micro.RegisterTTL(time.Second*60),
		micro.RegisterInterval(time.Second*30),
		utils.ConfigureTransport(),
	)
	service.Init()

	option := edgegrid.WithEnv(true)
	if env := os.Getenv("AKAMAI_EDGE_RC"); env != "" {
		option = edgegrid.WithFile(env)
	} else if config.Global.AkamaiConfig.EdgeRC != "" {
		option = edgegrid.WithFile(config.Global.AkamaiConfig.EdgeRC)
	}

	edgerc := edgegrid.Must(edgegrid.New(option))
	s := session.Must(session.New(
		session.WithSigner(edgerc),
	))

	var identity struct {
		AccountID string `json:"accountId"`
		Active    bool   `json:"active"`
		Contracts []struct {
			ContractID  string   `json:"contractId"`
			Features    []string `json:"features"`
			Permissions []string `json:"permissions"`
		} `json:"contracts"`
		Email string `json:"email"`
	}

	req, _ := http.NewRequest(http.MethodGet, "/config-gtm/v1/identity", nil)
	if _, err := s.Exec(req, &identity); err != nil {
		panic(err)
	}

	if len(identity.Contracts) != 1 && config.Global.AkamaiConfig.ContractId == "" {
		logger.Fatalf("More than one contract detected, specificy contract_id.")
	}

	var domainType string
	for _, contract := range identity.Contracts {
		if config.Global.AkamaiConfig.ContractId != "" && contract.ContractID != config.Global.AkamaiConfig.ContractId {
			continue
		}

		domainType = DetectTypeFromFeatures(contract.Features)
		logger.Infof("Detected Akamai Contract '%s' with best features enabling '%s' domain type.",
			contract.ContractID, domainType)
		break
	}

	// Create F5 worker instance with Server RPC interface
	akamai := AkamaiAgent{
		gtm.Client(s),
		domainType,
		server.NewRPCServerService("andromeda.server", service.Client()),
	}

	if err := akamai.EnsureDomain(); err != nil {
		panic(err)
	}

	// Register callbacks
	if err := micro.RegisterSubscriber("andromeda.force_sync",
		service.Server(), &Sub{&akamai}); err != nil {
		logger.Error(err)
	}

	_listServices, err := service.Options().Registry.ListServices()
	if err != nil {
		logger.Error(err)
	}
	logger.Infof("%+v", _listServices)

	go akamai.periodicSync()
	go func() {
		if err := akamai.pendingSync(); err != nil {
			logger.Error(err)
		}
	}()
	return service.Run()
}

func (s *AkamaiAgent) pendingSync() error {
	// Akamai backend can only process one change at a time
	response, err := s.rpc.GetDomains(context.Background(), &server.SearchRequest{
		Provider:       "akamai",
		PageNumber:     0,
		ResultPerPage:  1,
		FullyPopulated: true,
		Pending:        true,
	})
	if err != nil {
		return err
	}

	res := response.GetResponse()
	for _, domain := range res {
		// Run Sync
		if err := s.SyncProperty(domain); err != nil {
			return err
		}

		// Check for running domain's propagation state
		status, err := s.gtm.GetDomainStatus(context.Background(), config.Global.AkamaiConfig.Domain)
		if err != nil {
			return err
		}

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

func (s *AkamaiAgent) periodicSync() {
	t := time.NewTicker(30 * time.Second)
	defer t.Stop()
	for {
		<-t.C // Activate periodically
		if err := s.pendingSync(); err != nil {
			logger.Error(err)
		}
	}

}

func (s *AkamaiAgent) ForceSync(ctx context.Context, request *worker.SyncRequest) error {
	md, _ := metadata.FromContext(ctx)
	var searchIds []string
	if domainId, ok := md.Get("domain"); ok {
		searchIds = []string{domainId}
		logger.Infof("Sync invoked for domain %s", md["domain"])
	} else {
		logger.Infof("Sync invoked for all domains")
	}

	var pageNumber int32 = 0
	var domains []*rpcmodels.Domain
	for {
		response, err := s.rpc.GetDomains(context.Background(), &server.SearchRequest{
			Provider:       "akamai",
			PageNumber:     pageNumber,
			ResultPerPage:  10,
			FullyPopulated: true,
			Pending:        false,
			Ids:            searchIds,
		})
		if err != nil {
			return err
		}

		domains = append(domains, response.Response...)

		if len(response.Response) < 100 {
			break
		} else {
			pageNumber++
		}
	}

	go func() {
		for i, domain := range domains {
			logger.Infof("Syncing Domain %d/%d: %s", i, len(domains), domain.Id)
			if err := s.SyncProperty(domain); err != nil {
				logger.Error(err)
			}
		}
	}()
	return nil
}
