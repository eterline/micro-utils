// Copyright (c) 2025 EterLine (Andrew)
// This file is part of micro-utils.
// Licensed under the MIT License. See the LICENSE file for details.

package doh

func InitDnsCloudflareProvider() *DnsDoHProvider {
	return &DnsDoHProvider{
		serviceName: "cloudflare",
		upstream:    "https://cloudflare-dns.com/dns-query?%s",
		httpClient:  setupHttpClient(),
	}
}

func InitDnsGoogleProvider() *DnsDoHProvider {
	return &DnsDoHProvider{
		serviceName: "google",
		upstream:    "https://dns.google/resolve?%s",
		httpClient:  setupHttpClient(),
	}
}
