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
	"github.com/jmoiron/sqlx"
	"go-micro.dev/v4"
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
}

func New(db *sqlx.DB) *Controller {
	sv := micro.NewService(
		micro.Name("andromeda-api"),
	)
	sv.Init()

	c := Controller{
		DomainController{db, sv},
		PoolController{db, sv},
		DatacenterController{db, sv},
		MemberController{db, sv},
		MonitorController{db, sv},
		ServiceController{db, sv},
		QuotaController{db},
		SyncController{sv},
	}
	return &c
}
