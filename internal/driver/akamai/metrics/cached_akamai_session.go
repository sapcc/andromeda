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

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	"github.com/apex/log"
	"github.com/hashicorp/golang-lru/v2/expirable"
)

type PropertiesWindow struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type CachedAkamaiSession struct {
	session.Session
	managementDomain    string
	lastProperties      []string
	lastPropertiesFetch time.Time
	lastWindowRefresh   time.Time
	lastWindow          PropertiesWindow
	lastPropertyRefresh map[string]time.Time
	dataRowCache        *expirable.LRU[string, []DataRows]
	rateLimiter         *AkamaiRateLimiter
	gtmLock             sync.Mutex
}

func NewCachedAkamaiSession(s session.Session, domain string, rl *AkamaiRateLimiter) *CachedAkamaiSession {
	return &CachedAkamaiSession{
		Session:             s,
		managementDomain:    domain,
		lastPropertyRefresh: make(map[string]time.Time),
		dataRowCache:        expirable.NewLRU[string, []DataRows](1000, nil, 5*time.Minute),
		rateLimiter:         rl,
	}
}

func (c *CachedAkamaiSession) get(uri string, out any) error {
	var err error
	var req *http.Request

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

// getPropertiesWindow returns the start and end time of the properties window
// to be called directly exclusively by CachedAkamaiSession
func (c *CachedAkamaiSession) refreshPropertiesWindow() error {
	uri := "/gtm-api/v1/reports/traffic/properties-window"

	// only refresh if older than 5 min
	if c.lastWindowRefresh.After(time.Now().Add(-5 * time.Minute)) {
		return nil
	}

	if err := c.get(uri, &c.lastWindow); err != nil {
		return err
	}
	c.lastWindowRefresh = time.Now()
	return nil
}

type DataRows struct {
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
	DataRows []DataRows `json:"dataRows"`
	Links    []struct {
		Rel  string `json:"rel"`
		Href string `json:"href"`
	} `json:"links"`
}

func (c *CachedAkamaiSession) getTrafficReport(property string) ([]DataRows, error) {
	if cached, ok := c.dataRowCache.Get(property); ok {
		return cached, nil
	}

	c.gtmLock.Lock()
	defer c.gtmLock.Unlock()

	path := fmt.Sprintf("/gtm-api/v1/reports/traffic/domains/%s/properties/%s", c.managementDomain, property)

	if err := c.refreshPropertiesWindow(); err != nil {
		return nil, err
	}

	start, ok := c.lastPropertyRefresh[property]
	end := c.lastWindow.End
	if end.After(time.Now().Add(-5 * time.Minute)) {
		// quirky api
		end = end.Add(-5 * time.Minute)
	}
	if ok {
		start = start.Add(time.Minute)

		if end.Before(start.Add(5 * time.Minute)) {
			start = start.Add(-5 * time.Minute)
		}

		if start.Before(c.lastWindow.Start) {
			// truncate
			start = c.lastWindow.Start
		}
	} else {
		// Get the last 30 min
		start = c.lastWindow.End.Add(-30 * time.Minute)
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

	c.lastPropertyRefresh[property] = end
	evicted := c.dataRowCache.Add(property, trafficReport.DataRows)
	log.Debugf("[CachedAkamaiSession] evicted: %t", evicted)
	return trafficReport.DataRows, nil
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
