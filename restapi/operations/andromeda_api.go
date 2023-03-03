// Code generated by go-swagger; DO NOT EDIT.

// Copyright 2020 SAP SE
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/runtime/security"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/sapcc/andromeda/restapi/operations/administrative"
	"github.com/sapcc/andromeda/restapi/operations/datacenters"
	"github.com/sapcc/andromeda/restapi/operations/domains"
	"github.com/sapcc/andromeda/restapi/operations/members"
	"github.com/sapcc/andromeda/restapi/operations/monitors"
	"github.com/sapcc/andromeda/restapi/operations/pools"
)

// NewAndromedaAPI creates a new Andromeda instance
func NewAndromedaAPI(spec *loads.Document) *AndromedaAPI {
	return &AndromedaAPI{
		handlers:            make(map[string]map[string]http.Handler),
		formats:             strfmt.Default,
		defaultConsumes:     "application/json",
		defaultProduces:     "application/json",
		customConsumers:     make(map[string]runtime.Consumer),
		customProducers:     make(map[string]runtime.Producer),
		PreServerShutdown:   func() {},
		ServerShutdown:      func() {},
		spec:                spec,
		useSwaggerUI:        false,
		ServeError:          errors.ServeError,
		BasicAuthenticator:  security.BasicAuth,
		APIKeyAuthenticator: security.APIKeyAuth,
		BearerAuthenticator: security.BearerAuth,

		JSONConsumer: runtime.JSONConsumer(),

		JSONProducer: runtime.JSONProducer(),

		DatacentersDeleteDatacentersDatacenterIDHandler: datacenters.DeleteDatacentersDatacenterIDHandlerFunc(func(params datacenters.DeleteDatacentersDatacenterIDParams) middleware.Responder {
			return middleware.NotImplemented("operation datacenters.DeleteDatacentersDatacenterID has not yet been implemented")
		}),
		DomainsDeleteDomainsDomainIDHandler: domains.DeleteDomainsDomainIDHandlerFunc(func(params domains.DeleteDomainsDomainIDParams) middleware.Responder {
			return middleware.NotImplemented("operation domains.DeleteDomainsDomainID has not yet been implemented")
		}),
		MembersDeleteMembersMemberIDHandler: members.DeleteMembersMemberIDHandlerFunc(func(params members.DeleteMembersMemberIDParams) middleware.Responder {
			return middleware.NotImplemented("operation members.DeleteMembersMemberID has not yet been implemented")
		}),
		MonitorsDeleteMonitorsMonitorIDHandler: monitors.DeleteMonitorsMonitorIDHandlerFunc(func(params monitors.DeleteMonitorsMonitorIDParams) middleware.Responder {
			return middleware.NotImplemented("operation monitors.DeleteMonitorsMonitorID has not yet been implemented")
		}),
		PoolsDeletePoolsPoolIDHandler: pools.DeletePoolsPoolIDHandlerFunc(func(params pools.DeletePoolsPoolIDParams) middleware.Responder {
			return middleware.NotImplemented("operation pools.DeletePoolsPoolID has not yet been implemented")
		}),
		AdministrativeDeleteQuotasProjectIDHandler: administrative.DeleteQuotasProjectIDHandlerFunc(func(params administrative.DeleteQuotasProjectIDParams) middleware.Responder {
			return middleware.NotImplemented("operation administrative.DeleteQuotasProjectID has not yet been implemented")
		}),
		DatacentersGetDatacentersHandler: datacenters.GetDatacentersHandlerFunc(func(params datacenters.GetDatacentersParams) middleware.Responder {
			return middleware.NotImplemented("operation datacenters.GetDatacenters has not yet been implemented")
		}),
		DatacentersGetDatacentersDatacenterIDHandler: datacenters.GetDatacentersDatacenterIDHandlerFunc(func(params datacenters.GetDatacentersDatacenterIDParams) middleware.Responder {
			return middleware.NotImplemented("operation datacenters.GetDatacentersDatacenterID has not yet been implemented")
		}),
		DomainsGetDomainsHandler: domains.GetDomainsHandlerFunc(func(params domains.GetDomainsParams) middleware.Responder {
			return middleware.NotImplemented("operation domains.GetDomains has not yet been implemented")
		}),
		DomainsGetDomainsDomainIDHandler: domains.GetDomainsDomainIDHandlerFunc(func(params domains.GetDomainsDomainIDParams) middleware.Responder {
			return middleware.NotImplemented("operation domains.GetDomainsDomainID has not yet been implemented")
		}),
		MembersGetMembersHandler: members.GetMembersHandlerFunc(func(params members.GetMembersParams) middleware.Responder {
			return middleware.NotImplemented("operation members.GetMembers has not yet been implemented")
		}),
		MembersGetMembersMemberIDHandler: members.GetMembersMemberIDHandlerFunc(func(params members.GetMembersMemberIDParams) middleware.Responder {
			return middleware.NotImplemented("operation members.GetMembersMemberID has not yet been implemented")
		}),
		MonitorsGetMonitorsHandler: monitors.GetMonitorsHandlerFunc(func(params monitors.GetMonitorsParams) middleware.Responder {
			return middleware.NotImplemented("operation monitors.GetMonitors has not yet been implemented")
		}),
		MonitorsGetMonitorsMonitorIDHandler: monitors.GetMonitorsMonitorIDHandlerFunc(func(params monitors.GetMonitorsMonitorIDParams) middleware.Responder {
			return middleware.NotImplemented("operation monitors.GetMonitorsMonitorID has not yet been implemented")
		}),
		PoolsGetPoolsHandler: pools.GetPoolsHandlerFunc(func(params pools.GetPoolsParams) middleware.Responder {
			return middleware.NotImplemented("operation pools.GetPools has not yet been implemented")
		}),
		PoolsGetPoolsPoolIDHandler: pools.GetPoolsPoolIDHandlerFunc(func(params pools.GetPoolsPoolIDParams) middleware.Responder {
			return middleware.NotImplemented("operation pools.GetPoolsPoolID has not yet been implemented")
		}),
		AdministrativeGetQuotasHandler: administrative.GetQuotasHandlerFunc(func(params administrative.GetQuotasParams) middleware.Responder {
			return middleware.NotImplemented("operation administrative.GetQuotas has not yet been implemented")
		}),
		AdministrativeGetQuotasDefaultsHandler: administrative.GetQuotasDefaultsHandlerFunc(func(params administrative.GetQuotasDefaultsParams) middleware.Responder {
			return middleware.NotImplemented("operation administrative.GetQuotasDefaults has not yet been implemented")
		}),
		AdministrativeGetQuotasProjectIDHandler: administrative.GetQuotasProjectIDHandlerFunc(func(params administrative.GetQuotasProjectIDParams) middleware.Responder {
			return middleware.NotImplemented("operation administrative.GetQuotasProjectID has not yet been implemented")
		}),
		AdministrativeGetServicesHandler: administrative.GetServicesHandlerFunc(func(params administrative.GetServicesParams) middleware.Responder {
			return middleware.NotImplemented("operation administrative.GetServices has not yet been implemented")
		}),
		DatacentersPostDatacentersHandler: datacenters.PostDatacentersHandlerFunc(func(params datacenters.PostDatacentersParams) middleware.Responder {
			return middleware.NotImplemented("operation datacenters.PostDatacenters has not yet been implemented")
		}),
		DomainsPostDomainsHandler: domains.PostDomainsHandlerFunc(func(params domains.PostDomainsParams) middleware.Responder {
			return middleware.NotImplemented("operation domains.PostDomains has not yet been implemented")
		}),
		MembersPostMembersHandler: members.PostMembersHandlerFunc(func(params members.PostMembersParams) middleware.Responder {
			return middleware.NotImplemented("operation members.PostMembers has not yet been implemented")
		}),
		MonitorsPostMonitorsHandler: monitors.PostMonitorsHandlerFunc(func(params monitors.PostMonitorsParams) middleware.Responder {
			return middleware.NotImplemented("operation monitors.PostMonitors has not yet been implemented")
		}),
		PoolsPostPoolsHandler: pools.PostPoolsHandlerFunc(func(params pools.PostPoolsParams) middleware.Responder {
			return middleware.NotImplemented("operation pools.PostPools has not yet been implemented")
		}),
		AdministrativePostSyncHandler: administrative.PostSyncHandlerFunc(func(params administrative.PostSyncParams) middleware.Responder {
			return middleware.NotImplemented("operation administrative.PostSync has not yet been implemented")
		}),
		DatacentersPutDatacentersDatacenterIDHandler: datacenters.PutDatacentersDatacenterIDHandlerFunc(func(params datacenters.PutDatacentersDatacenterIDParams) middleware.Responder {
			return middleware.NotImplemented("operation datacenters.PutDatacentersDatacenterID has not yet been implemented")
		}),
		DomainsPutDomainsDomainIDHandler: domains.PutDomainsDomainIDHandlerFunc(func(params domains.PutDomainsDomainIDParams) middleware.Responder {
			return middleware.NotImplemented("operation domains.PutDomainsDomainID has not yet been implemented")
		}),
		MembersPutMembersMemberIDHandler: members.PutMembersMemberIDHandlerFunc(func(params members.PutMembersMemberIDParams) middleware.Responder {
			return middleware.NotImplemented("operation members.PutMembersMemberID has not yet been implemented")
		}),
		MonitorsPutMonitorsMonitorIDHandler: monitors.PutMonitorsMonitorIDHandlerFunc(func(params monitors.PutMonitorsMonitorIDParams) middleware.Responder {
			return middleware.NotImplemented("operation monitors.PutMonitorsMonitorID has not yet been implemented")
		}),
		PoolsPutPoolsPoolIDHandler: pools.PutPoolsPoolIDHandlerFunc(func(params pools.PutPoolsPoolIDParams) middleware.Responder {
			return middleware.NotImplemented("operation pools.PutPoolsPoolID has not yet been implemented")
		}),
		AdministrativePutQuotasProjectIDHandler: administrative.PutQuotasProjectIDHandlerFunc(func(params administrative.PutQuotasProjectIDParams) middleware.Responder {
			return middleware.NotImplemented("operation administrative.PutQuotasProjectID has not yet been implemented")
		}),
	}
}

