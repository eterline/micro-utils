// Copyright (c) 2025 EterLine (Andrew)
// This file is part of micro-utils.
// Licensed under the MIT License. See the LICENSE file for details.

package ips2ubntest

type Configuration struct {
	Addrs     []string `arg:"-a,--addrs" help:"Search ip address or domain. Can be list or single value."`
	OutFile   string   `arg:"-o,--out" help:"Output file with subnets."`
	InputFile string   `arg:"-i,--in" help:"Input file with ip addresses."`
	Separator string   `arg:"-s,--sep" help:"Output file separator."`
}
