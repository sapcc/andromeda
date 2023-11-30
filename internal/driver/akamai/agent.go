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

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/gtm"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/session"
	"github.com/go-micro/plugins/v4/sync/memory"
	lru "github.com/hashicorp/golang-lru/v2"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/metadata"
	"go-micro.dev/v4/sync"

	"github.com/sapcc/andromeda/internal/config"
	_ "github.com/sapcc/andromeda/internal/plugins"
	"github.com/sapcc/andromeda/internal/rpc/server"
	"github.com/sapcc/andromeda/internal/rpc/worker"
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
	session           *session.Session
	gtm               gtm.GTM
	gtmLock           sync.Sync
	domainType        string
	rpc               server.RPCServerService
	workerTicker      *time.Ticker
	lastSync          time.Time
	lastMemberStatus  time.Time
	forceSync         chan []string
	executing         bool
	datacenterIdCache *lru.Cache[string, int]
}

type Sync struct {
	akamai *AkamaiAgent
}

// Process Method can be of any name
func (s *Sync) Process(ctx context.Context, request *worker.SyncRequest) error {
	md, _ := metadata.FromContext(ctx)
	logger.Infof("[Sync] Received sync request %+v with metadata %+v", request.DomainIds, md)
	s.akamai.forceSync <- request.DomainIds
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
	syncer := &Sync{}
	service.Init(
		micro.AfterStart(func() error {
			// Create F5 worker instance with Server RPC interface
			s, domainType := NewAkamaiSession(&config.Global.AkamaiConfig)

			// Figure out minimal ticker interval
			interval := time.Duration(config.Global.AkamaiConfig.SyncInterval) + 1
			if time.Duration(config.Global.AkamaiConfig.MemberStatusInterval) < interval {
				interval = time.Duration(config.Global.AkamaiConfig.MemberStatusInterval) + 1
			}

			cache, _ := lru.New[string, int](64)

			akamai := AkamaiAgent{
				s,
				gtm.Client(*s),
				memory.NewSync(),
				domainType,
				server.NewRPCServerService("andromeda.server", service.Client()),
				time.NewTicker(interval * time.Second),
				time.Unix(0, 0),
				time.Unix(0, 0),
				make(chan []string),
				false,
				cache,
			}
			syncer.akamai = &akamai

			if err := akamai.EnsureDomain(domainType); err != nil {
				return err
			}

			go akamai.WorkerThread()
			// full sync immediately
			akamai.forceSync <- nil
			return nil
		}),
	)

	// Sync
	if err := micro.RegisterSubscriber("andromeda.sync",
		service.Server(), new(Sync)); err != nil {
		logger.Error(err)
	}

	return service.Run()
}

func (s *AkamaiAgent) WorkerThread() {
	syncInterval := time.Duration(config.Global.AkamaiConfig.SyncInterval) * time.Second
	memberStatusInterval := time.Duration(config.Global.AkamaiConfig.MemberStatusInterval) * time.Second

	for {
		select {
		case domains := <-s.forceSync:
			if err := s.FetchAndSyncDatacenters(nil, true); err != nil {
				logger.Error(err)
			}

			if err := s.FetchAndSyncGeomaps(nil, true); err != nil {
				logger.Error(err)
			}

			if err := s.FetchAndSyncDomains(domains); err != nil {
				logger.Error(err)
			}
		case <-s.workerTicker.C: // Activate periodically
			if time.Since(s.lastSync) > syncInterval {
				if err := s.FetchAndSyncDatacenters(nil, false); err != nil {
					logger.Error(err)
				}

				if err := s.FetchAndSyncGeomaps(nil, false); err != nil {
					logger.Error(err)
				}

				if err := s.FetchAndSyncDomains(nil); err != nil {
					logger.Error(err)
				}

				s.lastSync = time.Now()
			}
			if time.Since(s.lastMemberStatus) > memberStatusInterval {
				if err := s.memberStatusSync(); err != nil {
					logger.Error(err)
				}
				s.lastMemberStatus = time.Now()
			}
		}
	}

}

func (s *AkamaiAgent) ForceSync(ctx context.Context, request *worker.SyncRequest) error {
	logger.Infof("Got Sync request %v", request)
	md, _ := metadata.FromContext(ctx)
	if domainId, ok := md.Get("domain"); ok {
		s.forceSync <- []string{domainId}
	}
	return nil
}

func (s *AkamaiAgent) memberStatusSync() error {
	logger.Debugf("Running member status sync")
	response, err := s.rpc.GetDomains(context.Background(), &server.SearchRequest{
		Provider:       "akamai",
		PageNumber:     0,
		ResultPerPage:  1,
		FullyPopulated: true,
		Pending:        false,
	})
	if err != nil {
		return err
	}

	for _, domain := range response.GetResponse() {
		if err := s.syncMemberStatus(domain); err != nil {
			return err
		}
	}
	return nil
}
