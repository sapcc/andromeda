// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"time"

	"github.com/apex/log"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	prometheus.MustRegister(providerEntitiesGauge)
}

var providerEntitiesGauge = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "andromeda_provider_entities",
		Help: "Gauge of total DB entities",
	},
	[]string{"provider", "entity"},
)

var providers = []string{"akamai", "f5"}

func ProviderEntitiesMetricsWorker(interval time.Duration, db *sqlx.DB, collector func(db *sqlx.DB)) {
	collector(db)
	c := time.Tick(interval)
	for {
		<-c
		collector(db)
	}
}

func CollectProviderEntitiesMetrics(db *sqlx.DB) {
	for _, provider := range providers {
		CollectDomainsByProvider(provider, db)
	}
}

func CollectDomainsByProvider(provider string, db *sqlx.DB) {
	var count int
	query := `SELECT COUNT(*) FROM domain WHERE provider = ?`
	if err := db.Get(&count, query, provider); err != nil {
		log.Errorf("CollectDomainsByProvider(%s) failed: %v", provider, err)
		return
	}
	providerEntitiesGauge.WithLabelValues(provider, "domain").Set(float64(count))
}
