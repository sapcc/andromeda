// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"

	"github.com/sapcc/andromeda/db"
	"github.com/sapcc/andromeda/internal/auth"
	"github.com/sapcc/andromeda/internal/utils"
	"github.com/sapcc/andromeda/models"
	"github.com/sapcc/andromeda/restapi/operations/monitors"
)

type MonitorController struct {
	CommonController
}

// GetMonitors GET /monitors
func (c MonitorController) GetMonitors(params monitors.GetMonitorsParams) middleware.Responder {
	filter := make(map[string]any, 0)
	pagination := db.Pagination{
		HTTPRequest: params.HTTPRequest,
		Limit:       params.Limit,
		Marker:      params.Marker,
		PageReverse: params.PageReverse,
		Sort:        params.Sort,
	}
	if params.PoolID != nil {
		filter["pool_id"] = params.PoolID
	}
	rows, err := pagination.Query(c.db, "SELECT * from monitor", filter)
	if err != nil {
		if errors.Is(err, db.ErrInvalidMarker) {
			return monitors.NewGetMonitorsBadRequest().WithPayload(utils.InvalidMarker)
		}
		panic(err)
	}

	//goland:noinspection GoPreferNilSlice
	var _monitors = []*models.Monitor{}
	for rows.Next() {
		var monitor models.Monitor
		if err := rows.StructScan(&monitor); err != nil {
			panic(err)
		}
		// Filter result based on policy
		requestVars := map[string]string{"project_id": *monitor.ProjectID}
		if _, err = auth.Authenticate(params.HTTPRequest, requestVars); err == nil {
			_monitors = append(_monitors, &monitor)
		}
	}
	_links := pagination.GetLinks(_monitors)
	payload := monitors.GetMonitorsOKBody{Monitors: _monitors, Links: _links}
	return monitors.NewGetMonitorsOK().WithPayload(&payload)
}

// PostMonitors POST /monitors
func (c MonitorController) PostMonitors(params monitors.PostMonitorsParams) middleware.Responder {
	monitor := params.Monitor.Monitor
	if monitor.PoolID == nil {
		return monitors.NewPostMonitorsBadRequest().WithPayload(utils.PoolIDRequired)
	}

	if err := validateMonitor(monitor); err != nil {
		return monitors.NewPostMonitorsBadRequest().WithPayload(err)
	}

	if projectID, err := auth.Authenticate(params.HTTPRequest, nil); err != nil {
		return monitors.NewPostMonitorsDefault(403).WithPayload(utils.PolicyForbidden)
	} else if projectID != "" {
		monitor.ProjectID = &projectID
	}

	pool := models.Pool{ID: *monitor.PoolID}
	if err := PopulatePool(c.db, &pool, []string{"project_id"}, false); err != nil || *pool.ProjectID != *monitor.ProjectID {
		return monitors.NewPostMonitorsNotFound().WithPayload(utils.GetErrorPoolNotFound(monitor.PoolID))
	}

	// Set default values
	if err := utils.SetModelDefaults(monitor); err != nil {
		panic(err)
	}

	if err := db.TxExecute(c.db, func(tx *sqlx.Tx) error {
		sql := `
			INSERT INTO monitor 
				(name, admin_state_up, type, "interval", timeout, pool_id, send, receive, http_method, domain_name, project_id)
			VALUES 
				(:name, :admin_state_up, :type, :interval, :timeout, :pool_id, :send, :receive, :http_method, :domain_name, :project_id)
			RETURNING *
		`
		stmt, err := tx.PrepareNamed(sql)
		if err != nil {
			return err

		}
		if err := stmt.Get(monitor, monitor); err != nil {
			return err
		}

		return UpdateCascadePool(tx, *monitor.PoolID, "PENDING_UPDATE")
	}); err != nil {
		var pe *pgconn.PgError
		if errors.As(err, &pe) && pe.Code == pgerrcode.UniqueViolation {
			return monitors.NewPostMonitorsBadRequest().WithPayload(utils.GetErrorPoolHasAlreadyAMonitor(monitor.PoolID))
		}
		var me *mysql.MySQLError
		if errors.As(err, &me) && me.Number == 1062 {
			return monitors.NewPostMonitorsBadRequest().WithPayload(utils.GetErrorPoolHasAlreadyAMonitor(monitor.PoolID))
		}
		panic(err)
	}

	_ = PendingSync(c.rpc)
	return monitors.NewPostMonitorsCreated().WithPayload(&monitors.PostMonitorsCreatedBody{Monitor: monitor})
}

// GetMonitorsMonitorID GET /monitors/:id
func (c MonitorController) GetMonitorsMonitorID(params monitors.GetMonitorsMonitorIDParams) middleware.Responder {
	monitor := models.Monitor{ID: params.MonitorID}
	if err := PopulateMonitor(c.db, &monitor, []string{"*"}); err != nil {
		return monitors.NewGetMonitorsMonitorIDNotFound().WithPayload(utils.NotFound)
	}
	requestVars := map[string]string{"project_id": *monitor.ProjectID}
	if _, err := auth.Authenticate(params.HTTPRequest, requestVars); err != nil {
		return monitors.NewGetMonitorsMonitorIDDefault(403).WithPayload(utils.PolicyForbidden)
	}
	return monitors.NewGetMonitorsMonitorIDOK().WithPayload(&monitors.GetMonitorsMonitorIDOKBody{Monitor: &monitor})
}

