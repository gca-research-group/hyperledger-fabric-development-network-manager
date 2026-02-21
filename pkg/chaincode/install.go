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
			basePath := fmt.Sprintf("/chaincodes/%[1]s", name)
			tarfile := fmt.Sprintf("%s/%s.tar.gz", basePath, name)

			args := []string{
				"compose", "-f", c.network, "-f", composefile, "run", "--rm", "-T", containerName,
				"peer", "lifecycle", "chaincode", "calculatepackageid", tarfile,
			}

			packageId, _ := c.executor.OutputCommand("docker", args...)

			args = []string{
				"compose", "-f", c.network, "-f", composefile, "run", "--rm", "-T", containerName,
				"peer", "lifecycle", "chaincode", "queryinstalled",
			}

			installed, _ := c.executor.OutputCommand("docker", args...)

			if strings.Contains(strings.TrimSpace(string(installed)), strings.TrimSpace(string(packageId))) {
				continue
			}

			args = []string{
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
