package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/docker/docker/client"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/configtx"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/cryptoconfig"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/docker"
)

func FolderExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}

func FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !info.IsDir()
}

func RemoveFolderIfExists(folderPath string) error {
	if FolderExists(folderPath) {
		return os.RemoveAll(folderPath)
	}

	return nil
}

func isDockerRunning() bool {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return false
	}
	defer cli.Close()

	_, err = cli.Ping(context.Background())
	return err == nil
}

func ExecCommand(name string, arg ...string) error {
	fmt.Printf("%s %s\n", name, strings.Join(arg, " "))
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func buildToolsContainerName(organization pkg.Organization) string {
	return fmt.Sprintf("hyperledger-fabric-tools-%s", strings.ToLower(organization.Name))
}

type Fabric struct {
	config               pkg.Config
	network              string
	crytpoConfigRenderer *cryptoconfig.Renderer
	configTxRenderer     *configtx.Renderer
	dockerRenderer       *docker.Renderer
}

func (f *Fabric) removeContainers() *Fabric {
	fmt.Print("\n=========== Removing peer containers in execution ===========\n")

	for _, organization := range f.config.Organizations {
		fmt.Printf(">> Removing peer of the org: %s...\n", organization.Name)

		entries, err := os.ReadDir(fmt.Sprintf("%s/%s", f.config.Output, organization.Domain))

		if err != nil {
			log.Fatalf("Error when reading the directory to the org %s: %v\n", organization.Name, err)
		}

		for _, entry := range entries {
			if !entry.IsDir() && strings.Contains(entry.Name(), "peer") {
				var args []string

				file := fmt.Sprintf("%s/%s/%s", f.config.Output, organization.Domain, entry.Name())
				args = append(args, "compose", "-f", f.network, "-f", file, "down", "-v")
				if err := ExecCommand("docker", args...); err != nil {
					log.Fatalf("Error when removing the container for the organization %s, peer %s: %v\n", organization.Name, entry.Name(), err)
				}
			}
		}
	}

	fmt.Print("\n=========== Removing orderer containers ===========\n")
	for _, organization := range f.config.Organizations {
		for _, orderer := range organization.Orderers {
			file := fmt.Sprintf("%s/%s/%s.yml", f.config.Output, organization.Domain, orderer.Hostname)

			if FileExists(file) {
				var args []string

				args = append(args, "compose", "-f", f.network, "-f", file, "down", "-v")
				if err := ExecCommand("docker", args...); err != nil {
					log.Fatalf("Error when removing the container for the organization %s, orderer %s: %v\n", organization.Name, orderer.Name, err)
				}
			}
		}
	}

	return f
}

func (f *Fabric) renderConfigFiles() *Fabric {
	fmt.Print(">> Cleaning output folder...\n")
	if err := RemoveFolderIfExists(f.config.Output); err != nil {
		log.Fatalf("Error when cleaning output folder: %v", err)
	}

	if err := f.crytpoConfigRenderer.Render(); err != nil {
		panic(err)
	}

	if err := f.configTxRenderer.Render(); err != nil {
		panic(err)
	}

	if err := f.dockerRenderer.Render(); err != nil {
		panic(err)
	}

	return f
}

func (f *Fabric) isDockerRunning() *Fabric {
	fmt.Print(">> Checking if docker is running...\n")

	if !isDockerRunning() {
		panic("Docker is not running...")
	}

	return f
}

func (f *Fabric) generateCryptoMaterial() *Fabric {
	for _, organization := range f.config.Organizations {
		tools := fmt.Sprintf("%s/%s/tools.yml", f.config.Output, organization.Domain)

		fmt.Printf("\n=========== Generating crypto materials to %s ===========\n", organization.Name)

		containerName := buildToolsContainerName(organization)

		var args []string

		args = append(args, "compose", "-f", f.network, "-f", tools, "run", "--rm", "-T", containerName)
		args = append(args, "cryptogen", "generate")
		args = append(args, "--config=/opt/gopath/src/github.com/hyperledger/fabric/crypto-config.yml")
		args = append(args, "--output=/opt/gopath/src/github.com/hyperledger/fabric/crypto-materials")

		if err := ExecCommand("docker", args...); err != nil {
			log.Fatalf("Error when generating the crypto materials for the organization %s: %v\n", organization.Name, err)
		}
	}

	return f
}

func (f *Fabric) generateGenesisBlock() *Fabric {
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
				args = append(args, "-outputBlock", fmt.Sprintf("/opt/gopath/src/github.com/hyperledger/fabric/channel/%s.block", strings.ToLower(profile.Name)))
				args = append(args, "-profile", profile.Name)
				args = append(args, "-channelID", strings.ToLower(profile.Name))
				args = append(args, "-configPath", "/opt/gopath/src/github.com/hyperledger/fabric/")

				if err := ExecCommand("docker", args...); err != nil {
					log.Fatalf("Error when generating the genesis block for the organization %s: %v", organization.Name, err)
				}
			}
		}
	}

	return f
}

