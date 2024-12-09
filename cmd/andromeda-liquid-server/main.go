package main

import (
	"context"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/go-api-declarations/liquid"
	"github.com/sapcc/go-bits/httpext"
	"github.com/sapcc/go-bits/liquidapi"
	"github.com/sapcc/go-bits/must"
	"github.com/urfave/cli/v2"
)

type liquidLogic struct {
}

func (p liquidLogic) Init(ctx context.Context, provider *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) error {
	return nil
}

func (p liquidLogic) BuildServiceInfo(ctx context.Context) (liquid.ServiceInfo, error) {
	return liquid.ServiceInfo{}, nil
}

func (p liquidLogic) ScanCapacity(ctx context.Context, req liquid.ServiceCapacityRequest, serviceInfo liquid.ServiceInfo) (liquid.ServiceCapacityReport, error) {
	return liquid.ServiceCapacityReport{}, nil
}

func (p liquidLogic) ScanUsage(ctx context.Context, projectUUID string, req liquid.ServiceUsageRequest, serviceInfo liquid.ServiceInfo) (liquid.ServiceUsageReport, error) {
	return liquid.ServiceUsageReport{}, nil
}

func (p liquidLogic) SetQuota(ctx context.Context, projectUUID string, req liquid.ServiceQuotaRequest, serviceInfo liquid.ServiceInfo) error {
	return nil
}

func main() {
	config.ParseArgsAndRun("andromeda-liquid-api", "andromeda liquid api",
		func(c *cli.Context) error {
			ctx := httpext.ContextWithSIGINT(c.Context, 10*time.Second)
			logic := liquidLogic{}
			opts := liquidapi.RunOpts{DefaultListenAddress: "127.0.0.1:8080"}
			must.Succeed(liquidapi.Run(ctx, logic, opts))
			return nil
		})
}
