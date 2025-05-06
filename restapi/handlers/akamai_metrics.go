package handlers

import (
	"context"
	"time"

	"github.com/apex/log"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"

	"github.com/sapcc/andromeda/internal/rpc"
	"github.com/sapcc/andromeda/internal/rpc/server"
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

	// Get RPC client
	_, err := rpc.GetRPCClient()
	if err != nil {
		logger.WithError(err).Error("Failed to get RPC client")
		errMsg := "Failed to get RPC client: " + err.Error()
		return metrics_ops.NewGetMetricsAkamaiTotalDNSRequestsOK().WithPayload(&models.AkamaiTotalDNSResult{
			Property:  params.Property,
			StartDate: startDateTime,
			EndDate:   endDateTime,
			Error:     &errMsg,
		})
	}

	// Create RPC handler directly since we can't use the generated client (which doesn't have the method yet)
	handler := &server.RPCHandler{}

	// Prepare the request to the RPC service
	rpcRequest := &server.AkamaiTotalDNSRequestsRequest{
		Domain:    domain,
		Property:  params.Property,
		StartTime: startTime.Format(time.RFC3339),
		EndTime:   endTime.Format(time.RFC3339),
	}

	// Call the RPC method directly
	resp, err := handler.GetAkamaiTotalDNSRequests(context.Background(), rpcRequest)
	if err != nil {
		logger.WithError(err).Error("RPC call failed")
		errMsg := "RPC call failed: " + err.Error()
		return metrics_ops.NewGetMetricsAkamaiTotalDNSRequestsOK().WithPayload(&models.AkamaiTotalDNSResult{
			Property:  params.Property,
			StartDate: startDateTime,
			EndDate:   endDateTime,
			Error:     &errMsg,
		})
	}

	// Check for errors in the response
	if resp.Error != "" {
		logger.Error("Error in RPC response: " + resp.Error)
		errMsg := resp.Error
		return metrics_ops.NewGetMetricsAkamaiTotalDNSRequestsOK().WithPayload(&models.AkamaiTotalDNSResult{
			Property:  params.Property,
			StartDate: startDateTime,
			EndDate:   endDateTime,
			Error:     &errMsg,
		})
	}

	// Convert datacenter stats to model
	datacentersMap := make(map[string]models.AkamaiDatacenterStats)
	for dcID, dcData := range resp.Datacenters {
		datacentersMap[dcID] = models.AkamaiDatacenterStats{
			DatacenterID:  dcData.DatacenterId,
			TrafficTarget: dcData.TrafficTarget,
			TotalRequests: dcData.TotalRequests,
			Percentage:    float32(dcData.Percentage),
		}
	}

	return metrics_ops.NewGetMetricsAkamaiTotalDNSRequestsOK().WithPayload(&models.AkamaiTotalDNSResult{
		Property:      params.Property,
		StartDate:     startDateTime,
		EndDate:       endDateTime,
		TotalRequests: resp.TotalRequests,
		Datacenters:   datacentersMap,
	})
}
