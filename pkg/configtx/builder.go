package configtx

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg"
)

type Builder struct {
	config pkg.Config

	ordererOrgs      []*Node
	appOrgs          []*Node
	ordererAliases   []*Node
	appAliases       []*Node
	ordererAddresses []string
}

func NewBuilder(config pkg.Config) *Builder {
	return &Builder{config: config}
}

func (c *Builder) BuildOrganizations() {
	// Build Orderer Orgs
	for _, orderer := range c.config.Orderers {
		mspID := BuildMSPID(orderer.Name)
		org := NewOrdererOrganization(orderer.Name, orderer.Domain, mspID).
			WithDefaultOrdererPolicies(mspID)

		c.ordererOrgs = append(c.ordererOrgs, org.Build())
		c.ordererAliases = append(c.ordererAliases, AliasNode(orderer.Name, org.Build()))

		for _, addr := range orderer.Addresses {
			c.ordererAddresses = append(c.ordererAddresses, fmt.Sprintf("%s:%d", addr.Host, addr.Port))
		}
	}

	// Build Application Orgs
	for _, organization := range c.config.Organizations {
		mspID := BuildMSPID(organization.Name)
		org := NewApplicationOrganization(organization.Name, organization.Domain, mspID).
			WithAnchorPeer(organization.AnchorPeer).
			WithDefaultApplicationPolicies(mspID)

		c.appOrgs = append(c.appOrgs, org.Build())
		c.appAliases = append(c.appAliases, AliasNode(organization.Name, org.Build()))
	}
}

func (c *Builder) Build() (*Node, error) {
	c.BuildOrganizations()

	appCapLabel, appCapVal := NewApplicationCapability()
	ordCapLabel, ordCapVal := NewOrdererCapability()
	chCapLabel, chCapVal := NewChannelCapability()

	orderer := NewOrderer().
		WithAddresses(c.ordererAddresses).
		WithCapabilities(ordCapVal).
		WithPolicies().
		WithOrganizations(c.ordererAliases).
		WithBatchConfig().
		WithRaftConfig(c.config.Orderers).
		WithAnchor(OrdererDefaultsKey)

	application := NewApplication().
		WithPolicies().
		WithCapabilities(appCapVal).
		WithOrganizations(c.appAliases).
		WithAnchor(ApplicationDefaultsKey)

	channel := NewChannel().
		WithPolicies().
		WithCapabilities(chCapVal).
		WithAnchor(ChannelDefaultsKey)

	profiles := NewDefaultProfiles(orderer, application, channel, c.ordererAliases, c.appAliases)

	return MappingNode(
		ScalarNode(CapabilitiesKey), MappingNode(appCapLabel, appCapVal, ordCapLabel, ordCapVal, chCapLabel, chCapVal),
		ScalarNode(OrganizationsKey), SequenceNode(append(c.ordererOrgs, c.appOrgs...)...),
		ScalarNode(OrdererKey), orderer,
		ScalarNode(ApplicationKey), application,
		ScalarNode(ChannelKey), channel,
		ScalarNode(ProfilesKey), profiles,
	), nil
}
