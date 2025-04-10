package iproute

import (
	"net"
	"testing"
)

func TestParseRouteLine(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected Route
		err      bool
	}{
		{
			name: "Default route",
			line: "default via 192.168.1.1 dev eth0 proto dhcp metric 100",
			expected: Route{
				Destination: &net.IPNet{
					IP:   net.IP{0, 0, 0, 0},
					Mask: net.CIDRMask(0, 32),
				},
				Via:    net.IP{192, 168, 1, 1},
				Dev:    "eth0",
				Proto:  "dhcp",
				Metric: &[]int{100}[0],
			},
			err: false,
		},
		{
			name: "CIDR route with gateway",
			line: "172.16.0.0/16 via 10.0.0.1 dev eth1",
			expected: Route{
				Destination: &net.IPNet{
					IP:   net.IP{172, 16, 0, 0},
					Mask: net.CIDRMask(16, 32),
				},
				Via:   net.IP{10, 0, 0, 1},
				Dev:   "eth1",
				Proto: "",
			},
			err: false,
		},
		{
			name: "CIDR route without gateway",
			line: "10.0.0.0/24 dev eth1 proto kernel scope link src 10.0.0.5",
			expected: Route{
				Destination: &net.IPNet{
					IP:   net.IP{10, 0, 0, 0},
					Mask: net.CIDRMask(24, 32),
				},
				Dev:   "eth1",
				Proto: "kernel",
				Scope: "link",
				Src:   net.IP{10, 0, 0, 5},
			},
			err: false,
		},
		{
			name: "Route with metric",
			line: "192.168.1.0/24 dev eth0 proto kernel scope link src 192.168.1.100 metric 100",
			expected: Route{
				Destination: &net.IPNet{
					IP:   net.IP{192, 168, 1, 0},
					Mask: net.CIDRMask(24, 32),
				},
				Dev:    "eth0",
				Proto:  "kernel",
				Scope:  "link",
				Src:    net.IP{192, 168, 1, 100},
				Metric: &[]int{100}[0],
			},
			err: false,
		},
		{
			name:     "Invalid CIDR format",
			line:     "not_a_cidr dev eth0",
			expected: Route{},
			err:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseRouteLine(tt.line)
			if (err != nil) != tt.err {
				t.Errorf("parseRouteLine() error = %v, wantErr %v", err, tt.err)
				return
			}

			if err == nil {
				// Compare expected struct with the returned struct
				if !compareRoutes(*got, tt.expected) {
					t.Errorf("parseRouteLine() = %v, want %v", *got, tt.expected)
				}
			}
		})
	}
}

// compareRoutes is a helper function to compare two Route structs
func compareRoutes(r1, r2 Route) bool {
	if !r1.Destination.IP.Equal(r2.Destination.IP) || r1.Destination.Mask.String() != r2.Destination.Mask.String() {
		return false
	}

	if !r1.Via.Equal(r2.Via) {
		return false
	}

	if r1.Dev != r2.Dev || r1.Proto != r2.Proto || r1.Scope != r2.Scope || !r1.Src.Equal(r2.Src) {
		return false
	}

	// Handle optional Metric field
	if (r1.Metric == nil && r2.Metric != nil) || (r1.Metric != nil && r2.Metric == nil) {
		return false
	}

	if r1.Metric != nil && r2.Metric != nil && *r1.Metric != *r2.Metric {
		return false
	}

	return true
}

func TestGetRoutes(t *testing.T) {
	routes, err := GetRoutes()
	if err != nil {
		t.Fatal(err)
	}

	if len(routes) != 4 {
		t.Fatalf("found '%d' routes, expected: 4", len(routes))
	}
}
