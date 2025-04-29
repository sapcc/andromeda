// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"os"
	"testing"

	"github.com/go-openapi/loads"
	"github.com/iancoleman/strcase"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
	"github.com/xo/dburl"

	migration "github.com/sapcc/andromeda/db"
	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/policy"
	"github.com/sapcc/andromeda/internal/utils"
	"github.com/sapcc/andromeda/restapi"
)

type SuiteTest struct {
	suite.Suite
	dbUrl string
	db    *sqlx.DB
	c     *Controller
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(SuiteTest))
}

// Setup db value
func (t *SuiteTest) SetupSuite() {
	var err error

	t.dbUrl = os.Getenv("DB_URL")
	if t.dbUrl == "" {
		t.dbUrl = "mysql://localhost/test_suite_controller?sql_mode=%27ANSI_QUOTES%27"
	}

	u, _ := dburl.Parse(t.dbUrl)
	if u.Driver == "postgres" {
		u.Driver = "pgx"
	}
	t.db, err = sqlx.Connect(u.Driver, u.DSN)
	if err != nil {
		t.FailNow("Failed connecting to Database", err)
	}

	t.db.MapperFunc(strcase.ToSnake)
	policy.SetPolicyEngine("noop")

	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		t.FailNow("Failed loading swagger spec", err)
	}

	// Use it globally
	utils.SwaggerSpec = swaggerSpec
	config.Global.ApiSettings.PaginationMaxLimit = 1000

	cc := CommonController{
		db: t.db,
	}
	// initialize controller
	t.c = &Controller{
		Domains:     DomainController{cc},
		Pools:       PoolController{cc},
		Datacenters: DatacenterController{cc},
		Members:     MemberController{cc},
		Monitors:    MonitorController{cc},
		Services:    ServiceController{cc},
		Quotas:      QuotaController{cc},
		Sync:        SyncController{cc},
	}

	if err := migration.Migrate(t.dbUrl); err != nil {
		t.FailNow("Failed migration", err)
	}
}

// Run After All Test Done
func (t *SuiteTest) TearDownSuite() {
	if err := migration.Rollback(t.dbUrl); err != nil {
		t.FailNow("Failed rollback", err)
	}
}
