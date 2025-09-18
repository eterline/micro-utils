// Copyright (c) 2025 EterLine (Andrew)
// This file is part of micro-utils.
// Licensed under the MIT License. See the LICENSE file for details.


package main

import (
	"bytes"
	"errors"
	"log/slog"
	"net"

	"github.com/eterline/micro-utils/cmd/minescan/conf"
	"github.com/google/uuid"
)

var (
	cfg = conf.Configuration{
		Listen: ":25565",
		FilterWords: []string{
			"hi honeypots",
		},
	}
)

func main() {

	err := conf.ParseArgs(&cfg)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	filter, err := NewPayloadFilter(cfg.FilterWords...)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	log := slog.With("listen", cfg.Listen)
	log.Info("starting scanner detector")

	lst, err := net.Listen("tcp", cfg.Listen)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer lst.Close()

	for {
		conn, err := lst.Accept()
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		go func() {

			host, port, err := net.SplitHostPort(conn.RemoteAddr().String())
			if err != nil {
				slog.Error(err.Error())
				return
			}

			log := log.With(
				"conn_id", uuid.New(),
				"request_ip", host,
				"request_port", port,
			)

			log.Info("new connection")
			defer log.Info("connection closed")

			ok, err := process(conn, filter)
			if err != nil {
				log.Error(err.Error())
			}

			if ok {
				log.Info("honeypot scan detected")
			}
		}()
	}
}

type Matcher interface {
	Match(p []byte) bool
}

func process(c net.Conn, m Matcher) (bool, error) {

	buf := make([]byte, 512)
	defer c.Close()

	n, err := c.Read(buf)
	if err != nil {
		return false, err
	}

	if !m.Match(buf[:n]) {
		return false, nil
	}

	s := []byte("do not scan me, bro")
	if _, err := c.Write(s); err != nil {
		return false, err
	}

	return true, nil
}

type PayloadFilter struct {
	triggers [][]byte
}

func NewPayloadFilter(w ...string) (*PayloadFilter, error) {

	if w == nil || len(w) < 1 {
		return nil, errors.New("invalid filter inputs")
	}

	trs := make([][]byte, len(w))
	for i, words := range w {
		trs[i] = []byte(words)
	}

	return &PayloadFilter{
		triggers: trs,
	}, nil
}

func (pf *PayloadFilter) Match(p []byte) bool {
	for _, value := range pf.triggers {
		if bytes.Contains(p, value) {
			return true
		}
	}

	return false
}
