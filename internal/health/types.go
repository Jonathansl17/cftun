package health

import "net/http"

type Hint int

const (
	HintNone Hint = iota
	HintConnectionRefused
	HintDNS
	HintTunnelDown
	HintGeneric
)

type Result struct {
	URL    string
	Status int
	Err    error
	dns    bool
}

type Checker struct {
	HTTP *http.Client
}
