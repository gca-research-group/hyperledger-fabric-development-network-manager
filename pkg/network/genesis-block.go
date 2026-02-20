package network

import (
	"fmt"
	"strings"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/constants"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/compose"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/config"
)

func (f *Network) GenerateGenesisBlock() error {
	for _, organization := range f.config.Organizations {

		if organization.Bootstrap {
			for _, channel := range f.config.Channels {
				fmt.Printf("\n=========== Generating orderer genesis block to %s ===========\n", organization.Name)

				tools := compose.ResolveToolsDockerComposeFile(f.config.Output, organization.Domain)
				containerName := compose.ResolveToolsContainerName(organization)

				args := []string{
					"compose", "-f", f.network, "-f", tools, "run", "--rm", "-T", containerName,
					"configtxgen",
					"-outputBlock", fmt.Sprintf("%s/channels/%s.block", constants.DEFAULT_FABRIC_DIRECTORY, strings.ToLower(channel.Name)),
					"-profile", channel.Profile.Name,
					"-channelID", strings.ToLower(channel.Name),
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

func (f *Network) FetchGenesisBlock() error {
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

		var channels []config.Channel

		for _, channel := range f.config.Channels {
			for _, organizationName := range channel.Profile.Organizations {
				if organizationName == organization.Name {
					channels = append(channels, channel)
					break
				}
			}
		}

		tools := compose.ResolveToolsDockerComposeFile(f.config.Output, organization.Domain)
		for _, channel := range channels {
			containerName := compose.ResolveToolsContainerName(organization)
			block := fmt.Sprintf("%s/channels/%s.block", constants.DEFAULT_FABRIC_DIRECTORY, strings.ToLower(channel.Name))

			args := []string{
				"compose", "-f", f.network, "-f", tools, "run", "--rm", "-T", containerName,
				"peer", "channel", "fetch", "0", block, "-c", strings.ToLower(channel.Name), "-o", ordererAddress, "--tls", "--cafile", caFile,
			}

			if err := f.executor.ExecCommand("docker", args...); err != nil {
				return fmt.Errorf("Error when fetching the orderer %s of the organization %s to the channel %s: %v", orderer.Name, organization.Name, channel.Name, err)
			}
		}
	}

	return nil
}
