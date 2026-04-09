package chaincode

import (
	"fmt"
)

func (c *Chaincode) Package() error {

	organization := c.config.Organizations[0]

	for _, chaincode := range c.config.Chaincodes {

		name := chaincode.Name
		label := ResolveLabel(chaincode)
		tarfile := ResolveChaincodeTar(chaincode)
		chaincodePath := ResolveChaincodePath(chaincode)

		steps := []struct {
			name    string
			message string
			args    []string
		}{
			{"Initialize", "Error when initializing the chaincode module %s: %v", []string{
				"sh", "-c", fmt.Sprintf("cd %s && [ -f go.mod ] || go mod init %s; go mod tidy", chaincodePath, name),
			}},
			{"Package", "Error when packaging the chaincode %s: %v", []string{
				"peer", "lifecycle", "chaincode", "package", tarfile,
				"--path", chaincodePath,
				"--lang", "golang",
				"--label", label,
			}},
		}

		for _, step := range steps {
			fmt.Printf(">>> Step: %s\n", step.name)
			_, err := c.ExecInTools(organization, step.args)
			if err != nil {
				return fmt.Errorf(step.message, name, err)
			}
		}
	}

	return nil
}
