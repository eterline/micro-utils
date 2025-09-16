// Copyright (c) 2025 EterLine (Andrew)
// This file is part of My-Go-Project.
// Licensed under the MIT License. See the LICENSE file for details.


package api

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

var (
	apiApiCom, _ = url.Parse("https://ip-api.api.eterline.space/json")
)

type ApiResponseJSON struct {
	Message     string  `json:"message"`
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Query       net.IP  `json:"query"`
}

func (i ApiResponseJSON) Failed() (failedFrom error) {
	if i.Status == "fail" {
		return fmt.Errorf("ip information error: %s", i.Message)
	}
	return nil
}

func (i ApiResponseJSON) Latitude() float64 {
	return i.Lat
}

func (i ApiResponseJSON) Longitude() float64 {
	return i.Lon
}

func GetIpInfo(ip net.IP) (ApiResponseJSON, error) {
	return GetIpInfoWithContext(context.Background(), ip)
}

func GetIpInfoWithContext(ctx context.Context, ip net.IP) (ApiResponseJSON, error) {

	cl := http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	target := apiApiCom.JoinPath(ip.String()).String()

	req, err := http.NewRequestWithContext(ctx, "GET", target, nil)
	if err != nil {
		return ApiResponseJSON{}, fmt.Errorf("request error: %w", err)
	}

	res, err := cl.Do(req)
	if err != nil {
		return ApiResponseJSON{}, fmt.Errorf("failed to get info about (%s): %w", ip.String(), err)
	}
	defer res.Body.Close()

	var info ApiResponseJSON
	if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
		return ApiResponseJSON{}, fmt.Errorf("failed to decode response for (%s): %w", ip.String(), err)
	}

	return info, nil
}
