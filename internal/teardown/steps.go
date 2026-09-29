package teardown

import (
	"context"

	"github.com/Jonathansl17/cftun/internal/paths"
	"github.com/Jonathansl17/cftun/internal/routes"
)

func (p Procedure) stopService(ctx context.Context, registered bool) error {
	if !registered {
		return nil
	}
	return p.fail(StepStopService, "", p.Service.Stop(ctx))
}

func (p Procedure) removeRoutes(ctx context.Context, installed bool) []error {
	doc, err := p.Config.Load()
	if err != nil {
		p.Observer.Observe(Event{Kind: EventNoConfig, Err: err})
		return nil
	}
	var errs []error
	for _, rule := range doc.Rules() {
		removal, err := p.DNS.Delete(ctx, rule.Hostname)
		p.announceDNS(rule.Hostname, removal)
		errs = append(errs, p.fail(StepDeleteDNS, rule.Hostname, err))
	}
	if tunnel := doc.Tunnel(); tunnel != "" && installed {
		errs = append(errs, p.fail(StepDeleteTunnel, "", p.Tunnels.DeleteTunnel(ctx, tunnel)))
	}
	return errs
}

func (p Procedure) announceDNS(host string, removal routes.DNSRemoval) {
	switch removal.State {
	case routes.DNSManual:
		p.Observer.Observe(Event{Kind: EventDNSManual, Subject: host})
	case routes.DNSDeleted:
		p.Observer.Observe(Event{Kind: EventDNSDeleted, Subject: host, Count: removal.Deleted})
	}
}

func (p Procedure) removeCloudflared(ctx context.Context, installed, registered bool) []error {
	if !installed {
		return nil
	}
	var errs []error
	if registered {
		errs = append(errs, p.fail(StepUninstallService, "", p.Tunnels.UninstallService(ctx)))
	}
	return append(errs, p.fail(StepRemovePackage, "", p.Package.Uninstall(ctx)))
}

func (p Procedure) removeFiles(ctx context.Context, params Params) error {
	targets := targetsFor(params)
	p.announceRemoving(targets)
	return p.fail(StepRemoveFiles, "", p.Files.RemoveAll(ctx, targets))
}

func (p Procedure) removeTemp(ctx context.Context) error {
	matches, err := p.Temp.TempMatches(paths.QuickConfigPattern)
	if err != nil {
		return p.fail(StepRemoveTemp, "", err)
	}
	if len(matches) == 0 {
		return nil
	}
	targets := make([]Target, len(matches))
	for i, match := range matches {
		targets[i] = Target{Path: match, Kind: TargetFile}
	}
	p.announceRemoving(targets)
	return p.fail(StepRemoveTemp, "", p.Files.RemoveAll(ctx, targets))
}

func (p Procedure) announceRemoving(targets []Target) {
	for _, target := range targets {
		p.Observer.Observe(Event{Kind: EventRemoving, Subject: target.Path})
	}
}
