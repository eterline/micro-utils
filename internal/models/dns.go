// Copyright (c) 2025 EterLine (Andrew)
// This file is part of micro-utils.
// Licensed under the MIT License. See the LICENSE file for details.

package models

import (
	"context"
	"net"
	"time"
)

type DnsProviderType string

const (
	DnsCloudflareProvider DnsProviderType = "cloudflare"
	DnsGoogleProvider     DnsProviderType = "google"
	DnsLocalProvider      DnsProviderType = "local"
)

type Resolver interface {
	ResolveIP(ctx context.Context, s string) ([]net.IP, error)
	ResolveNS(ctx context.Context, s string) ([]string, error)
}

type AboutResolve struct {
	IPs               []net.IP `json:"ip,omitempty" yaml:"ip,omitempty"`
	NameServers       []string `json:"ns,omitempty" yaml:"ns,omitempty"`
	ErrorIPs          string   `json:"ip_error,omitempty" yaml:"ip_error,omitempty"`
	ErrorNS           string   `json:"ns_error,omitempty" yaml:"ns_error,omitempty"`
	ResolveDurationMs int64    `json:"resolve_duration_ms" yaml:"resolve_duration_ms"`
}

func (ar *AboutResolve) CalcDuration(start time.Time) {
	end := time.Now()
	ar.ResolveDurationMs = end.Sub(start).Milliseconds()
}
