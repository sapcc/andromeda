/*
 *   Copyright 2022 SAP SE
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package controller

import (
	"net/http"
	"net/http/httptest"

	"github.com/go-openapi/runtime"
	"github.com/stretchr/testify/assert"

	"github.com/sapcc/andromeda/restapi/operations/domains"
)

func (t *SuiteTest) TestDomains() {
	dc := t.c.Domains
	rr := httptest.NewRecorder()

	res := dc.GetDomainsDomainID(domains.GetDomainsDomainIDParams{DomainID: "test123"})
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), rr.Code, http.StatusNotFound, rr.Body)

	domain := domains.PostDomainsBody{}
	_ = domain.UnmarshalBinary([]byte(`
{
    "domain": 
        {
            "fqdn": "test.com",
            "name": "test",
            "provider": "f5"
        }
}`))

	// Write new domain
	res = dc.PostDomains(domains.PostDomainsParams{Domain: domain})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), http.StatusCreated, rr.Code, rr.Body)

	// Get all domains
	res = dc.GetDomains(domains.GetDomainsParams{})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), http.StatusOK, rr.Code, rr.Body)

	domainsResponse := domains.GetDomainsOKBody{}
	_ = domainsResponse.UnmarshalBinary(rr.Body.Bytes())
	assert.Equal(t.T(), len(domainsResponse.Domains), 1, rr.Body)
	assert.Equal(t.T(), domainsResponse.Domains[0].ID, domain.Domain.ID, rr.Body)

	// Get specific domain
	res = dc.GetDomainsDomainID(domains.GetDomainsDomainIDParams{DomainID: domain.Domain.ID})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), rr.Code, http.StatusOK, rr.Body)

	// Delete specific domain
	res = dc.DeleteDomainsDomainID(domains.DeleteDomainsDomainIDParams{DomainID: domain.Domain.ID})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), rr.Code, http.StatusNoContent, rr.Body)
}
