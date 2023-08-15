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

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/session"
	"go-micro.dev/v4/logger"

	"github.com/sapcc/andromeda/internal/config"
)

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
		logger.Fatalf("More than one contract detected, specificy contract_id.")
	}

	var domainType string
	for _, contract := range identity.Contracts {
		if akamaiConfig.ContractId != "" && contract.ContractID != akamaiConfig.ContractId {
			continue
		}

		domainType = DetectTypeFromFeatures(contract.Features)
		logger.Infof("Detected Akamai Contract '%s' with best features enabling '%s' domain type.",
			contract.ContractID, domainType)
		break
	}
	return &s, domainType
}
