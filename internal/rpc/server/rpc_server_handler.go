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

package server

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-openapi/strfmt"
	"net"
	"strings"

	"github.com/jmoiron/sqlx"
	"go-micro.dev/v4/logger"

	"github.com/sapcc/andromeda/db"
	"github.com/sapcc/andromeda/internal/rpcmodels"
	"github.com/sapcc/andromeda/internal/utils"
)

type RPCHandler struct {
	DB *sqlx.DB
}

// QueryxWithIds run sql query if optional WHERE condition based on IDs
func (u *RPCHandler) QueryxWithIds(sql string, request *SearchRequest) (*sqlx.Rows, error) {
	args := []interface{}{request.Provider}
	if len(request.Ids) > 0 {
		if strings.Contains(sql, "WHERE") {
			sql += ` AND id IN (?)`
		} else {
			sql += ` WHERE id IN (?)`
		}

		query, inArgs, err := sqlx.In(sql, request.Ids)
		if err != nil {
			return nil, err
		}
		sql = query
		args = append(inArgs, args...)
	}

	// Provider filter
	if strings.Contains(sql, "WHERE") {
		sql += ` AND provider = ?`
	} else {
		sql += ` WHERE provider = ?`
	}

	sql = u.DB.Rebind(sql)
	return u.DB.Queryx(sql, args...)
}

func (u *RPCHandler) GetMembers(ctx context.Context, request *SearchRequest, response *MembersResponse) error {
	sql := u.DB.Rebind(`SELECT id, admin_state_up, address, port, provisioning_status FROM member;`)
	rows, err := u.QueryxWithIds(sql, request)
	if err != nil {
		return err
	}
	for rows.Next() {
		var member rpcmodels.Member
		var address string
		if err := rows.Scan(&member.Id, &member.AdminStateUp, &address,
			&member.Port, &member.ProvisioningStatus); err != nil {
			return err
		}

		member.Address, err = utils.InetAton(net.ParseIP(address))
		if err != nil {
			return err
		}
		response.Response = append(response.Response, &member)
	}
	return nil
}

func (u *RPCHandler) GetPools(ctx context.Context, request *SearchRequest, response *PoolsResponse) error {
	sql := u.DB.Rebind(`SELECT id, admin_state_up FROM pool;`)
	rows, err := u.QueryxWithIds(sql, request)
	if err != nil {
		return err
	}
	for rows.Next() {
		var pool rpcmodels.Pool
		if err := rows.Scan(&pool.Id, &pool.AdminStateUp); err != nil {
			return err
		}
		response.Response = append(response.Response, &pool)
	}
	return nil
}

func (u *RPCHandler) GetDatacenters(ctx context.Context, request *SearchRequest, response *DatacentersResponse) error {
	q := sq.
		Select("id", "admin_state_up", "city", "state_or_province", "continent", "country", "latitude",
			"longitude", "scope", "provisioning_status", "provider", "meta").
		From("datacenter").
		Where("provider = ?", request.Provider)

	if request.Pending {
		q = q.Where("provisioning_status LIKE 'PENDING%'")
	}
	if request.Ids != nil {
		q = q.Where(sq.Eq{"id": request.Ids})
	}

	sql, args := q.MustSql()
	rows, err := u.DB.Queryx(u.DB.Rebind(sql), args...)
	if err != nil {
		return err
	}
	for rows.Next() {
		var datacenter rpcmodels.Datacenter
		if err := rows.StructScan(&datacenter); err != nil {
			logger.Error(err)
			return err
		}
		response.Response = append(response.Response, &datacenter)
	}
	return nil
}

