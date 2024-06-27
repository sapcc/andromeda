/*
 *   Copyright 2022 SAP SE
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

package akamai

import (
	"encoding/json"

	"github.com/apex/log"
)

// Max returns the larger of x or y.
func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// DetectTypeFromFeatures detects best domain type according to matrix
// https://techdocs.akamai.com/gtm/reference/api-workflow
func DetectTypeFromFeatures(features []string) string {
	const (
		FailoverOnly int = iota
		Static
		Weighted
		Full
	)

	bestDomainType := 0
	for _, feature := range features {
		switch feature {
		case "PERFORMANCE":
			bestDomainType = Full
		case "WEIGHTED_ROUND_ROBIN_LOAD_FEEDBACK":
			bestDomainType = Max(bestDomainType, Weighted)
		case "WEIGHTED_ROUND_ROBIN":
			bestDomainType = Max(bestDomainType, Weighted)
		case "WEIGHTED_HASHED":
			bestDomainType = Max(bestDomainType, Weighted)
		case "ASMAPPING":
			bestDomainType = Max(bestDomainType, Static)
		case "CIDRMAPPING":
			bestDomainType = Max(bestDomainType, Static)
		case "GEOGRAPHIC":
			bestDomainType = Max(bestDomainType, Static)
		case "STATIC":
			bestDomainType = Max(bestDomainType, Static)
		case "QTR":
			bestDomainType = Max(bestDomainType, FailoverOnly)
		case "FAILOVER":
			bestDomainType = Max(bestDomainType, FailoverOnly)
		}
	}

	// Resolve type
	switch bestDomainType {
	case FailoverOnly:
		return "failover-only"
	case Static:
		return "static"
	case Weighted:
		return "weighted"
	case Full:
		return "full"
	default:
		return "basic"
	}
}

func PrettyJson(data interface{}) string {
	val, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	return string(val)
}
