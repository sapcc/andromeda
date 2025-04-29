// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/apex/log"

	"github.com/sapcc/andromeda/internal/driver/f5"
)

func main() {
	if err := f5.ExecuteF5StatusAgent(); err != nil {
		log.Fatal(err.Error())
	}
}
