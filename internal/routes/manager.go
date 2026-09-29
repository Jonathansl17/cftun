package routes

import (
	"context"
	"errors"
	"fmt"

	"github.com/Jonathansl17/cloudfare/internal/ingress"
	"github.com/Jonathansl17/cloudfare/internal/msg"
	"github.com/Jonathansl17/cloudfare/internal/service"
)

// Manager applies route changes.
type Manager struct {
	ConfigPath string
	Store      Store
	Tunnels    Tunnels
	Service    Service
	DNS        DNS
	Report     Reporter
}

// List returns the configured hostname rules.
func (m Manager) List() ([]ingress.Rule, error) {
	doc, err := m.Store.Load()
	if err != nil {
		return nil, err
	}
	return doc.Rules(), nil
}

// Add publishes rule.
func (m Manager) Add(ctx context.Context, rule ingress.Rule, opts Options) error {
	doc, err := m.apply(ctx, func(d *ingress.Document) error { return d.Add(rule) })
	if err != nil {
		return err
	}
	if !opts.SkipDNS {
		if err := m.Tunnels.RouteDNS(ctx, doc.Tunnel(), rule.Hostname); err != nil {
			return err
		}
	}
	return m.restart(ctx, opts)
}

// Remove unpublishes host.
func (m Manager) Remove(ctx context.Context, host string, opts Options) error {
	if _, err := m.apply(ctx, func(d *ingress.Document) error { return d.Remove(host) }); err != nil {
		return err
	}
	if !opts.SkipDNS {
		if err := m.DeleteDNS(ctx, host); err != nil {
			return err
		}
	}
	return m.restart(ctx, opts)
}

// Edit replaces the rule of host with rule, moving DNS when the name changes.
func (m Manager) Edit(ctx context.Context, host string, rule ingress.Rule, opts Options) error {
	doc, err := m.apply(ctx, func(d *ingress.Document) error { return d.Update(host, rule) })
	if err != nil {
		return err
	}
	if host != rule.Hostname && !opts.SkipDNS {
		if err := m.Tunnels.RouteDNS(ctx, doc.Tunnel(), rule.Hostname); err != nil {
			return err
		}
		if err := m.DeleteDNS(ctx, host); err != nil {
			return err
		}
	}
	return m.restart(ctx, opts)
}

// DeleteDNS removes the CNAME of host, or tells the user to do it by hand
// when no API token is configured.
func (m Manager) DeleteDNS(ctx context.Context, host string) error {
	if !m.DNS.Configured() {
		m.Report.Printf(msg.WarnManualDNS, host)
		return nil
	}
	count, err := m.DNS.DeleteCNAME(ctx, host)
	if err != nil {
		return err
	}
	m.Report.Printf(msg.InfoDNSDeleted, count, host)
	return nil
}

// Validate runs cloudflared's ingress validation on the config.
func (m Manager) Validate(ctx context.Context) error {
	out, err := m.Tunnels.ValidateIngress(ctx, m.ConfigPath)
	if err != nil {
		return err
	}
	m.Report.Printf(msg.InfoValidated, out)
	return nil
}

// apply mutates and saves the config, rolling back if validation fails.
func (m Manager) apply(ctx context.Context, mutate func(*ingress.Document) error) (*ingress.Document, error) {
	doc, err := m.Store.Load()
	if err != nil {
		return nil, err
	}
	if err := mutate(doc); err != nil {
		return nil, err
	}
	if err := m.Store.Save(ctx, doc); err != nil {
		return nil, err
	}
	if err := m.Validate(ctx); err != nil {
		if rerr := m.Store.Restore(ctx); rerr != nil {
			return nil, errors.Join(err, fmt.Errorf("restore backup: %w", rerr))
		}
		return nil, fmt.Errorf("%s: %w", msg.ErrRolledBack, err)
	}
	return doc, nil
}

func (m Manager) restart(ctx context.Context, opts Options) error {
	if opts.SkipRestart {
		m.Report.Printf(msg.InfoRestartSkipped)
		return nil
	}
	return m.Service.Run(ctx, service.Restart)
}
