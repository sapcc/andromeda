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

	"github.com/sapcc/andromeda/client/pools"
	"github.com/sapcc/andromeda/models"
)

var PoolOptions struct {
	PoolList   `command:"list" description:"List Pools"`
	PoolShow   `command:"show" description:"Show Pool"`
	PoolCreate `command:"create" description:"Create Pool"`
	PoolDelete `command:"delete" description:"Delete Pool"`
}

type PoolList struct {
	Domain *strfmt.UUID `short:"a" long:"domain" description:"Filter by Domain ID"`
	Pool   *strfmt.UUID `short:"p" long:"pool" description:"Filter by Pool ID"`
}

type PoolShow struct {
	Positional struct {
		UUID strfmt.UUID `description:"UUID of the pool"`
	} `positional-args:"yes" required:"yes"`
}

type PoolCreate struct {
	Name    string   `short:"n" long:"name" description:"Name of the Pool"`
	Domain  []string `short:"a" long:"domain" description:"ID(s) of the associated Domain (multiple domains allowed)"`
	Disable bool     `short:"d" long:"disable" description:"Disable Pool" optional:"true" optional-value:"false"`
}

type PoolDelete struct {
	Positional struct {
		UUID strfmt.UUID `description:"UUID of the pool"`
	} `positional-args:"yes" required:"yes"`
}

func (*PoolList) Execute(_ []string) error {
	resp, err := AndromedaClient.Pools.GetPools(pools.NewGetPoolsParams().
		WithDomainID(PoolOptions.PoolList.Domain).
		WithPoolID(PoolOptions.PoolList.Pool))
	if err != nil {
		return err
	}
	return WriteTable(resp.GetPayload().Pools)
}

func (*PoolCreate) Execute(_ []string) error {
	adminStateUp := !PoolOptions.Disable
	domains := []strfmt.UUID{}
	for _, d := range PoolOptions.PoolCreate.Domain {
		domains = append(domains, strfmt.UUID(d))
	}
	pool := pools.PostPoolsBody{Pool: &models.Pool{
		AdminStateUp: &adminStateUp,
		Name:         &PoolOptions.Name,
		Domains:      domains,
	}}
	resp, err := AndromedaClient.Pools.PostPools(pools.NewPostPoolsParams().WithPool(pool))
	if err != nil {
		return err
	}
	return WriteTable(resp.GetPayload().Pool)
}

func (*PoolShow) Execute(_ []string) error {
	params := pools.
		NewGetPoolsPoolIDParams().
		WithPoolID(PoolOptions.PoolShow.Positional.UUID)
	resp, err := AndromedaClient.Pools.GetPoolsPoolID(params)
	if err != nil {
		return err
	}
	return WriteTable(resp.GetPayload().Pool)
}

func (*PoolDelete) Execute(_ []string) error {
	params := pools.
		NewDeletePoolsPoolIDParams().
		WithPoolID(PoolOptions.PoolDelete.Positional.UUID)

	if _, err := AndromedaClient.Pools.DeletePoolsPoolID(params); err != nil {
		return err
	}
	return nil
}

func init() {
	_, err := Parser.AddCommand("pool", "Pools", "Pool Commands.", &PoolOptions)
	if err != nil {
		panic(err)
	}
}
