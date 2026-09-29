package installcmd

import (
	"fmt"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
	"github.com/Jonathansl17/cftun/internal/msg"
	"github.com/Jonathansl17/cftun/internal/teardown"
)

func (o Observer) Observe(e teardown.Event) {
	switch e.Kind {
	case teardown.EventStepFailed:
		o.Session.Printf(msg.WarnStepFailed, stepLabel(e.Step, e.Subject), uikit.Describe(e.Err))
	case teardown.EventNoConfig:
		o.Session.Printf(msg.WarnNoConfig, uikit.Describe(e.Err))
	case teardown.EventRemoving:
		o.Session.Printf(msg.InfoRemoving, e.Subject)
	case teardown.EventDNSManual:
		o.Session.Printf(msg.WarnManualDNS, e.Subject)
	case teardown.EventDNSDeleted:
		o.Session.Printf(msg.InfoDNSDeleted, e.Count, e.Subject)
	}
}

func stepLabel(step teardown.Step, subject string) string {
	if subject == "" {
		return stepTexts[step]
	}
	return fmt.Sprintf(stepSubjectFormat, stepTexts[step], subject)
}
