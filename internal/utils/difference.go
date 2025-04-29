// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package utils

import "github.com/go-openapi/strfmt"

func UUIDsDifference(a, b []strfmt.UUID) []strfmt.UUID {
	// calculates the difference (a - b) of two slices of UUIDs
	mb := make(map[strfmt.UUID]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []strfmt.UUID
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
