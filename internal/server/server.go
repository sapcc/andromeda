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

package server

import (
	"net/http"
	"os"
	"time"

	"github.com/dlmiddlecote/sqlstats"
	"github.com/getsentry/sentry-go"
	"github.com/go-openapi/loads"
	"github.com/iancoleman/strcase"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sapcc/go-bits/logg"
	"github.com/xo/dburl"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"

	_ "github.com/sapcc/andromeda/db/plugins"
	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/controller"
	"github.com/sapcc/andromeda/internal/policy"
	"github.com/sapcc/andromeda/internal/rpc/server"
	"github.com/sapcc/andromeda/internal/utils"
	"github.com/sapcc/andromeda/middlewares"
	"github.com/sapcc/andromeda/restapi"
	"github.com/sapcc/andromeda/restapi/operations"
	"github.com/sapcc/andromeda/restapi/operations/administrative"
	"github.com/sapcc/andromeda/restapi/operations/datacenters"
	"github.com/sapcc/andromeda/restapi/operations/domains"
	"github.com/sapcc/andromeda/restapi/operations/geographic_maps"
	"github.com/sapcc/andromeda/restapi/operations/members"
	"github.com/sapcc/andromeda/restapi/operations/monitors"
	"github.com/sapcc/andromeda/restapi/operations/pools"
)

