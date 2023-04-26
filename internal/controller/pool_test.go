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

	"github.com/go-openapi/strfmt"

	"github.com/go-openapi/runtime"
	"github.com/stretchr/testify/assert"

	"github.com/sapcc/andromeda/restapi/operations/pools"
)

func (t *SuiteTest) createPool() strfmt.UUID {
	pool := pools.PostPoolsBody{}
	_ = pool.UnmarshalBinary([]byte(`{ "pool": { "name": "test" } }`))

	// Write new pool
	res := t.c.Pools.PostPools(pools.PostPoolsParams{Pool: pool})
	rr := httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), rr.Code, http.StatusCreated, rr.Body)

	poolResponse := pools.GetPoolsPoolIDOKBody{}
	_ = poolResponse.UnmarshalBinary(rr.Body.Bytes())
	assert.Equal(t.T(), "test", *poolResponse.Pool.Name, rr.Body)
	return poolResponse.Pool.ID
}

func (t *SuiteTest) TestPools() {
	rr := httptest.NewRecorder()

	// Handy alias for controller instance
	pc := t.c.Pools

	res := pc.GetPoolsPoolID(pools.GetPoolsPoolIDParams{PoolID: "test123"})
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), rr.Code, http.StatusNotFound, rr.Body)

	// Write new pool
	poolID := t.createPool()

	// Get all pools
	res = pc.GetPools(pools.GetPoolsParams{})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), rr.Code, http.StatusOK, rr.Body)

	poolsResponse := pools.GetPoolsOKBody{}
	_ = poolsResponse.UnmarshalBinary(rr.Body.Bytes())
	assert.Equal(t.T(), 1, len(poolsResponse.Pools), rr.Body)
	assert.Equal(t.T(), poolID, poolsResponse.Pools[0].ID, rr.Body)
	assert.Equal(t.T(), "test", *poolsResponse.Pools[0].Name, rr.Body)

	// Get specific pool
	res = pc.GetPoolsPoolID(pools.GetPoolsPoolIDParams{PoolID: poolID})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), rr.Code, http.StatusOK, rr.Body)

	// cleanup pool
	if _, err := t.db.Exec("DELETE FROM pool"); err != nil {
		t.FailNow(err.Error())
	}
}
