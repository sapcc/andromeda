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
	"github.com/sapcc/andromeda/client/members"
	"github.com/sapcc/andromeda/models"
)

var MemberOptions struct {
	MemberList   `command:"list" description:"List Members"`
	MemberShow   `command:"show" description:"Show Member"`
	MemberCreate `command:"create" description:"Create Member"`
	MemberDelete `command:"delete" description:"Delete Member"`
}

type MemberList struct {
	PositionalMemberList struct {
		PoolID strfmt.UUID `description:"UUID of the pool"`
	} `positional-args:"yes" required:"yes"`
}

type MemberShow struct {
	PositionalMemberShow struct {
		PoolID   strfmt.UUID `description:"UUID of the pool"`
		MemberID strfmt.UUID `description:"UUID of the member"`
	} `positional-args:"yes" required:"yes"`
}

type MemberCreate struct {
	PositionalMemberCreate struct {
		PoolID strfmt.UUID `description:"UUID of the pool"`
	} `positional-args:"yes" required:"yes"`
	Name    string      `short:"n" long:"name" description:"Name of the Member"`
	Address strfmt.IPv4 `short:"a" long:"address" description:"Address of the Member" required:"true"`
	Port    int64       `short:"p" long:"port" description:"Port of the Member" required:"true"`
	Disable bool        `short:"d" long:"disable" description:"Disable Member" optional:"true" optional-value:"false"`
}

type MemberDelete struct {
	PositionalMemberDelete struct {
		PoolID strfmt.UUID `description:"UUID of the pool"`
		UUID   strfmt.UUID `description:"UUID of the member"`
	} `positional-args:"yes" required:"yes"`
}

func (*MemberList) Execute(_ []string) error {
	resp, err := AndromedaClient.Members.GetPoolsPoolIDMembers(members.
		NewGetPoolsPoolIDMembersParams().
		WithPoolID(MemberOptions.PositionalMemberList.PoolID))
	if err != nil {
		return err
	}
	return WriteTable(resp.GetPayload().Members)
}

func (*MemberCreate) Execute(_ []string) error {
	adminStateUp := !MemberOptions.Disable
	member := members.PostPoolsPoolIDMembersBody{&models.Member{
		AdminStateUp: &adminStateUp,
		Name:         &MemberOptions.Name,
		Address:      &MemberOptions.Address,
		Port:         &MemberOptions.Port,
	}}
	resp, err := AndromedaClient.Members.PostPoolsPoolIDMembers(members.
		NewPostPoolsPoolIDMembersParams().
		WithPoolID(MemberOptions.PositionalMemberCreate.PoolID).
		WithMember(member))
	if err != nil {
		return err
	}
	return WriteTable(resp.GetPayload().Member)
}

func (*MemberShow) Execute(_ []string) error {
	params := members.
		NewGetPoolsPoolIDMembersMemberIDParams().
		WithPoolID(MemberOptions.PositionalMemberShow.PoolID).
		WithMemberID(MemberOptions.MemberShow.PositionalMemberShow.MemberID)
	resp, err := AndromedaClient.Members.GetPoolsPoolIDMembersMemberID(params)
	if err != nil {
		return err
	}
	return WriteTable(resp.GetPayload().Member)
}

func (*MemberDelete) Execute(_ []string) error {
	params := members.
		NewDeletePoolsPoolIDMembersMemberIDParams().
		WithPoolID(MemberOptions.MemberDelete.PositionalMemberDelete.PoolID).
		WithMemberID(MemberOptions.MemberDelete.PositionalMemberDelete.UUID)

	if _, err := AndromedaClient.Members.DeletePoolsPoolIDMembersMemberID(params); err != nil {
		return err
	}
	return nil
}

func init() {
	_, err := Parser.AddCommand("member", "Members", "Member Commands.", &MemberOptions)
	if err != nil {
		panic(err)
	}
}
