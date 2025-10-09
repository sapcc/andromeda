// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package f5

import (
	"encoding/json"
	"fmt"

	"github.com/apex/log"
	"github.com/f5devcentral/go-bigip"

	"github.com/sapcc/andromeda/internal/rpc/server"
	"github.com/sapcc/andromeda/internal/rpcmodels"
)

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
			StatusAvailabilityState struct {
				Description string `json:"description"`
			} `json:"status.availabilityState"`
		} `json:"entries"`
	} `json:"nestedStats"`
}

type ServerStats struct {
	NestedStats struct {
		Kind     string `json:"kind"`
		SelfLink string `json:"selfLink"`
		Entries  struct {
			VirtualServerPicks struct {
				Value uint64 `json:"value"`
			} `json:"vsPicks"`
		} `json:"entries"`
	} `json:"nestedStats"`
}

func buildMemberStatusUpdateRequest(session bigIPSession, store AndromedaF5Store) (*server.MemberStatusRequest, error) {
	datacenters, err := store.GetDatacenters()
	if err != nil {
		return nil, err
	}
	datacentersByID := make(map[string]*rpcmodels.Datacenter, len(datacenters))
	for _, dc := range datacenters {
		datacentersByID[dc.Id] = dc
	}
	domains, err := store.GetDomains()
	if err != nil {
		return nil, err
	}
	updates := []*server.MemberStatusRequest_MemberStatus{}
	for _, d := range domains {
		for _, p := range d.Pools {
			for _, m := range p.Members {
				if _, exists := datacentersByID[m.DatacenterId]; !exists {
					log.Warnf("invalid datacenter ID for member [datacenter ID = %s, member ID = %s]", m.DatacenterId, m.Id)
					continue
				}
				if datacentersByID[m.DatacenterId] == nil {
					log.Warnf("nil datacenter for member [member ID = %s]", m.Id)
					continue
				}
				urlPath := poolTypeAMemberStatsURL(
					as3DeclarationGSLBDomainTenantKey(d.Id),
					as3DeclarationGSLBPoolKey(p.Id),
					as3DeclarationGSLBServerKey(m.Address, datacentersByID[m.DatacenterId].Name),
					as3DeclarationGSLBVirtualServerName(m.Address, m.Port),
				)
				membersStats, err := fetchPoolTypeAMemberStats(session, urlPath)
				if err != nil {
					log.Warnf("failed to determine GSLB_Pool_Member_A status [BigIP URL path = %s]: %s", urlPath, err)
				}
				updates = append(updates, &server.MemberStatusRequest_MemberStatus{
					Id:     m.Id,
					Status: memberStatusFromPoolTypeAMemberAvailabilityState(membersStats.NestedStats.Entries.StatusAvailabilityState.Description),
				})
			}
		}
	}
	return &server.MemberStatusRequest{MemberStatus: updates}, nil
}

func poolTypeAMemberStatsURL(gslbDomainTenantKey, gslbPoolKey, gslbServerKey, gslbVirtualServerName string) string {
	return fmt.Sprintf("gtm/pool/a/~%s~application~%s/members/~Common~%s:%s/stats", gslbDomainTenantKey, gslbPoolKey, gslbServerKey, gslbVirtualServerName)
}

func serverStatsURL(gslbServerKey string) string {
	return fmt.Sprintf("gtm/server/~Common~%s/stats", gslbServerKey)
}

func memberStatusFromPoolTypeAMemberAvailabilityState(availabilityState string) server.MemberStatusRequest_MemberStatus_StatusType {
	switch availabilityState {
	case "available":
		return server.MemberStatusRequest_MemberStatus_ONLINE
	case "offline":
		return server.MemberStatusRequest_MemberStatus_OFFLINE
	default:
		return server.MemberStatusRequest_MemberStatus_UNKNOWN
	}
}

func fetchPoolTypeAMemberStats(s bigIPSession, urlPath string) (MembersStats, error) {
	var membersStats MembersStats
	mcs := MembersCollectionStats{}
	req := &bigip.APIRequest{
		Method:      "get",
		URL:         urlPath,
		ContentType: "application/json",
	}
	resp, err := s.APICall(req)
	if err != nil {
		var reqError bigip.RequestError
		_ = json.Unmarshal(resp, &reqError)
		if reqError.Code == 404 {
			return membersStats, fmt.Errorf("entity not found [BigIP URL Path = %s]: %w", urlPath, err)
		}
		return membersStats, err
	}
	err = json.Unmarshal(resp, &mcs)
	if err != nil {
		return membersStats, err
	}
	if len(mcs.Entries) != 1 {
		return membersStats, fmt.Errorf("expected exactly 1 key in `.entries`, got %d", len(mcs.Entries))
	}
	// stats.Entries, if valid, will always be a size 1 map and its only
	// key is always the pool type A member stats we need, so we iterate
	// just once in order to unmarshal its raw JSON value as a struct.
	for _, rawEntry := range mcs.Entries {
		if err := json.Unmarshal(rawEntry, &membersStats); err != nil {
			return membersStats, fmt.Errorf("could not decode nested member stats `entries.???.nestedStats`: %s", err)
		}
	}
	return membersStats, nil
}

func fetchServerStats(s bigIPSession, urlPath string) (ServerStats, error) {
	var serverStats ServerStats
	mcs := MembersCollectionStats{}
	req := &bigip.APIRequest{
		Method:      "get",
		URL:         urlPath,
		ContentType: "application/json",
	}
	resp, err := s.APICall(req)
	if err != nil {
		var reqError bigip.RequestError
		_ = json.Unmarshal(resp, &reqError)
		if reqError.Code == 404 {
			return serverStats, fmt.Errorf("entity not found [BigIP URL Path = %s]: %w", urlPath, err)
		}
		return serverStats, err
	}
	err = json.Unmarshal(resp, &mcs)
	if err != nil {
		return serverStats, err
	}
	if len(mcs.Entries) != 1 {
		return serverStats, fmt.Errorf("expected exactly 1 key in `.entries`, got %d", len(mcs.Entries))
	}
	// stats.Entries, if valid, will always be a size 1 map and its only
	// key is always the pool type A member stats we need, so we iterate
	// just once in order to unmarshal its raw JSON value as a struct.
	for _, rawEntry := range mcs.Entries {
		if err := json.Unmarshal(rawEntry, &serverStats); err != nil {
			return serverStats, fmt.Errorf("could not decode nested member stats `entries.???.nestedStats`: %s", err)
		}
	}
	return serverStats, nil
}
