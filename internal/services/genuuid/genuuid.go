// Copyright (c) 2025 EterLine (Andrew)
// This file is part of micro-utils.
// Licensed under the MIT License. See the LICENSE file for details.

package genuuid

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrCount         = errors.New("UUID count must be positive and above 0")
	ErrPayloadEmpty  = errors.New("hash payload must be not empty")
	ErrUnknownDomain = errors.New("unknown UUID domain area")
	ErrUnknownVer    = errors.New("unsupported UUID version for name-based generation")
)

func getNamespace(d string) (uuid.UUID, bool) {

	switch d = strings.ToLower(d); d {
	case "dns":
		return uuid.NameSpaceDNS, true
	case "url":
		return uuid.NameSpaceURL, true
	case "oid":
		return uuid.NameSpaceOID, true
	case "x500":
		return uuid.NameSpaceX500, true
	}

	return uuid.Nil, false
}

type GenerateUUID struct{}

func (gd GenerateUUID) GenerateOf(n int) ([]uuid.UUID, error) {
	if n <= 0 {
		return nil, ErrCount
	}

	uuidList := make([]uuid.UUID, n)
	for i := 0; i < n; i++ {
		v, err := uuid.NewRandom()
		if err != nil {
			return nil, fmt.Errorf("failed to generate UUID: %w", err)
		}
		uuidList[i] = v
	}

	return uuidList, nil
}

func (gd GenerateUUID) GenerateFrom(domain, payload string, version int) ([]uuid.UUID, error) {
	if payload == "" {
		return nil, ErrPayloadEmpty
	}

	if version != 3 && version != 5 {
		return nil, ErrUnknownVer
	}

	ns, ok := getNamespace(domain)
	if !ok {
		return nil, ErrUnknownDomain
	}

	value := uuid.NewMD5(ns, []byte(payload))

	return []uuid.UUID{value}, nil

}
