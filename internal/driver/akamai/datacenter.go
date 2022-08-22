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
	"fmt"

	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/configgtm"
	"go-micro.dev/v4/logger"

	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/models"
	"github.com/sapcc/andromeda/internal/rpc/server"
)

func (s *AkamaiAgent) GetDatacenter(datacenterID string) (*models.Datacenter, error) {
	response, err := s.rpc.GetDatacenters(context.Background(), &server.SearchRequest{
		Provider:      "akamai",
		PageNumber:    0,
		ResultPerPage: 1,
		Ids:           []string{datacenterID},
	})
	if err != nil {
		return nil, err
	}

	res := response.GetResponse()
	if len(res) != 1 {
		return nil, fmt.Errorf("Failed fetching datacenter '%s': len(res) = %d != 1", datacenterID,
			len(res))
	}
	return res[0], nil
}

func (s *AkamaiAgent) fetchOrCreateDatacenter(datacenter *models.Datacenter) (*gtm.Datacenter, error) {
	datacenters, err := s.gtm.ListDatacenters(context.Background(), config.Global.AkamaiConfig.Domain)
	if err != nil {
		return nil, err
	}

	for _, d := range datacenters {
		if d.Nickname == datacenter.Id {
			return d, nil
		}
	}

	// Create datacenter
	akamaiDatacenter := gtm.Datacenter{
		City:            datacenter.GetCity(),
		Continent:       datacenter.GetContinent(),
		Country:         datacenter.GetCountry(),
		StateOrProvince: datacenter.GetStateOrProvince(),
		Latitude:        float64(datacenter.GetLatitude()),
		Longitude:       float64(datacenter.GetLongitude()),
		Nickname:        datacenter.Id,
	}

	res, err := s.gtm.CreateDatacenter(context.Background(), &akamaiDatacenter, config.Global.AkamaiConfig.Domain)
	if err != nil {
		logger.Errorf("CreateDatacenter(%s) for domain %s failed", akamaiDatacenter.Nickname,
			config.Global.AkamaiConfig.Domain)
		return nil, err
	} else {
		logger.Infof("CreateDatacenter(%s) for domain %s", akamaiDatacenter.Nickname,
			config.Global.AkamaiConfig.Domain)
	}
	return res.Resource, nil
}

func (s *AkamaiAgent) SyncDatacenter(datacenter *models.Datacenter, force bool) (*models.Datacenter, error) {
	logger.Debugf("SyncDatacenter('%s')", datacenter.Id)

	// akamai datacenterId is a unique numeric reference to a domain specific datacenter
	meta := int(datacenter.GetMeta())

	// Consider synced
	if !force && meta != 0 {
		return datacenter, nil
	}

	backendDatacenter, err := s.fetchOrCreateDatacenter(datacenter)
	if err != nil {
		return nil, err
	}

	if backendDatacenter.DatacenterId != meta {
		req := &server.DatacenterMetaRequest{Id: datacenter.Id, Meta: int32(backendDatacenter.DatacenterId)}
		return s.rpc.UpdateDatacenterMeta(context.Background(), req)
	}
	return datacenter, nil
}
