// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

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
