// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package models

import "fmt"

// Error implements the error interface so *Error can be returned from functions
// that signal both HTTP-level validation failures and lower-level errors.
func (m *Error) Error() string {
	return fmt.Sprintf("HTTP %d: %s", m.Code, m.Message)
}
