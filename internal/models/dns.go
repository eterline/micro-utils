package models

import (
	"context"
	"strings"
	"vendor/golang.org/x/net/idna"
)

// Domain - dns query domain
type Domain string

// Type - dns query type
type DnsRecordType string

func (t DnsRecordType) String() string {
	return string(t)
}

// Question - dns query question
type Question struct {
	Name string `json:"name"`
	Type int    `json:"type"`
}

// Answer - dns query answer
type Answer struct {
	Name string `json:"name"`
	Type int    `json:"type"`
	TTL  int    `json:"TTL"`
	Data string `json:"data"`
}

// Response - dns query response
type DnsResponse struct {
	Status   int        `json:"Status"`
	TC       bool       `json:"TC"`
	RD       bool       `json:"RD"`
	RA       bool       `json:"RA"`
	AD       bool       `json:"AD"`
	CD       bool       `json:"CD"`
	Question []Question `json:"Question"`
	Answer   []Answer   `json:"Answer"`
}

// Supported dns query type
var (
	TypeA     = DnsRecordType("A")
	TypeAAAA  = DnsRecordType("AAAA")
	TypeCNAME = DnsRecordType("CNAME")
	TypeMX    = DnsRecordType("MX")
	TypeTXT   = DnsRecordType("TXT")
	TypeSPF   = DnsRecordType("SPF")
	TypeNS    = DnsRecordType("NS")
	TypeSOA   = DnsRecordType("SOA")
	TypePTR   = DnsRecordType("PTR")
	TypeANY   = DnsRecordType("ANY")
)

// Punycode - returns punycode of domain
func (d Domain) Punycode() (string, error) {
	name := strings.TrimSpace(string(d))

	return idna.New(
		idna.MapForLookup(),
		idna.Transitional(true),
		idna.StrictDomainName(false),
	).ToASCII(name)
}

type ProviderDoH interface {
	Query(ctx context.Context, d Domain, t DnsRecordType) (DnsResponse, error)
}
