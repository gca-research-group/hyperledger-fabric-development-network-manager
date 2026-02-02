package fabric

import (
	"fmt"
)

func (f *Fabric) GenerateCryptoMaterial() error {
	for _, organization := range f.config.Organizations {
		tools := fmt.Sprintf("%s/%s/tools.yml", f.config.Output, organization.Domain)

		fmt.Printf("\n=========== Generating crypto materials to %s ===========\n", organization.Name)

		containerName := buildToolsContainerName(organization)

		var args []string

		args = append(args, "compose", "-f", f.network, "-f", tools, "run", "--rm", "-T", containerName)
		args = append(args, "cryptogen", "generate")
		args = append(args, fmt.Sprintf("--config=%s/crypto-config.yml", DEFAULT_FABRIC_DIRECTORY))
		args = append(args, fmt.Sprintf("--output=%s/crypto-materials", DEFAULT_FABRIC_DIRECTORY))

		if err := f.executor.ExecCommand("docker", args...); err != nil {
			return fmt.Errorf("Error when generating the crypto materials for the organization %s: %v\n", organization.Name, err)
		}
	}

	return nil
}
