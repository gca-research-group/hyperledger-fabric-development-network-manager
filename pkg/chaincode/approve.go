package chaincode

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/compose"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/network"
)

func (c *Chaincode) Approve() error {
	ordererAddress, caFile := network.ResolveOrdererTLSConnection(c.config.Organizations)
	signaturePolicy := ""

	for _, organization := range c.config.Organizations {
		if signaturePolicy == "" {
			signaturePolicy = fmt.Sprintf("'%sMSP.peer'", organization.Name)
			continue
		}

		signaturePolicy = strings.Join([]string{signaturePolicy, fmt.Sprintf("'%sMSP.peer'", organization.Name)}, ",")
	}

	for _, organization := range c.config.Organizations {
		composefile := compose.ResolveToolsDockerComposeFile(c.config.Output, organization.Domain)
		containerName := compose.ResolveToolsContainerName(organization)

		for _, channel := range c.config.Channels {
			for _, chaincode := range channel.Chaincodes {
				name := filepath.Base(chaincode.Path)
				basePath := fmt.Sprintf("/chaincodes/%[1]s", name)
				tarfile := fmt.Sprintf("%s/%s.tar.gz", basePath, name)

				args := []string{
					"compose", "-f", c.network, "-f", composefile, "run", "--rm", "-T", containerName,
					"peer", "lifecycle", "chaincode", "calculatepackageid", tarfile,
				}

				packageId, _ := c.executor.OutputCommand("docker", args...)

				args = []string{
					"compose", "-f", c.network, "-f", composefile, "run", "--rm", "-T", containerName,
					"peer", "lifecycle", "chaincode", "queryapproved", "--channelID", strings.ToLower(channel.Name), "--name", name,
				}

				approved, _ := c.executor.OutputCommand("docker", args...)

				if strings.Contains(strings.TrimSpace(string(approved)), "Approved") {
					continue
				}

				version := "1.0"
				sequence := "1"

				args = []string{
					"compose", "-f", c.network, "-f", composefile, "run", "--rm", "-T", containerName,
					"peer", "lifecycle", "chaincode", "approveformyorg",
					"--channelID", strings.ToLower(channel.Name),
					"--name", name,
					"--version", version,
					"--sequence", sequence,
					"--package-id", string(packageId),
					"--signature-policy", fmt.Sprintf("AND(%s)", signaturePolicy),
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
