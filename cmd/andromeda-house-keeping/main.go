// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/urfave/cli/v2"

	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/housekeeping"
)

func main() {
	config.ParseArgsAndRun("andromeda", "andromeda house keeping",
		func(c *cli.Context) error {
			return housekeeping.HouseKeeping()
		},
	)
}
