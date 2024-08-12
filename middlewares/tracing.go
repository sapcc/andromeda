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
