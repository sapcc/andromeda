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
	insertQuery := db.Rebind(`INSERT INTO datacenter (admin_state_up, project_id, provider, name, city, continent, country)
	                VALUES (1, 0, "f5", ?, ?, ?, ?)`)
	selectQuery := db.Rebind(`SELECT id FROM datacenter WHERE provider = "f5" AND name = ? LIMIT 1`)
	updateQuery := db.Rebind(`UPDATE datacenter SET city = ?, continent = ?, country = ? WHERE id = ? LIMIT 1`)
	log.Infof("Upserting %d F5 datacenters...", len(f5DCs))
	for _, f5DC := range f5DCs {
		row := db.QueryRow(selectQuery, f5DC.Name)
		var id string
		if err := row.Scan(&id); err != nil {
			if _, err := db.Exec(insertQuery, f5DC.Name, f5DC.City, f5DC.Continent, f5DC.Country); err != nil {
				log.Errorf("Could not create F5 datacenter %q: %s", f5DC.Name, err)
				continue
			}
			log.Infof("Inserted F5 datacenter %q", f5DC.Name)
			continue
		}
		if _, err := db.Exec(updateQuery, f5DC.City, f5DC.Continent, f5DC.Country, id); err != nil {
			log.Errorf("Could not update F5 datacenter %q: %s", f5DC.Name, err)
			continue
		}
		log.Infof("Updated F5 datacenter %q", f5DC.Name)
	}
}
