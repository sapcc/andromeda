// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"fmt"
	"log"

	"github.com/Boostport/migration"
	"github.com/Boostport/migration/driver/mysql"
	"github.com/Boostport/migration/driver/postgres"
	"github.com/xo/dburl"
)

func Migrate(dbUrl string) error {
	u, err := dburl.Parse(dbUrl)
	if err != nil {
		return err
	}
	migration.SetLogger(log.Default())

	embedSource := &migration.EmbedMigrationSource{
		EmbedFS: bindata,
	}
	var driver migration.Driver

	switch u.Driver {
	case "mysql":
		// Create driver
		driver, err = mysql.New(u.DSN)
		embedSource.Dir = "migrations/mysql"
		if err != nil {
			return err
		}
	case "postgres":
		driver, err = postgres.New(u.DSN)
		embedSource.Dir = "migrations/postgresql"
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("database driver not supported: %s", u.Driver)
	}

	// Run all up migrations
	if _, err := migration.Migrate(driver, embedSource, migration.Up, 0); err != nil {
		return err
	}

	return nil
}

func Rollback(dbUrl string) error {
	u, err := dburl.Parse(dbUrl)
	if err != nil {
		return err
	}
	migration.SetLogger(log.Default())

	embedSource := &migration.EmbedMigrationSource{
		EmbedFS: bindata,
	}
	var driver migration.Driver

	switch u.Driver {
	case "mysql":
		// Create driver
		driver, err = mysql.New(u.DSN)
		embedSource.Dir = "migrations/mysql"
		if err != nil {
			return err
		}
	case "postgres":
		driver, err = postgres.New(u.DSN)
		embedSource.Dir = "migrations/postgresql"
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("database driver not supported: %s", u.Driver)
	}

	// Run all up migrations
	if _, err := migration.Migrate(driver, embedSource, migration.Down, 0); err != nil {
		return err
	}

	return nil
}
