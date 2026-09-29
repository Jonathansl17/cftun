package installcmd

import (
	"errors"
	"fmt"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
	"github.com/Jonathansl17/cftun/internal/teardown"
)

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
