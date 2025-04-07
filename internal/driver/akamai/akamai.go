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

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	"github.com/apex/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/sapcc/andromeda/internal/config"
)

var (
	counter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "akamai_edgegrid_requests_total",
			Help: "Total number of Akamai EdgeGrid API requests.",
		},
		[]string{"method", "code"},
	)
)

func init() {
	prometheus.MustRegister(counter)
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
		session.WithClient(&http.Client{
			Transport: promhttp.InstrumentRoundTripperCounter(counter, http.DefaultTransport),
		}),
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

type AkamaiCIDRBlock struct {
	CidrId        int    `json:"cidrId,omitempty"`
	ServiceId     int    `json:"serviceId,omitempty"`
	ServiceName   string `json:"serviceName,omitempty"`
	Description   string `json:"description,omitempty"`
	Cidr          string `json:"cidr,omitempty"`
	CidrMask      string `json:"cidrMask,omitempty"`
	Port          string `json:"port,omitempty"`
	CreationDate  string `json:"creationDate,omitempty"`
	EffectiveDate string `json:"effectiveDate,omitempty"`
	ChangeDate    string `json:"changeDate,omitempty"`
	MinIp         string `json:"minIp,omitempty"`
	MaxIp         string `json:"maxIp,omitempty"`
	LastAction    string `json:"lastAction,omitempty"`
}
