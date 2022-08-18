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
	"github.com/sapcc/andromeda/internal/models"
)

var MONITOR_LIVENESS_TYPE_MAP = map[models.Monitor_MonitorType]string{
	models.Monitor_HTTP:  "HTTP",
	models.Monitor_HTTPS: "HTTPS",
	models.Monitor_TCP:   "TCP",
}

func (s *AkamaiAgent) EnsureDomain() error {
	if _, err := s.gtm.GetDomain(context.Background(), config.Global.AkamaiConfig.Domain); err != nil {
		logger.Warnf("Akamai Domain %s doesn't exist, creating...")
		domain := s.gtm.NewDomain(context.Background(),
			config.Global.AkamaiConfig.Domain, config.Global.AkamaiConfig.DomainType)
		if _, err := s.gtm.CreateDomain(context.Background(), domain, nil); err != nil {
			return err
		}
	}
	return nil
}

func (s *AkamaiAgent) SyncDomain(domain *models.Domain) error {
	/*
		fieldsToCompare := []string{
			"Type",
			"Datacenters",
			"Datacenters.City",
			"Datacenters.Continent",
			"Datacenters.Country"}
		if !utils.DeepEqualFields(expectedDomain, akamaiDomain, fieldsToCompare) {
			// Update
			if _, err := expectedDomain.Update(nil); err != nil {
				return err
			}
		}
	*/
	return nil
}
