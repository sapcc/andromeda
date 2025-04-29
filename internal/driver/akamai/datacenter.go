// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package akamai

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/gtm"
	"github.com/apex/log"

	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/driver"
	"github.com/sapcc/andromeda/internal/rpc/server"
	"github.com/sapcc/andromeda/internal/rpcmodels"
	"github.com/sapcc/andromeda/internal/utils"
	"github.com/sapcc/andromeda/models"
)

var errDatacenterNotFound = errors.New("datacenter/GetDatacenterMeta: datacenter not found")

// GetDatacenterMeta return the meta id for a given datacenter uuid
func (s *AkamaiAgent) GetDatacenterMeta(datacenterUUID string, datacenters []*rpcmodels.Datacenter) (int, error) {
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

	var meta int
	for _, datacenter := range datacenters {
		if datacenter.GetId() == datacenterUUID {
			meta = int(datacenter.GetMeta())
			if meta != 0 {
				return meta, nil
			}
		}
	}

	// Refresh datacenter meta id from akamai
	request := gtm.ListDatacentersRequest{DomainName: config.Global.AkamaiConfig.Domain}
	dcs, err := s.gtm.ListDatacenters(context.Background(), request)
	if err != nil {
		return 0, nil
	}
	for _, d := range dcs {
		if d.Nickname == datacenterUUID {
			meta = d.DatacenterID
			break
		}
	}
	if meta == 0 {
		return 0, errDatacenterNotFound
	}

	// try syncing meta id with andromeda database
	req := &server.DatacenterMetaRequest{Id: datacenterUUID, Meta: int32(meta)}
	if _, err = s.rpc.UpdateDatacenterMeta(context.Background(), req); err != nil {
		log.Errorf("UpdateDatacenterMeta(%s) failed: %v", datacenterUUID, err)
	}

	// cache meta id
	s.datacenterIdCache.Add(datacenterUUID, meta)
	return meta, nil
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
	log.Debugf("Running FetchAndSyncDatacenters(%+v, force=%t)", datacenters, force)

	res, err := s.GetDatacenters(datacenters)
	if err != nil {
		return err
	}
	var provRequests []*server.ProvisioningStatusRequest_ProvisioningStatus
	for _, datacenter := range res {
		if datacenter.ProvisioningStatus == models.DatacenterProvisioningStatusPENDINGDELETE {
			provRequests = append(provRequests,
				driver.GetProvisioningStatusRequest(datacenter.Id, "DATACENTER", "DELETED"))
		} else if datacenter.ProvisioningStatus != models.DatacenterProvisioningStatusACTIVE {
			provRequests = append(provRequests,
				driver.GetProvisioningStatusRequest(datacenter.Id, "DATACENTER",
					models.DatacenterProvisioningStatusACTIVE))
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

func (s *AkamaiAgent) SyncDatacenter(datacenter *rpcmodels.Datacenter, force bool) (*rpcmodels.Datacenter, error) {
	log.Debugf("SyncDatacenter(%s, force=%t)", datacenter.Id, force)

	// akamai datacenterId is a unique numeric reference to a domain specific datacenter
	var meta = int(datacenter.GetMeta())
	var backendDatacenter *gtm.Datacenter
	var err error

	if meta == 0 {
		meta, err = s.GetDatacenterMeta(datacenter.Id, nil)
		if err != nil && !errors.Is(err, errDatacenterNotFound) {
			return nil, err
		}
	}
	if meta != 0 {
		request := gtm.GetDatacenterRequest{
			DatacenterID: meta,
			DomainName:   config.Global.AkamaiConfig.Domain,
		}
		if backendDatacenter, err = s.gtm.GetDatacenter(context.Background(), request); err != nil {
			// check if datacenter is not found
			var gtmErr *gtm.Error
			if errors.As(err, &gtmErr) && gtmErr.StatusCode != 404 {
				return nil, err
			}
		}
	}

	if datacenter.ProvisioningStatus == models.DatacenterProvisioningStatusPENDINGDELETE {
		// nothing to delete?
		if backendDatacenter == nil {
			// datacenter already deleted
			return nil, nil
		}

		// run Delete
		request := gtm.DeleteDatacenterRequest{
			DatacenterID: backendDatacenter.DatacenterID,
			DomainName:   config.Global.AkamaiConfig.Domain,
		}
		if _, err = s.gtm.DeleteDatacenter(context.Background(), request); err != nil {
			return nil, fmt.Errorf("DeleteDatacenter(%s) failed: %w", datacenter.Id, err)
		}
		return nil, nil
	}

	// consider synced
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
	if utils.DeepEqualFields(&referenceDatacenter, backendDatacenter, fieldsToCompare) {
		// no change
		return datacenter, nil
	}

	if backendDatacenter != nil {
		// Run Update
		referenceDatacenter.DatacenterID = backendDatacenter.DatacenterID
		request := gtm.UpdateDatacenterRequest{
			Datacenter: &referenceDatacenter,
			DomainName: config.Global.AkamaiConfig.Domain,
		}
		_, err = s.gtm.UpdateDatacenter(context.Background(), request)
		if err != nil {
			log.Errorf("UpdateDatacenter(%s) for domain %s failed", referenceDatacenter.Nickname,
				config.Global.AkamaiConfig.Domain)
			return nil, err
		} else {
			log.Infof("UpdateDatacenter(%s) for domain %s", referenceDatacenter.Nickname,
				config.Global.AkamaiConfig.Domain)
		}
	} else {
		var res *gtm.CreateDatacenterResponse
		request := gtm.CreateDatacenterRequest{
			Datacenter: &referenceDatacenter,
			DomainName: config.Global.AkamaiConfig.Domain,
		}
		res, err = s.gtm.CreateDatacenter(context.Background(), request)
		if err != nil {
			log.Errorf("CreateDatacenter(%s) for domain %s failed", referenceDatacenter.Nickname,
				config.Global.AkamaiConfig.Domain)
			return nil, err
		} else {
			log.Infof("CreateDatacenter(%s) for domain %s", referenceDatacenter.Nickname,
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
