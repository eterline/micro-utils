// Copyright (c) 2025 EterLine (Andrew)
// This file is part of micro-utils.
// Licensed under the MIT License. See the LICENSE file for details.

package smi

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	smiFormatValue = "--format=csv,noheader,nounits"
)

var (
	measureTrimRegexp = regexp.MustCompile(`^\s*([\d\.\-eE]+)\s*([a-zA-Z%]+)\s*$`)
)

type SmiQuery []string

func (sq SmiQuery) EncodeQuery() string {
	argList := strings.Join(sq, ",")
	return fmt.Sprintf("--query-gpu=%s", argList)
}

var (
	SmiIdentify = SmiQuery{
		"timestamp",
		"index",
		"name",
		"gpu_name",
		"uuid",
		"serial",
		"driver_version",
		"vbios_version",
		"pci.bus_id",
		"pci.domain",
		"pci.bus",
		"pci.device",
		"pci.device_id",
		"pci.sub_device_id",
		"pci.baseClass",
		"pci.subClass",
	}

	SmiStates = SmiQuery{
		"display_mode",
		"display_attached",
		"display_active",
		"persistence_mode",
		"accounting.mode",
		"accounting.buffer_size",
		"addressing_mode",
		"compute_mode",
		"compute_cap",
		"power.management",
		"mig.mode.current",
		"mig.mode.pending",
		"gsp.mode.current",
		"gsp.mode.default",
		"c2c.mode",
		"dramEncryption.mode.current",
		"dramEncryption.mode.pending",
		"ecc.mode.current",
		"ecc.mode.pending",
	}

	SmiTemperatures = SmiQuery{
		"temperature.gpu",
		"temperature.memory",
		"temperature.gpu.tlimit",
	}

	SmiPower = SmiQuery{
		"power.draw",
		"power.draw.average",
		"power.draw.instant",
		"power.limit",
		"enforced.power.limit",
		"power.default_limit",
		"power.min_limit",
		"power.max_limit",
		"module.power.draw.average",
		"module.power.draw.instant",
		"module.power.limit",
		"module.enforced.power.limit",
		"module.power.default_limit",
		"module.power.min_limit",
		"module.power.max_limit",
	}

	SmiUtilization = SmiQuery{
		"utilization.gpu",
		"utilization.memory",
		"utilization.encoder",
		"utilization.decoder",
		"utilization.jpeg",
		"utilization.ofa",
	}

	SmiMemory = SmiQuery{
		"memory.total",
		"memory.used",
		"memory.free",
		"memory.reserved",
		"protected_memory.total",
		"protected_memory.used",
		"protected_memory.free",
	}

	SmiPciBus = SmiQuery{
		"pcie.link.gen.gpucurrent",
		"pcie.link.gen.gpumax",
		"pcie.link.gen.hostmax",
		"pcie.link.gen.max",
		"pcie.link.width.current",
		"pcie.link.width.max",
	}

	SmiClocks = SmiQuery{
		"clocks.max.graphics",
		"clocks.max.memory",
		"clocks.max.sm",
		"clocks.default_applications.graphics",
		"clocks.default_applications.memory",
		"clocks.current.graphics",
		"clocks.current.memory",
		"clocks.current.sm",
		"clocks.current.video",
		"clocks.applications.graphics",
		"clocks.applications.memory",
	}
)

type SmiFetcher struct {
	bin string
}

func InitSMI() (*SmiFetcher, error) {
	s, err := exec.LookPath("nvidia-smi")
	if err != nil {
		return nil, fmt.Errorf("nvidia-smi binary searching error: %w", err)
	}

	f := &SmiFetcher{
		bin: s,
	}

	return f, nil
}

func (sf *SmiFetcher) Lookup(ctx context.Context, query SmiQuery) (map[string]string, error) {

	cmd := exec.CommandContext(
		ctx, sf.bin,
		query.EncodeQuery(),
		smiFormatValue,
	)

	data, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("nvidia-smi call query error: %w", err)
	}

	vs, err := getLines(data, query)
	if err != nil {
		return nil, fmt.Errorf("nvidia-smi invalid output: %w", err)
	}

	res := make(map[string]string, len(vs))
	for i, arg := range query {
		key := strings.ReplaceAll(arg, ".", "__")
		res[key] = cleanNA(vs[i])
	}

	return res, nil
}

func getLines(data []byte, query SmiQuery) (values []string, err error) {

	data = data[:len(data)-1]
	segments := bytes.Split(data, []byte(", "))
	lenQ := len(query)
	values = make([]string, lenQ)

	if lenQ == len(segments) {
		for i, seg := range segments {
			values[i] = string(seg)
		}

		return values, nil
	}

	return nil, errors.New("not consistent arguments and values")
}

func CamelFormat(v map[string]interface{}) map[string]interface{} {

	newV := make(map[string]interface{}, len(v))

	for keyI, valueI := range v {

		keyI = strings.ReplaceAll(keyI, "__", "_")
		splittedKey := strings.Split(keyI, "_")

		for indexJ, valueJ := range splittedKey {
			splittedKey[indexJ] = cases.Title(language.Und).String(valueJ)
		}

		keyI = strings.Join(splittedKey, "")
		newKeyI := firstToLower(keyI)

		newV[newKeyI] = valueI
	}

	return newV
}

func BuildTree(flat map[string]interface{}, splitKey string) map[string]interface{} {
	tree := make(map[string]interface{})

	for key, value := range flat {
		parts := strings.Split(key, splitKey)
		current := tree

		for i, part := range parts {
			if i == len(parts)-1 {
				// leaf node
				current[part] = value
			} else {
				// inner node
				if _, exists := current[part]; !exists {
					current[part] = make(map[string]interface{})
				}

				// assert to descend into map[string]interface{}
				if next, ok := current[part].(map[string]interface{}); ok {
					current = next
				} else {
					// overwrite conflicting leaf with map
					temp := make(map[string]interface{})
					current[part] = temp
					current = temp
				}
			}
		}
	}

	return tree
}

func cleanNA(s string) string {

	if s == "[N/A]" || s == "N/A" {
		return "null"
	}

	return s
}

func firstToLower(s string) string {
	r, size := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError && size <= 1 {
		return s
	}
	lc := unicode.ToLower(r)
	if r == lc {
		return s
	}
	return string(lc) + s[size:]
}
