// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package middlewares

import (
	"context"
	"time"

	"github.com/actatum/stormrpc"
	"github.com/actatum/stormrpc/middleware"
	"github.com/apex/log"
)

func Logging(next stormrpc.HandlerFunc) stormrpc.HandlerFunc {
	return func(ctx context.Context, r stormrpc.Request) stormrpc.Response {
		id := middleware.RequestIDFromContext(ctx)
		start := time.Now()

		resp := next(ctx, r)

		logCtx := log.WithFields(log.Fields{
			"request_id": id,
			"duration":   time.Since(start).String(),
			"subject":    r.Msg.Subject,
		})

		msg := "Success"
		if resp.Err != nil {
			msg = "Server Error"
			code := stormrpc.CodeFromErr(resp.Err)
			logCtx.WithError(resp.Err).WithField("code", code.String()).Error(msg)
		} else {
			logCtx.Debug(msg)
		}
		return resp
	}
}
