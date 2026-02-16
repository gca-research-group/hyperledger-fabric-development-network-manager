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
	for o := range config.Organizations {
		for i := range config.Organizations[o].Orderers {
			if config.Organizations[o].Orderers[i].Port == 0 {
				config.Organizations[o].Orderers[i].Port = 7050
			}
		}
	}

	return &Builder{config: config, appAliases: make(map[string]*yaml.Node)}
}

func (c *Builder) BuildOrganizations() {
	for _, organization := range c.config.Organizations {
		for _, orderer := range organization.Orderers {
			mspID := BuildMSPID(orderer.Name)
			org := NewOrdererOrganization(orderer.Name, organization.Domain, mspID).
				WithDefaultOrdererPolicies(mspID)

			c.ordererOrgs = append(c.ordererOrgs, org.Build())
			c.ordererAliases = append(c.ordererAliases, yaml.AliasNode(orderer.Name, org.Build()))

			ordererAddress := fmt.Sprintf("%s.%s:%d", orderer.Subdomain, organization.Domain, orderer.Port)
			c.ordererAddresses = append(c.ordererAddresses, ordererAddress)
		}
	}

	for _, organization := range c.config.Organizations {
		mspID := BuildMSPID(organization.Name)

		var anchorPeerHost string
		var anchorPeerPort int
		defaultAnchorPeerPort := 7051
		for i, peer := range organization.Peers {
			if peer.IsAnchor {
				anchorPeerHost = fmt.Sprintf("%s.%s", peer.Subdomain, organization.Domain)
				anchorPeerPort = peer.Port
			}

			if anchorPeerHost == "" && i == 0 {
				anchorPeerHost = fmt.Sprintf("%s.%s", peer.Subdomain, organization.Domain)
				anchorPeerPort = defaultAnchorPeerPort
			}
		}

		if anchorPeerPort == 0 {
			anchorPeerPort = defaultAnchorPeerPort
		}

		org := NewApplicationOrganization(organization.Name, organization.Domain, mspID, c.ordererAddresses).
			WithAnchorPeer(anchorPeerHost, anchorPeerPort).
			WithDefaultApplicationPolicies(mspID)

		c.appOrgs = append(c.appOrgs, org.Build())
		c.appAliases[organization.Name] = yaml.AliasNode(organization.Name, org.Build())
	}
}

func (c *Builder) BuildProfiles(
	orderer *yaml.Node,
	application *yaml.Node,
	channel *yaml.Node,
	appAliases []*yaml.Node,
	appCapability *yaml.Node,
) []*yaml.Node {
	var profiles []*yaml.Node

	for _, profile := range c.config.Profiles {
		var appAliases []*yaml.Node
		for _, organization := range profile.Organizations {
			appAliases = append(appAliases, c.appAliases[organization])
		}

		currentProfile := NewProfile(profile.Name, orderer, application, channel, appAliases, appCapability).Build()

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

	profiles := c.BuildProfiles(orderer, application, channel, appAliases, appCapVal)

	return yaml.MappingNode(
		yaml.ScalarNode(CapabilitiesKey), yaml.MappingNode(appCapLabel, appCapVal, ordCapLabel, ordCapVal, chCapLabel, chCapVal),
		yaml.ScalarNode(OrganizationsKey), yaml.SequenceNode(append(c.ordererOrgs, c.appOrgs...)...),
		yaml.ScalarNode(OrdererKey), orderer,
		yaml.ScalarNode(ApplicationKey), application,
		yaml.ScalarNode(ChannelKey), channel,
		yaml.ScalarNode(ProfilesKey), yaml.MappingNode(profiles...),
	), nil
}
