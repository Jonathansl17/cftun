package ingress

import "gopkg.in/yaml.v3"

type Rule struct {
	Hostname string
	Service  string
}

type Document struct {
	doc *yaml.Node
}
