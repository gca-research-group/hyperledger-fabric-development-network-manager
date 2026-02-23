package chaincode

import (
	"fmt"
	"path/filepath"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/compose"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/network"
)

func (c *Chaincode) Commit() error {
	ordererAddress, caFile := network.ResolveOrdererTLSConnection(c.config.Organizations)
	signaturePolicy := c.ResolveSignaturePolicy()

	organization := c.config.Organizations[0]

	composefile := compose.ResolveToolsDockerComposeFile(c.config.Output, organization.Domain)
	containerName := compose.ResolveToolsContainerName(organization)

	for _, channel := range c.config.Channels {
		for _, chaincode := range channel.Chaincodes {
			name := filepath.Base(chaincode.Path)
			version := DEFAULT_CHAINCODE_VERSION
			sequence := DEFAULT_CHAINCODE_SEQUENCE

			if c.IsChaincodeCommitted(composefile, containerName, network.ResolveChannelID(channel), name) {
				continue
			}

			peers := network.ResolvePeersTLSConnection(c.config.Organizations)

			args := []string{
				"compose", "-f", c.network, "-f", composefile, "run", "--rm", "-T", containerName,
				"peer", "lifecycle", "chaincode", "commit",
				"--channelID", network.ResolveChannelID(channel),
				"--name", name,
				"--version", version,
				"--sequence", sequence,
				"--signature-policy", signaturePolicy,
				"--orderer", ordererAddress,
				"--tls", "--cafile", caFile,
			}

			for _, peer := range peers {
				args = append(args, []string{"--peerAddresses", peer[0], "--tlsRootCertFiles", peer[1]}...)
			}

			if err := c.executor.ExecCommand("docker", args...); err != nil {
				return fmt.Errorf("Error when approving the chaincode %s in the organization %s: %v", name, organization.Name, err)
			}
		}
	}

	return nil
}