/*AndromedaAPI Platform agnostic GSLB frontend */
type AndromedaAPI struct {
	spec            *loads.Document
	context         *middleware.Context
	handlers        map[string]map[string]http.Handler
	formats         strfmt.Registry
	customConsumers map[string]runtime.Consumer
	customProducers map[string]runtime.Producer
	defaultConsumes string
	defaultProduces string
	Middleware      func(middleware.Builder) http.Handler
	useSwaggerUI    bool

	// BasicAuthenticator generates a runtime.Authenticator from the supplied basic auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	BasicAuthenticator func(security.UserPassAuthentication) runtime.Authenticator

	// APIKeyAuthenticator generates a runtime.Authenticator from the supplied token auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	APIKeyAuthenticator func(string, string, security.TokenAuthentication) runtime.Authenticator

	// BearerAuthenticator generates a runtime.Authenticator from the supplied bearer token auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	BearerAuthenticator func(string, security.ScopedTokenAuthentication) runtime.Authenticator

	// JSONConsumer registers a consumer for the following mime types:
	//   - application/json
	JSONConsumer runtime.Consumer

	// JSONProducer registers a producer for the following mime types:
	//   - application/json
	JSONProducer runtime.Producer

	// DatacentersDeleteDatacentersDatacenterIDHandler sets the operation handler for the delete datacenters datacenter ID operation
	DatacentersDeleteDatacentersDatacenterIDHandler datacenters.DeleteDatacentersDatacenterIDHandler
	// DomainsDeleteDomainsDomainIDHandler sets the operation handler for the delete domains domain ID operation
	DomainsDeleteDomainsDomainIDHandler domains.DeleteDomainsDomainIDHandler
	// MembersDeleteMembersMemberIDHandler sets the operation handler for the delete members member ID operation
	MembersDeleteMembersMemberIDHandler members.DeleteMembersMemberIDHandler
	// MonitorsDeleteMonitorsMonitorIDHandler sets the operation handler for the delete monitors monitor ID operation
	MonitorsDeleteMonitorsMonitorIDHandler monitors.DeleteMonitorsMonitorIDHandler
	// PoolsDeletePoolsPoolIDHandler sets the operation handler for the delete pools pool ID operation
	PoolsDeletePoolsPoolIDHandler pools.DeletePoolsPoolIDHandler
	// AdministrativeDeleteQuotasProjectIDHandler sets the operation handler for the delete quotas project ID operation
	AdministrativeDeleteQuotasProjectIDHandler administrative.DeleteQuotasProjectIDHandler
	// DatacentersGetDatacentersHandler sets the operation handler for the get datacenters operation
	DatacentersGetDatacentersHandler datacenters.GetDatacentersHandler
	// DatacentersGetDatacentersDatacenterIDHandler sets the operation handler for the get datacenters datacenter ID operation
	DatacentersGetDatacentersDatacenterIDHandler datacenters.GetDatacentersDatacenterIDHandler
	// DomainsGetDomainsHandler sets the operation handler for the get domains operation
	DomainsGetDomainsHandler domains.GetDomainsHandler
	// DomainsGetDomainsDomainIDHandler sets the operation handler for the get domains domain ID operation
	DomainsGetDomainsDomainIDHandler domains.GetDomainsDomainIDHandler
	// MembersGetMembersHandler sets the operation handler for the get members operation
	MembersGetMembersHandler members.GetMembersHandler
	// MembersGetMembersMemberIDHandler sets the operation handler for the get members member ID operation
	MembersGetMembersMemberIDHandler members.GetMembersMemberIDHandler
	// MonitorsGetMonitorsHandler sets the operation handler for the get monitors operation
	MonitorsGetMonitorsHandler monitors.GetMonitorsHandler
	// MonitorsGetMonitorsMonitorIDHandler sets the operation handler for the get monitors monitor ID operation
	MonitorsGetMonitorsMonitorIDHandler monitors.GetMonitorsMonitorIDHandler
	// PoolsGetPoolsHandler sets the operation handler for the get pools operation
	PoolsGetPoolsHandler pools.GetPoolsHandler
	// PoolsGetPoolsPoolIDHandler sets the operation handler for the get pools pool ID operation
	PoolsGetPoolsPoolIDHandler pools.GetPoolsPoolIDHandler
	// AdministrativeGetQuotasHandler sets the operation handler for the get quotas operation
	AdministrativeGetQuotasHandler administrative.GetQuotasHandler
	// AdministrativeGetQuotasDefaultsHandler sets the operation handler for the get quotas defaults operation
	AdministrativeGetQuotasDefaultsHandler administrative.GetQuotasDefaultsHandler
	// AdministrativeGetQuotasProjectIDHandler sets the operation handler for the get quotas project ID operation
	AdministrativeGetQuotasProjectIDHandler administrative.GetQuotasProjectIDHandler
	// AdministrativeGetServicesHandler sets the operation handler for the get services operation
	AdministrativeGetServicesHandler administrative.GetServicesHandler
	// DatacentersPostDatacentersHandler sets the operation handler for the post datacenters operation
	DatacentersPostDatacentersHandler datacenters.PostDatacentersHandler
	// DomainsPostDomainsHandler sets the operation handler for the post domains operation
	DomainsPostDomainsHandler domains.PostDomainsHandler
	// MembersPostMembersHandler sets the operation handler for the post members operation
	MembersPostMembersHandler members.PostMembersHandler
	// MonitorsPostMonitorsHandler sets the operation handler for the post monitors operation
	MonitorsPostMonitorsHandler monitors.PostMonitorsHandler
	// PoolsPostPoolsHandler sets the operation handler for the post pools operation
	PoolsPostPoolsHandler pools.PostPoolsHandler
	// AdministrativePostSyncHandler sets the operation handler for the post sync operation
	AdministrativePostSyncHandler administrative.PostSyncHandler
	// DatacentersPutDatacentersDatacenterIDHandler sets the operation handler for the put datacenters datacenter ID operation
	DatacentersPutDatacentersDatacenterIDHandler datacenters.PutDatacentersDatacenterIDHandler
	// DomainsPutDomainsDomainIDHandler sets the operation handler for the put domains domain ID operation
	DomainsPutDomainsDomainIDHandler domains.PutDomainsDomainIDHandler
	// MembersPutMembersMemberIDHandler sets the operation handler for the put members member ID operation
	MembersPutMembersMemberIDHandler members.PutMembersMemberIDHandler
	// MonitorsPutMonitorsMonitorIDHandler sets the operation handler for the put monitors monitor ID operation
	MonitorsPutMonitorsMonitorIDHandler monitors.PutMonitorsMonitorIDHandler
	// PoolsPutPoolsPoolIDHandler sets the operation handler for the put pools pool ID operation
	PoolsPutPoolsPoolIDHandler pools.PutPoolsPoolIDHandler
	// AdministrativePutQuotasProjectIDHandler sets the operation handler for the put quotas project ID operation
	AdministrativePutQuotasProjectIDHandler administrative.PutQuotasProjectIDHandler

	// ServeError is called when an error is received, there is a default handler
	// but you can set your own with this
	ServeError func(http.ResponseWriter, *http.Request, error)

	// PreServerShutdown is called before the HTTP(S) server is shutdown
	// This allows for custom functions to get executed before the HTTP(S) server stops accepting traffic
	PreServerShutdown func()

	// ServerShutdown is called when the HTTP(S) server is shut down and done
	// handling all active connections and does not accept connections any more
	ServerShutdown func()

	// Custom command line argument groups with their descriptions
	CommandLineOptionsGroups []swag.CommandLineOptionsGroup

	// User defined logger function.
	Logger func(string, ...interface{})
}

