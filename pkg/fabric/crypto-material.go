package fabric

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/internal/constants"
)

// Deprecated
func (f *Fabric) GenerateCryptoMaterial() error {
	for _, organization := range f.config.Organizations {
		tools := fmt.Sprintf("%s/%s/tools.yml", f.config.Output, organization.Domain)

		fmt.Printf("\n=========== Generating crypto materials to %s ===========\n", organization.Name)

		containerName := buildToolsContainerName(organization)

		args := []string{
			"compose", "-f", f.network, "-f", tools, "run", "--rm", "-T", containerName,
			"cryptogen", "generate",
			fmt.Sprintf("--config=%s/crypto-config.yml", constants.DEFAULT_FABRIC_DIRECTORY),
			fmt.Sprintf("--output=%s/crypto-materials", constants.DEFAULT_FABRIC_DIRECTORY),
		}

		if err := f.executor.ExecCommand("docker", args...); err != nil {
			return fmt.Errorf("Error when generating the crypto materials for the organization %s: %v\n", organization.Name, err)
		}
	}

	return nil
}
