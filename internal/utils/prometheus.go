/*
 *   Copyright 2024 SAP SE
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

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
