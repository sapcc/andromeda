// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package f5

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/sapcc/andromeda/internal/config"
	"github.com/scottdware/go-bigip"
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

func GetBigIPSession() (*bigip.BigIP, error) {
	parsedURL, err := url.Parse(config.Global.Default.Host)
	if err != nil {
		return nil, err
	}

	// check for password
	pw, ok := parsedURL.User.Password()
	if !ok {
		return nil, fmt.Errorf("password required for host '%s'", parsedURL.Hostname())
	}

	// todo: make configurable
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	session, err := bigip.NewTokenSession(
		parsedURL.Host,
		parsedURL.User.Username(),
		pw,
		"tmos",
		nil)
	if err != nil {
		return nil, err
	}
	return session, nil
}
