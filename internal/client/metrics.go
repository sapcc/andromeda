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
