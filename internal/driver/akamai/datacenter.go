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
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v5/pkg/gtm"
	"go-micro.dev/v4/logger"

	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/driver"
	"github.com/sapcc/andromeda/internal/rpc/server"
	"github.com/sapcc/andromeda/internal/rpcmodels"
)

func (s *AkamaiAgent) GetDatacenterReference(datacenterUUID string, datacenters []*rpcmodels.Datacenter) (int, error) {
	// Fetch from cache
	if id, ok := s.datacenterIdCache.Get(datacenterUUID); ok {
		return id, nil
	}

	if datacenters == nil {
		var err error
		if datacenters, err = s.GetDatacenters([]string{datacenterUUID}); err != nil {
			return 0, err
		}
	}

	var aDatacenter *rpcmodels.Datacenter
	for _, datacenter := range datacenters {
		if datacenter.GetId() == datacenterUUID {
			aDatacenter = datacenter
			break
		}
	}

	// DatacenterId is a unique number for an akamai datacenter
	s.datacenterIdCache.Add(datacenterUUID, int(aDatacenter.GetMeta()))
	return int(aDatacenter.GetMeta()), nil
}

func (s *AkamaiAgent) GetDatacenters(datacenterIDs []string) ([]*rpcmodels.Datacenter, error) {
	response, err := s.rpc.GetDatacenters(context.Background(), &server.SearchRequest{
		Provider:       "akamai",
		PageNumber:     0,
		ResultPerPage:  1,
		FullyPopulated: true,
		Ids:            datacenterIDs,
	})
	if err != nil {
		return nil, err
	}

	return response.GetResponse(), nil
}

func (s *AkamaiAgent) uploadDatacenter(datacenter *rpcmodels.Datacenter) (*gtm.Datacenter, error) {
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

func (s *AkamaiAgent) FetchAndSyncDatacenters(datacenters []string) error {
	logger.Debugf("Running FetchAndSyncDatacenters()")

	res, err := s.GetDatacenters(datacenters)
	if err != nil {
		return err
	}
	for _, datacenter := range res {
		if _, err = s.SyncDatacenter(datacenter, false); err != nil {
			return err
		}

		// Wait for status propagation
		var status string
		for ok := true; ok; ok = status == "PENDING" {
			time.Sleep(5 * time.Second)
			status, err = s.syncProvisioningStatus(nil)
			if err != nil {
				return err
			}
		}

		if status == "COMPLETE" {
			driver.UpdateProvisioningStatus(s.rpc,
				[]*server.ProvisioningStatusRequest_ProvisioningStatus{
					driver.GetProvisioningStatusRequest(datacenter.Id, "DATACENTER", "ACTIVE"),
				})
		}
	}
	return nil
}

func (s *AkamaiAgent) SyncDatacenter(datacenter *rpcmodels.Datacenter, force bool) (*rpcmodels.Datacenter, error) {
	logger.Debugf("SyncDatacenter('%s')", datacenter.Id)

	// akamai datacenterId is a unique numeric reference to a domain specific datacenter
	meta := int(datacenter.GetMeta())

	// Consider synced
	if !force && meta != 0 {
		return datacenter, nil
	}

	backendDatacenter, err := s.uploadDatacenter(datacenter)
	if err != nil {
		return nil, err
	}

	if backendDatacenter.DatacenterId != meta {
		req := &server.DatacenterMetaRequest{Id: datacenter.Id, Meta: int32(backendDatacenter.DatacenterId)}
		return s.rpc.UpdateDatacenterMeta(context.Background(), req)
	}
	return datacenter, nil
}
