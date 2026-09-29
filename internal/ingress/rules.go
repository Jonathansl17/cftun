package ingress

import (
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/Jonathansl17/cftun/internal/apperr"
)

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
		return apperr.Wrap(r.Hostname, ErrDuplicate)
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
	return apperr.Wrap(host, ErrNotFound)
}

func (d *Document) Update(host string, r Rule) error {
	n := d.node(host)
	if n == nil {
		return apperr.Wrap(host, ErrNotFound)
	}
	if !sameHost(host, r.Hostname) && d.node(r.Hostname) != nil {
		return apperr.Wrap(r.Hostname, ErrDuplicate)
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
