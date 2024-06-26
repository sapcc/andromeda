/*
 *   Copyright 2023 SAP SE
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

func (s *AkamaiAgent) GetGeomap(geomapID string) (*rpcmodels.Geomap, error) {
	response, err := s.rpc.GetGeomaps(context.Background(), &server.SearchRequest{
		Provider:      "akamai",
		PageNumber:    0,
		ResultPerPage: 1,
		Ids:           []string{geomapID},
	})
	if err != nil {
		return nil, err
	}

	res := response.GetResponse()
	if len(res) != 1 {
		return nil, fmt.Errorf("Failed fetching geomap '%s': len(res) = %d != 1", geomapID,
			len(res))
	}
	return res[0], nil
}

func (s *AkamaiAgent) constructAkamaiGeoMap(geomap *rpcmodels.Geomap) (*gtm.GeoMap, error) {
	defaultDatacenterReference, err := s.GetDatacenterReference(geomap.DefaultDatacenter, nil)
	if err != nil {
		return nil, err
	}

	var assignments = make([]*gtm.GeoAssignment, 0)
	for _, dc := range geomap.Assignment {
		if dc.Datacenter == geomap.DefaultDatacenter {
			// skip assignments to default datacenter
			continue
		}

		datacenterReference, err := s.GetDatacenterReference(dc.Datacenter, nil)
		if err != nil {
			return nil, err
		}

		assignments = append(assignments, &gtm.GeoAssignment{
			DatacenterBase: gtm.DatacenterBase{
				Nickname:     dc.Datacenter,
				DatacenterID: datacenterReference,
			},
			Countries: dc.Countries,
		})
	}

	// Create geomap
	akamaiGeoMap := gtm.GeoMap{
		DefaultDatacenter: &gtm.DatacenterBase{
			Nickname:     geomap.DefaultDatacenter,
			DatacenterID: defaultDatacenterReference},
		Assignments: assignments,
		Name:        geomap.Id,
	}

	return &akamaiGeoMap, nil
}

func (s *AkamaiAgent) FetchAndSyncGeomaps(geomaps []string, force bool) error {
	logger.Debugf("Running FetchAndSyncGeomaps(%+v, force=%t)", geomaps, force)
	response, err := s.rpc.GetGeomaps(context.Background(), &server.SearchRequest{
		Provider:       "akamai",
		PageNumber:     0,
		ResultPerPage:  1,
		FullyPopulated: true,
		Pending:        geomaps == nil && !force,
		Ids:            geomaps,
	})
	if err != nil {
		return err
	}

	var provRequests []*server.ProvisioningStatusRequest_ProvisioningStatus
	for _, geomap := range response.GetResponse() {
		if geomap.ProvisioningStatus == models.DatacenterProvisioningStatusPENDINGDELETE {
			provRequests = append(provRequests,
				driver.GetProvisioningStatusRequest(geomap.Id, "GEOGRAPHIC_MAP", "DELETED"))
		} else {
			provRequests = append(provRequests,
				driver.GetProvisioningStatusRequest(geomap.Id, "GEOGRAPHIC_MAP", "ACTIVE"))
		}

		if _, err = s.SyncGeomap(geomap, force); err != nil {
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

func (s *AkamaiAgent) SyncGeomap(geomap *rpcmodels.Geomap, force bool) (*gtm.GeoMap, error) {
	logger.Debugf("SyncGeomap(%s, force=%t)", geomap.Id, force)

	if geomap.ProvisioningStatus == models.GeomapProvisioningStatusPENDINGDELETE {
		// Run Delete
		toDelete, err := s.gtm.GetGeoMap(context.Background(), geomap.Id, config.Global.AkamaiConfig.Domain)
		if err != nil {
			return nil, err
		}

		if toDelete == nil {
			// geomap already deleted
			return nil, nil
		}

		_, err = s.gtm.DeleteGeoMap(context.Background(), toDelete, config.Global.AkamaiConfig.Domain)
		return nil, err
	}

	referenceGeoMap, err := s.constructAkamaiGeoMap(geomap)
	if err != nil {
		return nil, err
	}

	var backendGeoMap *gtm.GeoMap
	backendGeoMap, err = s.gtm.GetGeoMap(context.Background(), geomap.Id, config.Global.AkamaiConfig.Domain)
	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{
		"Name",
		"DefaultDatacenter",
		"DefaultDatacenter.Nickname",
		"DefaultDatacenter.DatacenterId",
		"Assignments",
		"Assignments.DatacenterBase",
		"Assignments.Countries",
	}
	if utils.DeepEqualFields(backendGeoMap, referenceGeoMap, fieldsToCompare) {
		// everything's equal, nothing to do
		return backendGeoMap, nil
	}

	logger.Debugf("SyncGeomap(%s) for domain %s, changes identified",
		geomap.Id, config.Global.AkamaiConfig.Domain)

	// create or update
	res, err := s.gtm.CreateGeoMap(context.Background(), referenceGeoMap, config.Global.AkamaiConfig.Domain)
	if err != nil {
		logger.Errorf("CreateGeomap(%s) for domain %s failed", referenceGeoMap.Name,
			config.Global.AkamaiConfig.Domain)
		return nil, err
	}

	logger.Infof("CreateGeomap(%s) for domain %s", referenceGeoMap.Name, config.Global.AkamaiConfig.Domain)
	return res.Resource, nil
}
