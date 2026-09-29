package cli

import (
	"github.com/Jonathansl17/cftun/internal/msg"
	"github.com/Jonathansl17/cftun/internal/routes"
)

func printOutcome(r Reporter, host string, out routes.Outcome) {
	r.Printf(msg.InfoValidated, out.Validation)
	switch out.DNS.State {
	case routes.DNSManual:
		r.Printf(msg.WarnManualDNS, host)
	case routes.DNSDeleted:
		r.Printf(msg.InfoDNSDeleted, out.DNS.Deleted, host)
	}
	if out.RestartSkipped {
		r.Printf(msg.InfoRestartSkipped)
	}
}
