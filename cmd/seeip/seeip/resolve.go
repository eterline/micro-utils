// Copyright (c) 2025 EterLine (Andrew)
// This file is part of My-Go-Project.
// Licensed under the MIT License. See the LICENSE file for details.


package seeip

import (
	"context"
	"net"
	"sync"

	microutils "github.com/eterline/micro-utils"
)

type NsLookup struct {
	IP       net.IP   `json:"ip" yaml:"ip"`
	Hostname []string `json:"hostname" yaml:"hostname"`
}

func lookupTableFromIp(ip net.IP) NsLookup {

	res := NsLookup{
		IP: ip,
	}

	names, err := net.LookupAddr(ip.String())
	if err == nil {
		res.Hostname = names
	}

	return res
}

func ResolveHost(host string) ([]NsLookup, error) {

	if ip := net.ParseIP(host); ip != nil {
		return []NsLookup{lookupTableFromIp(ip)}, nil
	}

	ips, err := net.LookupIP(host)
	if err != nil {
		return nil, err
	}

	var (
		tableList = []NsLookup{}
		wg        = &sync.WaitGroup{}
		mu        = sync.Mutex{}
	)

	for _, ip := range ips {
		wg.Add(1)

		go func(ipAddr net.IP) {
			defer wg.Done()

			hostNames := []string{host}
			if n, err := net.LookupAddr(ipAddr.String()); err == nil {
				hostNames = append(hostNames, n...)
			}

			table := NsLookup{
				IP:       ipAddr,
				Hostname: hostNames,
			}

			mu.Lock()
			tableList = append(tableList, table)
			mu.Unlock()
		}(ip)
	}

	wg.Wait()
	return tableList, nil
}

func (r Resolvs) isOk() bool {
	for _, v := range r {
		if v != nil && len(v) > 0 {
			return true
		}
	}
	return false
}

func GetInfos(resolved Resolvs, geoStamp bool) ResolvedInfos {

	var (
		data = ResolvedInfos{}
		mu   = sync.Mutex{}
		wg   = sync.WaitGroup{}
	)

	for hostname, ips := range resolved {
		for _, ip := range ips {

			wg.Add(1)
			ip := ip

			go func() {
				defer wg.Done()

				info, err := GetIpInfoWithContext(context.Background(), ip)
				if err != nil {
					microutils.PrintErr(err)
					return
				}

				if geoStamp {
					info.AddGeo()
				}

				mu.Lock()
				data.add(hostname, info)
				mu.Unlock()
			}()
		}
	}

	wg.Wait()
	return data
}
