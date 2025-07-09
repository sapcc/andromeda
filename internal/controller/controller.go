// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"context"
	"time"

	"github.com/actatum/stormrpc"
	"github.com/apex/log"
	"github.com/jmoiron/sqlx"
	"github.com/nats-io/nats.go"

	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/rpc/agent"
	"github.com/sapcc/andromeda/internal/rpcmodels"
)

type Controller struct {
	Domains     DomainController
	Pools       PoolController
	Datacenters DatacenterController
	Members     MemberController
	Monitors    MonitorController
	Services    ServiceController
	Quotas      QuotaController
	Sync        SyncController
	GeoMaps     GeoMapController
	CidrBlocks  CidrBlocksController
	Metrics     AkamaiMetricsController
}

type CommonController struct {
	db    *sqlx.DB
	nc    *nats.Conn
	rpc   *stormrpc.Client
	agent agent.RPCAgentClient
}

func New(db *sqlx.DB) *Controller {
	nc, err := nats.Connect(config.Global.Default.TransportURL)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to NATS")
	}

	rpcClient, err := stormrpc.NewClient("", stormrpc.WithNatsConn(nc))
	if err != nil {
		log.WithError(err).Fatal("Failed to create RPC client")
	}

	cc := CommonController{
		db:    db,
		nc:    nc,
		rpc:   rpcClient,
		agent: agent.NewRPCAgentClient(rpcClient),
	}

	c := Controller{
		DomainController{cc},
		PoolController{cc},
		DatacenterController{cc},
		MemberController{cc},
		MonitorController{cc},
		ServiceController{cc},
		QuotaController{cc},
		SyncController{cc},
		GeoMapController{cc},
		NewCidrBlocksController(cc),
		NewAkamaiMetricsController(cc),
	}
	return &c
}

func (cc CommonController) PendingSync() error {
	// Skip sync if agent is not initialized (e.g., in tests)
	if cc.agent == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := cc.agent.Sync(ctx, &rpcmodels.SyncRequest{})
	return err
}
