/*
 *   Copyright 2022 SAP SE
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

package akamai

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"go-micro.dev/v4/logger"

	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/driver"
	"github.com/sapcc/andromeda/internal/rpc/server"
	"github.com/sapcc/andromeda/internal/rpcmodels"
	"github.com/sapcc/andromeda/internal/utils"
)

type IPAvailability struct {
	DataRows []struct {
		CutOff      float64 `json:"cutOff"`
		Datacenters []struct {
			IPs []struct {
				Alive     bool    `json:"alive"`
				HandedOut bool    `json:"handedOut"`
				IP        string  `json:"ip"`
				Score     float64 `json:"score"`
			} `json:"IPs"`
			DatacenterID      int    `json:"datacenterId"`
			Nickname          string `json:"nickname"`
			TrafficTargetName string `json:"trafficTargetName"`
		} `json:"datacenters"`
		Timestamp time.Time `json:"timestamp"`
	}
}

func (s *AkamaiAgent) syncMemberStatus(domain *rpcmodels.Domain) error {
	trafficManagementDomain := config.Global.AkamaiConfig.Domain
	hostURL := fmt.Sprintf("/gtm-api/v1/reports/ip-availability/domains/%s/properties/%s?mostRecent=true",
		trafficManagementDomain, domain.GetFqdn())

	memberMap := make(map[string]string)
	for _, pool := range domain.Pools {
		for _, member := range pool.Members {
			memberMap[utils.InetNtoa(member.Address).String()] = member.GetId()
		}
	}

	stat := &IPAvailability{}
	req, _ := http.NewRequest(http.MethodGet, hostURL, nil)
	if out, err := (*s.session).Exec(req, &stat); err != nil {
		return err
	} else {
		bytes, _ := io.ReadAll(out.Body)
		logger.Debugf("%s", bytes)
	}

	if len(stat.DataRows) == 0 {
		// Nothing to do
		return nil
	}

	// we expect only one datarow due to mostRecent=true
	var memberStatusRequests []*server.MemberStatusRequest_MemberStatus
	for _, datacenter := range stat.DataRows[0].Datacenters {
		for _, ip := range datacenter.IPs {
			status := "OFFLINE"
			if ip.HandedOut {
				status = "ONLINE"
			}

			memberStatusRequests = append(memberStatusRequests,
				driver.GetMemberStatusRequest(memberMap[ip.IP], status))

			logger.Infof("status of %s: Alive: %+v ,HandedOut: %+v, %f", domain.Id, ip.Alive, ip.HandedOut, ip.Score)
		}
	}

	driver.UpdateMemberStatus(s.rpc, memberStatusRequests)
	return nil
}
