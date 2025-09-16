// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package metrics

import (
	"sync"
	"time"

	"github.com/apex/log"
)

type AkamaiRateLimiter struct {
	maxReqPerMinute int
	ticker          <-chan time.Time
	tokensCh        chan struct{}
	mu              sync.Mutex
}

func NewAkamaiRateLimiter(maxReqPerMinute int) *AkamaiRateLimiter {
	rl := &AkamaiRateLimiter{
		maxReqPerMinute: maxReqPerMinute,
		ticker:          time.Tick(1 * time.Minute),
		tokensCh:        make(chan struct{}, maxReqPerMinute),
	}

	// ensure rate limiter is operational before returning it
	rl.refreshTokensChan()

	go func() {
		for range rl.ticker {
			rl.refreshTokensChan()
		}
	}()

	return rl
}

func (rl *AkamaiRateLimiter) UseToken() {
	startTime := time.Now()
	defer func() {
		rateLimitingDurationSeconds.Add(time.Since(startTime).Seconds())
	}()

	// block until a token is available
	<-rl.tokensCh
	log.Debug("[AkamaiRateLimiter] Token used")
}

func (rl *AkamaiRateLimiter) refreshTokensChan() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	log.Debugf("[AkamaiRateLimiter] Adding tokens to channel (channel size = %d)", rl.maxReqPerMinute)
LOOP:
	for {
		select {
		case rl.tokensCh <- struct{}{}:
			// empty: successfully wrote to channel
		default:
			log.Debug("[AkamaiRateLimiter] Finished adding tokens to channel")
			break LOOP
		}
	}
}
