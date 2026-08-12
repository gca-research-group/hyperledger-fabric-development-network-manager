package validate

import "github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"

func ValidateTopologyFn(configuration spec.Config) error {
	organizationNames := make(map[string]struct{})
	exposedPorts := make(map[int]portOwner)
	bootstrapCount := 0
	hasOrderer := false

	for _, organization := range configuration.Organizations {
		if err := DuplicateOrganizationNameFn(organization, organizationNames); err != nil {
			return err
		}

		organizationNames[organization.Name] = struct{}{}

		if organization.Bootstrap {
			bootstrapCount++
		}

		if len(organization.Orderers) > 0 {
			hasOrderer = true
		}

		if err := ExposedPortConflictFn(organization.CertificateAuthority.ExposePort, portOwner{ownerType: "Certificate Authority", name: organization.Name}, exposedPorts); err != nil {
			return err
		}

		for _, peer := range organization.Peers {
			if err := InvalidPeerVersionFn(peer, organization.Name, configuration.Capabilities.Channel); err != nil {
				return err
			}

			if err := ExposedPortConflictFn(peer.ExposePort, portOwner{ownerType: "peer", name: peer.Name}, exposedPorts); err != nil {
				return err
			}
		}

		for _, orderer := range organization.Orderers {
			if err := InvalidOrdererVersionFn(orderer, organization.Name, configuration.Capabilities.Channel); err != nil {
				return err
			}

			if err := ExposedPortConflictFn(orderer.ExposePort, portOwner{ownerType: "orderer", name: orderer.Name}, exposedPorts); err != nil {
				return err
			}
		}
	}

	if err := NoOrdererTopologyFn(hasOrderer); err != nil {
		return err
	}

	if err := MultipleBootstrapOrganizationsFn(bootstrapCount); err != nil {
		return err
	}

	for _, profile := range configuration.Profiles {
		if err := InvalidConsensusTypeFn(profile); err != nil {
			return err
		}

		if err := UndefinedProfileOrganizationFn(profile, organizationNames); err != nil {
			return err
		}
	}
	return nil
}
