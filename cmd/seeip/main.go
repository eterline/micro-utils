// Copyright (c) 2025 EterLine (Andrew)
// This file is part of micro-utils.
// Licensed under the MIT License. See the LICENSE file for details.

package main

import (
	"context"
	"errors"

	microutils "github.com/eterline/micro-utils"
	ipDataAdapters "github.com/eterline/micro-utils/internal/adapters/ipdata"
	configSeeip "github.com/eterline/micro-utils/internal/config/seeip"
	ipDataService "github.com/eterline/micro-utils/internal/services/ipdata"

	"github.com/eterline/micro-utils/internal/config/cfgutil"
	"github.com/eterline/micro-utils/internal/models"
)

var (
	initArgs = cfgutil.UsualConfig[configSeeip.Configuration]{
		Config: &configSeeip.Configuration{
			Address:         []string{},
			IsJson:          false,
			Pretty:          false,
			ResolverService: "local",
		},
		Name: "seeip",
	}
)

func main() {

	cfg, err := initArgs.ParseArgs()
	if err != nil {
		microutils.PrintFatalErr(err)
	}

	// TODO: correct sql request
	// storeIP, err := ipDataAdapters.NewIpInfoSqlite(context.Background(), "./seeip.db")
	// if err != nil {
	// 	microutils.PrintFatalErr(err)
	// }

	rslv, err := selectResolver(cfg.ResolverService)
	if err != nil {
		microutils.PrintFatalErr(err)
	}

	resumer := ipDataAdapters.NewExternalApi()
	scr := ipDataService.NewNetworkScrapeService(cfg.Workers, rslv, resumer, nil)

	resolvs, err := scr.ResolveDNS(context.Background(), cfg.Address)
	if err != nil {
		microutils.PrintFatalErr(err)
	}

	d := []models.ResumeAboutIP{}

	for _, abouts := range resolvs {
		resume, err := scr.FetchAboutIP(abouts.IPs)
		if err != nil {
			microutils.PrintErr(err)
		}
		d = append(d, resume...)
	}

	resulted := ipDataAdapters.SortResolvedAndResume(resolvs, d)
	if cfg.IsJson {
		microutils.PrintJSON(cfg.Pretty, resulted)
	} else {
		microutils.PrintYaml(resulted)
	}
}

func selectResolver(name string) (models.Resolver, error) {
	switch {

	case name == "cloudflare":
		return ipDataAdapters.NewCloudflareResolver(), nil

	case name == "google":
		return ipDataAdapters.NewGoogleResolver(), nil

	case name == "local":
		return ipDataAdapters.NewLocalResolver(), nil

	case microutils.IsAddressString(name):
		return ipDataAdapters.NewRemoteResolver(name), nil

	default:
		return nil, errors.New("unknown DNS resolver name")
	}
}