// UseRedoc for documentation at /docs
func (o *AndromedaAPI) UseRedoc() {
	o.useSwaggerUI = false
}

// UseSwaggerUI for documentation at /docs
func (o *AndromedaAPI) UseSwaggerUI() {
	o.useSwaggerUI = true
}

// SetDefaultProduces sets the default produces media type
func (o *AndromedaAPI) SetDefaultProduces(mediaType string) {
	o.defaultProduces = mediaType
}

// SetDefaultConsumes returns the default consumes media type
func (o *AndromedaAPI) SetDefaultConsumes(mediaType string) {
	o.defaultConsumes = mediaType
}

// SetSpec sets a spec that will be served for the clients.
func (o *AndromedaAPI) SetSpec(spec *loads.Document) {
	o.spec = spec
}

// DefaultProduces returns the default produces media type
func (o *AndromedaAPI) DefaultProduces() string {
	return o.defaultProduces
}

// DefaultConsumes returns the default consumes media type
func (o *AndromedaAPI) DefaultConsumes() string {
	return o.defaultConsumes
}

// Formats returns the registered string formats
func (o *AndromedaAPI) Formats() strfmt.Registry {
	return o.formats
}

// RegisterFormat registers a custom format validator
func (o *AndromedaAPI) RegisterFormat(name string, format strfmt.Format, validator strfmt.Validator) {
	o.formats.Add(name, format, validator)
}

