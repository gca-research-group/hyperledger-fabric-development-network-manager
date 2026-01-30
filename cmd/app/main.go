package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

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

	fmt.Print(">> Cleaning output folder...\n")
	if err := RemoveFolderIfExists(config.Output); err != nil {
		log.Fatalf("Error when cleaning output folder: %v", err)
	}

	if err := cryptoconfig.NewRenderer(config).Render(); err != nil {
		panic(err)
	}

	if err := configtx.NewRenderer(config).Render(); err != nil {
		panic(err)
	}

	renderer := docker.NewRenderer(config)

	if err := renderer.Render(); err != nil {
		panic(err)
	}

	fmt.Print(">> Checking if docker is running...\n")

	if !isDockerRunning() {
		panic("Docker is not running...")
	}

	network := fmt.Sprintf("%s/network.yml", config.Output)

	// **************** Removing containers

	fmt.Print("\n=========== Removing peer containers in execution ===========\n")

	for _, organization := range config.Organizations {
		fmt.Printf(">> Removing peer of the org: %s...\n", organization.Name)

		entries, err := os.ReadDir(fmt.Sprintf("%s/%s", config.Output, organization.Domain))

		if err != nil {
			log.Fatalf("Error when reading the directory to the org %s: %v\n", organization.Name, err)
		}

		for _, entry := range entries {
			if !entry.IsDir() && strings.Contains(entry.Name(), "peer") {
				var args []string

				file := fmt.Sprintf("%s/%s/%s", config.Output, organization.Domain, entry.Name())
				args = append(args, "compose", "-f", network, "-f", file, "down", "-v")
				if err := ExecCommand("docker", args...); err != nil {
					log.Fatalf("Error when removing the container for the organization %s, peer %s: %v\n", organization.Name, entry.Name(), err)
				}
			}
		}
	}

	fmt.Print("\n=========== Removing orderer containers ===========\n")
	for _, organization := range config.Organizations {
		for _, orderer := range organization.Orderers {
			file := fmt.Sprintf("%s/%s/%s.yml", config.Output, organization.Domain, orderer.Hostname)

			if FileExists(file) {
				var args []string

				args = append(args, "compose", "-f", network, "-f", file, "down", "-v")
				if err := ExecCommand("docker", args...); err != nil {
					log.Fatalf("Error when removing the container for the organization %s, orderer %s: %v\n", organization.Name, orderer.Name, err)
				}
			}
		}
	}

	for _, organization := range config.Organizations {
		tools := fmt.Sprintf("%s/%s/tools.yml", config.Output, organization.Domain)

		fmt.Printf("\n=========== Generating crypto materials to %s ===========\n", organization.Name)

		containerName := fmt.Sprintf("hyperledger-fabric-%s-tools", organization.Domain)

		var args []string

		args = append(args, "compose", "-f", network, "-f", tools, "run", "--rm", "-T", containerName)
		args = append(args, "cryptogen", "generate")
		args = append(args, "--config=/opt/gopath/src/github.com/hyperledger/fabric/crypto-config.yml")
		args = append(args, "--output=/opt/gopath/src/github.com/hyperledger/fabric/crypto-materials")

		if err := ExecCommand("docker", args...); err != nil {
			log.Fatalf("Error when generating the crypto materials for the organization %s: %v\n", organization.Name, err)
		}
	}

	for _, organization := range config.Organizations {
		if organization.Bootstrap {
			fmt.Printf("\n=========== Generating orderer genesis block to %s ===========\n", organization.Name)
			renderer.RenderToolsWithMSP(organization)

			containerName := fmt.Sprintf("hyperledger-fabric-%s-tools", organization.Domain)
			tools := fmt.Sprintf("%s/%s/tools.yml", config.Output, organization.Domain)

			var args []string

			args = append(args, "compose", "-f", network, "-f", tools, "run", "--rm", "-T", containerName)
			args = append(args, "configtxgen")
			args = append(args, "-outputBlock", "/opt/gopath/src/github.com/hyperledger/fabric/channel/orderer.genesis.block")
			args = append(args, "-profile", configtx.OrdererGenesisProfileKey)
			args = append(args, "-channelID", "defaultchannel")
			args = append(args, "-configPath", "/opt/gopath/src/github.com/hyperledger/fabric/")

			if err := ExecCommand("docker", args...); err != nil {
				log.Fatalf("Error when generating the genesis block for the organization %s: %v", organization.Name, err)
			}

			for _, profile := range config.Profiles {
				fmt.Printf("\n=========== Generating the channel.tx files to the profile %s ===========\n", profile.Name)
				var args []string

				channelId := strings.ToLower(profile.Name)

				args = append(args, "compose", "-f", network, "-f", tools, "run", "--rm", "-T", containerName)
				args = append(args, "configtxgen")
				args = append(args, "-outputCreateChannelTx", fmt.Sprintf("/opt/gopath/src/github.com/hyperledger/fabric/channel/%s.tx", channelId))
				args = append(args, "-profile", profile.Name)
				args = append(args, "-channelID", channelId)
				args = append(args, "-configPath", "/opt/gopath/src/github.com/hyperledger/fabric/")

				if err := ExecCommand("docker", args...); err != nil {
					log.Fatalf("Error when generating the channel transaction file for the organization %s, profile %s: %v\n", organization.Name, profile.Name, err)
				}
			}
		}
	}

	fmt.Print("\n=========== Executing orderer containeres ===========\n")
	for _, organization := range config.Organizations {
		for _, orderer := range organization.Orderers {
			file := fmt.Sprintf("%s/%s/%s.yml", config.Output, organization.Domain, orderer.Hostname)

			if FileExists(file) {
				var args []string

				args = append(args, "compose", "-f", network, "-f", file, "up", "--build", "-d")
				if err := ExecCommand("docker", args...); err != nil {
				}
			}
		}
	}

	fmt.Print("\n=========== Executing peer containeres ===========\n")
	for _, organization := range config.Organizations {
		entries, err := os.ReadDir(fmt.Sprintf("%s/%s", config.Output, organization.Domain))

		if err != nil {
			log.Fatalf("Error when reading the directory to the org %s: %v\n", organization.Name, err)
		}

		for _, entry := range entries {
			if !entry.IsDir() && strings.Contains(entry.Name(), "peer") {
				var args []string

				file := fmt.Sprintf("%s/%s/%s", config.Output, organization.Domain, entry.Name())
				args = append(args, "compose", "-f", network, "-f", file, "up", "--build", "-d")
				if err := ExecCommand("docker", args...); err != nil {
					log.Fatalf("Error when executing the container for the organization %s, peer %s: %v\n", organization.Name, entry.Name(), err)
				}
			}
		}
	}

	time.Sleep(5 * time.Second)

	for _, organization := range config.Organizations {
		containerName := fmt.Sprintf("hyperledger-fabric-%s-tools", organization.Domain)
		tools := fmt.Sprintf("%s/%s/tools.yml", config.Output, organization.Domain)

		if organization.Bootstrap {
			for _, profile := range config.Profiles {
				var args []string

				channelId := strings.ToLower(profile.Name)
				orderer := organization.Orderers[0]
				ordererAddress := fmt.Sprintf("%s.%s:%d", orderer.Hostname, organization.Domain, orderer.Port)
				cafile := fmt.Sprintf("/opt/gopath/src/github.com/hyperledger/fabric/crypto-materials/ordererOrganizations/%s/orderers/%s.%s/msp/tlscacerts/tlsca.%s-cert.pem", organization.Domain, orderer.Hostname, organization.Domain, organization.Domain)

				args = append(args, "compose", "-f", network, "-f", tools, "run", "--rm", "-T", "-w", "/opt/gopath/src/github.com/hyperledger/fabric/channel", containerName)

				args = append(args, "peer", "channel", "create", "-o", ordererAddress, "-c", channelId, "-f", fmt.Sprintf("%s.tx", channelId), "--tls", "true", "--cafile", cafile)

				if err := ExecCommand("docker", args...); err != nil {
					log.Fatalf("Error when generating the channel file for the organization %s, profile %s: %v\n", organization.Name, profile.Name, err)
				}
			}
		}
	}

	// fmt.Print("\n=========== Joining peers to the channel ===========\n")

	// for _, organization := range config.Organizations {
	// 	tools := fmt.Sprintf("%s/%s/tools.yml", config.Output, organization.Domain)

	// 	containerName := fmt.Sprintf("hyperledger-fabric-%s-tools", organization.Domain)
	// 	msconfigPath := fmt.Sprintf("/opt/gopath/src/github.com/hyperledger/fabric/crypto-materials/peerOrganizations/%s/users/Admin@%s/msp", organization.Domain, organization.Domain)
	// 	localMSPID := fmt.Sprintf("%sMSP", organization.Name)

	// 	peers := 1

	// 	if organization.Peers > 0 {
	// 		peers = organization.Peers
	// 	}

	// 	for i := 0; i < peers; i++ {
	// 		address := fmt.Sprintf("peer%d.%s:7051", i, organization.Domain)
	// 		tlsRootCertFile := fmt.Sprintf("/opt/gopath/src/github.com/hyperledger/fabric/crypto-materials/peerOrganizations/%s/peers/peer%d.%s/tls/ca.crt", organization.Domain, i, organization.Domain)

	// 		var args []string
	// 		args = append(args, "compose", "-f", network, "-f", tools, "run", "--rm", "-T")
	// 		args = append(args, "-e", fmt.Sprintf("CORE_PEER_MSPCONFIGPATH=%s", msconfigPath))
	// 		args = append(args, "-e", fmt.Sprintf("CORE_PEER_ADDRESS=%s", address))
	// 		args = append(args, "-e", fmt.Sprintf("CORE_PEER_LOCALMSPID=%s", localMSPID))
	// 		args = append(args, "-e", fmt.Sprintf("CORE_PEER_TLS_ROOTCERT_FILE=%s", tlsRootCertFile))
	// 		args = append(args, containerName, "peer", "channel", "join", "-b", "/opt/gopath/src/github.com/hyperledger/fabric/channel/orderer.genesis.block")

	// 		if err := ExecCommand("docker", args...); err != nil {
	// 			log.Fatalf("Error when joining peer %d to the channel for the organization %s: %v\n", i, organization.Name, err)
	// 		}
	// 	}
}
