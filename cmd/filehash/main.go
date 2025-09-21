// Copyright (c) 2025 EterLine (Andrew)
// This file is part of micro-utils.
// Licensed under the MIT License. See the LICENSE file for details.

package main

import (
	"fmt"

	microutils "github.com/eterline/micro-utils"
	"github.com/eterline/micro-utils/cmd/filehash/filehash"
)

const (
	tableTmp = `===========================
 SHA256 = %s
 SHA1   = %s
 MD5    = %s
`
)

var (
	cfg = filehash.Configuration{
		File:       "",
		JsonFormat: false,
		Pretty:     false,
	}
)

func main() {

	res := &filehash.FileData{}

	err := filehash.ParseArgs(&cfg)
	if err != nil {
		microutils.PrintFatalErr(err)
	}

	if microutils.IsInputFromPipe() {

		if !cfg.JsonFormat {
			fmt.Println("Target file: stdin")
		}

		res, err = filehash.CalcPipelineHash()
		if err != nil {
			microutils.PrintFatalErr(err)
		}

	} else {

		if !cfg.JsonFormat {
			fmt.Printf("Target file: %s\n", cfg.File)
		}

		res, err = filehash.CalcFileHash(cfg.File)
		if err != nil {
			microutils.PrintFatalErr(err)
		}

	}

	if cfg.JsonFormat {
		if err := microutils.PrintJSON(cfg.Pretty, res); err != nil {
			microutils.PrintFatalErr(err)
		}
		return
	}

	fmt.Printf("File size: %s\n", microutils.BytesToSizeString(res.Size))
	fmt.Printf(tableTmp, res.SHA256, res.SHA1, res.MD5)
}
