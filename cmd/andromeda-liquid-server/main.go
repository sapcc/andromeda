// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/go-openapi/runtime"
	runtimeclient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/gophercloud/gophercloud/v2"
	"github.com/majewsky/gg/option"
	"github.com/majewsky/gg/options"
	"github.com/sapcc/go-api-declarations/bininfo"
	"github.com/sapcc/go-api-declarations/liquid"
	"github.com/sapcc/go-bits/httpext"
	"github.com/sapcc/go-bits/liquidapi"
	"github.com/sapcc/go-bits/logg"
	"github.com/sapcc/go-bits/must"
	"github.com/sapcc/go-bits/respondwith"
	"github.com/urfave/cli/v2"

	"github.com/sapcc/andromeda/client"
	"github.com/sapcc/andromeda/client/administrative"
	"github.com/sapcc/andromeda/models"
)

type liquidLogic struct {
	andromedaClient *client.Andromeda

	// See "--gmt-endpoint" CLI flag in main() function (e.g. "http://localhost:8080/v1")
	GTMEndpoint string
}

type ReauthTransport struct {
	transport http.RoundTripper
	context   context.Context
	provider  *gophercloud.ProviderClient
}

func (rt *ReauthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := rt.transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusUnauthorized {
		if err := rt.provider.Reauthenticate(rt.context, rt.provider.Token()); err != nil {
			logg.Error("ReauthTransport failed to reauthenticate: %s", err)
		}
	}
	return resp, nil
}

var defaultServiceUsageReport = liquid.ServiceUsageReport{
	InfoVersion: 1,
	Resources: map[liquid.ResourceName]*liquid.ResourceUsageReport{
		"datacenters": {
			Quota: option.Some(int64(0)),
			PerAZ: map[liquid.AvailabilityZone]*liquid.AZResourceUsageReport{
				liquid.AvailabilityZoneAny: {Usage: 0},
			},
		},
		"domains_akamai": {
			Quota: option.Some(int64(0)),
			PerAZ: map[liquid.AvailabilityZone]*liquid.AZResourceUsageReport{
				liquid.AvailabilityZoneAny: {Usage: 0},
			},
		},
		"domains_f5": {
			Quota: option.Some(int64(0)),
			PerAZ: map[liquid.AvailabilityZone]*liquid.AZResourceUsageReport{
				liquid.AvailabilityZoneAny: {Usage: 0},
			},
		},
		"members": {
			Quota: option.Some(int64(0)),
			PerAZ: map[liquid.AvailabilityZone]*liquid.AZResourceUsageReport{
				liquid.AvailabilityZoneAny: {Usage: 0},
			},
		},
		"monitors": {
			Quota: option.Some(int64(0)),
			PerAZ: map[liquid.AvailabilityZone]*liquid.AZResourceUsageReport{
				liquid.AvailabilityZoneAny: {Usage: 0},
			},
		},
		"pools": {
			Quota: option.Some(int64(0)),
			PerAZ: map[liquid.AvailabilityZone]*liquid.AZResourceUsageReport{
				liquid.AvailabilityZoneAny: {Usage: 0},
			},
		},
	},
}

func (l *liquidLogic) Init(context context.Context, provider *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) error {
	l.andromedaClient = client.Default
	if l.GTMEndpoint == "" {
		endpointOpts := gophercloud.EndpointOpts{
			Region: eo.Region,
		}
		endpointOpts.ApplyDefaults("gtm")
		endpoint, err := provider.EndpointLocator(endpointOpts)
		if err != nil {
			return err
		}
		l.GTMEndpoint = endpoint
	}
	logg.Debug("Using GTM endpoint %s", l.GTMEndpoint)
	uri, err := url.Parse(l.GTMEndpoint)
	if err != nil {
		return err
	}
	rt := runtimeclient.NewWithClient(
		uri.Host, uri.Path, []string{uri.Scheme},
		&http.Client{
			Transport: &ReauthTransport{http.DefaultTransport, context, provider},
		},
	)
	rt.DefaultAuthentication = runtime.ClientAuthInfoWriterFunc(func(req runtime.ClientRequest, reg strfmt.Registry) error {
		if err := req.SetHeaderParam("X-Auth-Token", provider.Token()); err != nil {
			return err
		}
		return nil
	})
	l.andromedaClient.SetTransport(rt)
	return nil
}