// Validate validates the registrations in the AndromedaAPI
func (o *AndromedaAPI) Validate() error {
	var unregistered []string

	if o.JSONConsumer == nil {
		unregistered = append(unregistered, "JSONConsumer")
	}

	if o.JSONProducer == nil {
		unregistered = append(unregistered, "JSONProducer")
	}

	if o.DatacentersDeleteDatacentersDatacenterIDHandler == nil {
		unregistered = append(unregistered, "datacenters.DeleteDatacentersDatacenterIDHandler")
	}
	if o.DomainsDeleteDomainsDomainIDHandler == nil {
		unregistered = append(unregistered, "domains.DeleteDomainsDomainIDHandler")
	}
	if o.MembersDeleteMembersMemberIDHandler == nil {
		unregistered = append(unregistered, "members.DeleteMembersMemberIDHandler")
	}
	if o.MonitorsDeleteMonitorsMonitorIDHandler == nil {
		unregistered = append(unregistered, "monitors.DeleteMonitorsMonitorIDHandler")
	}
	if o.PoolsDeletePoolsPoolIDHandler == nil {
		unregistered = append(unregistered, "pools.DeletePoolsPoolIDHandler")
	}
	if o.AdministrativeDeleteQuotasProjectIDHandler == nil {
		unregistered = append(unregistered, "administrative.DeleteQuotasProjectIDHandler")
	}
	if o.DatacentersGetDatacentersHandler == nil {
		unregistered = append(unregistered, "datacenters.GetDatacentersHandler")
	}
	if o.DatacentersGetDatacentersDatacenterIDHandler == nil {
		unregistered = append(unregistered, "datacenters.GetDatacentersDatacenterIDHandler")
	}
	if o.DomainsGetDomainsHandler == nil {
		unregistered = append(unregistered, "domains.GetDomainsHandler")
	}
	if o.DomainsGetDomainsDomainIDHandler == nil {
		unregistered = append(unregistered, "domains.GetDomainsDomainIDHandler")
	}
	if o.MembersGetMembersHandler == nil {
		unregistered = append(unregistered, "members.GetMembersHandler")
	}
	if o.MembersGetMembersMemberIDHandler == nil {
		unregistered = append(unregistered, "members.GetMembersMemberIDHandler")
	}
	if o.MonitorsGetMonitorsHandler == nil {
		unregistered = append(unregistered, "monitors.GetMonitorsHandler")
	}
	if o.MonitorsGetMonitorsMonitorIDHandler == nil {
		unregistered = append(unregistered, "monitors.GetMonitorsMonitorIDHandler")
	}
	if o.PoolsGetPoolsHandler == nil {
		unregistered = append(unregistered, "pools.GetPoolsHandler")
	}
	if o.PoolsGetPoolsPoolIDHandler == nil {
		unregistered = append(unregistered, "pools.GetPoolsPoolIDHandler")
	}
	if o.AdministrativeGetQuotasHandler == nil {
		unregistered = append(unregistered, "administrative.GetQuotasHandler")
	}
	if o.AdministrativeGetQuotasDefaultsHandler == nil {
		unregistered = append(unregistered, "administrative.GetQuotasDefaultsHandler")
	}
	if o.AdministrativeGetQuotasProjectIDHandler == nil {
		unregistered = append(unregistered, "administrative.GetQuotasProjectIDHandler")
	}
	if o.AdministrativeGetServicesHandler == nil {
		unregistered = append(unregistered, "administrative.GetServicesHandler")
	}
	if o.DatacentersPostDatacentersHandler == nil {
		unregistered = append(unregistered, "datacenters.PostDatacentersHandler")
	}
	if o.DomainsPostDomainsHandler == nil {
		unregistered = append(unregistered, "domains.PostDomainsHandler")
	}
	if o.MembersPostMembersHandler == nil {
		unregistered = append(unregistered, "members.PostMembersHandler")
	}
	if o.MonitorsPostMonitorsHandler == nil {
		unregistered = append(unregistered, "monitors.PostMonitorsHandler")
	}
	if o.PoolsPostPoolsHandler == nil {
		unregistered = append(unregistered, "pools.PostPoolsHandler")
	}
	if o.AdministrativePostSyncHandler == nil {
		unregistered = append(unregistered, "administrative.PostSyncHandler")
	}
	if o.DatacentersPutDatacentersDatacenterIDHandler == nil {
		unregistered = append(unregistered, "datacenters.PutDatacentersDatacenterIDHandler")
	}
	if o.DomainsPutDomainsDomainIDHandler == nil {
		unregistered = append(unregistered, "domains.PutDomainsDomainIDHandler")
	}
	if o.MembersPutMembersMemberIDHandler == nil {
		unregistered = append(unregistered, "members.PutMembersMemberIDHandler")
	}
	if o.MonitorsPutMonitorsMonitorIDHandler == nil {
		unregistered = append(unregistered, "monitors.PutMonitorsMonitorIDHandler")
	}
	if o.PoolsPutPoolsPoolIDHandler == nil {
		unregistered = append(unregistered, "pools.PutPoolsPoolIDHandler")
	}
	if o.AdministrativePutQuotasProjectIDHandler == nil {
		unregistered = append(unregistered, "administrative.PutQuotasProjectIDHandler")
	}

	if len(unregistered) > 0 {
		return fmt.Errorf("missing registration: %s", strings.Join(unregistered, ", "))
	}

	return nil
}

