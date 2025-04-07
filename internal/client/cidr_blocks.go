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
	"encoding/json"
	"fmt"

	"github.com/sapcc/andromeda/client/administrative"
)

type CidrBlocks struct {
	Provider string `short:"v" long:"provider" description:"Provider name" required:"false"`
}

func (cb CidrBlocks) Execute(_ []string) error {
	params := administrative.NewGetCidrBlocksParams()
	if cb.Provider != "" {
		params.SetProvider(&cb.Provider)
	}
	resp, err := AndromedaClient.Administrative.GetCidrBlocks(params)
	if err != nil {
		return err
	}

	type cidrBlocks struct {
		Description  string `json:"description,omitempty"`
		Cidr         string `json:"cidr,omitempty"`
		CidrMask     string `json:"cidrMask,omitempty"`
		CreationDate string `json:"creationDate,omitempty"`
		ChangeDate   string `json:"changeDate,omitempty"`
	}

	var dat []cidrBlocks
	if err = json.Unmarshal([]byte(fmt.Sprintf("%s", resp.Payload)), &dat); err != nil {
		return err
	}
	return WriteTable(dat)
}

func init() {
	_, err := Parser.AddCommand("cidr-blocks", "CIDR Blocks",
		"Show CIDR Blocks used by provided Agents.", &CidrBlocks{})
	if err != nil {
		panic(err)
	}
}
