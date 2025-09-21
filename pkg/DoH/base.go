// Copyright (c) 2025 EterLine (Andrew)
// This file is part of micro-utils.
// Licensed under the MIT License. See the LICENSE file for details.

package doh

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

// DnsDoHProvider - object of default DoH provider
type DnsDoHProvider struct {
	serviceName string
	upstream    string
	httpClient  *http.Client
}

// Service - get DoH provider name
func (c *DnsDoHProvider) Service() string {
	return c.serviceName
}

// Query - make DNS request to service
func (c *DnsDoHProvider) Query(ctx context.Context, d Domain, t Record) (DnsResponse, error) {
	name, err := d.Punycode()
	if err != nil {
		return DnsResponse{}, err
	}

	param := url.Values{}
	param.Add("name", name)
	param.Add("type", t.String())
	dnsURL := fmt.Sprintf(c.upstream, param.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, dnsURL, nil)
	if err != nil {
		return DnsResponse{}, err
	}

	req.Header.Set("Accept", "application/dns-json")

	r, err := c.httpClient.Do(req)
	if err != nil {
		return DnsResponse{}, err
	}

	res := DnsResponse{}

	defer r.Body.Close()
	if err := json.
		NewDecoder(r.Body).
		Decode(&res); err != nil {
		return DnsResponse{}, err
	}

	if res.Success() {
		return res, nil
	}

	return res, res.Status
}

func setupHttpClient() *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   3 * time.Second,
				KeepAlive: 60 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 3 * time.Second,
			DisableKeepAlives:   false,
			MaxIdleConns:        256,
			MaxIdleConnsPerHost: 256,
		},
	}
}
