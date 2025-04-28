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

package main

import (
	"github.com/urfave/cli/v3"

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
