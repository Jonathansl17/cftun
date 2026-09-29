package cli

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/Jonathansl17/cftun/internal/msg"
	"github.com/Jonathansl17/cftun/internal/teardown"
)

func (p teardownPrinter) Observe(e teardown.Event) {
	switch e.Kind {
	case teardown.EventStepFailed:
		p.Report.Printf(msg.WarnStepFailed, stepLabel(e.Step, e.Subject), e.Err)
	case teardown.EventNoConfig:
		p.Report.Printf(msg.WarnNoConfig, e.Err)
	case teardown.EventRemoving:
		p.Report.Printf(msg.InfoRemoving, e.Subject)
	case teardown.EventDNSManual:
		p.Report.Printf(msg.WarnManualDNS, e.Subject)
	case teardown.EventDNSDeleted:
		p.Report.Printf(msg.InfoDNSDeleted, e.Count, e.Subject)
	}
}

func stepLabel(step teardown.Step, subject string) string {
	if subject == "" {
		return stepTexts[step]
	}
	return fmt.Sprintf(stepSubjectFormat, stepTexts[step], subject)
}

func teardownParams(a *App) (teardown.Params, error) {
	config, err := filepath.Abs(a.ConfigPath)
	if err != nil {
		return teardown.Params{}, fmt.Errorf(resolvePathFormat, err)
	}
	backup, err := filepath.Abs(a.Store.BackupPath())
	if err != nil {
		return teardown.Params{}, fmt.Errorf(resolvePathFormat, err)
	}
	return teardown.Params{ConfigPath: config, BackupPath: backup, UserDir: a.UserCloudflaredDir()}, nil
}

func teardownFailure(err error) error {
	var joined interface{ Unwrap() []error }
	if !errors.As(err, &joined) {
		return err
	}
	failures := make([]error, 0, len(joined.Unwrap()))
	for _, failure := range joined.Unwrap() {
		failures = append(failures, stepFailure(failure))
	}
	return errors.Join(failures...)
}

func stepFailure(err error) error {
	var step *teardown.StepError
	if !errors.As(err, &step) {
		return present(err)
	}
	return fmt.Errorf(stepErrorFormat, stepLabel(step.Step, step.Subject), present(step.Err))
}
