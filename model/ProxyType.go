package model

import "fmt"

type ProxyType int

const (
	ProxyTypeHttp  ProxyType = 1
	ProxyTypeHttps ProxyType = 2
	ProxyTypeTcp   ProxyType = 3
	ProxyTypeUdp   ProxyType = 4
)

func (x ProxyType) String() string {
	switch x {
	case ProxyTypeHttp:
		return "http"
	case ProxyTypeHttps:
		return "https"
	case ProxyTypeTcp:
		return "tcp"
	case ProxyTypeUdp:
		return "udp"
	default:
		return fmt.Sprintf("ProxyType(%d)", x)
	}
}
