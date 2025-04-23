package handlers

import (
	"time"

	"github.com/apex/log"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"

	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/metrics"
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

	// Initialize Akamai session
	akamaiConfig := metrics.AkamaiConfig{
		ClientToken:  config.Global.Metrics.Akamai.ClientToken,
		ClientSecret: config.Global.Metrics.Akamai.ClientSecret,
		AccessToken:  config.Global.Metrics.Akamai.AccessToken,
		Host:         config.Global.Metrics.Akamai.Host,
	}
	akamai, err := metrics.NewAkamaiSession(akamaiConfig)
	if err != nil {
		logger.WithError(err).Error("Failed to initialize Akamai session")
		errMsg := "Failed to initialize Akamai session"
		return metrics_ops.NewGetMetricsAkamaiTotalDNSRequestsOK().WithPayload(&models.AkamaiTotalDNSResult{
			Property:  params.Property,
			StartDate: startDateTime,
			EndDate:   endDateTime,
			Error:     &errMsg,
		})
	}

	// Get total DNS requests
	totalRequests, datacenterStats, err := akamai.GetTotalDNSRequests(domain, params.Property, startTime, endTime)
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

	// Convert datacenter stats to model
	datacentersMap := make(map[string]models.AkamaiDatacenterStats)
	for _, stat := range datacenterStats {
		datacentersMap[stat.DatacenterID] = models.AkamaiDatacenterStats{
			DatacenterID:  stat.DatacenterID,
			TrafficTarget: stat.TrafficTarget,
			TotalRequests: stat.TotalRequests,
			Percentage:    stat.Percentage,
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
