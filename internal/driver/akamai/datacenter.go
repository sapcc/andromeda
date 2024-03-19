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

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/gtm"
	"go-micro.dev/v4/logger"

	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/driver"
	"github.com/sapcc/andromeda/internal/rpc/server"
	"github.com/sapcc/andromeda/internal/rpcmodels"
	"github.com/sapcc/andromeda/internal/utils"
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

func (s *AkamaiAgent) FetchAndSyncDatacenters(datacenters []string, force bool) error {
	logger.Debugf("Running FetchAndSyncDatacenters(%+v, force=%t)", datacenters, force)

	res, err := s.GetDatacenters(datacenters)
	if err != nil {
		return err
	}
	for _, datacenter := range res {
		if _, err = s.SyncDatacenter(datacenter, force); err != nil {
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
	logger.Debugf("SyncDatacenter(%s, force=%t)", datacenter.Id, force)

	// akamai datacenterId is a unique numeric reference to a domain specific datacenter
	meta := int(datacenter.GetMeta())

	// Consider synced
	if !force && meta != 0 {
		return datacenter, nil
	}

	var backendDatacenter *gtm.Datacenter
	if meta != 0 {
		var err error

		backendDatacenter, err = s.gtm.GetDatacenter(context.Background(), meta, config.Global.AkamaiConfig.Domain)
		if err != nil {
			// try to find datacenter by nickname
			datacenters, err := s.gtm.ListDatacenters(context.Background(), config.Global.AkamaiConfig.Domain)
			if err != nil {
				return nil, err
			}

			for _, d := range datacenters {
				if d.Nickname == datacenter.Id {
					backendDatacenter = d
					break
				}
			}
		}
	}

	// reference datacenter
	referenceDatacenter := gtm.Datacenter{
		City:            datacenter.GetCity(),
		Continent:       datacenter.GetContinent(),
		Country:         datacenter.GetCountry(),
		StateOrProvince: datacenter.GetStateOrProvince(),
		Latitude:        float64(datacenter.GetLatitude()),
		Longitude:       float64(datacenter.GetLongitude()),
		Nickname:        datacenter.Id,
	}

	// compare to reference datacenter
	fieldsToCompare := []string{
		"City",
		"Continent",
		"Country",
		"StateOrProvince",
		"Longitude",
		"Nickname",
	}
	if utils.DeepEqualFields(&referenceDatacenter, backendDatacenter, fieldsToCompare) {
		// no change
		return datacenter, nil
	}

	res, err := s.gtm.CreateDatacenter(context.Background(), &referenceDatacenter, config.Global.AkamaiConfig.Domain)
	if err != nil {
		logger.Errorf("CreateDatacenter(%s) for domain %s failed", referenceDatacenter.Nickname,
			config.Global.AkamaiConfig.Domain)
		return nil, err
	} else {
		logger.Infof("CreateDatacenter(%s) for domain %s", referenceDatacenter.Nickname,
			config.Global.AkamaiConfig.Domain)
	}

	// update backend datacenter
	backendDatacenter = res.Resource

	if backendDatacenter.DatacenterId != meta {
		req := &server.DatacenterMetaRequest{Id: datacenter.Id, Meta: int32(backendDatacenter.DatacenterId)}
		return s.rpc.UpdateDatacenterMeta(context.Background(), req)
	}
	return datacenter, nil
}
