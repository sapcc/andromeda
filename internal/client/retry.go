/*
 *   Copyright 2024 SAP SE
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

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
