// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package metrics

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/actatum/stormrpc"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	"github.com/apex/log"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	fQNameNumReq = prometheus.BuildFQName("andromeda", "akamai", "requests_5m")
	fQNameStatus = prometheus.BuildFQName("andromeda", "akamai", "status_5m")
	labels       = []string{"domain", "datacenter_id", "project_id", "target_ip"}
	descNumReq   = prometheus.NewDesc(fQNameNumReq, "Number of requests/5m per domain", labels, nil)
	descStatus   = prometheus.NewDesc(fQNameStatus, "Status per domain", labels, nil)
)

type AndromedaAkamaiCollector struct {
	session          *CachedAkamaiSession
	rpc              *CachedRPCClient
	managementDomain string
	fqName           string
}

func NewAndromedaAkamaiCollector(session session.Session, client *stormrpc.Client, managementDomain string) *AndromedaAkamaiCollector {
	andromedaAkamaiCollector := AndromedaAkamaiCollector{
		session:          NewCachedAkamaiSession(session, managementDomain, NewAkamaiRateLimiter(40)),
		rpc:              NewCachedRPCClient(client),
		managementDomain: managementDomain,
		fqName:           "akamai",
	}

	return &andromedaAkamaiCollector
}

func (p *AndromedaAkamaiCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc(p.fqName, "Andromeda Akamai Exporter", nil, nil)
}

func (p *AndromedaAkamaiCollector) Collect(ch chan<- prometheus.Metric) {
	properties, err := p.session.getProperties()
	if err != nil {
		log.Error(err.Error())
		return
	}

	for _, property := range properties {
		datarows, err := p.session.getTrafficReport(property)
		if err != nil {
			log.Error(err.Error())
			continue
		}

		for _, dataRow := range datarows {
			var projectID string
			if projectID, err = p.rpc.GetProject(dataRow.Datacenters[0].Nickname); err != nil {
				log.Error(err.Error())
				continue
			}

			for _, datacenter := range dataRow.Datacenters {
				// expose via prometheus
				target := strings.Split(datacenter.TrafficTargetName, " - ")[1]
				m := prometheus.MustNewConstMetric(descNumReq, prometheus.GaugeValue,
					float64(datacenter.Requests), property, datacenter.Nickname, projectID,
					target)
				ch <- prometheus.NewMetricWithTimestamp(dataRow.Timestamp, m)

				status, _ := strconv.Atoi(datacenter.Status)
				m = prometheus.MustNewConstMetric(descStatus, prometheus.GaugeValue,
					float64(status), property, datacenter.Nickname, projectID, target)
				ch <- prometheus.NewMetricWithTimestamp(dataRow.Timestamp, m)
			}
		}
	}
}

// Experimental metrics:
//
// * A Prometheus counter derived from Akamai "Report traffic per property" endpoint
// * Helper metrics for both Andromeda operators as well as customers
//   to understand the freshness of the Akamai reporting

var requestsCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "andromeda_akamai_requests",
		Help: "Total requests per target IP (derived from Akamai reporting API)",
	},
	[]string{"domain", "datacenter_id", "project_id", "target_ip"},
)

var requestsLastSyncGauge = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "andromeda_akamai_requests_last_sync",
		Help: "Last time the respective counter metric was incremented by Andromeda (Unix timestamp)",
	},
	[]string{"domain", "datacenter_id", "project_id", "target_ip"},
)

var requestsSyncErrorsCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "andromeda_akamai_requests_sync_errors",
		Help: "Total failures to retrieve requests from Akamai reporting API ('domain' is the only resolution possible)",
	},
	[]string{"domain"},
)

var requestsLastReportPeriodGauge = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "andromeda_akamai_requests_last_report_period",
		Help: "Latest 5-minute time interval for which Akamai could report the respective metric (Unix timestamp)",
	},
	[]string{"domain", "datacenter_id", "project_id", "target_ip"},
)

var rateLimitingDurationSeconds = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "andromeda_akamai_ratelimiting_duration_seconds",
		Help: "Time that Go routines spend waiting on the rate limiter",
	},
)

func AkamaiCustomMetricsSync(akamaiSession *CachedAkamaiSession, rpcClient *CachedRPCClient) {
	syncStart := time.Now()

	properties, err := akamaiSession.getProperties()
	if err != nil {
		log.Error(err.Error())
		return
	}

	log.Debugf("[AkamaiCustomMetricsSync] Got %d properties", len(properties))

	// each go routine handles a single property, then terminates
	var wg sync.WaitGroup
	for _, property := range properties {
		wg.Go(func() {
			AkamaiPropertyMetricsSync(akamaiSession, rpcClient, property)
		})
	}
	wg.Wait()

	log.Infof("[AkamaiCustomMetricsSync] Sync completed in %s", time.Since(syncStart))
}

func AkamaiPropertyMetricsSync(akamaiSession *CachedAkamaiSession, rpcClient *CachedRPCClient, property string) {
	datarows, err := akamaiSession.getTrafficReport(property)
	if err != nil {
		log.Errorf("Failed to fetch traffic reports for property %s: %v", property, err.Error())
		requestsSyncErrorsCounter.WithLabelValues(property).Inc()
		return
	}

	if len(datarows) == 0 {
		log.Debugf("Skipping empty report (zero entries) for property %s", property)
		return
	}

	// by dropping all but the latest report, we ensure
	// the same reported count isn't re-added to the respective
	// Prometheus counter in future iterations of propertyCh
	dataRow := datarows[len(datarows)-1]

	var projectID string
	if projectID, err = rpcClient.GetProject(dataRow.Datacenters[0].Nickname); err != nil {
		log.Errorf("%s Failed to extract report project ID for property %s: %v", property, err.Error())
		return
	}
	for _, datacenter := range dataRow.Datacenters {
		target := strings.Split(datacenter.TrafficTargetName, " - ")[1]

		log.Debugf("[PROPERTY = %s | DC = %s | PROJECT = %s | TARGET = %s | TIMESTAMP = %s | REQUESTS = %d ]",
			property, datacenter.Nickname,
			projectID, target, dataRow.Timestamp.Format(time.RFC3339),
			datacenter.Requests)

		requestsCounter.
			WithLabelValues(property, datacenter.Nickname, projectID, target).
			Add(float64(datacenter.Requests))
		requestsLastSyncGauge.
			WithLabelValues(property, datacenter.Nickname, projectID, target).
			Set(float64(time.Now().Unix()))
		requestsLastReportPeriodGauge.
			WithLabelValues(property, datacenter.Nickname, projectID, target).
			Set(float64(dataRow.Timestamp.Unix()))
	}
}
