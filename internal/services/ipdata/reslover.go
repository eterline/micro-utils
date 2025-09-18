// Copyright (c) 2025 EterLine (Andrew)
// This file is part of micro-utils.
// Licensed under the MIT License. See the LICENSE file for details.


package ipdata

import (
	"context"
	"net"
	"regexp"
	"runtime"
	"strings"
	"sync"

	microutils "github.com/eterline/micro-utils"
	"github.com/eterline/micro-utils/internal/models"
)

var (
	ipRegex = regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
)

/*
ResolverService - implements DNS resolving logic

	If resolver got nil, will use system DNS resolving
*/
type NetworkScrapeService struct {
	res        models.Resolver
	maxWorkers int
}

/*
NewResolverService - init resolver service

	If resolver got nil, will use system DNS resolving
*/
func NewNetworkScrapeService(r models.Resolver, workers int) *NetworkScrapeService {
	return &NetworkScrapeService{
		res:        r,
		maxWorkers: microutils.Clamp(workers, 1, runtime.NumCPU()),
	}
}

func isIP(s string) bool {
	return ipRegex.MatchString(strings.TrimSpace(s))
}

// ResolveAs - resolving of A and AAAA records in DNS
func (rs *NetworkScrapeService) ResolveDNS(ctx context.Context, names []string) map[string]models.AboutResolve {
	resolvPool := map[string]models.AboutResolve{}
	mu := sync.Mutex{}
	wg := &sync.WaitGroup{}
	tickets := make(chan struct{}, rs.maxWorkers)

	for _, name := range names {

		wg.Go(func() {
			tickets <- struct{}{}
			defer func() { <-tickets }()

			if isIP(name) {
				res := models.AboutResolve{
					IPs:         []net.IP{net.ParseIP(name)},
					NameServers: []string{},
				}
				mu.Lock()
				resolvPool[name] = res
				mu.Unlock()
				return
			}

			var (
				res      models.AboutResolve
				wgWorker sync.WaitGroup
			)

			wgWorker.Go(func() {
				ips, err := rs.res.ResolveIP(ctx, name)
				if err != nil {
					res.ErrorIPs = err.Error()
					return
				}
				res.IPs = ips
			})

			wgWorker.Go(func() {
				ns, err := rs.res.ResolveNS(ctx, name)
				if err != nil {
					res.ErrorNS = err.Error()
					return
				}
				res.NameServers = ns
			})

			wgWorker.Wait()

			mu.Lock()
			resolvPool[name] = res
			mu.Unlock()
		})
	}

	wg.Wait()
	return resolvPool
}
