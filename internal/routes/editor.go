package routes

import (
	"context"
	"errors"

	"github.com/Jonathansl17/cftun/internal/ingress"
)

func (e Editor) List() ([]ingress.Rule, error) {
	doc, err := e.Store.Load()
	if err != nil {
		return nil, err
	}
	return doc.Rules(), nil
}

func (e Editor) Add(ctx context.Context, rule ingress.Rule, opts Options) (Outcome, error) {
	doc, validation, err := e.applier().apply(ctx, func(d *ingress.Document) error { return d.Add(rule) })
	out := Outcome{Validation: validation}
	if err != nil {
		return out, err
	}
	if !opts.SkipDNS {
		if err := e.Router.RouteDNS(ctx, doc.Tunnel(), rule.Hostname); err != nil {
			return out, e.applier().undo(ctx, err)
		}
	}
	return e.restart(ctx, out, nil, opts)
}

func (e Editor) Remove(ctx context.Context, host string, opts Options) (Outcome, error) {
	_, validation, err := e.applier().apply(ctx, func(d *ingress.Document) error { return d.Remove(host) })
	out := Outcome{Validation: validation}
	if err != nil {
		return out, err
	}
	var dnsErr error
	if !opts.SkipDNS {
		out.DNS, dnsErr = e.DNS.Delete(ctx, host)
	}
	return e.restart(ctx, out, dnsErr, opts)
}

func (e Editor) Edit(ctx context.Context, host string, rule ingress.Rule, opts Options) (Outcome, error) {
	doc, validation, err := e.applier().apply(ctx, func(d *ingress.Document) error { return d.Update(host, rule) })
	out := Outcome{Validation: validation}
	if err != nil {
		return out, err
	}
	var dnsErr error
	if host != rule.Hostname && !opts.SkipDNS {
		if err := e.Router.RouteDNS(ctx, doc.Tunnel(), rule.Hostname); err != nil {
			return out, e.applier().undo(ctx, err)
		}
		out.DNS, dnsErr = e.DNS.Delete(ctx, host)
	}
	return e.restart(ctx, out, dnsErr, opts)
}

func (e Editor) restart(ctx context.Context, out Outcome, prior error, opts Options) (Outcome, error) {
	if opts.SkipRestart {
		out.RestartSkipped = true
		return out, prior
	}
	return out, errors.Join(prior, e.Restarter.Restart(ctx))
}

func (e Editor) applier() applier {
	return applier{store: e.Store, checker: e.Checker}
}
