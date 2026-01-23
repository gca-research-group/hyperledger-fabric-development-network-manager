package configtx

import "gopkg.in/yaml.v3"

type Node yaml.Node

func ScalarNode(value string) *Node {
	return (*Node)(&yaml.Node{Kind: yaml.ScalarNode, Value: value})
}

func MappingNode(content ...*Node) *Node {
	_content := make([]*yaml.Node, len(content))
	for i, c := range content {
		_content[i] = (*yaml.Node)(c)
	}

	return (*Node)(&yaml.Node{Kind: yaml.MappingNode, Content: _content})
}

func SequenceNode(content ...*Node) *Node {
	_content := make([]*yaml.Node, len(content))
	for i, c := range content {
		_content[i] = (*yaml.Node)(c)
	}

	return (*Node)(&yaml.Node{Kind: yaml.SequenceNode, Content: _content})
}

func AliasNode(value string, alias *Node) *Node {
	return (*Node)(&yaml.Node{Kind: yaml.AliasNode, Value: value, Alias: (*yaml.Node)(alias)})
}

func (n *Node) WithAnchor(name string) *Node {
	n.Anchor = name
	return n
}

func (n *Node) WithTag(name string) *Node {
	n.Tag = name
	return n
}

func (n *Node) MarshalYAML() (*yaml.Node, error) {
	if n == nil {
		return nil, nil
	}
	return (*yaml.Node)(n), nil
}

func (n *Node) GetOrCreateValue(keyName string, defaultNode *Node) *Node {
	for i := 0; i < len(n.Content); i += 2 {
		if n.Content[i].Value == keyName {
			return (*Node)(n.Content[i+1])
		}
	}

	n.Content = append(n.Content, (*yaml.Node)(ScalarNode(keyName)), (*yaml.Node)(defaultNode))
	return defaultNode
}

func (n *Node) WithStyle(style yaml.Style) *Node {
	n.Style = style
	return n
}
