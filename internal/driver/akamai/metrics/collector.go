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
