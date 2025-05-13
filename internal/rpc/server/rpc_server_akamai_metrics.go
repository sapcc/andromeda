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

package server

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	"github.com/apex/log"

	"github.com/sapcc/andromeda/internal/config"
)

// AkamaiTotalDNSRequestsRequest defines the request for Akamai DNS metrics
type AkamaiTotalDNSRequestsRequest struct {
	Domain    string `protobuf:"bytes,1,opt,name=domain,proto3" json:"domain,omitempty"`
	Property  string `protobuf:"bytes,2,opt,name=property,proto3" json:"property,omitempty"`
	StartTime string `protobuf:"bytes,3,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	EndTime   string `protobuf:"bytes,4,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`
}

// AkamaiDatacenterStats defines statistics for a datacenter
type AkamaiDatacenterStats struct {
	DatacenterId  string  `protobuf:"bytes,1,opt,name=datacenter_id,json=datacenterId,proto3" json:"datacenter_id,omitempty"`
	TrafficTarget string  `protobuf:"bytes,2,opt,name=traffic_target,json=trafficTarget,proto3" json:"traffic_target,omitempty"`
	TotalRequests int64   `protobuf:"varint,3,opt,name=total_requests,json=totalRequests,proto3" json:"total_requests,omitempty"`
	Percentage    float32 `protobuf:"fixed32,4,opt,name=percentage,proto3" json:"percentage,omitempty"`
}

// AkamaiTotalDNSRequestsResponse defines the response for Akamai DNS metrics
type AkamaiTotalDNSRequestsResponse struct {
	Property      string                            `protobuf:"bytes,1,opt,name=property,proto3" json:"property,omitempty"`
	StartDate     string                            `protobuf:"bytes,2,opt,name=start_date,json=startDate,proto3" json:"start_date,omitempty"`
	EndDate       string                            `protobuf:"bytes,3,opt,name=end_date,json=endDate,proto3" json:"end_date,omitempty"`
	TotalRequests int64                             `protobuf:"varint,4,opt,name=total_requests,json=totalRequests,proto3" json:"total_requests,omitempty"`
	Datacenters   map[string]*AkamaiDatacenterStats `protobuf:"bytes,5,rep,name=datacenters,proto3" json:"datacenters,omitempty"`
	Error         string                            `protobuf:"bytes,6,opt,name=error,proto3" json:"error,omitempty"`
}

// DataRows represents a data row in the Akamai traffic report
type DataRows struct {
	Timestamp   time.Time `json:"timestamp"`
	Datacenters []struct {
		Nickname          string `json:"nickname"`
		TrafficTargetName string `json:"trafficTargetName"`
		Requests          int    `json:"requests"`
		Status            string `json:"status"`
	} `json:"datacenters"`
}

// TrafficReport represents the Akamai traffic report
type TrafficReport struct {
	Metadata struct {
		Domain   string    `json:"domain"`
		Property string    `json:"property"`
		Start    time.Time `json:"start"`
		End      time.Time `json:"end"`
		Interval string    `json:"interval"`
		Uri      string    `json:"uri"`
	} `json:"metadata"`
	DataRows []DataRows `json:"dataRows"`
	Links    []struct {
		Rel  string `json:"rel"`
		Href string `json:"href"`
	} `json:"links"`
}

// createAkamaiSession creates a new session with the Akamai API
func createAkamaiSession(akamaiConfig *config.AkamaiConfig) (*session.Session, error) {
	option := edgegrid.WithEnv(true)
	if env := os.Getenv("AKAMAI_EDGE_RC"); env != "" {
		option = edgegrid.WithFile(env)
	} else if akamaiConfig.EdgeRC != "" {
		option = edgegrid.WithFile(akamaiConfig.EdgeRC)
	}

	edgerc, err := edgegrid.New(option)
	if err != nil {
		return nil, fmt.Errorf("failed to create EdgeGrid credentials: %w", err)
	}

	s, err := session.New(session.WithSigner(edgerc))
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &s, nil
}

