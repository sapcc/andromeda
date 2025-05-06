// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"github.com/go-openapi/strfmt"

	"github.com/sapcc/andromeda/client/administrative"
)

type SyncDomain struct {
	Domains []strfmt.UUID `short:"d" long:"domain" required:"true" description:"Domain IDs to sync, can be specified multiple times"`
}

func (sd SyncDomain) Execute(_ []string) error {
	body := administrative.PostSyncBody{Domains: sd.Domains}
	params := administrative.NewPostSyncParams()
	params.Domains = body
	_, err := AndromedaClient.Administrative.PostSync(params)
	return err
}

func init() {
	_, err := Parser.AddCommand("sync", "Sync", "Sync Commands.", &SyncDomain{})
	if err != nil {
		panic(err)
	}
}
