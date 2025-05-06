// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package metrics

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/actatum/stormrpc"
	"github.com/apex/log"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/driver/akamai"
	"github.com/sapcc/andromeda/internal/utils"
)

func ExecuteAkamaiReports() error {
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

	// Create a cached Akamai session
	cachedSession := NewCachedAkamaiSession(*session, config.Global.AkamaiConfig.Domain)

	// Register prometheus exporters
	prometheus.MustRegister(NewAndromedaAkamaiCollector(*session, client, config.Global.AkamaiConfig.Domain))
	
	// Register the total DNS requests collector
	prometheus.MustRegister(NewTotalDNSRequestsCollector(cachedSession, config.Global.AkamaiConfig.Domain))

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
	log.Info("Shutting down")

	return nil
}
