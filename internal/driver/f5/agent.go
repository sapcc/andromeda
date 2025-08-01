// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package f5

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/actatum/stormrpc"
	"github.com/apex/log"
	"github.com/nats-io/nats.go"

	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/rpc"
	"github.com/sapcc/andromeda/internal/rpc/server"
	"github.com/sapcc/andromeda/internal/utils"

	"github.com/f5devcentral/go-bigip"

	"github.com/sapcc/andromeda/internal/driver/f5/as3"
)

type F5Agent struct {
	bigIP *bigip.BigIP
	rpc   server.RPCServerClient
}

type FullSync struct {
	f5 *F5Agent
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
	f5 := F5Agent{
		activeF5Session,
		server.NewRPCServerClient(client),
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
		if err := f5.Sync(); err != nil {
			log.Errorf("Sync failed (next iteration in %s): %s", syncInterval, err.Error())
			return
		}
		log.Infof("Sync completed (next iteration in %s)", syncInterval)
		// TODO: write a `andromeda_last_sync{agent: f5}` gauge (timestamp) metric that can be alerted on
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
	log.Info("Skipping for now.")
	return nil
	adc := as3.ADC{SchemaVersion: "3.22.0"}

	//
	response, err := f5.rpc.GetDomains(context.Background(), &server.SearchRequest{
		Provider:       "f5",
		PageNumber:     0,
		ResultPerPage:  1000,
		FullyPopulated: true,
	})
	if err != nil {
		return err
	}

	if response.GetResponse() == nil {
		return nil
	}
	common, err := f5.GetCommonDeclaration()
	if err != nil {
		return err
	}
	tenant, err := f5.GetTenantDeclaration(response.GetResponse())
	if err != nil {
		return err
	}

	adc.AddTenant("ExampleTenant", *tenant)
	adc.AddTenant("Common", *common)
	if err := postDeclaration(adc); err != nil {
		return err
	}

	var prov []*server.ProvisioningStatusRequest_ProvisioningStatus
	for _, domain := range response.GetResponse() {
		log.Infof("%+v", domain)
		prov = append(prov, &server.ProvisioningStatusRequest_ProvisioningStatus{
			Id:     domain.GetId(),
			Model:  server.ProvisioningStatusRequest_ProvisioningStatus_DOMAIN,
			Status: server.ProvisioningStatusRequest_ProvisioningStatus_ACTIVE,
		})

		for _, pool := range domain.GetPools() {
			prov = append(prov, &server.ProvisioningStatusRequest_ProvisioningStatus{
				Id:     pool.GetId(),
				Model:  server.ProvisioningStatusRequest_ProvisioningStatus_POOL,
				Status: server.ProvisioningStatusRequest_ProvisioningStatus_ACTIVE,
			})
			for _, member := range pool.GetMembers() {
				prov = append(prov, &server.ProvisioningStatusRequest_ProvisioningStatus{
					Id:     member.GetId(),
					Model:  server.ProvisioningStatusRequest_ProvisioningStatus_MEMBER,
					Status: server.ProvisioningStatusRequest_ProvisioningStatus_ACTIVE,
				})
			}
			for _, monitor := range pool.GetMonitors() {
				prov = append(prov, &server.ProvisioningStatusRequest_ProvisioningStatus{
					Id:     monitor.GetId(),
					Model:  server.ProvisioningStatusRequest_ProvisioningStatus_MONITOR,
					Status: server.ProvisioningStatusRequest_ProvisioningStatus_ACTIVE,
				})
			}
		}
	}
	resp, err := f5.rpc.UpdateProvisioningStatus(context.Background(),
		&server.ProvisioningStatusRequest{ProvisioningStatus: prov})
	if err != nil {
		return err
	}
	log.Infof("%+v", resp)
	return nil
}

func postDeclaration(v interface{}) error {
	return fmt.Errorf("TO BE FIXED")
	/*
		js, err := json.MarshalIndent(v, "", "\t")
		if err != nil {
			return err
		}

		log.Infof("\n%s\n", string(js))

		session, err := GetBigIPSession()
		if err != nil {
			return err
		}
		req := &bigip.APIRequest{
			Method:      "post",
			URL:         "mgmt/shared/appsvcs/declare",
			Body:        string(js),
			ContentType: "application/json",
		}
		resp, err := session.APICall(req)
		if err != nil {
			print(err)
		}

		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, resp, "", "\t"); err != nil {
			return err
		}
		log.Info(prettyJSON.String())

		var response as3.Response
		if err := json.Unmarshal(resp, &response); err != nil {
			return err
		}

		for _, result := range response.Results {
			if result.Code != 200 {
				return fmt.Errorf("failed post! %s", result.Message)
			} else {
				log.Infof("succeeded: %s", result.Tenant)
			}
		}

		return nil
	*/
}
