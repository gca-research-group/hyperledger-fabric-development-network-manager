package chaincode

import (
	"fmt"
	"strings"
)

const DEFAULT_CHAINCODE_VERSION = "1.0"
const DEFAULT_CHAINCODE_SEQUENCE = "1"

func (c *Chaincode) ResolveSignaturePolicy() string {
	signaturePolicy := ""

	for _, organization := range c.config.Organizations {
		if signaturePolicy == "" {
			signaturePolicy = fmt.Sprintf("'%sMSP.peer'", organization.Name)
			continue
		}

		signaturePolicy = strings.Join([]string{signaturePolicy, fmt.Sprintf("'%sMSP.peer'", organization.Name)}, ",")
	}

	return fmt.Sprintf("AND(%s)", signaturePolicy)
}

func ResolveLabel(name string, version string) string {
	return fmt.Sprintf("%[1]s_%[2]s", name, version)
}

func ResolveChaincodePath(name string) string {
	return fmt.Sprintf("/chaincodes/%[1]s", name)
}

func ResolveChaincode(name string, version string) string {
	return fmt.Sprintf("%[1]s/%[2]s.tar.gz", ResolveChaincodePath(name), ResolveLabel(name, version))
}

func ResolveChecksum(name string, version string) string {
	return fmt.Sprintf("%[1]s/%[2]s.sha256sum", ResolveChaincodePath(name), ResolveLabel(name, version))
}

func (c *Chaincode) QueryPackageId(composefile string, containerName string, chaincodeTarFile string) string {

	args := []string{
		"compose", "-f", c.network, "-f", composefile, "run", "--rm", "-T", containerName,
		"peer", "lifecycle", "chaincode", "calculatepackageid", chaincodeTarFile,
	}

	packageId, _ := c.executor.OutputCommand("docker", args...)

	return strings.TrimSpace(string(packageId))
}

func (c *Chaincode) IsChaincodeUpToDate(composefile string, containerName string, basePath string, name string, version string) bool {
	args := []string{
		"compose", "-f", c.network, "-f", composefile, "run", "--rm", "-T", containerName,
		"sh", "-c", fmt.Sprintf("sha256sum -c %[1]s", ResolveChecksum(name, version)),
	}

	output, _ := c.executor.OutputCommand("docker", args...)

	return strings.Contains(strings.TrimSpace(string(output)), "OK")
}

func (c *Chaincode) ChaincodeFileExists(composefile string, containerName string, tarfile string) bool {
	args := []string{
		"compose", "-f", c.network, "-f", composefile, "run", "--rm", "-T", containerName,
		"sh", "-c", fmt.Sprintf("[ -f %s ]", tarfile),
	}

	err := c.executor.ExecCommand("docker", args...)

	return err == nil
}

func (c *Chaincode) IsChaincodeApproved(composefile string, containerName string, channelID string, name string) bool {
	args := []string{
		"compose", "-f", c.network, "-f", composefile, "run", "--rm", "-T", containerName,
		"peer", "lifecycle", "chaincode", "queryapproved", "--channelID", channelID, "--name", name,
	}

	approved, _ := c.executor.OutputCommand("docker", args...)

	return strings.Contains(strings.TrimSpace(string(approved)), "Approved")
}

func (c *Chaincode) IsChaincodeCommitted(composefile string, containerName string, channelID string, name string) bool {
	args := []string{
		"compose", "-f", c.network, "-f", composefile, "run", "--rm", "-T", containerName,
		"peer", "lifecycle", "chaincode", "querycommitted", "--channelID", channelID, "--name", name,
	}

	approved, _ := c.executor.OutputCommand("docker", args...)

	return strings.Contains(strings.TrimSpace(string(approved)), "Approvals")
}
