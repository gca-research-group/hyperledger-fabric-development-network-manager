package configtx

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/constants"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/yaml"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/config"
)

type Builder struct {
	config *config.Config

	ordererOrgs      []*yaml.Node
	appOrgs          []*yaml.Node
	ordererAliases   []*yaml.Node
	appAliases       map[string]*yaml.Node
	ordererAddresses []string
}

func NewBuilder(config *config.Config) *Builder {
	return &Builder{config: config, appAliases: make(map[string]*yaml.Node)}
}

func (c *Builder) BuildOrganizations() {
	for _, organization := range c.config.Organizations {
		if len(organization.Orderers) > 0 {
			ordererAddresses := []string{}
			for _, orderer := range organization.Orderers {
				ordererAddresses = append(ordererAddresses, fmt.Sprintf("%s.%s:%d", orderer.Subdomain, organization.Domain, orderer.Port))
			}

			organizationName := fmt.Sprintf("%sOrderer", organization.Name)
			mspID := config.ResolveOrdererMSPID(organization)
			org := NewOrdererOrganization(organizationName, organization.Domain, mspID, ordererAddresses).
				WithDefaultOrdererPolicies(mspID)

			c.ordererOrgs = append(c.ordererOrgs, org.Build())
			c.ordererAliases = append(c.ordererAliases, yaml.AliasNode(organizationName, org.Build()))

			c.ordererAddresses = append(c.ordererAddresses, ordererAddresses...)
		}
	}

	for _, organization := range c.config.Organizations {
		mspID := config.ResolveOrganizationMSPID(organization)

		var anchorPeerHost string
		var anchorPeerPort int

		for i, peer := range organization.Peers {
			if peer.IsAnchor {
				anchorPeerHost = fmt.Sprintf("%s.%s", peer.Subdomain, organization.Domain)
				anchorPeerPort = peer.Port
			}

			if anchorPeerHost == "" && i == 0 {
				anchorPeerHost = fmt.Sprintf("%s.%s", peer.Subdomain, organization.Domain)
				anchorPeerPort = constants.DEFAULT_PEER_PORT
			}
		}

		if anchorPeerPort == 0 {
			anchorPeerPort = constants.DEFAULT_PEER_PORT
		}

		org := NewApplicationOrganization(organization.Name, organization.Domain, mspID, c.ordererAddresses, c.config.Capabilities).
			WithAnchorPeer(anchorPeerHost, anchorPeerPort).
			WithDefaultApplicationPolicies(mspID)

		c.appOrgs = append(c.appOrgs, org.Build())
		c.appAliases[organization.Name] = yaml.AliasNode(organization.Name, org.Build())
	}
}

// func (c *Builder) BuildProfiles(
// 	orderer *yaml.Node,
// 	application *yaml.Node,
// 	channel *yaml.Node,
// 	appAliases []*yaml.Node,
// 	appCapability *yaml.Node,
// ) []*yaml.Node {
// 	var profiles []*yaml.Node

// 	for _, profile := range c.config.Profiles {
// 		var appAliases []*yaml.Node
// 		for _, organization := range profile.Organizations {
// 			appAliases = append(appAliases, c.appAliases[organization])
// 		}

// 		currentProfile := NewProfile(profile.Name, orderer, application, channel, appAliases, appCapability).Build()

// 		for _, node := range currentProfile.Content {
// 			profiles = append(profiles, (*yaml.Node)(node))
// 		}
// 	}

// 	return profiles
// }

func (c *Builder) Build() (*yaml.Node, error) {
	c.BuildOrganizations()

	appCapLabel, appCapVal := NewApplicationCapability(c.config.Capabilities.Application)
	ordCapLabel, ordCapVal := NewOrdererCapability(c.config.Capabilities.Orderer)
	chCapLabel, chCapVal := NewChannelCapability(c.config.Capabilities.Channel)

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

	profiles := []*yaml.Node{}

	orderer := NewOrderer(c.config.Capabilities).
		WithAddresses(c.ordererAddresses, c.config.Capabilities).
		WithCapabilities(ordCapVal).
		WithPolicies().
		WithOrganizations(c.ordererAliases)

	for _, profile := range c.config.Profiles {
		var appAliases []*yaml.Node
		for _, organization := range profile.Organizations {
			appAliases = append(appAliases, c.appAliases[organization])
		}

		currentProfile := NewProfile(profile, c.config.Organizations, orderer, application, channel, appAliases, appCapVal).Build()

		for _, node := range currentProfile.Content {
			profiles = append(profiles, (*yaml.Node)(node))
		}
	}

	nodes := []*yaml.Node{
		yaml.ScalarNode(CapabilitiesKey), yaml.MappingNode(appCapLabel, appCapVal, ordCapLabel, ordCapVal, chCapLabel, chCapVal),
		yaml.ScalarNode(OrganizationsKey), yaml.SequenceNode(append(c.ordererOrgs, c.appOrgs...)...),
		yaml.ScalarNode(OrdererKey), orderer.WithAnchor(OrdererDefaultsKey),
		yaml.ScalarNode(ApplicationKey), application,
		yaml.ScalarNode(ChannelKey), channel,
		yaml.ScalarNode(ProfilesKey), yaml.MappingNode(profiles...),
	}

	return yaml.MappingNode(
		nodes...,
	), nil
}
