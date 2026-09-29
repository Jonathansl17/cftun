package ingress

import "gopkg.in/yaml.v3"

// lookup returns the value node stored under key in a mapping node.
func lookup(mapping *yaml.Node, key string) *yaml.Node {
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		if mapping.Content[i].Value == key {
			return mapping.Content[i+1]
		}
	}
	return nil
}

// scalarValue returns the scalar under key, or "" when absent.
func scalarValue(mapping *yaml.Node, key string) string {
	if n := lookup(mapping, key); n != nil {
		return n.Value
	}
	return ""
}

// setScalar overwrites the scalar under key, appending the pair when absent.
func setScalar(mapping *yaml.Node, key, value string) {
	if n := lookup(mapping, key); n != nil {
		n.Kind, n.Tag, n.Value = yaml.ScalarNode, "", value
		return
	}
	mapping.Content = append(mapping.Content, newScalar(key), newScalar(value))
}

// ensureChild returns the node under key, appending an empty one of kind.
func ensureChild(mapping *yaml.Node, key string, kind yaml.Kind) *yaml.Node {
	if n := lookup(mapping, key); n != nil {
		return n
	}
	child := &yaml.Node{Kind: kind}
	mapping.Content = append(mapping.Content, newScalar(key), child)
	return child
}

func newScalar(value string) *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Value: value}
}

func newMapping(pairs ...string) *yaml.Node {
	m := &yaml.Node{Kind: yaml.MappingNode}
	for i := 0; i+1 < len(pairs); i += 2 {
		setScalar(m, pairs[i], pairs[i+1])
	}
	return m
}
