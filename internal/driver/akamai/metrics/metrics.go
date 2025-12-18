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

	// Custom metrics collector
	prometheus.MustRegister(availabilityAliveGauge)
	prometheus.MustRegister(availabilityHandedOutGauge)
	prometheus.MustRegister(availabilityScoreGauge)
	prometheus.MustRegister(availabilityLastSyncGauge)
	prometheus.MustRegister(availabilitySyncErrorsCounter)
	prometheus.MustRegister(availabilityLastReportPeriodGauge)

	prometheus.MustRegister(requestsCounter)
	prometheus.MustRegister(requestsLastSyncGauge)
	prometheus.MustRegister(requestsSyncErrorsCounter)
	prometheus.MustRegister(requestsLastReportPeriodGauge)

	prometheus.MustRegister(rateLimitingDurationSeconds)

	go func() {
		akamaiSession := NewCachedAkamaiSession(*session, config.Global.AkamaiConfig.Domain, NewAkamaiRateLimiter(100))
		rpcClient := NewCachedRPCClient(client)

		// Akamai API limitation, see <https://techdocs.akamai.com/gtm-reporting/reference/get-traffic-property>
		// Akamai will update the traffic report **at best** every 5 minutes.
		interval := 5 * time.Minute
		for {
			AkamaiCustomMetricsSync(akamaiSession, rpcClient)
			time.Sleep(interval)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
	log.Info("Shutting down")

	return nil
}
