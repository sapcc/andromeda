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
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/apex/log"
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/mcuadros/go-defaults"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

var (
	Global Andromeda
)

// ParseArgsAndRun Evaluate --config-file flags and remove them from env
func ParseArgsAndRun(name string, usage string, action cli.ActionFunc, flags ...cli.Flag) {
	// Append --config-file
	flags = append(flags, &cli.StringSliceFlag{
		Name:  "config-file",
		Usage: "Use config file",
	})

	app := &cli.App{
		Name:   name,
		Usage:  usage,
		Flags:  flags,
		Action: action,
		Before: func(c *cli.Context) error {
			if !c.IsSet("config-file") {
				return errors.New("no config files specified")
			}

			// Set defaults
			defaults.SetDefaults(&Global)

			// Parse config flags
			if err := parseConfigFlags(c.StringSlice("config-file")); err != nil {
				return err
			}

			// Ugly workaround, remove flags from osArgs
			// because micro wants to parse them and fails miserably
			i := 0
			newOsArgs := []string{}
			for {
				if i >= len(os.Args) {
					break
				}

				reaped := false

				for _, flagName := range c.FlagNames() {
					if strings.HasSuffix(os.Args[i], flagName) {
						// Reap the flag from osArgs
						reaped = true
						i++
						break
					}
				}
				if !reaped {
					newOsArgs = append(newOsArgs, os.Args[i])
				}
				i++
			}
			os.Args = newOsArgs
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		currentErr := err
		for errors.Unwrap(currentErr) != nil {
			log.Fatal(currentErr.Error())
			currentErr = errors.Unwrap(currentErr)
		}
		log.Fatal(currentErr.Error())
	}
}

func parseConfigFlags(flags []string) error {
	log.Infof("Config files: %+v", flags)
	for _, path := range flags {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		decoder := yaml.NewDecoder(file)
		if err = decoder.Decode(&Global); err != nil {
			return err
		}
		if err = file.Close(); err != nil {
			return err
		}
	}

	if Global.Default.Debug {
		log.SetLevel(log.DebugLevel)
	}
	return nil
}

type Andromeda struct {
	Default      Default               `yaml:"DEFAULT"`
	Database     Database              `yaml:"database"`
	ApiSettings  ApiSettings           `yaml:"api_settings"`
	ServiceAuth  clientconfig.AuthInfo `yaml:"service_auth"`
	Quota        Quota                 `yaml:"quota"`
	F5Config     F5Config              `yaml:"f5"`
	AkamaiConfig AkamaiConfig          `yaml:"akamai"`
	Audit        Audit                 `yaml:"audit_middleware_notifications"`
	HouseKeeping HouseKeeping          `yaml:"house_keeping"`
}

type ApiSettings struct {
	PolicyFile                string  `yaml:"policy-file" description:"Use policy file" default:"policy.json"`
	AuthStrategy              string  `yaml:"auth_strategy" description:"The auth strategy for API requests, currently supported: [keystone, none]"`
	PolicyEngine              string  `yaml:"policy_engine" description:"Policy engine to use, currently supported: [goslo, noop]"`
	DisablePagination         bool    `yaml:"disable_pagination" description:"Disable the usage of pagination"`
	DisableSorting            bool    `yaml:"disable_sorting" description:"Disable the usage of sorting"`
	PaginationMaxLimit        int64   `yaml:"pagination_max_limit" default:"1000" description:"The maximum number of items returned in a single response."`
	RateLimit                 float64 `yaml:"rate_limit" default:"100" description:"Maximum number of requests to limit per second."`
	DisableCors               bool    `yaml:"disable_cors" description:"Stops sending Access-Control-Allow-Origin Header to allow cross-origin requests."`
	EnableProxyHeadersParsing bool    `yaml:"enable_proxy_headers_parsing" default:"true" description:"Try parsing proxy headers for http scheme and base url."`
}

type Quota struct {
	Enabled                bool  `yaml:"enabled" description:"Enable quotas."`
	DefaultQuotaDomain     int64 `yaml:"domains" default:"0" description:"Default quota of domain per project."`
	DefaultQuotaPool       int64 `yaml:"pools" default:"0" description:"Default quota of pool per project."`
	DefaultQuotaMember     int64 `yaml:"members" default:"0" description:"Default quota of member per project."`
	DefaultQuotaMonitor    int64 `yaml:"monitors" default:"0" description:"Default quota of monitor per project."`
	DefaultQuotaDatacenter int64 `yaml:"datacenters" default:"0" description:"Default quota of datacenter per project."`
}

type Default struct {
	Debug            bool   `yaml:"debug" description:"Enable debug mode."`
	Host             string `yaml:"host" description:"Hostname used by the server/agent. Defaults to auto-discovery."`
	ApiBaseURL       string `yaml:"api_base_uri" description:"Base URI for the API for use in pagination links. This will be autodetected from the request if not overridden here."`
	TransportURL     string `yaml:"transport_url" description:"The network address and optional user credentials for connecting to the messaging backend."`
	Prometheus       bool   `yaml:"prometheus" description:"Enable prometheus exporter."`
	PrometheusListen string `yaml:"prometheus_listen" default:"127.0.0.1:9090" description:"Prometheus listen TCP network address."`
	SentryDSN        string `yaml:"sentry_dsn" description:"Sentry Data Source Name."`
}

type Database struct {
	Connection string `yaml:"connection" description:"Connection string to use to connect to the database."`
}

type F5Config struct {
	DNSServerAddress string `yaml:"dns_server_address"`
	ValidateCert     bool   `yaml:"validate_certificates"`
}

type AkamaiConfig struct {
	EdgeRC               string `yaml:"edgerc" description:"Path to akamai edgerc file, else sourced from AKAMAI_EDGE_RC env variable."`
	Domain               string `yaml:"domain" description:"Traffic Management Domain to use (e.g. production.akadns.net)."`
	DomainType           string `yaml:"domain_type" description:"Indicates the type of domain available based on your contract, defaults to autodetect. Either failover-only, static, weighted, basic, or full."`
	ContractId           string `yaml:"contract_id" description:"Indicated the contract id to use, autodetects if only one contract is associated."`
	SyncInterval         int64  `yaml:"sync_interval" default:"30" description:"Sync interval for checking for pending updates"`
	MemberStatusInterval int64  `yaml:"member_status_interval" default:"60" description:"Sync interval for checking for member status"`
}

type Audit struct {
	Enabled      bool   `yaml:"enabled" description:"Enables message notification bus."`
	TransportURL string `yaml:"transport_url" description:"The network address and optional user credentials for connecting to the messaging backend."`
	QueueName    string `yaml:"queue_name" description:"RabbitMQ queue name"`
}

type HouseKeeping struct {
	Enabled     bool  `yaml:"enabled" description:"Enables house keeping."`
	DeleteAfter int64 `yaml:"delete_after" default:"600" description:"Minimum seconds elapsed after cleanup of a deleted domain."`
}

func GetApiBaseUrl(r *http.Request) string {
	var baseUrl url.URL

	baseUrl.Scheme = "http"
	if r.TLS != nil {
		baseUrl.Scheme = "https"
	}
	baseUrl.Host = Global.Default.Host
	if Global.ApiSettings.EnableProxyHeadersParsing {
		if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
			baseUrl.Scheme = proto
		}
		if host := r.Header.Get("X-Forwarded-Host"); host != "" {
			baseUrl.Host = host
		}
	}

	return baseUrl.String()
}
