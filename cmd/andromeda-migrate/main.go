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
	"github.com/urfave/cli/v2"

	"github.com/sapcc/andromeda/db"
	"github.com/sapcc/andromeda/internal/config"
)

func main() {
	config.ParseArgsAndRun("andromeda-migrate", "database migration tool",
		func(c *cli.Context) error {
			dbUrl := config.Global.Database.Connection
			if c.String("db-url") != "" {
				dbUrl = c.String("db-url")
			}
			return db.Migrate(dbUrl)
		},
		&cli.StringFlag{
			Name:    "db-url",
			Usage:   "Database connection URL (overrides value provided in config file)",
			Value:   "",
			EnvVars: []string{"DB_URL"},
		},
	)
}
