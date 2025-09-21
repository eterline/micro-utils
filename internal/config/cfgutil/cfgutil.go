// Copyright (c) 2025 EterLine (Andrew)
// This file is part of micro-utils.
// Licensed under the MIT License. See the LICENSE file for details.

package cfgutil

import (
	"os"
	"path/filepath"

	"github.com/alexflint/go-arg"
)

type UsualConfig[T any] struct {
	Config *T
	Name   string
}

func (uc *UsualConfig[T]) ParseArgs() (T, error) {

	parserConfig := arg.Config{
		Program:           uc.Name,
		IgnoreEnv:         false,
		IgnoreDefault:     false,
		StrictSubcommands: true,
	}

	p, err := arg.NewParser(parserConfig, uc.Config)
	if err != nil {
		return *new(T), err
	}

	err = p.Parse(os.Args[1:])
	switch err {
	case arg.ErrHelp:
		p.WriteHelp(os.Stdout)
		os.Exit(1)
	case nil:
		return *uc.Config, nil
	}

	return *new(T), err
}

func (uc *UsualConfig[T]) selfExec() string {
	exePath, err := os.Executable()
	if err != nil {
		return uc.Name
	}

	return filepath.Base(exePath)
}