func (f *Fabric) joinOrdererToTheChannel() *Fabric {
	for _, organization := range f.config.Organizations {
		tools := fmt.Sprintf("%s/%s/tools.yml", f.config.Output, organization.Domain)
		for _, orderer := range organization.Orderers {
			for _, profile := range f.config.Profiles {
				var args []string
				containerName := buildToolsContainerName(organization)
				caFile := fmt.Sprintf("/opt/gopath/src/github.com/hyperledger/fabric/crypto-materials/ordererOrganizations/%s/orderers/%s.%s/tls/ca.crt", organization.Domain, orderer.Hostname, organization.Domain)
				clientCert := fmt.Sprintf("/opt/gopath/src/github.com/hyperledger/fabric/crypto-materials/ordererOrganizations/%s/orderers/%s.%s/tls/server.crt", organization.Domain, orderer.Hostname, organization.Domain)
				clientKey := fmt.Sprintf("/opt/gopath/src/github.com/hyperledger/fabric/crypto-materials/ordererOrganizations/%s/orderers/%s.%s/tls/server.key", organization.Domain, orderer.Hostname, organization.Domain)

				args = append(args, "compose", "-f", f.network, "-f", tools, "run", "--rm", "-T", containerName)
				args = append(args, "osnadmin", "channel", "join")
				args = append(args, "--channelID", strings.ToLower(profile.Name))
				args = append(args, "--config-block", fmt.Sprintf("/opt/gopath/src/github.com/hyperledger/fabric/channel/%s.block", strings.ToLower(profile.Name)))
				args = append(args, "-o", fmt.Sprintf("%s.%s:7053", orderer.Hostname, organization.Domain))
				args = append(args, "--ca-file", caFile)
				args = append(args, "--client-cert", clientCert)
				args = append(args, "--client-key", clientKey)

				if err := ExecCommand("docker", args...); err != nil {
					log.Fatalf("Error when joining the orderer %s of the organization %s to the channel %s: %v", orderer.Name, organization.Name, profile.Name, err)
				}
			}
		}
	}

	return f
}

func (f *Fabric) fetchTheGenesisBlock() *Fabric {

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
	caFile := fmt.Sprintf("/opt/gopath/src/github.com/hyperledger/fabric/crypto-materials/ordererOrganizations/%s/orderers/%s.%s/tls/ca.crt", ordererDomain, orderer.Hostname, ordererDomain)

	for _, organization := range f.config.Organizations {
		if organization.Bootstrap {
			continue
		}

		tools := fmt.Sprintf("%s/%s/tools.yml", f.config.Output, organization.Domain)
		for _, profile := range f.config.Profiles {
			containerName := buildToolsContainerName(organization)
			block := fmt.Sprintf("/opt/gopath/src/github.com/hyperledger/fabric/channel/%s.block", strings.ToLower(profile.Name))

			var args []string
			args = append(args, "compose", "-f", f.network, "-f", tools, "run", "--rm", "-T", containerName)
			args = append(args, "peer", "channel", "fetch", "0", block, "-c", strings.ToLower(profile.Name), "-o", ordererAddress, "--tls", "--cafile", caFile)

			if err := ExecCommand("docker", args...); err != nil {
				log.Fatalf("Error when fetching the orderer %s of the organization %s to the channel %s: %v", orderer.Name, organization.Name, profile.Name, err)
			}
		}
	}

	return f
}