func (u *RPCHandler) GetGeomaps(ctx context.Context, request *SearchRequest, response *GeomapsResponse) error {
	/*
		SELECT
		    id, scope, provider, default_datacenter, provisioning_status
		FROM geographic_map
		WHERE
		    provider = 'akamai' AND provisioning_status LIKE 'PENDING%'
	*/

	q := sq.
		Select("id", "default_datacenter").
		From("geographic_map").
		Where("provider = ?", request.Provider)

	if request.Pending {
		q = q.Where("provisioning_status LIKE 'PENDING%'")
	}
	if request.Ids != nil {
		q = q.Where(sq.Eq{"id": request.Ids})
	}

	if err := db.TxExecute(u.DB, func(tx *sqlx.Tx) error {
		sql, args := q.MustSql()
		rows, err := tx.Queryx(tx.Rebind(sql), args...)
		if err != nil {
			return err
		}
		for rows.Next() {
			var geomap rpcmodels.Geomap
			if err := rows.StructScan(&geomap); err != nil {
				return err
			}
			response.Response = append(response.Response, &geomap)
		}

		for _, geomap := range response.Response {
			var aggregate = "string_agg(country, ',') AS countries"
			if u.DB.DriverName() == "mysql" {
				aggregate = "GROUP_CONCAT(country) AS countries"
			}

			// Populate assignments
			sql, args := sq.Select("datacenter", aggregate).
				From("geographic_map_assignment").
				Where("geographic_map_id = ?", geomap.Id).
				GroupBy("datacenter").
				MustSql()
			rows, err := tx.Queryx(tx.Rebind(sql), args...)
			if err != nil {
				return err
			}
			for rows.Next() {
				var datacenter, countries string
				if err := rows.Scan(&datacenter, &countries); err != nil {
					return err
				}
				geomap.Assignment = append(geomap.Assignment, &rpcmodels.GeomapAssignment{
					Datacenter: datacenter,
					Countries:  strings.Split(countries, ","),
				})
			}
		}
		return nil
	}); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (u *RPCHandler) GetMonitors(ctx context.Context, request *SearchRequest, response *MonitorsResponse) error {
	panic("Todo")
}

func (u *RPCHandler) UpdateDatacenterMeta(ctx context.Context, req *DatacenterMetaRequest, res *rpcmodels.Datacenter) error {
	if err := db.TxExecute(u.DB, func(tx *sqlx.Tx) error {
		sql := tx.Rebind(`UPDATE datacenter SET meta = ? WHERE id = ?`)
		if _, err := tx.Exec(sql, req.GetMeta(), req.GetId()); err != nil {
			return err
		}

		sql = tx.Rebind(`SELECT id, admin_state_up, city, state_or_province, continent, country, 
               latitude, longitude, scope, provisioning_status, provider, meta FROM datacenter WHERE id = ?`)
		if err := tx.Get(res, sql, req.GetId()); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}
	return nil
}

func populatePools(u *RPCHandler, fullyPopulated bool, domainID string) ([]*rpcmodels.Pool, error) {
	sql := u.DB.Rebind(`SELECT id, admin_state_up, provisioning_status 
            FROM pool p 
            JOIN domain_pool_relation dpr ON p.id = dpr.pool_id
            WHERE dpr.domain_id = ?`)
	rows, err := u.DB.Queryx(sql, domainID)
	if err != nil {
		return nil, err
	}

	var pools []*rpcmodels.Pool
	for rows.Next() {
		var pool rpcmodels.Pool
		if err := rows.StructScan(&pool); err != nil {
			logger.Error(err)
			return nil, err
		}

		if fullyPopulated {
			if pool.Members, err = populateMembers(u, pool.Id); err != nil {
				return nil, err
			}
			if pool.Monitors, err = populateMonitors(u, pool.Id); err != nil {
				return nil, err
			}
		}

		pools = append(pools, &pool)
	}
	return pools, nil
}

func populateMonitors(u *RPCHandler, poolID string) ([]*rpcmodels.Monitor, error) {
	sql := u.DB.Rebind(`SELECT id, admin_state_up, "interval", send, receive, timeout, type, provisioning_status 
		FROM monitor WHERE pool_id = ?`)
	rows, err := u.DB.Queryx(sql, poolID)
	if err != nil {
		return nil, err
	}

	var monitors []*rpcmodels.Monitor
	for rows.Next() {
		var monitor rpcmodels.Monitor
		var send, receive *string
		var monitorType string

		if err := rows.Scan(&monitor.Id, &monitor.AdminStateUp, &monitor.Interval, &send,
			&receive, &monitor.Timeout, &monitorType, &monitor.ProvisioningStatus); err != nil {
			return nil, err
		}
		if send != nil {
			monitor.Send = *send
		}
		if receive != nil {
			monitor.Receive = *receive
		}
		monitor.Type = rpcmodels.Monitor_MonitorType(rpcmodels.Monitor_MonitorType_value[monitorType])
		monitors = append(monitors, &monitor)
	}
	return monitors, err
}

func populateMembers(u *RPCHandler, poolID string) ([]*rpcmodels.Member, error) {
	sql := u.DB.Rebind(`SELECT id, admin_state_up, address, port, datacenter_id, provisioning_status 
		FROM member WHERE pool_id = ?`)
	rows, err := u.DB.Queryx(sql, poolID)
	if err != nil {
		return nil, err
	}

	var members []*rpcmodels.Member
	for rows.Next() {
		var member rpcmodels.Member
		var address string
		if err := rows.Scan(&member.Id, &member.AdminStateUp, &address,
			&member.Port, &member.Datacenter, &member.ProvisioningStatus); err != nil {
			logger.Error(err)
			return nil, err
		}

		member.Address, err = utils.InetAton(net.ParseIP(address))
		if err != nil {
			return nil, err
		}

		members = append(members, &member)
	}
	return members, nil
}

func (u *RPCHandler) GetDomains(ctx context.Context, request *SearchRequest, response *DomainsResponse) error {
	sql := `SELECT id, admin_state_up, fqdn, mode, record_type, provisioning_status
               FROM domain WHERE provisioning_status != 'DELETED'`
	if request.Pending {
		sql += ` AND provisioning_status in ('PENDING_CREATE', 'PENDING_UPDATE', 'PENDING_DELETE')`
	}
	rows, err := u.QueryxWithIds(u.DB.Rebind(sql), request)

	if err != nil {
		logger.Error(err)
		return err
	}
	for rows.Next() {
		var domain rpcmodels.Domain
		if err := rows.StructScan(&domain); err != nil {
			logger.Error(err)
			return err
		}

		datacenterIds := []string{}

		if request.GetFullyPopulated() {
			if domain.Pools, err = populatePools(u, request.GetFullyPopulated(), domain.Id); err != nil {
				logger.Error(err)
				return err
			}
			for _, pool := range domain.Pools {
				for _, member := range pool.Members {
					found := false
					for _, datacenter := range datacenterIds {
						if datacenter == member.Datacenter {
							found = true
						}
					}
					if !found {
						datacenterIds = append(datacenterIds, member.Datacenter)
					}
				}
			}

			s := SearchRequest{
				PageNumber:    0,
				ResultPerPage: 100,
				Pending:       false,
				Ids:           datacenterIds,
				Provider:      request.Provider,
			}
			r := DatacentersResponse{}
			if err := u.GetDatacenters(ctx, &s, &r); err != nil {
				return err
			}
			domain.Datacenters = r.Response
		}

		response.Response = append(response.Response, &domain)
	}
	return nil
}

func (u *RPCHandler) UpdateProvisioningStatus(ctx context.Context, req *ProvisioningStatusRequest, res *ProvisioningStatusResponse) error {
	var statusResult []*StatusResult
	for _, provStatusReq := range req.GetProvisioningStatus() {
		table := strings.ToLower(provStatusReq.GetModel().String())
		provStatus := provStatusReq.GetStatus().String()

		var sql string
		var err error
		if provStatus == "DELETED" {
			sql = u.DB.Rebind(fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, table))
			_, err = u.DB.Exec(sql, provStatusReq.GetId())
		} else {
			sql = u.DB.Rebind(fmt.Sprintf(`UPDATE %s SET provisioning_status = ?, updated_at = NOW() WHERE id = ?`, table))
			if provStatus != "ERROR" {
				// only update non-error prov-status
				sql += " AND provisioning_status != 'ERROR'"
			}
			_, err = u.DB.Exec(sql, provStatus, provStatusReq.GetId())
		}
		if err != nil {
			logger.Error(err)
		}

		statusResult = append(statusResult, &StatusResult{
			Id:      provStatusReq.GetId(),
			Success: err == nil,
		})
	}
	res.ProvisioningStatusResult = statusResult
	return nil
}

// UpdateMemberStatus Updates member status according to the requests, also updates dependend pool and domain status.
func (u *RPCHandler) UpdateMemberStatus(ctx context.Context, req *MemberStatusRequest, res *MemberStatusResponse) error {
	var statusResult []*StatusResult
	tx, err := u.DB.Beginx()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	for _, memberStatusReq := range req.GetMemberStatus() {
		status := memberStatusReq.GetStatus().String()
		// Set all related objects to Online
		sql, args := sq.Update("member").
			Set("status", status).
			Where("id = ?", memberStatusReq.GetId()).
			MustSql()
		if _, err := tx.Exec(tx.Rebind(sql), args...); err != nil {
			return err
		}

		sql, args = sq.Select("d.id", "p.id").
			From("member m").
			LeftJoin("pool p ON m.pool_id = p.id").
			LeftJoin("domain_pool_relation dpr ON p.id = dpr.pool_id").
			LeftJoin("domain d ON dpr.domain_id = d.id").
			Where("m.id = ?", memberStatusReq.GetId()).
			MustSql()
		var domainID, poolID strfmt.UUID
		if err := tx.QueryRowx(tx.Rebind(sql), args...).Scan(&domainID, &poolID); err != nil {
			return err
		}

		sql, args = sq.Select("COUNT(m2.id)").
			From("member m").
			InnerJoin("member m2 ON m2.pool_id = m.pool_id").
			Where(sq.And{
				sq.Eq{"m.id": memberStatusReq.GetId()},
				sq.Expr("m2.id != m.id"),
				sq.Eq{"m2.status": "ONLINE"},
			}).MustSql()
		var membersOnline int
		if err := tx.Get(&membersOnline, tx.Rebind(sql), args...); err != nil {
			return err
		}

		if status == "ONLINE" || (membersOnline == 0 && status == "OFFLINE") {
			sql, args = sq.Update("pool").Set("status", status).Where("id = ?", poolID).MustSql()
			if _, err := tx.Exec(tx.Rebind(sql), args...); err != nil {
				return err
			}

			sql, args = sq.Update("domain").Set("status", status).Where("id = ?", domainID).MustSql()
			if _, err := tx.Exec(tx.Rebind(sql), args...); err != nil {
				return err
			}
		}

		statusResult = append(statusResult, &StatusResult{
			Id:      memberStatusReq.GetId(),
			Success: err == nil,
		})
	}
	res.MemberStatusResult = statusResult
	return tx.Commit()
}