func (l *liquidLogic) BuildServiceInfo(_ context.Context) (liquid.ServiceInfo, error) {
	return liquid.ServiceInfo{
		Version: 1,
		Resources: map[liquid.ResourceName]liquid.ResourceInfo{
			"datacenters": {
				HasCapacity: false,
				HasQuota:    true,
				Topology:    liquid.FlatTopology,
			},
			"domains_akamai": {
				HasCapacity: false,
				HasQuota:    true,
				Topology:    liquid.FlatTopology,
			},
			"domains_f5": {
				HasCapacity: false,
				HasQuota:    true,
				Topology:    liquid.FlatTopology,
			},
			"members": {
				HasCapacity: false,
				HasQuota:    true,
				Topology:    liquid.FlatTopology,
			},
			"monitors": {
				HasCapacity: false,
				HasQuota:    true,
				Topology:    liquid.FlatTopology,
			},
			"pools": {
				HasCapacity: false,
				HasQuota:    true,
				Topology:    liquid.FlatTopology,
			},
		},
		Rates:                  map[liquid.RateName]liquid.RateInfo{},
		CapacityMetricFamilies: map[liquid.MetricName]liquid.MetricFamilyInfo{},
	}, nil
}

func (l *liquidLogic) ScanCapacity(ctx context.Context, req liquid.ServiceCapacityRequest, serviceInfo liquid.ServiceInfo) (liquid.ServiceCapacityReport, error) {
	return liquid.ServiceCapacityReport{InfoVersion: 1}, nil
}

func (l *liquidLogic) ScanUsage(ctx context.Context, projectUUID string, req liquid.ServiceUsageRequest, serviceInfo liquid.ServiceInfo) (liquid.ServiceUsageReport, error) {
	params := administrative.NewGetQuotasProjectIDParams().WithDefaults().WithContext(ctx)
	params.ProjectID = projectUUID
	resp, err := l.andromedaClient.Administrative.GetQuotasProjectID(params)
	if err != nil {
		var notFoundErr *administrative.GetQuotasProjectIDNotFound
		if errors.As(err, &notFoundErr) {
			return defaultServiceUsageReport, nil
		}
		return liquid.ServiceUsageReport{}, err
	}
	return liquid.ServiceUsageReport{
		InfoVersion: 1,
		Resources: map[liquid.ResourceName]*liquid.ResourceUsageReport{
			"datacenters": {
				Quota: options.FromPointer(resp.Payload.Quota.Datacenter),
				PerAZ: map[liquid.AvailabilityZone]*liquid.AZResourceUsageReport{
					liquid.AvailabilityZoneAny: {Usage: uint64(resp.Payload.Quota.QuotaUsage.InUseDatacenter)},
				},
			},
			"domains_akamai": {
				Quota: options.FromPointer(resp.Payload.Quota.DomainAkamai),
				PerAZ: map[liquid.AvailabilityZone]*liquid.AZResourceUsageReport{
					liquid.AvailabilityZoneAny: {Usage: uint64(resp.Payload.Quota.QuotaUsage.InUseDomainAkamai)},
				},
			},
			"domains_f5": {
				Quota: options.FromPointer(resp.Payload.Quota.DomainF5),
				PerAZ: map[liquid.AvailabilityZone]*liquid.AZResourceUsageReport{
					liquid.AvailabilityZoneAny: {Usage: uint64(resp.Payload.Quota.QuotaUsage.InUseDomainF5)},
				},
			},
			"members": {
				Quota: options.FromPointer(resp.Payload.Quota.Member),
				PerAZ: map[liquid.AvailabilityZone]*liquid.AZResourceUsageReport{
					liquid.AvailabilityZoneAny: {Usage: uint64(resp.Payload.Quota.QuotaUsage.InUseMember)},
				},
			},
			"monitors": {
				Quota: options.FromPointer(resp.Payload.Quota.Monitor),
				PerAZ: map[liquid.AvailabilityZone]*liquid.AZResourceUsageReport{
					liquid.AvailabilityZoneAny: {Usage: uint64(resp.Payload.Quota.QuotaUsage.InUseMonitor)},
				},
			},
			"pools": {
				Quota: options.FromPointer(resp.Payload.Quota.Pool),
				PerAZ: map[liquid.AvailabilityZone]*liquid.AZResourceUsageReport{
					liquid.AvailabilityZoneAny: {Usage: uint64(resp.Payload.Quota.QuotaUsage.InUsePool)},
				},
			},
		},
	}, nil
}

