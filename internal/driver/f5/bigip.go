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

// bigIPSession is a minimal subset of struct type github.com/f5devcentral/go-bigip/BigIP
type bigIPSession interface {
	// APICall is defined by BigIP
	// It's best suited for iControlREST GET requests
	APICall(options *bigip.APIRequest) ([]byte, error)

	// PostAs3Bigip is defined by BigIP
	// It's necessary for POST /mgmt/shared/appsvcs/declare custom response payload handling
	PostAs3Bigip(as3NewJson, tenantFilter, queryParam string) (error, string, string)

	// GetDevices is defined by BigIP
	GetDevices() ([]bigip.Device, error)

	// GetHost is unique to this interface
	// It's a test-friendly shorthand for BigIP.Host
	GetHost() string
}

type bigIP struct {
	*bigip.BigIP
}

func newBigIPSession(b *bigip.BigIP) *bigIP {
	return &bigIP{BigIP: b}
}

func (b *bigIP) GetHost() string {
	return b.BigIP.Host
}

type activeDeviceMatcher func(bigIPSession) (*bigip.Device, error)
type deviceSessionFactory func(url string) (bigIPSession, error)

func getActiveDeviceSession(
	conf config.F5Config,
	factory deviceSessionFactory,
	matcher activeDeviceMatcher) (bigIPSession, *bigip.Device, error) {
	var s bigIPSession
	var d *bigip.Device
	for _, url := range conf.Devices {
		session, err := factory(url)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create session: %s", err)
		}
		device, err := matchActiveDevice(session)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get active device: %s", err)
		}
		if device != nil {
			s = session
			d = device
			break
		}
	}
	return s, d, nil
}

func matchActiveDevice(session bigIPSession) (*bigip.Device, error) {
	hostname, err := getSessionHostnameFromURL(session.GetHost())
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname from F5 device session: %s", err)
	}
	devices, err := session.GetDevices()
	if err != nil {
		return nil, fmt.Errorf("failed to get devices from F5 device session: %s", err)
	}
	device, err := filterDeviceMatchingHostnameSuffix(devices, hostname)
	if err != nil {
		return nil, fmt.Errorf("failed to filter F5 device matching session hostname: %v", err)
	}
	if device.FailoverState != "active" {
		return nil, nil
	}
	return device, nil
}

func getSessionHostnameFromURL(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	if parsedURL.Hostname() != "" {
		return parsedURL.Hostname(), nil
	}
	return rawURL, nil
}

func filterDeviceMatchingHostnameSuffix(devices []bigip.Device, hostname string) (*bigip.Device, error) {
	for _, device := range devices {
		if strings.HasSuffix(hostname, device.Hostname) {
			return &device, nil
		}
	}
	return nil, fmt.Errorf("device %s not found", hostname)
}

func getBigIPSession(rawURL string) (bigIPSession, error) {
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
	return newBigIPSession(bigip.NewSession(&bigip.Config{
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
	})), nil
}
