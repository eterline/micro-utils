// Copyright (c) 2025 EterLine (Andrew)
// This file is part of micro-utils.
// Licensed under the MIT License. See the LICENSE file for details.

package cli

import (
	"context"
	"errors"
	"fmt"

	microutils "github.com/eterline/micro-utils"
	"github.com/eterline/micro-utils/cmd/gpufo/conf"
	"github.com/eterline/micro-utils/cmd/gpufo/external/smi"
)

func HandleCLI(ctx context.Context, c conf.Configuration) error {

	q, err := parseTables(c.SmiTables)
	if err != nil {
		return err
	}

	smiUse, err := smi.InitSMI()
	if err != nil {
		return err
	}

	smiTable := make(map[string]interface{})

	for _, query := range q {
		res, err := smiUse.Lookup(ctx, query)
		if err == nil {
			for key, value := range res {
				smiTable[key] = value
			}
		}
	}

	if c.Tree && !c.Camel {
		smiTable = smi.BuildTree(smiTable, "__")
	}

	if !c.Tree && c.Camel {
		smiTable = smi.CamelFormat(smiTable)
	}

	if c.JsonFormat {
		return microutils.PrintJSON(c.Pretty, smiTable)
	}

	return microutils.PrintYaml(smiTable)
}

func parseTables(t []string) (q []smi.SmiQuery, err error) {

	q = make([]smi.SmiQuery, len(t))

	for i, n := range t {
		switch n {
		case "identify":
			q[i] = smi.SmiIdentify
		case "clocks":
			q[i] = smi.SmiClocks
		case "memory":
			q[i] = smi.SmiMemory
		case "pci":
			q[i] = smi.SmiPciBus
		case "power":
			q[i] = smi.SmiPower
		case "states":
			q[i] = smi.SmiStates
		case "temp":
			q[i] = smi.SmiTemperatures
		case "utilization":
			q[i] = smi.SmiUtilization
		default:
			return nil, fmt.Errorf("invalid table: %w", errors.ErrUnsupported)
		}
	}

	return q, nil
}
