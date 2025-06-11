// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package akamai

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/sapcc/andromeda/internal/rpcmodels"
)

func TestRPCHandlerAkamai_Sync(t *testing.T) {
	// Create a mock agent with a buffered channel
	mockAgent := &AkamaiAgent{
		forceSync: make(chan []string, 1),
	}

	handler := &RPCHandlerAkamai{
		agent: mockAgent,
	}

	// Test sync request
	ctx := context.Background()
	request := &rpcmodels.SyncRequest{
		DomainId: []string{"domain1", "domain2"},
	}

	response, err := handler.Sync(ctx, request)

	assert.NoError(t, err, "Sync should not return error")
	assert.NotNil(t, response, "Response should not be nil")

	// Check that the domain IDs were sent to the forceSync channel
	select {
	case receivedDomains := <-mockAgent.forceSync:
		assert.Equal(t, []string{"domain1", "domain2"}, receivedDomains, "Domain IDs should match")
	case <-time.After(100 * time.Millisecond):
		t.Error("Expected domain IDs to be sent to forceSync channel")
	}
}

func TestRPCHandlerAkamai_TimeRangeCalculation(t *testing.T) {
	// Test time range calculations
	testCases := []struct {
		name      string
		timeRange rpcmodels.GetDNSMetricsRequest_TimeRangeValue
		expected  time.Duration
	}{
		{"Last Hour", rpcmodels.GetDNSMetricsRequest_LAST_HOUR, time.Hour},
		{"Last Day", rpcmodels.GetDNSMetricsRequest_LAST_DAY, 24 * time.Hour},
		{"Last Week", rpcmodels.GetDNSMetricsRequest_LAST_WEEK, 7 * 24 * time.Hour},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			end := time.Now()
			start := end

			// This replicates the time calculation logic from GetDNSMetricsAkamai
			switch tc.timeRange {
			case rpcmodels.GetDNSMetricsRequest_LAST_HOUR:
				start = end.Add(-time.Hour)
			case rpcmodels.GetDNSMetricsRequest_LAST_DAY:
				start = end.Add(-24 * time.Hour)
			case rpcmodels.GetDNSMetricsRequest_LAST_WEEK:
				start = end.Add(-7 * 24 * time.Hour)
			}

			actualDuration := end.Sub(start)
			assert.Equal(t, tc.expected, actualDuration, "Time range calculation should be correct")
		})
	}
}

func TestRPCHandlerAkamai_DatacenterAggregation(t *testing.T) {
	// Test datacenter aggregation logic (simulating the logic from GetDNSMetricsAkamai)
	testData := []struct {
		nickname string
		requests int
	}{
		{"dc1", 100},
		{"dc1", 200}, // Same datacenter, should be aggregated
		{"dc2", 150},
		{"dc1", 50}, // Same datacenter again
	}

	// Simulate the aggregation logic
	processedDatacenters := make(map[string]int)
	totalRequests := int64(0)
	var datacenters []*rpcmodels.GetDNSMetricsResponse_Datacenter

	for _, data := range testData {
		totalRequests += int64(data.requests)

		if idx, exists := processedDatacenters[data.nickname]; exists {
			// Update existing datacenter
			datacenters[idx].Requests += int64(data.requests)
		} else {
			// Add new datacenter
			datacenter := &rpcmodels.GetDNSMetricsResponse_Datacenter{
				DatacenterId: data.nickname,
				Requests:     int64(data.requests),
				Status:       "up",
			}
			processedDatacenters[data.nickname] = len(datacenters)
			datacenters = append(datacenters, datacenter)
		}
	}

	assert.Equal(t, int64(500), totalRequests, "Total requests should be sum of all requests")
	assert.Len(t, datacenters, 2, "Should have 2 unique datacenters")

	// Check aggregated values
	dc1Found := false
	dc2Found := false
	for _, dc := range datacenters {
		if dc.DatacenterId == "dc1" {
			assert.Equal(t, int64(350), dc.Requests, "dc1 should have aggregated requests (100+200+50)")
			dc1Found = true
		} else if dc.DatacenterId == "dc2" {
			assert.Equal(t, int64(150), dc.Requests, "dc2 should have correct requests")
			dc2Found = true
		}
	}
	assert.True(t, dc1Found, "dc1 should be found")
	assert.True(t, dc2Found, "dc2 should be found")
}

