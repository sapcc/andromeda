// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package akamai

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/actatum/stormrpc"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/gtm"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	"github.com/apex/log"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/nats-io/nats.go"

	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/rpc"
	"github.com/sapcc/andromeda/internal/rpc/server"
	"github.com/sapcc/andromeda/internal/utils"
	"github.com/sapcc/andromeda/models"
)

var PROPERTY_TYPE_MAP = map[string]string{
	models.DomainModeAVAILABILITY: "weighted-round-robin",
	models.DomainModeGEOGRAPHIC:   "geographic",
	models.DomainModeWEIGHTED:     "weighted-round-robin",
	models.DomainModeROUNDROBIN:   "weighted-round-robin",
}

type AkamaiAgent struct {
	session           *session.Session
	gtm               gtm.GTM
	gtmLock           sync.Mutex
	domainType        string
	rpc               server.RPCServerClient
	workerTicker      *time.Ticker
	lastSync          time.Time
	lastMemberStatus  time.Time
	forceSync         chan []string
	executing         bool
	datacenterIdCache *lru.Cache[string, int]
}

var akamaiAgent *AkamaiAgent

func Sync(ctx context.Context, req stormrpc.Request) stormrpc.Response {
	var domainIDs []string
	if err := req.Decode(&domainIDs); err != nil {
		return stormrpc.NewErrorResponse(req.Reply, err)
	}
	log.WithField("domainIDs", domainIDs).Info("[Sync] Syncing domains")

	akamaiAgent.forceSync <- domainIDs
	resp, err := stormrpc.NewResponse(req.Reply, nil)
	if err != nil {
		return stormrpc.NewErrorResponse(req.Reply, err)
	}

	return resp
}

func GetCidrs(ctx context.Context, req stormrpc.Request) stormrpc.Response {
	cidrBlocksReq, _ := http.NewRequest(http.MethodGet, "/firewall-rules-manager/v1/cidr-blocks", nil)
	var cidrBlocks []AkamaiCIDRBlock
	if _, err := (*akamaiAgent.session).Exec(cidrBlocksReq, &cidrBlocks); err != nil {
		panic(err)
	}

	resp, err := stormrpc.NewResponse(req.Reply, cidrBlocks)
	if err != nil {
		return stormrpc.NewErrorResponse(req.Reply, err)
	}

	return resp
}

func ExecuteAkamaiAgent() error {
	nc, err := nats.Connect(config.Global.Default.TransportURL)
	if err != nil {
		return err
	}
	client, err := stormrpc.NewClient("", stormrpc.WithNatsConn(nc))
	if err != nil {
		return err
	}

	// Create F5 worker instance with Server RPC interface
	s, domainType := NewAkamaiSession(&config.Global.AkamaiConfig)

	// Figure out minimal ticker interval
	interval := time.Duration(config.Global.AkamaiConfig.SyncInterval) + 1
	if time.Duration(config.Global.AkamaiConfig.MemberStatusInterval) < interval {
		interval = time.Duration(config.Global.AkamaiConfig.MemberStatusInterval) + 1
	}

	cache, _ := lru.New[string, int](64)

	akamaiAgent = &AkamaiAgent{
		s,
		gtm.Client(*s),
		sync.Mutex{},
		domainType,
		server.NewRPCServerClient(client),
		time.NewTicker(interval * time.Second),
		time.Unix(0, 0),
		time.Unix(0, 0),
		make(chan []string),
		false,
		cache,
	}

	if err := akamaiAgent.EnsureDomain(domainType); err != nil {
		return err
	}

	srv := rpc.NewServer("andromeda-akamai-agent", stormrpc.WithNatsConn(nc))
	srv.Handle("andromeda.sync", Sync)
	srv.Handle("andromeda.get_cidrs.akamai", GetCidrs)

	go func() {
		_ = srv.Run()
	}()
	go akamaiAgent.WorkerThread()
	if config.Global.Default.Prometheus {
		go utils.PrometheusListen()
	}
	log.WithField("subjects", srv.Subjects()).Info("Subscribed")

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	// full sync immediately

	akamaiAgent.forceSync <- nil
	<-done
	log.Info("Shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return srv.Shutdown(ctx)
}

func (s *AkamaiAgent) WorkerThread() {
	syncInterval := time.Duration(config.Global.AkamaiConfig.SyncInterval) * time.Second
	memberStatusInterval := time.Duration(config.Global.AkamaiConfig.MemberStatusInterval) * time.Second

	for {
		select {
		case domains := <-s.forceSync:
			log.Debug("Running force sync")
			if err := s.FetchAndSyncDatacenters(nil, true); err != nil {
				log.Error(err.Error())
			}

			if err := s.FetchAndSyncGeomaps(nil, true); err != nil {
				log.Error(err.Error())
			}

			if err := s.FetchAndSyncDomains(domains, true); err != nil {
				log.Error(err.Error())
			}
		case <-s.workerTicker.C: // Activate periodically
			if time.Since(s.lastSync) > syncInterval {
				log.Debug("Running periodic sync")
				if err := s.FetchAndSyncDatacenters(nil, false); err != nil {
					log.Error(err.Error())
				}

				if err := s.FetchAndSyncGeomaps(nil, false); err != nil {
					log.Error(err.Error())
				}

				if err := s.FetchAndSyncDomains(nil, false); err != nil {
					log.Error(err.Error())
				}

				s.lastSync = time.Now()
			}
			if time.Since(s.lastMemberStatus) > memberStatusInterval {
				if err := s.memberStatusSync(); err != nil {
					log.Error(err.Error())
				}
				s.lastMemberStatus = time.Now()
			}
		}
	}
}

func (s *AkamaiAgent) memberStatusSync() error {
	log.Debugf("Running member status sync")
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
