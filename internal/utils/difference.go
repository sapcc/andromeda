/*
 *   Copyright 2020 SAP SE
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
