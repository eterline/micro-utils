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
	}
)

func main() {

	err := seeip.ParseArgs(&cfg)
	if err != nil {
		microutils.PrintFatalErr(err)
	}

	resolved, err := seeip.ResolveIP(cfg.Address...)
	if err != nil {
		microutils.PrintFatalErr(err)
	}

	infoSet := seeip.GetInfos(resolved)

	if cfg.JsonFormat {
		microutils.PrintJSON(cfg.Pretty, infoSet)
		os.Exit(0)
	}

	PrintSetText(infoSet)
}

func PrintSetText(infoSet seeip.ResolvedInfos) {
	b, _ := yaml.Marshal(infoSet)
	fmt.Println(string(b))
}
