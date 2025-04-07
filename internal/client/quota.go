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

package client

import (
	"github.com/go-openapi/strfmt"
	"github.com/jedib0t/go-pretty/table"

	"github.com/sapcc/andromeda/client/administrative"
	"github.com/sapcc/andromeda/models"
)

var QuotaOptions struct {
	QuotaList     `command:"list" description:"List Quotas"`
	QuotaDefaults `command:"defaults" description:"Show default project quota"`
	QuotaShow     `command:"show" description:"Show project quota"`
	QuotaUpdate   `command:"update" description:"Update project quota"`
	QuotaDelete   `command:"reset" description:"Reset project quota to defaults"`
}

type QuotaList struct{}
type QuotaDefaults struct{}
type QuotaShow struct {
	Positional struct {
		ProjectID strfmt.UUID `description:"The ID of the project to query"`
	} `positional-args:"yes" required:"yes"`
}
type QuotaUpdate struct {
	Positional struct {
		ProjectID strfmt.UUID `description:"The ID of the project to update"`
	} `positional-args:"yes" required:"yes"`
	DomainAkamai *int64 `long:"domain_akamai" description:"Domains (provider Akamai) integer value"`
	DomainF5     *int64 `long:"domain_f5" description:"Domains (provider F5) integer value"`
	Datacenter   *int64 `long:"datacenter" description:"Datacenter integer value"`
	Pool         *int64 `long:"pool" description:"Pool integer value"`
	Member       *int64 `long:"member" description:"Member integer value"`
	Monitor      *int64 `long:"monitor" description:"Monitor integer value"`
}
type QuotaDelete struct {
	Positional struct {
		ProjectID strfmt.UUID `description:"The ID of the project to reset"`
	} `positional-args:"yes" required:"yes"`
}

func (*QuotaList) Execute(_ []string) error {
	resp, err := AndromedaClient.Administrative.GetQuotas(nil)
	if err != nil {
		return err
	}

	Table.AppendHeader(table.Row{"Project ID", "DomainsAkamai", "DomainsF5", "Datacenters", "Pools", "Members", "Monitors"})
	for _, quota := range resp.Payload.Quotas {
		Table.AppendRow(table.Row{*quota.ProjectID, *quota.DomainAkamai, *quota.DomainF5, *quota.Datacenter, *quota.Pool,
			*quota.Member, *quota.Monitor})
	}
	Table.Render()
	return nil
}

func (*QuotaDefaults) Execute(_ []string) error {
	resp, err := AndromedaClient.Administrative.GetQuotasDefaults(nil)
	if err != nil {
		return err
	}

	Table.AppendHeader(table.Row{"DomainsAkamai", "DomainsF5", "Datacenters", "Pools", "Members", "Monitors"})
	domains_akamai := int(*resp.Payload.Quota.DomainAkamai)
	domains_f5 := int(*resp.Payload.Quota.DomainF5)
	datacenters := int(*resp.Payload.Quota.Datacenter)
	pools := int(*resp.Payload.Quota.Pool)
	members := int(*resp.Payload.Quota.Member)
	monitors := int(*resp.Payload.Quota.Monitor)
	Table.AppendRow(table.Row{domains_akamai, domains_f5, datacenters, pools, members, monitors})
	Table.Render()
	return nil
}

func (*QuotaShow) Execute(_ []string) error {
	params := administrative.NewGetQuotasProjectIDParams()
	params.ProjectID = QuotaOptions.QuotaShow.Positional.ProjectID.String()

	resp, err := AndromedaClient.Administrative.GetQuotasProjectID(params)
	if err != nil {
		return err
	}

	return WriteTable(resp.GetPayload().Quota)
}

func (*QuotaUpdate) Execute(_ []string) error {
	params := administrative.NewPutQuotasProjectIDParams().
		WithProjectID(QuotaOptions.QuotaUpdate.Positional.ProjectID.String()).
		WithQuota(administrative.PutQuotasProjectIDBody{
			Quota: &models.Quota{
				Datacenter:   QuotaOptions.QuotaUpdate.Datacenter,
				DomainAkamai: QuotaOptions.QuotaUpdate.DomainAkamai,
				DomainF5:     QuotaOptions.QuotaUpdate.DomainF5,
				Member:       QuotaOptions.QuotaUpdate.Member,
				Monitor:      QuotaOptions.QuotaUpdate.Monitor,
				Pool:         QuotaOptions.QuotaUpdate.Pool,
			},
		})

	resp, err := AndromedaClient.Administrative.PutQuotasProjectID(params)
	if err != nil {
		return err
	}

	return WriteTable(resp.GetPayload().Quota)
}

func (*QuotaDelete) Execute(_ []string) error {
	params := administrative.NewDeleteQuotasProjectIDParams().
		WithProjectID(QuotaOptions.QuotaDelete.Positional.ProjectID.String())

	_, err := AndromedaClient.Administrative.DeleteQuotasProjectID(params)
	return err
}

func init() {
	_, _ = Parser.AddCommand("quota", "Quotas", "Quota Commands.", &QuotaOptions)
}
