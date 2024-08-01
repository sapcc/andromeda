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
	"fmt"

	"github.com/go-openapi/strfmt"

	"github.com/sapcc/andromeda/client/members"
	"github.com/sapcc/andromeda/models"
)

var MemberOptions struct {
	MemberList   `command:"list" description:"List Members"`
	MemberShow   `command:"show" description:"Show Member"`
	MemberCreate `command:"create" description:"Create Member"`
	MemberDelete `command:"delete" description:"Delete Member"`
	MemberSet    `command:"set" description:"Update Member"`
}

type MemberList struct {
	PositionalMemberList struct {
		PoolID strfmt.UUID `description:"UUID of the pool"`
	} `positional-args:"yes"`
}

type MemberShow struct {
	PositionalMemberShow struct {
		MemberID strfmt.UUID `description:"UUID of the member"`
	} `positional-args:"yes" required:"yes"`
}

type MemberCreate struct {
	PositionalMemberCreate struct {
		PoolID strfmt.UUID `description:"UUID of the pool"`
	} `positional-args:"yes" required:"yes"`
	Name         string      `short:"n" long:"name" description:"Name of the Member"`
	Address      strfmt.IPv4 `short:"a" long:"address" description:"Address of the Member" required:"true"`
	Port         int64       `short:"p" long:"port" description:"Port of the Member" required:"true"`
	Disable      bool        `short:"d" long:"disable" description:"Disable Member" optional:"true" optional-value:"false"`
	DatacenterID strfmt.UUID `short:"i" long:"datacenter" description:"Datacenter ID"`
}

type MemberDelete struct {
	PositionalMemberDelete struct {
		UUID strfmt.UUID `description:"UUID of the member"`
	} `positional-args:"yes" required:"yes"`
}

type MemberSet struct {
	PositionalMemberSet struct {
		UUID strfmt.UUID `description:"UUID of the member"`
	} `positional-args:"yes" required:"yes"`
	Name    *string      `short:"n" long:"name" description:"Name of the Member"`
	Address *strfmt.IPv4 `short:"a" long:"address" description:"Address of the Member"`
	Port    *int64       `short:"p" long:"port" description:"Port of the Member"`
	Disable bool         `short:"d" long:"disable" description:"Disable Member"`
	Enable  bool         `short:"e" long:"enable" description:"Enable Member"`
}

func (*MemberList) Execute(_ []string) error {
	resp, err := AndromedaClient.Members.GetMembers(members.
		NewGetMembersParams().
		WithPoolID(&MemberOptions.PositionalMemberList.PoolID))
	if err != nil {
		return err
	}
	return WriteTable(resp.GetPayload().Members)
}

func (*MemberCreate) Execute(_ []string) error {
	adminStateUp := !MemberOptions.MemberCreate.Disable
	member := members.PostMembersBody{Member: &models.Member{
		AdminStateUp: &adminStateUp,
		Name:         &MemberOptions.MemberCreate.Name,
		Address:      &MemberOptions.MemberCreate.Address,
		Port:         &MemberOptions.MemberCreate.Port,
		DatacenterID: &MemberOptions.DatacenterID,
		PoolID:       &MemberOptions.PositionalMemberCreate.PoolID,
	}}
	resp, err := AndromedaClient.Members.PostMembers(members.
		NewPostMembersParams().WithMember(member))
	if err != nil {
		return err
	}
	return WriteTable(resp.GetPayload().Member)
}

func (*MemberShow) Execute(_ []string) error {
	params := members.
		NewGetMembersMemberIDParams().
		WithMemberID(MemberOptions.MemberShow.PositionalMemberShow.MemberID)
	resp, err := AndromedaClient.Members.GetMembersMemberID(params)
	if err != nil {
		return err
	}
	return WriteTable(resp.GetPayload().Member)
}

func (*MemberDelete) Execute(_ []string) error {
	params := members.
		NewDeleteMembersMemberIDParams().
		WithMemberID(MemberOptions.MemberDelete.PositionalMemberDelete.UUID)

	if _, err := AndromedaClient.Members.DeleteMembersMemberID(params); err != nil {
		return err
	}
	return nil
}

func (*MemberSet) Execute(_ []string) error {
	if MemberOptions.MemberSet.Disable && MemberOptions.MemberSet.Enable {
		return fmt.Errorf("cannot enable and disable member at the same time")
	}

	member := members.PutMembersMemberIDBody{Member: &models.Member{
		Name:    MemberOptions.MemberSet.Name,
		Address: MemberOptions.MemberSet.Address,
		Port:    MemberOptions.MemberSet.Port,
	}}
	if MemberOptions.MemberSet.Disable {
		adminStateUp := false
		member.Member.AdminStateUp = &adminStateUp
	} else if MemberOptions.MemberSet.Enable {
		adminStateUp := true
		member.Member.AdminStateUp = &adminStateUp
	}

	params := members.
		NewPutMembersMemberIDParams().
		WithMemberID(MemberOptions.MemberSet.PositionalMemberSet.UUID).
		WithMember(member)

	resp, err := AndromedaClient.Members.PutMembersMemberID(params)
	if err != nil {
		return err
	}
	return WriteTable(resp.GetPayload().Member)
}

func init() {
	_, err := Parser.AddCommand("member", "Members", "Member Commands.", &MemberOptions)
	if err != nil {
		panic(err)
	}
}
