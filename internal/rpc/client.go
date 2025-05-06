/*
 *   Copyright 2025 SAP SE
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
	"strings"

	"github.com/actatum/stormrpc"
	"github.com/apex/log"
	"github.com/sapcc/go-api-declarations/bininfo"

	"github.com/sapcc/andromeda/internal/config"
)

var rpcClient *stormrpc.Client

// GetRPCClient returns a singleton RPC client instance
func GetRPCClient() (*stormrpc.Client, error) {
	if rpcClient != nil {
		return rpcClient, nil
	}

	version := strings.TrimLeft(bininfo.VersionOr("-unknown"), "v")
	version = version[:strings.LastIndex(version, "-")]

	// Create client with just the NATS URL, other options like name and version are not needed
	client, err := stormrpc.NewClient(config.Global.Default.TransportURL)
	if err != nil {
		log.WithError(err).Error("Failed to create RPC client")
		return nil, err
	}

	rpcClient = client
	return rpcClient, nil
}
