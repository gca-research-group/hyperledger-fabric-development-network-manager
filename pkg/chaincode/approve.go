package chaincode

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/network"
)

func (c *Chaincode) Approve() error {
	ordererAddress, caFile := network.ResolveOrdererTLSConnection(c.config.Organizations)
	for _, organization := range c.config.Organizations {
		for _, channel := range c.config.Channels {
			channelID := network.ResolveChannelID(channel)

			for _, chaincode := range channel.Chaincodes {
				name := chaincode.Name
				version := ResolveChaincodeVersion(chaincode)
				sequence := c.ComputeCurrentApprovedSequence(organization, channelID, name)
				tarfile := ResolveChaincodeTar(chaincode)

				if c.IsChaincodeApproved(organization, channelID, chaincode, version) {
					continue
				}

				packageId := c.QueryPackageId(organization, tarfile)

				args := []string{
					"peer", "lifecycle", "chaincode", "approveformyorg",
					"--channelID", channelID,
					"--name", name,
					"--version", version,
					"--sequence", sequence,
					"--package-id", packageId,
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

				_, err := c.ExecInTools(organization, args)

				if err != nil {
					return fmt.Errorf("Error when approving the chaincode %s in the organization %s: %v", chaincode.Name, organization.Name, err)
				}
			}
		}
	}

	return nil
}
