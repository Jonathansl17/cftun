package ingress

import (
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	ErrDuplicate = errors.New("hostname already has a rule")
	ErrNotFound  = errors.New("hostname has no rule")
)

type Rule struct {
	Hostname string
	Service  string
}

func (d *Document) Rules() []Rule {
	var rules []Rule
	for _, n := range d.list().Content {
		if host := scalarValue(n, keyHostname); host != "" {
			rules = append(rules, Rule{Hostname: host, Service: scalarValue(n, keyService)})
		}
	}
	return rules
}

func (d *Document) Add(r Rule) error {
	if d.node(r.Hostname) != nil {
		return fmt.Errorf("%s: %w", r.Hostname, ErrDuplicate)
	}
	d.ensureCatchAll()
	list := d.list()
	at := len(list.Content) - 1
	entry := newMapping(keyHostname, r.Hostname, keyService, r.Service)
	list.Content = append(list.Content[:at], append([]*yaml.Node{entry}, list.Content[at:]...)...)
	return nil
}

func (d *Document) Remove(host string) error {
	list := d.list()
	for i, n := range list.Content {
		if sameHost(scalarValue(n, keyHostname), host) {
			list.Content = append(list.Content[:i], list.Content[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("%s: %w", host, ErrNotFound)
}

func (d *Document) Update(host string, r Rule) error {
	n := d.node(host)
	if n == nil {
		return fmt.Errorf("%s: %w", host, ErrNotFound)
	}
	if !sameHost(host, r.Hostname) && d.node(r.Hostname) != nil {
		return fmt.Errorf("%s: %w", r.Hostname, ErrDuplicate)
	}
	setScalar(n, keyHostname, r.Hostname)
	setScalar(n, keyService, r.Service)
	return nil
}

func (d *Document) ensureCatchAll() {
	list := d.list()
	if last := len(list.Content) - 1; last >= 0 && lookup(list.Content[last], keyHostname) == nil {
		return
	}
	list.Content = append(list.Content, newMapping(keyService, CatchAllService))
}

func (d *Document) node(host string) *yaml.Node {
	for _, n := range d.list().Content {
		if sameHost(scalarValue(n, keyHostname), host) {
			return n
		}
	}
	return nil
}

func sameHost(a, b string) bool {
	return a != "" && strings.EqualFold(a, b)
}
