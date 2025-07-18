// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package akamai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/apex/log"

	"github.com/sapcc/andromeda/internal/rpcmodels"
)

type RPCHandlerAkamai struct {
	agent *AkamaiAgent
}

// Sync is a method that handles synchronization requests for Akamai domains.
func (h *RPCHandlerAkamai) Sync(ctx context.Context, request *rpcmodels.SyncRequest) (*rpcmodels.SyncResponse, error) {
	domainIDs := request.GetDomainId()
	log.WithField("domainIDs", domainIDs).Info("[Sync] Syncing domains")

	h.agent.forceSync <- domainIDs
	return &rpcmodels.SyncResponse{}, nil
}

// GetCidrs retrieves CIDR blocks from Akamai's Firewall Rules Manager API.
func (h *RPCHandlerAkamai) GetCidrs(ctx context.Context, request *rpcmodels.GetCidrsRequest) (*rpcmodels.GetCidrsResponse, error) {
	var cidrBlocks *rpcmodels.GetCidrsResponse

	cidrBlocksReq, _ := http.NewRequest(http.MethodGet, "/firewall-rules-manager/v1/cidr-blocks", nil)
	if ret, err := (*h.agent.session).Exec(cidrBlocksReq, cidrBlocks); err != nil {
		log.
			WithError(err).
			WithField("body", ret.Body).
			Error("Error fetching CIDR blocks")
		return nil, err
	}
	return cidrBlocks, nil
}

// DomainInfo represents domain information from the database
type DomainInfo struct {
	ID                 string `db:"id"`
	FQDN               string `db:"fqdn"`
	Provider           string `db:"provider"`
	ProjectID          string `db:"project_id"`
	ProvisioningStatus string `db:"provisioning_status"`
}
