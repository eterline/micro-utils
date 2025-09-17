// Copyright (c) 2025 EterLine (Andrew)
// This file is part of My-Go-Project.
// Licensed under the MIT License. See the LICENSE file for details.


package ipdata

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/eterline/micro-utils/internal/models"
)

func InitDnsCloudflareProvider() models.ProviderDoH {
	return &DnsDoHProvider{
		upstream: "https://cloudflare-dns.com/dns-query?%s",
		httpClient: &http.Client{
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
		},
	}
}

func InitDnsGoogleProvider() models.ProviderDoH {
	return &DnsDoHProvider{
		upstream: "https://dns.google/resolve?%s",
		httpClient: &http.Client{
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
		},
	}
}

// DnsDoHProvider - provider for DoH cloudflare | google
type DnsDoHProvider struct {
	upstream   string
	httpClient *http.Client
}

func (c *DnsDoHProvider) Query(ctx context.Context, d models.Domain, t models.DnsRecordType) (models.DnsResponse, error) {
	name, err := d.Punycode()
	if err != nil {
		return models.DnsResponse{}, err
	}

	param := url.Values{}
	param.Add("name", name)
	param.Add("type", t.String())
	dnsURL := fmt.Sprintf(c.upstream, param.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, dnsURL, nil)
	if err != nil {
		return models.DnsResponse{}, err
	}

	req.Header.Set("Accept", "application/dns-json")

	r, err := c.httpClient.Do(req)
	if err != nil {
		return models.DnsResponse{}, err
	}

	res := models.DnsResponse{}

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		return models.DnsResponse{}, err
	}

	if res.Status != 0 {
		return res, fmt.Errorf("cloudflare: bad response code: %d", r.Status)
	}

	return res, nil
}
