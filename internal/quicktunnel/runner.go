package quicktunnel

import (
	"context"
	"errors"
	"fmt"

	"github.com/Jonathansl17/cftun/internal/hostfs"
	"github.com/Jonathansl17/cftun/internal/paths"
	"github.com/Jonathansl17/cftun/internal/sysexec"
)

func (r Runner) Run(ctx context.Context, localURL string) error {
	config, err := r.Files.EmptyTemp(paths.QuickConfigPattern)
	if err != nil {
		return fmt.Errorf("create quick tunnel config: %w", err)
	}
	runErr := r.run(ctx, config, localURL)
	if removeErr := config.Remove(); removeErr != nil {
		if errors.Is(runErr, ErrInterrupted) {
			return removeErr
		}
		return errors.Join(runErr, removeErr)
	}
	return runErr
}

func (r Runner) run(ctx context.Context, config hostfs.Temp, localURL string) error {
	if err := config.Write(configContent); err != nil {
		return err
	}
	err := r.Tunnels.QuickTunnel(ctx, localURL, config.Path)
	if sysexec.Interrupted(err) {
		return ErrInterrupted
	}
	return err
}
