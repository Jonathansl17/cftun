package ingress

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

const (
	minPort            = 1
	maxPort            = 65535
	localServiceFormat = "http://localhost:%d"
)

var (
	// ErrInvalidHostname reports a value that is not a fully qualified hostname.
	ErrInvalidHostname = errors.New("invalid hostname, expected something like app.example.com")
	// ErrInvalidPort reports a value outside the TCP port range.
	ErrInvalidPort = fmt.Errorf("invalid port, expected a number between %d and %d", minPort, maxPort)

	hostnamePattern = regexp.MustCompile(`^([a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z]{2,63}$`)
)

// NormalizeHostname trims, lowercases and validates a hostname.
func NormalizeHostname(raw string) (string, error) {
	host := strings.ToLower(strings.TrimSuffix(strings.TrimSpace(raw), "."))
	if !hostnamePattern.MatchString(host) {
		return "", fmt.Errorf("%q: %w", raw, ErrInvalidHostname)
	}
	return host, nil
}

// ParsePort validates a TCP port.
func ParsePort(raw string) (int, error) {
	port, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || port < minPort || port > maxPort {
		return 0, fmt.Errorf("%q: %w", raw, ErrInvalidPort)
	}
	return port, nil
}

// LocalService builds the origin URL for a port on this machine.
func LocalService(port int) string {
	return fmt.Sprintf(localServiceFormat, port)
}

// Port extracts the port from a service URL, if it has one.
func Port(service string) (int, bool) {
	u, err := url.Parse(service)
	if err != nil || u.Port() == "" {
		return 0, false
	}
	port, err := strconv.Atoi(u.Port())
	return port, err == nil
}
