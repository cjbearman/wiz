package wiz

import (
	"fmt"
	"net"
	"time"
)

var TimeFormat = time.RFC3339

// ipFromAddr extracts the IP address from a net.Addr, handling cases where the address includes a port.
func ipFromAddr(addr net.Addr) string {
	host, _, err := net.SplitHostPort(addr.String())
	if err != nil {
		// If it fails, the string likely doesn't have a port (e.g., "192.168.1.1")
		return addr.String()
	}
	return host
}

// formatMac takes a 12-character MAC address string and formats it as XX:XX:XX:XX:XX:XX. If the input is not 12 characters long, it marks it as invalid.
func formatMac(mac string) string {
	if len(mac) != 12 {
		return fmt.Sprintf("%s (invalid)", mac)
	}
	return fmt.Sprintf("%s:%s:%s:%s:%s:%s",
		mac[0:2], mac[2:4], mac[4:6], mac[6:8], mac[8:10], mac[10:12])
}
