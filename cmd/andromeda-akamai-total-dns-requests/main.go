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

package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/urfave/cli/v2"

	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/driver/akamai"
	"github.com/sapcc/andromeda/internal/driver/akamai/metrics"
)

func main() {
	cliFlags := []cli.Flag{
		&cli.StringFlag{
			Name:  "domain",
			Value: "andromeda.akadns.net",
			Usage: "GTM domain name",
		},
		&cli.StringFlag{
			Name:     "property",
			Required: true,
			Usage:    "GTM property ID",
		},
		&cli.StringFlag{
			Name:  "start",
			Usage: "Start date in ISO format (default: 2 days before end date)",
		},
		&cli.StringFlag{
			Name:  "end",
			Usage: "End date in ISO format (default: 15 minutes ago)",
		},
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Value:   "text",
			Usage:   "Output format (json, csv, text)",
		},
	}

	config.ParseArgsAndRun("andromeda-akamai-total-dns-requests", "Get total DNS requests for an Akamai GTM property",
		func(c *cli.Context) error {
			// Set up default dates if not provided
			var endTime time.Time
			var startTime time.Time
			var err error

			if c.String("end") != "" {
				endTime, err = time.Parse(time.RFC3339, c.String("end"))
				if err != nil {
					return fmt.Errorf("invalid end time format, use RFC3339: %w", err)
				}
			} else {
				endTime = time.Now().Add(-15 * time.Minute) // Default end time: 15 minutes ago
			}

			if c.String("start") != "" {
				startTime, err = time.Parse(time.RFC3339, c.String("start"))
				if err != nil {
					return fmt.Errorf("invalid start time format, use RFC3339: %w", err)
				}
			} else {
				startTime = endTime.Add(-48 * time.Hour) // Default start time: 2 days before end time
			}

			// Make sure contract_id is set to bypass the check
			// This is needed because the credentials don't have access to the contract information - Can be removed
			if config.Global.AkamaiConfig.ContractId == "" {
				config.Global.AkamaiConfig.ContractId = "BYPASS_CHECK_VALUE"
				log.Info("Setting ContractId to BYPASS_CHECK_VALUE to skip contract check")
			}

			// Initialize Akamai session using patched version
			session, domainType := akamai.NewAkamaiSessionPatched(&config.Global.AkamaiConfig)
			log.Infof("Connected to Akamai API with domain type: %s", domainType)

			// Create cached session
			cachedSession := metrics.NewCachedAkamaiSession(*session, c.String("domain"))

			// Get the data
			log.Infof("Fetching total DNS requests for %s from %s to %s",
				c.String("property"), startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))

			result, err := metrics.GetTotalDNSRequests(cachedSession, c.String("domain"), c.String("property"), startTime, endTime)
			if err != nil {
				return fmt.Errorf("failed to get total DNS requests: %w", err)
			}

			// Output the data in the requested format
			switch c.String("output") {
			case "json":
				jsonBytes, err := json.MarshalIndent(result, "", "  ")
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
					result["property"].(string),
					result["start_date"].(string),
					result["end_date"].(string),
					fmt.Sprintf("%d", result["total_requests"].(int)),
				})
				w.Write([]string{}) // Empty row for separation

				// Write header for datacenter info
				w.Write([]string{"Datacenter", "Datacenter ID", "Traffic Target", "Requests", "Percentage"})

				// Write datacenter data
				datacenters := result["datacenters"].(map[string]map[string]interface{})
				for dcName, dcData := range datacenters {
					w.Write([]string{
						dcName,
						dcData["datacenter_id"].(string),
						dcData["traffic_target"].(string),
						fmt.Sprintf("%d", dcData["total_requests"].(int)),
						fmt.Sprintf("%.2f%%", dcData["percentage"].(float64)*100),
					})
				}

			default: // text format
				fmt.Printf("Total DNS Requests for Property: %s\n", result["property"])
				fmt.Printf("Time Range: %s to %s\n", result["start_date"], result["end_date"])
				fmt.Printf("Total Requests: %d\n", result["total_requests"])

				fmt.Println(strings.Repeat("-", 100))
				fmt.Println("Datacenter                     | ID              | Traffic Target                | Requests        | Percentage")
				fmt.Println(strings.Repeat("-", 100))

				datacenters := result["datacenters"].(map[string]map[string]interface{})
				for dcName, dcData := range datacenters {
					// Format each field to fit column width
					dcNameStr := dcName
					if len(dcNameStr) > 30 {
						dcNameStr = dcNameStr[:27] + "..."
					} else {
						dcNameStr = fmt.Sprintf("%-30s", dcNameStr)
					}

					dcIDStr := dcData["datacenter_id"].(string)
					if len(dcIDStr) > 15 {
						dcIDStr = dcIDStr[:12] + "..."
					} else {
						dcIDStr = fmt.Sprintf("%-15s", dcIDStr)
					}

					trafficTargetStr := dcData["traffic_target"].(string)
					if len(trafficTargetStr) > 30 {
						trafficTargetStr = trafficTargetStr[:27] + "..."
					} else {
						trafficTargetStr = fmt.Sprintf("%-30s", trafficTargetStr)
					}

					requestsStr := fmt.Sprintf("%-15d", dcData["total_requests"].(int))
					percentageStr := fmt.Sprintf("%.2f%%", dcData["percentage"].(float64)*100)

					fmt.Printf("%s | %s | %s | %s | %s\n",
						dcNameStr, dcIDStr, trafficTargetStr, requestsStr, percentageStr)
				}
			}

			return nil
		}, cliFlags...)
}
