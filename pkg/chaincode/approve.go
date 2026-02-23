package chaincode

import (
	"fmt"
	"path/filepath"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/compose"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/network"
)

func (c *Chaincode) Approve() error {
	ordererAddress, caFile := network.ResolveOrdererTLSConnection(c.config.Organizations)
	signaturePolicy := c.ResolveSignaturePolicy()

	for _, organization := range c.config.Organizations {
		composefile := compose.ResolveToolsDockerComposeFile(c.config.Output, organization.Domain)
		containerName := compose.ResolveToolsContainerName(organization)

		for _, channel := range c.config.Channels {
			for _, chaincode := range channel.Chaincodes {
				name := filepath.Base(chaincode.Path)
				version := DEFAULT_CHAINCODE_VERSION
				sequence := DEFAULT_CHAINCODE_SEQUENCE
				tarfile := ResolveChaincode(name, version)

				if c.IsChaincodeApproved(composefile, containerName, network.ResolveChannelID(channel), name) {
					continue
				}

				packageId := c.QueryPackageId(composefile, containerName, tarfile)

				args := []string{
					"compose", "-f", c.network, "-f", composefile, "run", "--rm", "-T", containerName,
					"peer", "lifecycle", "chaincode", "approveformyorg",
					"--channelID", network.ResolveChannelID(channel),
					"--name", name,
					"--version", version,
					"--sequence", sequence,
					"--package-id", packageId,
					"--signature-policy", signaturePolicy,
					"--orderer", ordererAddress,
					"--tls", "--cafile", caFile,
				}

				if err := c.executor.ExecCommand("docker", args...); err != nil {
					return fmt.Errorf("Error when approving the chaincode %s in the organization %s: %v", name, organization.Name, err)
				}
			}
		}
	}

	return nil
}
