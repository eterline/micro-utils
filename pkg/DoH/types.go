// Copyright (c) 2025 EterLine (Andrew)
// This file is part of micro-utils.
// Licensed under the MIT License. See the LICENSE file for details.


package doh

import (
	"strings"

	"golang.org/x/net/idna"
)

// Domain - dns query domain
type Domain string

// Record - dns query type
type Record string

// Supported dns query type
const (
	// TypeA — IPv4 host address record.
	// Maps a domain name to an IPv4 address.
	// RFC 1035: https://www.rfc-editor.org/rfc/rfc1035
	TypeA = Record("A")

	// TypeAAAA — IPv6 host address record.
	// Maps a domain name to an IPv6 address.
	// RFC 3596: https://www.rfc-editor.org/rfc/rfc3596
	TypeAAAA = Record("AAAA")

	// TypeCNAME — Canonical Name record.
	// Maps an alias name to the true (canonical) domain name.
	// RFC 1035: https://www.rfc-editor.org/rfc/rfc1035
	TypeCNAME = Record("CNAME")

	// TypeMX — Mail Exchange record.
	// Defines mail servers responsible for accepting email on behalf of the domain.
	// RFC 1035: https://www.rfc-editor.org/rfc/rfc1035
	// RFC 7505 (Null MX): https://www.rfc-editor.org/rfc/rfc7505
	TypeMX = Record("MX")

	// TypeTXT — Text record.
	// Carries arbitrary human-readable or machine-readable text.
	// Commonly used for SPF, DKIM, and domain verification.
	// RFC 1035: https://www.rfc-editor.org/rfc/rfc1035
	TypeTXT = Record("TXT")

	// TypeSPF — Sender Policy Framework record (deprecated).
	// Used to specify authorized mail servers. Now replaced by TXT records.
	// RFC 7208: https://www.rfc-editor.org/rfc/rfc7208
	TypeSPF = Record("SPF")

	// TypeNS — Name Server record.
	// Delegates a DNS zone to use the given authoritative name servers.
	// RFC 1035: https://www.rfc-editor.org/rfc/rfc1035
	TypeNS = Record("NS")

	// TypeSOA — Start of Authority record.
	// Specifies authoritative information about a DNS zone,
	// including primary name server, serial number, and timers.
	// RFC 1035: https://www.rfc-editor.org/rfc/rfc1035
	TypeSOA = Record("SOA")

	// TypePTR — Pointer record.
	// Used mainly for reverse DNS lookups (IP address → domain name).
	// RFC 1035: https://www.rfc-editor.org/rfc/rfc1035
	TypePTR = Record("PTR")

	// TypeANY — Special query type.
	// Requests all available record types for a domain (discouraged in practice).
	// RFC 1035: https://www.rfc-editor.org/rfc/rfc1035
	TypeANY = Record("ANY")
)

func (t Record) String() string {
	return string(t)
}

type DoHStatusCode int

func (c DoHStatusCode) Error() string {
	switch c {
	case FORMERR:
		return "DoH query error - there was a format error in the DNS query itself"
	case SERVFAIL:
		return "DoH query error - the DNS server encountered an internal error and failed to process the request"
	case NXDOMAIN:
		return "DoH query error - the requested domain name does not exist"
	case REFUSED:
		return "DoH query error - the DNS server refused to answer the query"
	}

	return "DoH query completed successfully"
}

const (
	NOERROR DoHStatusCode = iota
	FORMERR
	SERVFAIL
	NXDOMAIN
	_
	REFUSED
)

// Question - dns query question for DoH providers
type Question struct {
	Name string `json:"name"`
	Type int    `json:"type"`
}

// Answer - dns query answer about record
type Answer struct {
	Name string `json:"name"`
	Type int    `json:"type"`
	TTL  int    `json:"TTL"`
	Data string `json:"data"`
}

// Response - dns query response from DoH providers
type DnsResponse struct {
	Status   DoHStatusCode `json:"Status"`
	TC       bool          `json:"TC"`
	RD       bool          `json:"RD"`
	RA       bool          `json:"RA"`
	AD       bool          `json:"AD"`
	CD       bool          `json:"CD"`
	Question []Question    `json:"Question"`
	Answer   []Answer      `json:"Answer"`
}

func (dr DnsResponse) Success() bool {
	return dr.Status == NOERROR
}

// Punycode - returns punycode of domain
func (d Domain) Punycode() (string, error) {
	name := strings.TrimSpace(string(d))

	return idna.New(
		idna.MapForLookup(),
		idna.Transitional(true),
		idna.StrictDomainName(false),
	).ToASCII(name)
}
