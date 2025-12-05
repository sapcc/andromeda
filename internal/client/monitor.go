// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"errors"
	"fmt"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag/conv"

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
	if err = waitForActiveMonitor(resp.Payload.Monitor.ID, false); err != nil {
		return fmt.Errorf("failed to wait for monitor %s to be active", resp.Payload.Monitor.ID)
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
	if err := waitForActiveMonitor(MonitorOptions.MonitorDelete.Positional.UUID, true); err != nil {
		return fmt.Errorf("failed to wait for monitor %s to be deleted",
			MonitorOptions.MonitorDelete.Positional.UUID)
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
		monitor.Monitor.AdminStateUp = conv.Pointer(false)
	} else if MonitorOptions.MonitorSet.Enable {
		monitor.Monitor.AdminStateUp = conv.Pointer(true)
	}

	params := monitors.
		NewPutMonitorsMonitorIDParams().
		WithMonitorID(MonitorOptions.MonitorSet.Positional.UUID).
		WithMonitor(monitor)
	resp, err := AndromedaClient.Monitors.PutMonitorsMonitorID(params)
	if err != nil {
		return err
	}
	if err = waitForActiveMonitor(MonitorOptions.MonitorSet.Positional.UUID, false); err != nil {
		return fmt.Errorf("failed to wait for monitor %s to be active",
			MonitorOptions.MonitorSet.Positional.UUID)
	}
	return WriteTable(resp.GetPayload().Monitor)
}

// waitForActiveMonitor waits for the monitor to be active, or optionally be deleted
func waitForActiveMonitor(id strfmt.UUID, deleted bool) error {
	// if not waiting, return immediately
	if !opts.Wait {
		return nil
	}

	return RetryWithBackoffMax(func() error {
		params := monitors.NewGetMonitorsMonitorIDParams().WithMonitorID(id)
		r, err := AndromedaClient.Monitors.GetMonitorsMonitorID(params)
		if err != nil {
			var getIDNotFound *monitors.GetMonitorsMonitorIDNotFound
			if errors.As(err, &getIDNotFound) && deleted {
				return nil
			}
			return err
		}

		res := r.GetPayload()
		if deleted || res.Monitor.ProvisioningStatus != models.MonitorProvisioningStatusACTIVE {
			return fmt.Errorf("monitor %s is not active yet", id)
		}
		return nil
	})
}

func init() {
	_, err := Parser.AddCommand("monitor", "Monitors", "Monitor Commands.", &MonitorOptions)
	if err != nil {
		panic(err)
	}
}
