// Copyright (c) 2025 EterLine (Andrew)
// This file is part of My-Go-Project.
// Licensed under the MIT License. See the LICENSE file for details.


package seeip

import (
	"os"
	"path/filepath"

	"github.com/alexflint/go-arg"
)

type Configuration struct {
	Address    []string `arg:"-a,--addr" help:"Search ip address or domain"`
	MapUrl     bool     `arg:"-m,--map" help:"Return map location url"`
	JsonFormat bool     `arg:"-j,--json" help:"JSON output format"`
	Pretty     bool     `arg:"-p,--pretty" help:"JSON output pretty style"`
	GeoStamp   bool     `arg:"-g,--geo" help:"Map geo position links"`
}

var (
	parserConfig = arg.Config{
		Program:           selfExec(),
		IgnoreEnv:         false,
		IgnoreDefault:     false,
		StrictSubcommands: true,
	}
)

func ParseArgs(c *Configuration) error {
	p, err := arg.NewParser(parserConfig, c)
	if err != nil {
		return err
	}

	err = p.Parse(os.Args[1:])
	if err == arg.ErrHelp {
		p.WriteHelp(os.Stdout)
		os.Exit(1)
	}
	return err
}

func selfExec() string {
	exePath, err := os.Executable()
	if err != nil {
		return "monita"
	}

	return filepath.Base(exePath)
}
