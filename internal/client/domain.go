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

	"github.com/sapcc/andromeda/client/domains"
	"github.com/sapcc/andromeda/models"
)

var DomainOptions struct {
	DomainList   `command:"list" description:"List Domains"`
	DomainShow   `command:"show" description:"Show Domain"`
	DomainCreate `command:"create" description:"Create Domain"`
	DomainDelete `command:"delete" description:"Delete Domain"`
}

type DomainList struct {
	Long bool `long:"long" description:"List additional fields in output"`
}

type DomainShow struct {
	Positional struct {
		UUID strfmt.UUID `description:"UUID of the domain"`
	} `positional-args:"yes" required:"yes"`
}

type DomainCreate struct {
	Name       string        `short:"n" long:"name" description:"Name of the Domain"`
	Provider   string        `short:"v" long:"provider" description:"Provider name" required:"true"`
	FQDN       string        `short:"q" long:"fqdn" description:"Fully qualified domain name" required:"true"`
	RecordType string        `short:"r" long:"recordtype" description:"Record type" default:"A"`
	Pools      []strfmt.UUID `short:"p" long:"pool" description:"Pool ID to associate, can be specified multiple times"`
	Disable    bool          `short:"d" long:"disable" description:"Disable Domain" optional:"true" optional-value:"false"`
}

type DomainDelete struct {
	Positional struct {
		UUID strfmt.UUID `description:"UUID of the domain"`
	} `positional-args:"yes" required:"yes"`
}

func (*DomainList) Execute(_ []string) error {
	resp, err := AndromedaClient.Domains.GetDomains(nil)

	if err != nil {
		return err
	}
	if !DomainOptions.Long && len(opts.Formatters.Columns) == 0 {
		opts.Formatters.Columns = []string{"id", "name", "admin_state_up", "fqdn", "project_id", "provisioning_status"}
	}
	return WriteTable(resp.GetPayload().Domains)
}

func (*DomainCreate) Execute(_ []string) error {
	fqdn := strfmt.Hostname(DomainOptions.FQDN)

	domain := domains.PostDomainsBody{&models.Domain{
		Name:       &DomainOptions.Name,
		Fqdn:       &fqdn,
		Provider:   &DomainOptions.Provider,
		RecordType: &DomainOptions.RecordType,
		Pools:      DomainOptions.Pools,
	}}
	resp, err := AndromedaClient.Domains.PostDomains(domains.NewPostDomainsParams().WithDomain(domain))
	if err != nil {
		return err
	}
	return WriteTable(resp.GetPayload().Domain)
}

func (*DomainShow) Execute(_ []string) error {
	params := domains.
		NewGetDomainsDomainIDParams().
		WithDomainID(DomainOptions.DomainShow.Positional.UUID)
	resp, err := AndromedaClient.Domains.GetDomainsDomainID(params)
	if err != nil {
		return err
	}
	return WriteTable(resp.GetPayload().Domain)
}

func (*DomainDelete) Execute(_ []string) error {
	params := domains.
		NewDeleteDomainsDomainIDParams().
		WithDomainID(DomainOptions.DomainDelete.Positional.UUID)

	if _, err := AndromedaClient.Domains.DeleteDomainsDomainID(params); err != nil {
		return err
	}
	return nil
}

func init() {
	_, err := Parser.AddCommand("domain", "Domains", "Domain Commands.", &DomainOptions)
	if err != nil {
		panic(err)
	}
}
