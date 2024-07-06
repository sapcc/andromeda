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
