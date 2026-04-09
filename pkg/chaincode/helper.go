package chaincode

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gca-research-group/fabric-network-orchestrator/pkg/compose"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/config"
)

func ResolveFilename(chaincode config.Chaincode) string {
	filename := filepath.Base(chaincode.Path)
	return filename
}

func ResolveLabel(chaincode config.Chaincode) string {
	return fmt.Sprintf("%[1]s_%[2]s", chaincode.Name, ResolveChaincodeVersion(chaincode))
}

func ResolveChaincodePath(chaincode config.Chaincode) string {
	return fmt.Sprintf("/chaincodes/%[1]s", filepath.Base(chaincode.Path))
}

func ResolveChaincodeTar(chaincode config.Chaincode) string {
	return fmt.Sprintf("%[1]s/%[2]s.tar.gz", ResolveChaincodePath(chaincode), ResolveLabel(chaincode))
}

func ResolveCollectionsConfig(chaincode config.Chaincode) string {
	return fmt.Sprintf("%s/%s", ResolveChaincodePath(chaincode), filepath.Base(chaincode.CollectionsConfig))
}

func (c *Chaincode) QueryPackageId(organization config.Organization, tarfile string) string {

	args := []string{
		"peer", "lifecycle", "chaincode", "calculatepackageid", tarfile,
	}

	packageId, _ := c.ExecInTools(organization, args)

	return strings.TrimSpace(string(packageId))
}

func (c *Chaincode) ChaincodeFileExists(organization config.Organization, tarfile string) bool {
	args := []string{
		"sh", "-c", fmt.Sprintf("[ -f %s ]", tarfile),
	}

	_, err := c.ExecInTools(organization, args)

	return err == nil
}

func (c *Chaincode) IsChaincodeApproved(organization config.Organization, channelID string, chaincode config.Chaincode, version string) bool {
	sequence := strconv.Itoa(c.QueryCurrentApprovedSequence(organization, channelID, chaincode.Name))

	args := []string{"peer", "lifecycle", "chaincode", "checkcommitreadiness",
		"-C", channelID,
		"-n", chaincode.Name,
		"--version", version,
		"--sequence", sequence,
		"--output", "json",
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

	output, err := c.ExecInTools(organization, args)

	if err != nil {
		if strings.Contains(err.Error(), "but new definition must be sequence") {
			re := regexp.MustCompile(`(\d+)\D*$`)
			matches := re.FindStringSubmatch(err.Error())

			if len(matches) > 1 {
				sequence = matches[1]
				for i, v := range args {
					if v == "--sequence" && i+1 < len(args) {
						args[i+1] = sequence
						break
					}
				}

				output, err = c.ExecInTools(organization, args)

				if err != nil {
					panic(err)
				}
			}
		}
	}

	var result map[string]map[string]bool
	json.Unmarshal(output, &result)

	return result["approvals"][config.ResolveOrganizationMSPID(organization)]
}

func (c *Chaincode) IsChaincodeCommitted(organization config.Organization, channelID string, name string, version string) bool {
	args := []string{
		"peer", "lifecycle", "chaincode", "querycommitted", "--channelID", channelID, "--name", name, "--output", "json",
	}

	output, _ := c.ExecInTools(organization, args)

	var result struct {
		Sequence  int             `json:"sequence"`
		Version   string          `json:"version"`
		Approvals map[string]bool `json:"approvals"`
	}

	json.Unmarshal(output, &result)

	currentApprovedSequence := c.QueryCurrentApprovedSequence(organization, channelID, name)

	isApproved := result.Approvals[config.ResolveOrganizationMSPID(organization)]
	isVersionUpToDate := result.Version == version
	isSequenceUpToDate := result.Sequence == currentApprovedSequence

	return isApproved && isVersionUpToDate && isSequenceUpToDate
}

func (c *Chaincode) QueryCurrentApprovedSequence(organization config.Organization, channelID string, name string) int {
	args := []string{
		"peer", "lifecycle", "chaincode", "queryapproved", "--channelID", channelID, "--name", name, "--output", "json",
	}

	output, _ := c.ExecInTools(organization, args)

	var result struct {
		Sequence int `json:"sequence"`
	}

	json.Unmarshal(output, &result)

	return result.Sequence
}

func (c *Chaincode) ComputeCurrentApprovedSequence(organization config.Organization, channelID string, name string) string {
	sequence := c.QueryCurrentApprovedSequence(organization, channelID, name)

	return strconv.Itoa(sequence + 1)
}

func (c *Chaincode) QueryCurrentCommittedSequence(organization config.Organization, channelID string, name string) int {
	args := []string{
		"peer", "lifecycle", "chaincode", "querycommitted", "--channelID", channelID, "--name", name, "--output", "json",
	}

	output, _ := c.ExecInTools(organization, args)

	var result struct {
		Sequence int `json:"sequence"`
	}

	json.Unmarshal(output, &result)

	return result.Sequence
}

func (c *Chaincode) QueryInstalled(organization config.Organization) string {

	args := []string{
		"peer", "lifecycle", "chaincode", "queryinstalled",
	}

	output, _ := c.ExecInTools(organization, args)

	return strings.TrimSpace(string(output))
}

func (c *Chaincode) IsChaincodeInstalled(organization config.Organization, tarfile string) bool {
	packageId := c.QueryPackageId(organization, tarfile)
	installed := c.QueryInstalled(organization)
	return strings.Contains(installed, packageId)
}

func (c *Chaincode) ExecInTools(organization config.Organization, args []string) ([]byte, error) {
	containerName := compose.ResolveToolsContainerName(organization)

	baseArgs := []string{
		"exec", containerName,
	}

	return c.executor.OutputCommand("docker", append(baseArgs, args...)...)
}

func ResolveChaincodeVersion(chaincode config.Chaincode) string {
	version := chaincode.Version

	if version == "" {
		version = "1.0"
	}
	return version
}
