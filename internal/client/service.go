// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"github.com/jedib0t/go-pretty/table"
)

var ServiceOptions struct {
	ServiceList `command:"list" description:"List Services"`
}

type ServiceList struct{}

func (*ServiceList) Execute(_ []string) error {
	resp, err := AndromedaClient.Administrative.GetServices(nil)
	if err != nil {
		return err
	}

	Table.AppendHeader(table.Row{"ID", "Type", "Provider", "Host"})
	for _, service := range resp.Payload.Services {
		Table.AppendRow(table.Row{service.ID, service.Type, service.Provider, service.Host})
	}
	Table.Render()
	return nil
}

func init() {
	_, _ = Parser.AddCommand("service", "Services", "Service Commands.", &ServiceOptions)
}
