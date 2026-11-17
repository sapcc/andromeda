// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"net/http"
	"net/http/httptest"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag/conv"
	"github.com/stretchr/testify/assert"

	"github.com/sapcc/andromeda/models"
	"github.com/sapcc/andromeda/restapi/operations/datacenters"
	"github.com/sapcc/andromeda/restapi/operations/domains"
	"github.com/sapcc/andromeda/restapi/operations/geographic_maps"
	"github.com/sapcc/andromeda/restapi/operations/members"
)

func (t *SuiteTest) createDomain() strfmt.UUID {
	fqdn := strfmt.Hostname("test.com")
	domain := domains.PostDomainsBody{
		Domain: &models.Domain{
			Fqdn:     &fqdn,
			Name:     conv.Pointer("test"),
			Provider: conv.Pointer("akamai"),
		},
	}

	// Write new domain
	res := t.c.Domains.PostDomains(domains.PostDomainsParams{Domain: domain})
	rr := httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), http.StatusCreated, rr.Code, rr.Body)

	domainResponse := domains.PostDomainsCreatedBody{}
	_ = domainResponse.UnmarshalBinary(rr.Body.Bytes())
	assert.Equal(t.T(), "test", *domainResponse.Domain.Name, rr.Body)
	return domainResponse.Domain.ID
}

// cleanupDomains deletes all domains from the database
func (t *SuiteTest) cleanupDomains() {
	_, err := t.db.Exec("DELETE FROM domain")
	if err != nil {
		t.FailNow(err.Error())
	}
}

func (t *SuiteTest) TestDomains() {
	dc := t.c.Domains
	rr := httptest.NewRecorder()

	res := dc.GetDomainsDomainID(domains.GetDomainsDomainIDParams{DomainID: "test123"})
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), rr.Code, http.StatusNotFound, rr.Body)

	// Write new domain
	domainID := t.createDomain()
	defer t.cleanupDomains()

	// Get all domains
	res = dc.GetDomains(domains.GetDomainsParams{})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), http.StatusOK, rr.Code, rr.Body)

	domainsResponse := domains.GetDomainsOKBody{}
	_ = domainsResponse.UnmarshalBinary(rr.Body.Bytes())
	assert.Equal(t.T(), len(domainsResponse.Domains), 1, rr.Body)
	assert.Equal(t.T(), domainsResponse.Domains[0].ID, domainID, rr.Body)

	// Get specific domain
	res = dc.GetDomainsDomainID(domains.GetDomainsDomainIDParams{DomainID: domainID})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), rr.Code, http.StatusOK, rr.Body)

	// Delete specific domain
	res = dc.DeleteDomainsDomainID(domains.DeleteDomainsDomainIDParams{DomainID: domainID})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), rr.Code, http.StatusNoContent, rr.Body)
}

// TestGeographicDomainMemberCoverage verifies that updating a geographic domain's pools
// is rejected when no pool member covers the geomap's default datacenter.
func (t *SuiteTest) TestGeographicDomainMemberCoverage() {
	defer t.cleanupDomains()
	defer func() { _, _ = t.db.Exec("DELETE FROM member") }()
	defer func() { _, _ = t.db.Exec("DELETE FROM pool") }()
	defer func() { _, _ = t.db.Exec("DELETE FROM geographic_map") }()
	defer func() { _, _ = t.db.Exec("DELETE FROM datacenter") }()

	// Create two datacenters: default DC and another
	dcBody := datacenters.PostDatacentersBody{}
	_ = dcBody.UnmarshalBinary([]byte(`{"datacenter": {"name": "dc-default", "provider": "akamai"}}`))
	res := t.c.Datacenters.PostDatacenters(datacenters.PostDatacentersParams{Datacenter: dcBody})
	rr := httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), http.StatusCreated, rr.Code, rr.Body)
	dcDefault := datacenters.PostDatacentersCreatedBody{}
	_ = dcDefault.UnmarshalBinary(rr.Body.Bytes())
	defaultDCID := dcDefault.Datacenter.ID

	dcBody2 := datacenters.PostDatacentersBody{}
	_ = dcBody2.UnmarshalBinary([]byte(`{"datacenter": {"name": "dc-other", "provider": "akamai"}}`))
	res = t.c.Datacenters.PostDatacenters(datacenters.PostDatacentersParams{Datacenter: dcBody2})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), http.StatusCreated, rr.Code, rr.Body)
	dcOther := datacenters.PostDatacentersCreatedBody{}
	_ = dcOther.UnmarshalBinary(rr.Body.Bytes())
	otherDCID := dcOther.Datacenter.ID

	// Create geomap with defaultDCID as default_datacenter
	geomapBody := geographic_maps.PostGeomapsBody{
		Geomap: &models.Geomap{
			DefaultDatacenter: &defaultDCID,
			Provider:          "akamai",
		},
	}
	res = t.c.GeoMaps.PostGeomaps(geographic_maps.PostGeomapsParams{Geomap: geomapBody})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), http.StatusCreated, rr.Code, rr.Body)
	geomapResp := geographic_maps.PostGeomapsCreatedBody{}
	_ = geomapResp.UnmarshalBinary(rr.Body.Bytes())
	geomapID := geomapResp.Geomap.ID

	// Create a geographic domain
	fqdn := strfmt.Hostname("geo.test.com")
	mode := models.DomainModeGEOGRAPHIC
	domainBody := domains.PostDomainsBody{
		Domain: &models.Domain{
			Fqdn:      &fqdn,
			Provider:  conv.Pointer("akamai"),
			Mode:      &mode,
			GeomapID:  &geomapID,
		},
	}
	res = t.c.Domains.PostDomains(domains.PostDomainsParams{Domain: domainBody})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), http.StatusCreated, rr.Code, rr.Body)
	domainResp := domains.PostDomainsCreatedBody{}
	_ = domainResp.UnmarshalBinary(rr.Body.Bytes())
	domainID := domainResp.Domain.ID

	// Create a pool with a member in the OTHER datacenter (not default)
	poolID := t.createPool(nil)
	memberBody := members.PostMembersBody{}
	_ = memberBody.UnmarshalBinary([]byte(`{"member": {"address": "1.2.3.4", "port": 80}}`))
	memberBody.Member.PoolID = &poolID
	memberBody.Member.DatacenterID = &otherDCID
	res = t.c.Members.PostMembers(members.PostMembersParams{Member: memberBody})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), http.StatusCreated, rr.Code, rr.Body)

	// Attaching the pool to the geographic domain should fail: no member in default DC
	updateBody := domains.PutDomainsDomainIDBody{
		Domain: &models.Domain{
			Pools: []strfmt.UUID{poolID},
		},
	}
	res = t.c.Domains.PutDomainsDomainID(domains.PutDomainsDomainIDParams{
		DomainID: domainID,
		Domain:   updateBody,
	})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), http.StatusBadRequest, rr.Code, rr.Body)

	// Now add a member in the default datacenter and retry — should succeed
	memberBody2 := members.PostMembersBody{}
	_ = memberBody2.UnmarshalBinary([]byte(`{"member": {"address": "5.6.7.8", "port": 80}}`))
	memberBody2.Member.PoolID = &poolID
	memberBody2.Member.DatacenterID = &defaultDCID
	res = t.c.Members.PostMembers(members.PostMembersParams{Member: memberBody2})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), http.StatusCreated, rr.Code, rr.Body)

	res = t.c.Domains.PutDomainsDomainID(domains.PutDomainsDomainIDParams{
		DomainID: domainID,
		Domain:   updateBody,
	})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), http.StatusAccepted, rr.Code, rr.Body)
}
