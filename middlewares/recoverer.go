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
	"time"

	"github.com/actatum/stormrpc"
	"github.com/getsentry/sentry-go"
)

// Recoverer is a stormrpc middleware that recovers from panics and sends them to Sentry.
func Recoverer(next stormrpc.HandlerFunc) stormrpc.HandlerFunc {
	return func(ctx context.Context, r stormrpc.Request) (resp stormrpc.Response) {
		defer func() {
			err := recover()
			if err != nil {
				// Log the panic to Sentry
				sentry.CurrentHub().RecoverWithContext(ctx, err)
				sentry.Flush(time.Second * 5)

				// pass back to the caller
				resp = stormrpc.NewErrorResponse(
					r.Reply,
					stormrpc.Errorf(stormrpc.ErrorCodeInternal, "%v", err),
				)
			}
		}()

		resp = next(ctx, r)

		return resp
	}
}
