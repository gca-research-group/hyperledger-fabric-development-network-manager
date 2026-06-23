package config

import (
	"fmt"
	"regexp"
	"strings"
)

var channelNameRegex = regexp.MustCompile(`^[a-z][a-z0-9.-]{0,248}$`)

func ValidateConfig(config Config) error {
	if config.Output == "" {
		return &ValidationError{
			RuleID: "R01",
			Rule:   "Empty Output",
			Detail: "the output directory cannot be empty",
		}
	}

	if len(config.Organizations) == 0 {
		return &ValidationError{
			RuleID: "R02",
			Rule:   "No Organizations",
			Detail: "at least one organization must be defined",
		}
	}

	for i := range config.Organizations {
		organization := &config.Organizations[i]

		if organization.Name == "" {
			return &ValidationError{
				RuleID: "R03",
				Rule:   "Empty Org Name",
				Detail: fmt.Sprintf("name of the organization index %d is undefined", i),
			}
		}

		if organization.Domain == "" {
			return &ValidationError{
				RuleID: "R04",
				Rule:   "Empty Org Domain",
				Detail: fmt.Sprintf("domain of the organization index %d is undefined", i),
			}
		}

		if organization.CertificateAuthority.ExposePort < 0 {
			return &ValidationError{
				RuleID: "R05",
				Rule:   "Invalid CA Port",
				Detail: fmt.Sprintf("expose port of the certificate authority of the organization %s should be greater than zero", organization.Name),
			}
		}

		for j, peer := range organization.Peers {
			if peer.Name == "" {
				return &ValidationError{
					RuleID: "R06",
					Rule:   "Empty Peer Name",
					Detail: fmt.Sprintf("name of the peer index %d of the organization %s is undefined", j, organization.Name),
				}
			}

			if peer.Subdomain == "" {
				return &ValidationError{
					RuleID: "R07",
					Rule:   "Empty Peer Subdomain",
					Detail: fmt.Sprintf("subdomain of the peer %s of the organization %s is undefined", peer.Name, organization.Name),
				}
			}

			if peer.ExposePort < 0 {
				return &ValidationError{
					RuleID: "R08",
					Rule:   "Invalid Peer Port",
					Detail: fmt.Sprintf("expose port of the peer %s of the organization %s should be greater than zero", peer.Name, organization.Name),
				}
			}
		}

		for j, orderer := range organization.Orderers {
			if orderer.Name == "" {
				return &ValidationError{
					RuleID: "R09",
					Rule:   "Empty Orderer Name",
					Detail: fmt.Sprintf("name of the orderer index %d of the organization %s is undefined", j, organization.Name),
				}
			}

			if orderer.Subdomain == "" {
				return &ValidationError{
					RuleID: "R10",
					Rule:   "Empty Orderer Subdomain",
					Detail: fmt.Sprintf("subdomain of the orderer %s of the organization %s is undefined", orderer.Name, organization.Name),
				}
			}

			if orderer.ExposePort < 0 {
				return &ValidationError{
					RuleID: "R11",
					Rule:   "Invalid Orderer Port",
					Detail: fmt.Sprintf("expose port of the orderer %s of the organization %s should be greater than zero", orderer.Name, organization.Name),
				}
			}
		}
	}

	for _, profile := range config.Profiles {
		if len(profile.Organizations) == 0 {
			return &ValidationError{
				RuleID: "R15",
				Rule:   "Empty Profile Orgs",
				Detail: fmt.Sprintf("profile %s must include at least one organization", profile.Name),
			}
		}
	}

	for _, ch := range config.Channels {
		if ch.Name == "" {
			return &ValidationError{
				RuleID: "R17",
				Rule:   "Empty Channel Name",
				Detail: "channel name cannot be empty",
			}
		}

		if ch.Profile == "" {
			return &ValidationError{
				RuleID: "R16",
				Rule:   "Empty Channel Profile",
				Detail: fmt.Sprintf("channel %s must reference a profile", ch.Name),
			}
		}

		if !channelNameRegex.MatchString(ch.Name) {
			return &ValidationError{
				RuleID: "R29",
				Rule:   "Invalid Channel Name",
				Detail: fmt.Sprintf("invalid channel name: %s (must be lowercase alphanumeric, start with a letter, and contain only '.', '-', or alphanumeric characters, max 249 characters)", ch.Name),
			}
		}

		for i := range ch.Chaincodes {
			chaincode := &ch.Chaincodes[i]

			if chaincode.Name == "" {
				return &ValidationError{
					RuleID: "R12",
					Rule:   "Empty Chaincode Name",
					Detail: fmt.Sprintf("name of the chaincode %d is empty", i),
				}
			}

			if chaincode.Path == "" {
				return &ValidationError{
					RuleID: "R13",
					Rule:   "Empty Chaincode Path",
					Detail: fmt.Sprintf("path of the chaincode %d is empty", i),
				}
			}

			if chaincode.Version == "" {
				return &ValidationError{
					RuleID: "R14",
					Rule:   "Empty Chaincode Version",
					Detail: fmt.Sprintf("version of the chaincode %d is empty", i),
				}
			}
		}
	}

	if _, ok := CapabilityMap[config.Capabilities.Channel]; !ok {
		return &ValidationError{
			RuleID: "R18",
			Rule:   "Unsupported Channel Capability",
			Detail: fmt.Sprintf("unsupported channel capability: %s", config.Capabilities.Channel),
		}
	}

	if _, ok := CapabilityMap[config.Capabilities.Application]; !ok {
		return &ValidationError{
			RuleID: "R19",
			Rule:   "Unsupported Application Capability",
			Detail: fmt.Sprintf("unsupported application capability: %s", config.Capabilities.Application),
		}
	}

	if _, ok := CapabilityMap[config.Capabilities.Orderer]; !ok {
		return &ValidationError{
			RuleID: "R20",
			Rule:   "Unsupported Orderer Capability",
			Detail: fmt.Sprintf("unsupported orderer capability: %s", config.Capabilities.Orderer),
		}
	}

	organizationNames := make(map[string]struct{})
	exposedPorts := make(map[int]portOwner)
	bootstrapCount := 0
	hasOrderer := false

	for i := range config.Organizations {
		organization := &config.Organizations[i]

		if _, exists := organizationNames[organization.Name]; exists {
			return &ValidationError{
				RuleID: "R21",
				Rule:   "Duplicate Org Name",
				Detail: fmt.Sprintf("duplicate organization name: %s", organization.Name),
			}
		}

		organizationNames[organization.Name] = struct{}{}

		if organization.Bootstrap {
			bootstrapCount++
		}

		if len(organization.Orderers) > 0 {
			hasOrderer = true
		}

		caPort := organization.CertificateAuthority.ExposePort
		if caPort > 0 {
			if existingOwner, exists := exposedPorts[caPort]; exists {
				newOwner := portOwner{ownerType: "Certificate Authority", name: organization.Name}
				return &ValidationError{
					RuleID: "R28",
					Rule:   "Exposed Port Conflict",
					Detail: fmt.Sprintf("Port %d is assigned to both %s and %s.", caPort, existingOwner, newOwner),
				}
			}
			exposedPorts[caPort] = portOwner{ownerType: "Certificate Authority", name: organization.Name}
		}

		for _, peer := range organization.Peers {
			if err := validateBinary(peer.Version, MinBinaryVersion[config.Capabilities.Channel]); err != nil {
				return &ValidationError{
					RuleID: "R22",
					Rule:   "Invalid Peer Version",
					Detail: fmt.Sprintf("peer version of org %s invalid: %v", organization.Name, err),
				}
			}

			if peer.ExposePort > 0 {
				if existingOwner, exists := exposedPorts[peer.ExposePort]; exists {
					newOwner := portOwner{ownerType: "peer", name: peer.Name}
					return &ValidationError{
						RuleID: "R28",
						Rule:   "Exposed Port Conflict",
						Detail: fmt.Sprintf("Port %d is assigned to both %s and %s.", peer.ExposePort, existingOwner, newOwner),
					}
				}
				exposedPorts[peer.ExposePort] = portOwner{ownerType: "peer", name: peer.Name}
			}
		}

		for _, orderer := range organization.Orderers {
			if err := validateBinary(orderer.Version, MinBinaryVersion[config.Capabilities.Channel]); err != nil {
				return &ValidationError{
					RuleID: "R23",
					Rule:   "Invalid Orderer Version",
					Detail: fmt.Sprintf("orderer version of org %s invalid: %v", organization.Name, err),
				}
			}

			if orderer.ExposePort > 0 {
				if existingOwner, exists := exposedPorts[orderer.ExposePort]; exists {
					newOwner := portOwner{ownerType: "orderer", name: orderer.Name}
					return &ValidationError{
						RuleID: "R28",
						Rule:   "Exposed Port Conflict",
						Detail: fmt.Sprintf("Port %d is assigned to both %s and %s.", orderer.ExposePort, existingOwner, newOwner),
					}
				}
				exposedPorts[orderer.ExposePort] = portOwner{ownerType: "orderer", name: orderer.Name}
			}
		}
	}

	if !hasOrderer {
		return &ValidationError{
			RuleID: "R24",
			Rule:   "No Orderer Topology",
			Detail: "at least one orderer must be configured",
		}
	}

	if bootstrapCount > 1 {
		return &ValidationError{
			RuleID: "R25",
			Rule:   "Multiple Bootstrap Orgs",
			Detail: "exactly one bootstrap organization must be defined",
		}
	}

	for _, profile := range config.Profiles {
		consensusType := profile.Consensus.Type

		if consensusType != "" && consensusType != "etcdraft" && consensusType != "BFT" {
			return &ValidationError{
				RuleID: "R26",
				Rule:   "Invalid Consensus Type",
				Detail: fmt.Sprintf("invalid consensus type for the profile %s", profile.Name),
			}
		}

		for _, orgName := range profile.Organizations {
			if _, ok := organizationNames[orgName]; !ok {
				return &ValidationError{
					RuleID: "R27",
					Rule:   "Profile References Undefined Org",
					Detail: fmt.Sprintf("organization not defined: %s", orgName),
				}
			}
		}
	}

	return nil
}

func validateBinary(version string, minVersion string) error {
	if version == "" {
		return nil
	}

	vParts := parseVersion(version)
	minParts := parseVersion(minVersion)

	for i := 0; i < 3; i++ {
		if vParts[i] > minParts[i] {
			return nil
		}
		if vParts[i] < minParts[i] {
			return fmt.Errorf("version %s is lower than required %s", version, minVersion)
		}
	}

	return nil
}

func parseVersion(v string) [3]int {
	var result [3]int
	parts := strings.Split(v, ".")

	for i := 0; i < len(parts) && i < 3; i++ {
		fmt.Sscanf(parts[i], "%d", &result[i])
	}

	return result
}
