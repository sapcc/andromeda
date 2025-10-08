// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package f5

import (
	"errors"
	"strings"
	"testing"

	"github.com/f5devcentral/go-bigip"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/sapcc/andromeda/internal/rpcmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCollectVirtualServerMetrics(t *testing.T) {
	picksCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "picks_counter", Help: "Picks counter"},
		[]string{"domain", "datacenter_id", "project_id", "target_ip"},
	)
	prometheus.MustRegister(picksCounter)
	assert := assert.New(t)

	t.Run("Fails if it cannot fetch datacenters over RPC", func(t *testing.T) {
		session := new(mockedBigIPSession)
		store := new(mockedStore)
		store.On("GetDatacenters").Return([]*rpcmodels.Datacenter{}, errors.New("RPC failed for datacenters"))
		err := collectVirtualServerMetrics(session, store, picksCounter)
		assert.ErrorContains(err, "RPC failed for datacenters")
		store.AssertCalled(t, "GetDatacenters")
		store.AssertNotCalled(t, "GetDomains")
		session.AssertNotCalled(t, "APIRequest", mock.Anything)
	})

	t.Run("Fails if it cannot fetch domains over RPC", func(t *testing.T) {
		session := new(mockedBigIPSession)
		store := new(mockedStore)
		store.On("GetDatacenters").Return([]*rpcmodels.Datacenter{{Id: "dc1-uuid", Name: "dc1-name"}}, nil)
		store.On("GetDomains").Return([]*rpcmodels.Domain{}, errors.New("RPC failed for domains"))
		err := collectVirtualServerMetrics(session, store, picksCounter)
		assert.ErrorContains(err, "RPC failed for domains")
		store.AssertCalled(t, "GetDatacenters")
		store.AssertCalled(t, "GetDomains")
		session.AssertNotCalled(t, "APIRequest", mock.Anything)
	})

	t.Run("Succeeds by incrementing the given counter metric accordingly", func(t *testing.T) {
		expectedURLPaths := []string{
			"gtm/server/~Common~cc_andromeda_srv_10.10.0.11_dc1-name/stats",
			"gtm/server/~Common~cc_andromeda_srv_10.10.0.12_dc2-name/stats",
			"gtm/server/~Common~cc_andromeda_srv_10.10.0.13_dc2-name/stats",
			"gtm/server/~Common~cc_andromeda_srv_10.10.0.14_dc2-name/stats"}
		session := new(mockedBigIPSession)
		session.
			On("APICall", &bigip.APIRequest{Method: "get", ContentType: "application/json", URL: expectedURLPaths[0]}).
			Return([]byte(`{"entries": {"theKey": {"nestedStats": {"entries": {"vsPicks": {"value": 51}}}}}}`), nil)
		session.
			On("APICall", &bigip.APIRequest{Method: "get", ContentType: "application/json", URL: expectedURLPaths[1]}).
			Return([]byte(`{"entries": {"theKey": {"nestedStats": {"entries": {"vsPicks": {"value": 52}}}}}}`), nil)
		session.
			On("APICall", &bigip.APIRequest{Method: "get", ContentType: "application/json", URL: expectedURLPaths[2]}).
			Return([]byte(`{"entries": {"theKey": {"nestedStats": {"entries": {"vsPicks": {"value": 53}}}}}}`), nil)
		session.
			On("APICall", &bigip.APIRequest{Method: "get", ContentType: "application/json", URL: expectedURLPaths[3]}).
			Return([]byte(`{"code": 404}`), errors.New("entity not found"))
		store := new(mockedStore)
		store.On("GetDatacenters").Return([]*rpcmodels.Datacenter{
			{Id: "dc1-uuid", Name: "dc1-name"},
			{Id: "dc2-uuid", Name: "dc2-name"}}, nil)
		store.On("GetDomains").Return([]*rpcmodels.Domain{{Id: "dom1-uuid", Fqdn: "foo.com", Pools: []*rpcmodels.Pool{
			{Id: "pool2-uuid", Members: []*rpcmodels.Member{
				{Id: "member1-uuid", Address: "10.10.0.11", Port: 80, DatacenterId: "dc1-uuid", ProjectId: "p1-uuid"},
				{Id: "member2-uuid", Address: "10.10.0.12", Port: 80, DatacenterId: "dc2-uuid", ProjectId: "p1-uuid"},
				{Id: "member3-uuid", Address: "10.10.0.13", Port: 80, DatacenterId: "dc2-uuid", ProjectId: "p1-uuid"},
				{Id: "member4-uuid", Address: "10.10.0.14", Port: 80, DatacenterId: "dc2-uuid", ProjectId: "p1-uuid"}}}}}}, nil)
		err := collectVirtualServerMetrics(session, store, picksCounter)
		assert.Nil(err)
		expected := strings.NewReader(`
			# HELP picks_counter Picks counter
			# TYPE picks_counter counter
			picks_counter{datacenter_id="dc1-uuid", domain="foo.com", project_id="p1-uuid", target_ip="10.10.0.11"} 51
			picks_counter{datacenter_id="dc2-uuid", domain="foo.com", project_id="p1-uuid", target_ip="10.10.0.12"} 52
			picks_counter{datacenter_id="dc2-uuid", domain="foo.com", project_id="p1-uuid", target_ip="10.10.0.13"} 53
			`)
		err = testutil.CollectAndCompare(picksCounter, expected, "picks_counter")
		// if the following assertion fails, the tip to debugging the issue is to visually identify the embedded diff (-/+ leading characters) in the error string
		// TODO: make this diff obvious in case of failure
		assert.Nil(err)
	})
}
