package uikit

import (
	"github.com/Jonathansl17/cftun/internal/ingress"
	"github.com/Jonathansl17/cftun/internal/msg"
	"github.com/Jonathansl17/cftun/internal/routes"
	"github.com/Jonathansl17/cftun/internal/store"
)

const (
	annotationNeedsCloudflared = "cftun/needs-cloudflared"
	annotationEnabled          = "true"

	tableMinWidth = 0
	tableTabWidth = 4
	tablePadding  = 2
	tablePadChar  = ' '
	tableFlags    = 0
	columnSep     = "\t"
	joinSeparator = "\n"

	subjectTextFormat = "%s: %s"
	causeTextFormat   = "%s: %v"
)

var matchers = []Matcher{
	matchRollback,
	matchDNSRollback,
	matchUnsafeTarget,
	sentinel(routes.ErrTunnelNotFound, msg.ErrTunnelNotFound),
	sentinel(routes.ErrCredentialsMissing, msg.ErrCredentialsMissing),
	sentinel(ingress.ErrMalformed, msg.ErrMalformed),
	sentinel(ingress.ErrDuplicate, msg.ErrDuplicate),
	sentinel(ingress.ErrNotFound, msg.ErrNotFound),
	sentinel(ingress.ErrInvalidHostname, msg.ErrInvalidHostname),
	sentinel(ingress.ErrInvalidPort, invalidPortText()),
	sentinel(store.ErrMissing, msg.ErrConfigMissing),
}
