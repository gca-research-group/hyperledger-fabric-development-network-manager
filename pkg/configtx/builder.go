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
	appAliases       map[string]*yaml.Node
	ordererAddresses []string
}

func NewBuilder(config pkg.Config) *Builder {
	return &Builder{config: config, appAliases: make(map[string]*yaml.Node)}
}

func (c *Builder) BuildOrganizations() {
	for _, organization := range c.config.Organizations {
		mspID := BuildMSPID(organization.Name)
		org := NewApplicationOrganization(organization.Name, organization.Domain, mspID).
			WithAnchorPeer(organization.AnchorPeer).
			WithDefaultApplicationPolicies(mspID)

		c.appOrgs = append(c.appOrgs, org.Build())
		c.appAliases[organization.Name] = yaml.AliasNode(organization.Name, org.Build())

		for _, orderer := range organization.Orderers {
			mspID := BuildMSPID(orderer.Name)
			org := NewOrdererOrganization(orderer.Name, organization.Domain, mspID).
				WithDefaultOrdererPolicies(mspID)

			c.ordererOrgs = append(c.ordererOrgs, org.Build())
			c.ordererAliases = append(c.ordererAliases, yaml.AliasNode(orderer.Name, org.Build()))

			c.ordererAddresses = append(c.ordererAddresses, fmt.Sprintf("%s.%s:%d", orderer.Hostname, organization.Domain, orderer.Port))
		}
	}
}

func (c *Builder) BuildProfiles(
	orderer *yaml.Node,
	application *yaml.Node,
	channel *yaml.Node,
	appAliases []*yaml.Node,
) []*yaml.Node {
	var profiles []*yaml.Node

	defaultProfile := NewDefaultProfile(orderer, channel, c.ordererAliases, appAliases).Build()

	for _, node := range defaultProfile.Content {
		profiles = append(profiles, (*yaml.Node)(node))
	}

	for _, profile := range c.config.Profiles {
		var appAliases []*yaml.Node
		for _, organization := range profile.Organizations {
			appAliases = append(appAliases, c.appAliases[organization])
		}

		currentProfile := NewProfile(profile.Name, orderer, application, channel, c.ordererAliases, appAliases).Build()

		for _, node := range currentProfile.Content {
			profiles = append(profiles, (*yaml.Node)(node))
		}
	}

	return profiles
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
		WithRaftConfig(c.config.Organizations).
		WithAnchor(OrdererDefaultsKey)

	var appAliases []*yaml.Node

	for _, alias := range c.appAliases {
		appAliases = append(appAliases, alias)
	}

	application := NewApplication().
		WithPolicies().
		WithCapabilities(appCapVal).
		WithOrganizations(appAliases).
		WithAnchor(ApplicationDefaultsKey)

	channel := NewChannel().
		WithPolicies().
		WithCapabilities(chCapVal).
		WithAnchor(ChannelDefaultsKey)

	profiles := c.BuildProfiles(orderer, application, channel, appAliases)

	return yaml.MappingNode(
		yaml.ScalarNode(CapabilitiesKey), yaml.MappingNode(appCapLabel, appCapVal, ordCapLabel, ordCapVal, chCapLabel, chCapVal),
		yaml.ScalarNode(OrganizationsKey), yaml.SequenceNode(append(c.ordererOrgs, c.appOrgs...)...),
		yaml.ScalarNode(OrdererKey), orderer,
		yaml.ScalarNode(ApplicationKey), application,
		yaml.ScalarNode(ChannelKey), channel,
		yaml.ScalarNode(ProfilesKey), yaml.MappingNode(profiles...),
	), nil
}
