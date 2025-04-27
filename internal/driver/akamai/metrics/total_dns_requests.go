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

package metrics

import (
	"fmt"
	"net/url"
	"time"

	"github.com/apex/log"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	fQNameTotalDNSRequests = prometheus.BuildFQName("andromeda", "akamai", "total_dns_requests")
	totalDNSRequestsLabels = []string{"domain", "property", "datacenter_id", "datacenter_name", "traffic_target"}
	descTotalDNSRequests   = prometheus.NewDesc(fQNameTotalDNSRequests, "Total DNS requests per property and datacenter", totalDNSRequestsLabels, nil)
)

// TotalDNSRequestsCollector collects total DNS request metrics from Akamai GTM
type TotalDNSRequestsCollector struct {
	session          *CachedAkamaiSession
	managementDomain string
	startTime        time.Time
	endTime          time.Time
}

// NewTotalDNSRequestsCollector creates a new collector for total DNS requests
func NewTotalDNSRequestsCollector(session *CachedAkamaiSession, managementDomain string) *TotalDNSRequestsCollector {
	// Default to last 2 days
	endTime := time.Now().Add(-15 * time.Minute) // Offset by 15 mins to avoid data delays
	startTime := endTime.Add(-48 * time.Hour)

	return &TotalDNSRequestsCollector{
		session:          session,
		managementDomain: managementDomain,
		startTime:        startTime,
		endTime:          endTime,
	}
}

// Describe implements prometheus.Collector
func (c *TotalDNSRequestsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- descTotalDNSRequests
}

// Collect implements prometheus.Collector
func (c *TotalDNSRequestsCollector) Collect(ch chan<- prometheus.Metric) {
	properties, err := c.session.getProperties()
	if err != nil {
		log.Errorf("Failed to get properties: %v", err)
		return
	}

	for _, property := range properties {
		c.collectPropertyMetrics(ch, property)
	}
}

// collectPropertyMetrics collects metrics for a specific property
func (c *TotalDNSRequestsCollector) collectPropertyMetrics(ch chan<- prometheus.Metric, property string) {
	datarows, err := c.session.getTrafficReport(property)
	if err != nil {
		log.Errorf("Failed to get traffic report for %s: %v", property, err)
		return
	}

	// Use maps to accumulate totals by datacenter
	datacenterTotals := make(map[string]struct {
		ID            string
		Name          string
		TrafficTarget string
		Requests      int
	})

	// Process all data rows to accumulate totals by datacenter
	for _, dataRow := range datarows {
		for _, datacenter := range dataRow.Datacenters {
			key := fmt.Sprintf("%s-%s", property, datacenter.Nickname)
			entry := datacenterTotals[key]

			entry.ID = datacenter.Nickname
			entry.Name = datacenter.Nickname
			entry.TrafficTarget = datacenter.TrafficTargetName
			entry.Requests += datacenter.Requests

			datacenterTotals[key] = entry
		}
	}

	// Expose metrics for each datacenter
	for _, entry := range datacenterTotals {
		ch <- prometheus.MustNewConstMetric(
			descTotalDNSRequests,
			prometheus.GaugeValue,
			float64(entry.Requests),
			c.managementDomain,
			property,
			entry.ID,
			entry.Name,
			entry.TrafficTarget,
		)
	}
}

// GetTotalDNSRequests retrieves DNS request data programmatically (for CLI or API use)
func GetTotalDNSRequests(session *CachedAkamaiSession, domain string, property string, startTime, endTime time.Time) (map[string]interface{}, error) {
	// Build the URL for direct API access to get data for the specific time range
	path := fmt.Sprintf("/gtm-api/v1/reports/traffic/domains/%s/properties/%s", domain, property)
	params := url.Values{}
	params.Add("start", startTime.UTC().Format(time.RFC3339))
	params.Add("end", endTime.UTC().Format(time.RFC3339))
	uri := fmt.Sprintf("%s?%s", path, params.Encode())

	log.Infof("Fetching traffic report with custom time range: %s", uri)

	// Make the direct API request
	var trafficReport TrafficReport
	err := session.get(uri, &trafficReport)
	if err != nil {
		return nil, fmt.Errorf("failed to get traffic report: %v", err)
	}

	totalRequests := 0
	datacenterRequests := make(map[string]map[string]interface{})

	// Process all data rows
	for _, dataRow := range trafficReport.DataRows {
		for _, datacenter := range dataRow.Datacenters {
			dcName := datacenter.Nickname
			requests := datacenter.Requests

			// Update total requests count
			totalRequests += requests

			// Update datacenter-specific counts
			if _, exists := datacenterRequests[dcName]; !exists {
				datacenterRequests[dcName] = map[string]interface{}{
					"datacenter_id":  dcName,
					"total_requests": 0,
					"traffic_target": datacenter.TrafficTargetName,
				}
			}
			dcData := datacenterRequests[dcName]
			dcData["total_requests"] = dcData["total_requests"].(int) + requests
			datacenterRequests[dcName] = dcData
		}
	}

	// Calculate percentages for each datacenter
	for dcName, dcData := range datacenterRequests {
		if totalRequests > 0 {
			percentage := float64(dcData["total_requests"].(int)) / float64(totalRequests)
			dcData["percentage"] = percentage
		} else {
			dcData["percentage"] = 0.0
		}
		datacenterRequests[dcName] = dcData
	}

	// Build result
	result := map[string]interface{}{
		"property":       property,
		"start_date":     startTime.Format(time.RFC3339),
		"end_date":       endTime.Format(time.RFC3339),
		"total_requests": totalRequests,
		"datacenters":    datacenterRequests,
	}

	return result, nil
}
