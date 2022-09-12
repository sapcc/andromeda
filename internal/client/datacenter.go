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
	"github.com/go-openapi/strfmt"

	"github.com/sapcc/andromeda/client/datacenters"
	"github.com/sapcc/andromeda/models"
)

var DatacenterOptions struct {
	DatacenterList   `command:"list" description:"List Datacenters"`
	DatacenterShow   `command:"show" description:"Show Datacenter"`
	DatacenterCreate `command:"create" description:"Create Datacenter"`
	DatacenterDelete `command:"delete" description:"Delete Datacenter"`
}

type DatacenterList struct {
}

type DatacenterShow struct {
	Positional struct {
		UUID strfmt.UUID `description:"UUID of the datacenter"`
	} `positional-args:"yes" required:"yes"`
}

type DatacenterCreate struct {
	Name            string   `short:"n" long:"name" description:"Name of the Datacenter"`
	Provider        string   `long:"provider" description:"Provider name" required:"true"`
	Continent       *string  `long:"continent" description:"A two-letter code that specifies the continent where the data center maps to"`
	Country         *string  `long:"country" description:"A two-letter ISO 3166 country code that specifies the country where the data center maps to"`
	StateOrProvince *string  `long:"state_or_province" description:"Specifies a two-letter ISO 3166 country code for the state or province where the data center is located"`
	City            *string  `long:"city" description:"The name of the city where the data center is located"`
	Latitude        *float32 `long:"latitude" description:"Specifies the geographical latitude of the data center's position"`
	Longitude       *float32 `long:"longitude" description:"Specifies the geographical longitude of the data center's position"`
	Disable         bool     `short:"d" long:"disable" description:"Disable Datacenter" optional:"true" optional-value:"false"`
}

type DatacenterDelete struct {
	Positional struct {
		UUID strfmt.UUID `description:"UUID of the datacenter"`
	} `positional-args:"yes" required:"yes"`
}

func (*DatacenterList) Execute(_ []string) error {
	resp, err := AndromedaClient.Datacenters.GetDatacenters(nil)
	if err != nil {
		return err
	}
	return WriteTable(resp.GetPayload().Datacenters)
}

func (*DatacenterCreate) Execute(_ []string) error {
	adminStateUp := !DatacenterOptions.Disable
	datacenter := datacenters.PostDatacentersBody{Datacenter: &models.Datacenter{
		AdminStateUp:    &adminStateUp,
		City:            DatacenterOptions.City,
		Continent:       DatacenterOptions.Continent,
		Country:         DatacenterOptions.Country,
		Latitude:        DatacenterOptions.Latitude,
		Longitude:       DatacenterOptions.Longitude,
		Name:            &DatacenterOptions.Name,
		Provider:        DatacenterOptions.Provider,
		StateOrProvince: DatacenterOptions.StateOrProvince,
	}}
	resp, err := AndromedaClient.Datacenters.PostDatacenters(datacenters.NewPostDatacentersParams().WithDatacenter(datacenter))
	if err != nil {
		return err
	}
	return WriteTable(resp.GetPayload().Datacenter)
}

func (*DatacenterShow) Execute(_ []string) error {
	params := datacenters.
		NewGetDatacentersDatacenterIDParams().
		WithDatacenterID(DatacenterOptions.DatacenterShow.Positional.UUID)
	resp, err := AndromedaClient.Datacenters.GetDatacentersDatacenterID(params)
	if err != nil {
		return err
	}
	return WriteTable(resp.GetPayload().Datacenter)
}

func (*DatacenterDelete) Execute(_ []string) error {
	params := datacenters.
		NewDeleteDatacentersDatacenterIDParams().
		WithDatacenterID(DatacenterOptions.DatacenterDelete.Positional.UUID)

	if _, err := AndromedaClient.Datacenters.DeleteDatacentersDatacenterID(params); err != nil {
		return err
	}
	return nil
}

func init() {
	/*_, err := Parser.AddCommand("datacenter", "Datacenters", "Datacenter Commands.", &DatacenterOptions)
	if err != nil {
		panic(err)
	}*/
}
