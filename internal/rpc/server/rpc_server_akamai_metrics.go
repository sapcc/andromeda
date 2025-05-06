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
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

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

	// Check for .edgerc file
	edgercPath := config.Global.AkamaiConfig.EdgeRC
	if edgercPath == "" {
		// Try common locations
		homeDir, err := os.UserHomeDir()
		if err == nil {
			candidates := []string{
				".edgerc",                         // Current directory
				filepath.Join(homeDir, ".edgerc"), // Home directory
			}

			for _, path := range candidates {
				if _, err := os.Stat(path); err == nil {
					edgercPath = path
					break
				}
			}
		}
	}

	if edgercPath == "" {
		return &AkamaiTotalDNSRequestsResponse{
			Property: req.Property,
			Error:    "No .edgerc file found for Akamai authentication",
		}, nil
	}

	// Use the CLI tool to avoid import cycles
	cmdPath, err := exec.LookPath("./build/andromeda-akamai-total-dns-requests")
	if err != nil {
		// Try to find it in the current directory
		if _, err := os.Stat("./build/andromeda-akamai-total-dns-requests"); err == nil {
			cmdPath = "./build/andromeda-akamai-total-dns-requests"
		} else {
			return &AkamaiTotalDNSRequestsResponse{
				Property: req.Property,
				Error:    "Could not find andromeda-akamai-total-dns-requests tool: " + err.Error(),
			}, nil
		}
	}

	// Prepare CLI arguments
	args := []string{
		"--domain", domain,
		"--property", req.Property,
		"--start", startTime.Format(time.RFC3339),
		"--end", endTime.Format(time.RFC3339),
		"--output", "json",
		"--edgerc", edgercPath,
	}

	logger.Infof("Executing: %s %s", cmdPath, strings.Join(args, " "))

	// Run the CLI tool
	cmd := exec.Command(cmdPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.WithError(err).Error("Failed to execute CLI tool")
		return &AkamaiTotalDNSRequestsResponse{
			Property:  req.Property,
			StartDate: startTime.Format(time.RFC3339),
			EndDate:   endTime.Format(time.RFC3339),
			Error:     "Failed to execute CLI tool: " + err.Error() + " - Output: " + string(output),
		}, nil
	}

	// Parse the JSON output
	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		logger.WithError(err).Error("Failed to parse CLI output")
		return &AkamaiTotalDNSRequestsResponse{
			Property:  req.Property,
			StartDate: startTime.Format(time.RFC3339),
			EndDate:   endTime.Format(time.RFC3339),
			Error:     "Failed to parse CLI output: " + err.Error() + " - Output: " + string(output),
		}, nil
	}

	// Extract data from result
	totalRequests := int64(result["total_requests"].(float64))
	datacenters := result["datacenters"].(map[string]interface{})

	// Convert datacenter stats to model
	datacentersMap := make(map[string]*AkamaiDatacenterStats)
	for dcID, dcData := range datacenters {
		dcMap := dcData.(map[string]interface{})
		datacentersMap[dcID] = &AkamaiDatacenterStats{
			DatacenterId:  dcMap["datacenter_id"].(string),
			TrafficTarget: dcMap["traffic_target"].(string),
			TotalRequests: int64(dcMap["total_requests"].(float64)),
			Percentage:    float32(dcMap["percentage"].(float64)),
		}
	}

	return &AkamaiTotalDNSRequestsResponse{
		Property:      req.Property,
		StartDate:     result["start_date"].(string),
		EndDate:       result["end_date"].(string),
		TotalRequests: totalRequests,
		Datacenters:   datacentersMap,
	}, nil
}
