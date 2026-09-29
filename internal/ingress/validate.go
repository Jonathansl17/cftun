package ingress

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/Jonathansl17/cftun/internal/apperr"
)

func NormalizeHostname(raw string) (string, error) {
	host := strings.ToLower(strings.TrimSuffix(strings.TrimSpace(raw), "."))
	if !hostnamePattern.MatchString(host) {
		return "", apperr.Wrap(strconv.Quote(raw), ErrInvalidHostname)
	}
	return host, nil
}

func ParsePort(raw string) (int, error) {
	port, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || port < MinPort || port > MaxPort {
		return 0, apperr.Wrap(strconv.Quote(raw), ErrInvalidPort)
	}
	return port, nil
}

func LocalService(port int) string {
	return fmt.Sprintf(localServiceFormat, port)
}

func Port(service string) (int, bool) {
	u, err := url.Parse(service)
	if err != nil || u.Port() == "" {
		return 0, false
	}
	port, err := strconv.Atoi(u.Port())
	return port, err == nil
}
