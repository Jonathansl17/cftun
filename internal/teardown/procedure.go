package teardown

import (
	"context"
	"errors"
)

func (p Procedure) Run(ctx context.Context, params Params) error {
	_, versionErr := p.Tunnels.Version(ctx)
	installed := versionErr == nil
	registered := installed && p.Service.Registered()
	errs := []error{p.stopService(ctx, registered)}
	errs = append(errs, p.removeRoutes(ctx, installed)...)
	errs = append(errs, p.removeCloudflared(ctx, installed, registered)...)
	errs = append(errs, p.removeFiles(ctx, params), p.removeTemp(ctx))
	return errors.Join(errs...)
}

func (p Procedure) fail(step Step, subject string, err error) error {
	if err == nil {
		return nil
	}
	p.Observer.Observe(Event{Kind: EventStepFailed, Step: step, Subject: subject, Err: err})
	return &StepError{Step: step, Subject: subject, Err: err}
}
