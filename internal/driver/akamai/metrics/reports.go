/*
 *   Copyright 2025 SAP SE
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
