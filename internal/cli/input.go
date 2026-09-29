package cli

import (
	"errors"

	"github.com/Jonathansl17/cftun/internal/ingress"
	"github.com/Jonathansl17/cftun/internal/msg"
)

// errNoRules reports that there is nothing to pick from.
var errNoRules = errors.New(msg.ErrNoRules)

// valueOrAsk validates value when given, or prompts until a valid one arrives.
func valueOrAsk(a *App, value, label string, validate func(string) error) (string, error) {
	if value != "" {
		return value, validate(value)
	}
	return a.Prompt.Ask(label, validate)
}

// askHostname returns a normalized hostname from flag or prompt.
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

// askPort returns a valid port from flag or prompt.
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

// pickRule returns the rule named by value, or lets the user choose one.
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

// firstArg returns args[0] or "".
func firstArg(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return ""
}

// notEmpty rejects blank answers.
func notEmpty(s string) error {
	if s == "" {
		return errors.New(msg.ErrEmpty)
	}
	return nil
}
