// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package akamai

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/sapcc/andromeda/internal/rpcmodels"
)

func TestRPCHandlerAkamai_Sync(t *testing.T) {
	testCases := []struct {
		name             string
		channelSize      int
		domainIds        []string
		expectedDomains  []string
		expectNil        bool
	}{
		{
			name:            "sync with multiple domains",
			channelSize:     1,
			domainIds:       []string{"domain1", "domain2"},
			expectedDomains: []string{"domain1", "domain2"},
			expectNil:       false,
		},
		{
			name:            "sync with nil domain IDs",
			channelSize:     1,
			domainIds:       nil,
			expectedDomains: nil,
			expectNil:       true,
		},
		{
			name:            "sync with empty domain IDs",
			channelSize:     1,
			domainIds:       []string{},
			expectedDomains: []string{},
			expectNil:       false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock agent with a buffered channel
			mockAgent := &AkamaiAgent{
				forceSync: make(chan []string, tc.channelSize),
			}

			handler := &RPCHandlerAkamai{
				agent: mockAgent,
			}

			ctx := context.Background()
			request := &rpcmodels.SyncRequest{
				DomainId: tc.domainIds,
			}

			response, err := handler.Sync(ctx, request)

			assert.NoError(t, err, "Sync should not return error")
			assert.NotNil(t, response, "Response should not be nil")

			// Check that the domain IDs were sent to the forceSync channel
			select {
			case receivedDomains := <-mockAgent.forceSync:
				if tc.expectNil {
					assert.Nil(t, receivedDomains, "Domain IDs should be nil")
				} else {
					assert.Equal(t, tc.expectedDomains, receivedDomains, "Domain IDs should match")
				}
			case <-time.After(100 * time.Millisecond):
				t.Error("Expected domain IDs to be sent to forceSync channel")
			}
		})
	}
}