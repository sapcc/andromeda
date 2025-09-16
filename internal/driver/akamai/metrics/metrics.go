// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package metrics

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/actatum/stormrpc"
	"github.com/apex/log"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/driver/akamai"
	"github.com/sapcc/andromeda/internal/utils"
)

func ExecuteAkamaiMetrics() error {
	go utils.PrometheusListen()

	// Akamai connection
	session, _ := akamai.NewAkamaiSession(&config.Global.AkamaiConfig)

	// RPC connection
	nc, err := nats.Connect(config.Global.Default.TransportURL)
	if err != nil {
		return err
	}
	client, err := stormrpc.NewClient("", stormrpc.WithNatsConn(nc))
	if err != nil {
		return err
	}

	// Experimental: custom metrics collector
	prometheus.MustRegister(requestsCounter)
	prometheus.MustRegister(requestsLastSyncGauge)
	prometheus.MustRegister(requestsSyncErrorsCounter)
	prometheus.MustRegister(requestsLastReportPeriodGauge)

	go func() {
		akamaiSession := NewCachedAkamaiSession(*session, config.Global.AkamaiConfig.Domain)
		rpcClient := NewCachedRPCClient(client)
		AkamaiCustomMetricsSync(akamaiSession, rpcClient)

		// Akamai API limitation, see <https://techdocs.akamai.com/gtm-reporting/reference/get-traffic-property>
		interval := 5 * time.Minute

		c := time.Tick(interval)
		for {
			<-c
			AkamaiCustomMetricsSync(akamaiSession, rpcClient)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
	log.Info("Shutting down")

	return nil
}
