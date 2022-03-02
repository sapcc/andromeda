/*
 *   Copyright 2021 SAP SE
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

package cli

import (
	"os"

	"github.com/asim/go-micro/v3/logger"
	"github.com/jedib0t/go-pretty/table"
	"github.com/jessevdk/go-flags"
	"github.com/jmoiron/sqlx/reflectx"

	"github.com/sapcc/andromeda/client"
)

var (
	Parser          = flags.NewParser(&opts, flags.Default)
	Table           = table.NewWriter()
	Mapper          = reflectx.NewMapper("json")
	AndromedaClient *client.Andromeda
)
var opts struct {
	Debug      bool             `short:"d" long:"debug" description:"Show verbose debug information"`
	Formatters outputFormatters `group:"Output formatters"`
}

func SetupClient() {
	Table.SetOutputMirror(os.Stdout)

	transportConfig := client.DefaultTransportConfig()
	transportConfig.WithHost("localhost:8000")
	AndromedaClient = client.NewHTTPClientWithConfig(nil, transportConfig)

	if _, err := Parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			} else {
				logger.Error(err)
			}
		}
		os.Exit(code)
	}
}
