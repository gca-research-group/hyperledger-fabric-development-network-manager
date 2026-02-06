package fabric

import (
	"fmt"
	"strings"
)

func (f *Fabric) JoinOrdererToTheChannel() error {
	for _, organization := range f.config.Organizations {
		tools := fmt.Sprintf("%s/%s/tools.yml", f.config.Output, organization.Domain)
		for _, orderer := range organization.Orderers {
			for _, profile := range f.config.Profiles {
				containerName := buildToolsContainerName(organization)
				caFile := fmt.Sprintf("%s/%s/ordererOrganizations/%s/orderers/%s.%s/tls/ca.crt", DEFAULT_FABRIC_DIRECTORY, organization.Domain, organization.Domain, orderer.Hostname, organization.Domain)
				clientCert := fmt.Sprintf("%s/%s/ordererOrganizations/%s/orderers/%s.%s/tls/server.crt", DEFAULT_FABRIC_DIRECTORY, organization.Domain, organization.Domain, orderer.Hostname, organization.Domain)
				clientKey := fmt.Sprintf("%s/%s/ordererOrganizations/%s/orderers/%s.%s/tls/server.key", DEFAULT_FABRIC_DIRECTORY, organization.Domain, organization.Domain, orderer.Hostname, organization.Domain)

				args := []string{
					"compose", "-f", f.network, "-f", tools, "run", "--rm", "-T", containerName,
					"osnadmin", "channel", "join",
					"--channelID", strings.ToLower(profile.Name),
					"--config-block", fmt.Sprintf("%s/channel/%s.block", DEFAULT_FABRIC_DIRECTORY, strings.ToLower(profile.Name)),
					"-o", fmt.Sprintf("%s.%s:7053", orderer.Hostname, organization.Domain),
					"--ca-file", caFile,
					"--client-cert", clientCert,
					"--client-key", clientKey,
				}

				if err := f.executor.ExecCommand("docker", args...); err != nil {
					return fmt.Errorf("Error when joining the orderer %s of the organization %s to the channel %s: %v", orderer.Name, organization.Name, profile.Name, err)
				}
			}
		}
	}

	return nil
}

func (f *Fabric) JoinPeersToTheChannels() error {
	for _, organization := range f.config.Organizations {

		tools := fmt.Sprintf("%s/%s/tools.yml", f.config.Output, organization.Domain)

		for i := range organization.Peers {
			for _, profile := range f.config.Profiles {
				containerName := buildToolsContainerName(organization)
				block := fmt.Sprintf("%s/channel/%s.block", DEFAULT_FABRIC_DIRECTORY, strings.ToLower(profile.Name))
				tlsCertFile := fmt.Sprintf("%s/crypto-materials/%s/peerOrganizations/peers/peer%d.%s/tls/server.crt", DEFAULT_FABRIC_DIRECTORY, organization.Domain, i, organization.Domain)
				tlsKeyFile := fmt.Sprintf("%s/crypto-materials/%s/peerOrganizations/peers/peer%d.%s/tls/server.key", DEFAULT_FABRIC_DIRECTORY, organization.Domain, i, organization.Domain)
				mspConfigPath := fmt.Sprintf("%s/crypto-materials/peerOrganizations/%s/users/Admin@%s/msp", DEFAULT_FABRIC_DIRECTORY, organization.Domain, organization.Domain)

				args := []string{
					"compose", "-f", f.network, "-f", tools, "run", "--rm", "-T",
					"-e", fmt.Sprintf("CORE_PEER_ADDRESS=peer%d.%s:7051", i, organization.Domain),
					"-e", fmt.Sprintf("CORE_PEER_TLS_CERT_FILE=%s", tlsCertFile),
					"-e", fmt.Sprintf("CORE_PEER_TLS_KEY_FILE=%s", tlsKeyFile),
					"-e", fmt.Sprintf("CORE_PEER_MSPCONFIGPATH=%s", mspConfigPath),
					containerName, "peer", "channel", "join", "-b", block,
				}

				if err := f.executor.ExecCommand("docker", args...); err != nil {
					return fmt.Errorf("Error when joining the peer %d of the organization %s to the channel %s: %v", i, organization.Name, profile.Name, err)
				}
			}
		}
	}

	return nil
}
