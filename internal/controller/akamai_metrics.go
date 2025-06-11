// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	lru "github.com/hashicorp/golang-lru/v2"

	"github.com/sapcc/andromeda/internal/auth"
	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/rpc/agent/akamai"
	"github.com/sapcc/andromeda/internal/rpcmodels"
	"github.com/sapcc/andromeda/internal/utils"
	"github.com/sapcc/andromeda/models"
	apiMetrics "github.com/sapcc/andromeda/restapi/operations/metrics"
)

// AkamaiMetricsController handles the metrics operations for Akamai
type AkamaiMetricsController struct {
	CommonController
	// Cache for DNS metrics to improve performance
	metricsCache *lru.Cache[string, *models.AkamaiTotalDNSRequests]
}

// NewAkamaiMetricsController initializes the controller with cache
func NewAkamaiMetricsController(cc CommonController) AkamaiMetricsController {
	// Create a cache with 100 entries max
	cache, err := lru.New[string, *models.AkamaiTotalDNSRequests](100)
	if err != nil {
		log.WithError(err).Warn("Failed to create metrics cache, continuing without caching")
	} else {
		log.Info("Initialized metrics cache for DNS requests")
	}

	return AkamaiMetricsController{
		CommonController: cc,
		metricsCache:     cache,
	}
}

// GetMetricsAkamaiTotalDNSRequestsHandler handles the GET /metrics/akamai/total-dns-requests endpoint
func (c AkamaiMetricsController) GetMetricsAkamaiTotalDNSRequestsHandler(params apiMetrics.GetMetricsAkamaiTotalDNSRequestsParams) middleware.Responder {
	projectID, err := auth.Authenticate(params.HTTPRequest, nil)
	if err != nil {
		return apiMetrics.NewGetMetricsAkamaiTotalDNSRequestsDefault(403).WithPayload(utils.PolicyForbidden)
	} else if projectID == "" && params.ProjectID != nil {
		projectID = *params.ProjectID
	}

	domain := models.Domain{ID: params.DomainID, Pools: []strfmt.UUID{}}
	{
		// Resolve domain object to option FQDN as it's the property name in Akamai

		// resolve domain
		if err := PopulateDomain(c.db, &domain, []string{"*"}); err != nil {
			return apiMetrics.NewGetMetricsAkamaiTotalDNSRequestsDefault(404).WithPayload(utils.NotFound)
		}

		// This validates the domain's project ID against the user's scoped project ID
		requestVars := map[string]string{"project_id": *domain.ProjectID}
		if _, err := auth.Authenticate(params.HTTPRequest, requestVars); err != nil {
			return apiMetrics.NewGetMetricsAkamaiTotalDNSRequestsDefault(403).WithPayload(utils.PolicyForbidden)
		}
	}

	// Set default time range if not specified
	if params.TimeRange == nil {
		params.TimeRange = swag.String("last_hour")
	}

	// Create cache key based on parameters
	cacheKey := fmt.Sprintf("dns_metrics:%s:%s:%v",
		params.DomainID,
		projectID,
		params.TimeRange)

	// Check cache first
	if c.metricsCache != nil {
		if cachedResult, found := c.metricsCache.Get(cacheKey); found {
			log.WithField("cache_key", cacheKey).Debug("Cache hit for DNS metrics")
			return apiMetrics.NewGetMetricsAkamaiTotalDNSRequestsOK().WithPayload(
				&apiMetrics.GetMetricsAkamaiTotalDNSRequestsOKBody{
					TotalDNSRequests: cachedResult,
				})
		}
	}

	requestFields := log.Fields{
		"domain_id":  params.DomainID,
		"project_id": projectID,
	}
	log.WithFields(requestFields).Debug("Fetching DNS metrics via RPC")

	timeRange := rpcmodels.GetDNSMetricsRequest_TimeRangeValue_value[strings.ToUpper(*params.TimeRange)]
	req := &rpcmodels.GetDNSMetricsRequest{
		Domain:    config.Global.AkamaiConfig.Domain,
		Property:  domain.Fqdn.String(),
		TimeRange: rpcmodels.GetDNSMetricsRequest_TimeRangeValue(timeRange),
	}

	// Create RPC request
	ctx, cancel := context.WithTimeout(params.HTTPRequest.Context(), 30*time.Second)
	defer cancel()

	client := akamai.NewRPCAgentAkamaiClient(c.rpc)
	akamaiDNSMetrics, err := client.GetDNSMetricsAkamai(ctx, req)
	if err != nil {
		log.WithError(err).WithFields(requestFields).Error("Failed to fetch DNS metrics via RPC")
		return apiMetrics.
			NewGetMetricsAkamaiTotalDNSRequestsDefault(408).
			WithPayload(utils.TryAgainLater)
	}

	datacenters := make([]*models.AkamaiTotalDNSRequestsDatacentersItems0, len(akamaiDNSMetrics.Datacenters))
	for i, dc := range akamaiDNSMetrics.Datacenters {
		datacenters[i] = &models.AkamaiTotalDNSRequestsDatacentersItems0{
			DatacenterID:       dc.GetDatacenterId(),
			DatacenterNickname: dc.GetDatacenterId(),
			Percentage:         dc.GetPercentage(),
			Requests:           dc.GetRequests(),
			Status:             dc.GetStatus(),
			TargetIP:           dc.GetTargetIp(),
		}
	}

	result := &models.AkamaiTotalDNSRequests{
		Datacenters:   datacenters,
		PropertyName:  domain.Fqdn.String(),
		TimeRange:     *params.TimeRange,
		TotalRequests: akamaiDNSMetrics.GetTotalRequests(),
	}

	// Cache the result
	if c.metricsCache != nil {
		c.metricsCache.Add(cacheKey, result)
		log.WithField("cache_key", cacheKey).Debug("Cached DNS metrics result")
	}

	return apiMetrics.NewGetMetricsAkamaiTotalDNSRequestsOK().WithPayload(
		&apiMetrics.GetMetricsAkamaiTotalDNSRequestsOKBody{
			TotalDNSRequests: result,
		})
}
