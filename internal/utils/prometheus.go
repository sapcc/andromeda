// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"fmt"
	"net/http"

	"github.com/apex/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/sapcc/andromeda/internal/config"
)

func PrometheusListen() {
	log.WithField("url", fmt.Sprintf("http://%s/metrics", config.Global.Default.PrometheusListen)).
		Info("Serving prometheus metrics")
	http.Handle("/metrics", promhttp.Handler())

	if err := http.ListenAndServe(config.Global.Default.PrometheusListen, nil); err != nil {
		log.WithError(err).Fatal("Failed to start prometheus listener")
	}
}
