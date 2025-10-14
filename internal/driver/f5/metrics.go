// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package f5

import (
	"github.com/apex/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sapcc/andromeda/internal/rpcmodels"
)

// these metrics are intended to be added to all F5 agents
func init() {
	prometheus.MustRegister(lastSyncTimestampGauge)
	prometheus.MustRegister(lastSyncDurationSecondsGauge)
}

// Operational metrics

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

// Federated metrics

var virtualServerPicksCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "andromeda_f5_virtual_server_picks",
		Help: "Total picks of the virtual server by the DNS resolver (derived from F5 iControlRest tm:gtm:server:serverstats endpoint)",
	},
	[]string{"domain", "datacenter_id", "project_id", "target_ip"},
)

func collectVirtualServerMetrics(session bigIPSession, store AndromedaF5Store, picksCounter *prometheus.CounterVec) error {
	datacenters, err := store.GetDatacenters()
	if err != nil {
		return err
	}
	datacentersByID := make(map[string]*rpcmodels.Datacenter, len(datacenters))
	for _, dc := range datacenters {
		datacentersByID[dc.Id] = dc
	}
	domains, err := store.GetDomains()
	if err != nil {
		return err
	}
	for _, d := range domains {
		for _, p := range d.Pools {
			for _, m := range p.Members {
				if _, exists := datacentersByID[m.DatacenterId]; !exists {
					log.Warnf("invalid datacenter ID for member [datacenter ID = %s, member ID = %s]", m.DatacenterId, m.Id)
					continue
				}
				if datacentersByID[m.DatacenterId] == nil {
					log.Warnf("nil datacenter for member [member ID = %s]", m.Id)
					continue
				}
				urlPath := serverStatsURL(
					as3DeclarationGSLBServerKey(m.Address, datacentersByID[m.DatacenterId].Name),
				)
				serverStats, err := fetchServerStats(session, urlPath)
				if err != nil {
					log.Warnf("failed to determine GSLB_Server picks [BigIP URL path = %s]: %s", urlPath, err)
					continue
				}
				picksCounter.WithLabelValues(d.Fqdn, m.DatacenterId, m.ProjectId, m.Address).Add(float64(serverStats.NestedStats.Entries.VirtualServerPicks.Value))
			}
		}
	}
	return nil
}
