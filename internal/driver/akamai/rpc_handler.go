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

package akamai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/apex/log"

	"github.com/sapcc/andromeda/internal/rpcmodels"
)

type RPCHandlerAkamai struct {
	agent *AkamaiAgent
}

// Sync is a method that handles synchronization requests for Akamai domains.
func (h *RPCHandlerAkamai) Sync(ctx context.Context, request *rpcmodels.SyncRequest) (*rpcmodels.SyncResponse, error) {
	domainIDs := request.GetDomainId()
	log.WithField("domainIDs", domainIDs).Info("[Sync] Syncing domains")

	h.agent.forceSync <- domainIDs
	return &rpcmodels.SyncResponse{}, nil
}

// GetCidrs retrieves CIDR blocks from Akamai's Firewall Rules Manager API.
func (h *RPCHandlerAkamai) GetCidrs(ctx context.Context, request *rpcmodels.GetCidrsRequest) (*rpcmodels.GetCidrsResponse, error) {
	var cidrBlocks *rpcmodels.GetCidrsResponse

	cidrBlocksReq, _ := http.NewRequest(http.MethodGet, "/firewall-rules-manager/v1/cidr-blocks", nil)
	if ret, err := (*h.agent.session).Exec(cidrBlocksReq, cidrBlocks); err != nil {
		log.
			WithError(err).
			WithField("body", ret.Body).
			Error("Error fetching CIDR blocks")
		return nil, err
	}
	return cidrBlocks, nil
}

// DomainInfo represents domain information from the database
type DomainInfo struct {
	ID                 string `db:"id"`
	FQDN               string `db:"fqdn"`
	Provider           string `db:"provider"`
	ProjectID          string `db:"project_id"`
	ProvisioningStatus string `db:"provisioning_status"`
}

// GetDNSMetricsAkamai retrieves DNS metrics
func (h *RPCHandlerAkamai) GetDNSMetricsAkamai(ctx context.Context, req *rpcmodels.GetDNSMetricsRequest) (*rpcmodels.GetDNSMetricsResponse, error) {
	// Now get the traffic report for the property
	// First, get the properties window to determine the time range
	windowURI := "/gtm-api/v1/reports/traffic/properties-window"
	windowReq, err := http.NewRequest(http.MethodGet, windowURI, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating properties window request: %w", err)
	}

	// Define the properties window struct
	var propertiesWindow struct {
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
	}

	// Execute the properties window request
	_, err = (*h.agent.session).Exec(windowReq, &propertiesWindow)
	if err != nil {
		return nil, fmt.Errorf("error fetching properties window: %w", err)
	}

	// Calculate start and end times based on the time range
	end := propertiesWindow.End
	start := propertiesWindow.Start
	switch req.GetTimeRange() {
	case rpcmodels.GetDNSMetricsRequest_LAST_HOUR:
		start = end.Add(-time.Hour)
	case rpcmodels.GetDNSMetricsRequest_LAST_DAY:
		start = end.Add(-24 * time.Hour)
	case rpcmodels.GetDNSMetricsRequest_LAST_WEEK:
		start = end.Add(-7 * 24 * time.Hour)
	}

	// Create request for the traffic report
	trafficURI := fmt.Sprintf("/gtm-api/v1/reports/traffic/domains/%s/properties/%s",
		url.QueryEscape(req.GetDomain()), url.QueryEscape(req.GetProperty()))
	params := fmt.Sprintf("?start=%s&end=%s", start.UTC().Format(time.RFC3339), end.UTC().Format(time.RFC3339))
	trafficURI += params

	trafficReq, err := http.NewRequest(http.MethodGet, trafficURI, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating traffic report request: %w", err)
	}

	// Define traffic report struct
	var trafficReport struct {
		DataRows []struct {
			Timestamp   time.Time `json:"timestamp"`
			Datacenters []struct {
				Nickname          string `json:"nickname"`
				TrafficTargetName string `json:"trafficTargetName"`
				Requests          int    `json:"requests"`
				Status            string `json:"status"`
			} `json:"datacenters"`
		} `json:"dataRows"`
	}

	// Execute the traffic report request
	_, err = (*h.agent.session).Exec(trafficReq, &trafficReport)
	if err != nil {
		return nil, fmt.Errorf("error fetching traffic report: %w", err)
	}

	// Build response
	totalRequests := int64(0)
	var datacenters []*rpcmodels.GetDNSMetricsResponse_Datacenter

	// Keep track of datacenters we've already processed
	processedDatacenters := make(map[string]int)

	// Process traffic data
	for _, dataRow := range trafficReport.DataRows {
		for _, dc := range dataRow.Datacenters {
			totalRequests += int64(dc.Requests)

			// Check if we've already processed this datacenter
			if idx, exists := processedDatacenters[dc.Nickname]; exists {
				// Update existing datacenter
				datacenters[idx].Requests += int64(dc.Requests)
				continue
			}

			// Add datacenter to response
			datacenter := &rpcmodels.GetDNSMetricsResponse_Datacenter{
				DatacenterId: dc.Nickname,
				Requests:     int64(dc.Requests),
				Status:       dc.Status,
			}

			// Extract target IP from traffic target name
			if strings.Contains(dc.TrafficTargetName, " - ") {
				parts := strings.Split(dc.TrafficTargetName, " - ")
				if len(parts) > 1 {
					datacenter.TargetIp = parts[1]
				}
			}

			// Store the index of this datacenter in our map
			processedDatacenters[dc.Nickname] = len(datacenters)
			datacenters = append(datacenters, datacenter)
		}
	}

	// Calculate percentages
	if totalRequests > 0 {
		for _, dc := range datacenters {
			// Explicitly convert to float32 to match the model type
			dc.Percentage = float32(float64(dc.Requests) / float64(totalRequests) * 100)
		}
	}

	// Log the metrics data for debugging
	debugLog := log.WithFields(log.Fields{
		"property":       req.GetProperty(),
		"domain_id":      req.GetDomain(),
		"total_requests": totalRequests,
	})
	if len(datacenters) > 0 {
		jsonData, _ := json.MarshalIndent(datacenters, "", "  ")
		debugLog = debugLog.WithField("datacenters", string(jsonData))
	}
	debugLog.Debug("DNS metrics data")

	return &rpcmodels.GetDNSMetricsResponse{
		Datacenters:   datacenters,
		TotalRequests: totalRequests,
	}, nil
}
