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

	if u.Driver == "mysql" {
		// Create driver
		driver, err = mysql.New(u.DSN)
		embedSource.Dir = "migrations/mysql"
		if err != nil {
			return err
		}

	} else if u.Driver == "postgres" {
		driver, err = postgres.New(u.DSN)
		embedSource.Dir = "migrations/postgresql"
		if err != nil {
			return err
		}
	} else {
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

	if u.Driver == "mysql" {
		// Create driver
		driver, err = mysql.New(u.DSN)
		embedSource.Dir = "migrations/mysql"
		if err != nil {
			return err
		}

	} else if u.Driver == "postgres" {
		driver, err = postgres.New(u.DSN)
		embedSource.Dir = "migrations/postgresql"
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("database driver not supported: %s", u.Driver)
	}

	// Run all up migrations
	if _, err := migration.Migrate(driver, embedSource, migration.Down, 0); err != nil {
		return err
	}

	return nil
}
