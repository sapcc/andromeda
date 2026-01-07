// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package metrics

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/session"
	"github.com/apex/log"
	"github.com/hashicorp/golang-lru/v2/expirable"
)

type PropertiesWindow struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type CachedAkamaiSession struct {
	session.Session
	managementDomain string
	rateLimiter      *AkamaiRateLimiter
	gtmLock          sync.Mutex

	// shared by all cached report fetchers
	lastProperties []string

	// shared by all cached report fetchers
	lastPropertiesFetch time.Time

	availabilityReportDataRowCache        *expirable.LRU[string, []AvailabilityReportDataRows]
	availabilityReportLastPropertyRefresh map[string]time.Time
	availabilityReportLastWindow          PropertiesWindow
	availabilityReportLastWindowReport    time.Time

	trafficReportDataRowCache        *expirable.LRU[string, []TrafficReportDataRows]
	trafficReportLastPropertyRefresh map[string]time.Time
	trafficReportLastWindow          PropertiesWindow
	trafficReportLastWindowReport    time.Time
}

func NewCachedAkamaiSession(s session.Session, domain string, rl *AkamaiRateLimiter) *CachedAkamaiSession {
	return &CachedAkamaiSession{
		Session:                               s,
		managementDomain:                      domain,
		availabilityReportDataRowCache:        expirable.NewLRU[string, []AvailabilityReportDataRows](1000, nil, 5*time.Minute),
		availabilityReportLastPropertyRefresh: make(map[string]time.Time),
		trafficReportDataRowCache:             expirable.NewLRU[string, []TrafficReportDataRows](1000, nil, 5*time.Minute),
		trafficReportLastPropertyRefresh:      make(map[string]time.Time),
		rateLimiter:                           rl,
	}
}

func (c *CachedAkamaiSession) get(uri string, out any) error {
	var err error
	var req *http.Request

	log.Debugf("Retrieving %s", uri)
	req, err = http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return err
	}

	c.rateLimiter.UseToken()

	resp, err := c.Exec(req, out)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		var dump []byte
		dump, err = httputil.DumpResponse(resp, true)
		if err != nil {
			return err
		}
		return fmt.Errorf("failed to get properties window: %s", string(dump))
	}

	return nil
}

// refreshPropertiesAvailabilityReportWindow returns the start and end time of the properties window
// to be called directly exclusively by CachedAkamaiSession
func (c *CachedAkamaiSession) refreshPropertiesAvailabilityReportWindow() error {
	// WARNING: the respective "ip-availability" endpoint does not exist so we defer to the "traffic" endpoint.
	// WARNING: undocumented endpoint.
	uri := "/gtm-api/v1/reports/traffic/properties-window"

	// only refresh if older than 5 min
	if c.availabilityReportLastWindowReport.After(time.Now().Add(-5 * time.Minute)) {
		return nil
	}

	if err := c.get(uri, &c.availabilityReportLastWindow); err != nil {
		return err
	}
	c.availabilityReportLastWindowReport = time.Now()
	log.Debugf("Set AVAILABILITY report last window report to %v", c.availabilityReportLastWindowReport)
	return nil
}

// refreshPropertiesTrafficReportWindow returns the start and end time of the properties window
// to be called directly exclusively by CachedAkamaiSession
func (c *CachedAkamaiSession) refreshPropertiesTrafficReportWindow() error {
	// WARNING: undocumented endpoint.
	uri := "/gtm-api/v1/reports/traffic/properties-window"

	// only refresh if older than 5 min
	if c.trafficReportLastWindowReport.After(time.Now().Add(-5 * time.Minute)) {
		return nil
	}

	if err := c.get(uri, &c.trafficReportLastWindow); err != nil {
		return err
	}
	c.trafficReportLastWindowReport = time.Now()
	log.Debugf("Set TRAFFIC report last window report to %v", c.trafficReportLastWindowReport)
	return nil
}

