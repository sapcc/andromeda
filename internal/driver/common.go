// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"context"

	"github.com/apex/log"

	"github.com/sapcc/andromeda/internal/rpc/server"
)

func GetProvisioningStatusRequest(id string, model string, status string) *server.ProvisioningStatusRequest_ProvisioningStatus {
	provisioningStatus := server.ProvisioningStatusRequest_ProvisioningStatus_StatusType(
		server.ProvisioningStatusRequest_ProvisioningStatus_StatusType_value[status])
	provisioningModel := server.ProvisioningStatusRequest_ProvisioningStatus_Model(
		server.ProvisioningStatusRequest_ProvisioningStatus_Model_value[model])

	return &server.ProvisioningStatusRequest_ProvisioningStatus{
		Id:     id,
		Model:  provisioningModel,
		Status: provisioningStatus,
	}
}

func UpdateProvisioningStatus(rpc server.RPCServerClient, statusRequests []*server.ProvisioningStatusRequest_ProvisioningStatus) {
	if _, err := rpc.UpdateProvisioningStatus(context.Background(),
		&server.ProvisioningStatusRequest{ProvisioningStatus: statusRequests}); err != nil {
		log.Error(err.Error())
	}
}

func GetMemberStatusRequest(id string, status server.MemberStatusRequest_MemberStatus_StatusType) *server.MemberStatusRequest_MemberStatus {
	return &server.MemberStatusRequest_MemberStatus{
		Id:     id,
		Status: status,
	}
}

func UpdateMemberStatus(rpc server.RPCServerClient, statusRequests []*server.MemberStatusRequest_MemberStatus) {
	if _, err := rpc.UpdateMemberStatus(context.Background(),
		&server.MemberStatusRequest{MemberStatus: statusRequests}); err != nil {
		log.Error(err.Error())
	}
}

const (
	AgentAkamai = "andromeda-akamai-agent"
	AgentF5     = "andromeda-f5-agent"
	Server      = "andromeda-server"
)

func GetServiceType(name string) string {
	switch name {
	case AgentAkamai:
		return "akamai"
	case AgentF5:
		return "f5"
	case Server:
		return "server"
	default:
		return "unknown"
	}
}
