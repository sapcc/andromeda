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

package f5

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/actatum/stormrpc"
	"github.com/apex/log"
	"github.com/nats-io/nats.go"
	"github.com/scottdware/go-bigip"

	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/rpc/server"
	"github.com/sapcc/andromeda/internal/utils"
)

type StatusController struct {
	bigIP *bigip.BigIP
	rpc   server.RPCServerClient
}

func ExecuteF5StatusAgent() error {
	session, err := GetBigIPSession()
	if err != nil {
		return fmt.Errorf("BigIP: %v", err)
	}

	device, err := session.GetCurrentDevice()
	if err != nil {
		return err
	}
	log.Infof("Connected to %s %s (%s)", device.MarketingName, device.Name, device.Version)

	// check if DNS module activated
	if err := BigIPSupportsDNS(device); err != nil {
		return err
	}

	// RPC server
	nc, err := nats.Connect(config.Global.Default.TransportURL)
	if err != nil {
		return err
	}
	client, err := stormrpc.NewClient("", stormrpc.WithNatsConn(nc))
	if err != nil {
		return err
	}

	sc := StatusController{
		session,
		server.NewRPCServerClient(client),
	}

	if config.Global.Default.Prometheus {
		go utils.PrometheusListen()
	}

	go func() {
		t := time.NewTicker(30 * time.Second)
		defer t.Stop()
		for {
			<-t.C // Activate periodically
			// Refresh token if needed
			if time.Until(sc.bigIP.TokenExpiry) <= 60 {
				if err := sc.bigIP.RefreshTokenSession(36000); err != nil {
					log.Error(err.Error())
				}
			}
			if err := sc.StatusHandler(); err != nil {
				log.Error(err.Error())
				_ = sc.bigIP.RefreshTokenSession(36000)
			}
		}
	}()

	return nil
}

type MembersCollectionStats struct {
	Kind     string                     `json:"kind"`
	SelfLink string                     `json:"selfLink"`
	Entries  map[string]json.RawMessage `json:"entries"`
}

type MembersStats struct {
	NestedStats struct {
		Kind     string `json:"kind"`
		SelfLink string `json:"selfLink"`
		Entries  struct {
			Alternate struct {
				Value int `json:"value"`
			} `json:"alternate"`
			Fallback struct {
				Value int `json:"value"`
			} `json:"fallback"`
			PoolName struct {
				Description string `json:"description"`
			} `json:"poolName"`
			PoolType struct {
				Description string `json:"description"`
			} `json:"poolType"`
			Preferred struct {
				Value int `json:"value"`
			} `json:"preferred"`
			ServerName struct {
				Description string `json:"description"`
			} `json:"serverName"`
			StatusAvailabilityState struct {
				Description string `json:"description"`
			} `json:"status.availabilityState"`
			StatusEnabledState struct {
				Description string `json:"description"`
			} `json:"status.enabledState"`
			StatusStatusReason struct {
				Description string `json:"description"`
			} `json:"status.statusReason"`
			VsName struct {
				Description string `json:"description"`
			} `json:"vsName"`
		} `json:"entries"`
	} `json:"nestedStats"`
}

func (c StatusController) StatusHandler() error {
	pools, err := c.bigIP.GetGTMAPools()
	if err != nil {
		return err
	}

	var stats MembersCollectionStats
	var msrs []*server.MemberStatusRequest_MemberStatus

	for _, pool := range pools.GTMAPools {
		partition := ConvertPartitionPath(pool.FullPath)
		err, ok := GetForEntity(c.bigIP, &stats, fmt.Sprintf("gtm/pool/a/%s/members/stats", partition))
		if err != nil {
			return err
		}
		if ok {
			for _, partitons := range stats.Entries {
				memberStats := MembersStats{}
				if err := json.Unmarshal(partitons, &memberStats); err != nil {
					log.Warnf("Could not decode nested member stats: %s", err)
					continue
				}

				memberID := strings.TrimPrefix(memberStats.NestedStats.Entries.VsName.Description, "member_")
				memberStatus := server.MemberStatusRequest_MemberStatus_UNKNOWN
				switch memberStats.NestedStats.Entries.StatusAvailabilityState.Description {
				case "available":
					memberStatus = server.MemberStatusRequest_MemberStatus_ONLINE
				case "offline":
					memberStatus = server.MemberStatusRequest_MemberStatus_OFFLINE
				}

				msr := &server.MemberStatusRequest_MemberStatus{
					Id:     memberID,
					Status: memberStatus,
				}
				log.Debugf("Member %s has status %s", memberID, msr.GetStatus().String())

				msrs = append(msrs, msr)
			}
		}
	}
	_, err = c.rpc.UpdateMemberStatus(context.Background(),
		&server.MemberStatusRequest{MemberStatus: msrs})
	if err != nil {
		return err
	}
	log.Infof("Refreshed status of %d members", len(msrs))

	return nil
}
