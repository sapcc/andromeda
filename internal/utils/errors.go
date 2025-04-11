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

package utils

import (
	"fmt"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/go-sql-driver/mysql"

	"github.com/sapcc/andromeda/models"
)

var (
	PolicyForbidden          = &models.Error{Code: 403, Message: "policy does not allow this request to be performed"}
	NotFound                 = &models.Error{Code: 404, Message: "not found"}
	ProviderUnchangeable     = &models.Error{Code: 400, Message: "provider cannot be changed"}
	InvalidMarker            = &models.Error{Code: 400, Message: "invalid marker"}
	PoolIDImmutable          = &models.Error{Code: 400, Message: "invalid value for 'pool_id': change of immutable attribute 'pool_id' not allowed"}
	PoolIDRequired           = &models.Error{Code: 400, Message: "invalid value for 'pool_id': 'pool_id' is required"}
	DuplicateMember          = &models.Error{Code: 409, Message: "invalid value for 'address' and 'port': endpoint already exists in pool"}
	DuplicateDomain          = &models.Error{Code: 409, Message: "invalid value for 'domain' and 'provider': domain already exists"}
	DatacenterInUse          = &models.Error{Code: 409, Message: "datacenter is in use"}
	InvalidSendString        = &models.Error{Code: 409, Message: "invalid value for 'send': must be a URL path"}
	MissingFQDN              = &models.Error{Code: 400, Message: "invalid value for 'fqdn': 'fqdn' is required"}
	MissingProvider          = &models.Error{Code: 400, Message: "invalid value for 'provider': 'provider' is required"}
	MissingAddressOrPort     = &models.Error{Code: 400, Message: "invalid value for 'address' and 'port': 'address' and 'port' are required"}
	FQDNImmutable            = &models.Error{Code: 400, Message: "invalid value for 'fqdn': change of immutable attribute 'fqdn' not allowed"}
	MySQLForeignKeyViolation = &mysql.MySQLError{Number: 1451}
)

func Unauthorized(err error) *models.Error {
	return &models.Error{Code: 401, Message: fmt.Sprintf("Unauthorized: %s", err.Error())}
}

func GetQuotaMetResponse(resource string) *models.Error {
	return &models.Error{Code: 403, Message: fmt.Sprintf(
		"Quota has been met for Resource: %s", resource)}
}

func GetInvalidProviderBoundResourceResponse(resource string) *models.Error {
	return &models.Error{Code: 403, Message: fmt.Sprintf(
		"Cannot determine quota: provider missing for resource: %s", resource)}
}

func GetErrorPoolNotFound(poolID *strfmt.UUID) *models.Error {
	return &models.Error{Code: 404, Message: fmt.Sprintf(
		"invalid value for 'pool_id': Pool '%s' not found", poolID)}
}

func GetErrorPoolHasAlreadyAMonitor(poolID *strfmt.UUID) *models.Error {
	return &models.Error{Code: 400, Message: fmt.Sprintf(
		"invalid value for 'pool_id': Pool '%s' already has a monitor", poolID)}
}

func GetErrorImmutable(self string, related string) *models.Error {
	return &models.Error{Code: 409, Message: fmt.Sprintf(
		"%s is currently immutable due to ongoing operations on related %s", self, related)}

}

type ResourcesNotFoundError struct {
	Ids      []strfmt.UUID
	Resource string
}

func (rnf ResourcesNotFoundError) Error() string {
	var err strings.Builder
	err.WriteString(rnf.Resource + "(s) not found: [")
	for i, id := range rnf.Ids {
		err.WriteString(id.String())
		if i < len(rnf.Ids)-1 {
			err.WriteString(", ")
		}
	}
	err.WriteString("]")
	return err.String()
}
