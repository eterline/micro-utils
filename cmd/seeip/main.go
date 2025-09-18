// Copyright (c) 2025 EterLine (Andrew)
// This file is part of micro-utils.
// Licensed under the MIT License. See the LICENSE file for details.


package main

import (
	"context"
	"errors"

	microutils "github.com/eterline/micro-utils"
	ipDataAdapters "github.com/eterline/micro-utils/internal/adapters/ipdata"
	"github.com/eterline/micro-utils/internal/config/cfgutil"
	"github.com/eterline/micro-utils/internal/config/seeip"
	"github.com/eterline/micro-utils/internal/models"
	ipDataService "github.com/eterline/micro-utils/internal/services/ipdata"
)

var (
	initArgs = cfgutil.UsualConfig[seeip.Configuration]{
		Config: &seeip.Configuration{
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

	rslv, err := selectResolver(cfg.ResolverService)
	if err != nil {
		microutils.PrintFatalErr(err)
	}

	scr := ipDataService.NewNetworkScrapeService(rslv, cfg.Workers)
	about := scr.ResolveDNS(context.Background(), cfg.Address)

	if cfg.IsJson {
		microutils.PrintJSON(cfg.Pretty, about)
		return
	}

	microutils.PrintYaml(about)
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
