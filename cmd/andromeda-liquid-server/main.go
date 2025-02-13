package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/go-openapi/runtime"
	runtimeclient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/gophercloud/gophercloud/v2"
	"github.com/sapcc/go-api-declarations/liquid"
	"github.com/sapcc/go-bits/httpext"
	"github.com/sapcc/go-bits/liquidapi"
	"github.com/sapcc/go-bits/logg"
	"github.com/sapcc/go-bits/must"
	"github.com/urfave/cli/v2"

	"github.com/sapcc/andromeda/client"
	"github.com/sapcc/andromeda/client/administrative"
	"github.com/sapcc/andromeda/models"
)

type liquidLogic struct {
	andromedaClient *client.Andromeda
}

func (l *liquidLogic) Init(_ context.Context, provider *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) error {
	l.andromedaClient = client.Default
	endpointOpts := gophercloud.EndpointOpts{
		Region: eo.Region,
	}
	endpointOpts.ApplyDefaults("gtm")
	endpoint, err := provider.EndpointLocator(endpointOpts)
	if err != nil {
		return err
	}
	uri, err := url.Parse(endpoint)
	if err != nil {
		return err
	}
	rt := runtimeclient.New(uri.Host, uri.Path, []string{uri.Scheme})
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
				Topology:    liquid.FlatResourceTopology,
			},
			"domains": {
				HasCapacity: false,
				HasQuota:    true,
				Topology:    liquid.FlatResourceTopology,
			},
			"members": {
				HasCapacity: false,
				HasQuota:    true,
				Topology:    liquid.FlatResourceTopology,
			},
			"monitors": {
				HasCapacity: false,
				HasQuota:    true,
				Topology:    liquid.FlatResourceTopology,
			},
			"pools": {
				HasCapacity: false,
				HasQuota:    true,
				Topology:    liquid.FlatResourceTopology,
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
		return liquid.ServiceUsageReport{}, err
	}
	return liquid.ServiceUsageReport{
		InfoVersion: 1,
		Resources: map[liquid.ResourceName]*liquid.ResourceUsageReport{
			"datacenters": {
				Quota: resp.Payload.Quota.Datacenter,
				PerAZ: map[liquid.AvailabilityZone]*liquid.AZResourceUsageReport{
					liquid.AvailabilityZoneAny: {Usage: uint64(resp.Payload.Quota.QuotaUsage.InUseDatacenter)},
				},
			},
			"domains": {
				Quota: resp.Payload.Quota.Domain,
				PerAZ: map[liquid.AvailabilityZone]*liquid.AZResourceUsageReport{
					liquid.AvailabilityZoneAny: {Usage: uint64(resp.Payload.Quota.QuotaUsage.InUseDomain)},
				},
			},
			"members": {
				Quota: resp.Payload.Quota.Member,
				PerAZ: map[liquid.AvailabilityZone]*liquid.AZResourceUsageReport{
					liquid.AvailabilityZoneAny: {Usage: uint64(resp.Payload.Quota.QuotaUsage.InUseMember)},
				},
			},
			"monitors": {
				Quota: resp.Payload.Quota.Monitor,
				PerAZ: map[liquid.AvailabilityZone]*liquid.AZResourceUsageReport{
					liquid.AvailabilityZoneAny: {Usage: uint64(resp.Payload.Quota.QuotaUsage.InUseMonitor)},
				},
			},
			"pools": {
				Quota: resp.Payload.Quota.Pool,
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
		Datacenter: func(num uint64) *int64 { i := int64(num); return &i }(req.Resources["datacenters"].Quota),
		Domain:     func(num uint64) *int64 { i := int64(num); return &i }(req.Resources["domains"].Quota),
		Member:     func(num uint64) *int64 { i := int64(num); return &i }(req.Resources["members"].Quota),
		Monitor:    func(num uint64) *int64 { i := int64(num); return &i }(req.Resources["monitors"].Quota),
		Pool:       func(num uint64) *int64 { i := int64(num); return &i }(req.Resources["pools"].Quota),
	}}
	_, err := l.andromedaClient.Administrative.PutQuotasProjectID(params)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	app := &cli.App{
		Name:  "andromeda-liquid-api",
		Usage: "andromeda liquid api",
		Flags: []cli.Flag{
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
			logic := &liquidLogic{}
			opts := liquidapi.RunOpts{DefaultListenAddress: fmt.Sprintf("%s:%d", host, port)}
			return liquidapi.Run(ctx, logic, opts)
		},
	}
	must.Succeed(app.Run(os.Args))
}
