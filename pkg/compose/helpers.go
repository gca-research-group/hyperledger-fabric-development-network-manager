package compose

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/docker/docker/client"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/constants"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/executor"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/config"
)

func ResolvePeerPort(port int) int {
	if port == 0 {
		return constants.DEFAULT_PEER_PORT
	}

	return port
}

func ResolveOrdererPort(port int) int {
	if port == 0 {
		return constants.DEFAULT_ORDERER_PORT
	}

	return port
}

func ResolvePeerVersion(version string) string {
	if version == "" {
		return constants.DEFAULT_FABRIC_VERSION
	}

	return version
}

func ResolvePeerDomain(subdomain string, domain string) string {
	return fmt.Sprintf("%s.%s", subdomain, domain)
}

func ResolveOrdererDomain(subdomain string, domain string) string {
	return fmt.Sprintf("%s.%s", subdomain, domain)
}

func ResolveOrdererVersion(version string) string {
	if version == "" {
		return constants.DEFAULT_FABRIC_VERSION
	}

	return version
}

func ResolveCertificateAuthorityDomain(domain string) string {
	return fmt.Sprintf("ca.%s", domain)
}

func ResolveCertificateAuthorityVersion(version string) string {
	if version == "" {
		return "latest"
	}

	return version
}

func ResolveDockerNetworkName(network string) string {
	if network == "" {
		return constants.DEFAULT_NETWORK
	}

	return network
}

func IsDockerRunning() bool {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return false
	}

	defer cli.Close()

	_, err = cli.Ping(context.Background())
	return err == nil
}

func RemoveContainersInNetwork(network string) error {
	list := exec.Command(
		"docker",
		"network",
		"inspect",
		network,
		"--format",
		"{{ range .Containers }}{{ .Name }}{{ \"\\n\" }}{{ end }}",
	)

	out, err := list.Output()
	if err != nil {
		return err
	}

	containers := strings.Fields(string(out))

	if len(containers) == 0 {
		return nil
	}

	for _, c := range containers {
		rm := exec.Command("docker", "rm", "-f", c)
		rm.Stderr = os.Stderr

		if err := rm.Run(); err != nil {
			return err
		}
	}

	return nil
}

func ResolvePeerContainerName(domain string, subdomain string) string {
	return fmt.Sprintf("%s.%s", subdomain, domain)
}

func ResolvePeerDockerComposeFile(output string, domain string, subdomain string) string {
	return fmt.Sprintf("%[1]s/%[2]s/peers/%[3]s/%[3]s.yml", output, domain, subdomain)
}

func ResolvePeerCouchDBDockerComposeFile(output string, domain string, subdomain string) string {
	return fmt.Sprintf("%[1]s/%[2]s/peers/%[3]s/couchdb.yml", output, domain, subdomain)
}

func ResolveCertificateAuthorityDockerComposeFile(output string, domain string) string {
	return fmt.Sprintf("%s/%s/certificate-authority/certificate-authority.yml", output, domain)
}

func ResolveOrdererDockerComposeFile(output string, domain string, subdomain string) string {
	return fmt.Sprintf("%[1]s/%[2]s/orderers/%[3]s/%[3]s.yml", output, domain, subdomain)
}

func ResolveToolsDockerComposeFile(output string, domain string) string {
	return fmt.Sprintf("%s/%s/tools.yml", output, domain)
}

func ResolveNetworkDockerComposeFile(output string) string {
	return fmt.Sprintf("%s/network.yml", output)
}

func ResolveCertificateAuthorityContainerName(domain string) string {
	return fmt.Sprintf("ca.%s", domain)
}

func ResolveToolsContainerName(organization config.Organization) string {
	return fmt.Sprintf("hyperledger-fabric-tools-%s", strings.ToLower(organization.Name))
}

func RunContainerFromTheDockerComposeFile(network string, file string) error {
	args := []string{"compose", "-f", network, "-f", file, "up", "--build", "-d"}
	executor := &executor.DefaultExecutor{}

	return executor.ExecCommand("docker", args...)
}
