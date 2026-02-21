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
	ordererAddress, caFile := ResolveOrdererTLSConnection(f.config.Organizations)

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
				return fmt.Errorf("Error when fetching genesis block for the organization %s to the channel %s: %v", organization.Name, channel.Name, err)
			}
		}
	}

	return nil
}
