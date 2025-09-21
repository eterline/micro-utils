// Copyright (c) 2025 EterLine (Andrew)
// This file is part of micro-utils.
// Licensed under the MIT License. See the LICENSE file for details.

package ipdata

import (
	"context"
	"fmt"
	"net"

	doh "github.com/eterline/micro-utils/pkg/DoH"
	dns "github.com/miekg/dns"
)

type CloudflareResolve struct {
	rs *doh.DnsDoHProvider
}

func NewCloudflareResolver() CloudflareResolve {
	return CloudflareResolve{
		rs: doh.InitDnsCloudflareProvider(),
	}
}

func (rs CloudflareResolve) ResolveIP(ctx context.Context, s string) ([]net.IP, error) {

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

func (rs CloudflareResolve) ResolveNS(ctx context.Context, s string) ([]string, error) {
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

type GoogleResolve struct {
	rs *doh.DnsDoHProvider
}

func NewGoogleResolver() GoogleResolve {
	return GoogleResolve{
		rs: doh.InitDnsCloudflareProvider(),
	}
}

func (rs GoogleResolve) ResolveIP(ctx context.Context, s string) ([]net.IP, error) {

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

func (rs GoogleResolve) ResolveNS(ctx context.Context, s string) ([]string, error) {
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

type RemoteResolve struct {
	dnsSocket string
}

func NewRemoteResolver(dns string) RemoteResolve {
	return RemoteResolve{
		dnsSocket: dns,
	}
}

func (rs RemoteResolve) ResolveIP(ctx context.Context, s string) ([]net.IP, error) {

	var ips []net.IP
	var errs []error

	resMsg := dns.Msg{}
	fqdn := dns.Fqdn(s)

	resMsg.SetQuestion(fqdn, dns.TypeA)
	resA, errA := dns.ExchangeContext(ctx, &resMsg, rs.dnsSocket)
	if errA == nil {
		for _, ans := range resA.Answer {
			if a, ok := ans.(*dns.A); ok {
				ips = append(ips, a.A)
			}
		}
	} else {
		errs = append(errs, errA)
	}

	resMsg.SetQuestion(fqdn, dns.TypeA)
	resAAAA, errAAAA := dns.ExchangeContext(ctx, &resMsg, rs.dnsSocket)
	if errAAAA == nil {
		for _, ans := range resAAAA.Answer {
			if a, ok := ans.(*dns.A); ok {
				ips = append(ips, a.A)
			}
		}
	} else {
		errs = append(errs, errAAAA)
	}

	if len(ips) > 0 {
		return ips, nil
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("no IP resolved for %s: %v", s, errs)
	}

	return nil, fmt.Errorf("no IP resolved for %s", s)
}

func (rs RemoteResolve) ResolveNS(ctx context.Context, s string) ([]string, error) {
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
