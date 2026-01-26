package configtx

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/internal/yaml"
)

type Policy struct {
	Qualifier string
	Rule      string
}

func NewSignaturePolicy(mspID string, role string) *yaml.Node {
	return yaml.MappingNode(
		yaml.ScalarNode(TypeKey),
		yaml.ScalarNode(SignatureKey),
		yaml.ScalarNode(RuleKey),
		yaml.ScalarNode(fmt.Sprintf("OR('%s.%s')", mspID, role)).WithDoubleQuotedStyle(),
	)
}

func NewImplicitMetaPolicy(policy Policy) *yaml.Node {
	if policy.Qualifier == "" {
		policy.Qualifier = "ANY"
	}

	return yaml.MappingNode(
		yaml.ScalarNode(TypeKey),
		yaml.ScalarNode(ImplicitMetaKey),
		yaml.ScalarNode(RuleKey),
		yaml.ScalarNode(fmt.Sprintf("%s %s", policy.Qualifier, policy.Rule)).WithDoubleQuotedStyle(),
	)
}

func NewMemberPolicy(mspID string) *yaml.Node {
	return NewSignaturePolicy(mspID, memberKey)
}

func NewPeerPolicy(mspID string) *yaml.Node {
	return NewSignaturePolicy(mspID, peerKey)
}

func NewAdminPolicy(mspID string) *yaml.Node {
	return NewSignaturePolicy(mspID, adminKey)
}
