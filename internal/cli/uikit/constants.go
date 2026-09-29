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
	annotationMenuGroup        = "cftun/menu-group"
	annotationMenuLabel        = "cftun/menu-label"
	annotationMenuOrder        = "cftun/menu-order"

	OrderSetup     = 10
	OrderQuick     = 20
	OrderRoutes    = 30
	OrderService   = 40
	OrderTunnels   = 50
	OrderAccount   = 60
	OrderUninstall = 70

	slotErrorFormat = "menu slot of %s: %w"

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
