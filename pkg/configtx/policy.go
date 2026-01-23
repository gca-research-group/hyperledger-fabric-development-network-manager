package configtx

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Policy struct {
	Qualifier string
	Rule      string
}

func NewSignaturePolicy(mspID string, role string) *Node {
	return MappingNode(
		ScalarNode(TypeKey),
		ScalarNode(SignatureKey),
		ScalarNode(RuleKey),
		ScalarNode(fmt.Sprintf("OR('%s.%s')", mspID, role)).WithStyle(yaml.DoubleQuotedStyle),
	)
}

func NewImplicitMetaPolicy(policy Policy) *Node {
	if policy.Qualifier == "" {
		policy.Qualifier = "ANY"
	}

	return MappingNode(
		ScalarNode(TypeKey),
		ScalarNode(ImplicitMetaKey),
		ScalarNode(RuleKey),
		ScalarNode(fmt.Sprintf("%s %s", policy.Qualifier, policy.Rule)).WithStyle(yaml.DoubleQuotedStyle),
	)
}

func NewMemberPolicy(mspID string) *Node {
	return NewSignaturePolicy(mspID, memberKey)
}

func NewPeerPolicy(mspID string) *Node {
	return NewSignaturePolicy(mspID, peerKey)
}

func NewAdminPolicy(mspID string) *Node {
	return NewSignaturePolicy(mspID, adminKey)
}
