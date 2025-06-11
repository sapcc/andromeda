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
}

type CommonController struct {
	db  *sqlx.DB
	nc  *nats.Conn
	rpc *stormrpc.Client
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
		db:  db,
		nc:  nc,
		rpc: rpcClient,
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
		CidrBlocksController{cc, make(map[string]cidrBlocks)},
	}
	return &c
}

func PendingSync(client *stormrpc.Client) error {
	if client == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r, err := stormrpc.NewRequest("andromeda.sync", []string{})
	if err != nil {
		return err
	}

	resp := client.Do(ctx, r)
	return resp.Err
}
