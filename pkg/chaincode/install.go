package chaincode

import (
	"fmt"

	"github.com/gca-research-group/fabric-network-orchestrator/pkg/config"
)

func (c *Chaincode) Install() error {

	chaincodes := config.ResolveChaincodes(*c.config)

	for _, organization := range c.config.Organizations {
		for _, chaincode := range chaincodes {
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
