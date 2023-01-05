/*
 *   Copyright 2021 SAP SE
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

package client

import (
	"github.com/go-openapi/strfmt"
	"github.com/sapcc/andromeda/client/administrative"
)

type SyncDomain struct {
	Domains []strfmt.UUID `short:"d" long:"domain" required:"true" description:"Domain IDs to sync, can be specified multiple times"`
}

func (sd SyncDomain) Execute(_ []string) error {
	body := administrative.PostSyncBody{Domains: sd.Domains}
	params := administrative.NewPostSyncParams()
	params.Domains = body
	_, err := AndromedaClient.Administrative.PostSync(params)
	return err
}

func init() {
	_, err := Parser.AddCommand("sync", "Sync", "Sync Commands.", &SyncDomain{})
	if err != nil {
		panic(err)
	}
}
