package health

import (
	"errors"
	"net/http"
	"syscall"
)

func (r Result) OK() bool {
	return r.Err == nil && r.Status < http.StatusInternalServerError
}

func (r Result) Hint() Hint {
	switch {
	case r.OK():
		return HintNone
	case errors.Is(r.Err, syscall.ECONNREFUSED):
		return HintConnectionRefused
	case r.dns:
		return HintDNS
	case r.Status == http.StatusBadGateway || r.Status == statusTunnelDown:
		return HintTunnelDown
	default:
		return HintGeneric
	}
}
