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
	"github.com/sapcc/andromeda/models"
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
	var provRequests []*server.ProvisioningStatusRequest_ProvisioningStatus
	for _, datacenter := range res {
		if datacenter.ProvisioningStatus == models.DatacenterProvisioningStatusPENDINGDELETE {
			provRequests = append(provRequests,
				driver.GetProvisioningStatusRequest(datacenter.Id, "DATACENTER", "DELETED"))
		} else {
			provRequests = append(provRequests,
				driver.GetProvisioningStatusRequest(datacenter.Id, "DATACENTER", "ACTIVE"))
		}

		if _, err = s.SyncDatacenter(datacenter, force); err != nil {
			return err
		}
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
		driver.UpdateProvisioningStatus(s.rpc, provRequests)
	}

	return nil
}

func getBackendDatacenter(meta int, datacenter *rpcmodels.Datacenter, s *AkamaiAgent) *gtm.Datacenter {
	if meta != 0 {
		if backendDatacenter, err := s.gtm.GetDatacenter(context.Background(), meta, config.Global.AkamaiConfig.Domain); err == nil {
			return backendDatacenter
		}
	}

	// try to find datacenter by nickname
	datacenters, err := s.gtm.ListDatacenters(context.Background(), config.Global.AkamaiConfig.Domain)
	if err != nil {
		return nil
	}
	for _, d := range datacenters {
		if d.Nickname == datacenter.Id {
			return d
		}
	}
	return nil
}

func (s *AkamaiAgent) SyncDatacenter(datacenter *rpcmodels.Datacenter, force bool) (*rpcmodels.Datacenter, error) {
	logger.Debugf("SyncDatacenter(%s, force=%t)", datacenter.Id, force)

	// akamai datacenterId is a unique numeric reference to a domain specific datacenter
	meta := int(datacenter.GetMeta())

	if datacenter.ProvisioningStatus == models.DatacenterProvisioningStatusPENDINGDELETE {
		// Run Delete
		toDelete := getBackendDatacenter(meta, datacenter, s)
		if toDelete == nil {
			// datacenter already deleted
			return nil, nil
		}

		_, err := s.gtm.DeleteDatacenter(context.Background(), toDelete, config.Global.AkamaiConfig.Domain)
		return nil, err
	}

	// Consider synced
	if !force && meta != 0 {
		return datacenter, nil
	}

	// reference datacenter
	referenceDatacenter := gtm.Datacenter{
		City:            datacenter.GetCity(),
		Continent:       datacenter.GetContinent(),
		Country:         datacenter.GetCountry(),
		StateOrProvince: datacenter.GetStateOrProvince(),
		Latitude:        datacenter.GetLatitude(),
		Longitude:       datacenter.GetLongitude(),
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
	backendDatacenter := getBackendDatacenter(meta, datacenter, s)
	if utils.DeepEqualFields(&referenceDatacenter, backendDatacenter, fieldsToCompare) {
		// no change
		return datacenter, nil
	}

	if backendDatacenter != nil {
		// Run Update
		referenceDatacenter.DatacenterID = backendDatacenter.DatacenterID
		_, err := s.gtm.UpdateDatacenter(context.Background(), &referenceDatacenter, config.Global.AkamaiConfig.Domain)
		if err != nil {
			logger.Errorf("UpdateDatacenter(%s) for domain %s failed", referenceDatacenter.Nickname,
				config.Global.AkamaiConfig.Domain)
			return nil, err
		} else {
			logger.Infof("UpdateDatacenter(%s) for domain %s", referenceDatacenter.Nickname,
				config.Global.AkamaiConfig.Domain)
		}
	} else {
		res, err := s.gtm.CreateDatacenter(context.Background(), &referenceDatacenter, config.Global.AkamaiConfig.Domain)
		if err != nil {
			logger.Errorf("CreateDatacenter(%s) for domain %s failed", referenceDatacenter.Nickname,
				config.Global.AkamaiConfig.Domain)
			return nil, err
		} else {
			logger.Infof("CreateDatacenter(%s) for domain %s", referenceDatacenter.Nickname,
				config.Global.AkamaiConfig.Domain)
		}
		backendDatacenter = res.Resource
	}
	// update backend datacenter

	if backendDatacenter.DatacenterID != meta {
		req := &server.DatacenterMetaRequest{Id: datacenter.Id, Meta: int32(backendDatacenter.DatacenterID)}
		return s.rpc.UpdateDatacenterMeta(context.Background(), req)
	}
	return datacenter, nil
}
