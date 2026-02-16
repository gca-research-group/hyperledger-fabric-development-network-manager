package fabric

import (
	"fmt"
	"os"
	"strings"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/docker"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/internal/file"
)

func (f *Fabric) IsDockerRunning() error {
	fmt.Print(">> Checking if Docker is running...\n")
	if !docker.IsDockerRunning() {
		return fmt.Errorf("Docker is not running")
	}
	return nil
}

func (f *Fabric) RemoveContainers() error {
	fmt.Print("\n=========== Removing peer containers ===========\n")

	for _, organization := range f.config.Organizations {
		fmt.Printf(">> Removing peer of the org: %s...\n", organization.Name)

		entries, err := os.ReadDir(fmt.Sprintf("%s/%s", f.config.Output, organization.Domain))

		if err != nil {
			return fmt.Errorf("Error when reading the directory to the organization %s: %v\n", organization.Name, err)
		}

		for _, entry := range entries {
			if !entry.IsDir() && strings.Contains(entry.Name(), "peer") {
				file := fmt.Sprintf("%s/%s/%s", f.config.Output, organization.Domain, entry.Name())
				args := []string{"compose", "-f", f.network, "-f", file, "down", "-v"}
				if err := f.executor.ExecCommand("docker", args...); err != nil {
					return fmt.Errorf("Error when removing the container for the organization %s, peer %s: %v\n", organization.Name, entry.Name(), err)
				}
			}
		}
	}

	fmt.Print("\n=========== Removing orderer containers ===========\n")
	for _, organization := range f.config.Organizations {
		for _, orderer := range organization.Orderers {
			config := fmt.Sprintf("%s/%s/%s.yml", f.config.Output, organization.Domain, orderer.Subdomain)

			if file.FileExists(config) {
				args := []string{"compose", "-f", f.network, "-f", config, "down", "-v"}
				if err := f.executor.ExecCommand("docker", args...); err != nil {
					return fmt.Errorf("Error when removing the container for the organization %s, orderer %s: %v\n", organization.Name, orderer.Name, err)
				}
			}
		}
	}

	fmt.Print("\n=========== Removing certificate authority containers ===========\n")
	for _, organization := range f.config.Organizations {
		file := fmt.Sprintf("%s/%s/ca.yml", f.config.Output, organization.Domain)
		args := []string{"compose", "-f", f.network, "-f", file, "down", "-v"}
		if err := f.executor.ExecCommand("docker", args...); err != nil {
			return fmt.Errorf("Error when removing the CA container for the organization %s: %v\n", organization.Name, err)
		}
	}

	return nil
}

func (f *Fabric) RunOrdererContainers() error {
	fmt.Print("\n=========== Executing orderer containeres ===========\n")
	for _, organization := range f.config.Organizations {
		for _, orderer := range organization.Orderers {
			config := fmt.Sprintf("%s/%s/%s.yml", f.config.Output, organization.Domain, orderer.Subdomain)

			if file.FileExists(config) {
				args := []string{"compose", "-f", f.network, "-f", config, "up", "--build", "-d"}
				if err := f.executor.ExecCommand("docker", args...); err != nil {
					return fmt.Errorf("Error when executing the orderer container for the organization %s, orderer %s: %v\n", organization.Name, orderer.Name, err)
				}
			}
		}
	}

	return nil
}

func (f *Fabric) RunPeerContainers() error {
	fmt.Print("\n=========== Executing peer containeres ===========\n")
	for _, organization := range f.config.Organizations {
		entries, err := os.ReadDir(fmt.Sprintf("%s/%s", f.config.Output, organization.Domain))

		if err != nil {
			return fmt.Errorf("Error when reading the directory to the org %s: %v\n", organization.Name, err)
		}

		for _, entry := range entries {
			if !entry.IsDir() && strings.Contains(entry.Name(), "peer") {
				file := fmt.Sprintf("%s/%s/%s", f.config.Output, organization.Domain, entry.Name())

				args := []string{"compose", "-f", f.network, "-f", file, "up", "--build", "-d"}

				if err := f.executor.ExecCommand("docker", args...); err != nil {
					return fmt.Errorf("Error when executing the container for the organization %s, peer %s: %v\n", organization.Name, entry.Name(), err)
				}
			}
		}
	}

	return nil
}

func (f *Fabric) RunCAContainers() error {
	fmt.Print("\n=========== Executing certificate authority containers ===========\n")
	for _, organization := range f.config.Organizations {
		file := fmt.Sprintf("%s/%s/ca.yml", f.config.Output, organization.Domain)

		args := []string{"compose", "-f", f.network, "-f", file, "up", "--build", "-d"}

		if err := f.executor.ExecCommand("docker", args...); err != nil {
			return fmt.Errorf("Error when executing the CA container for the organization %s: %v\n", organization.Name, err)
		}
	}

	return nil
}
