/*
 *   Copyright 2022 SAP SE
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

package akamai

import (
	"net/http"
	"os"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	"github.com/apex/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/sapcc/andromeda/internal/config"
)

var (
	edgegridAPICalls = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "akamai_edgegrid_api_calls_total",
			Help: "Total number of Akamai EdgeGrid API calls",
		},
		[]string{"method", "status_code"},
	)

	edgegridErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "akamai_edgegrid_errors_total",
			Help: "Total number of Akamai EdgeGrid errors",
		},
	)
)

func init() {
	prometheus.MustRegister(edgegridAPICalls, edgegridErrors)
}

func instrumentedRoundTripper(next http.RoundTripper) http.RoundTripper {
	return promhttp.RoundTripperFunc(func(req *http.Request) (resp *http.Response, err error) {
		defer func() {
			if err != nil {
				edgegridErrors.Inc()
			}
		}()
		resp, err = next.RoundTrip(req)
		edgegridAPICalls.WithLabelValues(req.Method, strconv.Itoa(resp.StatusCode)).Inc()
		return resp, err
	})
}

func NewAkamaiSession(akamaiConfig *config.AkamaiConfig) (*session.Session, string) {
	option := edgegrid.WithEnv(true)
	if env := os.Getenv("AKAMAI_EDGE_RC"); env != "" {
		option = edgegrid.WithFile(env)
	} else if akamaiConfig.EdgeRC != "" {
		option = edgegrid.WithFile(akamaiConfig.EdgeRC)
	}

	edgerc := edgegrid.Must(edgegrid.New(option))
	s := session.Must(session.New(
		session.WithSigner(edgerc),
		session.WithClient(&http.Client{Transport: instrumentedRoundTripper(http.DefaultTransport)}),
	))

	var identity struct {
		AccountID string `json:"accountId"`
		Active    bool   `json:"active"`
		Contracts []struct {
			ContractID  string   `json:"contractId"`
			Features    []string `json:"features"`
			Permissions []string `json:"permissions"`
		} `json:"contracts"`
		Email string `json:"email"`
	}

	req, _ := http.NewRequest(http.MethodGet, "/config-gtm/v1/identity", nil)
	if _, err := s.Exec(req, &identity); err != nil {
		panic(err)
	}

	if len(identity.Contracts) != 1 && akamaiConfig.ContractId == "" {
		log.Fatalf("More than one contract detected, specificy contract_id.")
	}

	var domainType string
	for _, contract := range identity.Contracts {
		if akamaiConfig.ContractId != "" && contract.ContractID != akamaiConfig.ContractId {
			continue
		}

		domainType = DetectTypeFromFeatures(contract.Features)
		log.Infof("Detected Akamai Contract '%s' with best features enabling '%s' domain type.",
			contract.ContractID, domainType)
		break
	}
	return &s, domainType
}
