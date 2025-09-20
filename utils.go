// Copyright (c) 2025 EterLine (Andrew)
// This file is part of micro-utils.
// Licensed under the MIT License. See the LICENSE file for details.

package microutils

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"regexp"
	"runtime"
	"strconv"

	"golang.org/x/exp/constraints"
	"gopkg.in/yaml.v3"
)

var (
	cpuCountAtStart = 0
)

func init() {
	cpuCountAtStart = runtime.NumCPU()
}

const (
	prefixB  = "Bytes"
	prefixKB = "KBytes"
	prefixMB = "MBytes"
	prefixGB = "GBytes"
	prefixTB = "TBytes"
)

type FileInfo struct {
	Name string
	Path string
	Data []byte
}

func (f *FileInfo) Size() int {
	if f == nil {
		return 0
	}
	return len(f.Data)
}

func BytesToSizeString(s int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
		TB = 1024 * GB
	)

	size := float64(s)

	switch {
	case s < KB:
		return fmt.Sprintf("%.4f B", size)
	case s < MB:
		return fmt.Sprintf("%.4f KB", size/KB)
	case s < GB:
		return fmt.Sprintf("%.4f MB", size/MB)
	case s < TB:
		return fmt.Sprintf("%.4f GB", size/GB)
	default:
		return fmt.Sprintf("%.4f TB", size/TB)
	}
}

func PrintFatalErr(err error) {
	fmt.Printf("ERROR: %v\n", err)
	os.Exit(1)
}

func PrintErr(err error) {
	fmt.Printf("ERROR: %v\n", err)
}

func PrintJSON(pretty bool, v any) error {
	if pretty {
		return jsonPrintPretty(v)
	}
	return jsonPrint(v)
}

func jsonPrint(v any) error {
	return json.NewEncoder(os.Stdout).Encode(v)
}

func jsonPrintPretty(v any) error {
	data, err := json.MarshalIndent(v, " ", "  ")
	if err == nil {
		fmt.Println(string(data))
	}
	return err
}

func IsInputFromPipe() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	return (fi.Mode() & os.ModeCharDevice) == 0
}

func PrintYaml(v any) (err error) {
	data, err := yaml.Marshal(v)
	if err != nil {
		return err
	}

	_, err = fmt.Print(string(data))
	return
}

func Clamp[T constraints.Ordered](x, min, max T) T {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

/*
	 InitWorkersCount - get maximum gorutines workers count in one time.
		Calculates from 1 as minimum, `i` as wanted and runtime.NumCPU() as maximum from at starting program moment.
*/
func InitWorkersCount(i int) int {
	return Clamp(i, 1, cpuCountAtStart)
}

/*
	 InitWorkersCountCurrently - get maximum gorutines workers count in one time.
		Calculates from 1 as minimum, `i` as wanted and runtime.NumCPU() as maximum from actual time.
*/
func InitWorkersCountCurrently(i int) int {
	return Clamp(i, 1, runtime.NumCPU())
}

var domainRegex = regexp.MustCompile(`^(?i:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?)(?:\.(?i:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?))*$`)

func IsAddressString(s string) bool {
	host := s

	h, port, err := net.SplitHostPort(s)
	if err == nil {
		host = h
	}

	if port != "" {
		p, err := strconv.ParseInt(port, 10, 16)
		if err != nil {
			return false
		}

		if p < 0 || p > 65_535 {
			return false
		}
	}

	if ip := net.ParseIP(host); ip != nil {
		if ip.To4() != nil {
			return true
		}
		return true
	}

	if len(host) > 0 && len(host) <= 253 && domainRegex.MatchString(host) {
		return true
	}

	return false
}

type TicketPool struct {
	ticketCh chan struct{}
}

func NewTicketPool(workers int) *TicketPool {
	if workers < 0 {
		workers = 1
	}

	return &TicketPool{
		ticketCh: make(chan struct{}, workers),
	}
}

func (tp *TicketPool) PutTicket() {
	if tp.ticketCh == nil {
		panic("use of closed workers ticket pool")
	}
	<-tp.ticketCh
}

func (tp *TicketPool) CatchTicket() {
	if tp.ticketCh == nil {
		panic("use of closed workers ticket pool")
	}
	tp.ticketCh <- struct{}{}
}

func (tp *TicketPool) ClosePool() {
	close(tp.ticketCh)
}
