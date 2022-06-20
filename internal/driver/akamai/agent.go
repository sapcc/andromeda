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
	"os"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_4"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/metadata"

	"github.com/sapcc/andromeda/internal/rpc/server"
	"github.com/sapcc/andromeda/internal/rpc/worker"
	"github.com/sapcc/andromeda/models"
)

var DOMAIN_MODE_MAP = map[string]string{
	models.DomainModeAVAILABILITY: "failover-only",
	models.DomainModeGEOGRAPHIC:   "static",
	models.DomainModeWEIGHTED:     "weighted",
	models.DomainModeROUNDROBIN:   "weighted",
}

var PROPERTY_TYPE_MAP = map[string]string{
	models.DomainModeAVAILABILITY: "failover",
	models.DomainModeGEOGRAPHIC:   "geographic",
	models.DomainModeWEIGHTED:     "weighted-round-robin",
	models.DomainModeROUNDROBIN:   "weighted-round-robin",
}

type AkamaiAgent struct {
	config *edgegrid.Config
	rpc    server.RPCServerService
}

// All methods of Sub will be executed when
// a message is received
type Sub struct {
	akamai *AkamaiAgent
}

func ExecuteAkamaiAgent() error {
	path, _ := os.LookupEnv("AKAMAI_EDGE_RC")
	config, _ := edgegrid.Init(path, "default")
	configgtm.Init(config)

	meta := map[string]string{
		"type":    "Akamai",
		"host":    config.Host,
		"version": "1.4",
	}
	service := micro.NewService(
		micro.Name("andromeda.agent.akamai"),
		micro.Metadata(meta),
		micro.RegisterTTL(time.Second*60),
		micro.RegisterInterval(time.Second*30),
	)
	service.Init()

	// Create F5 worker instance with Server RPC interface
	akamai := AkamaiAgent{
		&config,
		server.NewRPCServerService("andromeda.server", service.Client()),
	}

	// Register callbacks
	/*if err := worker.RegisterRPCWorkerHandler(service.Server(), &f5); err != nil {
		return err
	}*/
	if err := micro.RegisterSubscriber("andromeda.sync_all",
		service.Server(), &Sub{&akamai}); err != nil {
		logger.Error(err)
	}

	go akamai.periodicSync()
	go func() {
		if err := akamai.fullSync(); err != nil {
			logger.Error(err)
		}
	}()
	return service.Run()
}

func (s *AkamaiAgent) fullSync() error {
	var pageNumber int32 = 0
	for true {
		response, err := s.rpc.GetDomains(context.Background(), &server.SearchRequest{
			Provider:       "akamai",
			PageNumber:     pageNumber,
			ResultPerPage:  100,
			FullyPopulated: true,
			Pending:        true,
		})
		if err != nil {
			return err
		}

		res := response.GetResponse()
		for _, domain := range res {
			if err := s.SyncDomain(domain); err != nil {
				return err
			}
		}

		if len(res) < 100 {
			break
		} else {
			pageNumber++
		}
	}

	return nil
}

func (s *AkamaiAgent) periodicSync() {
	t := time.NewTicker(30 * time.Second)
	defer t.Stop()
	for {
		<-t.C // Activate periodically
		if err := s.fullSync(); err != nil {
			logger.Error(err)
		}
	}

}

func (s *AkamaiAgent) SyncAll(ctx context.Context, request *worker.SyncRequest) error {
	logger.Info("Sync invoked!")
	md, _ := metadata.FromContext(ctx)
	logger.Info("[pubsub.1] Received event %+v with metadata %+v\n", request, md)
	// do something with event
	return nil
}
