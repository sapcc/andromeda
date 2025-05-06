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
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/actatum/stormrpc"
	"github.com/apex/log"
	"github.com/dlmiddlecote/sqlstats"
	"github.com/go-openapi/loads"
	"github.com/iancoleman/strcase"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/xo/dburl"

	_ "github.com/sapcc/andromeda/db/plugins"
	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/controller"
	"github.com/sapcc/andromeda/internal/policy"
	"github.com/sapcc/andromeda/internal/rpc"
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
	metrics_ops "github.com/sapcc/andromeda/restapi/operations/metrics"
	"github.com/sapcc/andromeda/restapi/operations/monitors"
	"github.com/sapcc/andromeda/restapi/operations/pools"
)

func ExecuteServer(server *restapi.Server) error {
	log.Info("Starting up andromeda-server")

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
	db, err := sqlx.Connect(u.Driver, u.DSN)
	if err != nil {
		log.WithError(err).WithField("driver", u.Driver).Fatal("Failed to connect to database")
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
	api.Logger = log.Infof

	// Prometheus Metrics
	if config.Global.Default.Prometheus {
		go utils.PrometheusListen()

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

	// Akamai Metrics
	api.MetricsGetMetricsAkamaiTotalDNSRequestsHandler = metrics_ops.GetMetricsAkamaiTotalDNSRequestsHandlerFunc(c.AkamaiMetrics.GetTotalDNSRequests)

	// Administrative
	api.AdministrativeGetServicesHandler = administrative.GetServicesHandlerFunc(c.Services.GetServices)
	api.AdministrativePostSyncHandler = administrative.PostSyncHandlerFunc(c.Sync.PostSync)
	api.AdministrativeGetCidrBlocksHandler = administrative.GetCidrBlocksHandlerFunc(c.CidrBlocks.GetCidrBlocks)

	// Quota Middleware
	if config.Global.Quota.Enabled {
		log.Info("Initializing quota middleware")

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

	//_rpc worker
	_rpc := RPCServer(db)

	// run the api and _rpc server
	go func() {
		log.Infof("RPC Listening on %v", _rpc.Subjects())
		if err = _rpc.Run(); err != nil {
			log.WithError(err).Fatal("listening on RPC")
		}
	}()
	go func() {
		if err = server.Serve(); err != nil {
			log.WithError(err).Fatal("listening on API")
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	<-done
	log.Infof("Shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	defer func() {
		if err = _rpc.Shutdown(ctx); err != nil {
			log.WithError(err).Fatal("shutting down RPC")
		}
		if err = server.Shutdown(); err != nil {
			log.WithError(err).Fatal("shutting down API")
		}
	}()

	return nil
}

func RPCServer(db *sqlx.DB) *stormrpc.Server {
	srv := rpc.NewServer("andromeda-server")
	svc := &server.RPCHandler{DB: db}
	server.RegisterRPCServerServer(srv, svc)
	return srv
}
