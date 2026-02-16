package fabric

import (
	"fmt"
	"strings"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/constants"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/config"
)

func (f *Fabric) GenerateGenesisBlock() error {
	for _, organization := range f.config.Organizations {

		if organization.Bootstrap {
			for _, profile := range f.config.Profiles {
				fmt.Printf("\n=========== Generating orderer genesis block to %s ===========\n", organization.Name)

				containerName := buildToolsContainerName(organization)
				tools := fmt.Sprintf("%s/%s/tools.yml", f.config.Output, organization.Domain)

				args := []string{
					"compose", "-f", f.network, "-f", tools, "run", "--rm", "-T", containerName,
					"configtxgen",
					"-outputBlock", fmt.Sprintf("%s/channel/%s.block", constants.DEFAULT_FABRIC_DIRECTORY, strings.ToLower(profile.Name)),
					"-profile", profile.Name,
					"-channelID", strings.ToLower(profile.Name),
					"-configPath", fmt.Sprintf("%s/", constants.DEFAULT_FABRIC_DIRECTORY),
				}

				if err := f.executor.ExecCommand("docker", args...); err != nil {
					return fmt.Errorf("Error when generating the genesis block for the organization %s: %v", organization.Name, err)
				}
			}
		}
	}

	return nil
}

func (f *Fabric) FetchGenesisBlock() error {
	var orderer config.Orderer
	var ordererDomain string
	var ordererPort int

	for _, organization := range f.config.Organizations {
		if len(organization.Orderers) > 0 {
			orderer = organization.Orderers[0]
			ordererDomain = organization.Domain
			ordererPort = orderer.Port

			if ordererPort == 0 {
				ordererPort = 7050
			}
			break
		}
	}

	ordererAddress := fmt.Sprintf("%s.%s:%d", orderer.Subdomain, ordererDomain, ordererPort)
	caFile := fmt.Sprintf("%[1]s/%[2]s/ordererOrganizations/%[2]s/orderers/%[3]s.%[2]s/tls/ca.crt", constants.DEFAULT_FABRIC_DIRECTORY, ordererDomain, orderer.Subdomain)

	for _, organization := range f.config.Organizations {
		if organization.Bootstrap {
			continue
		}

		tools := fmt.Sprintf("%s/%s/tools.yml", f.config.Output, organization.Domain)
		for _, profile := range f.config.Profiles {
			containerName := buildToolsContainerName(organization)
			block := fmt.Sprintf("%s/channel/%s.block", constants.DEFAULT_FABRIC_DIRECTORY, strings.ToLower(profile.Name))

			args := []string{
				"compose", "-f", f.network, "-f", tools, "run", "--rm", "-T", containerName,
				"peer", "channel", "fetch", "0", block, "-c", strings.ToLower(profile.Name), "-o", ordererAddress, "--tls", "--cafile", caFile,
			}

			if err := f.executor.ExecCommand("docker", args...); err != nil {
				return fmt.Errorf("Error when fetching the orderer %s of the organization %s to the channel %s: %v", orderer.Name, organization.Name, profile.Name, err)
			}
		}
	}

	return nil
}
