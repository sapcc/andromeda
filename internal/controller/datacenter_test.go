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

	"github.com/sapcc/andromeda/restapi/operations/datacenters"
)

func (t *SuiteTest) TestDatacenters() {
	dc := t.c.Datacenters
	rr := httptest.NewRecorder()

	res := dc.GetDatacentersDatacenterID(datacenters.GetDatacentersDatacenterIDParams{DatacenterID: "test123"})
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), rr.Code, http.StatusNotFound, rr.Body)

	datacenter := datacenters.PostDatacentersBody{}
	_ = datacenter.UnmarshalBinary([]byte(`{ "datacenter": { "name": "test" } }`))

	// Write new datacenter
	res = dc.PostDatacenters(datacenters.PostDatacentersParams{Datacenter: datacenter})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), rr.Code, http.StatusCreated, rr.Body)

	// Get all datacenters
	res = dc.GetDatacenters(datacenters.GetDatacentersParams{})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), rr.Code, http.StatusOK, rr.Body)

	datacentersResponse := datacenters.GetDatacentersOKBody{}
	_ = datacentersResponse.UnmarshalBinary(rr.Body.Bytes())
	assert.Equal(t.T(), len(datacentersResponse.Datacenters), 1, rr.Body)
	assert.Equal(t.T(), datacentersResponse.Datacenters[0].ID, datacenter.Datacenter.ID, rr.Body)

	// Get specific datacenter
	res = dc.GetDatacentersDatacenterID(datacenters.GetDatacentersDatacenterIDParams{DatacenterID: datacenter.Datacenter.ID})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), rr.Code, http.StatusOK, rr.Body)

	// Delete specific datacenter
	res = dc.DeleteDatacentersDatacenterID(datacenters.DeleteDatacentersDatacenterIDParams{DatacenterID: datacenter.Datacenter.ID})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), rr.Code, http.StatusNoContent, rr.Body)
}
