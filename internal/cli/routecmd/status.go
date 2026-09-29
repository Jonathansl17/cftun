package routecmd

import (
	"fmt"

	"github.com/Jonathansl17/cftun/internal/health"
	"github.com/Jonathansl17/cftun/internal/msg"
)

func statusLabel(r health.Result) string {
	switch {
	case r.Err != nil:
		return msg.StatusDown
	case r.OK():
		return fmt.Sprintf(msg.StatusOKFormat, r.Status)
	default:
		return fmt.Sprintf(msg.StatusFailFormat, r.Status)
	}
}
