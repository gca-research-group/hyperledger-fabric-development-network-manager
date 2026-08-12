package validate

import "github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"

func Config(configuration spec.Config) error {

	if err := InvalidOutputDirectoryNameFn(configuration); err != nil {
		return err
	}

	if err := NoOrganizationFn(configuration); err != nil {
		return err
	}

	for i, organization := range configuration.Organizations {
		if err := EmptyOrganizationNameFn(organization, i); err != nil {
			return err
		}

		if err := EmptyOrganizationDomainFn(organization, i); err != nil {
			return err
		}

		if err := InvalidCertificateAuthorityPortFn(organization); err != nil {
			return err
		}

		for j, peer := range organization.Peers {
			if err := EmptyPeerNameFn(peer, j, organization.Name); err != nil {
				return err
			}

			if err := EmptyPeerSubdomainFn(peer, organization.Name); err != nil {
				return err
			}

			if err := InvalidPeerPortFn(peer, organization.Name); err != nil {
				return err
			}
		}

		for j, orderer := range organization.Orderers {
			if err := EmptyOrdererNameFn(orderer, j, organization.Name); err != nil {
				return err
			}

			if err := EmptyOrdererSubdomainFn(orderer, organization.Name); err != nil {
				return err
			}

			if err := InvalidOrdererPortFn(orderer, organization.Name); err != nil {
				return err
			}
		}
	}

	for _, profile := range configuration.Profiles {
		if err := EmptyProfileOrganizationsFn(profile); err != nil {
			return err
		}
	}

	for _, channel := range configuration.Channels {
		if err := EmptyChannelNameFn(channel); err != nil {
			return err
		}

		if err := EmptyChannelProfileFn(channel); err != nil {
			return err
		}

		if err := InvalidChannelNameFn(channel); err != nil {
			return err
		}

		for i, chaincode := range channel.Chaincodes {
			if err := EmptyChaincodeNameFn(chaincode, i); err != nil {
				return err
			}

			if err := EmptyChaincodePathFn(chaincode, i); err != nil {
				return err
			}

			if err := EmptyChaincodeVersionFn(chaincode, i); err != nil {
				return err
			}
		}
	}

	if err := UnsupportedChannelCapabilityFn(configuration); err != nil {
		return err
	}

	if err := UnsupportedApplicationCapabilityFn(configuration); err != nil {
		return err
	}

	if err := UnsupportedOrdererCapabilityFn(configuration); err != nil {
		return err
	}

	if err := ValidateTopologyFn(configuration); err != nil {
		return err
	}

	return nil
}
