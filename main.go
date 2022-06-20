/*
 *   Copyright 2020 SAP SE
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
	"os"

	"github.com/jessevdk/go-flags"
	"go-micro.dev/v4/logger"

	"github.com/sapcc/andromeda/db"
	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/driver/akamai"
	"github.com/sapcc/andromeda/internal/driver/f5"
	"github.com/sapcc/andromeda/internal/server"
	"github.com/sapcc/andromeda/restapi"
)

func main() {
	parser := flags.NewParser(&config.Global, flags.Default|flags.IgnoreUnknown)
	parser.ShortDescription = "Andromeda"
	parser.LongDescription = "Platform agnostic GSLB frontend"

	API := restapi.NewServer(nil)
	if _, err := parser.AddCommand("server",
		"Run andromeda server",
		"Runs the andromeda server component (API + RPC)",
		API); err != nil {
		logger.Fatal(err)
	}
	if _, err := parser.AddCommand("f5-agent",
		"Run andromeda f5 agent",
		"Runs the andromeda f5 Agent component", &struct{}{}); err != nil {
		logger.Fatal(err)
	}
	if _, err := parser.AddCommand("akamai-agent",
		"Run andromeda akamai agent",
		"Runs the andromeda akamai Agent component", &struct{}{}); err != nil {
		logger.Fatal(err)
	}
	if _, err := parser.AddCommand("f5-status-agent",
		"Run andromeda f5 status agent",
		"Runs the andromeda f5 Status Agent component", &struct{}{}); err != nil {
		logger.Fatal(err)
	}
	if _, err := parser.AddCommand("migrate",
		"Run andromeda database migration",
		"Runs the automatic database migration for andromeda", &struct{}{}); err != nil {
		logger.Fatal(err)
	}

	if _, err := parser.Parse(); err != nil {
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

	if len(config.Global.ConfigFile) > 0 {
		iniParser := flags.NewIniParser(parser)
		for _, file := range config.Global.ConfigFile {
			if err := iniParser.ParseFile(file); err != nil {
				logger.Fatal(err)
			}
		}
	}

	// Custom command handler for restapi server
	switch parser.Command.Active.Name {
	case "server":
		if err := server.ExecuteServer(API); err != nil {
			logger.Fatal(err)
		}
	case "f5-status-agent":
		if err := f5.ExecuteF5StatusAgent(); err != nil {
			logger.Fatal(err)
		}
	case "f5-agent":
		if err := f5.ExecuteF5Agent(); err != nil {
			logger.Fatal(err)
		}
	case "akamai-agent":
		if err := akamai.ExecuteAkamaiAgent(); err != nil {
			logger.Fatal(err)
		}
	case "migrate":
		dbUrl := config.Global.Database.Connection
		if err := db.Migrate(dbUrl); err != nil {
			logger.Fatal(err)
		}
	}
}
