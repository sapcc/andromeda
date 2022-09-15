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

	"github.com/sapcc/andromeda/client/monitors"
	"github.com/sapcc/andromeda/models"
)

var MonitorOptions struct {
	MonitorList   `command:"list" description:"List Monitors"`
	MonitorShow   `command:"show" description:"Show Monitor"`
	MonitorCreate `command:"create" description:"Create Monitor"`
	MonitorDelete `command:"delete" description:"Delete Monitor"`
}

type MonitorList struct {
}

type MonitorShow struct {
	Positional struct {
		UUID strfmt.UUID `description:"UUID of the monitor"`
	} `positional-args:"yes" required:"yes"`
}

type MonitorCreate struct {
	Name     string  `short:"n" long:"name" description:"Name of the Monitor"`
	Pool     string  `short:"p" long:"pool" description:"ID of the pool to check members" required:"true"`
	Type     *string `short:"t" long:"type" description:"Type of the health check monitor"`
	Interval *int64  `short:"i" long:"interval" description:"The interval, in seconds, between health checks." optional:"true"`
	Timeout  *int64  `short:"o" long:"timeout" description:"The time in total, in seconds, after which a health check times out" optional:"true"`
	Send     *string `short:"s" long:"send" description:"Specifies the text string that the monitor sends to the target member." optional:"true"`
	Receive  *string `short:"r" long:"receive" description:"Specifies the text string that the monitor expects to receive from the target member." optional:"true"`
	Disable  bool    `short:"d" long:"disable" description:"Disable Monitor" optional:"true" optional-value:"false"`
}

type MonitorDelete struct {
	Positional struct {
		UUID strfmt.UUID `description:"UUID of the monitor"`
	} `positional-args:"yes" required:"yes"`
}

func (*MonitorList) Execute(_ []string) error {
	resp, err := AndromedaClient.Monitors.GetMonitors(nil)
	if err != nil {
		return err
	}
	return WriteTable(resp.GetPayload().Monitors)
}

func (*MonitorCreate) Execute(_ []string) error {
	adminStateUp := !MonitorOptions.Disable
	poolID := strfmt.UUID(MonitorOptions.Pool)
	monitor := monitors.PostMonitorsBody{&models.Monitor{
		AdminStateUp: &adminStateUp,
		Name:         &MonitorOptions.Name,
		PoolID:       &poolID,
		Type:         MonitorOptions.Type,
		Interval:     MonitorOptions.Interval,
		Timeout:      MonitorOptions.Timeout,
		Send:         MonitorOptions.Send,
		Receive:      MonitorOptions.Receive,
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

func init() {
	_, err := Parser.AddCommand("monitor", "Monitors", "Monitor Commands.", &MonitorOptions)
	if err != nil {
		panic(err)
	}
}
