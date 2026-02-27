// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package akamai

import (
	"net/http"
	"os"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/session"
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

type AkamaiCIDRBlock map[string]any
