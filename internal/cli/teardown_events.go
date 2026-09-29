package cli

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
	"github.com/Jonathansl17/cftun/internal/msg"
	"github.com/Jonathansl17/cftun/internal/teardown"
)

func (p teardownPrinter) Observe(e teardown.Event) {
	switch e.Kind {
	case teardown.EventStepFailed:
		p.Report.Printf(msg.WarnStepFailed, stepLabel(e.Step, e.Subject), uikit.Describe(e.Err))
	case teardown.EventNoConfig:
		p.Report.Printf(msg.WarnNoConfig, uikit.Describe(e.Err))
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
	config, err := filepath.Abs(a.Store.Path())
	if err != nil {
		return teardown.Params{}, fmt.Errorf(resolvePathFormat, err)
	}
	backup, err := filepath.Abs(a.Store.BackupPath())
	if err != nil {
		return teardown.Params{}, fmt.Errorf(resolvePathFormat, err)
	}
	userDir, err := a.userCloudflaredDir()
	if err != nil {
		return teardown.Params{}, err
	}
	return teardown.Params{ConfigPath: config, BackupPath: backup, UserDir: userDir}, nil
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
		return err
	}
	text := fmt.Sprintf(stepTextFormat, stepLabel(step.Step, step.Subject), uikit.Describe(step.Err))
	return uikit.NewDisplayError(text, err)
}
