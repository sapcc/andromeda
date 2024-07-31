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
	"github.com/go-openapi/swag"

	"github.com/sapcc/andromeda/client/monitors"
	"github.com/sapcc/andromeda/models"
)

var MonitorOptions struct {
	MonitorList   `command:"list" description:"List Monitors"`
	MonitorShow   `command:"show" description:"Show Monitor"`
	MonitorCreate `command:"create" description:"Create Monitor"`
	MonitorDelete `command:"delete" description:"Delete Monitor"`
	MonitorSet    `command:"set" description:"Update Monitor"`
}

type MonitorList struct {
	PositionalMonitorList struct {
		PoolID strfmt.UUID `description:"UUID of the pool"`
	} `positional-args:"yes"`
}

type MonitorShow struct {
	Positional struct {
		UUID strfmt.UUID `description:"UUID of the monitor"`
	} `positional-args:"yes" required:"yes"`
}

type MonitorCreate struct {
	Name       string  `short:"n" long:"name" description:"Name of the Monitor"`
	Pool       string  `short:"p" long:"pool" description:"ID of the pool to check members" required:"true"`
	Type       *string `short:"t" long:"type" description:"Type of the health check monitor"`
	Interval   *int64  `short:"i" long:"interval" description:"The interval, in seconds, between health checks." optional:"true"`
	Timeout    *int64  `short:"o" long:"timeout" description:"The time in total, in seconds, after which a health check times out" optional:"true"`
	Send       *string `short:"s" long:"send" description:"Specifies the text string that the monitor sends to the target member." optional:"true"`
	Receive    *string `short:"r" long:"receive" description:"Specifies the text string that the monitor expects to receive from the target member." optional:"true"`
	Disable    bool    `short:"d" long:"disable" description:"Disable Monitor" optional:"true" optional-value:"false"`
	HTTPMethod *string `short:"m" long:"http-method" description:"HTTP method to use for monitor checks. Only used for HTTP/S monitors." choice:"GET" choice:"POST" choice:"PUT" choice:"HEAD" choice:"PATH" choice:"DELETE" choice:"OPTIONS"`
	DomainName *string `short:"D" long:"domain-name" description:"Domain name to use for monitor checks. Only used for HTTP/S monitors."`
}

type MonitorDelete struct {
	Positional struct {
		UUID strfmt.UUID `description:"UUID of the monitor"`
	} `positional-args:"yes" required:"yes"`
}

type MonitorSet struct {
	Positional struct {
		UUID strfmt.UUID `description:"UUID of the monitor"`
	} `positional-args:"yes" required:"yes"`
	Name       *string `short:"n" long:"name" description:"Name of the Monitor"`
	Disable    bool    `short:"d" long:"disable" description:"Enable Monitor" optional:"true" optional-value:"true"`
	Enable     bool    `short:"e" long:"enable" description:"Enable Monitor" optional:"true" optional-value:"true"`
	HTTPMethod *string `short:"m" long:"http-method" description:"HTTP method to use for monitor checks. Only used for HTTP/S monitors." choice:"GET" choice:"POST" choice:"PUT" choice:"HEAD" choice:"PATH" choice:"DELETE" choice:"OPTIONS"`
	DomainName *string `short:"D" long:"domain-name" description:"Domain name to use for monitor checks. Only used for HTTP/S monitors."`
}

func (*MonitorList) Execute(_ []string) error {
	resp, err := AndromedaClient.Monitors.GetMonitors(monitors.
		NewGetMonitorsParams().WithPoolID(&MonitorOptions.PositionalMonitorList.PoolID))
	if err != nil {
		return err
	}
	return WriteTable(resp.GetPayload().Monitors)
}

func (*MonitorCreate) Execute(_ []string) error {
	adminStateUp := !MonitorOptions.MonitorCreate.Disable
	poolID := strfmt.UUID(MonitorOptions.Pool)
	monitor := monitors.PostMonitorsBody{Monitor: &models.Monitor{
		AdminStateUp: &adminStateUp,
		Name:         &MonitorOptions.MonitorCreate.Name,
		PoolID:       &poolID,
		Type:         MonitorOptions.Type,
		Interval:     MonitorOptions.Interval,
		Timeout:      MonitorOptions.Timeout,
		Send:         MonitorOptions.Send,
		Receive:      MonitorOptions.Receive,
		HTTPMethod:   MonitorOptions.MonitorCreate.HTTPMethod,
		DomainName:   (*strfmt.Hostname)(MonitorOptions.MonitorCreate.DomainName),
	}}
	resp, err := AndromedaClient.Monitors.PostMonitors(monitors.NewPostMonitorsParams().WithMonitor(monitor))
	if err != nil {
		return err
	}
	return WriteTable(resp.GetPayload().Monitor)
}

func (*MonitorShow) Execute(_ []string) error {
	params := monitors.
		NewGetMonitorsMonitorIDParams().
		WithMonitorID(MonitorOptions.MonitorShow.Positional.UUID)
	resp, err := AndromedaClient.Monitors.GetMonitorsMonitorID(params)
	if err != nil {
		return err
	}
	return WriteTable(resp.GetPayload().Monitor)
}

func (*MonitorDelete) Execute(_ []string) error {
	params := monitors.
		NewDeleteMonitorsMonitorIDParams().
		WithMonitorID(MonitorOptions.MonitorDelete.Positional.UUID)

	if _, err := AndromedaClient.Monitors.DeleteMonitorsMonitorID(params); err != nil {
		return err
	}
	return nil
}

func (*MonitorSet) Execute(_ []string) error {
	if MonitorOptions.MonitorSet.Disable && MonitorOptions.MonitorSet.Enable {
		return fmt.Errorf("cannot enable and disable monitor at the same time")
	}

	monitor := monitors.PutMonitorsMonitorIDBody{Monitor: &models.Monitor{
		Name:       MonitorOptions.MonitorSet.Name,
		HTTPMethod: MonitorOptions.MonitorSet.HTTPMethod,
		DomainName: (*strfmt.Hostname)(MonitorOptions.MonitorSet.DomainName),
	}}
	if MonitorOptions.MonitorSet.Disable {
		monitor.Monitor.AdminStateUp = swag.Bool(false)
	} else if MonitorOptions.MonitorSet.Enable {
		monitor.Monitor.AdminStateUp = swag.Bool(true)
	}

	params := monitors.
		NewPutMonitorsMonitorIDParams().
		WithMonitorID(MonitorOptions.MonitorSet.Positional.UUID).
		WithMonitor(monitor)
	resp, err := AndromedaClient.Monitors.PutMonitorsMonitorID(params)
	if err != nil {
		return err
	}
	return WriteTable(resp.GetPayload().Monitor)
}

func init() {
	_, err := Parser.AddCommand("monitor", "Monitors", "Monitor Commands.", &MonitorOptions)
	if err != nil {
		panic(err)
	}
}
