// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package f5

import (
	"context"
	"encoding/json"
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
)

func init() {
	prometheus.MustRegister(lastSyncTimestampGauge)
	prometheus.MustRegister(lastSyncDurationSecondsGauge)
}

var lastSyncTimestampGauge = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "andromeda_agent_last_sync_timestamp",
		Help: "Last time an agent has successfully completed its sync loop (sync completion timestamp)",
	},
	[]string{"agent"},
)
var lastSyncDurationSecondsGauge = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "andromeda_agent_last_sync_duration_seconds",
		Help: "Last time an agent has successfully completed its sync loop (sync duration in seconds)",
	},
	[]string{"agent"},
)

type syncFunc func(session bigIPSession, rpc server.RPCServerClient) error
type instrumentedSyncFunc func(session bigIPSession, rpc server.RPCServerClient)

func syncWorker(syncInterval time.Duration, session bigIPSession, rpc server.RPCServerClient, syncFn instrumentedSyncFunc) {
	syncFn(session, rpc)
	c := time.Tick(syncInterval)
	for {
		<-c
		syncFn(session, rpc)
	}
}

func newInstrumentedSyncFunc(agentName string, syncInterval time.Duration, syncFn syncFunc) instrumentedSyncFunc {
	return func(session bigIPSession, rpc server.RPCServerClient) {
		syncStart := time.Now()
		err := syncFn(session, rpc)
		elapsed := time.Since(syncStart)
		if err != nil {
			log.Errorf("Sync failed after %s (next iteration in %s): %s", elapsed, syncInterval, err.Error())
			return
		}
		log.Infof("Sync completed in %s (next iteration in %s)", elapsed, syncInterval)
		lastSyncTimestampGauge.WithLabelValues(agentName).Set(float64(time.Now().Unix()))
		lastSyncDurationSecondsGauge.WithLabelValues(agentName).Set(elapsed.Seconds())
	}
}

func ExecuteF5Agent(agentName string, syncInterval time.Duration, syncFn syncFunc) error {
	log.Debugf("Enabled=%+v Devices=%v VCMPs=%v PhysicalNetwork=%v",
		config.Global.F5Config.Enabled,
		config.Global.F5Config.Devices,
		config.Global.F5Config.VCMPs,
		config.Global.F5Config.PhysicalNetwork,
	)

	activeF5Session, activeF5Device, err := getActiveDeviceSession(config.Global.F5Config, getBigIPSession, matchActiveDevice)
	if err != nil {
		return fmt.Errorf("failed to determine active F5 session: %w", err)
	}

	log.Infof("Connected to F5 device [marketing name = %q, name = %q, version = %s, edition = %q, failover state = %q]",
		activeF5Device.MarketingName, activeF5Device.Name, activeF5Device.Version, activeF5Device.Edition, activeF5Device.FailoverState)

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
	srv := rpc.NewServer(fmt.Sprintf("andromeda-%s-agent", agentName), stormrpc.WithNatsConn(nc))

	instrumentedSyncFunc := newInstrumentedSyncFunc(agentName, syncInterval, syncFn)

	// Allows the sync to be invoked over RPC via an HTTP handler in Andromeda Server
	// see `m31ctl sync`
	// TODO it doesn't work; the first handler to receive the message
	// runs the sync, while the other will not receive any event.
	srv.Handle("andromeda.sync", func(ctx context.Context, req stormrpc.Request) stormrpc.Response {
		log.WithField("request", req).Info("[pubsub.1] Received event")
		instrumentedSyncFunc(activeF5Session, rpcClient)
		resp, err := stormrpc.NewResponse(req.Reply, nil)

		// TODO if setting resp.Err works, then instrumentedSyncFunc
		// should return an error and make this possible.
		// Update: if set to an error, it will cause the Andromeda Server to panic
		// resp.Err = syncErr

		if err != nil {
			return stormrpc.NewErrorResponse(req.Reply, err)
		}
		return resp
	})

	go syncWorker(syncInterval, activeF5Session, rpcClient, instrumentedSyncFunc)
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

func ExecuteF5DeclarationAgent() error {
	return ExecuteF5Agent("f5-declaration", 5*time.Minute, declarationSync)
}

func ExecuteF5StatusAgent() error {
	return ExecuteF5Agent("f5-status", 5*time.Minute, statusSync)
}

func declarationSync(session bigIPSession, rpc server.RPCServerClient) error {
	decl, rpcRequest, err := buildAS3Declaration(NewAndromedaF5Store(rpc), buildAS3CommonTenant, buildAS3DomainTenant)
	if err != nil {
		return err
	}
	jsonDoc, err := json.Marshal(decl)
	if err != nil {
		return err
	}
	log.Debugf("AS3 declaration: %s", string(jsonDoc))
	log.Debugf("RPC provisioning status updates: %v", rpcRequest.ProvisioningStatus)
	if err := postAS3Declaration(decl, session, sanityCheckAS3Declaration); err != nil {
		return err
	}
	log.Debugf("Posted AS3 declaration successfully")
	if _, err := rpc.UpdateProvisioningStatus(context.Background(), rpcRequest); err != nil {
		return err
	}
	log.Debugf("Posted RPC provisioning status updates successfully")
	return nil
}

func statusSync(session bigIPSession, rpc server.RPCServerClient) error {
	req, err := buildMemberStatusUpdateRequest(session, NewAndromedaF5Store(rpc))
	if err != nil {
		return err
	}
	log.Debugf("RPC member status updates: %v", req.MemberStatus)
	if _, err = rpc.UpdateMemberStatus(context.Background(), req); err != nil {
		return err
	}
	log.Debugf("Posted RPC member status updates successfully")
	return nil
}
