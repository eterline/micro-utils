// Copyright (c) 2025 EterLine (Andrew)
// This file is part of micro-utils.
// Licensed under the MIT License. See the LICENSE file for details.


package main

import (
	"encoding/json"
	"log/slog"
	"net/http"

	microutils "github.com/eterline/micro-utils"
)

type (
	StubCfg struct {
		Listen string `arg:"-l,--listen" help:"Listen connection addr"`
	}

	ResponseBody struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
)

var (
	cfg = &microutils.UsualConfig[StubCfg]{
		Config: StubCfg{
			Listen: ":80",
		},
	}
)

func runServer(c *StubCfg) error {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {

		res := ResponseBody{
			Code:    http.StatusOK,
			Message: "OK",
		}

		w.Header().Add("Content-type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(res)
	})

	err := http.ListenAndServe(c.Listen, nil)
	switch err {
	case http.ErrServerClosed:
		return nil
	case nil:
		return nil
	}

	return err
}

func main() {

	slog.Info("start service")
	defer slog.Info("exiting service")

	argConf, err := cfg.ParseArgs("stubit")
	if err != nil {
		slog.Error("arg parse error", "error", err.Error())
		return
	}

	slog.Info("http wed server started", "listen", argConf.Listen)

	if err := runServer(argConf); err != nil {
		slog.Error("server error", "error", err.Error())
		return
	}
}
