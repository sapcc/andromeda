// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package middlewares

import (
	"context"

	"github.com/actatum/stormrpc"
	"github.com/prometheus/client_golang/prometheus"
)

// Tracing is a stormrpc middleware that records the duration of the request.
func Tracing(httpRequestHistogram *prometheus.HistogramVec) func(next stormrpc.HandlerFunc) stormrpc.HandlerFunc {
	return func(next stormrpc.HandlerFunc) stormrpc.HandlerFunc {
		return func(ctx context.Context, r stormrpc.Request) stormrpc.Response {
			if httpRequestHistogram != nil {
				timer := prometheus.NewTimer(httpRequestHistogram.WithLabelValues(r.Subject()))
				defer timer.ObserveDuration()
			}
			return next(ctx, r)
		}
	}
}
