package main

import (
	"fmt"
	"io"
	"os"

	microutils "github.com/eterline/micro-utils"
	"github.com/eterline/micro-utils/internal/config/cfgutil"
	"github.com/eterline/micro-utils/internal/config/ips2ubntest"
	"github.com/eterline/micro-utils/internal/services/netparse"
)

var (
	initArgs = cfgutil.UsualConfig[ips2ubntest.Configuration]{
		Config: &ips2ubntest.Configuration{
			Addrs:     []string{},
			OutFile:   "",
			InputFile: "",
			Separator: "\\n",
		},
		Name: "seeip",
	}
)

func main() {

	cfg, err := initArgs.ParseArgs()
	if err != nil {
		microutils.PrintFatalErr(err)
	}

	if err := validateConfig(cfg); err != nil {
		microutils.PrintFatalErr(err)
	}

	prs := netparse.NewNetParser()

	if err := prs.ParseAddrs(cfg.Addrs); err != nil {
		microutils.PrintErr(err)
	}

	if cfg.InputFile != "" {
		err, _ := prs.ParseFromFile(cfg.InputFile, cfg.Separator)
		if err != nil {
			microutils.PrintErr(err)
		}
	}

	output, err := selectOutput(cfg.OutFile)
	if err != nil {
		microutils.PrintFatalErr(err)
	}
	defer output.Close()

	subV4, subV6 := prs.Subnets()

	switch {
	case len(subV4) > 0:
		err := netparse.ExportSubnetsTo(subV4, output)
		if err != nil {
			microutils.PrintErr(err)
		}
	case len(subV6) > 0:
		err := netparse.ExportSubnetsTo(subV6, output)
		if err != nil {
			microutils.PrintErr(err)
		}
	}
}

func selectOutput(file string) (io.WriteCloser, error) {
	if file == "" {
		return os.Stdout, nil
	}

	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open output file: %s - %w", file, err)
	}

	return f, nil
}

func validateConfig(cfg ips2ubntest.Configuration) error {
	if len(cfg.Addrs) == 0 && cfg.InputFile == "" {
		return fmt.Errorf("no input addresses or input file provided")
	}
	return nil
}
