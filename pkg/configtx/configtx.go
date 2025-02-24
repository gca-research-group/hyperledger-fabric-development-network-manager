package configtx

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/configtx/application"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/configtx/capabilities"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/configtx/channel"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/configtx/orderer"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/configtx/organization"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/configtx/profiles"
)

type ConfigTx struct {
	Capabilities  capabilities.Capabilities   `yaml:"Capabilities"`
	Organizations []organization.Organization `yaml:"Organizations"`
	Orderer       orderer.Orderer             `yaml:"Orderer"`
	Application   application.Application     `yaml:"Application"`
	Channel       channel.Channel             `yaml:"Channel"`
	Profiles      map[string]interface{}      `yaml:"Profiles"`
}

func Build(config pkg.Config) ConfigTx {
	var _configtx ConfigTx
	var _organizations []string
	var _ordererAddresses []string

	_configtx.Capabilities = capabilities.NewCapabilities()

	for _, orderer := range config.Orderers {

		MSPID := fmt.Sprintf("%sMSP", orderer.Name)
		_configtx.Organizations = append(_configtx.Organizations, organization.Orderer{
			Name:   orderer.Name,
			ID:     MSPID,
			MSPDir: fmt.Sprintf("./crypto-materials/ordererOrganizations/%s/msp", orderer.Domain),
			Policies: organization.OrdererPolicies{
				Readers: organization.Policy{
					Type: "Signature",
					Rule: fmt.Sprintf("\"OR('%s.member')\"", MSPID),
				},
				Writers: organization.Policy{
					Type: "Signature",
					Rule: fmt.Sprintf("\"OR('%s.member')\"", MSPID),
				},
				Admins: organization.Policy{
					Type: "Signature",
					Rule: fmt.Sprintf("\"OR('%s.admin')\"", MSPID),
				},
			},
		})

		_ordererAddresses = append(_ordererAddresses, fmt.Sprintf("%s:%d", orderer.Domain, orderer.Port))
	}

	_configtx.Orderer = orderer.NewOrderer(_ordererAddresses)

	for _, peer := range config.Peers {

		MSPID := fmt.Sprintf("%sMSP", peer.Name)

		anchorPeerPort := peer.Port

		if anchorPeerPort < 1 {
			anchorPeerPort = 7051
		}

		_configtx.Organizations = append(_configtx.Organizations, organization.Peer{
			Name:   peer.Name,
			ID:     MSPID,
			MSPDir: fmt.Sprintf("./crypto-materials/peerOrganizations/%s/msp", peer.Domain),
			Policies: organization.PeerPolicies{
				Readers: organization.Policy{
					Type: "Signature",
					Rule: fmt.Sprintf("\"OR('%s.member')\"", MSPID),
				},
				Writers: organization.Policy{
					Type: "Signature",
					Rule: fmt.Sprintf("\"OR('%s.member')\"", MSPID),
				},
				Admins: organization.Policy{
					Type: "Signature",
					Rule: fmt.Sprintf("\"OR('%s.admin')\"", MSPID),
				},
				Endorsement: organization.Policy{
					Type: "Signature",
					Rule: fmt.Sprintf("\"OR('%s.member')\"", MSPID),
				},
			},
			AnchorPeers: []organization.AnchorPeer{{
				Host: peer.Domain,
				Port: anchorPeerPort,
			}},
		})

		_organizations = append(_organizations, peer.Name)
	}

	_configtx.Application = application.NewApplication(_organizations)
	_configtx.Channel = channel.NewChannel()
	_configtx.Profiles = map[string]interface{}{
		"OrdererProfile": profiles.OrdererProfile{
			Channel: "<<: *ChannelDefaults",
			Orderer: profiles.Orderer{
				Orderer:       "<<: *OrdererDefaults",
				Organizations: []string{"- *Orderer"},
			},
			Consortiums: map[string]profiles.Organizations{
				"Consortium": {
					Organizations: _organizations,
				},
			},
			Application: profiles.Application{
				Application:   "<<: *ApplicationDefaults",
				Organizations: _organizations,
			},
		},
	}

	for _, channel := range config.Channels {
		_configtx.Profiles[channel.Name] = profiles.ChannelProfile{
			Channel:    "<<: *ChannelDefaults",
			Consortium: "Consortium",
			Application: profiles.Application{
				Application:   "<<: *ApplicationDefaults",
				Organizations: channel.Organizations,
			},
		}
	}

	return _configtx
}

func UpdateAnchors(content string, organizations []string) string {
	var pattern string
	var re *regexp.Regexp

	/* Capabilities */
	content = strings.Replace(content, `Application:`, "Application: &ApplicationCapabilities", 1)
	content = strings.Replace(content, `Orderer:`, "Orderer: &OrdererCapabilities", 1)
	content = strings.Replace(content, `Channel:`, "Channel: &ChannelCapabilities", 1)

	/* Organizations */
	pattern = `(?m)^Organizations:`
	re = regexp.MustCompile(pattern)
	content = re.ReplaceAllString(content, "\nOrganizations:")

	pattern = `-?\s*Name:\s*(\S+)`
	re = regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		content = strings.ReplaceAll(content, fmt.Sprintf("- Name: %s", match[1]), fmt.Sprintf("- &%s\n    Name: %s", match[1], match[1]))
	}

	for _, organization := range organizations {
		re = regexp.MustCompile(fmt.Sprintf(`(?m)^  - &%s`, organization))
		content = re.ReplaceAllString(content, fmt.Sprintf("\n  - &%s", organization))
	}

	/* Orderer */
	content = strings.Replace(content, `- '<<: *OrdererCapabilities'`, "<<: *OrdererCapabilities", -1)
	pattern = `(?m)^Orderer:`
	re = regexp.MustCompile(pattern)
	content = re.ReplaceAllString(content, "\nOrderer: &OrdererDefaults")

	/* Application */
	content = strings.Replace(content, `- '<<: *ApplicationCapabilities'`, "<<: *ApplicationCapabilities", -1)

	for _, organization := range organizations {
		content = strings.ReplaceAll(content, fmt.Sprintf("- %s", organization), fmt.Sprintf("- *%s", organization))
	}

	pattern = `(?m)^Application:`
	re = regexp.MustCompile(pattern)
	content = re.ReplaceAllString(content, "\nApplication: &ApplicationDefaults")

	/* Channel */
	content = strings.Replace(content, `- '<<: *ChannelCapabilities'`, "<<: *ChannelCapabilities", -1)

	pattern = `(?m)^Channel:`
	re = regexp.MustCompile(pattern)
	content = re.ReplaceAllString(content, "\nChannel: &ChannelDefaults")

	/* Profiles */
	content = strings.Replace(content, `Channel: '<<: *ChannelDefaults'`, "<<: *ChannelDefaults", -1)
	content = strings.Replace(content, `Application: '<<: *ApplicationDefaults'`, "<<: *ApplicationDefaults", -1)
	content = strings.Replace(content, `Orderer: '<<: *OrdererDefaults'`, "<<: *OrdererDefaults", -1)
	content = strings.Replace(content, `- '- *Orderer'`, "- *Orderer", -1)

	pattern = `(?m)^Profiles:`
	re = regexp.MustCompile(pattern)
	content = re.ReplaceAllString(content, "\nProfiles:")

	return content
}
