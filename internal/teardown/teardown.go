// Package teardown removes every trace of cloudflared from the machine.
package teardown

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/Jonathansl17/cftun/internal/cloudflared"
	"github.com/Jonathansl17/cftun/internal/ingress"
	"github.com/Jonathansl17/cftun/internal/msg"
	"github.com/Jonathansl17/cftun/internal/service"
	"github.com/Jonathansl17/cftun/internal/sysexec"
)

// systemPaths are the locations cloudflared reads or writes outside the
// user's home.
var systemPaths = []string{
	"/etc/cloudflared",
	"/usr/local/etc/cloudflared",
	"/root/" + cloudflared.HomeDir,
	"/var/log/cloudflared.log",
	"/var/log/cloudflared.err",
	"/etc/systemd/system/cloudflared.service",
	"/etc/systemd/system/cloudflared-update.service",
	"/etc/systemd/system/cloudflared-update.timer",
	"/etc/init.d/cloudflared",
}

// Steps are the collaborators the teardown drives.
type Steps struct {
	Runner    sysexec.Runner
	Load      func() (*ingress.Document, error)
	DeleteDNS func(ctx context.Context, host string) error
	Tunnels   *cloudflared.Client
	Service   interface {
		Run(context.Context, service.Action) error
	}
	Uninstall         func(ctx context.Context) error
	Installed         func(ctx context.Context) bool
	ServiceRegistered func() bool
	Report            interface{ Printf(string, ...any) }
	// Paths are extra files to delete: config, backup, user credentials.
	Paths []string
}

// Run executes every step, continuing past failures so one broken piece does
// not leave the rest behind. All failures are returned together.
func (s Steps) Run(ctx context.Context) error {
	var errs []error
	record := func(step string, err error) {
		if err != nil {
			s.Report.Printf(msg.WarnStepFailed, step, err)
			errs = append(errs, fmt.Errorf("%s: %w", step, err))
		}
	}
	installed := s.Installed(ctx)
	if installed && s.ServiceRegistered() {
		record(msg.StepStopService, s.Service.Run(ctx, service.Stop))
	}
	if doc, err := s.Load(); err == nil {
		s.removeRoutes(ctx, doc, installed, record)
	} else {
		s.Report.Printf(msg.WarnNoConfig, err)
	}
	if installed {
		record(msg.StepUninstallService, s.Tunnels.UninstallService(ctx))
		record(msg.StepRemovePackage, s.Uninstall(ctx))
	}
	record(msg.StepRemoveFiles, s.removeFiles(ctx))
	return errors.Join(errs...)
}

// removeRoutes deletes DNS records and the tunnel while the cert still exists.
func (s Steps) removeRoutes(ctx context.Context, doc *ingress.Document, installed bool, record func(string, error)) {
	for _, rule := range doc.Rules() {
		record(msg.StepDeleteDNS+" "+rule.Hostname, s.DeleteDNS(ctx, rule.Hostname))
	}
	if tunnel := doc.Tunnel(); tunnel != "" && installed {
		record(msg.StepDeleteTunnel, s.Tunnels.DeleteTunnel(ctx, tunnel))
	}
}

func (s Steps) removeFiles(ctx context.Context) error {
	paths := append(append([]string{}, systemPaths...), s.Paths...)
	for _, p := range paths {
		s.Report.Printf(msg.InfoRemoving, filepath.Clean(p))
	}
	return s.Runner.Stream(ctx, sysexec.Command{
		Name:       "rm",
		Args:       append([]string{"-rf", "--"}, paths...),
		Privileged: true,
	})
}
