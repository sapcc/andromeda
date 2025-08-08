// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package f5

import (
	"os"
	"testing"

	"github.com/f5devcentral/go-bigip"
	"github.com/stretchr/testify/assert"
)

func TestGetActiveDeviceSession(t *testing.T) {
	t.Skip()
}

func TestMatchActiveDevice(t *testing.T) {
	t.Skip()
}

func TestGetSessionHostname(t *testing.T) {
	assert := assert.New(t)
	t.Run("Fails if malformed URL", func(t *testing.T) {
		_, err := getSessionHostname(&bigip.BigIP{Host: "://foo.bar"})
		assert.Error(err)
	})
	t.Run("Succeeds with parsed hostname, if any", func(t *testing.T) {
		hostname, err := getSessionHostname(&bigip.BigIP{Host: "https://foo.bar"})
		assert.Nil(err)
		assert.Equal("foo.bar", hostname)
	})
	t.Run("Succeeds but falls back to BigIP.Host if parsed URL has no hostname", func(t *testing.T) {
		// TODO not sure what's the point of supporting this scenario
		hostname, err := getSessionHostname(&bigip.BigIP{Host: "---"})
		assert.Nil(err)
		assert.Equal("---", hostname)
	})
}

func TestFilterDeviceMatchingHostnameSuffix(t *testing.T) {
	assert := assert.New(t)
	t.Run("Fails if no matches", func(t *testing.T) {
		devices := []bigip.Device{
			{Hostname: "foo"},
			{Hostname: "bar"},
		}
		_, err := filterDeviceMatchingHostnameSuffix(devices, "a-very-long-hostname")
		assert.Error(err)
	})
	t.Run("Succeeds if suffix matches", func(t *testing.T) {
		devices := []bigip.Device{
			{Hostname: "hostname-foo"},
			{Hostname: "long-hostname"},
		}
		d, err := filterDeviceMatchingHostnameSuffix(devices, "a-very.long-hostname")
		assert.Nil(err)
		assert.Equal(d.Hostname, "long-hostname")
	})
}

func TestGetBigIPSession(t *testing.T) {
	assert := assert.New(t)
	t.Run("Rejects mal-formed URL", func(t *testing.T) {
		_, err := getBigIPSession("badURL")
		assert.Error(err)
	})
	t.Run("Fails if credentials are missing from both URL and environment", func(t *testing.T) {
		_, err := getBigIPSession("f5-device-1.local")
		if assert.Error(err) {
			assert.ErrorContains(err, "BIGIP_USER")
		}
	})
	t.Run("Succeeds without environment variables", func(t *testing.T) {
		s, err := getBigIPSession("https://johndoe:insecure@f5-device-1.local")
		assert.Nil(err)
		assert.Equal("https://f5-device-1.local", s.Host)
		assert.Equal("johndoe", s.User)
		assert.Equal("insecure", s.Password)
	})
	t.Run("Succeeds with environment variables (fallback if credentials not in URL)", func(t *testing.T) {
		assert.Nil(os.Setenv("BIGIP_USER", "johndoe"))
		assert.Nil(os.Setenv("BIGIP_PASSWORD", "insecure"))
		defer func() {
			assert.Nil(os.Unsetenv("BIGIP_USER"))
			assert.Nil(os.Unsetenv("BIGIP_PASSWORD"))
		}()
		s, err := getBigIPSession("https://f5-device-1.local")
		assert.Nil(err)
		assert.Equal("https://f5-device-1.local", s.Host)
		assert.Equal("johndoe", s.User)
		assert.Equal("insecure", s.Password)
	})
}
