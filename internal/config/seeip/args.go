// Copyright (c) 2025 EterLine (Andrew)
// This file is part of micro-utils.
// Licensed under the MIT License. See the LICENSE file for details.


package seeip

type Configuration struct {
	Address         []string `arg:"-a,--addr" help:"Search ip address or domain. Can be list or single value."`
	ResolverService string   `arg:"-r,--reslov" help:"Resolver service name | DNS server address."`
	IsJson          bool     `arg:"-j,--json" help:"JSON object output."`
	Pretty          bool     `arg:"-f,--format" help:"JSON formatted object output."`
	Workers         int      `arg:"-w,--workers" help:"Process worker count."`
}
