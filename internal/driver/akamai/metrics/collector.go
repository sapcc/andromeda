// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package metrics

import (
	"strconv"
	"strings"

	"github.com/actatum/stormrpc"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	"github.com/apex/log"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	fQNameNumReq  = prometheus.BuildFQName("andromeda", "akamai", "requests_5m")
	fQNameStatus  = prometheus.BuildFQName("andromeda", "akamai", "status_5m")
	fQNameUpdated = prometheus.BuildFQName("andromeda", "akamai", "last_update_timestamp")
	labels        = []string{"domain", "datacenter_id", "project_id", "target_ip"}
	descNumReq    = prometheus.NewDesc(fQNameNumReq, "Number of requests/5m per domain (most recent data point)", labels, nil)
	descStatus    = prometheus.NewDesc(fQNameStatus, "Status per domain (most recent data point)", labels, nil)
	descUpdated   = prometheus.NewDesc(fQNameUpdated, "Unix timestamp of when this metric was last updated from Akamai API", labels, nil)
)

type AndromedaAkamaiCollector struct {
	session          *CachedAkamaiSession
	rpc              *CachedRPCClient
	managementDomain string
	fqName           string
}

func NewAndromedaAkamaiCollector(session session.Session, client *stormrpc.Client, managementDomain string) *AndromedaAkamaiCollector {
	andromedaAkamaiCollector := AndromedaAkamaiCollector{
		session:          NewCachedAkamaiSession(session, managementDomain),
		rpc:              NewCachedRPCClient(client),
		managementDomain: managementDomain,
		fqName:           "akamai",
	}

	return &andromedaAkamaiCollector
}

func (p *AndromedaAkamaiCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc(p.fqName, "Andromeda Akamai Exporter", nil, nil)
}

// Collect fetches metrics from Akamai API and exposes them to Prometheus.
// Since Akamai provides historical data (5-minute aggregated), we handle this by:
// 1. Fetching all available data points from the API
// 2. Keeping only the most recent value for each metric combination
// 3. Exposing these as current gauge values (without historical timestamps)
// 4. Including a separate metric to track when the data was collected by Akamai
//
// This approach avoids issues with Prometheus rejecting old timestamps and creating
// gaps in the data, while still providing visibility into data freshness.
func (p *AndromedaAkamaiCollector) Collect(ch chan<- prometheus.Metric) {
	properties, err := p.session.getProperties()
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Track the most recent data point for each metric combination
	type metricKey struct {
		property   string
		datacenter string
		projectID  string
		target     string
	}
	
	latestRequests := make(map[metricKey]float64)
	latestStatus := make(map[metricKey]float64)
	latestTimestamp := make(map[metricKey]float64)

	for _, property := range properties {
		datarows, err := p.session.getTrafficReport(property)
		if err != nil {
			log.Error(err.Error())
			continue
		}

		// Process all data rows but only keep the most recent values
		for _, dataRow := range datarows {
			// Check if datacenters is empty to prevent panic
			if len(dataRow.Datacenters) == 0 {
				log.Warnf("Property %s has no datacenters in dataRow, skipping", property)
				continue
			}
			
			var projectID string
			if projectID, err = p.rpc.GetProject(dataRow.Datacenters[0].Nickname); err != nil {
				log.Error(err.Error())
				continue
			}

			for _, datacenter := range dataRow.Datacenters {
				target := strings.Split(datacenter.TrafficTargetName, " - ")[1]
				key := metricKey{
					property:   property,
					datacenter: datacenter.Nickname,
					projectID:  projectID,
					target:     target,
				}

				// Update with the latest values (datarows are ordered by time)
				latestRequests[key] = float64(datacenter.Requests)
				status, _ := strconv.Atoi(datacenter.Status)
				latestStatus[key] = float64(status)
				latestTimestamp[key] = float64(dataRow.Timestamp.Unix())
			}
		}
	}

	// Emit the most recent values as current metrics
	for key, requests := range latestRequests {
		ch <- prometheus.MustNewConstMetric(descNumReq, prometheus.GaugeValue,
			requests, key.property, key.datacenter, key.projectID, key.target)
		
		if status, ok := latestStatus[key]; ok {
			ch <- prometheus.MustNewConstMetric(descStatus, prometheus.GaugeValue,
				status, key.property, key.datacenter, key.projectID, key.target)
		}
		
		if timestamp, ok := latestTimestamp[key]; ok {
			ch <- prometheus.MustNewConstMetric(descUpdated, prometheus.GaugeValue,
				timestamp, key.property, key.datacenter, key.projectID, key.target)
		}
	}
}
