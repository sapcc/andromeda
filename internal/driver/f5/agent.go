/*
 *   Copyright 2021 SAP SE
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

package f5

import (
	"bytes"
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

	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/rpc/server"

	"github.com/scottdware/go-bigip"

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
	session, err := GetBigIPSession()
	if err != nil {
		return fmt.Errorf("acquiring BigIP session: %v", err)
	}

	device, err := session.GetCurrentDevice()
	if err != nil {
		return err
	}
	log.Infof("connected to %s %s (%s)", device.MarketingName, device.Name, device.Version)

	if err = BigIPSupportsDNS(device); err != nil {
		return err
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
		session,
		server.NewRPCServerClient(client),
	}

	srv, err := stormrpc.NewServer(&stormrpc.ServerConfig{}, stormrpc.WithNatsConn(nc))
	if err != nil {
		return err
	}
	fs := &FullSync{&f5}
	srv.Handle("andromeda.sync", fs.FullSync)

	go f5.periodicSync()
	go f5.fullSync()
	go func() {
		_ = srv.Run()
	}()
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
	if err := f5.Sync(false); err != nil {
		log.Error(err.Error())
	}
}

func (f5 *F5Agent) periodicSync() {
	t := time.NewTicker(10 * time.Second)
	defer t.Stop()
	for {
		<-t.C // Activate periodically
		if err := f5.Sync(true); err != nil {
			log.Error(err.Error())
		}
	}

}

func (f5 *F5Agent) Sync(pending bool) error {
	adc := as3.ADC{SchemaVersion: "3.22.0"}

	// Tenants
	response, err := f5.rpc.GetDomains(context.Background(), &server.SearchRequest{
		Provider:       "f5",
		PageNumber:     0,
		ResultPerPage:  1000,
		FullyPopulated: true,
		Pending:        pending,
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
			return fmt.Errorf("Failed post! %s", result.Message)
		} else {
			log.Infof("Succeeded: %s", result.Tenant)
		}
	}

	return nil
}
