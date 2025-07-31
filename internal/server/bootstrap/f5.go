package bootstrap

import (
	"github.com/apex/log"
	"github.com/jmoiron/sqlx"
	"github.com/sapcc/andromeda/internal/config"
)

// TODO: upsert will not be possible due to missing UNIQUE KEY INDEX.
// TODO: perform two queries: select matching provider/name, then either insert or update
// TODO: use compatible SQL syntax for MariaDB/Postgres
// UpsertF5Datacenters ensures that all F5 datacenters defined in the config file
// are present and up to date in the `datacenter` table; only upserts are
// performed (no deletes)
func UpsertF5Datacenters(db *sqlx.DB, f5DCs []config.F5Datacenter) {
	if len(f5DCs) == 0 {
		log.Warn("No F5 datacenters provided in config file; skipping 'upsert F5 datacenters' operation")
		return
	}
	log.Infof("Upserting %d F5 datacenters...", len(f5DCs))
	// MariaDB syntax
	query := `INSERT INTO datacenter
	    (admin_state_up, project_id, provider, name, city, continent, country)
	    VALUES (1, "f5", 0, :name, :city, :continent, :country)`
	for _, f5DC := range f5DCs {
		if _, err := db.NamedExec(query, f5DC); err != nil {
			log.Errorf("Could not upsert F5 datacenter %q: %s", f5DC.Name, err)
			continue
		}
		log.Infof("Successfully upserted F5 datacenter %q", f5DC.Name)
	}
}
