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

func GetMemberStatusRequest(id string, status string) *server.MemberStatusRequest_MemberStatus {
	memberStatus := server.MemberStatusRequest_MemberStatus_StatusType(
		server.MemberStatusRequest_MemberStatus_StatusType_value[status])

	return &server.MemberStatusRequest_MemberStatus{
		Id:     id,
		Status: memberStatus,
	}
}

func UpdateMemberStatus(rpc server.RPCServerClient, statusRequests []*server.MemberStatusRequest_MemberStatus) {
	if _, err := rpc.UpdateMemberStatus(context.Background(),
		&server.MemberStatusRequest{MemberStatus: statusRequests}); err != nil {
		log.Error(err.Error())
	}
}
