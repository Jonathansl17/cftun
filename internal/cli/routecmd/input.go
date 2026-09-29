package routecmd

import (
	"strconv"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
	"github.com/Jonathansl17/cftun/internal/ingress"
	"github.com/Jonathansl17/cftun/internal/msg"
)

func (g *Group) pickRule(value string) (ingress.Rule, error) {
	rules, err := g.deps.Routes.List()
	if err != nil {
		return ingress.Rule{}, err
	}
	if value != "" {
		return findRule(rules, value)
	}
	if len(rules) == 0 {
		return ingress.Rule{}, ErrNoRules
	}
	labels := make([]string, len(rules))
	for i, r := range rules {
		labels[i] = r.Hostname + msg.RuleArrow + r.Service
	}
	i, err := g.deps.Session.Prompt.Select(msg.PromptPickRule, labels)
	if err != nil {
		return ingress.Rule{}, err
	}
	return rules[i], nil
}

func findRule(rules []ingress.Rule, value string) (ingress.Rule, error) {
	host, err := ingress.NormalizeHostname(value)
	if err != nil {
		return ingress.Rule{}, err
	}
	for _, r := range rules {
		if r.Hostname == host {
			return r, nil
		}
	}
	return ingress.Rule{}, ingress.ErrNotFound
}

func (g *Group) askEdit(rule ingress.Rule, flags routeFlags) (ingress.Rule, error) {
	newHost, port := flags.NewHost, flags.Port
	if newHost == "" && port == "" {
		var err error
		if newHost, port, err = g.askEditInteractive(rule); err != nil {
			return ingress.Rule{}, err
		}
	}
	hostname, err := uikit.AskHostname(g.deps.Session, orDefault(newHost, rule.Hostname), noLabel)
	if err != nil {
		return ingress.Rule{}, err
	}
	service := rule.Service
	if port != "" {
		p, err := uikit.AskPort(g.deps.Session, port, noLabel)
		if err != nil {
			return ingress.Rule{}, err
		}
		service = ingress.LocalService(p)
	}
	return ingress.Rule{Hostname: hostname, Service: service}, nil
}

func (g *Group) askEditInteractive(rule ingress.Rule) (string, string, error) {
	currentPort := ""
	if p, ok := ingress.Port(rule.Service); ok {
		currentPort = strconv.Itoa(p)
	}
	host, err := g.keepIfBlank(msg.PromptNewHostname, rule.Hostname)
	if err != nil {
		return "", "", err
	}
	port, err := g.keepIfBlank(msg.PromptNewPort, currentPort)
	return host, port, err
}

func (g *Group) keepIfBlank(label, current string) (string, error) {
	answer, err := g.deps.Session.Prompt.Ask(label+msg.CurrentValue(current), acceptAny)
	if err != nil {
		return "", err
	}
	return orDefault(answer, current), nil
}

func acceptAny(string) error {
	return nil
}

func orDefault(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
