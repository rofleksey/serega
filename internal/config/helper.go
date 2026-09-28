package config

import (
	"fmt"
	"net"
	"strings"
)

func parseCIDRs(raw string) ([]*net.IPNet, error) {
	if raw == "" {
		return nil, nil
	}

	parts := strings.Split(raw, ",")

	proxies := make([]*net.IPNet, 0, len(parts))
	for _, part := range parts {
		_, cidr, err := net.ParseCIDR(strings.TrimSpace(part))
		if err != nil {
			return nil, fmt.Errorf("SEREGA_TRUSTED_PROXY_CIDRS contains invalid CIDR %q", part)
		}

		proxies = append(proxies, cidr)
	}

	return proxies, nil
}
