package iproute

import (
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
)

// GetRoutes fetches the routes using the `ip route list` command and parses the output.
func GetRoutes() ([]Route, error) {
	// Run the `ip route list` command
	cmd := exec.Command("ip", "route", "list")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error running ip route list: %v", err)
	}

	// Split output into lines
	lines := strings.Split(string(output), "\n")

	var routes []Route
	for _, line := range lines {
		if line != "" {
			route, err := parseRouteLine(line)
			if err != nil {
				fmt.Printf("Failed to parse route: %v\n", err)
				continue
			}
			routes = append(routes, *route)
		}
	}

	return routes, nil
}

// parseRouteLine parses a single route line into a Route struct.
func parseRouteLine(line string) (*Route, error) {
	tokens := strings.Fields(line)
	route := &Route{}
	i := 0

	if tokens[i] == "default" {
		_, ipnet, _ := net.ParseCIDR("0.0.0.0/0")
		route.Destination = ipnet
		i++
	} else {
		if _, ipnet, err := net.ParseCIDR(tokens[i]); err == nil {
			route.Destination = ipnet
		} else {
			return nil, fmt.Errorf("invalid CIDR: %s", tokens[i])
		}
		i++
	}

	for i < len(tokens) {
		switch tokens[i] {
		case "via":
			i++
			route.Via = net.ParseIP(tokens[i])
		case "dev":
			i++
			route.Dev = tokens[i]
		case "proto":
			i++
			route.Proto = tokens[i]
		case "scope":
			i++
			route.Scope = tokens[i]
		case "src":
			i++
			route.Src = net.ParseIP(tokens[i])
		case "metric":
			i++
			if m, err := strconv.Atoi(tokens[i]); err == nil {
				route.Metric = &m
			}
		}
		i++
	}

	return route, nil
}