// PutMonitorsMonitorID PUT /monitors/:id
func (c MonitorController) PutMonitorsMonitorID(params monitors.PutMonitorsMonitorIDParams) middleware.Responder {
	monitor := models.Monitor{ID: params.MonitorID}

	if err := validateMonitor(params.Monitor.Monitor); err != nil {
		return monitors.NewPutMonitorsMonitorIDBadRequest().WithPayload(err)
	}

	if err := PopulateMonitor(c.db, &monitor, []string{"project_id", "pool_id"}); err != nil {
		return monitors.NewPutMonitorsMonitorIDNotFound().WithPayload(utils.NotFound)
	}
	requestVars := map[string]string{"project_id": *monitor.ProjectID}
	if _, err := auth.Authenticate(params.HTTPRequest, requestVars); err != nil {
		return monitors.NewPutMonitorsMonitorIDDefault(403).WithPayload(utils.PolicyForbidden)
	}

	if params.Monitor.Monitor.PoolID != nil && *params.Monitor.Monitor.PoolID != *monitor.PoolID {
		return monitors.NewPutMonitorsMonitorIDBadRequest().WithPayload(utils.PoolIDImmutable)
	}

	params.Monitor.Monitor.ID = params.MonitorID
	if err := db.TxExecute(c.db, func(tx *sqlx.Tx) error {
		sql := `
		UPDATE monitor SET
			name = COALESCE(:name, name),
			admin_state_up = COALESCE(:admin_state_up, admin_state_up),
		    "interval" = COALESCE(:interval, "interval"),
			receive = COALESCE(:receive, receive),
			send = COALESCE(:send, send),
			timeout = COALESCE(:timeout, timeout),
			type = COALESCE(:type, type),
			http_method = COALESCE(:http_method, http_method),
			domain_name = COALESCE(:domain_name, domain_name),
		    updated_at = NOW(),
		    provisioning_status = 'PENDING_UPDATE'
		WHERE id = :id
	`

		if _, err := c.db.NamedExec(sql, params.Monitor.Monitor); err != nil {
			return err
		}
		return UpdateCascadePool(tx, *monitor.PoolID, "PENDING_UPDATE")
	}); err != nil {
		panic(err)
	}

	if err := PopulateMonitor(c.db, &monitor, []string{"*"}); err != nil {
		panic(err)
	}
	_ = PendingSync(c.rpc)
	return monitors.NewPutMonitorsMonitorIDAccepted().WithPayload(
		&monitors.PutMonitorsMonitorIDAcceptedBody{Monitor: &monitor})
}

// DeleteMonitorsMonitorID DELETE /monitors/:id
func (c MonitorController) DeleteMonitorsMonitorID(params monitors.DeleteMonitorsMonitorIDParams) middleware.Responder {
	monitor := models.Monitor{ID: params.MonitorID}
	if err := PopulateMonitor(c.db, &monitor, []string{"project_id", "pool_id"}); err != nil {
		return monitors.NewDeleteMonitorsMonitorIDNotFound().WithPayload(utils.NotFound)
	}
	requestVars := map[string]string{"project_id": *monitor.ProjectID}
	if _, err := auth.Authenticate(params.HTTPRequest, requestVars); err != nil {
		return monitors.NewDeleteMonitorsMonitorIDDefault(403).WithPayload(utils.PolicyForbidden)
	}

	if err := db.TxExecute(c.db, func(tx *sqlx.Tx) error {
		sql := tx.Rebind(`UPDATE monitor SET provisioning_status = 'PENDING_DELETE', updated_at = NOW() WHERE id = ?`)
		res := c.db.MustExec(sql, params.MonitorID)
		if deleted, _ := res.RowsAffected(); deleted != 1 {
			return ErrEmptyResult
		}
		return UpdateCascadePool(tx, *monitor.PoolID, "PENDING_UPDATE")
	}); err != nil {
		if errors.Is(err, ErrEmptyResult) {
			return monitors.NewDeleteMonitorsMonitorIDNotFound().WithPayload(utils.NotFound)
		}
		panic(err)
	}
	_ = PendingSync(c.rpc)
	return monitors.NewDeleteMonitorsMonitorIDNoContent()
}

// PopulateMonitor populates attributes of a monitor instance based on it's ID
func PopulateMonitor(db *sqlx.DB, monitor *models.Monitor, fields []string) error {
	sql := fmt.Sprintf(`SELECT %s FROM monitor WHERE id = ?`, strings.Join(fields, ", "))
	if err := db.Get(monitor, db.Rebind(sql), monitor.ID); err != nil {
		return err
	}
	return nil
}

// validateMonitor 2nd pass checks if the monitor is valid
func validateMonitor(monitor *models.Monitor) *models.Error {
	if monitor.Send != nil && (*monitor.Type == models.MonitorTypeHTTP || *monitor.Type == models.MonitorTypeHTTPS) {
		u, err := url.ParseRequestURI(*monitor.Send)
		if err != nil || u.Host != "" || u.Scheme != "" {
			return utils.InvalidSendString
		}
	}
	return nil
}
