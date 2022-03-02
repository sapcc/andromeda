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

package config

import (
	"github.com/gophercloud/utils/openstack/clientconfig"
)

var (
	Global Andromeda
)

type Andromeda struct {
	Verbose     bool                  `short:"v" long:"verbose" description:"Show verbose debug information"`
	ConfigFile  []string              `short:"c" long:"config-file" description:"Use config file" no-ini:"true"`
	Default     Default               `group:"DEFAULT"`
	Database    Database              `group:"database"`
	ApiSettings ApiSettings           `group:"api_settings"`
	ServiceAuth clientconfig.AuthInfo `group:"service_auth"`
	AgentConfig AgentConfig           `group:"agent"`
	Quota       Quota                 `group:"quota"`
}

type ApiSettings struct {
	PolicyFile         string  `short:"p" long:"policy-file" description:"Use policy file" default:"policy.json"`
	AuthStrategy       string  `long:"auth_strategy" description:"The auth strategy for API requests, currently supported: [keystone, none]"`
	PolicyEngine       string  `long:"policy_engine" description:"Policy engine to use, currently supported: [goslo, noop]"`
	DisablePagination  bool    `long:"disable_pagination" description:"Disable the usage of pagination"`
	DisableSorting     bool    `long:"disable_sorting" description:"Disable the usage of sorting"`
	PaginationMaxLimit int64   `long:"pagination_max_limit" default:"1000" description:"The maximum number of items returned in a single response."`
	RateLimit          float64 `long:"rate_limit" default:"100" description:"Maximum number of requests to limit per second."`
	DisableCors        bool    `long:"disable_cors" description:"Stops sending Access-Control-Allow-Origin Header to allow cross-origin requests."`
}

type Quota struct {
	Enabled                bool  `long:"enabled" description:"Enable quotas."`
	DefaultQuotaDomain     int64 `long:"default_quota_domain" default:"0" description:"Default quota of domain per project."`
	DefaultQuotaPool       int64 `long:"default_quota_pool" default:"0" description:"Default quota of pool per project."`
	DefaultQuotaMember     int64 `long:"default_quota_member" default:"0" description:"Default quota of member per project."`
	DefaultQuotaMonitor    int64 `long:"default_quota_monitor" default:"0" description:"Default quota of monitor per project."`
	DefaultQuotaDatacenter int64 `long:"default_quota_datacenter" default:"0" description:"Default quota of datacenter per project."`
}

type Default struct {
	ApiBaseURL   string `long:"api_base_uri" description:"Base URI for the API for use in pagination links. This will be autodetected from the request if not overridden here."`
	TransportURL string `long:"transport_url" description:"The network address and optional user credentials for connecting to the messaging backend."`
}

type Database struct {
	Connection string `long:"connection" description:"Connection string to use to connect to the database."`
}

type AgentConfig struct {
	Host             string `long:"host"`
	DNSServerAddress string `long:"dns_server_address"`
	ValidateCert     bool   `long:"validate_certificates"`
}
