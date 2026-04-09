package chaincode

import (
	"fmt"
)

func (c *Chaincode) Install() error {

	for _, organization := range c.config.Organizations {
		for _, chaincode := range c.config.Chaincodes {
			name := chaincode.Name
			tarfile := ResolveChaincodeTar(chaincode)

			if c.IsChaincodeInstalled(organization, tarfile) {
				continue
			}

			args := []string{
				"peer", "lifecycle", "chaincode", "install", tarfile,
			}

			_, err := c.ExecInTools(organization, args)

			if err != nil {
				return fmt.Errorf("Error when installing the chaincode %s in the organization %s: %v", name, organization.Name, err)
			}
		}
	}

	return nil
}