// ServeErrorFor gets a error handler for a given operation id
func (o *AndromedaAPI) ServeErrorFor(operationID string) func(http.ResponseWriter, *http.Request, error) {
	return o.ServeError
}

// AuthenticatorsFor gets the authenticators for the specified security schemes
func (o *AndromedaAPI) AuthenticatorsFor(schemes map[string]spec.SecurityScheme) map[string]runtime.Authenticator {
	return nil
}

// Authorizer returns the registered authorizer
func (o *AndromedaAPI) Authorizer() runtime.Authorizer {
	return nil
}

// ConsumersFor gets the consumers for the specified media types.
// MIME type parameters are ignored here.
func (o *AndromedaAPI) ConsumersFor(mediaTypes []string) map[string]runtime.Consumer {
	result := make(map[string]runtime.Consumer, len(mediaTypes))
	for _, mt := range mediaTypes {
		switch mt {
		case "application/json":
			result["application/json"] = o.JSONConsumer
		}

		if c, ok := o.customConsumers[mt]; ok {
			result[mt] = c
		}
	}
	return result
}

// ProducersFor gets the producers for the specified media types.
// MIME type parameters are ignored here.
func (o *AndromedaAPI) ProducersFor(mediaTypes []string) map[string]runtime.Producer {
	result := make(map[string]runtime.Producer, len(mediaTypes))
	for _, mt := range mediaTypes {
		switch mt {
		case "application/json":
			result["application/json"] = o.JSONProducer
		}

		if p, ok := o.customProducers[mt]; ok {
			result[mt] = p
		}
	}
	return result
}

