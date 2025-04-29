// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

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
