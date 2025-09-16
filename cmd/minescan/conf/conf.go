// Copyright (c) 2025 EterLine (Andrew)
// This file is part of My-Go-Project.
// Licensed under the MIT License. See the LICENSE file for details.


package conf

import (
	"os"
	"path/filepath"

	"github.com/alexflint/go-arg"
)

type Configuration struct {
	Listen      string   `arg:"-l,--listen" help:"Listen connection addr"`
	FilterWords []string `arg:"-f,--filter" help:"Words in packets that will trigger scanner"`
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
