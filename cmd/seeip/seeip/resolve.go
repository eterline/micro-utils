package seeip

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	microutils "github.com/eterline/micro-utils"
)

type Resolvs map[string][]net.IP

func (r Resolvs) add(key string, addr net.IP) {
	_, ok := r[key]
	if !ok {
		r[key] = make([]net.IP, 0)
	}
	r[key] = append(r[key], addr)
}

func (r Resolvs) isOk() bool {
	for _, v := range r {
		if v != nil && len(v) > 0 {
			return true
		}
	}
	return false
}

func ResolveIP(hosts ...string) (Resolvs, error) {
	resolv := Resolvs{}

	for _, addr := range hosts {
		ip := net.ParseIP(addr)
		if ip != nil {
			resolv.add(addr, ip)
			continue
		}

		ips, err := net.LookupIP(addr)
		if err != nil {
			continue
		}

		for _, ip := range ips {
			resolv.add(addr, ip)
		}
	}

	if resolv.isOk() {
		return resolv, nil
	}

	return nil, fmt.Errorf("no any hosts resolved")
}

type ResolvedInfos map[string][]InfoIP

func (r ResolvedInfos) add(key string, info InfoIP) {
	_, ok := r[key]
	if !ok {
		r[key] = make([]InfoIP, 0)
	}
	r[key] = append(r[key], info)
}

func GetInfos(resolved Resolvs) ResolvedInfos {

	var (
		data = ResolvedInfos{}
		mu   = sync.Mutex{}
		wg   = sync.WaitGroup{}
	)

	for hostname, ips := range resolved {
		for _, ip := range ips {
			time.Sleep(100 * time.Millisecond)
			wg.Add(1)
			ip := ip

			go func() {
				defer wg.Done()

				info, err := GetIpInfoWithContext(context.Background(), ip)
				if err != nil {
					microutils.PrintErr(err)
					return
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
