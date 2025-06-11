// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package akamai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/actatum/stormrpc"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/gtm"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	"github.com/apex/log"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/nats-io/nats.go"

	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/rpc"
	"github.com/sapcc/andromeda/internal/rpc/server"
	"github.com/sapcc/andromeda/internal/utils"
	"github.com/sapcc/andromeda/models"
)

var PROPERTY_TYPE_MAP = map[string]string{
	models.DomainModeAVAILABILITY: "weighted-round-robin",
	models.DomainModeGEOGRAPHIC:   "geographic",
	models.DomainModeWEIGHTED:     "weighted-round-robin",
	models.DomainModeROUNDROBIN:   "weighted-round-robin",
}

type AkamaiAgent struct {
	session           *session.Session
	gtm               gtm.GTM
	gtmLock           sync.Mutex
	domainType        string
	rpc               server.RPCServerClient
	workerTicker      *time.Ticker
	lastSync          time.Time
	lastMemberStatus  time.Time
	forceSync         chan []string
	executing         bool
	datacenterIdCache *lru.Cache[string, int]
}

var akamaiAgent *AkamaiAgent

// getDNSMetricsData retrieves DNS metrics
func getDNSMetricsData(ctx context.Context, session *session.Session, domain string, propertyName, projectID, timeRange *string) (*models.AkamaiTotalDNSRequests, error) {
	// Set default time range if not specified
	timeRangeValue := "last_hour"
	if timeRange != nil {
		timeRangeValue = *timeRange
	}

	// First, we need to get the list of properties
	// Create a request to get the domain list
	domainListURI := fmt.Sprintf("/gtm-api/v1/reports/domain-list/%s", domain)
	domainListReq, err := http.NewRequest(http.MethodGet, domainListURI, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating domain list request: %w", err)
	}

	// Define the domain summary struct to unmarshal response
	var domainSummary struct {
		Name       string   `json:"name"`
		Properties []string `json:"properties"`
	}

	// Execute the domain list request
	_, err = (*session).Exec(domainListReq, &domainSummary)
	if err != nil {
		return nil, fmt.Errorf("error fetching domain list: %w", err)
	}

	// Get the properties from the response
	properties := domainSummary.Properties

	// Filter properties by name if specified
	if propertyName != nil && *propertyName != "" {
		var filteredProperties []string
		for _, prop := range properties {
			if prop == *propertyName {
				filteredProperties = append(filteredProperties, prop)
				break
			}
		}

		// If exact match not found, try partial match
		if len(filteredProperties) == 0 {
			for _, prop := range properties {
				if strings.Contains(prop, *propertyName) {
					filteredProperties = append(filteredProperties, prop)
				}
			}
		}

		if len(filteredProperties) > 0 {
			properties = filteredProperties
		} else {
			log.WithField("property_name", *propertyName).Info("No property found matching the provided name")
		}
	}

	// Filter properties by project ID if specified
	// Note: This would require a domain lookup to match project ID to properties
	// For initial implementation, we're just logging this info
	if projectID != nil && *projectID != "" {
		log.WithField("project_id", *projectID).Info("Project ID filtering requested but not yet implemented")
	}

	// If no properties found, return empty result
	if len(properties) == 0 {
		return nil, fmt.Errorf("no Akamai GTM properties found matching the criteria")
	}

	// Use first property for initial implementation
	property := properties[0]

	// Store the requested property name for later use
	requestedProperty := property
	if propertyName != nil && *propertyName != "" {
		requestedProperty = *propertyName
	}

	// Now get the traffic report for the property
	// First, get the properties window to determine the time range
	windowURI := "/gtm-api/v1/reports/traffic/properties-window"
	windowReq, err := http.NewRequest(http.MethodGet, windowURI, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating properties window request: %w", err)
	}

	// Define the properties window struct
	var propertiesWindow struct {
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
	}

	// Execute the properties window request
	_, err = (*session).Exec(windowReq, &propertiesWindow)
	if err != nil {
		return nil, fmt.Errorf("error fetching properties window: %w", err)
	}

	// Calculate start and end times based on the time range
	end := propertiesWindow.End

	// Use the last 30 minutes by default
	start := end.Add(-30 * time.Minute)

	// Create request for the traffic report
	trafficURI := fmt.Sprintf("/gtm-api/v1/reports/traffic/domains/%s/properties/%s", domain, property)
	params := fmt.Sprintf("?start=%s&end=%s", start.UTC().Format(time.RFC3339), end.UTC().Format(time.RFC3339))
	trafficURI += params

	trafficReq, err := http.NewRequest(http.MethodGet, trafficURI, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating traffic report request: %w", err)
	}

	// Define traffic report struct
	var trafficReport struct {
		DataRows []struct {
			Timestamp   time.Time `json:"timestamp"`
			Datacenters []struct {
				Nickname          string `json:"nickname"`
				TrafficTargetName string `json:"trafficTargetName"`
				Requests          int    `json:"requests"`
				Status            string `json:"status"`
			} `json:"datacenters"`
		} `json:"dataRows"`
	}

	// Execute the traffic report request
	_, err = (*session).Exec(trafficReq, &trafficReport)
	if err != nil {
		return nil, fmt.Errorf("error fetching traffic report: %w", err)
	}

	// Build response
	totalRequests := int64(0)
	datacenters := []*models.AkamaiTotalDNSRequestsDatacentersItems0{}

	// Keep track of datacenters we've already processed
	processedDatacenters := make(map[string]int)

	// Process traffic data
	for _, dataRow := range trafficReport.DataRows {
		for _, dc := range dataRow.Datacenters {
			totalRequests += int64(dc.Requests)

			// Check if we've already processed this datacenter
			if idx, exists := processedDatacenters[dc.Nickname]; exists {
				// Update existing datacenter
				datacenters[idx].Requests += int64(dc.Requests)
				continue
			}

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

			// Store the index of this datacenter in our map
			processedDatacenters[dc.Nickname] = len(datacenters)
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

	// Log the metrics data for debugging
	debugLog := log.WithFields(log.Fields{
		"actual_property":    property,
		"requested_property": requestedProperty,
	})
	if len(datacenters) > 0 {
		jsonData, _ := json.MarshalIndent(datacenters, "", "  ")
		debugLog = debugLog.WithField("datacenters", string(jsonData))
	}
	debugLog.Debug("DNS metrics data")

	return &models.AkamaiTotalDNSRequests{
		PropertyName:  requestedProperty, // Use the property name requested by the client
		TimeRange:     timeRangeValue,
		TotalRequests: totalRequests,
		Datacenters:   datacenters,
	}, nil
}

// GetTotalDNSRequests handles RPC requests for DNS metrics data
func GetTotalDNSRequests(ctx context.Context, req stormrpc.Request) stormrpc.Response {
	type dnsMetricsRequest struct {
		PropertyName *string `json:"property_name,omitempty"`
		ProjectID    *string `json:"project_id,omitempty"`
		TimeRange    *string `json:"time_range,omitempty"`
	}

	var params dnsMetricsRequest
	if err := req.Decode(&params); err != nil {
		return stormrpc.NewErrorResponse(req.Reply, err)
	}

	// Get metrics data using our helper function
	dnsMetrics, err := getDNSMetricsData(ctx, akamaiAgent.session, config.Global.AkamaiConfig.Domain,
		params.PropertyName, params.ProjectID, params.TimeRange)
	if err != nil {
		return stormrpc.NewErrorResponse(req.Reply, err)
	}

	resp, err := stormrpc.NewResponse(req.Reply, dnsMetrics)
	if err != nil {
		return stormrpc.NewErrorResponse(req.Reply, err)
	}

	return resp
}

func Sync(ctx context.Context, req stormrpc.Request) stormrpc.Response {
	var domainIDs []string
	if err := req.Decode(&domainIDs); err != nil {
		return stormrpc.NewErrorResponse(req.Reply, err)
	}
	log.WithField("domainIDs", domainIDs).Info("[Sync] Syncing domains")

	akamaiAgent.forceSync <- domainIDs
	resp, err := stormrpc.NewResponse(req.Reply, nil)
	if err != nil {
		return stormrpc.NewErrorResponse(req.Reply, err)
	}

	return resp
}

func GetCidrs(ctx context.Context, req stormrpc.Request) stormrpc.Response {
	// Local type declaration for AkamaiCIDRBlock
	type AkamaiCIDRBlock map[string]any

	cidrBlocksReq, _ := http.NewRequest(http.MethodGet, "/firewall-rules-manager/v1/cidr-blocks", nil)
	var cidrBlocks []AkamaiCIDRBlock
	if _, err := (*akamaiAgent.session).Exec(cidrBlocksReq, &cidrBlocks); err != nil {
		panic(err)
	}

	resp, err := stormrpc.NewResponse(req.Reply, cidrBlocks)
	if err != nil {
		return stormrpc.NewErrorResponse(req.Reply, err)
	}

	return resp
}

func ExecuteAkamaiAgent() error {
	nc, err := nats.Connect(config.Global.Default.TransportURL)
	if err != nil {
		return err
	}
	client, err := stormrpc.NewClient("", stormrpc.WithNatsConn(nc))
	if err != nil {
		return err
	}

	// Create F5 worker instance with Server RPC interface
	s, domainType := NewAkamaiSession(&config.Global.AkamaiConfig)

	// Figure out minimal ticker interval
	interval := time.Duration(config.Global.AkamaiConfig.SyncInterval) + 1
	if time.Duration(config.Global.AkamaiConfig.MemberStatusInterval) < interval {
		interval = time.Duration(config.Global.AkamaiConfig.MemberStatusInterval) + 1
	}

	cache, _ := lru.New[string, int](64)

	akamaiAgent = &AkamaiAgent{
		s,
		gtm.Client(*s),
		sync.Mutex{},
		domainType,
		server.NewRPCServerClient(client),
		time.NewTicker(interval * time.Second),
		time.Unix(0, 0),
		time.Unix(0, 0),
		make(chan []string),
		false,
		cache,
	}

	if err := akamaiAgent.EnsureDomain(domainType); err != nil {
		return err
	}

	srv := rpc.NewServer("andromeda-akamai-agent", stormrpc.WithNatsConn(nc))
	srv.Handle("andromeda.sync", Sync)
	srv.Handle("andromeda.get_cidrs.akamai", GetCidrs)
	srv.Handle("andromeda.get_dns_metrics.akamai", GetTotalDNSRequests)

	go func() {
		_ = srv.Run()
	}()
	go akamaiAgent.WorkerThread()
	if config.Global.Default.Prometheus {
		go utils.PrometheusListen()
	}
	log.WithField("subjects", srv.Subjects()).Info("Subscribed")

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	// full sync immediately

	akamaiAgent.forceSync <- nil
	<-done
	log.Info("Shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return srv.Shutdown(ctx)
}

func (s *AkamaiAgent) WorkerThread() {
	syncInterval := time.Duration(config.Global.AkamaiConfig.SyncInterval) * time.Second
	memberStatusInterval := time.Duration(config.Global.AkamaiConfig.MemberStatusInterval) * time.Second

	for {
		select {
		case domains := <-s.forceSync:
			log.Debug("Running force sync")
			if err := s.FetchAndSyncDatacenters(nil, true); err != nil {
				log.Error(err.Error())
			}

			if err := s.FetchAndSyncGeomaps(nil, true); err != nil {
				log.Error(err.Error())
			}

			if err := s.FetchAndSyncDomains(domains, true); err != nil {
				log.Error(err.Error())
			}
		case <-s.workerTicker.C: // Activate periodically
			if time.Since(s.lastSync) > syncInterval {
				log.Debug("Running periodic sync")
				if err := s.FetchAndSyncDatacenters(nil, false); err != nil {
					log.Error(err.Error())
				}

				if err := s.FetchAndSyncGeomaps(nil, false); err != nil {
					log.Error(err.Error())
				}

				if err := s.FetchAndSyncDomains(nil, false); err != nil {
					log.Error(err.Error())
				}

				s.lastSync = time.Now()
			}
			if time.Since(s.lastMemberStatus) > memberStatusInterval {
				if err := s.memberStatusSync(); err != nil {
					log.Error(err.Error())
				}
				s.lastMemberStatus = time.Now()
			}
		}
	}
}

func (s *AkamaiAgent) memberStatusSync() error {
	log.Debugf("Running member status sync")
	response, err := s.rpc.GetDomains(context.Background(), &server.SearchRequest{
		Provider:       "akamai",
		PageNumber:     0,
		ResultPerPage:  1,
		FullyPopulated: true,
		Pending:        false,
	})
	if err != nil {
		return err
	}

	for _, domain := range response.GetResponse() {
		if err := s.syncMemberStatus(domain); err != nil {
			return err
		}
	}
	return nil
}
