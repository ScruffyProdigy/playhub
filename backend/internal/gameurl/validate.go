package gameurl

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"
)

// ValidateOutboundURL checks a game API or play URL before Lobby fetches or POSTs to it.
func ValidateOutboundURL(ctx context.Context, raw string, requireHTTPS bool) error {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return fmt.Errorf("gameurl: URL is required")
	}
	u, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("gameurl: invalid URL: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("gameurl: URL scheme must be http or https")
	}
	if requireHTTPS && u.Scheme != "https" {
		return fmt.Errorf("gameurl: HTTPS is required in production")
	}
	if u.User != nil {
		return fmt.Errorf("gameurl: URL must not include userinfo")
	}
	host := strings.TrimSpace(u.Hostname())
	if host == "" {
		return fmt.Errorf("gameurl: URL host is required")
	}
	if strings.EqualFold(host, "localhost") || host == "127.0.0.1" || host == "::1" {
		if requireHTTPS {
			return fmt.Errorf("gameurl: localhost URLs are not allowed in production")
		}
		return nil
	}

	resolver := net.DefaultResolver
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
	}
	ips, err := resolver.LookupIPAddr(ctx, host)
	if err != nil {
		return fmt.Errorf("gameurl: DNS lookup failed for %q: %w", host, err)
	}
	if len(ips) == 0 {
		return fmt.Errorf("gameurl: no addresses for host %q", host)
	}
	for _, addr := range ips {
		if blockedIP(addr.IP) {
			return fmt.Errorf("gameurl: URL host %q resolves to blocked address %s", host, addr.IP)
		}
	}
	return nil
}

func blockedIP(ip net.IP) bool {
	if ip == nil {
		return true
	}
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsPrivate() {
		return true
	}
	if ip.IsMulticast() || ip.IsUnspecified() {
		return true
	}
	// CGNAT / shared address space 100.64.0.0/10
	if len(ip) == net.IPv4len && ip[0] == 100 && ip[1] >= 64 && ip[1] <= 127 {
		return true
	}
	return false
}
