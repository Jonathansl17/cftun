package cli

import (
	"errors"

	"github.com/Jonathansl17/cftun/internal/ingress"
	"github.com/Jonathansl17/cftun/internal/msg"
)

var errNoRules = errors.New(msg.ErrNoRules)

func valueOrAsk(a *App, value, label string, validate func(string) error) (string, error) {
	if value != "" {
		return value, validate(value)
	}
	return a.Prompt.Ask(label, validate)
}

func askHostname(a *App, value, label string) (string, error) {
	raw, err := valueOrAsk(a, value, label, func(s string) error {
		_, err := ingress.NormalizeHostname(s)
		return err
	})
	if err != nil {
		return "", err
	}
	return ingress.NormalizeHostname(raw)
}

func askPort(a *App, value, label string) (int, error) {
	raw, err := valueOrAsk(a, value, label, func(s string) error {
		_, err := ingress.ParsePort(s)
		return err
	})
	if err != nil {
		return 0, err
	}
	return ingress.ParsePort(raw)
}

func pickRule(a *App, value string) (ingress.Rule, error) {
	rules, err := a.Routes.List()
	if err != nil {
		return ingress.Rule{}, err
	}
	if value != "" {
		for _, r := range rules {
			if host, _ := ingress.NormalizeHostname(value); host == r.Hostname {
				return r, nil
			}
		}
		return ingress.Rule{}, ingress.ErrNotFound
	}
	if len(rules) == 0 {
		return ingress.Rule{}, errNoRules
	}
	labels := make([]string, len(rules))
	for i, r := range rules {
		labels[i] = r.Hostname + msg.RuleArrow + r.Service
	}
	i, err := a.Prompt.Select(msg.PromptPickRule, labels)
	if err != nil {
		return ingress.Rule{}, err
	}
	return rules[i], nil
}

func firstArg(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return ""
}

func notEmpty(s string) error {
	if s == "" {
		return errors.New(msg.ErrEmpty)
	}
	return nil
}
