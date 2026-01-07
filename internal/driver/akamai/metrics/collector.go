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
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/session"
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

// Akamai metrics:
//
// * Prometheus gauges derived from Akamai "Availability per property" endpoint
// * Prometheus counter derived from Akamai "Report traffic per property" endpoint
// * Helper metrics for both Andromeda operators as well as customers
//   to understand the freshness of the Akamai reporting

var availabilityAliveGauge = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "andromeda_akamai_availability_alive",
		Help: "Indicates whether GTM considered the server alive",
	},
	[]string{"domain", "datacenter_id", "project_id", "target_ip"},
)

var availabilityHandedOutGauge = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "andromeda_akamai_availability_handed_out",
		Help: "Indicates whether the server was handed out",
	},
	[]string{"domain", "datacenter_id", "project_id", "target_ip"},
)

var availabilityScoreGauge = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "andromeda_akamai_availability_score",
		Help: "The server score according to GTM Liveness Tests from all Web agents",
	},
	[]string{"domain", "datacenter_id", "project_id", "target_ip"},
)

var availabilityLastSyncGauge = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "andromeda_akamai_availability_last_sync",
		Help: "Last time the respective 'availability' metrics were updated by Andromeda (Unix timestamp)",
	},
	[]string{"domain", "datacenter_id", "project_id", "target_ip"},
)

var availabilitySyncErrorsCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "andromeda_akamai_availability_sync_errors",
		Help: "Total failures to retrieve 'availability' metrics from Akamai 'availability per property' API ('domain' is the only resolution possible)",
	},
	[]string{"domain"},
)

var availabilityLastReportPeriodGauge = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "andromeda_akamai_availability_last_report_period",
		Help: "Latest 5-minute time interval for which Akamai could report the respective metrics (Unix timestamp)",
	},
	[]string{"domain", "datacenter_id", "project_id", "target_ip"},
)

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
			AkamaiPropertyAvailabilityMetricsSync(akamaiSession, rpcClient, property)
			AkamaiPropertyTrafficMetricsSync(akamaiSession, rpcClient, property)
		})
	}
	wg.Wait()

	log.Infof("[AkamaiCustomMetricsSync] Sync completed in %s", time.Since(syncStart))
}

func AkamaiPropertyAvailabilityMetricsSync(akamaiSession *CachedAkamaiSession, rpcClient *CachedRPCClient, property string) {
	datarows, err := akamaiSession.getAvailabilityReport(property)
	if err != nil {
		log.Errorf("Failed to fetch AVAILABILITY reports for property %s: %v", property, err.Error())
		availabilitySyncErrorsCounter.WithLabelValues(property).Inc()
		return
	}

	if len(datarows) == 0 {
		log.Debugf("Skipping empty AVAILABILITY report (zero entries) for property %s", property)
		return
	}

	// by dropping all but the latest report, we ensure
	// the same reported count isn't re-added to the respective
	// Prometheus counter in future iterations of propertyCh
	dataRow := datarows[len(datarows)-1]

	if len(dataRow.Datacenters) == 0 {
		log.Warnf("No datacenters in AVAILABILITY report for property %s", property)
		return
	}

	var projectID string
	if projectID, err = rpcClient.GetProject(dataRow.Datacenters[0].Nickname); err != nil {
		log.Errorf("%s Failed to extract AVAILABILITY report project ID for property %s: %v", property, err.Error())
		return
	}

	boolToFloat64 := func(b bool) float64 {
		if b {
			return 1.0
		}
		return 0.0
	}

	for _, datacenter := range dataRow.Datacenters {
		if len(datacenter.IPs) == 0 {
			log.Errorf("incomplete datarow for property %s: datacenter.IPs has no entries", property)
			return
		}

		targetReport := datacenter.IPs[0]

		log.Debugf("AVAILABILITY [PROPERTY = %s | DC = %s | PROJECT = %s | TARGET = %s | TIMESTAMP = %s | "+
			"ALIVE = %v | HANDED OUT = %v | SCORE = %f ]",
			property, datacenter.Nickname,
			projectID, targetReport.IP, dataRow.Timestamp.Format(time.RFC3339),
			targetReport.Alive, targetReport.HandedOut, targetReport.Score)

		availabilityAliveGauge.
			WithLabelValues(property, datacenter.Nickname, projectID, targetReport.IP).
			Set(boolToFloat64(targetReport.Alive))
		availabilityHandedOutGauge.
			WithLabelValues(property, datacenter.Nickname, projectID, targetReport.IP).
			Set(boolToFloat64(targetReport.HandedOut))
		availabilityScoreGauge.
			WithLabelValues(property, datacenter.Nickname, projectID, targetReport.IP).
			Set(targetReport.Score)
		availabilityLastSyncGauge.
			WithLabelValues(property, datacenter.Nickname, projectID, targetReport.IP).
			Set(float64(time.Now().Unix()))
		availabilityLastReportPeriodGauge.
			WithLabelValues(property, datacenter.Nickname, projectID, targetReport.IP).
			Set(float64(dataRow.Timestamp.Unix()))
	}
}

func AkamaiPropertyTrafficMetricsSync(akamaiSession *CachedAkamaiSession, rpcClient *CachedRPCClient, property string) {
	datarows, err := akamaiSession.getTrafficReport(property)
	if err != nil {
		log.Errorf("Failed to fetch TRAFFIC reports for property %s: %v", property, err.Error())
		requestsSyncErrorsCounter.WithLabelValues(property).Inc()
		return
	}

	if len(datarows) == 0 {
		log.Debugf("Skipping empty TRAFFIC report (zero entries) for property %s", property)
		return
	}

	// by dropping all but the latest report, we ensure
	// the same reported count isn't re-added to the respective
	// Prometheus counter in future iterations of propertyCh
	dataRow := datarows[len(datarows)-1]

	if len(dataRow.Datacenters) == 0 {
		log.Warnf("No datacenters in TRAFFIC report for property %s", property)
		return
	}

	var projectID string
	if projectID, err = rpcClient.GetProject(dataRow.Datacenters[0].Nickname); err != nil {
		log.Errorf("%s Failed to extract TRAFFIC report project ID for property %s: %v", property, err.Error())
		return
	}
	for _, datacenter := range dataRow.Datacenters {
		target := strings.Split(datacenter.TrafficTargetName, " - ")[1]

		log.Debugf("TRAFFIC [PROPERTY = %s | DC = %s | PROJECT = %s | TARGET = %s | TIMESTAMP = %s | REQUESTS = %d ]",
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
