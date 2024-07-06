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

package rpc

import (
	"context"

	"github.com/actatum/stormrpc"
	"github.com/actatum/stormrpc/middleware"
	"github.com/apex/log"
	"github.com/sapcc/go-api-declarations/bininfo"

	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/middlewares"
)

func NewServer(name string, opts ...stormrpc.ServerOption) *stormrpc.Server {
	opts = append(opts, stormrpc.WithErrorHandler(func(ctx context.Context, err error) {
		log.WithError(err).Error("RPC Error")
	}))
	srv, err := stormrpc.NewServer(&stormrpc.ServerConfig{
		NatsURL: config.Global.Default.TransportURL,
		Name:    name,
		Version: bininfo.Version(),
	}, opts...)
	if err != nil {
		log.Fatal(err.Error())
	}

	if srv != nil {
		log.Info("Loading RPC middleware")
		srv.Use(middleware.RequestID, middlewares.Logging, middleware.Recoverer)
	}
	return srv
}
