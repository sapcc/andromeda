// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

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
