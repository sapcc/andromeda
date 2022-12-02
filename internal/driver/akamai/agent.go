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
	"go-micro.dev/v4/metadata"

	"github.com/sapcc/andromeda/internal/config"
	_ "github.com/sapcc/andromeda/internal/plugins"
	"github.com/sapcc/andromeda/internal/rpc/server"
	"github.com/sapcc/andromeda/internal/rpc/worker"
	"github.com/sapcc/andromeda/internal/rpcmodels"
	"github.com/sapcc/andromeda/internal/utils"
	"github.com/sapcc/andromeda/models"
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

	// Create F5 worker instance with Server RPC interface
	s, domainType := NewAkamaiSession()
	akamai := AkamaiAgent{
		gtm.Client(*s),
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
	status, err := s.gtm.GetDomainStatus(context.Background(), config.Global.AkamaiConfig.Domain)
	if err != nil {
		return err
	}
	if status.PropagationStatus == "PENDING" {
		return nil
	}

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
