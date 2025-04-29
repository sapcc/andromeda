// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"fmt"
	"time"
)

const (
	startBackoff = 5 * time.Second
	maxBackoff   = 30 * time.Second
	maxRetries   = 60
)

// RetryWithBackoffMax retries the given operation with exponential backoff
func RetryWithBackoffMax(fn func() error) error {
	var lastError error
	for i := 0; i < maxRetries; i++ {
		if lastError = fn(); lastError == nil {
			return nil
		}

		if i < maxRetries-1 {
			backoff := startBackoff << uint(i)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
			time.Sleep(backoff)
		}
	}
	return fmt.Errorf("max retries exceeded: %w", lastError)
}
