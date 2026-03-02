package network

import (
	"fmt"
	"strings"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/constants"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/compose"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/config"
)

func (f *Network) JoinOrdererToTheChannel() error {

	for _, organization := range f.config.Organizations {
		tools := compose.ResolveToolsDockerComposeFile(f.config.Output, organization.Domain)
		containerName := compose.ResolveToolsContainerName(organization)

		for _, orderer := range organization.Orderers {
			for _, channel := range f.config.Channels {

				caFile := fmt.Sprintf("%[1]s/%[2]s/ordererOrganizations/%[2]s/users/Admin@%[2]s/tls/ca.crt", constants.DEFAULT_FABRIC_DIRECTORY, organization.Domain)

				clientCert := fmt.Sprintf("%[1]s/%[2]s/ordererOrganizations/%[2]s/users/Admin@%[2]s/tls/client.crt", constants.DEFAULT_FABRIC_DIRECTORY, organization.Domain)

				clientKey := fmt.Sprintf("%[1]s/%[2]s/ordererOrganizations/%[2]s/users/Admin@%[2]s/tls/client.key", constants.DEFAULT_FABRIC_DIRECTORY, organization.Domain)

				blockFile := fmt.Sprintf("%s/channels/%s.block", constants.DEFAULT_FABRIC_DIRECTORY, ResolveChannelID(channel))

				args := []string{
					"compose", "-f", f.network, "-f", tools,
					"run", "--rm", "-T", containerName,
					"curl",
					"-X", "POST",
					fmt.Sprintf("https://%s.%s:7053/participation/v1/channels", orderer.Subdomain, organization.Domain),
					"--cacert", caFile,
					"--cert", clientCert,
					"--key", clientKey,
					"-F", fmt.Sprintf("config-block=@%s", blockFile),
				}

				if err := f.executor.ExecCommand("docker", args...); err != nil {
					return fmt.Errorf("error joining orderer %s of organization %s to channel %s: %v", orderer.Name, organization.Name, channel.Name, err)
				}
			}
		}
	}

	return nil
}

func (f *Network) JoinPeersToTheChannels() error {
	for _, organization := range f.config.Organizations {

		tools := compose.ResolveToolsDockerComposeFile(f.config.Output, organization.Domain)

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

			for _, peer := range organization.Peers {
				peerPort := compose.ResolvePeerPort(peer.Port)

				tlsCertFile := fmt.Sprintf("%[1]s/%[2]s/peerOrganizations/peers/%[3]s.%[2]s/tls/server.crt", constants.DEFAULT_FABRIC_DIRECTORY, organization.Domain, peer.Subdomain)
				tlsKeyFile := fmt.Sprintf("%[1]s/%[2]s/peerOrganizations/peers/%[3]s.%[2]s/tls/server.key", constants.DEFAULT_FABRIC_DIRECTORY, organization.Domain, peer.Subdomain)
				mspConfigPath := fmt.Sprintf("%[1]s/%[2]s/peerOrganizations/%[2]s/users/Admin@%[2]s/msp", constants.DEFAULT_FABRIC_DIRECTORY, organization.Domain)

				args := []string{
					"compose", "-f", f.network, "-f", tools, "run", "--rm", "-T",
					"-e", fmt.Sprintf("CORE_PEER_ADDRESS=%s.%s:%d", peer.Subdomain, organization.Domain, peerPort),
					"-e", fmt.Sprintf("CORE_PEER_TLS_CERT_FILE=%s", tlsCertFile),
					"-e", fmt.Sprintf("CORE_PEER_TLS_KEY_FILE=%s", tlsKeyFile),
					"-e", fmt.Sprintf("CORE_PEER_MSPCONFIGPATH=%s", mspConfigPath),
					containerName,
				}

				output, err := f.executor.OutputCommand("docker", append(args, "peer", "channel", "list")...)

				if err != nil {
					return fmt.Errorf("Error when listing the channels that the peer %s of the organization %s has joined to the channel %s: %v\n", peer.Subdomain, organization.Name, channel.Name, err)
				}

				if strings.Contains(string(output), ResolveChannelID(channel)) {
					fmt.Printf("Skipping: the peer %s of the organization %s has already joined to the channel %s\n", peer.Subdomain, organization.Name, channel.Name)
					continue
				}

				if err := f.executor.ExecCommand("docker", append(args, "peer", "channel", "join", "-b", block)...); err != nil {
					return fmt.Errorf("Error when joining the peer %s of the organization %s to the channel %s: %v\n", peer.Subdomain, organization.Name, channel.Name, err)
				}
			}
		}
	}

	return nil
}
