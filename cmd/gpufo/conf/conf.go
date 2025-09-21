// Copyright (c) 2025 EterLine (Andrew)
// This file is part of micro-utils.
// Licensed under the MIT License. See the LICENSE file for details.

package conf

import (
	"os"
	"path/filepath"

	"github.com/alexflint/go-arg"
)

type Configuration struct {
	SmiTables  []string `arg:"-t,--tables" help:"Search ip address or domain"`
	JsonFormat bool     `arg:"-j,--json" help:"JSON output format"`
	Pretty     bool     `arg:"-p,--pretty" help:"JSON output pretty style"`
	Tree       bool     `arg:"-n,--tree" help:"Tree nodes format"`
	Camel      bool     `arg:"-c,--camel" help:"Camel notation keys format"`
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
