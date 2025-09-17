// Copyright (c) 2025 EterLine (Andrew)
// This file is part of My-Go-Project.
// Licensed under the MIT License. See the LICENSE file for details.

package ipdata

import (
	"context"
	"net"
	"regexp"
	"strings"
	"sync"

	"github.com/eterline/micro-utils/internal/models"
)

var (
	ipRegex = regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
)

/*
ResolverService - implements DNS resolving logic

	If resolver got nil, will use system DNS resolving
*/
type ResolverService struct {
	provider models.ProviderDoH
}

/*
NewResolverService - init resolver service

	If resolver got nil, will use system DNS resolving
*/
func NewResolverService(r models.ProviderDoH) *ResolverService {
	return &ResolverService{
		provider: r,
	}
}

func isIP(s string) bool {
	return ipRegex.MatchString(strings.TrimSpace(s))
}

// ResolveAs - resolving of A and AAAA records in DNS
func (rs *ResolverService) ResolveAs(names []string) map[string][]net.IP {
	namePool := map[string][]net.IP{}
	mu := sync.Mutex{}
	wg := &sync.WaitGroup{}

	for _, name := range names {

		wg.Go(func() {

			if isIP(name) {
				mu.Lock()
				namePool[name] = []net.IP{net.ParseIP(name)}
				mu.Unlock()
				return
			}

			if rs.provider == nil {
				ips, _ := net.LookupIP(name)
				mu.Lock()
				namePool[name] = ips
				mu.Unlock()
				return
			}

			ips := []net.IP{}
			domain := models.Domain(name)

			aRecords, err := rs.provider.Query(context.TODO(), domain, models.TypeA)
			if err == nil {
				for _, ans := range aRecords.Answer {
					ips = append(ips, net.IP(ans.Data))
				}
			}

			aaaaRecords, err := rs.provider.Query(context.TODO(), domain, models.TypeAAAA)
			if err == nil {
				for _, ans := range aaaaRecords.Answer {
					ips = append(ips, net.IP(ans.Data))
				}
			}

			mu.Lock()
			namePool[name] = ips
			mu.Unlock()
		})

	}

	wg.Wait()
	return namePool
}

// ResolveNS - resolving of NS records in DNS
func (rs *ResolverService) ResolveNS(names []string) map[string][]string {
	nsPool := map[string][]string{}
	mu := sync.Mutex{}
	wg := &sync.WaitGroup{}

	for _, name := range names {

		wg.Go(func() {

			if isIP(name) {
				mu.Lock()
				nsPool[name] = []string{}
				mu.Unlock()
				return
			}

			if rs.provider == nil {
				var nsL []string
				data, _ := net.LookupNS(name)

				for _, NS := range data {
					nsL = append(nsL, NS.Host)
				}

				mu.Lock()
				nsPool[name] = nsL
				mu.Unlock()
				return
			}

			var nsL []string
			domain := models.Domain(name)

			nsRec, err := rs.provider.Query(context.TODO(), domain, models.TypeNS)
			if err == nil {
				for _, ans := range nsRec.Answer {
					nsL = append(nsL, ans.Data)
				}
			}

			mu.Lock()
			nsPool[name] = nsL
			mu.Unlock()
		})

	}

	wg.Wait()
	return nsPool
}
