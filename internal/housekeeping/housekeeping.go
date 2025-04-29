// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package housekeeping

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/apex/log"
	"github.com/go-openapi/strfmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/iancoleman/strcase"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sapcc/go-bits/jobloop"
	"github.com/xo/dburl"

	"github.com/sapcc/andromeda/internal/config"
)

type Executor struct {
	DB *sqlx.DB
}

func (e *Executor) findNextPoolToActivate(_ context.Context, tx *sqlx.Tx, _ prometheus.Labels) (*strfmt.UUID, error) {
	var poolID strfmt.UUID
	sql := `
		SELECT id 
		FROM pool
		LEFT JOIN domain_pool_relation dpr ON id = dpr.pool_id
		WHERE pool.provisioning_status LIKE 'PENDING_%' AND dpr.domain_id IS NULL
		LIMIT 1 FOR UPDATE`

	if tx.DriverName() != "mysql" {
		sql += " OF pool"
	}

	if err := tx.Get(&poolID, tx.Rebind(sql)); err != nil {
		return nil, err
	}

	return &poolID, nil
}

func (e *Executor) poolCascadeActive(_ context.Context, tx *sqlx.Tx, poolID *strfmt.UUID, _ prometheus.Labels) error {
	if _, err := tx.Exec(tx.Rebind(`UPDATE member SET provisioning_status = 'ACTIVE' WHERE pool_id = ?`),
		poolID); err != nil {
		return err
	}
	if _, err := tx.Exec(tx.Rebind(`UPDATE monitor SET provisioning_status = 'ACTIVE' WHERE pool_id = ?`),
		poolID); err != nil {
		return err
	}
	if _, err := tx.Exec(tx.Rebind(`UPDATE pool SET provisioning_status = 'ACTIVE' WHERE id = ?`),
		poolID); err != nil {
		return err
	}

	log.Infof("Unbound pool activation %s", poolID)
	return tx.Commit()
}

func (e *Executor) CleanupDeletedDomains(_ context.Context, labels prometheus.Labels) error {
	// Delete deleted domains
	sql := `DELETE FROM
           		domain 
			WHERE 
			    provisioning_status = 'DELETED' AND updated_at < (NOW() - INTERVAL %d SECOND)`
	res, err := e.DB.Exec(fmt.Sprintf(sql, config.Global.HouseKeeping.DeleteAfter))
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	labels["count"] = strconv.FormatInt(rowsAffected, 10)
	log.Infof("Cleaned up %d domains", rowsAffected)
	return nil
}

func (e *Executor) EventTranslationJob(registerer prometheus.Registerer) jobloop.Job {
	return (&jobloop.TxGuardedJob[*sqlx.Tx, *strfmt.UUID]{
		Metadata: jobloop.JobMetadata{
			ReadableName:    "unbound pool activation",
			ConcurrencySafe: true,
			CounterOpts:     prometheus.CounterOpts{Name: "unbound_pool_activation"},
			CounterLabels:   nil,
		},
		BeginTx:     e.DB.Beginx,
		DiscoverRow: e.findNextPoolToActivate,
		ProcessRow:  e.poolCascadeActive,
	}).Setup(registerer)
}

func (e *Executor) CleanupDeletedDomainsCronJob(registerer prometheus.Registerer) jobloop.Job {
	return (&jobloop.CronJob{
		Metadata: jobloop.JobMetadata{
			ReadableName:  "cleanup deleted domains",
			CounterOpts:   prometheus.CounterOpts{Name: "cleanup_deleted_domains"},
			CounterLabels: []string{"count"},
		},
		Interval: time.Second * time.Duration(config.Global.HouseKeeping.DeleteAfter),
		Task:     e.CleanupDeletedDomains,
	}).Setup(registerer)
}

func HouseKeeping() error {
	if !config.Global.HouseKeeping.Enabled {
		log.Fatal("Housekeeping disabled")
	}
	log.Info("Running housekeeping")

	// Database
	u, err := dburl.Parse(config.Global.Database.Connection)
	if err != nil {
		return err
	}
	if u.Driver == "postgres" {
		u.Driver = "pgx"
	}
	db, err := sqlx.Connect(u.Driver, u.DSN)
	if err != nil {
		log.WithError(err).WithField("driver", u.Driver).Fatal("Failed to connect to database")
	}

	// Mapper function for SQL name mapping, snake_case table names
	db.MapperFunc(strcase.ToSnake)

	executor := Executor{db}
	ctx, cancel := context.WithCancel(context.Background())
	go executor.EventTranslationJob(prometheus.DefaultRegisterer).Run(ctx)
	go executor.CleanupDeletedDomainsCronJob(prometheus.DefaultRegisterer).Run(ctx)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigs
	log.Infof("Got signal %s, terminating.", sig)
	cancel()
	return db.Close()
}
