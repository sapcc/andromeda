// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"fmt"

	"github.com/actatum/stormrpc"
	"github.com/nats-io/nats.go"

	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/models"
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
	// Create metrics client
	client, err := NewMetricsClient()
	if err != nil {
		return fmt.Errorf("failed to create metrics client: %w", err)
	}

	var propName, projID, timeRange *string
	if MetricsOptions.MetricsAkamaiDNS.PropertyName != "" {
		propName = &MetricsOptions.MetricsAkamaiDNS.PropertyName
	}
	if MetricsOptions.MetricsAkamaiDNS.ProjectID != "" {
		projID = &MetricsOptions.MetricsAkamaiDNS.ProjectID
	}
	if MetricsOptions.MetricsAkamaiDNS.TimeRange != "" {
		timeRange = &MetricsOptions.MetricsAkamaiDNS.TimeRange
	}

	result, err := client.GetAkamaiTotalDNSRequests(context.Background(), propName, projID, timeRange)
	if err != nil {
		return fmt.Errorf("failed to get DNS metrics: %w", err)
	}

	// Write result to table
	return WriteTable(result)
}

// MetricsClient provides access to metrics-related functionality
type MetricsClient struct {
	rpc *stormrpc.Client
}

// NewMetricsClient creates a new client for metrics operations
func NewMetricsClient() (*MetricsClient, error) {
	// Connect to NATS
	nc, err := nats.Connect(config.Global.Default.TransportURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	// Create RPC client
	rpcClient, err := stormrpc.NewClient("", stormrpc.WithNatsConn(nc))
	if err != nil {
		return nil, fmt.Errorf("failed to create RPC client: %w", err)
	}

	return &MetricsClient{
		rpc: rpcClient,
	}, nil
}

// GetAkamaiTotalDNSRequests retrieves total DNS request metrics from Akamai GTM
func (c *MetricsClient) GetAkamaiTotalDNSRequests(ctx context.Context, propertyName, projectID, timeRange *string) (*models.AkamaiTotalDNSRequests, error) {
	// Prepare request parameters
	params := struct {
		PropertyName *string `json:"property_name,omitempty"`
		ProjectID    *string `json:"project_id,omitempty"`
		TimeRange    *string `json:"time_range,omitempty"`
	}{
		PropertyName: propertyName,
		ProjectID:    projectID,
		TimeRange:    timeRange,
	}

	// Create request
	req, err := stormrpc.NewRequest("andromeda.get_dns_metrics.akamai", params)
	if err != nil {
		return nil, fmt.Errorf("failed to create RPC request: %w", err)
	}

	// Send request
	resp := c.rpc.Do(ctx, req)
	if resp.Err != nil {
		return nil, fmt.Errorf("RPC request failed: %w", resp.Err)
	}

	// Parse response
	var result models.AkamaiTotalDNSRequests
	if err := resp.Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// Initialize the command
func init() {
	_, err := Parser.AddCommand("metrics", "Metrics", "Metrics Commands.", &MetricsOptions)
	if err != nil {
		panic(err)
	}
}
