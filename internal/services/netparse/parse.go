package netparse

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/netip"
	"os"
	"strings"

	"github.com/eterline/micro-utils/pkg/netipuse"
)

type netParser struct {
	setV6 *netipuse.PoolIPBuilder
	setV4 *netipuse.PoolIPBuilder
}

func NewNetParser() *netParser {
	return &netParser{
		setV6: &netipuse.PoolIPBuilder{},
		setV4: &netipuse.PoolIPBuilder{},
	}
}

func (p *netParser) ParseAddrs(addrs []string) error {
	if len(addrs) == 0 {
		return errors.New("addrs for parsing is empty")
	}

	added := 0

	for _, addr := range addrs {
		addr, err := netipuse.ParsePrefixOrAddr(addr)
		if err != nil {
			continue
		}

		added++

		switch {
		case addr.Is4():
			p.setV4.Add(addr)

		case addr.Is4In6():
			a := netip.AddrFrom4(addr.As4())
			p.setV4.Add(a)
			continue

		case addr.Is6():
			p.setV6.Add(addr)
		}
	}

	if added == 0 {
		return errors.New("no valid addrs for parsing")
	}

	return nil
}

func (p *netParser) ParseFromFile(file, sep string) (error, bool) {
	if sep == "" {
		sep = "\n"
	}

	f, err := os.OpenFile(file, os.O_RDONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open IPs file: %s - %w", file, err), false
	}
	defer f.Close()

	var (
		added   = 0
		fRd     = bufio.NewReader(f)
		byteSep = []byte(sep)[0]
	)

	for {
		line, err := fRd.ReadString(byteSep)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("failed to read line from file: %s - %w", file, err), false
		}

		line = strings.TrimSpace(line)

		addr, err := netipuse.ParsePrefixOrAddr(line)
		if err != nil {
			continue
		}

		added++

		switch {
		case addr.Is4():
			p.setV4.Add(addr)

		case addr.Is4In6():
			a := netip.AddrFrom4(addr.As4())
			p.setV4.Add(a)
			continue

		case addr.Is6():
			p.setV6.Add(addr)
		}
	}

	if added == 0 {
		return errors.New("no valid addrs for parsing"), false
	}

	return nil, true
}

func (p *netParser) Subnets() (v4, v6 []netip.Prefix) {
	pool, _ := p.setV4.PoolIP()
	if pool != nil {
		v4 = pool.Prefixes()
	}

	pool, _ = p.setV6.PoolIP()
	if pool != nil {
		v6 = pool.Prefixes()
	}

	return v4, v6
}

func ExportSubnetsTo(pfxs []netip.Prefix, w io.Writer) error {
	for _, pfx := range pfxs {
		_, err := fmt.Fprintln(w, pfx.String())
		if err != nil {
			return fmt.Errorf("failed to write prefix to writer: %s - %w", pfx.String(), err)
		}
	}
	return nil
}
