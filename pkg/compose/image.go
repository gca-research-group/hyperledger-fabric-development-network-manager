package compose

import (
	"fmt"

	"github.com/gca-research-group/fabric-network-orchestrator/internal/executor"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/config"

	mapset "github.com/deckarep/golang-set/v2"
)

func PullImages(config config.Config, executor executor.Executor) error {

	images := mapset.NewSet[string]()
	images.Add(ResolveCouchDBImage())

	for _, organization := range config.Organizations {
		images.Add(ResolveCertificateAuthorityImage(organization.CertificateAuthority))

		for _, peer := range organization.Peers {
			images.Add(ResolvePeerImage(peer))
			images.Add(ResolveChaincodeCompilerImage(peer))
		}

		for _, orderer := range organization.Orderers {
			images.Add(ResolveOrdererImage(orderer))
		}

		images.Add(ResolveToolsImage(config.Capabilities))
	}

	for image := range images.Iter() {
		if _, err := executor.OutputCommand("docker", "pull", image); err != nil {
			return fmt.Errorf("Error when pulling the image %s: %v\n", image, err)
		}
	}

	return nil
}