func (l *liquidLogic) SetQuota(ctx context.Context, projectUUID string, req liquid.ServiceQuotaRequest, _ liquid.ServiceInfo) error {
	params := administrative.NewPutQuotasProjectIDParams().WithDefaults().WithContext(ctx)
	params.ProjectID = projectUUID
	params.Quota = administrative.PutQuotasProjectIDBody{Quota: &models.Quota{
		Datacenter:   func(num uint64) *int64 { i := int64(num); return &i }(req.Resources["datacenters"].Quota),
		DomainAkamai: func(num uint64) *int64 { i := int64(num); return &i }(req.Resources["domains_akamai"].Quota),
		DomainF5:     func(num uint64) *int64 { i := int64(num); return &i }(req.Resources["domains_f5"].Quota),
		Member:       func(num uint64) *int64 { i := int64(num); return &i }(req.Resources["members"].Quota),
		Monitor:      func(num uint64) *int64 { i := int64(num); return &i }(req.Resources["monitors"].Quota),
		Pool:         func(num uint64) *int64 { i := int64(num); return &i }(req.Resources["pools"].Quota),
	}}
	_, err := l.andromedaClient.Administrative.PutQuotasProjectID(params)
	if err != nil {
		return err
	}
	return nil
}

func (l *liquidLogic) ReviewCommitmentChange(ctx context.Context, req liquid.CommitmentChangeRequest, serviceInfo liquid.ServiceInfo) (liquid.CommitmentChangeResponse, error) {
	err := errors.New("this liquid does not manage commitments")
	return liquid.CommitmentChangeResponse{}, respondwith.CustomStatus(http.StatusBadRequest, err)
}

func main() {
	app := &cli.App{
		Name:  "andromeda-liquid-api",
		Usage: "andromeda liquid api",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "gtm-endpoint",
				Usage: "The GTM endpoint (useful for local testing against local Andromeda Server)",
				Value: "",
			},
			&cli.StringFlag{
				Name:  "host",
				Usage: "The IP to listen on",
				Value: "0.0.0.0",
			},
			&cli.IntFlag{
				Name:  "port",
				Usage: "Port to listen",
				Value: 8080,
			},
			&cli.BoolFlag{
				Name:    "liquid_debug",
				Usage:   "Enable verbose logging for Liquid HTTP handling",
				EnvVars: []string{"LIQUID_DEBUG"},
			},
		},
		Action: func(c *cli.Context) error {
			logg.ShowDebug = c.Bool("liquid_debug")
			ctx := httpext.ContextWithSIGINT(c.Context, 10*time.Second)
			host := c.String("host")
			port := c.Int("port")
			logic := &liquidLogic{GTMEndpoint: c.String("gtm-endpoint")}
			opts := liquidapi.RunOpts{DefaultListenAddress: fmt.Sprintf("%s:%d", host, port)}
			return liquidapi.Run(ctx, logic, opts)
		},
		Version: bininfo.Version(),
	}
	must.Succeed(app.Run(os.Args))
}
