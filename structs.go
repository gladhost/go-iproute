package iproute

import "net"

type Route struct {
	Destination *net.IPNet `json:"destination"`      // CIDR like 0.0.0.0/0 or 192.168.1.0/24
	Via         net.IP     `json:"via,omitempty"`    // gateway IP
	Dev         string     `json:"dev"`              // network interface name
	Proto       string     `json:"proto,omitempty"`  // protocol (e.g., kernel, dhcp)
	Scope       string     `json:"scope,omitempty"`  // scope (e.g., link)
	Src         net.IP     `json:"src,omitempty"`    // source IP
	Metric      *int       `json:"metric,omitempty"` // optional metric value
}
