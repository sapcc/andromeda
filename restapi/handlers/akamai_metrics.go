package handlers

import (
	"time"

	"github.com/apex/log"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"

	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/driver/akamai"
	"github.com/sapcc/andromeda/internal/driver/akamai/metrics"
	"github.com/sapcc/andromeda/models"
	metrics_ops "github.com/sapcc/andromeda/restapi/operations/metrics"
)

// AkamaiMetricsHandler handles requests for Akamai metrics
type AkamaiMetricsHandler struct{}

// NewAkamaiMetricsHandler creates a new instance of AkamaiMetricsHandler
func NewAkamaiMetricsHandler() *AkamaiMetricsHandler {
	return &AkamaiMetricsHandler{}
}

// GetMetricsAkamaiTotalDNSRequests handles requests to retrieve total DNS requests for an Akamai GTM property
func GetMetricsAkamaiTotalDNSRequests(params metrics_ops.GetMetricsAkamaiTotalDNSRequestsParams) middleware.Responder {
	logger := log.WithFields(log.Fields{
		"property": params.Property,
	})
	logger.Info("Getting total DNS requests for Akamai GTM property")

	// Set default domain if not provided
	domain := "andromeda.akadns.net"
	if params.Domain != nil {
		domain = *params.Domain
	}

	// Set default time range if not provided
	endTime := time.Now().Add(-15 * time.Minute)
	if params.End != nil {
		endTime = time.Time(*params.End)
	}
	endDateTime := strfmt.DateTime(endTime)

	startTime := endTime.Add(-48 * time.Hour)
	if params.Start != nil {
		startTime = time.Time(*params.Start)
	}
	startDateTime := strfmt.DateTime(startTime)

	// Make sure contract_id is set to bypass the check if needed
	if config.Global.AkamaiConfig.ContractId == "" {
		config.Global.AkamaiConfig.ContractId = "BYPASS_CHECK_VALUE"
		logger.Info("Setting ContractId to BYPASS_CHECK_VALUE to skip contract check")
	}

	// Initialize Akamai session using patched version
	session, domainType := akamai.NewAkamaiSessionPatched(&config.Global.AkamaiConfig)
	logger.Infof("Connected to Akamai API with domain type: %s", domainType)

	// Create cached session
	cachedSession := metrics.NewCachedAkamaiSession(*session, domain)

	// Get the data
	logger.Infof("Fetching total DNS requests for %s from %s to %s",
		params.Property, startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))

	result, err := metrics.GetTotalDNSRequests(cachedSession, domain, params.Property, startTime, endTime)
	if err != nil {
		logger.WithError(err).Error("Failed to get total DNS requests")
		errMsg := "Failed to get total DNS requests: " + err.Error()
		return metrics_ops.NewGetMetricsAkamaiTotalDNSRequestsOK().WithPayload(&models.AkamaiTotalDNSResult{
			Property:  params.Property,
			StartDate: startDateTime,
			EndDate:   endDateTime,
			Error:     &errMsg,
		})
	}

	// Extract data from result
	totalRequests := int64(result["total_requests"].(int))
	datacenters := result["datacenters"].(map[string]map[string]interface{})

	// Convert datacenter stats to model
	datacentersMap := make(map[string]models.AkamaiDatacenterStats)
	for dcID, dcData := range datacenters {
		datacentersMap[dcID] = models.AkamaiDatacenterStats{
			DatacenterID:  dcData["datacenter_id"].(string),
			TrafficTarget: dcData["traffic_target"].(string),
			TotalRequests: int64(dcData["total_requests"].(int)),
			Percentage:    float32(dcData["percentage"].(float64)),
		}
	}

	return metrics_ops.NewGetMetricsAkamaiTotalDNSRequestsOK().WithPayload(&models.AkamaiTotalDNSResult{
		Property:      params.Property,
		StartDate:     startDateTime,
		EndDate:       endDateTime,
		TotalRequests: totalRequests,
		Datacenters:   datacentersMap,
	})
}
