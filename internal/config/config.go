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
	Verbose      bool                  `short:"v" long:"verbose" description:"Show verbose debug information"`
	Default      Default               `json:"DEFAULT"`
	Database     Database              `json:"database"`
	ApiSettings  ApiSettings           `json:"api_settings"`
	ServiceAuth  clientconfig.AuthInfo `json:"service_auth"`
	Quota        Quota                 `json:"quota"`
	F5Config     F5Config              `json:"f5"`
	AkamaiConfig AkamaiConfig          `json:"akamai"`
}

type ApiSettings struct {
	PolicyFile         string  `short:"p" json:"policy-file" description:"Use policy file" default:"policy.json"`
	AuthStrategy       string  `json:"auth_strategy" description:"The auth strategy for API requests, currently supported: [keystone, none]"`
	PolicyEngine       string  `json:"policy_engine" description:"Policy engine to use, currently supported: [goslo, noop]"`
	DisablePagination  bool    `json:"disable_pagination" description:"Disable the usage of pagination"`
	DisableSorting     bool    `json:"disable_sorting" description:"Disable the usage of sorting"`
	PaginationMaxLimit int64   `json:"pagination_max_limit" default:"1000" description:"The maximum number of items returned in a single response."`
	RateLimit          float64 `json:"rate_limit" default:"100" description:"Maximum number of requests to limit per second."`
	DisableCors        bool    `json:"disable_cors" description:"Stops sending Access-Control-Allow-Origin Header to allow cross-origin requests."`
}

type Quota struct {
	Enabled                bool  `json:"enabled" description:"Enable quotas."`
	DefaultQuotaDomain     int64 `json:"default_quota_domain" default:"0" description:"Default quota of domain per project."`
	DefaultQuotaPool       int64 `json:"default_quota_pool" default:"0" description:"Default quota of pool per project."`
	DefaultQuotaMember     int64 `json:"default_quota_member" default:"0" description:"Default quota of member per project."`
	DefaultQuotaMonitor    int64 `json:"default_quota_monitor" default:"0" description:"Default quota of monitor per project."`
	DefaultQuotaDatacenter int64 `json:"default_quota_datacenter" default:"0" description:"Default quota of datacenter per project."`
}

type Default struct {
	ApiBaseURL   string `json:"api_base_uri" description:"Base URI for the API for use in pagination links. This will be autodetected from the request if not overridden here."`
	TransportURL string `json:"transport_url" description:"The network address and optional user credentials for connecting to the messaging backend."`
}

type Database struct {
	Connection string `json:"connection" description:"Connection string to use to connect to the database."`
}

type F5Config struct {
	Host             string `json:"host"`
	DNSServerAddress string `json:"dns_server_address"`
	ValidateCert     bool   `json:"validate_certificates"`
}

type AkamaiConfig struct {
	EdgeRC     string `json:"edgerc" description:"Path to akamai edgerc file, else sourced from AKAMAI_EDGE_RC env variable."`
	Domain     string `json:"domain" description:"Traffic Management Domain to use (e.g. production.akadns.net)."`
	DomainType string `json:"domain_type" description:"Indicates the type of domain available based on your contract, defaults to autodetect. Either failover-only, static, weighted, basic, or full."`
	ContractId string `json:"contract_id" description:"Indicated the contract id to use, autodetects if only one contract is associated."`
}
