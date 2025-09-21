// Copyright (c) 2025 EterLine (Andrew)
// This file is part of micro-utils.
// Licensed under the MIT License. See the LICENSE file for details.

package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"sync"

	microutils "github.com/eterline/micro-utils"
	"github.com/eterline/micro-utils/cmd/opnsense-filter-log-scan/conf"
	"gopkg.in/yaml.v3"
)

var (
	logNameReg = regexp.MustCompile(`^(filter_\d{8}|latest)\.log$`)
	cfg        = conf.Configuration{
		LogsDir:    "./",
		JsonFormat: false,
		Pretty:     false,
	}
)

func main() {

	jr := NewLogJournal()

	err := conf.ParseArgs(&cfg)
	if err != nil {
		microutils.PrintFatalErr(err)
	}

	entries, err := os.ReadDir(cfg.LogsDir)
	if err != nil {
		microutils.PrintFatalErr(err)
	}

	for _, e := range entries {

		if !logNameReg.MatchString(e.Name()) {
			continue
		}

		dst := filepath.Join(cfg.LogsDir, e.Name())

		fmt.Printf("open: %s\n", dst)

		if e.IsDir() {
			err := fmt.Errorf("path %s is not file: %w", dst, errors.ErrUnsupported)
			microutils.PrintErr(err)
			continue
		}

		f, err := os.Open(dst)
		if err != nil {
			microutils.PrintErr(err)
			continue
		}

		func() {
			defer f.Close()

			sc := bufio.NewScanner(f)

			for sc.Scan() {

				logLine := sc.Bytes()
				if err := jr.Add(logLine); err != nil {
					slog.Error(err.Error())
				}
			}
		}()
	}

	d, err := yaml.Marshal(jr.DataSortMore())
	if err != nil {
		slog.Error(err.Error())
		return
	}

	err = os.WriteFile("addrs.yaml", d, 0664)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

type (
	IpStats struct {
		Counter   int64
		Protocols map[string]int64
		Actions   map[string]int64
	}

	LogJournal struct {
		mu   sync.Mutex
		data map[string]IpStats
	}
)

func NewLogJournal() *LogJournal {
	return &LogJournal{
		data: make(map[string]IpStats),
	}
}

func (lj *LogJournal) Add(line []byte) error {
	l, err := ParseLogFilterLine(line)
	if err != nil {
		return err
	}

	key := l.SrcIP.String()
	proto := l.Protocol.String()
	action := l.Action.String()

	lj.mu.Lock()
	defer lj.mu.Unlock()

	stats := lj.data[key]
	stats.Counter++

	if stats.Protocols == nil {
		stats.Protocols = make(map[string]int64)
	}
	stats.Protocols[proto]++

	if stats.Actions == nil {
		stats.Actions = make(map[string]int64)
	}
	stats.Actions[action]++

	lj.data[key] = stats

	return nil
}

type IpSrcCounter struct {
	IP    string
	Stats IpStats
}

type IpCounterSlice []IpSrcCounter

func (a IpCounterSlice) Len() int      { return len(a) }
func (a IpCounterSlice) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a IpCounterSlice) Less(i, j int) bool {
	return a[i].Stats.Actions[Blocked.String()] < a[j].Stats.Actions[Blocked.String()]
}
func (a IpCounterSlice) More(i, j int) bool {
	return a[i].Stats.Actions[Blocked.String()] > a[j].Stats.Actions[Blocked.String()]
}

func (lj *LogJournal) Data() IpCounterSlice {

	list := make([]IpSrcCounter, 0, len(lj.data))

	for ip, stats := range lj.data {

		cnt := IpSrcCounter{
			IP:    ip,
			Stats: stats,
		}

		list = append(list, cnt)
	}

	return list
}

func (lj *LogJournal) DataSortLess() IpCounterSlice {
	d := lj.Data()
	sort.Slice(d, d.Less)
	return d
}

func (lj *LogJournal) DataSortMore() IpCounterSlice {
	d := lj.Data()
	sort.Slice(d, d.More)
	return d
}

type FilterAction int32

const (
	_ FilterAction = iota
	Passed
	Blocked
	Redirect
)

func (fa FilterAction) String() (s string) {
	switch fa {
	case Passed:
		s = "pass"
	case Blocked:
		s = "block"
	case Redirect:
		s = "rdr"
	default:
		s = fmt.Sprintf("unknown filter action: %d", fa)
	}

	return
}

func parseFilterAction(p []byte) (fa FilterAction, err error) {

	switch string(p) {
	case "pass":
		fa = Passed
	case "block":
		fa = Blocked
	case "rdr":
		fa = Redirect
	default:
		return 0, fmt.Errorf("unknown filter action: %s", p)
	}

	return fa, nil
}

type FilterDirection int32

const (
	_ FilterDirection = iota
	Input
	Output
)

func parseFilterDirection(p []byte) (fd FilterDirection, err error) {

	switch string(p) {
	case "in":
		fd = Input
	case "out":
		fd = Output
	default:
		return 0, fmt.Errorf("unknown filter direction: %s", p)
	}

	return fd, nil
}

