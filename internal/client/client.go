// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"os"

	"github.com/apex/log"
	"github.com/go-openapi/runtime"
	runtimeclient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	osclient "github.com/gophercloud/utils/v2/client"
	"github.com/gophercloud/utils/v2/openstack/clientconfig"
	"github.com/jedib0t/go-pretty/table"
	"github.com/jessevdk/go-flags"
	"github.com/jmoiron/sqlx/reflectx"

	"github.com/sapcc/andromeda/client"
)

var (
	Parser          = flags.NewParser(&opts, flags.Default)
	Table           = table.NewWriter()
	Mapper          = reflectx.NewMapper("json")
	AndromedaClient = client.Default
)

type outputFormatters struct {
	Format     string   `short:"f" long:"format" description:"The output format, defaults to table" choice:"table" choice:"csv" choice:"markdown" choice:"html" choice:"value" choice:"json" default:"table"`
	Columns    []string `short:"c" long:"column" description:"specify the column(s) to include, can be repeated to show multiple columns"`
	SortColumn []string `long:"sort-column" description:"specify the column(s) to sort the data (columns specified first have a priority, non-existing columns are ignored), can be repeated"`
}

var opts struct {
	Debug      bool             `long:"debug" description:"Show verbose debug information"`
	Formatters outputFormatters `group:"Output formatters"`
	Wait       bool             `long:"wait" description:"Wait for the operation to complete (if supported)"`

	OSEndpoint          string `long:"os-endpoint" env:"OS_ENDPOINT" description:"The endpoint that will always be used"`
	OSAuthUrl           string `long:"os-auth-url" env:"OS_AUTH_URL" description:"Authentication URL"`
	OSPassword          string `long:"os-password" env:"OS_PASSWORD" description:"User's password to use with"`
	OSUsername          string `long:"os-username" env:"OS_USERNAME" description:"User's username to use with"`
	OSProjectDomainName string `long:"os-project-domain-name" env:"OS_PROJECT_DOMAIN_NAME" description:"Domain name containing project"`
	OSProjectName       string `long:"os-project-name" env:"OS_PROJECT_NAME" description:"Project name to scope to"`
	OSRegionName        string `long:"os-region-name" env:"OS_REGION_NAME" description:"Authentication region name"`
	OSUserDomainName    string `long:"os-user-domain-name" env:"OS_USER_DOMAIN_NAME" description:"User's domain name"`
}

func SetupClient() {
	Table.SetOutputMirror(os.Stdout)

	Parser.CommandHandler = func(command flags.Commander, args []string) error {
		if command == nil {
			return nil
		}

		ao, err := clientconfig.AuthOptions(&clientconfig.ClientOpts{
			RegionName: opts.OSRegionName,
			AuthInfo: &clientconfig.AuthInfo{
				AuthURL:           opts.OSAuthUrl,
				Username:          opts.OSUsername,
				Password:          opts.OSPassword,
				ProjectName:       opts.OSProjectName,
				ProjectDomainName: opts.OSProjectDomainName,
				UserDomainName:    opts.OSUserDomainName,
				AllowReauth:       true,
			},
		})
		if err != nil {
			return err
		}

		provider, err := openstack.NewClient(opts.OSAuthUrl)
		if err != nil {
			return err
		}
		if opts.Debug {
			provider.HTTPClient = http.Client{
				Transport: &osclient.RoundTripper{
					Rt:     &http.Transport{},
					Logger: &osclient.DefaultLogger{},
				},
			}
		}

		err = openstack.Authenticate(context.Background(), provider, *ao)
		if err != nil {
			return err
		}

		endpointOpts := gophercloud.EndpointOpts{
			Region: opts.OSRegionName,
		}
		endpointOpts.ApplyDefaults("gtm")
		endpoint, err := provider.EndpointLocator(endpointOpts)
		if err != nil {
			return err
		}
		// Override endpoint?
		if opts.OSEndpoint != "" {
			endpoint = opts.OSEndpoint
		}

		uri, err := url.Parse(endpoint)
		if err != nil {
			return err
		}

		rt := runtimeclient.New(uri.Host, uri.Path, []string{uri.Scheme})
		rt.SetDebug(opts.Debug)
		rt.DefaultAuthentication = runtime.ClientAuthInfoWriterFunc(func(req runtime.ClientRequest, reg strfmt.Registry) error {
			if err := req.SetHeaderParam("X-Auth-Token", provider.Token()); err != nil {
				return err
			}
			return nil
		})
		AndromedaClient.SetTransport(rt)

		return command.Execute(args)
	}

	if _, err := Parser.Parse(); err != nil {
		code := 1
		var fe *flags.Error
		if errors.As(err, &fe) {
			if errors.Is(fe.Type, flags.ErrHelp) {
				code = 0
			} else {
				log.WithError(err).Error("while parsing command line arguments")
			}
		}
		os.Exit(code)
	}
}
