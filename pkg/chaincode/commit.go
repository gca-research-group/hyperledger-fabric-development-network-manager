package chaincode

import (
	"fmt"
	"strconv"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/network"
)

func (c *Chaincode) Commit() error {
	ordererAddress, caFile := network.ResolveOrdererTLSConnection(c.config.Organizations)
	organization := c.config.Organizations[0]

	for _, channel := range c.config.Channels {
		channelID := network.ResolveChannelID(channel)
		for _, chaincode := range channel.Chaincodes {
			version := ResolveChaincodeVersion(chaincode)
			sequence := c.QueryCurrentApprovedSequence(organization, channelID, chaincode.Name)

			if c.IsChaincodeCommitted(organization, channelID, chaincode.Name, version) {
				continue
			}

			peers := network.ResolvePeersTLSConnection(c.config.Organizations)

			args := []string{
				"peer", "lifecycle", "chaincode", "commit",
				"--channelID", channelID,
				"--name", chaincode.Name,
				"--version", version,
				"--sequence", strconv.Itoa(sequence),
				"--orderer", ordererAddress,
				"--tls", "--cafile", caFile,
			}

			if chaincode.SignaturePolicy != "" {
				args = append(args, "--signature-policy", chaincode.SignaturePolicy)
			}

			if chaincode.ChannelConfigPolicy != "" {
				args = append(args, "--channel-config-policy", chaincode.ChannelConfigPolicy)
			}

			if chaincode.CollectionsConfig != "" {
				args = append(args, "--collections-config", ResolveCollectionsConfig(chaincode))
			}

			for _, peer := range peers {
				args = append(args, []string{"--peerAddresses", peer[0], "--tlsRootCertFiles", peer[1]}...)
			}

			_, err := c.ExecInTools(organization, args)

			if err != nil {
				return fmt.Errorf("Error when committing the chaincode %s in the organization %s: %v", chaincode.Name, organization.Name, err)
			}
		}
	}

	return nil
}
