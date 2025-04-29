// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package rpc

import (
	"context"
	"strings"

	"github.com/actatum/stormrpc"
	"github.com/actatum/stormrpc/middleware"
	"github.com/apex/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sapcc/go-api-declarations/bininfo"

	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/middlewares"
)

// NewServer creates a new RPC server with the given name and options.
func NewServer(name string, opts ...stormrpc.ServerOption) *stormrpc.Server {
	opts = append(opts, stormrpc.WithErrorHandler(func(ctx context.Context, err error) {
		log.WithError(err).Error("RPC Error")
	}))

	version := strings.TrimLeft(bininfo.VersionOr("-unknown"), "v")
	srv, err := stormrpc.NewServer(&stormrpc.ServerConfig{
		NatsURL: config.Global.Default.TransportURL,
		Name:    name,
		Version: version[:strings.LastIndex(version, "-")],
	}, opts...)
	if err != nil {
		log.Fatal(err.Error())
	}

	if srv != nil {
		log.Info("Loading RPC middleware")

		if config.Global.Default.Prometheus && config.Global.Default.PrometheusRPCMetrics {
			log.Info("Enabling Prometheus RPC metrics")
			rpcRequestHistogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
				Subsystem: "rpc",
				Name:      "request_duration_seconds",
				Help:      "The latency of the RPC requests.",
				Buckets:   prometheus.DefBuckets,
			}, []string{"subject"})
			prometheus.MustRegister(rpcRequestHistogram)

			srv.Use(middlewares.Tracing(rpcRequestHistogram))
		}

		srv.Use(middleware.RequestID, middlewares.Logging, middlewares.Recoverer)
	}
	return srv
}
