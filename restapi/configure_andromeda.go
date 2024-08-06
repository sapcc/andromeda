// This file is safe to edit. Once it exists it will not be overwritten

// Copyright 2020 SAP SE
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package restapi

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/apex/log"
	"github.com/didip/tollbooth"
	"github.com/dre1080/recovr"
	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/rs/cors"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	"github.com/slok/go-http-metrics/middleware/std"

	"github.com/sapcc/andromeda/internal/auth"
	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/middlewares"
	"github.com/sapcc/andromeda/restapi/operations"
)

//go:generate swagger generate server --target ../../andromeda --name Andromeda --spec ../swagger.yml --principal interface{} --exclude-main

func configureFlags(api *operations.AndromedaAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.AndromedaAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError
	api.Logger = log.Infof

	// To continue using redoc as your UI, uncomment the following line
	api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()
	api.PreServerShutdown = func() {}
	api.ServerShutdown = func() {
		sentry.Flush(5 * time.Second)
	}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	if rl := config.Global.ApiSettings.RateLimit; rl > .0 {
		log.Infof("Initializing rate limit middleware (rate_limit=%f)", rl)
		limiter := tollbooth.NewLimiter(rl, nil)
		handler = tollbooth.LimitHandler(limiter, handler)
	}

	if config.Global.Audit.Enabled {
		log.Info("Initializing audit middleware")
		auditMiddleware := middlewares.NewAuditController()
		handler = auditMiddleware.AuditHandler(handler)
	}

	switch config.Global.ApiSettings.AuthStrategy {
	case "keystone":
		log.Info("Initializing keystone middleware")
		var err error
		handler, err = auth.KeystoneMiddleware(handler)
		if err != nil {
			log.Errorf("Error initializing keystone middleware: %w", err)
		}
	}

	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	handler = middlewares.HealthCheckMiddleware(handler)

	if !config.Global.ApiSettings.DisableCors {
		log.Info("Initializing CORS middleware")
		handler = cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"HEAD", "GET", "POST", "PUT", "DELETE"},
			AllowedHeaders: []string{"Content-Type", "User-Agent", "X-Auth-Token"},
		}).Handler(handler)
	}

	handler = sentryhttp.New(sentryhttp.Options{
		Repanic: true,
	}).Handle(handler)

	if config.Global.Default.Prometheus {
		handler = std.Handler("", middleware.New(middleware.Config{
			Recorder:      metrics.NewRecorder(metrics.Config{}),
			GroupedStatus: true,
		}), handler)
	}

	return recovr.New()(handler)
}