func ExecuteServer(server *restapi.Server) error {
	logger.Info("Starting up andromeda-server")

	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		return err
	}

	// Use it globally
	utils.SwaggerSpec = swaggerSpec

	// Database
	u, err := dburl.Parse(config.Global.Database.Connection)
	if err != nil {
		return err
	}
	if u.Driver == "postgres" {
		u.Driver = "pgx"
	}
	db := sqlx.MustConnect(u.Driver, u.DSN)

	if config.Global.Default.SentryDSN != "" {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:              config.Global.Default.SentryDSN,
			AttachStacktrace: true,
			Release:          "TODO Version",
		}); err != nil {
			logg.Fatal("Sentry initialization failed: %v", err)
		}

		logg.Info("Sentry is enabled")
	}

	// Mapper function for SQL name mapping, snake_case table names
	db.MapperFunc(strcase.ToSnake)

	// Policy Engine
	policy.SetPolicyEngine(config.Global.ApiSettings.PolicyEngine)

	// Controller
	c := controller.New(db)

	// Initialize API
	api := operations.NewAndromedaAPI(swaggerSpec)

	// Logger
	api.Logger = logger.Infof

	// Prometheus Metrics
	if config.Global.Default.Prometheus {
		http.Handle("/metrics", promhttp.Handler())
		go PrometheusServer()

		// Create a new collector, the name will be used as a label on the metrics
		collector := sqlstats.NewStatsCollector("db_name", db)

		// Register it with Prometheus
		prometheus.MustRegister(collector)
	}

	// Domains
	api.DomainsGetDomainsHandler = domains.GetDomainsHandlerFunc(c.Domains.GetDomains)
	api.DomainsPostDomainsHandler = domains.PostDomainsHandlerFunc(c.Domains.PostDomains)
	api.DomainsGetDomainsDomainIDHandler = domains.GetDomainsDomainIDHandlerFunc(c.Domains.GetDomainsDomainID)
	api.DomainsPutDomainsDomainIDHandler = domains.PutDomainsDomainIDHandlerFunc(c.Domains.PutDomainsDomainID)
	api.DomainsDeleteDomainsDomainIDHandler = domains.DeleteDomainsDomainIDHandlerFunc(c.Domains.DeleteDomainsDomainID)

	// Pools
	api.PoolsGetPoolsHandler = pools.GetPoolsHandlerFunc(c.Pools.GetPools)
	api.PoolsPostPoolsHandler = pools.PostPoolsHandlerFunc(c.Pools.PostPools)
	api.PoolsGetPoolsPoolIDHandler = pools.GetPoolsPoolIDHandlerFunc(c.Pools.GetPoolsPoolID)
	api.PoolsPutPoolsPoolIDHandler = pools.PutPoolsPoolIDHandlerFunc(c.Pools.PutPoolsPoolID)
	api.PoolsDeletePoolsPoolIDHandler = pools.DeletePoolsPoolIDHandlerFunc(c.Pools.DeletePoolsPoolID)

	// Members
	api.MembersGetMembersHandler = members.GetMembersHandlerFunc(c.Members.GetMembers)
	api.MembersPostMembersHandler = members.PostMembersHandlerFunc(c.Members.PostMembers)
	api.MembersGetMembersMemberIDHandler = members.GetMembersMemberIDHandlerFunc(c.Members.GetMembersMemberID)
	api.MembersPutMembersMemberIDHandler = members.PutMembersMemberIDHandlerFunc(c.Members.PutMembersMemberID)
	api.MembersDeleteMembersMemberIDHandler = members.DeleteMembersMemberIDHandlerFunc(c.Members.DeleteMembersMemberID)

	// Datacenters
	api.DatacentersGetDatacentersHandler = datacenters.GetDatacentersHandlerFunc(c.Datacenters.GetDatacenters)
	api.DatacentersPostDatacentersHandler = datacenters.PostDatacentersHandlerFunc(c.Datacenters.PostDatacenters)
	api.DatacentersGetDatacentersDatacenterIDHandler = datacenters.GetDatacentersDatacenterIDHandlerFunc(c.Datacenters.GetDatacentersDatacenterID)
	api.DatacentersPutDatacentersDatacenterIDHandler = datacenters.PutDatacentersDatacenterIDHandlerFunc(c.Datacenters.PutDatacentersDatacenterID)
	api.DatacentersDeleteDatacentersDatacenterIDHandler = datacenters.DeleteDatacentersDatacenterIDHandlerFunc(c.Datacenters.DeleteDatacentersDatacenterID)

	// Geographic Maps
	api.GeographicMapsGetGeomapsHandler = geographic_maps.GetGeomapsHandlerFunc(c.GeoMaps.GetGeomaps)
	api.GeographicMapsPostGeomapsHandler = geographic_maps.PostGeomapsHandlerFunc(c.GeoMaps.PostGeomaps)
	api.GeographicMapsGetGeomapsGeomapIDHandler = geographic_maps.GetGeomapsGeomapIDHandlerFunc(c.GeoMaps.GetGeomapsGeoMapID)
	api.GeographicMapsPutGeomapsGeomapIDHandler = geographic_maps.PutGeomapsGeomapIDHandlerFunc(c.GeoMaps.PutGeomapsGeoMapID)
	api.GeographicMapsDeleteGeomapsGeomapIDHandler = geographic_maps.DeleteGeomapsGeomapIDHandlerFunc(c.GeoMaps.DeleteGeomapsGeoMapID)

	// Monitors
	api.MonitorsGetMonitorsHandler = monitors.GetMonitorsHandlerFunc(c.Monitors.GetMonitors)
	api.MonitorsPostMonitorsHandler = monitors.PostMonitorsHandlerFunc(c.Monitors.PostMonitors)
	api.MonitorsGetMonitorsMonitorIDHandler = monitors.GetMonitorsMonitorIDHandlerFunc(c.Monitors.GetMonitorsMonitorID)
	api.MonitorsPutMonitorsMonitorIDHandler = monitors.PutMonitorsMonitorIDHandlerFunc(c.Monitors.PutMonitorsMonitorID)
	api.MonitorsDeleteMonitorsMonitorIDHandler = monitors.DeleteMonitorsMonitorIDHandlerFunc(c.Monitors.DeleteMonitorsMonitorID)

	// Administrative
	api.AdministrativeGetServicesHandler = administrative.GetServicesHandlerFunc(c.Services.GetServices)
	api.AdministrativePostSyncHandler = administrative.PostSyncHandlerFunc(c.Sync.PostSync)

	// Quota Middleware
	if config.Global.Quota.Enabled {
		logger.Info("Initializing quota middleware")

		// Admin handler
		api.AdministrativeGetQuotasHandler = administrative.GetQuotasHandlerFunc(c.Quotas.GetQuotas)
		api.AdministrativeGetQuotasProjectIDHandler = administrative.GetQuotasProjectIDHandlerFunc(c.Quotas.GetQuotasProjectID)
		api.AdministrativeGetQuotasDefaultsHandler = administrative.GetQuotasDefaultsHandlerFunc(c.Quotas.GetQuotasDefaults)
		api.AdministrativePutQuotasProjectIDHandler = administrative.PutQuotasProjectIDHandlerFunc(c.Quotas.PutQuotasProjectID)
		api.AdministrativeDeleteQuotasProjectIDHandler = administrative.DeleteQuotasProjectIDHandlerFunc(c.Quotas.DeleteQuotasProjectID)

		qc := middlewares.NewQuotaController(db)
		api.AddMiddlewareFor("POST", "/datacenters", qc.QuotaHandler)
		api.AddMiddlewareFor("POST", "/domains", qc.QuotaHandler)
		api.AddMiddlewareFor("POST", "/monitors", qc.QuotaHandler)
		api.AddMiddlewareFor("POST", "/pools", qc.QuotaHandler)
		api.AddMiddlewareFor("POST", "/pools/{pool_id}/members", qc.QuotaHandler)
		api.AddMiddlewareFor("DELETE", "/datacenters/{datacenter_id}", qc.QuotaHandler)
		api.AddMiddlewareFor("DELETE", "/domains/{domain_id}", qc.QuotaHandler)
		api.AddMiddlewareFor("DELETE", "/monitors/{monitor_id}", qc.QuotaHandler)
		api.AddMiddlewareFor("DELETE", "/pools/{pool_id}", qc.QuotaHandler)
		api.AddMiddlewareFor("DELETE", "/pools/{pool_id}/members/{member_id}", qc.QuotaHandler)
	}
	server.SetAPI(api)

	//rpc worker
	go RPCServer(db)

	defer func() {
		if err := server.Shutdown(); err != nil {
			logger.Fatal(err)
		}
	}()
	if err := server.Serve(); err != nil {
		return err
	}

	return nil
}

func RPCServer(db *sqlx.DB) {
	host, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	service := micro.NewService(
		micro.Name("andromeda.server"),
		micro.Metadata(map[string]string{
			"type": "andromeda",
			"host": host,
		}),
		micro.RegisterTTL(time.Second*60),
		micro.RegisterInterval(time.Second*30),
		utils.ConfigureTransport(),
	)
	service.Init()

	// Update handler for provisioning and status updates
	if err := server.RegisterRPCServerHandler(service.Server(), &server.RPCHandler{DB: db}); err != nil {
		panic(err)
	}

	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}

func PrometheusServer() {
	logg.Info("Serving prometheus metrics at http://%s/metrics", config.Global.Default.PrometheusListen)
	if err := http.ListenAndServe(config.Global.Default.PrometheusListen, nil); err != nil {
		logg.Fatal(err.Error())
	}
}
