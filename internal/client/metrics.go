// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"fmt"

	"github.com/sapcc/andromeda/client/metrics"
)

// CLI command options
var MetricsOptions struct {
	MetricsAkamaiDNS `command:"akamai-dns-requests" description:"Get Akamai GTM DNS request metrics"`
}

type MetricsAkamaiDNS struct {
	PropertyName string `short:"p" long:"property-name" description:"Filter by Akamai GTM property name"`
	ProjectID    string `short:"i" long:"project-id" description:"Filter by project ID"`
	TimeRange    string `short:"t" long:"time-range" description:"Time range for metrics data" default:"last_hour" choice:"last_hour" choice:"last_day" choice:"last_week"`
}

// Execute the metrics command
func (*MetricsAkamaiDNS) Execute(_ []string) error {
	// Create params for the API call
	params := metrics.NewGetMetricsAkamaiTotalDNSRequestsParams()

	// Set request parameters if provided
	if MetricsOptions.MetricsAkamaiDNS.PropertyName != "" {
		params.PropertyName = &MetricsOptions.MetricsAkamaiDNS.PropertyName
	}
	if MetricsOptions.MetricsAkamaiDNS.ProjectID != "" {
		params.ProjectID = &MetricsOptions.MetricsAkamaiDNS.ProjectID
	}
	if MetricsOptions.MetricsAkamaiDNS.TimeRange != "" {
		params.TimeRange = &MetricsOptions.MetricsAkamaiDNS.TimeRange
	}

	// Call the Andromeda API endpoint
	resp, err := AndromedaClient.Metrics.GetMetricsAkamaiTotalDNSRequests(params)
	if err != nil {
		return fmt.Errorf("failed to get DNS metrics: %w", err)
	}

	// Write result to table
	return WriteTable(resp.Payload.TotalDNSRequests)
}

// Initialize the command
func init() {
	_, err := Parser.AddCommand("metrics", "Metrics", "Metrics Commands.", &MetricsOptions)
	if err != nil {
		panic(err)
	}
}
