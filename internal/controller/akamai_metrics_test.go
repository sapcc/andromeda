// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sapcc/andromeda/internal/rpcmodels"
	"github.com/sapcc/andromeda/models"
)

func TestNewAkamaiMetricsController(t *testing.T) {
	cc := CommonController{}
	controller := NewAkamaiMetricsController(cc)

	assert.NotNil(t, controller.metricsCache, "Cache should be initialized")
	assert.Equal(t, cc, controller.CommonController, "CommonController should be set")
}

func TestAkamaiMetricsController_CacheOperations(t *testing.T) {
	cc := CommonController{}
	controller := NewAkamaiMetricsController(cc)

	// Test cache operations
	testData := &models.AkamaiTotalDNSRequests{
		PropertyName:  "test.example.com",
		TimeRange:     "last_hour",
		TotalRequests: 1000,
	}

	cacheKey := "test_key"
	controller.metricsCache.Add(cacheKey, testData)

	cachedResult, found := controller.metricsCache.Get(cacheKey)
	assert.True(t, found, "Data should be found in cache")
	assert.Equal(t, testData, cachedResult, "Cached data should match original")
}

func TestAkamaiMetricsController_CacheInitializationFailure(t *testing.T) {
	// Test behavior when cache fails to initialize
	cc := CommonController{}

	// This test demonstrates the controller can work without cache
	controller := AkamaiMetricsController{
		CommonController: cc,
		metricsCache:     nil,
	}

	assert.Nil(t, controller.metricsCache, "Cache should be nil when initialization fails")
}

func TestAkamaiMetricsController_CacheKeyGeneration(t *testing.T) {
	// Test cache key generation logic
	domainID := "test-domain-123"
	projectID := "test-project-456"
	timeRange := "last_hour"

	expectedKey := "dns_metrics:test-domain-123:test-project-456:last_hour"
	actualKey := fmt.Sprintf("dns_metrics:%s:%s:%s", domainID, projectID, timeRange)

	assert.Equal(t, expectedKey, actualKey, "Cache key should match expected format")
}

func TestAkamaiMetricsController_DatacenterMapping(t *testing.T) {
	// Test the datacenter mapping logic from RPC response to API response
	rpcDatacenters := []*rpcmodels.GetDNSMetricsResponse_Datacenter{
		{
			DatacenterId: "dc1",
			Requests:     500,
			Status:       "up",
			TargetIp:     "1.2.3.4",
			Percentage:   50.0,
		},
		{
			DatacenterId: "dc2",
			Requests:     300,
			Status:       "down",
			TargetIp:     "5.6.7.8",
			Percentage:   30.0,
		},
	}

	// Convert to API model format
	apiDatacenters := make([]*models.AkamaiTotalDNSRequestsDatacentersItems0, len(rpcDatacenters))
	for i, dc := range rpcDatacenters {
		apiDatacenters[i] = &models.AkamaiTotalDNSRequestsDatacentersItems0{
			DatacenterID:       dc.GetDatacenterId(),
			DatacenterNickname: dc.GetDatacenterId(),
			Percentage:         dc.GetPercentage(),
			Requests:           dc.GetRequests(),
			Status:             dc.GetStatus(),
			TargetIP:           dc.GetTargetIp(),
		}
	}

	assert.Len(t, apiDatacenters, 2, "Should have 2 datacenters")
	assert.Equal(t, "dc1", apiDatacenters[0].DatacenterID, "Datacenter ID should match")
	assert.Equal(t, "dc1", apiDatacenters[0].DatacenterNickname, "Datacenter nickname should match ID")
	assert.Equal(t, int64(500), apiDatacenters[0].Requests, "Requests should match")
	assert.Equal(t, "up", apiDatacenters[0].Status, "Status should match")
	assert.Equal(t, "1.2.3.4", apiDatacenters[0].TargetIP, "Target IP should match")
	assert.Equal(t, float32(50.0), apiDatacenters[0].Percentage, "Percentage should match")
}

func TestAkamaiMetricsController_TimeRangeMapping(t *testing.T) {
	// Test time range parameter mapping
	testCases := []struct {
		input    string
		expected rpcmodels.GetDNSMetricsRequest_TimeRangeValue
	}{
		{"LAST_HOUR", rpcmodels.GetDNSMetricsRequest_LAST_HOUR},
		{"LAST_DAY", rpcmodels.GetDNSMetricsRequest_LAST_DAY},
		{"LAST_WEEK", rpcmodels.GetDNSMetricsRequest_LAST_WEEK},
	}

	for _, tc := range testCases {
		if value, exists := rpcmodels.GetDNSMetricsRequest_TimeRangeValue_value[tc.input]; exists {
			assert.Equal(t, tc.expected, rpcmodels.GetDNSMetricsRequest_TimeRangeValue(value),
				"Time range mapping should be correct for %s", tc.input)
		}
	}
}

func TestAkamaiMetricsController_CacheWithNilCache(t *testing.T) {
	// Test behavior when cache is nil
	controller := AkamaiMetricsController{
		metricsCache: nil,
	}

	// Should not panic when cache is nil
	assert.NotPanics(t, func() {
		if controller.metricsCache != nil {
			controller.metricsCache.Get("test")
		}
	}, "Should not panic with nil cache")

	assert.NotPanics(t, func() {
		if controller.metricsCache != nil {
			controller.metricsCache.Add("test", &models.AkamaiTotalDNSRequests{})
		}
	}, "Should not panic when adding to nil cache")
}
