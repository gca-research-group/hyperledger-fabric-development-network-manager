package network

import (
	"fmt"
	"strings"
	"time"

	"github.com/gca-research-group/fabric-network-orchestrator/internal/constants"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/compose"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/config"
)

func (f *Network) GenerateGenesisBlock() error {
	for _, organization := range f.config.Organizations {

		if organization.Bootstrap {
			containerName := compose.ResolveToolsContainerName(organization)

			for _, channel := range f.config.Channels {
				fmt.Printf("\n=========== Generating orderer genesis block to %s ===========\n", organization.Name)

				script := strings.Join(
					[]string{
						"configtxgen",
						"-outputBlock", fmt.Sprintf("%s/channels/%s.block", constants.DEFAULT_FABRIC_DIRECTORY, ResolveChannelID(channel)),
						"-profile", channel.Profile.Name,
						"-channelID", ResolveChannelID(channel),
						"-configPath", fmt.Sprintf("%s/", constants.DEFAULT_FABRIC_DIRECTORY),
					}, " ",
				)

				args := []string{"exec", containerName, "sh", "-c", script}

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

		for _, channel := range channels {
			containerName := compose.ResolveToolsContainerName(organization)
			block := fmt.Sprintf("%s/channels/%s.block", constants.DEFAULT_FABRIC_DIRECTORY, ResolveChannelID(channel))

			args := []string{
				"exec", containerName,
				"peer", "channel", "fetch", "newest", block,
				"-c", ResolveChannelID(channel),
				"-o", ordererAddress,
				"--tls", "--cafile", caFile,
			}

			deadline := time.Now().Add(60 * time.Second)

			for {
				_, err := f.executor.OutputCommand("docker", args...)

				if err == nil {
					break
				}

				if !strings.Contains(err.Error(), "SERVICE_UNAVAILABLE") {
					return fmt.Errorf("Fatal error while waiting for orderer readiness: %w", err)
				}

				if time.Now().After(deadline) {
					return fmt.Errorf("Timeout waiting for genesis block: %w", err)
				}

				time.Sleep(5 * time.Second)
			}
		}
	}

	return nil
}
