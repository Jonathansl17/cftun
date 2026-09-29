// Package ingress edits cloudflared config files while preserving unknown keys.
package ingress

import (
	"bytes"
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

const (
	keyTunnel      = "tunnel"
	keyCredentials = "credentials-file"
	keyIngress     = "ingress"
	keyHostname    = "hostname"
	keyService     = "service"
	yamlIndent     = 2

	// CatchAllService answers requests that match no hostname rule.
	CatchAllService = "http_status:404"
)

// ErrMalformed reports a config whose shape cloudflared would not accept.
var ErrMalformed = errors.New("config must be a mapping with an ingress list")

// Document is a parsed cloudflared config.
type Document struct {
	doc *yaml.Node
}

// Parse reads a cloudflared config.
func Parse(data []byte) (*Document, error) {
	var n yaml.Node
	if err := yaml.Unmarshal(data, &n); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	if n.Kind != yaml.DocumentNode || len(n.Content) == 0 || n.Content[0].Kind != yaml.MappingNode {
		return nil, ErrMalformed
	}
	if list := lookup(n.Content[0], keyIngress); list != nil && list.Kind != yaml.SequenceNode {
		return nil, ErrMalformed
	}
	return &Document{doc: &n}, nil
}

// New builds a config for tunnel with only the catch-all rule.
func New(tunnel, credentials string) *Document {
	root := newMapping(keyTunnel, tunnel, keyCredentials, credentials)
	d := &Document{doc: &yaml.Node{Kind: yaml.DocumentNode, Content: []*yaml.Node{root}}}
	d.ensureCatchAll()
	return d
}

// Bytes serializes the document.
func (d *Document) Bytes() ([]byte, error) {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(yamlIndent)
	if err := enc.Encode(d.doc); err != nil {
		return nil, fmt.Errorf("encode config: %w", err)
	}
	if err := enc.Close(); err != nil {
		return nil, fmt.Errorf("encode config: %w", err)
	}
	return buf.Bytes(), nil
}

// Tunnel returns the tunnel name or UUID the config points to.
func (d *Document) Tunnel() string {
	return scalarValue(d.root(), keyTunnel)
}

// Credentials returns the credentials file path.
func (d *Document) Credentials() string {
	return scalarValue(d.root(), keyCredentials)
}

func (d *Document) root() *yaml.Node {
	return d.doc.Content[0]
}

func (d *Document) list() *yaml.Node {
	return ensureChild(d.root(), keyIngress, yaml.SequenceNode)
}
