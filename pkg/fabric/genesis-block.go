package fabric

import (
	"fmt"
	"strings"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg"
)

func (f *Fabric) GenerateGenesisBlock() error {
	for _, organization := range f.config.Organizations {
		f.dockerRenderer.RenderToolsWithMSP(organization)

		if organization.Bootstrap {
			for _, profile := range f.config.Profiles {
				fmt.Printf("\n=========== Generating orderer genesis block to %s ===========\n", organization.Name)

				containerName := buildToolsContainerName(organization)
				tools := fmt.Sprintf("%s/%s/tools.yml", f.config.Output, organization.Domain)

				var args []string

				args = append(args, "compose", "-f", f.network, "-f", tools, "run", "--rm", "-T", containerName)
				args = append(args, "configtxgen")
				args = append(args, "-outputBlock", fmt.Sprintf("%s/channel/%s.block", DEFAULT_FABRIC_DIRECTORY, strings.ToLower(profile.Name)))
				args = append(args, "-profile", profile.Name)
				args = append(args, "-channelID", strings.ToLower(profile.Name))
				args = append(args, "-configPath", fmt.Sprintf("%s/", DEFAULT_FABRIC_DIRECTORY))

				if err := f.executor.ExecCommand("docker", args...); err != nil {
					return fmt.Errorf("Error when generating the genesis block for the organization %s: %v", organization.Name, err)
				}
			}
		}
	}

	return nil
}

func (f *Fabric) FetchGenesisBlock() error {
	var orderer pkg.Orderer
	var ordererDomain string

	for _, organization := range f.config.Organizations {
		if len(organization.Orderers) > 0 {
			orderer = organization.Orderers[0]
			ordererDomain = organization.Domain
			break
		}
	}

	ordererAddress := fmt.Sprintf("%s.%s:%d", orderer.Hostname, ordererDomain, orderer.Port)
	caFile := fmt.Sprintf("%s/crypto-materials/ordererOrganizations/%s/orderers/%s.%s/tls/ca.crt", DEFAULT_FABRIC_DIRECTORY, ordererDomain, orderer.Hostname, ordererDomain)

	for _, organization := range f.config.Organizations {
		if organization.Bootstrap {
			continue
		}

		tools := fmt.Sprintf("%s/%s/tools.yml", f.config.Output, organization.Domain)
		for _, profile := range f.config.Profiles {
			containerName := buildToolsContainerName(organization)
			block := fmt.Sprintf("%s/channel/%s.block", DEFAULT_FABRIC_DIRECTORY, strings.ToLower(profile.Name))

			var args []string
			args = append(args, "compose", "-f", f.network, "-f", tools, "run", "--rm", "-T", containerName)
			args = append(args, "peer", "channel", "fetch", "0", block, "-c", strings.ToLower(profile.Name), "-o", ordererAddress, "--tls", "--cafile", caFile)

			if err := f.executor.ExecCommand("docker", args...); err != nil {
				return fmt.Errorf("Error when fetching the orderer %s of the organization %s to the channel %s: %v", orderer.Name, organization.Name, profile.Name, err)
			}
		}
	}

	return nil
}