func (f *Fabric) joinPeersToTheChannels() *Fabric {
	for _, organization := range f.config.Organizations {

		tools := fmt.Sprintf("%s/%s/tools.yml", f.config.Output, organization.Domain)
		for i := 0; i < organization.Peers; i++ {
			for _, profile := range f.config.Profiles {
				containerName := buildToolsContainerName(organization)
				block := fmt.Sprintf("/opt/gopath/src/github.com/hyperledger/fabric/channel/%s.block", strings.ToLower(profile.Name))
				tlsCertFile := fmt.Sprintf("/opt/gopath/src/github.com/hyperledger/fabric/crypto-materials/%s/peerOrganizations/peers/peer%d.%s/tls/server.crt", organization.Domain, i, organization.Domain)
				tlsKeyFile := fmt.Sprintf("/opt/gopath/src/github.com/hyperledger/fabric/crypto-materials/%s/peerOrganizations/peers/peer%d.%s/tls/server.key", organization.Domain, i, organization.Domain)
				mspConfigPath := fmt.Sprintf("/opt/gopath/src/github.com/hyperledger/fabric/crypto-materials/peerOrganizations/%s/users/Admin@%s/msp", organization.Domain, organization.Domain)

				var args []string
				args = append(args, "compose", "-f", f.network, "-f", tools, "run", "--rm", "-T")
				args = append(args, "-e", fmt.Sprintf("CORE_PEER_ADDRESS=peer%d.%s:7051", i, organization.Domain))
				args = append(args, "-e", fmt.Sprintf("CORE_PEER_TLS_CERT_FILE=%s", tlsCertFile))
				args = append(args, "-e", fmt.Sprintf("CORE_PEER_TLS_KEY_FILE=%s", tlsKeyFile))
				args = append(args, "-e", fmt.Sprintf("CORE_PEER_MSPCONFIGPATH=%s", mspConfigPath))
				args = append(args, containerName, "peer", "channel", "join", "-b", block)

				if err := ExecCommand("docker", args...); err != nil {
					log.Fatalf("Error when joining the peer %d of the organization %s to the channel %s: %v", i, organization.Name, profile.Name, err)
				}
			}
		}
	}

	return f
}

func (f *Fabric) runOrdererContainers() *Fabric {
	fmt.Print("\n=========== Executing orderer containeres ===========\n")
	for _, organization := range f.config.Organizations {
		for _, orderer := range organization.Orderers {
			file := fmt.Sprintf("%s/%s/%s.yml", f.config.Output, organization.Domain, orderer.Hostname)

			if FileExists(file) {
				var args []string

				args = append(args, "compose", "-f", f.network, "-f", file, "up", "--build", "-d")
				if err := ExecCommand("docker", args...); err != nil {
				}
			}
		}
	}

	return f
}

func (f *Fabric) runPeerContainers() *Fabric {
	fmt.Print("\n=========== Executing peer containeres ===========\n")
	for _, organization := range f.config.Organizations {
		entries, err := os.ReadDir(fmt.Sprintf("%s/%s", f.config.Output, organization.Domain))

		if err != nil {
			log.Fatalf("Error when reading the directory to the org %s: %v\n", organization.Name, err)
		}

		for _, entry := range entries {
			if !entry.IsDir() && strings.Contains(entry.Name(), "peer") {
				var args []string

				file := fmt.Sprintf("%s/%s/%s", f.config.Output, organization.Domain, entry.Name())
				args = append(args, "compose", "-f", f.network, "-f", file, "up", "--build", "-d")
				if err := ExecCommand("docker", args...); err != nil {
					log.Fatalf("Error when executing the container for the organization %s, peer %s: %v\n", organization.Name, entry.Name(), err)
				}
			}
		}
	}

	return f
}

func newFabric(config pkg.Config) *Fabric {
	network := fmt.Sprintf("%s/network.yml", config.Output)

	crytpoConfigRenderer := cryptoconfig.NewRenderer(config)
	configTxRenderer := configtx.NewRenderer(config)
	dockerRenderer := docker.NewRenderer(config)

	return &Fabric{config, network, crytpoConfigRenderer, configTxRenderer, dockerRenderer}
}

func main() {
	config := pkg.Config{
		Output: "./dist/hyperledger-fabric",
		Organizations: []pkg.Organization{
			{
				Name:   "Org1",
				Domain: "org1.com",
				Peers:  3,
				Orderers: []pkg.Orderer{
					{
						Name:     "Orderer",
						Hostname: "orderer",
						Port:     7050,
					},
				},
				Bootstrap: true,
			},
			{
				Name:   "Org2",
				Domain: "org2.com",
			},
			{
				Name:   "Org3",
				Domain: "org3.com",
			},
		},
		Profiles: []pkg.Profile{
			{Name: "MultiChannel", Organizations: []string{"Org1", "Org2", "Org3"}},
		},
	}

	fabric := newFabric(config)

	fabric.
		isDockerRunning().
		removeContainers().
		renderConfigFiles().
		generateCryptoMaterial().
		generateGenesisBlock().
		runOrdererContainers().
		runPeerContainers().
		joinOrdererToTheChannel().
		fetchTheGenesisBlock().
		joinPeersToTheChannels()
}
