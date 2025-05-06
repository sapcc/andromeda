// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package policy

import (
	"github.com/sapcc/go-bits/gopherpolicy"
	"net/http"
)

type noOpPolicyEngine struct{}

func (p noOpPolicyEngine) init() {}

func (p noOpPolicyEngine) AuthorizeRequest(r *http.Request, t *gopherpolicy.Token, _ string) bool {
	return true
}

func (p noOpPolicyEngine) AuthorizeGetAllRequest(r *http.Request, t *gopherpolicy.Token, _ string) bool {
	return true
}
