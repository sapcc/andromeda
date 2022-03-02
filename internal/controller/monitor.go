/*
 *   Copyright 2020 SAP SE
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

package controller

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"

	"github.com/jackc/pgconn"

	"github.com/asim/go-micro/v3"
	"github.com/go-openapi/runtime/middleware"
	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"

	"github.com/sapcc/andromeda/db"
	"github.com/sapcc/andromeda/internal/auth"
	"github.com/sapcc/andromeda/internal/policy"
	"github.com/sapcc/andromeda/internal/utils"
	"github.com/sapcc/andromeda/models"
	"github.com/sapcc/andromeda/restapi/operations/monitors"
)

type MonitorController struct {
	db *sqlx.DB
	sv micro.Service
}

//GetMonitors GET /monitors
func (c MonitorController) GetMonitors(params monitors.GetMonitorsParams) middleware.Responder {
	pagination := db.NewPagination("monitor", params.Limit, params.Marker, params.Sort, params.PageReverse)
	rows, err := pagination.Query(c.db, params.HTTPRequest, nil)
	if err != nil {
		if errors.Is(err, db.ErrInvalidMarker) {
			return monitors.NewGetMonitorsBadRequest().WithPayload(InvalidMarker)
		}
		if errors.Is(err, db.ErrPolicyForbidden) {
			return GetPolicyForbiddenResponse()
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
		_monitors = append(_monitors, &monitor)
	}
	_links := pagination.GetLinks(_monitors, params.HTTPRequest)
	payload := monitors.GetMonitorsOKBody{Monitors: _monitors, Links: _links}
	return monitors.NewGetMonitorsOK().WithPayload(&payload)
}

//PostMonitors POST /monitors
func (c MonitorController) PostMonitors(params monitors.PostMonitorsParams) middleware.Responder {
	monitor := params.Monitor.Monitor
	projectID, err := auth.ProjectScopeForRequest(params.HTTPRequest)
	if err != nil {
		panic(err)
	}
	if !policy.Engine.AuthorizeRequest(params.HTTPRequest, projectID) {
		return GetPolicyForbiddenResponse()
	}
	monitor.ProjectID = &projectID

	pool := models.Pool{ID: *monitor.PoolID}
	if err := PopulatePool(c.db, &pool, []string{"project_id"}, false); err != nil || *pool.ProjectID != projectID {
		return monitors.NewPostMonitorsNotFound().WithPayload(GetErrorPoolNotFound(monitor.PoolID))
	}

	// Set default values
	if err := utils.SetModelDefaults(monitor); err != nil {
		panic(err)
	}

	if err := db.TxExecute(c.db, func(tx *sqlx.Tx) error {
		sql := `
			INSERT INTO monitor 
				(name, admin_state_up, type, "interval", timeout, pool_id, send, receive, project_id)
			VALUES 
				(:name, :admin_state_up, :type, :interval, :timeout, :pool_id, :send, :receive, :project_id)
			RETURNING *
		`
		stmt, err := tx.PrepareNamed(sql)
		if err != nil {
			return err

		}
		if err := stmt.Get(monitor, monitor); err != nil {
			return err
		}

		// Invalidate parent pool and domains
		sql = `UPDATE pool SET provisioning_status = 'PENDING_UPDATE' where id = ?`
		if _, err := tx.Exec(sql, monitor.PoolID); err != nil {
			return err
		}
		sql = `
				UPDATE domain d 
				JOIN domain_pool_relation dpr 
				    ON d.id = dpr.domain_id 
				SET provisioning_status = 'PENDING_UPDATE' where dpr.pool_id = ?`
		if _, err := tx.Exec(sql, monitor.PoolID); err != nil {
			return err
		}

		return nil
	}); err != nil {
		var pe *pgconn.PgError
		if errors.As(err, &pe) && pe.Code == pgerrcode.UniqueViolation {
			return monitors.NewPostMonitorsBadRequest().WithPayload(GetErrorPoolHasAlreadyAMonitor(monitor.PoolID))
		}
		var me *mysql.MySQLError
		if errors.As(err, &me) && me.Number == 1062 {
			return monitors.NewPostMonitorsBadRequest().WithPayload(GetErrorPoolHasAlreadyAMonitor(monitor.PoolID))
		}
		panic(err)
	}

	return monitors.NewPostMonitorsCreated().WithPayload(&monitors.PostMonitorsCreatedBody{Monitor: monitor})
}

//GetMonitorsMonitorID GET /monitors/:id
func (c MonitorController) GetMonitorsMonitorID(params monitors.GetMonitorsMonitorIDParams) middleware.Responder {
	monitor := models.Monitor{ID: params.MonitorID}
	if err := PopulateMonitor(c.db, &monitor, []string{"*"}); err != nil {
		return monitors.NewGetMonitorsMonitorIDNotFound().WithPayload(NotFound)
	}

	if !policy.Engine.AuthorizeRequest(params.HTTPRequest, *monitor.ProjectID) {
		return GetPolicyForbiddenResponse()
	}
	return monitors.NewGetMonitorsMonitorIDOK().WithPayload(&monitors.GetMonitorsMonitorIDOKBody{Monitor: &monitor})
}

//PutMonitorsMonitorID PUT /monitors/:id
func (c MonitorController) PutMonitorsMonitorID(params monitors.PutMonitorsMonitorIDParams) middleware.Responder {
	monitor := models.Monitor{ID: params.MonitorID}
	if err := PopulateMonitor(c.db, &monitor, []string{"project_id", "pool_id"}); err != nil {
		return monitors.NewPutMonitorsMonitorIDNotFound().WithPayload(NotFound)
	}
	if !policy.Engine.AuthorizeRequest(params.HTTPRequest, *monitor.ProjectID) {
		return GetPolicyForbiddenResponse()
	}

	if *params.Monitor.Monitor.PoolID != *monitor.PoolID {
		return monitors.NewPutMonitorsMonitorIDBadRequest().WithPayload(PoolIDImmutable)
	}

	params.Monitor.Monitor.ID = params.MonitorID
	sql := `
		UPDATE monitor SET
			name = COALESCE(:name, name),
			admin_state_up = COALESCE(:admin_state_up, admin_state_up),
		    "interval" = COALESCE(:interval, "interval"),
			receive = COALESCE(:receive, receive),
			send = COALESCE(:send, send),
			timeout = COALESCE(:timeout, timeout),
			type = COALESCE(:type, type),
		    updated_at = NOW()
		WHERE id = :id
	`

	if _, err := c.db.NamedExec(sql, params.Monitor.Monitor); err != nil {
		panic(err)
	}
	if err := PopulateMonitor(c.db, &monitor, []string{"*"}); err != nil {
		panic(err)
	}
	return monitors.NewPutMonitorsMonitorIDAccepted().WithPayload(
		&monitors.PutMonitorsMonitorIDAcceptedBody{Monitor: &monitor})
}

//DeleteMonitorsMonitorID DELETE /monitors/:id
func (c MonitorController) DeleteMonitorsMonitorID(params monitors.DeleteMonitorsMonitorIDParams) middleware.Responder {
	monitor := models.Monitor{ID: params.MonitorID}
	if err := PopulateMonitor(c.db, &monitor, []string{"project_id"}); err != nil {
		return monitors.NewDeleteMonitorsMonitorIDNotFound().WithPayload(NotFound)
	}
	if !policy.Engine.AuthorizeRequest(params.HTTPRequest, *monitor.ProjectID) {
		return GetPolicyForbiddenResponse()
	}

	sql := `DELETE FROM monitor WHERE id = ?`
	res := c.db.MustExec(sql, params.MonitorID)
	if deleted, _ := res.RowsAffected(); deleted != 1 {
		monitors.NewDeleteMonitorsMonitorIDNotFound().WithPayload(NotFound)
	}
	return monitors.NewDeleteMonitorsMonitorIDNoContent()
}

//PopulateMonitor populates attributes of a monitor instance based on it's ID
func PopulateMonitor(db *sqlx.DB, monitor *models.Monitor, fields []string) error {
	sql := fmt.Sprintf(`SELECT %s FROM monitor WHERE id = ?`, strings.Join(fields, ", "))
	if err := db.Get(monitor, sql, monitor.ID); err != nil {
		return err
	}
	return nil
}
