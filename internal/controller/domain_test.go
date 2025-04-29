// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"net/http"
	"net/http/httptest"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"

	"github.com/sapcc/andromeda/models"
	"github.com/sapcc/andromeda/restapi/operations/domains"
)

func (t *SuiteTest) createDomain() strfmt.UUID {
	fqdn := strfmt.Hostname("test.com")
	domain := domains.PostDomainsBody{
		Domain: &models.Domain{
			Fqdn:     &fqdn,
			Name:     swag.String("test"),
			Provider: swag.String("akamai"),
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
