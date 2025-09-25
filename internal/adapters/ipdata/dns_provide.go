// Copyright (c) 2025 EterLine (Andrew)
// This file is part of micro-utils.
// Licensed under the MIT License. See the LICENSE file for details.

package ipdata

import (
	"context"
	"fmt"
	"net"
	"sync"

	doh "github.com/eterline/micro-utils/pkg/DoH"
	dns "github.com/miekg/dns"
)

type DoHqueryService interface {
	Query(ctx context.Context, d doh.Domain, t doh.Record) (doh.DnsResponse, error)
	Service() string
}

// DoHResolve - DNS over HTTP/s adapter type
type DoHResolve struct {
	rs DoHqueryService
}

// NewCloudflareResolver - cloudflare DNS over HTTP/s
func NewCloudflareResolver() *DoHResolve {
	return &DoHResolve{
		rs: doh.InitDnsCloudflareProvider(),
	}
}

// NewGoogleResolver - google DNS over HTTP/s
func NewGoogleResolver() *DoHResolve {
	return &DoHResolve{
		rs: doh.InitDnsCloudflareProvider(),
	}
}

func (rs *DoHResolve) ResolveIP(ctx context.Context, s string) ([]net.IP, error) {

	var ips []net.IP
	var errs []error

	resA, errA := rs.rs.Query(ctx, doh.Domain(s), doh.TypeA)
	if errA == nil && resA.Status == doh.NOERROR {
		for _, ans := range resA.Answer {
			if ip := net.ParseIP(ans.Data); ip != nil {
				ips = append(ips, ip)
			}
		}
	} else if errA != nil {
		errs = append(errs, fmt.Errorf("A query failed: %w", errA))
	} else {
		errs = append(errs, fmt.Errorf("A query returned status %v", resA.Status))
	}

	resAAAA, errAAAA := rs.rs.Query(ctx, doh.Domain(s), doh.TypeAAAA)
	if errAAAA == nil && resAAAA.Status == doh.NOERROR {
		for _, ans := range resAAAA.Answer {
			if ip := net.ParseIP(ans.Data); ip != nil {
				ips = append(ips, ip)
			}
		}
	} else if errAAAA != nil {
		errs = append(errs, fmt.Errorf("AAAA query failed: %w", errAAAA))
	} else {
		errs = append(errs, fmt.Errorf("AAAA query returned status %v", resAAAA.Status))
	}

	if len(ips) > 0 {
		return ips, nil
	}

	return nil, fmt.Errorf("no IP resolved for %s: %v", s, errs)

}

func (rs *DoHResolve) ResolveNS(ctx context.Context, s string) ([]string, error) {
	res, err := rs.rs.Query(ctx, doh.Domain(s), doh.TypeNS)
	if err != nil {
		return nil, err
	}

	var nss []string
	for _, ans := range res.Answer {
		nss = append(nss, ans.Data)
	}

	return nss, nil
}

// =======================================

// LocalResolve - use localhost or system DNS server as IP resolve server
type LocalResolve struct{}

func NewLocalResolver() LocalResolve {
	return LocalResolve{}
}

func (rs LocalResolve) ResolveIP(ctx context.Context, s string) ([]net.IP, error) {
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, network, address)
		},
	}

	ips, err := r.LookupIP(ctx, "ip", s)
	if err != nil {
		return nil, err
	}

	return ips, nil
}

func (rs LocalResolve) ResolveNS(ctx context.Context, s string) ([]string, error) {
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, network, address)
		},
	}

	nss, err := r.LookupNS(ctx, s)
	if err != nil {
		return nil, err
	}

	nssL := make([]string, len(nss))
	for i, ns := range nss {
		nssL[i] = ns.Host
	}

	return nssL, nil
}

// LocalResolve - use certain DNS server as IP resolve server.
type RemoteResolve struct {
	dnsSocket string
}

func correctDNSsrv(socket string) string {
	host, port, err := net.SplitHostPort(socket)
	if err != nil {
		return net.JoinHostPort(socket, "53")
	}
	return net.JoinHostPort(host, port)
}

func NewRemoteResolver(socket string) (*RemoteResolve, error) {
	// check DNS socket string.
	// replace empty port with :53 as DNS default port.
	//	"8.8.8.8" -> "8.8.8.8:53"
	//	"10.192.0.33:7890" - use as it is
	host, port, err := net.SplitHostPort(socket)
	if err != nil {
		socket = net.JoinHostPort(socket, "53")
	} else {
		socket = net.JoinHostPort(host, port)
	}

	return &RemoteResolve{
		dnsSocket: correctDNSsrv(socket),
	}, nil
}

func (rs *RemoteResolve) ResolveIP(ctx context.Context, s string) ([]net.IP, error) {

	var (
		fqdn = dns.Fqdn(s)
		ips  []net.IP
		errs []error
		mu   sync.Mutex
		wg   sync.WaitGroup
	)

	query := func(qtype uint16) {
		defer wg.Done()
		msg := new(dns.Msg)
		msg.SetQuestion(fqdn, qtype)

		res, err := dns.ExchangeContext(ctx, msg, rs.dnsSocket)
		if err != nil {
			mu.Lock()
			errs = append(errs, err)
			mu.Unlock()
			return
		}

		mu.Lock()
		defer mu.Unlock()

		for _, ans := range res.Answer {
			switch rr := ans.(type) {
			case *dns.A:
				ips = append(ips, rr.A)
			case *dns.AAAA:
				ips = append(ips, rr.AAAA)
			}
		}
	}

	wg.Add(2)
	go query(dns.TypeA)    // request for IPv4
	go query(dns.TypeAAAA) // request for IPv6
	wg.Wait()

	if len(ips) > 0 {
		return ips, nil
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("no IP resolved for %s: %v", s, errs)
	}

	return nil, fmt.Errorf("no IP resolved for %s", s)
}

func (rs *RemoteResolve) ResolveNS(ctx context.Context, s string) ([]string, error) {
	resMsg := dns.Msg{}
	fqdn := dns.Fqdn(s)

	resMsg.SetQuestion(fqdn, dns.TypeNS)

	res, err := dns.ExchangeContext(ctx, &resMsg, rs.dnsSocket)
	if err != nil {
		return nil, err
	}

	var nss []string
	for _, ans := range res.Answer {
		if a, ok := ans.(*dns.NS); ok {
			nss = append(nss, a.Ns)
		}
	}

	if len(nss) > 0 {
		return nss, nil
	}

	return nil, fmt.Errorf("no NS resolved for %s", s)
}
