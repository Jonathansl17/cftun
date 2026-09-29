package ingress

import (
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	// ErrDuplicate reports a hostname that already has a rule.
	ErrDuplicate = errors.New("hostname already has a rule")
	// ErrNotFound reports a hostname without a rule.
	ErrNotFound = errors.New("hostname has no rule")
)

// Rule routes one public hostname to one local service.
type Rule struct {
	Hostname string
	Service  string
}

// Rules returns every hostname rule, excluding the catch-all.
func (d *Document) Rules() []Rule {
	var rules []Rule
	for _, n := range d.list().Content {
		if host := scalarValue(n, keyHostname); host != "" {
			rules = append(rules, Rule{Hostname: host, Service: scalarValue(n, keyService)})
		}
	}
	return rules
}

// Find returns the rule for host.
func (d *Document) Find(host string) (Rule, bool) {
	n := d.node(host)
	if n == nil {
		return Rule{}, false
	}
	return Rule{Hostname: scalarValue(n, keyHostname), Service: scalarValue(n, keyService)}, true
}

// Add inserts r just before the catch-all rule.
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

// Remove deletes the rule for host.
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

// Update replaces the rule for host with r, keeping any extra keys it has.
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

// ensureCatchAll guarantees the list ends with a rule without hostname.
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
