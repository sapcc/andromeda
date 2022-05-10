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
	"time"

	"go-micro.dev/v4/metadata"

	"github.com/sapcc/andromeda/internal/rpc/server"

	"go-micro.dev/v4"

	"github.com/scottdware/go-bigip"
	"go-micro.dev/v4/logger"

	"github.com/sapcc/andromeda/internal/driver/f5/as3"
	"github.com/sapcc/andromeda/internal/rpc/worker"
)

type F5Agent struct {
	bigIP *bigip.BigIP
	rpc   server.RPCServerService
}

// All methods of Sub will be executed when
// a message is received
type Sub struct {
	f5 *F5Agent
}

// Method can be of any name
func (s *Sub) Process(ctx context.Context, request *worker.SyncRequest) error {
	md, _ := metadata.FromContext(ctx)
	logger.Info("[pubsub.1] Received event %+v with metadata %+v\n", request, md)
	// do something with event
	s.f5.fullSync()
	return nil
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
	logger.Infof("connected to %s %s (%s)", device.MarketingName, device.Name, device.Version)

	if err := BigIPSupportsDNS(device); err != nil {
		return err
	}

	metadata := map[string]string{
		"type":    "F5",
		"host":    device.Hostname,
		"version": device.Version,
	}
	service := micro.NewService(
		micro.Name("andromeda.agent.f5"),
		micro.Metadata(metadata),
		micro.RegisterTTL(time.Second*60),
		micro.RegisterInterval(time.Second*30),
	)
	service.Init()

	// Create F5 worker instance with Server RPC interface
	f5 := F5Agent{
		session,
		server.NewRPCServerService("andromeda.server", service.Client()),
	}

	// Register callbacks
	/*if err := worker.RegisterRPCWorkerHandler(service.Server(), &f5); err != nil {
		return err
	}*/
	if err := micro.RegisterSubscriber("andromeda.sync_all", service.Server(), &Sub{&f5}); err != nil {
		logger.Error(err)
	}

	go f5.periodicSync()
	go f5.fullSync()
	return service.Run()
}

func (f5 *F5Agent) fullSync() {
	if err := f5.Sync(false); err != nil {
		logger.Error(err)
	}
}

func (f5 *F5Agent) periodicSync() {
	t := time.NewTicker(10 * time.Second)
	defer t.Stop()
	for {
		<-t.C // Activate periodically
		if err := f5.Sync(true); err != nil {
			logger.Error(err)
		}
	}

}

func (s *F5Agent) SyncAll(ctx context.Context, request *worker.SyncRequest) error {
	logger.Info("Sync invoked!")
	md, _ := metadata.FromContext(ctx)
	logger.Info("[pubsub.1] Received event %+v with metadata %+v\n", request, md)
	// do something with event
	return nil
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
		logger.Info(domain)
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
	logger.Info(resp)
	return nil
}

func postDeclaration(v interface{}) error {
	js, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	logger.Info("\n", string(js), "\n")

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
	logger.Info(prettyJSON.String())

	var response as3.Response
	if err := json.Unmarshal(resp, &response); err != nil {
		return err
	}

	for _, result := range response.Results {
		if result.Code != 200 {
			return fmt.Errorf("Failed post! %s", result.Message)
		} else {
			logger.Info("Succeeded: ", result.Tenant)
		}
	}

	return nil
}
