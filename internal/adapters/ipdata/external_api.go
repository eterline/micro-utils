// Copyright (c) 2025 EterLine (Andrew)
// This file is part of micro-utils.
// Licensed under the MIT License. See the LICENSE file for details.


package ipdata

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/eterline/micro-utils/internal/models"
)

type IpInfoExternalApi struct {
	client *http.Client
}

func NewExternalApi() *IpInfoExternalApi {
	return &IpInfoExternalApi{
		client: &http.Client{},
	}
}

func (ea *IpInfoExternalApi) GetData(ip net.IP) (models.AboutIPobject, error) {

	api := fmt.Sprintf("http://ip-api.com/json/%s?fields=66846719", ip.String())

	resp, err := ea.client.Get(api)
	if err != nil {
		return models.AboutIPobject{}, err
	}
	defer resp.Body.Close()

	var ipObj models.AboutIPobject

	if err := json.NewDecoder(resp.Body).Decode(&ipObj); err != nil {
		return models.AboutIPobject{}, err
	}

	ipObj.RequestTime = time.Now()

	return ipObj, nil
}
