// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package f5

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/actatum/stormrpc"
	"github.com/apex/log"
	"github.com/f5devcentral/go-bigip"
	"github.com/nats-io/nats.go"

	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/rpc/server"
	"github.com/sapcc/andromeda/internal/utils"
)

type StatusController struct {
	bigIP *bigip.BigIP
	rpc   server.RPCServerClient
}

func ExecuteF5StatusAgent() error {
	log.Debugf("Enabled=%+v Devices=%v VCMPs=%v PhysicalNetwork=%v",
		config.Global.F5Config.Enabled,
		config.Global.F5Config.Devices,
		config.Global.F5Config.VCMPs,
		config.Global.F5Config.PhysicalNetwork,
	)

	var activeF5Session *bigip.BigIP
	for _, url := range config.Global.F5Config.Devices {
		deviceSession, err := GetBigIPSession(url)
		if err != nil {
			return fmt.Errorf("failed to acquire F5 device session: %v", err)
		}
		device, err := GetActiveDevice(deviceSession)
		if err != nil {
			return fmt.Errorf("failed to determine whether F5 device is active: %v", err)
		}
		if device != nil {
			activeF5Session = deviceSession
			log.Infof("Connected to F5 device [marketing name = %q, name = %q, version = %s, edition = %q, failover state = %q]",
				device.MarketingName, device.Name, device.Version, device.Edition, device.FailoverState)
		}
	}

	if activeF5Session == nil {
		return errors.New("failed to determine active F5 session")
	}

	// RPC server
	nc, err := nats.Connect(config.Global.Default.TransportURL)
	if err != nil {
		return err
	}
	client, err := stormrpc.NewClient("", stormrpc.WithNatsConn(nc))
	if err != nil {
		return err
	}

	sc := StatusController{
		activeF5Session,
		server.NewRPCServerClient(client),
	}

	if config.Global.Default.Prometheus {
		go utils.PrometheusListen()
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	// TODO: exit the go routine gracefully
	go func() {
		// TODO: activate immediately, then tick
		t := time.NewTicker(30 * time.Second)
		defer t.Stop()
		for {
			<-t.C // Activate periodically
			if err := sc.StatusHandler(); err != nil {
				log.Error(err.Error())
			}
		}
	}()

	<-done
	log.Info("Shutting down")
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return nil
}

type MembersCollectionStats struct {
	Kind     string                     `json:"kind"`
	SelfLink string                     `json:"selfLink"`
	Entries  map[string]json.RawMessage `json:"entries"`
}

type MembersStats struct {
	NestedStats struct {
		Kind     string `json:"kind"`
		SelfLink string `json:"selfLink"`
		Entries  struct {
			Alternate struct {
				Value int `json:"value"`
			} `json:"alternate"`
			Fallback struct {
				Value int `json:"value"`
			} `json:"fallback"`
			PoolName struct {
				Description string `json:"description"`
			} `json:"poolName"`
			PoolType struct {
				Description string `json:"description"`
			} `json:"poolType"`
			Preferred struct {
				Value int `json:"value"`
			} `json:"preferred"`
			ServerName struct {
				Description string `json:"description"`
			} `json:"serverName"`
			StatusAvailabilityState struct {
				Description string `json:"description"`
			} `json:"status.availabilityState"`
			StatusEnabledState struct {
				Description string `json:"description"`
			} `json:"status.enabledState"`
			StatusStatusReason struct {
				Description string `json:"description"`
			} `json:"status.statusReason"`
			VsName struct {
				Description string `json:"description"`
			} `json:"vsName"`
		} `json:"entries"`
	} `json:"nestedStats"`
}

func (c StatusController) StatusHandler() error {
	log.Debug("StatusHandler started")
	type gtmPoolA_ResPayload struct {
		Items []struct {
			FullPath string `json:"fullPath"`
		} `json:"items"`
	}
	var resPayload gtmPoolA_ResPayload
	err, ok := GetForEntity(c.bigIP, &resPayload, "gtm/pool/a")
	if !ok {
		return fmt.Errorf("gtm/pool/a not found")
	}
	if err != nil {
		return err
	}

	log.Debugf("Got %d 'gtm/pool/a' items", len(resPayload.Items))
	var stats MembersCollectionStats
	var msrs []*server.MemberStatusRequest_MemberStatus

	for _, pool := range resPayload.Items {
		partition := ConvertPartitionPath(pool.FullPath)
		log.Debugf("Got pool partition %q", partition)
		err, ok := GetForEntity(c.bigIP, &stats, fmt.Sprintf("gtm/pool/a/%s/members/stats", partition))
		if err != nil {
			return err
		}
		if !ok {
			// we skip a pool for which no member stats were found
			continue
		}
		for _, partitions := range stats.Entries {
			memberStats := MembersStats{}
			if err := json.Unmarshal(partitions, &memberStats); err != nil {
				log.Warnf("Could not decode nested member stats [partition = %q]: %s", partition, err)
				continue
			}

			memberID := strings.TrimPrefix(memberStats.NestedStats.Entries.VsName.Description, "member_")
			memberStatus := server.MemberStatusRequest_MemberStatus_UNKNOWN
			switch memberStats.NestedStats.Entries.StatusAvailabilityState.Description {
			case "available":
				memberStatus = server.MemberStatusRequest_MemberStatus_ONLINE
			case "offline":
				memberStatus = server.MemberStatusRequest_MemberStatus_OFFLINE
			}

			msr := &server.MemberStatusRequest_MemberStatus{
				Id:     memberID,
				Status: memberStatus,
			}
			log.Debugf("Member %s has status %s", memberID, msr.GetStatus().String())

			msrs = append(msrs, msr)
		}
	}
	_, err = c.rpc.UpdateMemberStatus(context.Background(),
		&server.MemberStatusRequest{MemberStatus: msrs})
	if err != nil {
		return err
	}
	log.Infof("Refreshed status of %d members", len(msrs))
	log.Debug("StatusHandler finished")

	return nil
}
