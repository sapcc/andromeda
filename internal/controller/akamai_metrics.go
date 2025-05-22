// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/actatum/stormrpc"
	"github.com/apex/log"
	"github.com/go-openapi/runtime/middleware"
	lru "github.com/hashicorp/golang-lru/v2"

	"github.com/sapcc/andromeda/models"
	apiMetrics "github.com/sapcc/andromeda/restapi/operations/metrics"
)

// AkamaiMetricsController handles the metrics operations for Akamai
type AkamaiMetricsController struct {
	CommonController
	// Cache for DNS metrics to improve performance
	metricsCache *lru.Cache[string, *models.AkamaiTotalDNSRequests]
}

// Init initializes the controller with cache
func (c *AkamaiMetricsController) Init() {
	// Create a cache with 100 entries max
	cache, err := lru.New[string, *models.AkamaiTotalDNSRequests](100)
	if err != nil {
		log.WithError(err).Warn("Failed to create metrics cache, continuing without caching")
	} else {
		c.metricsCache = cache
		log.Info("Initialized metrics cache for DNS requests")
	}
}

// GetTotalDNSRequests retrieves total DNS request metrics for Akamai GTM properties
func (c AkamaiMetricsController) GetTotalDNSRequests(ctx context.Context, params apiMetrics.GetMetricsAkamaiTotalDNSRequestsParams) (*models.AkamaiTotalDNSRequests, error) {
	// Set default time range if not specified
	timeRange := "last_hour"
	if params.TimeRange != nil {
		timeRange = *params.TimeRange
	}

	// Create cache key based on parameters
	cacheKey := fmt.Sprintf("dns_metrics:%s:%s:%s",
		getStringParam(params.PropertyName),
		getStringParam(params.ProjectID),
		timeRange)

	// Check cache first
	if c.metricsCache != nil {
		if cachedResult, found := c.metricsCache.Get(cacheKey); found {
			log.WithField("cache_key", cacheKey).Debug("Cache hit for DNS metrics")
			return cachedResult, nil
		}
	}

	requestFields := log.Fields{
		"property_name": getStringParam(params.PropertyName),
		"project_id":    getStringParam(params.ProjectID),
		"time_range":    timeRange,
	}

	log.WithFields(requestFields).Debug("Fetching DNS metrics via RPC")

	// Prepare RPC request parameters
	type dnsMetricsRequest struct {
		PropertyName *string `json:"property_name,omitempty"`
		ProjectID    *string `json:"project_id,omitempty"`
		TimeRange    *string `json:"time_range,omitempty"`
	}

	rpcParams := dnsMetricsRequest{
		PropertyName: params.PropertyName,
		ProjectID:    params.ProjectID,
		TimeRange:    &timeRange,
	}

	// Create RPC request
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	rpcReq, err := stormrpc.NewRequest("andromeda.get_dns_metrics.akamai", rpcParams)
	if err != nil {
		log.WithError(err).WithFields(requestFields).Error("Failed to create RPC request")
		return nil, fmt.Errorf("error creating RPC request: %w", err)
	}

	// Send RPC request to Akamai agent
	rpcResp := c.rpc.Do(ctx, rpcReq)
	if rpcResp.Err != nil {
		log.WithError(rpcResp.Err).WithFields(requestFields).Error("RPC request to Akamai agent failed")
		return nil, fmt.Errorf("error from Akamai agent: %w", rpcResp.Err)
	}

	// Parse RPC response
	var result models.AkamaiTotalDNSRequests
	if err := rpcResp.Decode(&result); err != nil {
		log.WithError(err).WithFields(requestFields).Error("Failed to decode RPC response")
		return nil, fmt.Errorf("error decoding RPC response: %w", err)
	}

	// Validate the response
	if result.PropertyName == "" {
		log.WithFields(requestFields).Warn("RPC returned empty property name")
		return nil, fmt.Errorf("no properties found matching the criteria")
	}

	if result.TotalRequests == 0 && len(result.Datacenters) == 0 {
		log.WithFields(requestFields).Warn("RPC returned no DNS metrics data")
		return nil, fmt.Errorf("no DNS metrics data available for the specified criteria")
	}

	// Cache the result
	if c.metricsCache != nil {
		c.metricsCache.Add(cacheKey, &result)
		log.WithField("cache_key", cacheKey).Debug("Cached DNS metrics result")
	}

	log.WithFields(log.Fields{
		"property_name":      result.PropertyName,
		"requested_property": getStringParam(params.PropertyName),
		"total_requests":     result.TotalRequests,
		"datacenters":        len(result.Datacenters),
	}).Debug("Successfully retrieved DNS metrics")

	return &result, nil
}

// GetMetricsAkamaiTotalDNSRequestsHandler handles the GET /metrics/akamai/total-dns-requests endpoint
func (c AkamaiMetricsController) GetMetricsAkamaiTotalDNSRequestsHandler(params apiMetrics.GetMetricsAkamaiTotalDNSRequestsParams) middleware.Responder {
	result, err := c.GetTotalDNSRequests(params.HTTPRequest.Context(), params)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"property_name": getStringParam(params.PropertyName),
			"project_id":    getStringParam(params.ProjectID),
			"time_range":    getStringParamWithDefault(params.TimeRange, "last_hour"),
		}).Error("Failed to get DNS metrics")

		return apiMetrics.NewGetMetricsAkamaiTotalDNSRequestsBadRequest().WithPayload(
			&models.Error{
				Message: err.Error(),
				Code:    400,
			})
	}

	return apiMetrics.NewGetMetricsAkamaiTotalDNSRequestsOK().WithPayload(
		&apiMetrics.GetMetricsAkamaiTotalDNSRequestsOKBody{
			TotalDNSRequests: result,
		})
}

// Helper function to get string value from pointer or empty string
func getStringParam(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Helper function to get string value from pointer or default value
func getStringParamWithDefault(s *string, defaultValue string) string {
	if s == nil {
		return defaultValue
	}
	return *s
}
