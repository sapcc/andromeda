// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package f5

import (
	"errors"
	"os"
	"testing"

	"github.com/f5devcentral/go-bigip"
	"github.com/sapcc/andromeda/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestGetActiveDeviceSession(t *testing.T) {
	assert := assert.New(t)
	conf := config.F5Config{
		Devices: []string{"https://a.local", "https://b.local"},
	}

	t.Run("Fails if device session factory fails", func(t *testing.T) {
		factory := func(url string) (bigIPSession, error) {
			return nil, errors.New("please let the caller know I failed")
		}
		_, _, err := getActiveDeviceSession(conf, factory, nil)
		if assert.Error(err) {
			assert.ErrorContains(err, "failed to create session")
		}
	})

	t.Run("Fails if device matcher fails", func(t *testing.T) {
		s := new(mockedBigIPSession)
		factory := func(url string) (bigIPSession, error) {
			s.On("GetHost").Return("http://c.local")
			s.On("GetDevices").Return([]bigip.Device{
				{Name: "a", Hostname: "a.local", FailoverState: "passive"},
				{Name: "b", Hostname: "b.local", FailoverState: "active"},
			}, nil)
			return s, nil
		}
		_, _, err := getActiveDeviceSession(conf, factory, matchActiveDevice)
		s.AssertCalled(t, "GetHost")
		s.AssertCalled(t, "GetDevices")
		if assert.Error(err) {
			assert.ErrorContains(err, "failed to get active device")
		}
	})

	t.Run("Succeeds otherwise", func(t *testing.T) {
		s := new(mockedBigIPSession)
		factory := func(url string) (bigIPSession, error) {
			s.On("GetHost").Return("http://b.local")
			s.On("GetDevices").Return([]bigip.Device{
				{Name: "a", Hostname: "a.local", FailoverState: "passive"},
				{Name: "b", Hostname: "b.local", FailoverState: "active"},
			}, nil)
			return s, nil
		}
		session, device, err := getActiveDeviceSession(conf, factory, matchActiveDevice)
		s.AssertCalled(t, "GetHost")
		s.AssertCalled(t, "GetDevices")
		if assert.Nil(err) {
			assert.Equal("b", device.Name)
			assert.Equal("b.local", device.Hostname)
			assert.Same(s, session)
		}
	})

}

func TestMatchActiveDevice(t *testing.T) {
	assert := assert.New(t)

	t.Run("Fails if session has a bad hostname", func(t *testing.T) {
		s := new(mockedBigIPSession)
		s.On("GetHost").Return("://b.local")
		_, err := matchActiveDevice(s)
		s.AssertCalled(t, "GetHost")
		s.AssertNotCalled(t, "GetDevices")
		if assert.Error(err) {
			assert.ErrorContains(err, "failed to get hostname")
		}
	})

	t.Run("Fails if it cannot get devices list", func(t *testing.T) {
		s := new(mockedBigIPSession)
		s.On("GetHost").Return("https://b.local")
		s.On("GetDevices").Return([]bigip.Device{}, errors.New("please let the caller now I failed"))
		_, err := matchActiveDevice(s)
		s.AssertCalled(t, "GetHost")
		s.AssertCalled(t, "GetDevices")
		if assert.Error(err) {
			assert.ErrorContains(err, "failed to get devices")
		}
	})

	t.Run("Fails if it cannot find one device whose hostname matches the session host (suffix matching)", func(t *testing.T) {
		s := new(mockedBigIPSession)
		devices := []bigip.Device{
			{Name: "a", Hostname: "x.local", FailoverState: "passive"},
			{Name: "b", Hostname: "y.local", FailoverState: "active"},
		}
		s.On("GetHost").Return("https://b.local")
		s.On("GetDevices").Return(devices, nil)
		_, err := matchActiveDevice(s)
		s.AssertCalled(t, "GetHost")
		s.AssertCalled(t, "GetDevices")
		if assert.Error(err) {
			assert.ErrorContains(err, "failed to filter")
		}
	})

	t.Run("Fails if no 'active' device found", func(t *testing.T) {
		s := new(mockedBigIPSession)
		devices := []bigip.Device{
			{Name: "a", Hostname: "a.local", FailoverState: "passive"},
			{Name: "b", Hostname: "b.local", FailoverState: "passive"},
		}
		s.On("GetHost").Return("https://b.local")
		s.On("GetDevices").Return(devices, nil)
		d, err := matchActiveDevice(s)
		s.AssertCalled(t, "GetHost")
		s.AssertCalled(t, "GetDevices")
		assert.Nil(err)
		assert.Nil(d, "no active device should have been found")
	})

	t.Run("Succeeds with first 'active' device (the other is expected to be 'passive')", func(t *testing.T) {
		s := new(mockedBigIPSession)
		devices := []bigip.Device{
			{Name: "a", Hostname: "a.local", FailoverState: "passive"},
			{Name: "b", Hostname: "b.local", FailoverState: "active"},
		}
		s.On("GetHost").Return("https://b.local")
		s.On("GetDevices").Return(devices, nil)
		d, err := matchActiveDevice(s)
		s.AssertCalled(t, "GetHost")
		s.AssertCalled(t, "GetDevices")
		assert.Nil(err)
		if assert.NotNil(d, "one device (Name: 'b') should have been detected as 'active'") {
			assert.Equal(d.Name, "b")
			assert.Equal(d.Hostname, "b.local")
		}
	})
}

func TestGetSessionHostnameFromURL(t *testing.T) {
	assert := assert.New(t)

	t.Run("Fails if malformed URL", func(t *testing.T) {
		_, err := getSessionHostnameFromURL("://foo.bar")
		assert.Error(err)
	})

	t.Run("Succeeds with parsed hostname, if any", func(t *testing.T) {
		hostname, err := getSessionHostnameFromURL("https://foo.bar")
		assert.Nil(err)
		assert.Equal("foo.bar", hostname)
	})

	t.Run("Succeeds but falls back to BigIP.Host if parsed URL has no hostname", func(t *testing.T) {
		hostname, err := getSessionHostnameFromURL("---")
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
		assert.Equal("https://f5-device-1.local", s.GetHost())
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
		assert.Equal("https://f5-device-1.local", s.GetHost())
	})
}
