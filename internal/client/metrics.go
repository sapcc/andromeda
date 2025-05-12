// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

/*
 *   Copyright 2025 SAP SE
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
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/jedib0t/go-pretty/table"

	"github.com/sapcc/andromeda/client/metrics"
	"github.com/sapcc/andromeda/models"
)

var MetricsOptions struct {
	MetricsAkamaiDNSRequests `command:"akamai-dns-requests" description:"Get Akamai GTM total DNS requests via Andromeda API"`
}

type MetricsAkamaiDNSRequests struct {
	DomainID strfmt.UUID `long:"domain-id" required:"true" description:"Andromeda Domain UUID"`
	Start    string      `short:"s" long:"start" description:"Start date/time in RFC3339 format (e.g., 2025-05-01T00:00:00Z)"`
	End      string      `short:"e" long:"end" description:"End date/time in RFC3339 format (e.g., 2025-05-10T15:00:00Z)"`
	Output   string      `short:"o" long:"output" description:"Output format (json, csv, table)" default:"table" choice:"json" choice:"csv" choice:"table"`
}

func (m *MetricsAkamaiDNSRequests) Execute(_ []string) error {
	params := metrics.NewGetMetricsAkamaiTotalDNSRequestsParams()
	params.SetDomainID(m.DomainID)

	// Parse and set Start time
	if m.Start != "" {
		startTime, err := time.Parse(time.RFC3339, m.Start)
		if err != nil {
			return fmt.Errorf("invalid start time format, use RFC3339 (e.g., 2025-05-01T00:00:00Z): %w", err)
		}
		startDateTime := strfmt.DateTime(startTime)
		params.SetStart(&startDateTime)
	}

	// Parse and set End time
	if m.End != "" {
		endTime, err := time.Parse(time.RFC3339, m.End)
		if err != nil {
			return fmt.Errorf("invalid end time format, use RFC3339 (e.g., 2025-05-10T15:00:00Z): %w", err)
		}
		endDateTime := strfmt.DateTime(endTime)
		params.SetEnd(&endDateTime)
	}

	// Call the Andromeda API
	resp, err := AndromedaClient.Metrics.GetMetricsAkamaiTotalDNSRequests(params)
	if err != nil {
		// Try to parse the error response for more details
		apiError, ok := err.(*runtime.APIError)
		if ok {
			errorPayload, ok := apiError.Response.(runtime.ClientResponse)
			if ok && errorPayload.Code() >= 400 {
				// Attempt to read the error model from the response body
				andromedaError := &models.Error{}
				if readErr := json.NewDecoder(errorPayload.Body()).Decode(andromedaError); readErr == nil && andromedaError.Message != "" {
					return fmt.Errorf("API error: %s (code: %d)", andromedaError.Message, errorPayload.Code())
				}
			}
			// Fallback to default error message
			return fmt.Errorf("failed to get Akamai DNS metrics from API: %w", err)
		}
		// Handle non-API errors
		return fmt.Errorf("failed to get Akamai DNS metrics: %w", err)
	}

	// Output the data in the requested format
	switch m.Output {
	case "json":
		jsonBytes, err := json.MarshalIndent(resp.Payload, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(jsonBytes))

	case "csv":
		// Create CSV writer
		w := csv.NewWriter(os.Stdout)
		defer w.Flush()

		// Write header for property info
		w.Write([]string{"Property", "Start Date", "End Date", "Total Requests"})
		w.Write([]string{
			resp.Payload.Property,
			resp.Payload.StartDate.String(),
			resp.Payload.EndDate.String(),
			fmt.Sprintf("%d", resp.Payload.TotalRequests),
		})
		w.Write([]string{}) // Empty row for separation

		// Write header for datacenter info
		w.Write([]string{"Datacenter", "Datacenter ID", "Traffic Target", "Requests", "Percentage"})

		// Write datacenter data
		for dcName, dcData := range resp.Payload.Datacenters {
			w.Write([]string{
				dcName,
				dcData.DatacenterID,
				dcData.TrafficTarget,
				fmt.Sprintf("%d", dcData.TotalRequests),
				fmt.Sprintf("%.2f%%", float64(dcData.Percentage)*100),
			})
		}

	default: // table format
		// Print summary information
		fmt.Printf("Total DNS Requests for Property: %s\n", resp.Payload.Property)
		fmt.Printf("Time Range: %s to %s\n", resp.Payload.StartDate.String(), resp.Payload.EndDate.String())
		fmt.Printf("Total Requests: %d\n\n", resp.Payload.TotalRequests)

		// Set up table for datacenters
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Datacenter", "Datacenter ID", "Traffic Target", "Requests", "Percentage"})

		// Add data rows
		for dcName, dcData := range resp.Payload.Datacenters {
			t.AppendRow(table.Row{
				dcName,
				dcData.DatacenterID,
				dcData.TrafficTarget,
				dcData.TotalRequests,
				fmt.Sprintf("%.2f%%", float64(dcData.Percentage)*100),
			})
		}
		t.Render()
	}

	return nil
}

func init() {
	_, err := Parser.AddCommand("metrics", "Metrics", "Metrics Commands.", &MetricsOptions)
	if err != nil {
		panic(err)
	}
}
