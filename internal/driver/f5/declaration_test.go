// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package f5

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildsA33Declaration(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("foo", "foo", "that makes sense")
	t.Log("Hello 👋")
}
