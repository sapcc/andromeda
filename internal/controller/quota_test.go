// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"net/http"
	"net/http/httptest"

	"github.com/go-openapi/runtime"
	"github.com/stretchr/testify/assert"

	"github.com/sapcc/andromeda/restapi/operations/administrative"
)

func (t *SuiteTest) TestQuotas() {
	dc := t.c.Quotas
	rr := httptest.NewRecorder()
	projectID := "test123"

	res := dc.GetQuotasProjectID(administrative.GetQuotasProjectIDParams{
		ProjectID: projectID})
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), rr.Code, http.StatusNotFound)

	quota := administrative.PutQuotasProjectIDBody{}
	_ = quota.UnmarshalBinary([]byte(`{ "quota": { "domain_akamai": 1234, "domain_f5": 2345 } }`))

	// Write new quota
	res = dc.PutQuotasProjectID(administrative.PutQuotasProjectIDParams{
		Quota:     quota,
		ProjectID: projectID})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), http.StatusAccepted, rr.Code)

	// Get all quotas
	res = dc.GetQuotas(administrative.GetQuotasParams{})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), http.StatusOK, rr.Code)

	quotasResponse := administrative.GetQuotasOKBody{}
	_ = quotasResponse.UnmarshalBinary(rr.Body.Bytes())
	assert.Equal(t.T(), len(quotasResponse.Quotas), 1, rr.Body)
	assert.Equal(t.T(), *quotasResponse.Quotas[0].DomainAkamai, int64(1234), rr.Body)
	assert.Equal(t.T(), *quotasResponse.Quotas[0].DomainF5, int64(2345), rr.Body)

	// Get specific quota
	res = dc.GetQuotasProjectID(administrative.GetQuotasProjectIDParams{
		ProjectID: projectID})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), rr.Code, http.StatusOK, rr.Body)

	// Delete specific quota
	res = dc.DeleteQuotasProjectID(administrative.DeleteQuotasProjectIDParams{
		ProjectID: projectID})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), rr.Code, http.StatusNoContent, rr.Body)
}

func (t *SuiteTest) TestQuotasUpdateSelective() {
	dc := t.c.Quotas
	projectID := "test123"

	quota := administrative.PutQuotasProjectIDBody{}
	_ = quota.UnmarshalBinary([]byte(`{ "quota": { "domain_akamai": 1234, "domain_f5": 2345 } }`))

	// Write new quota
	res := dc.PutQuotasProjectID(administrative.PutQuotasProjectIDParams{
		Quota:     quota,
		ProjectID: projectID})
	rr := httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), http.StatusAccepted, rr.Code)
	assert.JSONEq(t.T(), `{"quota":{"datacenter":0, "domain_akamai":1234, "domain_f5":2345, "member":0, "monitor":0, "pool":0}}`,
		rr.Body.String())

	// Update selective
	quota = administrative.PutQuotasProjectIDBody{}
	_ = quota.UnmarshalBinary([]byte(`{ "quota": { "datacenter": 1 } }`))
	res = dc.PutQuotasProjectID(administrative.PutQuotasProjectIDParams{
		Quota:     quota,
		ProjectID: projectID})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), http.StatusAccepted, rr.Code)
	assert.JSONEq(t.T(), `{"quota":{"datacenter":1, "domain_akamai":1234, "domain_f5":2345, "member":0, "monitor":0, "pool":0}}`,
		rr.Body.String())
}
