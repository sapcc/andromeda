// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag/conv"
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

	testDatacenter := datacenters.PostDatacentersBody{Datacenter: &models.Datacenter{Name: conv.Pointer("test")}}

	// Write new datacenter
	res := dc.PostDatacenters(datacenters.PostDatacentersParams{Datacenter: testDatacenter})
	rr := httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), rr.Code, http.StatusCreated, rr.Body)
	datacenter := datacenters.PostDatacentersBody{}
	_ = datacenter.UnmarshalBinary(rr.Body.Bytes())
	assert.Equal(t.T(), datacenter.Datacenter.Name, testDatacenter.Datacenter.Name)

	// Update datacenter
	testDatacenter.Datacenter.City = conv.Pointer("berlin")
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

func (st *SuiteTest) TestRestrictedDatacenterManagement() {
	assert := assert.New(st.T())

	st.T().Run("The API should forbid new datacenters with provider=f5", func(t *testing.T) {
		dc := st.c.Datacenters

		res := dc.GetDatacentersDatacenterID(datacenters.GetDatacentersDatacenterIDParams{DatacenterID: "f5_dc_test"})
		rr := httptest.NewRecorder()
		res.WriteResponse(rr, runtime.JSONProducer())
		assert.Equal(http.StatusNotFound, rr.Code, rr.Body)

		datacenter := datacenters.PostDatacentersBody{}
		_ = datacenter.UnmarshalBinary([]byte(`{ "datacenter": { "name": "f5_dc_test", "provider": "f5", "scope": "public" } }`))

		// Write new datacenter
		res = dc.PostDatacenters(datacenters.PostDatacentersParams{Datacenter: datacenter})
		rr = httptest.NewRecorder()
		res.WriteResponse(rr, runtime.JSONProducer())
		assert.Equal(http.StatusBadRequest, rr.Code, rr.Body)
	})

	st.T().Run("The API should forbid modification of datacenters with provider=f5", func(t *testing.T) {
		datacenter := mustCreateF5Datacenter(st)
		dcID := datacenter.ID

		dc := st.c.Datacenters
		rr := httptest.NewRecorder()

		// Attempt to change name (or any other field)
		res := dc.PutDatacentersDatacenterID(datacenters.PutDatacentersDatacenterIDParams{
			Datacenter: datacenters.PutDatacentersDatacenterIDBody{
				Datacenter: &models.Datacenter{
					Name: conv.Pointer("test_f5_dc_modified"),
				},
			},
			DatacenterID: dcID,
		})
		res.WriteResponse(rr, runtime.JSONProducer())
		assert.Equal(http.StatusBadRequest, rr.Code, rr.Body)

		// Cleanup
		if _, err := st.db.Exec("DELETE FROM datacenter"); err != nil {
			st.FailNow(err.Error())
		}
	})

	st.T().Run("The API should forbid existing non-f5 datacenters from setting provider=f5", func(t *testing.T) {
		dc := st.c.Datacenters

		datacenter := datacenters.PostDatacentersBody{}
		_ = datacenter.UnmarshalBinary([]byte(`{ "datacenter": { "name": "test_akamai_dc", "scope": "private", "provider": "akamai" } }`))

		// Write new datacenter
		res := dc.PostDatacenters(datacenters.PostDatacentersParams{Datacenter: datacenter})
		rr := httptest.NewRecorder()
		res.WriteResponse(rr, runtime.JSONProducer())
		assert.Equal(http.StatusCreated, rr.Code, rr.Body)

		// Ensure the datacenter ID can be retrieved
		var freshDC datacenters.PostDatacentersCreatedBody
		if err := json.Unmarshal(rr.Body.Bytes(), &freshDC); err != nil {
			st.FailNow(err.Error())
		}

		// Attempt to change provider to "f5"
		res = dc.PutDatacentersDatacenterID(datacenters.PutDatacentersDatacenterIDParams{
			Datacenter: datacenters.PutDatacentersDatacenterIDBody{
				Datacenter: &models.Datacenter{
					Name:     conv.Pointer("test_akamai_dc"),
					Provider: "f5",
				},
			},
			DatacenterID: freshDC.Datacenter.ID,
		})
		rr = httptest.NewRecorder()
		res.WriteResponse(rr, runtime.JSONProducer())
		assert.Equal(http.StatusBadRequest, rr.Code, rr.Body)

		// Cleanup
		if _, err := st.db.Exec("DELETE FROM datacenter"); err != nil {
			st.FailNow(err.Error())
		}
	})

	st.T().Run("The API should forbid existing f5 datacenters from being deleted", func(t *testing.T) {
		datacenter := mustCreateF5Datacenter(st)

		// Attempt to delete the fresh datacenter
		dc := st.c.Datacenters
		res := dc.DeleteDatacentersDatacenterID(datacenters.DeleteDatacentersDatacenterIDParams{
			DatacenterID: datacenter.ID,
		})
		rr := httptest.NewRecorder()
		res.WriteResponse(rr, runtime.JSONProducer())
		assert.Equal(http.StatusBadRequest, rr.Code, rr.Body)

		// Cleanup
		if _, err := st.db.Exec("DELETE FROM datacenter"); err != nil {
			st.FailNow(err.Error())
		}
	})
}

func mustCreateF5Datacenter(st *SuiteTest) *models.Datacenter {
	num, err := conv.ConvertFloat64("20")
	if err != nil {
		st.FailNow(err.Error())
	}
	datacenter := &models.Datacenter{
		Name:            conv.Pointer("test_f5_dc"),
		Provider:        "f5",
		Scope:           conv.Pointer("public"),
		AdminStateUp:    conv.Pointer(true),
		StateOrProvince: conv.Pointer("zz"),
		City:            conv.Pointer("zz"),
		Continent:       conv.Pointer("zz"),
		Country:         conv.Pointer("zz"),
		Latitude:        &num,
		Longitude:       &num,
		ProjectID:       conv.Pointer("-"),
	}

	sql := `
		INSERT INTO datacenter ( name,  admin_state_up,  continent,  country,  state_or_province,  city,  latitude,  longitude,  scope,  project_id,  provider)
		VALUES                 (:name, :admin_state_up, :continent, :country, :state_or_province, :city, :latitude, :longitude, :scope, :project_id, :provider)
		RETURNING *`

	stmt, _ := st.db.PrepareNamed(sql)
	if err := stmt.Get(datacenter, datacenter); err != nil {
		st.FailNow(err.Error())
	}

	return datacenter
}
