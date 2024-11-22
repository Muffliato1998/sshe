package utils

import (
	"fmt"
	"net"
)

func AssertIpAddressValid(ip string) error {
	if ip == "" {
		return fmt.Errorf("IP cannot be empty")
	}
	if net.ParseIP(ip) == nil {
		return fmt.Errorf("%s is an invalid IP address", ip)
	}

	return nil
}
