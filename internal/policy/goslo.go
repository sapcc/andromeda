/*
 *   Copyright 2020 SAP SE
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

package policy

import (
	"net/http"

	"github.com/asim/go-micro/v3/logger"

	"github.com/sapcc/andromeda/internal/auth"
	"github.com/sapcc/andromeda/internal/config"
)

type gosloPolicyEngine struct{}

func (p gosloPolicyEngine) init() {
	if config.Global.ApiSettings.AuthStrategy != "keystone" {
		logger.Fatal("Policy engine goslo supports only api_settings.auth_strategy = 'keystone'")
	}
}

func (p gosloPolicyEngine) AuthorizeRequest(r *http.Request, target string) bool {
	rule := RuleFromHTTPRequest(r)
	t := auth.TokenFrom(r)

	if t != nil {
		t.Context.Request = map[string]string{
			"project_id": target,
		}
		return t.Check(rule)
	}
	// Ignore disabled keystone middleware
	return true
}

func (p gosloPolicyEngine) AuthorizeGetAllRequest(r *http.Request, target string) bool {
	rule := RuleFromHTTPRequest(r)
	t := auth.TokenFrom(r)

	if t != nil {
		t.Context.Request = map[string]string{
			"project_id": target,
		}
		return t.Check(rule + "-global")
	}
	// Ignore disabled keystone middleware
	return true
}
