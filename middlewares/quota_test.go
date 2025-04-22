package middlewares

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProviderFromHTTPRequestAcceptsRegularPayload(t *testing.T) {
	jsonPayload := minimalNewDomainRequestPayload()
	w := httptest.NewRecorder()
	r, err := http.NewRequest("POST", "/domains", bytes.NewBuffer(jsonPayload))
	r.Header.Set("Content-Type", "application/json")
	assert.Nil(t, err)
	expected := "akamai"
	actual, err := providerFromHTTPRequest("domain", w, r)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestProviderFromHTTPRequestRejectsMaliciousPayload(t *testing.T) {
	jsonPayload := arbitrarilySizedNewDomainRequestPayload(2 * 1024 * 1024) // 2MB
	w := httptest.NewRecorder()
	r, err := http.NewRequest("POST", "/domains", bytes.NewBuffer(jsonPayload))
	r.Header.Set("Content-Type", "application/json")
	assert.Nil(t, err)
	_, err = providerFromHTTPRequest("domain", w, r)
	assert.NotNil(t, err, "Unreasonably large request payloads should be rejected")
}

func minimalNewDomainRequestPayload() []byte {
	return []byte(`{"domain": {"provider": "akamai"}}`)
}

func arbitrarilySizedNewDomainRequestPayload(size int) []byte {
	payload := make([]byte, size)
	copy(payload, minimalNewDomainRequestPayload())
	return payload
}
