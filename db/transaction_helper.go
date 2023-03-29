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
	"context"

	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbsqlx"
	"github.com/jmoiron/sqlx"
)

// TxExecute is a handy wrapper for safely executing queries inside a closure
func TxExecute(db *sqlx.DB, fn func(*sqlx.Tx) error) error {
	// CockroachDB transactions must be embedded within a special retry wrapper
	if db.DriverName() == "cockroachdb" {
		return crdbsqlx.ExecuteTx(context.Background(), db, nil, fn)
	}

	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		// Rollback transaction
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	return tx.Commit()
}
