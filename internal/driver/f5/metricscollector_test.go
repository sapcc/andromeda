// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package f5

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestCollectVirtualServerMetrics(t *testing.T) {
	picksCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "picks_counter", Help: "Picks counter"},
		[]string{"domain", "datacenter_id", "project_id", "target_ip"},
	)
	prometheus.MustRegister(picksCounter)
	assert := assert.New(t)

	t.Run("Succeeds by incrementing the given counter metric accordingly", func(t *testing.T) {
		session := new(mockedBigIPSession)
		store := new(mockedStore)
		err := collectVirtualServerMetrics(session, store, picksCounter)
		assert.Nil(err)
		expected := strings.NewReader(`
			# HELP picks_counter Picks counter
			# TYPE picks_counter counter
			picks_counter{datacenter_id="dc1", domain="foo.com", project_id="p1", target_ip="10.10.0.1"} 1
			`)
		err = testutil.CollectAndCompare(picksCounter, expected, "picks_counter")
		// if the following assertion fails, the tip to debugging the issue is to visually identify the embedded diff (-/+ leading characters) in the error string
		// TODO: make this diff obvious in case of failure
		assert.Nil(err)
	})
}
