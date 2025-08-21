// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package bootstrap

import (
	"github.com/apex/log"
	"github.com/jmoiron/sqlx"
	"github.com/sapcc/andromeda/internal/config"
)

// UpsertF5Datacenters ensures that all F5 datacenters defined in the config file
// are present and up to date in the `datacenter` table; only upserts are
// performed (no deletes)
func UpsertF5Datacenters(db *sqlx.DB, f5DCs []config.F5Datacenter) {
	if len(f5DCs) == 0 {
		log.Warn("No F5 datacenters provided in config file; skipping 'upsert F5 datacenters' operation")
		return
	}
	tx, err := db.Beginx()
	if err != nil {
		log.Errorf("Failed to start UpsertF5Datacenters() transaction: %w", err)
		return
	}
	defer func() { _ = tx.Rollback() }()
	log.Info("Upserting F5 datacenters...")
	insertQuery := tx.Rebind(`INSERT INTO datacenter
                        (provisioning_status, admin_state_up, project_id,
			 provider, name, city, continent, country)
	                VALUES ("ACTIVE", 1, 0, "f5", ?, ?, ?, ?)`)
	selectQuery := tx.Rebind(`SELECT id FROM datacenter WHERE provider = "f5" AND name = ? LIMIT 1`)
	updateQuery := tx.Rebind(`UPDATE datacenter SET city = ?, continent = ?, country = ? WHERE id = ? LIMIT 1`)
	for _, f5DC := range f5DCs {
		row := tx.QueryRow(selectQuery, f5DC.Name)
		var id string
		if err := row.Scan(&id); err != nil {
			if _, err := tx.Exec(insertQuery, f5DC.Name, f5DC.City, f5DC.Continent, f5DC.Country); err != nil {
				log.Errorf("Could not create F5 datacenter %q: %s", f5DC.Name, err)
				return
			}
			continue
		}
		if _, err := tx.Exec(updateQuery, f5DC.City, f5DC.Continent, f5DC.Country, id); err != nil {
			log.Errorf("Could not update F5 datacenter %q: %s", f5DC.Name, err)
			return
		}
	}
	if err := tx.Commit(); err != nil {
		log.Errorf("Failed to commit UpsertF5Datacenters() transaction: %w", err)
		return
	}
	log.Infof("Successfully upserted F5 datacenters")
}
