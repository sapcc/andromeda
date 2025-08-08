// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package f5

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/actatum/stormrpc"
	"github.com/apex/log"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/rpc"
	"github.com/sapcc/andromeda/internal/rpc/server"
	"github.com/sapcc/andromeda/internal/utils"

	"github.com/f5devcentral/go-bigip"
)

const agentName = "f5-as3-declaration"

var (
	lastSyncTimestampGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "andromeda_agent_last_sync_timestamp",
			Help: "Last time an agent has successfully completed its sync loop (sync completion timestamp)",
		},
		[]string{"agent"},
	)

	lastSyncDurationDurationSecondsGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "andromeda_agent_last_sync_duration_seconds",
			Help: "Last time an agent has successfully completed its sync loop (sync duration in seconds)",
		},
		[]string{"agent"},
	)
)

type F5Agent struct {
	bigIP              *bigip.BigIP
	declarationBuilder AS3DeclarationBuilder
	rpc                server.RPCServerClient
}

type FullSync struct {
	f5 *F5Agent
}

func init() {
	prometheus.MustRegister(lastSyncTimestampGauge)
	prometheus.MustRegister(lastSyncDurationDurationSecondsGauge)
}

// Method can be of any name
func (s *FullSync) FullSync(ctx context.Context, req stormrpc.Request) stormrpc.Response {
	log.WithField("request", req).Info("[pubsub.1] Received event")
	// do something with event
	s.f5.fullSync()

	resp, err := stormrpc.NewResponse(req.Reply, nil)
	if err != nil {
		return stormrpc.NewErrorResponse(req.Reply, err)
	}

	return resp
}

func ExecuteF5Agent() error {
	log.Debugf("Enabled=%+v Devices=%v VCMPs=%v PhysicalNetwork=%v",
		config.Global.F5Config.Enabled,
		config.Global.F5Config.Devices,
		config.Global.F5Config.VCMPs,
		config.Global.F5Config.PhysicalNetwork,
	)

	var activeF5Session *bigip.BigIP
	for _, url := range config.Global.F5Config.Devices {
		deviceSession, err := GetBigIPSession(url)
		if err != nil {
			return fmt.Errorf("failed to acquire F5 device session: %v", err)
		}
		device, err := GetActiveDevice(deviceSession)
		if err != nil {
			return fmt.Errorf("failed to determine whether F5 device is active: %v", err)
		}
		if device != nil {
			activeF5Session = deviceSession
			log.Infof("Connected to F5 device [marketing name = %q, name = %q, version = %s, edition = %q, failover state = %q]",
				device.MarketingName, device.Name, device.Version, device.Edition, device.FailoverState)
		}
	}

	if activeF5Session == nil {
		return errors.New("failed to determine active F5 session")
	}

	nc, err := nats.Connect(config.Global.Default.TransportURL)
	if err != nil {
		return err
	}
	client, err := stormrpc.NewClient("", stormrpc.WithNatsConn(nc))
	if err != nil {
		return err
	}

	// Create F5 worker instance with Server RPC interface
	rpcClient := server.NewRPCServerClient(client)
	f5 := F5Agent{
		bigIP:              activeF5Session,
		declarationBuilder: NewAS3DeclarationBuilder(NewAndromedaF5Store(rpcClient)),
		rpc:                rpcClient,
	}

	srv := rpc.NewServer("andromeda-f5-agent", stormrpc.WithNatsConn(nc))
	fs := &FullSync{&f5}

	// Allows the sync to be invoked over RPC via an HTTP handler in
	// Andromeda Server (see `m31ctl sync`)
	srv.Handle("andromeda.sync", fs.FullSync)

	go f5.fullSync()
	go func() {
		_ = srv.Run()
	}()
	if config.Global.Default.Prometheus {
		go utils.PrometheusListen()
	}

	log.Infof("👋 Listening on %v", srv.Subjects())

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
	log.Infof("💀 Shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(ctx)
}

func (f5 *F5Agent) fullSync() {
	syncInterval := 5 * time.Minute
	sync := func() {
		syncStart := time.Now()
		err := f5.Sync()
		elapsed := time.Since(syncStart)
		if err != nil {
			log.Errorf("Sync failed after %s (next iteration in %s): %s", elapsed, syncInterval, err.Error())
			return
		}
		log.Infof("Sync completed in %s (next iteration in %s)", elapsed, syncInterval)
		lastSyncTimestampGauge.WithLabelValues(agentName).Set(float64(time.Now().Unix()))
		lastSyncDurationDurationSecondsGauge.WithLabelValues(agentName).Set(elapsed.Seconds())
	}
	sync()
	c := time.Tick(syncInterval)
	for {
		<-c
		sync()
	}
}

// Sync relies on the AS3 `POST /declare` endpoint, therefore all entities must
// be included in the payload.
func (f5 *F5Agent) Sync() error {
	decl, rpcRequest, err := f5.declarationBuilder.Build()
	if err != nil {
		return err
	}
	jsonDoc, err := json.Marshal(decl)
	if err != nil {
		return err
	}
	log.Debugf("AS3 declaration: %s", string(jsonDoc))
	log.Debugf("RPC provisioning status updates: %v", rpcRequest.ProvisioningStatus)
	if err := postAS3Declaration(decl, f5.bigIP, sanityCheckAS3Declaration); err != nil {
		return err
	}
	log.Debugf("Posted AS3 declaration successfully")
	if _, err := f5.rpc.UpdateProvisioningStatus(context.Background(), rpcRequest); err != nil {
		return err
	}
	log.Debugf("Posted RPC provisioning status updates successfully")
	return nil
}
