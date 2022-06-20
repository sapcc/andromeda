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
	"net"
	"strings"

	"github.com/jmoiron/sqlx"
	"go-micro.dev/v4/logger"

	"github.com/sapcc/andromeda/internal/models"
	"github.com/sapcc/andromeda/internal/utils"
)

type RPCHandler struct {
	DB *sqlx.DB
}

//QueryxWithIds run sql query if optional WHERE condition based on IDs
func (u *RPCHandler) QueryxWithIds(sql string, request *SearchRequest) (*sqlx.Rows, error) {
	var rows *sqlx.Rows
	if len(request.Ids) > 0 {
		if strings.Contains(sql, "WHERE") {
			sql += ` AND id IN (?)`
		} else {
			sql += ` WHERE id IN (?)`
		}
		query, args, err := sqlx.In(sql, request.Ids)
		if err != nil {
			return nil, err
		}
		query = u.DB.Rebind(query)
		rows, err = u.DB.Queryx(query, args...)
		if err != nil {
			return nil, err
		}
	} else {
		var err error
		rows, err = u.DB.Queryx(sql)
		if err != nil {
			return nil, err
		}
	}
	return rows, nil
}

func (u *RPCHandler) GetMembers(ctx context.Context, request *SearchRequest, response *MembersResponse) error {
	sql := `SELECT id, admin_state_up, address, port FROM member;`
	rows, err := u.QueryxWithIds(sql, request)
	if err != nil {
		return err
	}
	for rows.Next() {
		var member models.Member
		var address string
		if err := rows.Scan(&member.Id, &member.AdminStateUp, &address, &member.Port); err != nil {
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
	sql := `SELECT id, admin_state_up FROM pool;`
	rows, err := u.QueryxWithIds(sql, request)
	if err != nil {
		return err
	}
	for rows.Next() {
		var pool models.Pool
		if err := rows.Scan(&pool.Id, &pool.AdminStateUp); err != nil {
			return err
		}
		response.Response = append(response.Response, &pool)
	}
	return nil
}

func (u *RPCHandler) GetDatacenters(ctx context.Context, request *SearchRequest, response *DatacentersResponse) error {
	sql := `SELECT id, admin_state_up, city, state_or_province, continent, country, latitude, longitude, scope, provisioning_status FROM datacenter`
	rows, err := u.QueryxWithIds(sql, request)
	if err != nil {
		return err
	}
	for rows.Next() {
		var datacenter models.Datacenter
		if err := rows.StructScan(&datacenter); err != nil {
			logger.Error(err)
			return err
		}
		response.Response = append(response.Response, &datacenter)
	}
	return nil
}

func (u *RPCHandler) GetMonitors(ctx context.Context, request *SearchRequest, response *MonitorsResponse) error {
	sql := "SELECT id, admin_state_up, `interval`, pool_id, send, receive, timeout, type, provisioning_status FROM monitor;"
	rows, err := u.QueryxWithIds(sql, request)
	if err != nil {
		return err
	}
	for rows.Next() {
		var monitor models.Monitor
		if err := rows.StructScan(&monitor); err != nil {
			logger.Error(err)
			return err
		}
		response.Response = append(response.Response, &monitor)
	}
	return nil
}

func populatePools(u *RPCHandler, fullyPopulated bool, domainID string) ([]*models.Pool, error) {
	sql := `SELECT id, admin_state_up FROM pool p JOIN domain_pool_relation dpr on p.id = dpr.pool_id WHERE dpr.domain_id = ?`
	rows, err := u.DB.Queryx(sql, domainID)
	if err != nil {
		return nil, err
	}

	var pools []*models.Pool
	for rows.Next() {
		var pool models.Pool
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

func populateMonitors(u *RPCHandler, poolID string) ([]*models.Monitor, error) {
	sql := `SELECT id, admin_state_up, "interval", send, receive, timeout, type FROM monitor WHERE pool_id = ?`
	rows, err := u.DB.Queryx(sql, poolID)
	if err != nil {
		return nil, err
	}

	var monitors []*models.Monitor
	for rows.Next() {
		var monitor models.Monitor
		var send, receive *string
		var monitorType string
		if err := rows.Scan(&monitor.Id, &monitor.AdminStateUp, &monitor.Interval, &send, &receive, &monitor.Timeout, &monitorType); err != nil {
			return nil, err
		}
		if send != nil {
			monitor.Send = *send
		}
		if receive != nil {
			monitor.Receive = *receive
		}
		monitor.Type = models.Monitor_MonitorType(models.Monitor_MonitorType_value[monitorType])

		monitors = append(monitors, &monitor)
	}
	return monitors, err
}

func populateMembers(u *RPCHandler, poolID string) ([]*models.Member, error) {
	sql := `SELECT id, admin_state_up, address, port, datacenter_id FROM member WHERE pool_id = ?`
	rows, err := u.DB.Queryx(sql, poolID)
	if err != nil {
		return nil, err
	}

	var members []*models.Member
	for rows.Next() {
		var member models.Member
		var address string
		if err := rows.Scan(&member.Id, &member.AdminStateUp, &address,
			&member.Port, &member.Datacenter); err != nil {
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
	var sql string
	if request.Pending {
		sql = `
		SELECT 
			id, admin_state_up, fqdn, mode, record_type, provisioning_status
		FROM domain
		WHERE 
			provider = ? AND provisioning_status in ('PENDING_CREATE', 'PENDING_UPDATE', 'PENDING_DELETE')
		`
	} else {
		sql = `SELECT id, admin_state_up, fqdn, mode, record_type FROM domain WHERE provider = ?`
	}
	rows, err := u.QueryxWithIds(sql, request)

	if err != nil {
		logger.Error(err)
		return err
	}
	for rows.Next() {
		var domain models.Domain
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
				Pending:       true,
				Ids:           datacenterIds,
			}
			r := DatacentersResponse{}
			if err := u.GetDatacenters(nil, &s, &r); err != nil {
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
		sql := fmt.Sprintf(`UPDATE %s SET provisioning_status = ? WHERE id = ?`, table)
		_, err := u.DB.Exec(sql, provStatusReq.GetStatus().String(), provStatusReq.GetId())
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

func (u *RPCHandler) UpdateMemberStatus(ctx context.Context, req *MemberStatusRequest, res *MemberStatusResponse) error {
	var statusResult []*StatusResult
	for _, memberStatusReq := range req.GetMemberStatus() {
		sql := `UPDATE member SET status = ? WHERE id = ?`
		_, err := u.DB.Exec(sql, memberStatusReq.GetStatus().String(), memberStatusReq.GetId())
		if err != nil {
			logger.Error(err)
		}
		statusResult = append(statusResult, &StatusResult{
			Id:      memberStatusReq.GetId(),
			Success: err == nil,
		})
	}
	res.MemberStatusResult = statusResult
	return nil
}