type AvailabilityReportDataRows struct {
	Timestamp   time.Time `json:"timestamp"`
	Datacenters []struct {
		IPs []struct {
			Alive     bool    `json:"alive"`
			HandedOut bool    `json:"handedOut"`
			IP        string  `json:"ip"`
			Score     float64 `json:"score"`
		} `json:"ips"`
		Nickname          string `json:"nickname"`
		TrafficTargetName string `json:"trafficTargetName"`
	} `json:"datacenters"`
}

type AvailabilityReport struct {
	Metadata struct {
		Domain   string    `json:"domain"`
		Property string    `json:"property"`
		Start    time.Time `json:"start"`
		End      time.Time `json:"end"`
		Interval string    `json:"interval"`
		Uri      string    `json:"uri"`
	} `json:"metadata"`
	AvailabilityReportDataRows []AvailabilityReportDataRows `json:"dataRows"`
	Links                      []struct {
		Rel  string `json:"rel"`
		Href string `json:"href"`
	} `json:"links"`
}

type TrafficReportDataRows struct {
	Timestamp   time.Time `json:"timestamp"`
	Datacenters []struct {
		Nickname          string `json:"nickname"`
		TrafficTargetName string `json:"trafficTargetName"`
		Requests          int    `json:"requests"`
		Status            string `json:"status"`
	} `json:"datacenters"`
}

type TrafficReport struct {
	Metadata struct {
		Domain   string    `json:"domain"`
		Property string    `json:"property"`
		Start    time.Time `json:"start"`
		End      time.Time `json:"end"`
		Interval string    `json:"interval"`
		Uri      string    `json:"uri"`
	} `json:"metadata"`
	TrafficReportDataRows []TrafficReportDataRows `json:"dataRows"`
	Links                 []struct {
		Rel  string `json:"rel"`
		Href string `json:"href"`
	} `json:"links"`
}

// getAvailabilityReport retrieves the Akamai "IP availability" report for a
// property (i.e. Andromeda domain).
//
// Unlike the "traffic" report, which Akamai updates on average every 5 minutes
// for most properties, the "IP availability" report has its own (much less
// reliable) cadence. A given property report may be as much as 90 minutes
// behind compared to another property report under the same Akamai domain.
func (c *CachedAkamaiSession) getAvailabilityReport(property string) ([]AvailabilityReportDataRows, error) {
	if cached, ok := c.availabilityReportDataRowCache.Get(property); ok {
		log.Debugf("found AVAILABILITY report for property %s in cache", property)
		return cached, nil
	}

	c.gtmLock.Lock()
	defer c.gtmLock.Unlock()

	// reference: <https://techdocs.akamai.com/gtm-reporting/reference/get-ip-availability-property>
	path := fmt.Sprintf("/gtm-api/v1/reports/ip-availability/domains/%s/properties/%s", c.managementDomain, property)

	if err := c.refreshPropertiesAvailabilityReportWindow(); err != nil {
		return nil, err
	}

	start, ok := c.availabilityReportLastPropertyRefresh[property]
	end := c.availabilityReportLastWindow.End
	if end.After(time.Now().Add(-5 * time.Minute)) {
		// quirky api
		end = end.Add(-5 * time.Minute)
	}
	if ok {
		start = start.Add(time.Minute)

		if end.Before(start.Add(5 * time.Minute)) {
			start = start.Add(-5 * time.Minute)
		}

		if start.Before(c.availabilityReportLastWindow.Start) {
			// truncate
			start = c.availabilityReportLastWindow.Start
		}
	} else {
		// Get the last 120 min (a property report could easily be delayed by 90 minutes)
		start = c.availabilityReportLastWindow.End.Add(-120 * time.Minute)
	}

	params := url.Values{}
	params.Add("start", start.UTC().Format(time.RFC3339))
	params.Add("end", end.UTC().Format(time.RFC3339))
	log.Debugf("[CachedAkamaiSession] Time interval for AVAILABILITY reports set to [ START = %s | END = %s ]",
		start.UTC().Format(time.RFC3339),
		end.UTC().Format(time.RFC3339),
	)
	uri := fmt.Sprintf("%s?%s", path, params.Encode())
	log.Infof("[CachedAkamaiSession] Retrieving %s", uri)

	var availabilityReport AvailabilityReport
	err := c.get(uri, &availabilityReport)
	if err != nil {
		return nil, err
	}

	c.availabilityReportLastPropertyRefresh[property] = end
	evicted := c.availabilityReportDataRowCache.Add(property, availabilityReport.AvailabilityReportDataRows)
	log.Debugf("[CachedAkamaiSession] evicted: %t", evicted)
	return availabilityReport.AvailabilityReportDataRows, nil
}

