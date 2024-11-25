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

package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/apex/log"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/utils/v2/openstack/clientconfig"
	"github.com/sapcc/go-bits/audittools"
	"github.com/sapcc/go-bits/gopherpolicy"

	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/policy"
	"github.com/sapcc/andromeda/internal/utils"
)

var (
	ErrForbidden      = errors.New("forbidden")
	projectContextKey = &contextKey{"os_token"}
)

type contextKey struct {
	name string
}

// Middleware Keystone token injector, also implements goslo policy checker
func KeystoneMiddleware(next http.Handler) (http.Handler, error) {
	authInfo := config.Global.ServiceAuth
	providerClient, err := clientconfig.AuthenticatedClient(context.Background(), &clientconfig.ClientOpts{
		AuthInfo: &authInfo})
	if err != nil {
		return nil, err
	}
	keystoneV3, err := openstack.NewIdentityV3(providerClient, gophercloud.EndpointOpts{})
	if err != nil {
		return nil, err
	}
	tv := gopherpolicy.TokenValidator{
		IdentityV3: keystoneV3,
		Cacher:     gopherpolicy.InMemoryCacher(),
	}
	if err := tv.LoadPolicyFile(config.Global.ApiSettings.PolicyFile, nil); err != nil {
		return nil, err
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := tv.CheckToken(r)
		if t.Err != nil {
			middleware.
				Error(403, utils.Unauthorized(t.Err), utils.JSONHeader).
				WriteResponse(w, runtime.JSONProducer())
			return
		}

		if t.Enforcer != nil && config.Global.ApiSettings.EnablePolicyTracing {
			t.Context.Logger = log.Debugf
			log.Debugf("token has auth = %v", t.Context.Auth)
			log.Debugf("token has roles = %v", t.Context.Roles)
		}

		ctx := context.WithValue(r.Context(), projectContextKey, t)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}), nil
}

func TokenFrom(r *http.Request) *gopherpolicy.Token {
	raw, ok := r.Context().Value(projectContextKey).(*gopherpolicy.Token)
	if ok {
		return raw
	}
	return nil
}

// ProjectScopeForRequest helper for getting project id
func ProjectScopeForRequest(r *http.Request) (string, error) {
	if config.Global.ApiSettings.AuthStrategy != "keystone" {
		return "", nil
	}
	ksToken := TokenFrom(r)
	if ksToken.Err != nil {
		return "", ksToken.Err
	}
	return ksToken.ProjectScopeUUID(), nil
}

func UserForRequest(r *http.Request) (audittools.UserInfo, error) {
	if config.Global.ApiSettings.AuthStrategy != "keystone" {
		return nil, nil
	}
	if ksToken := TokenFrom(r); ksToken != nil {
		return ksToken, nil
	}
	return nil, errors.New("failure accessing keystone token")
}

func Authenticate(r *http.Request, requestVars map[string]string) (string, error) {
	if config.Global.ApiSettings.AuthStrategy != "keystone" {
		return "", nil
	}

	if t := TokenFrom(r); t != nil {
		if requestVars == nil {
			// inject project_id from token if not provided
			requestVars = map[string]string{"project_id": t.ProjectScopeUUID()}
		}
		t.Context.Request = requestVars
		if t.Check(policy.RuleFromHTTPRequest(r)) {
			return t.ProjectScopeUUID(), nil
		}
		return "", ErrForbidden
	}

	return "", nil
}
