// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package as3

type Result struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Host    string `json:"host"`
	Tenant  string `json:"tenant"`
	RunTime int64  `json:"runTime"`
}

type Response struct {
	Results []Result `json:"results"`
}
