// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package utils

import "net/http"

var JSONHeader = http.Header{"Content-Type": {"application/json; charset=utf-8"}}
