package routecmd

import (
	"github.com/Jonathansl17/cftun/internal/msg"
	"github.com/Jonathansl17/cftun/internal/routes"
)

func (g *Group) printOutcome(host string, out routes.Outcome) {
	s := g.deps.Session
	s.Printf(msg.InfoValidated, out.Validation)
	switch out.DNS.State {
	case routes.DNSManual:
		s.Printf(msg.WarnManualDNS, host)
	case routes.DNSDeleted:
		s.Printf(msg.InfoDNSDeleted, out.DNS.Deleted, host)
	}
	if out.RestartSkipped {
		s.Printf(msg.InfoRestartSkipped)
	}
}
