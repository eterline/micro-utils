// Copyright (c) 2025 EterLine (Andrew)
// This file is part of My-Go-Project.
// Licensed under the MIT License. See the LICENSE file for details.


package seeip

import (
	"errors"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	parallelThreads = 50
)

func ParsePorts(s string) ([]NetPort, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, errors.New("empty port string")
	}

	subs := strings.Split(s, ",")
	ports := []NetPort{}
	seen := make(map[NetPort]bool)

	for _, sub := range subs {
		sub = strings.TrimSpace(sub)
		if sub == "" {
			continue
		}

		if strings.Contains(sub, "-") {
			rangeParts := strings.SplitN(sub, "-", 2)
			if len(rangeParts) != 2 {
				continue
			}
			start, ok1 := StringToNetPort(rangeParts[0])
			end, ok2 := StringToNetPort(rangeParts[1])
			if !ok1 || !ok2 || start > end {
				continue
			}
			for p := start; p <= end; p++ {
				if !seen[p] {
					ports = append(ports, p)
					seen[p] = true
				}
			}
		} else {
			p, ok := StringToNetPort(sub)
			if ok && !seen[p] {
				ports = append(ports, p)
				seen[p] = true
			}
		}
	}

	if len(ports) == 0 {
		return nil, errors.New("no valid ports found")
	}
	return ports, nil
}

type NetPort uint16

func StringToNetPort(s string) (NetPort, bool) {
	v, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return 0, false
	}

	if v < 1 || v > 65535 {
		return 0, false
	}

	return NetPort(v), true
}

func (p NetPort) Address(ip net.IP) string {
	return net.JoinHostPort(ip.String(), strconv.Itoa(int(p)))
}

type ScanMap map[NetPort]bool

func (m ScanMap) Allowed() []NetPort {
	r := []NetPort{}
	for port, ok := range m {
		if ok {
			r = append(r, port)
		}
	}
	return r
}

func (m ScanMap) Closed() []NetPort {
	r := []NetPort{}
	for port, ok := range m {
		if !ok {
			r = append(r, port)
		}
	}
	return r
}

func ScanTCP(ip net.IP, ports []NetPort) ScanMap {

	var (
		numPorts = len(ports)
		data     = make(ScanMap, numPorts)
		wg       = &sync.WaitGroup{}
		mu       = sync.Mutex{}
	)

	allowChan := make(chan struct{}, parallelThreads)
	defer close(allowChan)

	for range parallelThreads {
		allowChan <- struct{}{}
	}

	for _, port := range ports {

		<-allowChan
		wg.Add(1)

		go func(p NetPort) {
			defer func() {
				allowChan <- struct{}{}
				wg.Done()
			}()

			addr := p.Address(ip)
			conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
			if conn != nil {
				conn.Close()
			}

			mu.Lock()
			data[p] = err == nil
			mu.Unlock()
		}(port)
	}

	wg.Wait()
	return data
}