// HandlerFor gets a http.Handler for the provided operation method and path
func (o *AndromedaAPI) HandlerFor(method, path string) (http.Handler, bool) {
	if o.handlers == nil {
		return nil, false
	}
	um := strings.ToUpper(method)
	if _, ok := o.handlers[um]; !ok {
		return nil, false
	}
	if path == "/" {
		path = ""
	}
	h, ok := o.handlers[um][path]
	return h, ok
}

// Context returns the middleware context for the andromeda API
func (o *AndromedaAPI) Context() *middleware.Context {
	if o.context == nil {
		o.context = middleware.NewRoutableContext(o.spec, o, nil)
	}

	return o.context
}

func (o *AndromedaAPI) initHandlerCache() {
	o.Context() // don't care about the result, just that the initialization happened
	if o.handlers == nil {
		o.handlers = make(map[string]map[string]http.Handler)
	}

	if o.handlers["DELETE"] == nil {
		o.handlers["DELETE"] = make(map[string]http.Handler)
	}
	o.handlers["DELETE"]["/datacenters/{datacenter_id}"] = datacenters.NewDeleteDatacentersDatacenterID(o.context, o.DatacentersDeleteDatacentersDatacenterIDHandler)
	if o.handlers["DELETE"] == nil {
		o.handlers["DELETE"] = make(map[string]http.Handler)
	}
	o.handlers["DELETE"]["/domains/{domain_id}"] = domains.NewDeleteDomainsDomainID(o.context, o.DomainsDeleteDomainsDomainIDHandler)
	if o.handlers["DELETE"] == nil {
		o.handlers["DELETE"] = make(map[string]http.Handler)
	}
	o.handlers["DELETE"]["/members/{member_id}"] = members.NewDeleteMembersMemberID(o.context, o.MembersDeleteMembersMemberIDHandler)
	if o.handlers["DELETE"] == nil {
		o.handlers["DELETE"] = make(map[string]http.Handler)
	}
	o.handlers["DELETE"]["/monitors/{monitor_id}"] = monitors.NewDeleteMonitorsMonitorID(o.context, o.MonitorsDeleteMonitorsMonitorIDHandler)
	if o.handlers["DELETE"] == nil {
		o.handlers["DELETE"] = make(map[string]http.Handler)
	}
	o.handlers["DELETE"]["/pools/{pool_id}"] = pools.NewDeletePoolsPoolID(o.context, o.PoolsDeletePoolsPoolIDHandler)
	if o.handlers["DELETE"] == nil {
		o.handlers["DELETE"] = make(map[string]http.Handler)
	}
	o.handlers["DELETE"]["/quotas/{project_id}"] = administrative.NewDeleteQuotasProjectID(o.context, o.AdministrativeDeleteQuotasProjectIDHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/datacenters"] = datacenters.NewGetDatacenters(o.context, o.DatacentersGetDatacentersHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/datacenters/{datacenter_id}"] = datacenters.NewGetDatacentersDatacenterID(o.context, o.DatacentersGetDatacentersDatacenterIDHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/domains"] = domains.NewGetDomains(o.context, o.DomainsGetDomainsHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/domains/{domain_id}"] = domains.NewGetDomainsDomainID(o.context, o.DomainsGetDomainsDomainIDHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/members"] = members.NewGetMembers(o.context, o.MembersGetMembersHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/members/{member_id}"] = members.NewGetMembersMemberID(o.context, o.MembersGetMembersMemberIDHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/monitors"] = monitors.NewGetMonitors(o.context, o.MonitorsGetMonitorsHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/monitors/{monitor_id}"] = monitors.NewGetMonitorsMonitorID(o.context, o.MonitorsGetMonitorsMonitorIDHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/pools"] = pools.NewGetPools(o.context, o.PoolsGetPoolsHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/pools/{pool_id}"] = pools.NewGetPoolsPoolID(o.context, o.PoolsGetPoolsPoolIDHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/quotas"] = administrative.NewGetQuotas(o.context, o.AdministrativeGetQuotasHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/quotas/defaults"] = administrative.NewGetQuotasDefaults(o.context, o.AdministrativeGetQuotasDefaultsHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/quotas/{project_id}"] = administrative.NewGetQuotasProjectID(o.context, o.AdministrativeGetQuotasProjectIDHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/services"] = administrative.NewGetServices(o.context, o.AdministrativeGetServicesHandler)
	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/datacenters"] = datacenters.NewPostDatacenters(o.context, o.DatacentersPostDatacentersHandler)
	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/domains"] = domains.NewPostDomains(o.context, o.DomainsPostDomainsHandler)
	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/members"] = members.NewPostMembers(o.context, o.MembersPostMembersHandler)
	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/monitors"] = monitors.NewPostMonitors(o.context, o.MonitorsPostMonitorsHandler)
	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/pools"] = pools.NewPostPools(o.context, o.PoolsPostPoolsHandler)
	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/sync"] = administrative.NewPostSync(o.context, o.AdministrativePostSyncHandler)
	if o.handlers["PUT"] == nil {
		o.handlers["PUT"] = make(map[string]http.Handler)
	}
	o.handlers["PUT"]["/datacenters/{datacenter_id}"] = datacenters.NewPutDatacentersDatacenterID(o.context, o.DatacentersPutDatacentersDatacenterIDHandler)
	if o.handlers["PUT"] == nil {
		o.handlers["PUT"] = make(map[string]http.Handler)
	}
	o.handlers["PUT"]["/domains/{domain_id}"] = domains.NewPutDomainsDomainID(o.context, o.DomainsPutDomainsDomainIDHandler)
	if o.handlers["PUT"] == nil {
		o.handlers["PUT"] = make(map[string]http.Handler)
	}
	o.handlers["PUT"]["/members/{member_id}"] = members.NewPutMembersMemberID(o.context, o.MembersPutMembersMemberIDHandler)
	if o.handlers["PUT"] == nil {
		o.handlers["PUT"] = make(map[string]http.Handler)
	}
	o.handlers["PUT"]["/monitors/{monitor_id}"] = monitors.NewPutMonitorsMonitorID(o.context, o.MonitorsPutMonitorsMonitorIDHandler)
	if o.handlers["PUT"] == nil {
		o.handlers["PUT"] = make(map[string]http.Handler)
	}
	o.handlers["PUT"]["/pools/{pool_id}"] = pools.NewPutPoolsPoolID(o.context, o.PoolsPutPoolsPoolIDHandler)
	if o.handlers["PUT"] == nil {
		o.handlers["PUT"] = make(map[string]http.Handler)
	}
	o.handlers["PUT"]["/quotas/{project_id}"] = administrative.NewPutQuotasProjectID(o.context, o.AdministrativePutQuotasProjectIDHandler)
}

