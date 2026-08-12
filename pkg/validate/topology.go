package validate

import (
	"errors"

	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func ValidateTopologyFn(configuration spec.Config) error {
	var validationErrors []error
	collect := func(err error) {
		if err != nil {
			validationErrors = append(validationErrors, err)
		}
	}

	organizationNames := make(map[string]struct{})
	exposedPorts := make(map[int]portOwner)
	bootstrapCount := 0
	hasOrderer := false

	for _, organization := range configuration.Organizations {
		collect(DuplicateOrganizationNameFn(organization, organizationNames))

		organizationNames[organization.Name] = struct{}{}

		if organization.Bootstrap {
			bootstrapCount++
		}

		if len(organization.Orderers) > 0 {
			hasOrderer = true
		}

		collect(ExposedPortConflictFn(organization.CertificateAuthority.ExposePort, portOwner{ownerType: "Certificate Authority", name: organization.Name}, exposedPorts))

		for _, peer := range organization.Peers {
			collect(InvalidPeerVersionFn(peer, organization.Name, configuration.Capabilities.Channel))
			collect(ExposedPortConflictFn(peer.ExposePort, portOwner{ownerType: "peer", name: peer.Name}, exposedPorts))
		}

		for _, orderer := range organization.Orderers {
			collect(InvalidOrdererVersionFn(orderer, organization.Name, configuration.Capabilities.Channel))
			collect(ExposedPortConflictFn(orderer.ExposePort, portOwner{ownerType: "orderer", name: orderer.Name}, exposedPorts))
		}
	}

	collect(NoOrdererTopologyFn(hasOrderer))
	collect(MultipleBootstrapOrganizationsFn(bootstrapCount))

	for _, profile := range configuration.Profiles {
		collect(InvalidConsensusTypeFn(profile))
		collect(UndefinedProfileOrganizationFn(profile, organizationNames))
	}

	return errors.Join(validationErrors...)
}
