// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"net/http"
	"net/http/httptest"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"

	"github.com/sapcc/andromeda/models"
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

func (t *SuiteTest) TestUpdateDatacenterSelective() {
	dc := t.c.Datacenters

	testDatacenter := datacenters.PostDatacentersBody{Datacenter: &models.Datacenter{Name: swag.String("test")}}

	// Write new datacenter
	res := dc.PostDatacenters(datacenters.PostDatacentersParams{Datacenter: testDatacenter})
	rr := httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), rr.Code, http.StatusCreated, rr.Body)
	datacenter := datacenters.PostDatacentersBody{}
	_ = datacenter.UnmarshalBinary(rr.Body.Bytes())
	assert.Equal(t.T(), datacenter.Datacenter.Name, testDatacenter.Datacenter.Name)

	// Update datacenter
	testDatacenter.Datacenter.City = swag.String("berlin")
	testDatacenter.Datacenter.Name = nil
	res = dc.PutDatacentersDatacenterID(datacenters.PutDatacentersDatacenterIDParams{DatacenterID: datacenter.Datacenter.ID,
		Datacenter: datacenters.PutDatacentersDatacenterIDBody(testDatacenter)})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), http.StatusAccepted, rr.Code, rr.Body)
	_ = datacenter.UnmarshalBinary(rr.Body.Bytes())
	assert.Equal(t.T(), "test", *datacenter.Datacenter.Name)
	assert.Equal(t.T(), "berlin", *datacenter.Datacenter.City)
}
