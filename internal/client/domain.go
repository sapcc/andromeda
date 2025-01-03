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

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/sapcc/andromeda/client/domains"
	"github.com/sapcc/andromeda/models"
)

var DomainOptions struct {
	DomainList   `command:"list" description:"List Domains"`
	DomainShow   `command:"show" description:"Show Domain"`
	DomainCreate `command:"create" description:"Create Domain"`
	DomainDelete `command:"delete" description:"Delete Domain"`
	DomainSet    `command:"set" description:"Update Domain"`
}

type DomainList struct {
	Long   bool         `long:"long" description:"List additional fields in output"`
	Domain *strfmt.UUID `short:"a" long:"domain" description:"Filter by Domain ID"`
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
	Mode       string        `short:"m" long:"mode" description:"Load balancing method to use for the references pools." default:"ROUND_ROBIN" choice:"ROUND_ROBIN" choice:"WEIGHTED" choice:"GEOGRAPHIC" choice:"AVAILABILITY"`
	RecordType string        `short:"r" long:"recordtype" description:"Record type" default:"A"`
	Pools      []strfmt.UUID `short:"p" long:"pool" description:"Pool ID to associate, can be specified multiple times"`
	Disable    bool          `short:"d" long:"disable" description:"Disable Domain" optional:"true" optional-value:"false"`
}

type DomainDelete struct {
	Positional struct {
		UUID strfmt.UUID `description:"UUID of the domain"`
	} `positional-args:"yes" required:"yes"`
}

type DomainSet struct {
	Positional struct {
		UUID strfmt.UUID `description:"UUID of the domain"`
	} `positional-args:"yes" required:"yes"`
	Name       string        `short:"n" long:"name" description:"Name of the Domain"`
	FQDN       string        `short:"q" long:"fqdn" description:"Fully qualified domain name"`
	Mode       string        `short:"m" long:"mode" description:"Load balancing method to use for the references pools." optional:"true" choice:"ROUND_ROBIN" choice:"WEIGHTED" choice:"GEOGRAPHIC" choice:"AVAILABILITY"`
	RecordType string        `short:"r" long:"recordtype" description:"Record type"`
	Pools      []strfmt.UUID `short:"p" long:"pool" description:"Pool ID to associate, can be specified multiple times"`
	NoPools    bool          `long:"no-pools" description:"Remove all pools from domain" optional:"true" optional-value:"true"`
	Disable    bool          `short:"d" long:"disable" description:"Enable Domain" optional:"true" optional-value:"true"`
	Enable     bool          `short:"e" long:"enable" description:"Enable Domain" optional:"true" optional-value:"true"`
}

func (*DomainList) Execute(_ []string) error {
	resp, err := AndromedaClient.Domains.GetDomains(
		domains.NewGetDomainsParams().WithDomainID(DomainOptions.DomainList.Domain))

	if err != nil {
		return err
	}
	if !DomainOptions.Long && len(opts.Formatters.Columns) == 0 {
		opts.Formatters.Columns = []string{"id", "name", "admin_state_up", "fqdn", "project_id", "provisioning_status"}
	}
	return WriteTable(resp.GetPayload().Domains)
}

func (*DomainCreate) Execute(_ []string) error {
	fqdn := strfmt.Hostname(DomainOptions.DomainCreate.FQDN)
	domain := domains.PostDomainsBody{Domain: &models.Domain{
		Name:       &DomainOptions.DomainCreate.Name,
		Fqdn:       &fqdn,
		Mode:       &DomainOptions.DomainCreate.Mode,
		Provider:   &DomainOptions.Provider,
		RecordType: &DomainOptions.DomainCreate.RecordType,
		Pools:      DomainOptions.DomainCreate.Pools,
	}}
	resp, err := AndromedaClient.Domains.PostDomains(domains.NewPostDomainsParams().WithDomain(domain))
	if err != nil {
		return err
	}
	if err = waitForActiveDomain(resp.Payload.Domain.ID, false); err != nil {
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
	if err := waitForActiveDomain(DomainOptions.DomainDelete.Positional.UUID, true); err != nil {
		return fmt.Errorf("failed to wait for domain %s to be deleted", DomainOptions.DomainDelete.Positional.UUID)
	}
	return nil
}

func (*DomainSet) Execute(_ []string) error {
	if DomainOptions.DomainSet.Disable && DomainOptions.DomainSet.Enable {
		return fmt.Errorf("cannot enable and disable domain at the same time")
	}

	if DomainOptions.DomainSet.NoPools && len(DomainOptions.DomainSet.Pools) > 0 {
		return fmt.Errorf("cannot remove all pools and add new pools at the same time")
	}

	domain := domains.PutDomainsDomainIDBody{Domain: &models.Domain{}}
	if DomainOptions.DomainSet.Disable {
		domain.Domain.AdminStateUp = swag.Bool(false)
	} else if DomainOptions.DomainSet.Enable {
		domain.Domain.AdminStateUp = swag.Bool(true)
	}
	if DomainOptions.DomainSet.Name != "" {
		domain.Domain.Name = &DomainOptions.DomainSet.Name
	}
	if DomainOptions.DomainSet.FQDN != "" {
		fqdn := strfmt.Hostname(DomainOptions.DomainSet.FQDN)
		domain.Domain.Fqdn = &fqdn
	}
	if DomainOptions.DomainSet.Mode != "" {
		domain.Domain.Mode = &DomainOptions.DomainSet.Mode
	}
	if len(DomainOptions.DomainSet.Pools) > 0 {
		domain.Domain.Pools = DomainOptions.DomainSet.Pools
	} else if DomainOptions.DomainSet.NoPools {
		domain.Domain.Pools = []strfmt.UUID{}
	}

	params := domains.
		NewPutDomainsDomainIDParams().
		WithDomainID(DomainOptions.DomainSet.Positional.UUID).
		WithDomain(domain)

	resp, err := AndromedaClient.Domains.PutDomainsDomainID(params)
	if err != nil {
		return err
	}
	if err = waitForActiveDomain(DomainOptions.DomainSet.Positional.UUID, false); err != nil {
		return fmt.Errorf("failed to wait for domain %s to be active", DomainOptions.DomainSet.Positional.UUID)
	}
	return WriteTable(resp.GetPayload().Domain)
}

// waitForActiveDomain waits for the domain to be active, or optionally be deleted
func waitForActiveDomain(id strfmt.UUID, deleted bool) error {
	// if not waiting, return immediately
	if !opts.Wait {
		return nil
	}

	return RetryWithBackoffMax(func() error {
		params := domains.NewGetDomainsDomainIDParams().WithDomainID(id)
		r, err := AndromedaClient.Domains.GetDomainsDomainID(params)
		if err != nil {
			var getIDNotFound *domains.GetDomainsDomainIDNotFound
			if errors.As(err, &getIDNotFound) && deleted {
				return nil
			}
			return err
		}

		res := r.GetPayload()
		if deleted || res.Domain.ProvisioningStatus != models.DomainProvisioningStatusACTIVE {
			return fmt.Errorf("domain %s is not active yet", id)
		}
		return nil
	})
}

func init() {
	_, err := Parser.AddCommand("domain", "Domains", "Domain Commands.", &DomainOptions)
	if err != nil {
		panic(err)
	}
}
