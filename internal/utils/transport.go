/*
 *   Copyright 2022 SAP SE
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

package utils

import (
	"fmt"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/server"
	"net/url"

	nats_b "github.com/go-micro/plugins/v4/broker/nats"
	nats_r "github.com/go-micro/plugins/v4/registry/nats"
	nats_t "github.com/go-micro/plugins/v4/transport/nats"
	"go-micro.dev/v4/broker"
	"go-micro.dev/v4/transport"

	"github.com/sapcc/andromeda/internal/config"
)

func ConfigureTransport() micro.Option {
	return func(o *micro.Options) {
		uri, err := url.Parse(config.Global.Default.TransportURL)
		if err != nil {
			panic(err)
		}

		switch uri.Scheme {
		case "http":
			o.Transport = transport.NewHTTPTransport(
				transport.Addrs(config.Global.Default.TransportURL))
			o.Broker = broker.NewBroker(
				broker.Addrs(config.Global.Default.TransportURL))
			o.Registry = registry.NewRegistry(
				registry.Addrs(config.Global.Default.TransportURL))
		case "nats":
			o.Transport = nats_t.NewTransport(
				transport.Addrs(uri.Host))
			o.Broker = nats_b.NewBroker(
				broker.Addrs(uri.Host))
			o.Registry = nats_r.NewRegistry(
				registry.Addrs(uri.Host))
		default:
			panic(fmt.Errorf("unknown scheme for transport_url %s", uri.Scheme))
		}
		if err := o.Client.Init(client.Transport(o.Transport)); err != nil {
			panic(err)
		}
		if err := o.Server.Init(server.Transport(o.Transport)); err != nil {
			panic(err)
		}
		if err := o.Client.Init(client.Broker(o.Broker)); err != nil {
			panic(err)
		}
		if err := o.Server.Init(server.Broker(o.Broker)); err != nil {
			panic(err)
		}
		if err := o.Client.Init(client.Registry(o.Registry)); err != nil {
			panic(err)
		}
		if err := o.Server.Init(server.Registry(o.Registry)); err != nil {
			panic(err)
		}
		// Update Broker
		if err := o.Broker.Init(broker.Registry(o.Registry)); err != nil {
			panic(err)
		}
	}
}