func (c *CachedAkamaiSession) getTrafficReport(property string) ([]TrafficReportDataRows, error) {
	if cached, ok := c.trafficReportDataRowCache.Get(property); ok {
		return cached, nil
	}

	c.gtmLock.Lock()
	defer c.gtmLock.Unlock()

	// reference: <https://techdocs.akamai.com/gtm-reporting/reference/get-traffic-datacenter>
	path := fmt.Sprintf("/gtm-api/v1/reports/traffic/domains/%s/properties/%s", c.managementDomain, property)

	if err := c.refreshPropertiesTrafficReportWindow(); err != nil {
		return nil, err
	}

	start, ok := c.trafficReportLastPropertyRefresh[property]
	end := c.trafficReportLastWindow.End
	if end.After(time.Now().Add(-5 * time.Minute)) {
		// quirky api
		end = end.Add(-5 * time.Minute)
	}
	if ok {
		start = start.Add(time.Minute)

		if end.Before(start.Add(5 * time.Minute)) {
			start = start.Add(-5 * time.Minute)
		}

		if start.Before(c.trafficReportLastWindow.Start) {
			// truncate
			start = c.trafficReportLastWindow.Start
		}
	} else {
		// Get the last 30 min
		start = c.trafficReportLastWindow.End.Add(-30 * time.Minute)
	}

	params := url.Values{}
	params.Add("start", start.UTC().Format(time.RFC3339))
	params.Add("end", end.UTC().Format(time.RFC3339))
	log.Debugf("[CachedAkamaiSession] Time interval for traffic reports set to [ START = %s | END = %s ]",
		start.UTC().Format(time.RFC3339),
		end.UTC().Format(time.RFC3339),
	)
	uri := fmt.Sprintf("%s?%s", path, params.Encode())
	log.Infof("[CachedAkamaiSession] Retrieving %s", uri)

	var trafficReport TrafficReport
	err := c.get(uri, &trafficReport)
	if err != nil {
		return nil, err
	}

	c.trafficReportLastPropertyRefresh[property] = end
	evicted := c.trafficReportDataRowCache.Add(property, trafficReport.TrafficReportDataRows)
	log.Debugf("[CachedAkamaiSession] evicted: %t", evicted)
	return trafficReport.TrafficReportDataRows, nil
}

type DomainSummary struct {
	Name        string   `json:"name"`
	Properties  []string `json:"properties"`
	Datacenters []struct {
		DatacenterId       int    `json:"datacenterId"`
		DatacenterNickname string `json:"datacenterNickname"`
	} `json:"datacenters"`
	Resources []interface{} `json:"resources"`
	Links     []struct {
		Rel  string `json:"rel"`
		Href string `json:"href"`
	} `json:"links"`
}

// getProperties returns the properties
func (c *CachedAkamaiSession) getProperties() ([]string, error) {
	c.gtmLock.Lock()
	defer c.gtmLock.Unlock()

	uri := fmt.Sprintf("/gtm-api/v1/reports/domain-list/%s", c.managementDomain)

	// only refresh if older than 10 min, or never fetched
	if len(c.lastProperties) == 0 || c.lastPropertiesFetch.Before(time.Now().Add(-10*time.Minute)) {
		var domainSummary DomainSummary
		err := c.get(uri, &domainSummary)
		if err != nil {
			return nil, err
		}

		c.lastProperties = domainSummary.Properties
		c.lastPropertiesFetch = time.Now()
	}

	return c.lastProperties, nil
}