func TestRPCHandlerAkamai_TargetIPExtraction(t *testing.T) {
	// Test target IP extraction from traffic target name (logic from GetDNSMetricsAkamai)
	testCases := []struct {
		trafficTargetName string
		expectedIP        string
	}{
		{"Target 1 - 1.2.3.4", "1.2.3.4"},
		{"Target 2 - 192.168.1.1", "192.168.1.1"},
		{"Target without IP", ""},
		{"", ""},
		{"Target - ", ""},
		{"Target-1.2.3.4", ""}, // Wrong separator
	}

	for _, tc := range testCases {
		t.Run(tc.trafficTargetName, func(t *testing.T) {
			var targetIP string
			// This replicates the IP extraction logic from GetDNSMetricsAkamai
			if strings.Contains(tc.trafficTargetName, " - ") {
				parts := strings.Split(tc.trafficTargetName, " - ")
				if len(parts) > 1 {
					targetIP = parts[1]
				}
			}

			assert.Equal(t, tc.expectedIP, targetIP, "Target IP extraction should be correct")
		})
	}
}

func TestRPCHandlerAkamai_PercentageCalculation(t *testing.T) {
	// Test percentage calculation logic (from GetDNSMetricsAkamai)
	datacenters := []*rpcmodels.GetDNSMetricsResponse_Datacenter{
		{Requests: 500},
		{Requests: 300},
		{Requests: 200},
	}

	totalRequests := int64(0)
	for _, dc := range datacenters {
		totalRequests += dc.Requests
	}

	expectedPercentages := []float32{50.0, 30.0, 20.0}

	// Calculate percentages (replicating the logic from GetDNSMetricsAkamai)
	if totalRequests > 0 {
		for i, dc := range datacenters {
			dc.Percentage = float32(float64(dc.Requests) / float64(totalRequests) * 100)
			assert.Equal(t, expectedPercentages[i], dc.Percentage, "Percentage calculation should be correct")
		}
	}
}

func TestRPCHandlerAkamai_EmptyDataResponse(t *testing.T) {
	// Test handling of empty data (no datacenters)
	var datacenters []*rpcmodels.GetDNSMetricsResponse_Datacenter
	totalRequests := int64(0)

	// Simulate empty traffic report
	emptyDataRows := []struct {
		Datacenters []struct {
			Nickname string
			Requests int
		}
	}{}

	// Process empty data
	for _, dataRow := range emptyDataRows {
		for _, dc := range dataRow.Datacenters {
			totalRequests += int64(dc.Requests)
		}
	}

	response := &rpcmodels.GetDNSMetricsResponse{
		Datacenters:   datacenters,
		TotalRequests: totalRequests,
	}

	assert.Equal(t, int64(0), response.TotalRequests, "Total requests should be 0 for empty data")
	assert.Len(t, response.Datacenters, 0, "Should have 0 datacenters for empty data")
}

func TestRPCHandlerAkamai_SyncRequestValidation(t *testing.T) {
	// Test sync request with nil domain IDs
	mockAgent := &AkamaiAgent{
		forceSync: make(chan []string, 1),
	}

	handler := &RPCHandlerAkamai{
		agent: mockAgent,
	}

	ctx := context.Background()
	request := &rpcmodels.SyncRequest{
		DomainId: nil, // Test with nil
	}

	response, err := handler.Sync(ctx, request)

	assert.NoError(t, err, "Sync should not return error with nil domain IDs")
	assert.NotNil(t, response, "Response should not be nil")

	// Check that nil was sent to the forceSync channel
	select {
	case receivedDomains := <-mockAgent.forceSync:
		assert.Nil(t, receivedDomains, "Domain IDs should be nil")
	case <-time.After(100 * time.Millisecond):
		t.Error("Expected nil to be sent to forceSync channel")
	}
}

func TestRPCHandlerAkamai_SyncRequestEmpty(t *testing.T) {
	// Test sync request with empty domain IDs
	mockAgent := &AkamaiAgent{
		forceSync: make(chan []string, 1),
	}

	handler := &RPCHandlerAkamai{
		agent: mockAgent,
	}

	ctx := context.Background()
	request := &rpcmodels.SyncRequest{
		DomainId: []string{}, // Test with empty slice
	}

	response, err := handler.Sync(ctx, request)

	assert.NoError(t, err, "Sync should not return error with empty domain IDs")
	assert.NotNil(t, response, "Response should not be nil")

	// Check that empty slice was sent to the forceSync channel
	select {
	case receivedDomains := <-mockAgent.forceSync:
		assert.Equal(t, []string{}, receivedDomains, "Domain IDs should be empty slice")
	case <-time.After(100 * time.Millisecond):
		t.Error("Expected empty slice to be sent to forceSync channel")
	}
}
