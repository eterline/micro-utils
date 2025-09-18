// Copyright (c) 2025 EterLine (Andrew)
// This file is part of micro-utils.
// Licensed under the MIT License. See the LICENSE file for details.


package main

import (
	"context"

	microutils "github.com/eterline/micro-utils"
	"github.com/eterline/micro-utils/cmd/gpufo/adapt/cli"
	"github.com/eterline/micro-utils/cmd/gpufo/conf"
)

var (
	cfg = conf.Configuration{
		SmiTables: []string{
			"identify",
			"clocks",
			"memory",
			"pci",
			"power",
			"states",
			"temp",
			"utilization",
		},
		JsonFormat: false,
		Pretty:     false,
		Tree:       false,
		Camel:      false,
	}
)

func main() {
	ctx := context.Background()

	err := conf.ParseArgs(&cfg)
	if err != nil {
		microutils.PrintFatalErr(err)
	}

	err = cli.HandleCLI(ctx, cfg)
	if err != nil {
		microutils.PrintFatalErr(err)
	}
}
