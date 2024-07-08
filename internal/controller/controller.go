/*
 *   Copyright 2020 SAP SE
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
}

type CommonController struct {
	db  *sqlx.DB
	nc  *nats.Conn
	rpc *stormrpc.Client
}

func New(db *sqlx.DB) *Controller {
	nc, err := nats.Connect(config.Global.Default.TransportURL)
	if err != nil {
		log.Fatal(err.Error())
	}

	rpcClient, err := stormrpc.NewClient("", stormrpc.WithNatsConn(nc))
	if err != nil {
		log.Fatal(err.Error())
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
