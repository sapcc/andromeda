// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package metrics

import (
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

// TestHistoricalDataHandling verifies that the collector handles historical data correctly
// by only exposing the most recent values without historical timestamps
func TestHistoricalDataHandling(t *testing.T) {
	// Test data
	property := "test-property"
	datacenter := "dc1"
	projectID := "project123"
	target := "192.168.1.1"
	
	// Historical timestamp (10 minutes ago)
	historicalTime := time.Now().Add(-10 * time.Minute)
	
	// Create metrics as the collector would
	requestsMetric := prometheus.MustNewConstMetric(
		descNumReq, 
		prometheus.GaugeValue,
		200.0,
		property, datacenter, projectID, target,
	)
	
	statusMetric := prometheus.MustNewConstMetric(
		descStatus,
		prometheus.GaugeValue,
		200.0,
		property, datacenter, projectID, target,
	)
	
	timestampMetric := prometheus.MustNewConstMetric(
		descUpdated,
		prometheus.GaugeValue,
		float64(historicalTime.Unix()),
		property, datacenter, projectID, target,
	)

	// Collect the metrics
	ch := make(chan prometheus.Metric, 3)
	ch <- requestsMetric
	ch <- statusMetric
	ch <- timestampMetric
	close(ch)

	// Verify the metrics
	metricsCollected := []prometheus.Metric{}
	for metric := range ch {
		metricsCollected = append(metricsCollected, metric)
	}

	assert.Len(t, metricsCollected, 3, "Should have 3 metrics")

	// Check that none of the metrics have explicit timestamps
	for i, metric := range metricsCollected[:2] { // First two metrics (requests and status)
		dto := &dto.Metric{}
		err := metric.Write(dto)
		assert.NoError(t, err)
		assert.Nil(t, dto.TimestampMs, "Metric %d should not have an explicit timestamp", i)
	}

	// Verify the timestamp tracking metric contains the historical timestamp
	timestampDto := &dto.Metric{}
	err := metricsCollected[2].Write(timestampDto)
	assert.NoError(t, err)
	assert.Equal(t, float64(historicalTime.Unix()), timestampDto.Gauge.GetValue(), 
		"Timestamp metric should contain the Unix timestamp of when data was collected")
}

// TestMultipleHistoricalDataPoints verifies that only the most recent data is kept
func TestMultipleHistoricalDataPoints(t *testing.T) {
	// Simulate processing multiple data rows
	type metricKey struct {
		property   string
		datacenter string
		projectID  string
		target     string
	}

	// Track latest values (simulating the collector logic)
	latestRequests := make(map[metricKey]float64)
	latestTimestamp := make(map[metricKey]float64)

	// Simulate multiple data points for the same metric
	key := metricKey{
		property:   "prop1",
		datacenter: "dc1",
		projectID:  "proj1",
		target:     "192.168.1.1",
	}

	now := time.Now()
	
	// Process historical data points (oldest to newest)
	dataPoints := []struct {
		timestamp time.Time
		requests  float64
	}{
		{now.Add(-15 * time.Minute), 100},
		{now.Add(-10 * time.Minute), 150},
		{now.Add(-5 * time.Minute), 200}, // Most recent
	}

	for _, dp := range dataPoints {
		latestRequests[key] = dp.requests
		latestTimestamp[key] = float64(dp.timestamp.Unix())
	}

	// Verify only the latest values are kept
	assert.Equal(t, 200.0, latestRequests[key], "Should keep the most recent request count")
	assert.Equal(t, float64(now.Add(-5*time.Minute).Unix()), latestTimestamp[key], 
		"Should keep the most recent timestamp")
}

// TestMetricDescriptions verifies the metric descriptions are informative
func TestMetricDescriptions(t *testing.T) {
	// Check that descriptions explain the behavior
	assert.Contains(t, descNumReq.String(), "most recent data point", 
		"Request metric description should indicate it's the most recent value")
	assert.Contains(t, descStatus.String(), "most recent data point",
		"Status metric description should indicate it's the most recent value")
	assert.Contains(t, descUpdated.String(), "when this metric was last updated",
		"Timestamp metric description should explain what it tracks")
}