package chaincode

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/compose"
)

func (c *Chaincode) Install() error {

	for _, organization := range c.config.Organizations {
		composefile := compose.ResolveToolsDockerComposeFile(c.config.Output, organization.Domain)
		containerName := compose.ResolveToolsContainerName(organization)

		for _, chaincode := range c.config.Chaincodes {
			name := filepath.Base(chaincode.Path)
			tarfile := ResolveChaincode(name, DEFAULT_CHAINCODE_VERSION)

			if c.IsChaincodeInstalled(composefile, containerName, tarfile) {
				continue
			}

			args := []string{
				"compose", "-f", c.network, "-f", composefile, "run", "--rm", "-T", containerName,
				"peer", "lifecycle", "chaincode", "install", tarfile,
			}

			if err := c.executor.ExecCommand("docker", args...); err != nil {
				return fmt.Errorf("Error when installing the chaincode %s in the organization %s: %v", name, organization.Name, err)
			}
		}
	}

	return nil
}

func (c *Chaincode) QueryInstalled(toolsComposeFile, containerName string) string {

	args := []string{
		"compose", "-f", c.network, "-f", toolsComposeFile, "run", "--rm", "-T", containerName,
		"peer", "lifecycle", "chaincode", "queryinstalled",
	}

	installed, _ := c.executor.OutputCommand("docker", args...)

	return strings.TrimSpace(string(installed))
}

func (c *Chaincode) IsChaincodeInstalled(toolsComposeFile, containerName string, chaincodeTarFile string) bool {
	packageId := c.QueryPackageId(toolsComposeFile, containerName, chaincodeTarFile)
	installed := c.QueryInstalled(toolsComposeFile, containerName)
	return strings.Contains(installed, packageId)
}
