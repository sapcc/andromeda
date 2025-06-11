// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-openapi/runtime/middleware"

	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/driver/akamai"
	akamaiMetrics "github.com/sapcc/andromeda/internal/driver/akamai/metrics"
	"github.com/sapcc/andromeda/models"
	apiMetrics "github.com/sapcc/andromeda/restapi/operations/metrics"
)

// AkamaiMetricsController handles the metrics operations for Akamai
type AkamaiMetricsController struct {
	CommonController
}

// GetTotalDNSRequests retrieves total DNS request metrics for Akamai GTM properties
func (c AkamaiMetricsController) GetTotalDNSRequests(ctx context.Context, params apiMetrics.GetMetricsAkamaiTotalDNSRequestsParams) (*models.AkamaiTotalDNSRequests, error) {
	// Create a session with Akamai
	session, _ := akamai.NewAkamaiSession(&config.Global.AkamaiConfig)

	// Set default time range if not specified
	timeRange := "last_hour"
	if params.TimeRange != nil {
		timeRange = *params.TimeRange
	}

	// Create cached session for metrics
	cachedSession := akamaiMetrics.NewCachedAkamaiSession(*session, config.Global.AkamaiConfig.Domain)

	// Get properties
	properties, err := cachedSession.GetProperties()
	if err != nil {
		return nil, fmt.Errorf("error fetching Akamai properties: %w", err)
	}

	// Filter properties by name if specified
	if params.PropertyName != nil && *params.PropertyName != "" {
		var filteredProperties []string
		for _, prop := range properties {
			if strings.Contains(prop, *params.PropertyName) {
				filteredProperties = append(filteredProperties, prop)
			}
		}
		if len(filteredProperties) > 0 {
			properties = filteredProperties
		}
	}

	// If no properties found, return empty result
	if len(properties) == 0 {
		return nil, fmt.Errorf("no Akamai GTM properties found matching the criteria")
	}

	// Use first property for initial implementation
	property := properties[0]

	// Get traffic report for the property
	trafficData, err := cachedSession.GetTrafficReport(property)
	if err != nil {
		return nil, fmt.Errorf("error fetching traffic data: %w", err)
	}

	// Build response
	totalRequests := int64(0)
	datacenters := []*models.AkamaiTotalDNSRequestsDatacentersItems0{}

	// Process traffic data
	for _, dataRow := range trafficData {
		for _, dc := range dataRow.Datacenters {
			totalRequests += int64(dc.Requests)

			// Add datacenter to response
			datacenter := &models.AkamaiTotalDNSRequestsDatacentersItems0{
				DatacenterID:       dc.Nickname,
				DatacenterNickname: dc.Nickname,
				Requests:           int64(dc.Requests),
				Status:             dc.Status,
			}

			// Extract target IP from traffic target name
			if strings.Contains(dc.TrafficTargetName, " - ") {
				parts := strings.Split(dc.TrafficTargetName, " - ")
				if len(parts) > 1 {
					datacenter.TargetIP = parts[1]
				}
			}

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

	return &models.AkamaiTotalDNSRequests{
		PropertyName:  property,
		TimeRange:     timeRange,
		TotalRequests: totalRequests,
		Datacenters:   datacenters,
	}, nil
}

// GetMetricsAkamaiTotalDNSRequestsHandler handles the GET /metrics/akamai/total-dns-requests endpoint
func (c AkamaiMetricsController) GetMetricsAkamaiTotalDNSRequestsHandler(params apiMetrics.GetMetricsAkamaiTotalDNSRequestsParams) middleware.Responder {
	result, err := c.GetTotalDNSRequests(params.HTTPRequest.Context(), params)
	if err != nil {
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
