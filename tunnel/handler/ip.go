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

func isPrivateIP(addr string) bool {
	ip := net.ParseIP(addr)

	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}

	for _, block := range privateIPs {
		if block.Contains(ip) {
			return true
		}
	}

	return false
}
