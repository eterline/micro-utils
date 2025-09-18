// Copyright (c) 2025 EterLine (Andrew)
// This file is part of micro-utils.
// Licensed under the MIT License. See the LICENSE file for details.


package models

import (
	"context"
	"net"
)

type Resolver interface {
	ResolveIP(ctx context.Context, s string) ([]net.IP, error)
	ResolveNS(ctx context.Context, s string) ([]string, error)
}

type AboutResolve struct {
	IPs         []net.IP `json:"ip,omitempty" yaml:"ip,omitempty,omitempty"`
	NameServers []string `json:"ns,omitempty" yaml:"ns,omitempty,omitempty"`
	ErrorIPs    string   `json:"ip_error,omitempty" yaml:"ip_error,omitempty"`
	ErrorNS     string   `json:"ns_error,omitempty" yaml:"ns_error,omitempty"`
}
