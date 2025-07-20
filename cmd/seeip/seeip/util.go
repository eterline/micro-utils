package seeip

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

func GetIpInfoWithContext(ctx context.Context, ip net.IP) (InfoIP, error) {

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
		return InfoIP{}, fmt.Errorf("request error: %w", err)
	}

	res, err := cl.Do(req)
	if err != nil {
		return InfoIP{}, fmt.Errorf("failed to get info about (%s): %w", ip.String(), err)
	}
	defer res.Body.Close()

	var info InfoIP
	if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
		return InfoIP{}, fmt.Errorf("failed to decode response for (%s): %w", ip.String(), err)
	}

	return info, nil
}

type InfoIP struct {
	Query       string  `json:"query,omitempty"`
	Status      string  `json:"status,omitempty"`
	Country     string  `json:"country,omitempty"`
	CountryCode string  `json:"countryCode,omitempty"`
	Region      string  `json:"region,omitempty"`
	RegionName  string  `json:"regionName,omitempty"`
	Isp         string  `json:"isp,omitempty"`
	Org         string  `json:"org,omitempty"`
	As          string  `json:"as,omitempty"`
	City        string  `json:"city,omitempty"`
	Zip         string  `json:"zip,omitempty"`
	Lat         float64 `json:"lat,omitempty"`
	Lon         float64 `json:"lon,omitempty"`
	Timezone    string  `json:"timezone,omitempty"`
	Message     string  `json:"message,omitempty"`
}

func (i InfoIP) Failed() (failedFrom error) {
	if i.Status == "fail" {
		return fmt.Errorf("ip information error: %s", i.Message)
	}
	return nil
}

func (i InfoIP) Latitude() float64 {
	return i.Lat
}

func (i InfoIP) Longitude() float64 {
	return i.Lon
}
