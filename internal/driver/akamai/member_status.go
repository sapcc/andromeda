// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package akamai

import (
	"fmt"
	"net/http"
	"time"

	"github.com/apex/log"

	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/driver"
	"github.com/sapcc/andromeda/internal/rpc/server"
	"github.com/sapcc/andromeda/internal/rpcmodels"
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
	logger := log.WithField("domain.id", domain.Id)
	trafficManagementDomain := config.Global.AkamaiConfig.Domain
	hostURL := fmt.Sprintf("/gtm-api/v1/reports/ip-availability/domains/%s/properties/%s?mostRecent=true",
		trafficManagementDomain, domain.GetFqdn())

	memberMap := make(map[string]string)
	for _, pool := range domain.Pools {
		for _, member := range pool.Members {
			memberMap[member.Address] = member.GetId()
		}
	}

	if len(memberMap) == 0 {
		// Nothing to do
		return nil
	}

	stat := &IPAvailability{}
	req, _ := http.NewRequest(http.MethodGet, hostURL, nil)
	if _, err := (*s.session).Exec(req, &stat); err != nil {
		return err
	}

	if len(stat.DataRows) == 0 {
		logger.Warn("syncMemberStatus: no data rows found")
		return nil
	}

	// we expect only one datarow due to mostRecent=true
	var memberStatusRequests []*server.MemberStatusRequest_MemberStatus
	for _, datacenter := range stat.DataRows[0].Datacenters {
		for _, ip := range datacenter.IPs {
			status := server.MemberStatusRequest_MemberStatus_OFFLINE
			if ip.HandedOut && ip.Alive {
				status = server.MemberStatusRequest_MemberStatus_ONLINE
			} else if ip.HandedOut && !ip.Alive {
				status = server.MemberStatusRequest_MemberStatus_NO_MONITOR
			}

			if id, ok := memberMap[ip.IP]; ok {
				memberStatusRequests = append(memberStatusRequests,
					driver.GetMemberStatusRequest(id, status))
			} else {
				log.Warnf("unknown member with ip %s not found as port of domain %s", ip.IP, domain.Id)
			}

			log.Infof("status of domain %s: Alive: %+v, HandedOut: %+v, Score: %f -> Status: %s",
				domain.Id, ip.Alive, ip.HandedOut, ip.Score, status)
		}
	}

	driver.UpdateMemberStatus(s.rpc, memberStatusRequests)
	return nil
}