// GetAkamaiTotalDNSRequests implements the RPC method for retrieving Akamai total DNS requests
func (u *RPCHandler) GetAkamaiTotalDNSRequests(ctx context.Context, req *AkamaiTotalDNSRequestsRequest) (*AkamaiTotalDNSRequestsResponse, error) {
	logger := log.WithFields(log.Fields{
		"property": req.Property,
		"domain":   req.Domain,
	})
	logger.Info("RPC: Getting total DNS requests for Akamai GTM property")

	// Set default domain if not provided
	domain := "andromeda.akadns.net"
	if req.Domain != "" {
		domain = req.Domain
	}

	// Parse time values
	var startTime, endTime time.Time
	var err error

	if req.EndTime != "" {
		endTime, err = time.Parse(time.RFC3339, req.EndTime)
		if err != nil {
			logger.WithError(err).Error("Failed to parse end time")
			return &AkamaiTotalDNSRequestsResponse{
				Property: req.Property,
				Error:    "Failed to parse end time: " + err.Error(),
			}, nil
		}
	} else {
		endTime = time.Now().Add(-15 * time.Minute)
	}

	if req.StartTime != "" {
		startTime, err = time.Parse(time.RFC3339, req.StartTime)
		if err != nil {
			logger.WithError(err).Error("Failed to parse start time")
			return &AkamaiTotalDNSRequestsResponse{
				Property: req.Property,
				Error:    "Failed to parse start time: " + err.Error(),
			}, nil
		}
	} else {
		startTime = endTime.Add(-48 * time.Hour)
	}

	// Get Akamai configuration
	akamaiConfig := config.Global.AkamaiConfig
	if akamaiConfig.EdgeRC == "" {
		return &AkamaiTotalDNSRequestsResponse{
			Property: req.Property,
			Error:    "No EdgeRC configuration found for Akamai authentication",
		}, nil
	}

	// Initialize Akamai session
	logger.Info("Initializing Akamai session")
	session, err := createAkamaiSession(&akamaiConfig)
	if err != nil {
		logger.WithError(err).Error("Failed to create Akamai session")
		return &AkamaiTotalDNSRequestsResponse{
			Property: req.Property,
			Error:    "Failed to create Akamai session: " + err.Error(),
		}, nil
	}

	// Build the URL for direct API access to get data for the specific time range
	path := fmt.Sprintf("/gtm-api/v1/reports/traffic/domains/%s/properties/%s", domain, req.Property)
	params := url.Values{}
	params.Add("start", startTime.UTC().Format(time.RFC3339))
	params.Add("end", endTime.UTC().Format(time.RFC3339))
	uri := fmt.Sprintf("%s?%s", path, params.Encode())

	logger.Infof("Fetching traffic report: %s", uri)

	// Make the API request directly
	var trafficReport TrafficReport
	httpReq, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		logger.WithError(err).Error("Failed to create request")
		return &AkamaiTotalDNSRequestsResponse{
			Property:  req.Property,
			StartDate: startTime.Format(time.RFC3339),
			EndDate:   endTime.Format(time.RFC3339),
			Error:     "Failed to create request: " + err.Error(),
		}, nil
	}

	resp, err := (*session).Exec(httpReq, &trafficReport)
	if err != nil {
		logger.WithError(err).Error("Failed to execute request")
		return &AkamaiTotalDNSRequestsResponse{
			Property:  req.Property,
			StartDate: startTime.Format(time.RFC3339),
			EndDate:   endTime.Format(time.RFC3339),
			Error:     "Failed to get traffic report: " + err.Error(),
		}, nil
	}

	if resp.StatusCode != http.StatusOK {
		logger.WithField("status", resp.Status).Error("Unexpected status code")
		return &AkamaiTotalDNSRequestsResponse{
			Property:  req.Property,
			StartDate: startTime.Format(time.RFC3339),
			EndDate:   endTime.Format(time.RFC3339),
			Error:     fmt.Sprintf("Failed to get traffic report: unexpected status %d", resp.StatusCode),
		}, nil
	}

	// Process the traffic report data
	totalRequests := 0
	datacentersMap := make(map[string]*AkamaiDatacenterStats)

	// Process all data rows
	for _, dataRow := range trafficReport.DataRows {
		for _, datacenter := range dataRow.Datacenters {
			dcName := datacenter.Nickname
			requests := datacenter.Requests

			// Update total requests count
			totalRequests += requests

			// Update datacenter-specific counts
			if _, exists := datacentersMap[dcName]; !exists {
				datacentersMap[dcName] = &AkamaiDatacenterStats{
					DatacenterId:  dcName,
					TrafficTarget: datacenter.TrafficTargetName,
					TotalRequests: 0,
				}
			}
			datacentersMap[dcName].TotalRequests += int64(requests)
		}
	}

	// Calculate percentages for each datacenter
	if totalRequests > 0 {
		for dcName, dcData := range datacentersMap {
			percentage := float32(dcData.TotalRequests) / float32(totalRequests)
			dcData.Percentage = percentage
			datacentersMap[dcName] = dcData
		}
	}

	return &AkamaiTotalDNSRequestsResponse{
		Property:      req.Property,
		StartDate:     startTime.Format(time.RFC3339),
		EndDate:       endTime.Format(time.RFC3339),
		TotalRequests: int64(totalRequests),
		Datacenters:   datacentersMap,
	}, nil
}