// Serve creates a http handler to serve the API over HTTP
// can be used directly in http.ListenAndServe(":8000", api.Serve(nil))
func (o *AndromedaAPI) Serve(builder middleware.Builder) http.Handler {
	o.Init()

	if o.Middleware != nil {
		return o.Middleware(builder)
	}
	if o.useSwaggerUI {
		return o.context.APIHandlerSwaggerUI(builder)
	}
	return o.context.APIHandler(builder)
}

// Init allows you to just initialize the handler cache, you can then recompose the middleware as you see fit
func (o *AndromedaAPI) Init() {
	if len(o.handlers) == 0 {
		o.initHandlerCache()
	}
}

// RegisterConsumer allows you to add (or override) a consumer for a media type.
func (o *AndromedaAPI) RegisterConsumer(mediaType string, consumer runtime.Consumer) {
	o.customConsumers[mediaType] = consumer
}

// RegisterProducer allows you to add (or override) a producer for a media type.
func (o *AndromedaAPI) RegisterProducer(mediaType string, producer runtime.Producer) {
	o.customProducers[mediaType] = producer
}

// AddMiddlewareFor adds a http middleware to existing handler
func (o *AndromedaAPI) AddMiddlewareFor(method, path string, builder middleware.Builder) {
	um := strings.ToUpper(method)
	if path == "/" {
		path = ""
	}
	o.Init()
	if h, ok := o.handlers[um][path]; ok {
		o.handlers[method][path] = builder(h)
	}
}
