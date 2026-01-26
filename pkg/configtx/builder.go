package configtx

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/internal/yaml"
)

type Builder struct {
	config pkg.Config

	ordererOrgs      []*yaml.Node
	appOrgs          []*yaml.Node
	ordererAliases   []*yaml.Node
	appAliases       []*yaml.Node
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
		c.ordererAliases = append(c.ordererAliases, yaml.AliasNode(orderer.Name, org.Build()))

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
		c.appAliases = append(c.appAliases, yaml.AliasNode(organization.Name, org.Build()))
	}
}

func (c *Builder) Build() (*yaml.Node, error) {
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

	return yaml.MappingNode(
		yaml.ScalarNode(CapabilitiesKey), yaml.MappingNode(appCapLabel, appCapVal, ordCapLabel, ordCapVal, chCapLabel, chCapVal),
		yaml.ScalarNode(OrganizationsKey), yaml.SequenceNode(append(c.ordererOrgs, c.appOrgs...)...),
		yaml.ScalarNode(OrdererKey), orderer,
		yaml.ScalarNode(ApplicationKey), application,
		yaml.ScalarNode(ChannelKey), channel,
		yaml.ScalarNode(ProfilesKey), profiles,
	), nil
}
