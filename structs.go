package cidr4

import (
	"fmt"
	"strconv"
	"strings"
)

type IPv4 uint32

func (n IPv4) String() string {
	dd := [4]int{}
	for i := 0; i < 4; i++ {
		dd[len(dd)-1-i] = int((n >> (8 * i)) & 0xff)
	}
	return fmt.Sprintf("%d.%d.%d.%d", dd[0], dd[1], dd[2], dd[3])
}

func ParseIPv4(s string) (ip IPv4, err error) {
	tok := strings.Split(s, ".")
	if len(tok) != 4 {
		err = fmt.Errorf("ParseIPv4: cannot parse %s as IPv4 address, malformed", s)
		return
	}

	dd := [4]int{}
	for i := 0; i < 4; i++ {
		dd[i], err = strconv.Atoi(tok[i])
		if err != nil {
			err = fmt.Errorf("ParseIPv4: cannot parse %s as IPv4 address: %s", s, err.Error())
			return
		}
		if dd[i] < 0 || dd[i] > 0xff {
			err = fmt.Errorf("ParseIPv4: cannot parse %s as IPv4 address: %d not in 0-255", s, dd[i])
			return
		}
	}
	for i := 0; i < 4; i++ {
		ip |= IPv4(dd[len(dd)-1-i] << (8 * i))
	}
	return
}

type MASKv4 int

func (m MASKv4) String() string {
	return fmt.Sprintf("%d", m)
}

func ParseMASKv4(s string) (m MASKv4, err error) {
	var mi int
	mi, err = strconv.Atoi(s)
	if err != nil {
		err = fmt.Errorf("ParseMASKv4: cannot parse %s as IPv4 mask: %s", s, err.Error())
		return
	}
	if mi < 0 || mi > 32 {
		err = fmt.Errorf("ParseMASKv4: cannot parse %s as IPv4 mask: %d not in 0-32", s, mi)
		return
	}
	m = MASKv4(mi)
	return
}

func (m MASKv4) Bits() IPv4 {
	var bits IPv4
	bits = 1
	bits <<= 32 - m
	bits -= 1
	bits = ^bits
	return bits
}

type CIDRv4 struct {
	IP   IPv4
	Mask MASKv4
}

func (c CIDRv4) String() string {
	return fmt.Sprintf("%s/%s", c.IP.String(), c.Mask.String())
}

func ParseCIDRv4(s string) (c CIDRv4, err error) {
	tok := strings.Split(s, "/")
	if len(tok) != 2 {
		err = fmt.Errorf("ParseCIDRv4: cannot parse %s as IPv4 CIDR: malformed", s)
		return
	}

	var ip IPv4
	ip, err = ParseIPv4(tok[0])
	if err != nil {
		err = fmt.Errorf("ParseCIDRv4: cannot parse %s as IPv4 CIDR: %s", s, err.Error())
		return
	}
	c.IP = ip

	var m MASKv4
	m, err = ParseMASKv4(tok[1])
	if err != nil {
		err = fmt.Errorf("ParseCIDRv4: cannot parse %s as IPv4 CIDR: %s", s, err.Error())
		return
	}
	c.Mask = m

	c.Standardize()

	return
}

func (c *CIDRv4) Standardize() {
	c.IP = c.IP & c.Mask.Bits()
}

func (c CIDRv4) Contains(o CIDRv4) bool {
	if c.Mask > o.Mask {
		return false
	}
	return o.IP&c.Mask.Bits() == c.IP
}

func (c CIDRv4) LeastCommonCIDR(o CIDRv4) (r CIDRv4) {
	r = c
	if r.Mask < o.Mask {
		r.Mask = o.Mask
	}
	u := c.IP ^ o.IP
	for ; r.Mask >= 0; r.Mask-- {
		if r.Mask.Bits()&u == 0 {
			break
		}
	}
	r.Standardize()
	return
}

func (c CIDRv4) Merge(o CIDRv4) (ok bool, r CIDRv4) {
	if c.Contains(o) {
		ok = true
		r = c
		return
	}
	if o.Contains(c) {
		ok = true
		r = o
		return
	}
	if c.Mask == o.Mask {
		m := c.Mask - 1
		if c.IP != o.IP && c.IP&m.Bits() == o.IP&m.Bits() {
			ok = true
			r.IP = c.IP
			r.Mask = m
			r.Standardize()
			return
		}
	}
	return
}
