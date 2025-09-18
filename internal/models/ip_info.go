// Copyright (c) 2025 EterLine (Andrew)
// This file is part of micro-utils.
// Licensed under the MIT License. See the LICENSE file for details.


package models

import (
	"fmt"
	"time"
)

/*
AboutIPobject - response ip info object

	More about: https://ip-api.com/docs/api:json
*/
type AboutIPobject struct {
	Status        string    `json:"status" yaml:"status"`
	Continent     string    `json:"continent" yaml:"continent"`
	ContinentCode string    `json:"continentCode" yaml:"continentCode"`
	Country       string    `json:"country" yaml:"country"`
	CountryCode   string    `json:"countryCode" yaml:"countryCode"`
	Region        string    `json:"region" yaml:"region"`
	RegionName    string    `json:"regionName" yaml:"regionName"`
	City          string    `json:"city" yaml:"city"`
	District      string    `json:"district" yaml:"district"`
	Zip           string    `json:"zip" yaml:"zip"`
	Lat           float64   `json:"lat" yaml:"lat"`
	Lon           float64   `json:"lon" yaml:"lon"`
	Timezone      string    `json:"timezone" yaml:"timezone"`
	Offset        int       `json:"offset" yaml:"offset"`
	Currency      string    `json:"currency" yaml:"currency"`
	Isp           string    `json:"isp" yaml:"isp"`
	Org           string    `json:"org" yaml:"org"`
	As            string    `json:"as" yaml:"as"`
	Asname        string    `json:"asname" yaml:"asname"`
	Reverse       string    `json:"reverse" yaml:"reverse"`
	Mobile        bool      `json:"mobile" yaml:"mobile"`
	Proxy         bool      `json:"proxy" yaml:"proxy"`
	Hosting       bool      `json:"hosting" yaml:"hosting"`
	RequestTime   time.Time `json:"-"`
}

func (ip AboutIPobject) MapLinks() map[string]string {
	return map[string]string{
		"google": fmt.Sprintf("https://www.google.com/maps?q=%f,%f", ip.Lat, ip.Lon),
		"osm":    fmt.Sprintf("https://www.openstreetmap.org/?mlat=%f&mlon=%f", ip.Lat, ip.Lon),
		"yandex": fmt.Sprintf("https://yandex.com/maps/?ll=%f,%f&z=12", ip.Lon, ip.Lat),
	}
}

type DnsProviderType string

const (
	DnsCloudflareProvider DnsProviderType = "cloudflare"
	DnsGoogleProvider     DnsProviderType = "google"
	DnsLocalProvider      DnsProviderType = "local"
)
