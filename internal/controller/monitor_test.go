/*
 *   Copyright 2024 SAP SE
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
	"testing"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"

	"github.com/sapcc/andromeda/internal/utils"
	"github.com/sapcc/andromeda/models"
	"github.com/sapcc/andromeda/restapi/operations/monitors"
)

func (t *SuiteTest) TestMonitors() {
	mc := t.c.Monitors
	rr := httptest.NewRecorder()

	poolID := t.createPool()
	res := mc.GetMonitorsMonitorID(monitors.GetMonitorsMonitorIDParams{
		MonitorID: "test123",
	})
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), http.StatusNotFound, rr.Code, rr.Body)

	monitor := monitors.PostMonitorsBody{}
	_ = monitor.UnmarshalBinary([]byte(`{ "monitor": { "name": "test", "project_id": "" } }`))
	monitor.Monitor.PoolID = &poolID

	// Write new monitor
	res = mc.PostMonitors(monitors.PostMonitorsParams{Monitor: monitor})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), http.StatusCreated, rr.Code, rr.Body)

	// Get all monitors
	res = mc.GetMonitors(monitors.GetMonitorsParams{PoolID: &poolID})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), http.StatusOK, rr.Code, rr.Body)

	monitorsResponse := monitors.GetMonitorsOKBody{}
	_ = monitorsResponse.UnmarshalBinary(rr.Body.Bytes())
	assert.Equal(t.T(), len(monitorsResponse.Monitors), 1, rr.Body)
	assert.Equal(t.T(), monitorsResponse.Monitors[0].ID, monitor.Monitor.ID, rr.Body)
	assert.Equal(t.T(), *monitorsResponse.Monitors[0].Name, "test", rr.Body)

	// Get all monitors without pool_id filter
	res = mc.GetMonitors(monitors.GetMonitorsParams{})
	rr = httptest.NewRecorder()
	res.WriteResponse(rr, runtime.JSONProducer())
	assert.Equal(t.T(), http.StatusOK, rr.Code, rr.Body)

	monitorsResponse = monitors.GetMonitorsOKBody{}
	_ = monitorsResponse.UnmarshalBinary(rr.Body.Bytes())
	assert.Equal(t.T(), 1, len(monitorsResponse.Monitors), rr.Body)

	// cleanup pool
	if _, err := t.db.Exec("DELETE FROM pool"); err != nil {
		t.FailNow(err.Error())
	}
}

func TestValidateMonitor(t *testing.T) {
	monitor := models.Monitor{
		Send: swag.String("GET /"),
		Type: swag.String(models.MonitorTypeHTTP),
	}
	assert.Equal(t, validateMonitor(&monitor), utils.InvalidSendString)

	monitor.Send = swag.String("http://example.com/test")
	assert.Equal(t, validateMonitor(&monitor), utils.InvalidSendString)

	monitor.Send = swag.String("/test/site?param1=1&param2=2")
	assert.Nil(t, validateMonitor(&monitor))
}
