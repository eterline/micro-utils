// Copyright (c) 2025 EterLine (Andrew)
// This file is part of My-Go-Project.
// Licensed under the MIT License. See the LICENSE file for details.


package main

import (
	"fmt"
	"os"

	microutils "github.com/eterline/micro-utils"
	"github.com/eterline/micro-utils/cmd/seeip/seeip"
	"gopkg.in/yaml.v3"
)

var (
	cfg = seeip.Configuration{
		Address:    []string{},
		MapUrl:     false,
		JsonFormat: false,
		Pretty:     false,
		GeoStamp:   false,
	}
)

func main() {

	err := seeip.ParseArgs(&cfg)
	if err != nil {
		microutils.PrintFatalErr(err)
	}

	resolved, err := seeip.ResolveHost(cfg.Address...)
	if err != nil {
		microutils.PrintFatalErr(err)
	}

	infoSet := seeip.GetInfos(resolved, cfg.GeoStamp)

	if cfg.JsonFormat {
		microutils.PrintJSON(cfg.Pretty, infoSet)
		os.Exit(0)
	}

	PrintSetText(infoSet)
}

func ResolveHostList(hostnames ...string) []seeip.LookupTable {

	set := []seeip.LookupTable{}

	for _, host := range hostnames {
		lookups, err := seeip.ResolveHost(host)
		if err != nil {
			microutils.PrintErr(err)
			continue
		}
		set = append(set, lookups...)
	}

	return set
}

func PrintSetText(infoSet seeip.ResolvedInfos) {
	b, _ := yaml.Marshal(infoSet)
	fmt.Println(string(b))
}
