// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/urfave/cli/v2"

	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/server"
	"github.com/sapcc/andromeda/restapi"
)

func main() {
	config.ParseArgsAndRun("andromeda", "andromeda api server",
		func(c *cli.Context) error {
			API := restapi.NewServer(nil)
			API.Port = c.Int("port")
			API.Host = c.String("host")
			return server.ExecuteServer(API)
		},
		&cli.StringFlag{
			Name:  "host",
			Usage: "The IP to listen on",
			Value: "0.0.0.0",
		},
		&cli.IntFlag{
			Name:  "port",
			Usage: "Port to listen",
			Value: 8080,
		},
	)
}
