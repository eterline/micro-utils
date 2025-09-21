// Copyright (c) 2025 EterLine (Andrew)
// This file is part of micro-utils.
// Licensed under the MIT License. See the LICENSE file for details.

package main

import (
	"os"

	microutils "github.com/eterline/micro-utils"
	"github.com/eterline/micro-utils/internal/services/genuuid"
)

type Config struct {
	Count   int    `arg:"-c,--count" help:"UUID v4 count to generate"`
	Payload string `arg:"-p,--payload" help:"UUID hashed payload"`
	Domain  string `arg:"-d,--domain" help:"UUID domain namespace"`
	Version int    `arg:"-v,--version" help:"UUID version 3 or 5"`
}

var (
	Args = microutils.UsualConfig[Config]{
		Config: &Config{
			Count:   1,
			Payload: "",
			Domain:  "dns",
			Version: 3,
		},
		Name: "uuid",
	}
)

func main() {

	c, err := Args.ParseArgs()
	if err != nil {
		microutils.PrintFatalErr(err)
	}

	gen := genuuid.GenerateUUID{}

	if c.Payload == "" {
		uuids, err := gen.GenerateOf(c.Count)
		if err != nil {
			microutils.PrintFatalErr(err)
		}

		microutils.PrintYaml(uuids)
		os.Exit(0)
	}

	uuids, err := gen.GenerateFrom(c.Domain, c.Payload, c.Version)
	if err != nil {
		microutils.PrintFatalErr(err)
	}

	microutils.PrintYaml(uuids)
}
