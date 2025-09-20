// Copyright (c) 2025 EterLine (Andrew)
// This file is part of micro-utils.
// Licensed under the MIT License. See the LICENSE file for details.

package ipdata

import (
	"context"
	"errors"
	"net"
	"regexp"
	"strings"
	"sync"
	"time"

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

type IPstorage interface {
	Get(ctx context.Context, ip net.IP) (*models.AboutIPobject, error)
	Save(ctx context.Context, ip net.IP, obj models.AboutIPobject) error
}

type NetworkScrapeService struct {
	resolv     models.Resolver
	resumer    models.ResumerIP
	storage    IPstorage
	maxWorkers int
}

/*
NewResolverService - init resolver service

	If resolver got nil, will use system DNS resolving
*/
func NewNetworkScrapeService(
	workers int, rv models.Resolver, ru models.ResumerIP, st IPstorage,
) *NetworkScrapeService {
	return &NetworkScrapeService{
		resolv:     rv,
		resumer:    ru,
		storage:    st,
		maxWorkers: microutils.InitWorkersCountCurrently(workers),
	}
}

func isIP(s string) bool {
	return ipRegex.MatchString(strings.TrimSpace(s))
}

// ResolveAs - resolving of A and AAAA records in DNS
func (rs *NetworkScrapeService) ResolveDNS(ctx context.Context, names []string) (map[string]models.AboutResolve, error) {
	if names == nil {
		return map[string]models.AboutResolve{}, errors.New("resolving name pool is nil")
	}

	if len(names) < 1 {
		return map[string]models.AboutResolve{}, errors.New("resolving name pool is empty")
	}

	var (
		resolvPool = map[string]models.AboutResolve{}
		mu         = sync.Mutex{}
		wg         = &sync.WaitGroup{}
		tp         = microutils.NewTicketPool(rs.maxWorkers)
	)
	defer tp.ClosePool()

	for _, name := range names {
		wg.Go(func() {
			tp.CatchTicket()
			defer tp.PutTicket()

			startTime := time.Now()

			if isIP(name) {
				res := models.AboutResolve{
					IPs:         []net.IP{net.ParseIP(name)},
					NameServers: []string{},
				}
				res.CalcDuration(startTime)
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
				ips, err := rs.resolv.ResolveIP(ctx, name)
				if err != nil {
					res.ErrorIPs = err.Error()
					return
				}
				res.IPs = ips
			})

			wgWorker.Go(func() {
				ns, err := rs.resolv.ResolveNS(ctx, name)
				if err != nil {
					res.ErrorNS = err.Error()
					return
				}
				res.NameServers = ns
			})

			wgWorker.Wait()
			res.CalcDuration(startTime)

			mu.Lock()
			resolvPool[name] = res
			mu.Unlock()
		})
	}

	wg.Wait()
	return resolvPool, nil
}

func (rs *NetworkScrapeService) FetchAboutIP(ipPool []net.IP) ([]models.ResumeAboutIP, error) {
	if ipPool == nil {
		return []models.ResumeAboutIP{}, errors.New("ip pool is nil")
	}

	if len(ipPool) < 1 {
		return []models.ResumeAboutIP{}, errors.New("ip pool is empty")
	}

	var (
		resumes = make([]models.ResumeAboutIP, len(ipPool))
		mu      = sync.Mutex{}
		wg      = &sync.WaitGroup{}
		tp      = microutils.NewTicketPool(rs.maxWorkers)
	)
	defer tp.ClosePool()

	for i, ip := range ipPool {
		wg.Go(func() {
			tp.CatchTicket()
			defer tp.PutTicket()

			about := models.ResumeAboutIP{RequestIP: ip}

			if rs.storage != nil {
				obj, err := rs.storage.Get(context.Background(), ip)
				if obj != nil && err != nil {
					about.Resume = *obj
				}
			}

			obj, err := rs.resumer.ResumeIP(ip)
			if err != nil {
				about.Err = err.Error()
			} else {
				about.Resume = obj
			}

			mu.Lock()
			resumes[i] = about
			mu.Unlock()

			if rs.storage != nil {
				rs.storage.Save(context.Background(), ip, obj)
			}
		})
	}

	wg.Wait()
	return resumes, nil
}
