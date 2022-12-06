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
	"context"

	"go-micro.dev/v4/logger"

	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/rpcmodels"
)

func (s *AkamaiAgent) syncStatus(domain *rpcmodels.Domain) (string, error) {
	// Check for running domain's propagation state
	status, err := s.gtm.GetDomainStatus(context.Background(), config.Global.AkamaiConfig.Domain)
	if err != nil {
		return "UNKNOWN", err
	}

	// Tracks the status of the domain's propagation state. Either PENDING, COMPLETE, or DENIED.
	// A DENIED value indicates that the domain configuration is invalid,
	// and doesn't propagate until the validation errors are resolved.
	switch status.PropagationStatus {
	case "PENDING":
		logger.Debug("Akamai Backend: pending configuration change")
	case "DENIED":
		logger.Errorf("Domain %s failed syncing: %s", domain.Id, status.Message)
		if err := s.UpdateDomainProvisioningStatus(domain, "ERROR"); err != nil {
			return "UNKNOWN", err
		}
	case "COMPLETE":
		logger.Infof("Domain %s is in sync", domain.Id)
		if err := s.UpdateDomainProvisioningStatus(domain, "ACTIVE"); err != nil {
			return "UNKNOWN", err
		}
	}
	return status.PropagationStatus, nil
}
