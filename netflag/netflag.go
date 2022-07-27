package netflag

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/pizchen/gou/aflag"
)

const (
	DPORT_MSK = (1 << 0)
	SPORT_MSK = (1 << 1)
	EPORT_MSK = (1 << 2)
	DADDR_MSK = (1 << 4)
	SADDR_MSK = (1 << 5)
	EADDR_MSK = (1 << 6)
)

type Netaddr struct {
	Ip   net.IP
	Port uint16
}

type NetaddrConfig struct {
	Eaddr []*Netaddr
	Saddr []*Netaddr
	Daddr []*Netaddr
	V4msk []byte
	V6msk []byte
	Flag  byte
}

func procAddrString(addr string) (*Netaddr, error) {

	var err error
	var port int

	na := &Netaddr{Ip: make([]byte, 0), Port: 0}
	parts := strings.Split(addr, "@")

	switch len(parts) {
	case 2:
		port, err = strconv.Atoi(parts[1])
		if err != nil {
			goto INVALID
		}
		na.Port = uint16(port)
		fallthrough
	case 1:
		if strings.ContainsAny(parts[0], ".:") {
			na.Ip = net.ParseIP(parts[0])
			if na.Ip == nil {
				goto INVALID
			}
		} else if len(parts) == 1 {
			port, err = strconv.Atoi(parts[1])
			if err != nil || na.Port > 65535 {
				goto INVALID
			}
			na.Port = uint16(port)
		} else if len(parts[0]) > 0 {
			goto INVALID
		}
	default:
		goto INVALID
	}

	return na, nil

INVALID:
	return nil, fmt.Errorf("Invalid address: %v", addr)
}

func procAddrs(addrs aflag.ArrayFlagString) ([]*Netaddr, error) {

	hasIp := false
	hasPort := false
	res := make([]*Netaddr, 0)

	for i, s := range addrs {
		na, err := procAddrString(s)
		if err != nil {
			return nil, err
		}
		if i == 0 {
			if len(na.Ip) > 0 {
				hasIp = true
			}
			if na.Port > 0 {
				hasPort = true
			}
		} else {
			if hasIp && len(na.Ip) == 0 {
				return nil, fmt.Errorf("Conflicts found in addrs: missing IP [%v]", s)
			}
			if !hasIp && len(na.Ip) > 0 {
				return nil, fmt.Errorf("Conflicts found in addrs: extra IP [%v]", s)
			}
			if hasPort && na.Port == 0 {
				return nil, fmt.Errorf("Conflicts found in addrs: missing port [%v]", s)
			}
			if !hasPort && na.Port > 0 {
				return nil, fmt.Errorf("Conflicts found in addrs: extra port [%v]", s)
			}
		}
		res = append(res, na)
	}

	return res, nil
}

var (
	Args struct {
		Host  aflag.ArrayFlagString
		Src   aflag.ArrayFlagString
		Dst   aflag.ArrayFlagString
		V4msk uint
		V6msk uint
	}
)

func FlagAdd() {
	flag.Var(&Args.Src, "src", "Match source address: IP@PORT, comma separated list or repeated")
	flag.Var(&Args.Dst, "dst", "Match dest address: IP@PORT, comma separated list or repeated")
	flag.Var(&Args.Host, "host", "Match src/dst address: IP@PORT, comma separated list or repeated")
	flag.UintVar(&Args.V4msk, "msk4", 0, "ipv4 address netmask bits number")
	flag.UintVar(&Args.V6msk, "msk6", 0, "ipv6 address netmask bits number")
}

func FlagParse() (*NetaddrConfig, error) {

	cfg := &NetaddrConfig{
		Eaddr: make([]*Netaddr, 0),
		Saddr: make([]*Netaddr, 0),
		Daddr: make([]*Netaddr, 0),
		V4msk: net.CIDRMask(32, 32),
		V6msk: net.CIDRMask(128, 128),
		Flag:  0,
	}
	var err error

	if !flag.Parsed() {
		return nil, fmt.Errorf("command line arguments not parsed")
	}

	if len(Args.Host) > 0 && (len(Args.Src) > 0 || len(Args.Dst) > 0) {
		return nil, fmt.Errorf("Option -host is exclusive with -src/-dst")
	}

	if Args.V4msk > 0 {
		if Args.V4msk > 32 {
			Args.V4msk = 32
		}
		if Args.V4msk < 8 {
			return nil, fmt.Errorf("Unsupported netmask value: IPv4 < 8")
		}
		copy(cfg.V4msk[:], net.CIDRMask(int(Args.V4msk), 32))
	}

	if Args.V6msk > 0 {
		if Args.V6msk > 128 {
			Args.V6msk = 128
		}
		if Args.V6msk < 64 {
			return nil, fmt.Errorf("Unsupported netmask value: IPv6 < 64")
		}
		copy(cfg.V6msk[:], net.CIDRMask(int(Args.V6msk), 128))
	}

	if len(Args.Host) > 0 {
		if cfg.Eaddr, err = procAddrs(Args.Host); err != nil {
			return nil, err
		}
		if len(cfg.Eaddr[0].Ip) > 0 {
			cfg.Flag |= EADDR_MSK
		}
		if cfg.Eaddr[0].Port > 0 {
			cfg.Flag |= EPORT_MSK
		}
	}
	if len(Args.Src) > 0 {
		if cfg.Saddr, err = procAddrs(Args.Src); err != nil {
			return nil, err
		}
		if len(cfg.Saddr[0].Ip) > 0 {
			cfg.Flag |= SADDR_MSK
		}
		if cfg.Saddr[0].Port > 0 {
			cfg.Flag |= SPORT_MSK
		}
	}
	if len(Args.Dst) > 0 {
		if cfg.Daddr, err = procAddrs(Args.Dst); err != nil {
			return nil, err
		}
		if len(cfg.Daddr[0].Ip) > 0 {
			cfg.Flag |= DADDR_MSK
		}
		if cfg.Daddr[0].Port > 0 {
			cfg.Flag |= DPORT_MSK
		}
	}

	return cfg, nil
}
