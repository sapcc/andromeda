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

	"github.com/sapcc/andromeda/restapi/operations/administrative"
)

func (t *SuiteTest) TestQuotas() {
	dc := t.c.Quotas
	rr := httptest.NewRecorder()
	projectID := "test123"

	res := dc.GetQuotasProjectID(administrative.GetQuotasProjectIDParams{
		ProjectID: projectID})
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), rr.Code, http.StatusNotFound, rr.Body)

	quota := administrative.PutQuotasProjectIDBody{}
	_ = quota.UnmarshalBinary([]byte(`{ "quota": { "domain": 1234 } }`))

	// Write new quota
	res = dc.PutQuotasProjectID(administrative.PutQuotasProjectIDParams{
		Quota:     quota,
		ProjectID: projectID})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), http.StatusAccepted, rr.Code, rr.Body)

	// Get all quotas
	res = dc.GetQuotas(administrative.GetQuotasParams{})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), http.StatusOK, rr.Code, rr.Body)

	quotasResponse := administrative.GetQuotasOKBody{}
	_ = quotasResponse.UnmarshalBinary(rr.Body.Bytes())
	assert.Equal(t.T(), len(quotasResponse.Quotas), 1, rr.Body)
	assert.Equal(t.T(), *quotasResponse.Quotas[0].Domain, int64(1234), rr.Body)

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
