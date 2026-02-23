package chaincode

import (
	"fmt"
	"path/filepath"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/compose"
)

func (c *Chaincode) Package() error {

	organization := c.config.Organizations[0]

	composefile := compose.ResolveToolsDockerComposeFile(c.config.Output, organization.Domain)
	containerName := compose.ResolveToolsContainerName(organization)

	for _, chaincode := range c.config.Chaincodes {
		name := filepath.Base(chaincode.Path)
		version := DEFAULT_CHAINCODE_VERSION
		label := ResolveLabel(name, version)
		basePath := ResolveChaincodePath(name)
		tarfile := ResolveChaincode(name, version)

		isChaincodeUpToDate := c.IsChaincodeUpToDate(composefile, containerName, basePath, name, version)
		chaincodeFileExists := c.ChaincodeFileExists(composefile, containerName, tarfile)

		if chaincodeFileExists && isChaincodeUpToDate {
			continue
		}

		args := []string{
			"compose", "-f", c.network, "-f", composefile, "run", "--rm", "-T", containerName,
			"sh", "-c", fmt.Sprintf("cd %s && [ -f go.mod ] || go mod init %s; go mod tidy", basePath, name),
		}

		if err := c.executor.ExecCommand("docker", args...); err != nil {
			return fmt.Errorf("Error when initializing the chaincode module %s: %v", name, err)
		}

		args = []string{
			"compose", "-f", c.network, "-f", composefile, "run", "--rm", "-T", containerName,
			"sh", "-c", fmt.Sprintf("sha256sum %[1]s/%[2]s.go > %[3]s", basePath, name, ResolveChecksum(name, version)),
		}

		if err := c.executor.ExecCommand("docker", args...); err != nil {
			return fmt.Errorf("Error when computing the chaincode checksum %s: %v", name, err)
		}

		args = []string{
			"compose", "-f", c.network, "-f", composefile, "run", "--rm", "-T", containerName,
			"peer", "lifecycle", "chaincode", "package", tarfile,
			"--path", basePath,
			"--lang", "golang",
			"--label", label,
		}

		if err := c.executor.ExecCommand("docker", args...); err != nil {
			return fmt.Errorf("Error when packaging the chaincode %s: %v", name, err)
		}
	}

	return nil
}
