// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package bootstrap

import (
	"github.com/apex/log"
	"github.com/jmoiron/sqlx"

	andromedaDB "github.com/sapcc/andromeda/db"
	"github.com/sapcc/andromeda/internal/config"
)

// UpsertF5Datacenters ensures that all F5 datacenters defined in the config file
// are present and up to date in the `datacenter` table; only upserts are
// performed (no deletes)
func UpsertF5Datacenters(db *sqlx.DB, f5DCs []config.F5Datacenter) {
	if len(f5DCs) == 0 {
		return
	}

	log.Infof("Upserting %d F5 datacenters into the database", len(f5DCs))
	if err := andromedaDB.TxExecute(db, func(tx *sqlx.Tx) error {
		insertStmt, err := tx.PrepareNamed(`
			INSERT INTO datacenter 
				(name, continent, country, state_or_province, city, latitude, longitude, 
				provisioning_status, scope, admin_state_up, project_id, provider)
			VALUES 
				(:name, :continent, :country, :state_or_province, :city, :latitude, :longitude, 
				'ACTIVE', 'public', TRUE, '-', 'f5')
			;
		`)
		if err != nil {
			return err
		}

		updateStmt, err := tx.PrepareNamed(`
			UPDATE datacenter 
			SET 
				city = :city,
				continent = :continent,
				country = :country,
				latitude = :latitude,
				longitude = :longitude,
				state_or_province = :state_or_province,
				provider = 'f5',
				provisioning_status = 'ACTIVE',
				scope = 'public',
				admin_state_up = TRUE,
				project_id = '-',
				updated_at = NOW()
			WHERE 
				id = (SELECT id FROM datacenter WHERE provider = 'f5' AND name = :name)
		`)
		if err != nil {
			return err
		}

		for _, f5DC := range f5DCs {
			res, err := updateStmt.Exec(f5DC)
			if err != nil {
				return err
			}
			if rows, err := res.RowsAffected(); err != nil || rows == 0 {
				// no rows updated, try to insert
				if _, err := insertStmt.Exec(f5DC); err != nil {
					return err
				}
			} else if err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		log.Errorf("Failed to upsert F5 datacenters: %s", err)
	} else {
		log.Infof("Successfully upserted %d F5 datacenters into the database", len(f5DCs))
	}
}
