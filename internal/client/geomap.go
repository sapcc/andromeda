/*
 *   Copyright 2021 SAP SE
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

package client

import (
	"errors"
	"fmt"

	"github.com/apex/log"
	"github.com/go-openapi/strfmt"

	"github.com/sapcc/andromeda/client/geographic_maps"
	"github.com/sapcc/andromeda/models"
)

var GeomapOptions struct {
	GeomapList   `command:"list" description:"List Geomaps"`
	GeomapShow   `command:"show" description:"Show Geomap"`
	GeomapCreate `command:"create" description:"Create Geomap"`
	GeomapDelete `command:"delete" description:"Delete Geomap"`
}

type GeomapList struct {
	Datacenter        *strfmt.UUID `short:"d" long:"datacenter" description:"Filter by Datacenter ID"`
	DefaultDatacenter *strfmt.UUID `short:"e" long:"default_datacenter" description:"Filter by default Datacenter ID"`
}

type GeomapShow struct {
	Positional struct {
		UUID strfmt.UUID `description:"UUID of the geographic map"`
	} `positional-args:"yes" required:"yes"`
}

type GeomapCreate struct {
	// Array of datacenter assignments
	Assignment map[string]strfmt.UUID `short:"a" long:"assignment" description:"Datacenter assignment of 2-letter country code and datacenter id, e.g. --assignment=DE:UUID --assignment=US:UUID"`
	Default    *strfmt.UUID           `short:"d" long:"default-datacenter" description:"Default datacenter" required:"yes"`
	Name       *string                `short:"n" long:"name" description:"Name of the geographic map"`
	Provider   string                 `short:"p" long:"provider" description:"Provider name"`
	Scope      *string                `short:"s" long:"scope" description:"Scope of the geographic map" default:"private" choice:"private" choice:"shared"`
}

type GeomapDelete struct {
	Positional struct {
		UUID strfmt.UUID `description:"UUID of the geographic map"`
	} `positional-args:"yes" required:"yes"`
}

func (*GeomapList) Execute(_ []string) error {
	resp, err := AndromedaClient.GeographicMaps.GetGeomaps(
		geographic_maps.NewGetGeomapsParams().
			WithDatacenterID(GeomapOptions.GeomapList.Datacenter).
			WithDefaultDatacenterID(GeomapOptions.GeomapList.DefaultDatacenter),
	)
	if err != nil {
		return err
	}
	return WriteTable(resp.GetPayload().Geomaps)
}

func (*GeomapCreate) Execute(_ []string) error {
	var assignments []*models.GeomapAssignmentsItems0
	for country, datacenter := range GeomapOptions.Assignment {
		assignments = append(assignments, &models.GeomapAssignmentsItems0{
			Country:    country,
			Datacenter: datacenter,
		})
	}
	geomap := geographic_maps.PostGeomapsBody{Geomap: &models.Geomap{
		Assignments:       assignments,
		DefaultDatacenter: GeomapOptions.Default,
		Name:              GeomapOptions.Name,
		Provider:          GeomapOptions.Provider,
		Scope:             GeomapOptions.Scope,
	}}
	resp, err := AndromedaClient.GeographicMaps.PostGeomaps(geographic_maps.NewPostGeomapsParams().WithGeomap(geomap))
	if err != nil {
		return err
	}
	if err = waitForActiveGeomap(resp.Payload.Geomap.ID, false); err != nil {
		log.WithError(err).Error("Failed to wait for geomap to be active")
	}
	return WriteTable(resp.GetPayload().Geomap)
}

func (*GeomapShow) Execute(_ []string) error {
	params := geographic_maps.
		NewGetGeomapsGeomapIDParams().
		WithGeomapID(GeomapOptions.GeomapShow.Positional.UUID)
	resp, err := AndromedaClient.GeographicMaps.GetGeomapsGeomapID(params)
	if err != nil {
		return err
	}
	return WriteTable(resp.GetPayload().Geomap)
}

func (*GeomapDelete) Execute(_ []string) error {
	params := geographic_maps.
		NewDeleteGeomapsGeomapIDParams().
		WithGeomapID(GeomapOptions.GeomapDelete.Positional.UUID)

	if _, err := AndromedaClient.GeographicMaps.DeleteGeomapsGeomapID(params); err != nil {
		return err
	}
	if err := waitForActiveGeomap(GeomapOptions.GeomapDelete.Positional.UUID, true); err != nil {
		log.WithError(err).Error("Failed to wait for geomap to be deleted")
	}
	return nil
}

// waitForActiveGeomap waits for the geomap to be active, or optionally be deleted
func waitForActiveGeomap(id strfmt.UUID, deleted bool) error {
	// if not waiting, return immediately
	if !opts.Wait {
		return nil
	}

	return RetryWithBackoffMax(func() error {
		params := geographic_maps.NewGetGeomapsGeomapIDParams().WithGeomapID(id)
		r, err := AndromedaClient.GeographicMaps.GetGeomapsGeomapID(params)
		if err != nil {
			var getIDNotFound *geographic_maps.GetGeomapsGeomapIDNotFound
			if errors.As(err, &getIDNotFound) && deleted {
				return nil
			}
			return err
		}

		res := r.GetPayload()
		if deleted || res.Geomap.ProvisioningStatus != models.GeomapProvisioningStatusACTIVE {
			return fmt.Errorf("geomap %s is not active yet", id)
		}
		return nil
	})
}

func init() {
	_, err := Parser.AddCommand("geomap", "Geographical Maps", "Geomap Commands.", &GeomapOptions)
	if err != nil {
		panic(err)
	}
}