type ProtoType int32

func (pt ProtoType) PortRelative() bool {
	switch pt {
	case 1, 2, 4, 41, 47, 50:
		return false
	default:
		return true
	}
}

func (pt ProtoType) String() string {
	switch pt {
	case 0:
		return "hopopt"
	case 1:
		return "icmp"
	case 2:
		return "igmp"
	case 3:
		return "ggp"
	case 4:
		return "ip-in-ip"
	case 5:
		return "st"
	case 6:
		return "tcp"
	case 7:
		return "cbt"
	case 8:
		return "egp"
	case 9:
		return "igp"
	case 17:
		return "udp"
	case 50:
		return "esp"
	case 51:
		return "ah"
	case 132:
		return "sctp"
	case 136:
		return "udplite"
	case 137:
		return "mpls-in-ip"
	case 138:
		return "manet"
	case 253, 254:
		return "experimental"
	case 255:
		return "reserved"
	default:
		return fmt.Sprintf("unknown protocol (%d)", pt)
	}
}

func parseProtoType(p []byte) (pt ProtoType, err error) {
	val, err := strconv.Atoi(string(p))
	if err != nil {
		return 0, err
	}

	return ProtoType(val), nil
}

type Port int32

func parsePort(p []byte) (pt Port, err error) {
	val, err := strconv.Atoi(string(p))
	if err != nil {
		return 0, err
	}

	if val <= 0 || val > 65535 {
		return 0, errors.New("invalid port value")
	}

	return Port(val), nil
}

type VersionIP int32

func parseVersionIP(p []byte) (pt VersionIP, err error) {
	val, _ := strconv.Atoi(string(p))

	switch val {
	case 4, 6:
		return VersionIP(val), nil
	default:
		return 0, errors.New("invalid ip version")
	}
}

func parseRuleID(p []byte) (id int32, err error) {
	val, err := strconv.Atoi(string(p))
	if err != nil {
		return 0, err
	}

	return int32(val), nil
}

type LogFilterLine struct {
	RuleID         int32
	UUID           string
	Interface      string
	Match          bool
	Action         FilterAction
	Direction      FilterDirection
	VerIP          VersionIP
	ServiceType    string
	TTL            int32
	IPID           int32
	FragmentOffset int64
	FragFlags      string
	Protocol       ProtoType
	Len            int32
	SrcIP          net.IP
	DstIP          net.IP
	SrcPort        Port
	DstPort        Port
}

func ParseLogFilterLine(line []byte) (l LogFilterLine, err error) {

	segments := bytes.Split(line, []byte("\"] "))
	if len(segments) != 2 {
		return l, errors.New("invalid log line format")
	}

	fields := bytes.Split(segments[1], []byte{','})

	l.RuleID, err = parseRuleID(fields[0])
	if err != nil {
		return LogFilterLine{}, err
	}

	l.UUID = string(fields[3])

	l.Interface = string(fields[4])
	l.Match = (string(fields[5]) == "match")

	l.Action, err = parseFilterAction(fields[6])
	if err != nil {
		return LogFilterLine{}, err
	}

	l.Direction, err = parseFilterDirection(fields[7])
	if err != nil {
		return LogFilterLine{}, err
	}

	l.VerIP, err = parseVersionIP(fields[8])
	if err != nil {
		return LogFilterLine{}, err
	}

	l.ServiceType = string(fields[9])

	ttl, err := strconv.Atoi(string(fields[11]))
	if err != nil {
		return LogFilterLine{}, err
	}
	l.TTL = int32(ttl)

	ipid, err := strconv.Atoi(string(fields[12]))
	if err != nil {
		return LogFilterLine{}, err
	}
	l.IPID = int32(ipid)

	l.FragmentOffset, err = strconv.ParseInt(string(fields[13]), 10, 64)
	if err != nil {
		return LogFilterLine{}, err
	}

	l.FragFlags = string(fields[14])

	l.Protocol, err = parseProtoType(fields[15])
	if err != nil {
		return LogFilterLine{}, err
	}

	len, err := strconv.Atoi(string(fields[17]))
	if err != nil {
		return LogFilterLine{}, err
	}
	l.Len = int32(len)

	ip := net.ParseIP(string(fields[18]))
	if ip == nil {
		return LogFilterLine{}, errors.New("invalid source ip")
	}
	l.SrcIP = ip

	ip = net.ParseIP(string(fields[19]))
	if ip == nil {
		return LogFilterLine{}, errors.New("invalid destination ip")
	}
	l.DstIP = ip

	if l.Protocol.PortRelative() {

		l.SrcPort, err = parsePort(fields[20])
		if err != nil {
			return LogFilterLine{}, errors.New("invalid source port")
		}

		l.DstPort, err = parsePort(fields[21])
		if err != nil {
			return LogFilterLine{}, errors.New("invalid destination port")
		}

	}

	return l, nil
}
