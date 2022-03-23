package handler

import (
	"net"

	"github.com/micro/micro/v3/service/logger"
)

var (
	privateIPs []*net.IPNet
)

func init() {
	for _, cidr := range []string{
		"127.0.0.0/8",    // IPv4 loopback
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
		"169.254.0.0/16", // RFC3927 link-local
		"::1/128",        // IPv6 loopback
		"fe80::/10",      // IPv6 link-local
		"fc00::/7",       // IPv6 unique local addr
	} {
		_, block, err := net.ParseCIDR(cidr)
		if err != nil {
			logger.Errorf("Error parsing ip block %q: %v", cidr, err)
			continue
		}
		privateIPs = append(privateIPs, block)
	}
}

func isPrivateIP(host string) bool {
	var addr string

	// split on host port
	h, _, err := net.SplitHostPort(host)
	if err == nil {
		host = h
	}

	// resolve the host
	addrs, err := net.LookupHost(host)
	if err != nil {
		addr = host
	} else {
		addr = addrs[0]
	}

	logger.Infof("Checking host %v address %v", host, addr)
	ip := net.ParseIP(addr)

	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		logger.Infof("Blocked ip local host %v address %v", host, addr)
		return true
	}

	for _, block := range privateIPs {
		if block.Contains(ip) {
			logger.Infof("Blocked ip cidr host %v address %v", host, addr)
			return true
		}
	}

	logger.Infof("No match for ip host %v address %v", host, addr)
	return false
}
