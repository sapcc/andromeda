// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

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
