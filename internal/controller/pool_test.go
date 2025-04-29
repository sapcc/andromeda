// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"net/http"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/go-openapi/runtime"
	"github.com/stretchr/testify/assert"

	"github.com/sapcc/andromeda/models"
	"github.com/sapcc/andromeda/restapi/operations/pools"
)

func (t *SuiteTest) createPool(domains []strfmt.UUID) strfmt.UUID {
	pool := pools.PostPoolsBody{
		Pool: &models.Pool{
			Name:    swag.String("test"),
			Domains: domains,
		},
	}

	// Write new pool
	res := t.c.Pools.PostPools(pools.PostPoolsParams{Pool: pool})
	rr := httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), http.StatusCreated, rr.Code, rr.Body)

	poolResponse := pools.GetPoolsPoolIDOKBody{}
	_ = poolResponse.UnmarshalBinary(rr.Body.Bytes())
	assert.Equal(t.T(), "test", *poolResponse.Pool.Name, rr.Body)
	return poolResponse.Pool.ID
}

// cleanupPools deletes all pools from database
func (t *SuiteTest) cleanupPools() {
	_, err := t.db.Exec("DELETE FROM pool")
	if err != nil {
		t.FailNow(err.Error())
	}
}

func (t *SuiteTest) TestPools() {
	// Handy alias for controller instance
	pc := t.c.Pools

	res := pc.GetPoolsPoolID(pools.GetPoolsPoolIDParams{PoolID: "test123"})
	rr := httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), http.StatusNotFound, rr.Code, rr.Body)

	// Write new pool
	poolID := t.createPool(nil)

	// Get all pools
	res = pc.GetPools(pools.GetPoolsParams{})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), http.StatusOK, rr.Code, rr.Body)

	poolsResponse := pools.GetPoolsOKBody{}
	_ = poolsResponse.UnmarshalBinary(rr.Body.Bytes())
	assert.Equal(t.T(), 1, len(poolsResponse.Pools), rr.Body)
	assert.Equal(t.T(), poolID, poolsResponse.Pools[0].ID, rr.Body)
	assert.Equal(t.T(), "test", *poolsResponse.Pools[0].Name, rr.Body)

	// Get specific pool
	res = pc.GetPoolsPoolID(pools.GetPoolsPoolIDParams{PoolID: poolID})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), http.StatusOK, rr.Code, rr.Body)

	// cleanup pool
	if _, err := t.db.Exec("DELETE FROM pool"); err != nil {
		t.FailNow(err.Error())
	}
}

func (t *SuiteTest) TestPoolImmutable() {
	// Handy alias for controller instance
	pc := t.c.Pools

	domainID := t.createDomain()

	// Write new pool
	poolID := t.createPool([]strfmt.UUID{domainID})

	defer t.cleanupPools()
	defer t.cleanupDomains()

	// Update pool
	pool := pools.PutPoolsPoolIDBody{
		Pool: &models.Pool{
			Name: swag.String("test2"),
		},
	}
	res := pc.PutPoolsPoolID(pools.PutPoolsPoolIDParams{PoolID: poolID, Pool: pool})
	rr := httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), http.StatusConflict, rr.Code, rr.Body)
}
