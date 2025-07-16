// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package f5

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/f5devcentral/go-bigip"
	"github.com/sapcc/andromeda/internal/config"
)

func BigIPSupportsDNS(device *bigip.Device) error {
	// check if DNS module activated
	supported := false
	for _, module := range device.ActiveModules {
		// BigIP Appliance
		if strings.Contains(module, "DNS Services") {
			supported = true
		}
		// BigIP-VE Evaluation
		if strings.Contains(module, "DNS VE Lab") {
			supported = true
		}
	}
	if !supported {
		return fmt.Errorf("device %s does not support DNS Services", device.Name)
	}

	return nil
}

func GetBigIPSession(rawURL string) (*bigip.BigIP, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	user := parsedURL.User.Username()
	if user == "" {
		var ok bool
		user, ok = os.LookupEnv("BIGIP_USER")
		if !ok {
			return nil, fmt.Errorf("BIGIP_USER required for host '%s'", parsedURL.Hostname())
		}
	}

	// check for password
	password, ok := parsedURL.User.Password()
	if !ok {
		password, ok = os.LookupEnv("BIGIP_PASSWORD")
		if !ok {
			return nil, fmt.Errorf("BIGIP_PASSWORD required for host '%s'", parsedURL.Hostname())
		}
	}

	// todo: make configurable
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	return bigip.NewSession(&bigip.Config{
		Address:           parsedURL.Hostname(),
		Username:          user,
		Password:          password,
		LoginReference:    "tmos",
		CertVerifyDisable: !config.Global.F5Config.ValidateCert,
		ConfigOptions: &bigip.ConfigOptions{
			APICallTimeout: 60 * time.Second,
			TokenTimeout:   1200 * time.Second,
			APICallRetries: int(config.Global.F5Config.MaxRetries),
		},
	}), nil
}

func IsActiveF5DeviceSession(deviceSession *bigip.BigIP) (bool, error) {
	hostname, err := GetSessionHostname(deviceSession)
	if err != nil {
		return false, fmt.Errorf("failed to get hostname from F5 device session: %v", err)
	}
	devices, err := deviceSession.GetDevices()
	if err != nil {
		return false, fmt.Errorf("failed to get devices from F5 device session: %v", err)
	}
	device, err := FilterDeviceMatchingHostname(devices, hostname)
	if err != nil {
		return false, fmt.Errorf("failed to filter F5 device matching session hostname: %v", err)
	}
	return device.FailoverState == "active", nil
}

func GetActiveDevice(deviceSession *bigip.BigIP) (*bigip.Device, error) {
	hostname, err := GetSessionHostname(deviceSession)
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname from F5 device session: %v", err)
	}
	devices, err := deviceSession.GetDevices()
	if err != nil {
		return nil, fmt.Errorf("failed to get devices from F5 device session: %v", err)
	}
	device, err := FilterDeviceMatchingHostname(devices, hostname)
	if err != nil {
		return nil, fmt.Errorf("failed to filter F5 device matching session hostname: %v", err)
	}
	if device.FailoverState != "active" {
		return nil, nil
	}
	return device, nil
}
