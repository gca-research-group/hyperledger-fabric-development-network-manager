package config

import (
	"fmt"
	"regexp"
	"strings"
)

var channelNameRegex = regexp.MustCompile(`^[a-z][a-z0-9.-]{0,248}$`)

func ValidateSemantic(config Config) error {
	if _, ok := CapabilityMap[config.Capabilities.Channel]; !ok {
		return &ValidationError{
			Layer:  "Semantic Analysis",
			RuleID: "SM01",
			Rule:   "Unsupported Channel Capability",
			Detail: fmt.Sprintf("unsupported channel capability: %s", config.Capabilities.Channel),
		}
	}

	if _, ok := CapabilityMap[config.Capabilities.Application]; !ok {
		return &ValidationError{
			Layer:  "Semantic Analysis",
			RuleID: "SM02",
			Rule:   "Unsupported Application Capability",
			Detail: fmt.Sprintf("unsupported application capability: %s", config.Capabilities.Application),
		}
	}

	if _, ok := CapabilityMap[config.Capabilities.Orderer]; !ok {
		return &ValidationError{
			Layer:  "Semantic Analysis",
			RuleID: "SM03",
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
				Layer:  "Semantic Analysis",
				RuleID: "SM04",
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
					Layer:  "Semantic Analysis",
					RuleID: "SM11",
					Rule:   "Exposed Port Conflict",
					Detail: fmt.Sprintf("Port %d is assigned to both %s and %s.", caPort, existingOwner, newOwner),
				}
			}
			exposedPorts[caPort] = portOwner{ownerType: "Certificate Authority", name: organization.Name}
		}

		for _, peer := range organization.Peers {
			if err := validateBinary(peer.Version, MinBinaryVersion[config.Capabilities.Channel]); err != nil {
				return &ValidationError{
					Layer:  "Semantic Analysis",
					RuleID: "SM05",
					Rule:   "Invalid Peer Version",
					Detail: fmt.Sprintf("peer version of org %s invalid: %v", organization.Name, err),
				}
			}

			if peer.ExposePort > 0 {
				if existingOwner, exists := exposedPorts[peer.ExposePort]; exists {
					newOwner := portOwner{ownerType: "peer", name: peer.Name}
					return &ValidationError{
						Layer:  "Semantic Analysis",
						RuleID: "SM11",
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
					Layer:  "Semantic Analysis",
					RuleID: "SM06",
					Rule:   "Invalid Orderer Version",
					Detail: fmt.Sprintf("orderer version of org %s invalid: %v", organization.Name, err),
				}
			}

			if orderer.ExposePort > 0 {
				if existingOwner, exists := exposedPorts[orderer.ExposePort]; exists {
					newOwner := portOwner{ownerType: "orderer", name: orderer.Name}
					return &ValidationError{
						Layer:  "Semantic Analysis",
						RuleID: "SM11",
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
			Layer:  "Semantic Analysis",
			RuleID: "SM07",
			Rule:   "No Orderer Topology",
			Detail: "at least one orderer must be configured",
		}
	}

	if bootstrapCount > 1 {
		return &ValidationError{
			Layer:  "Semantic Analysis",
			RuleID: "SM08",
			Rule:   "Multiple Bootstrap Orgs",
			Detail: "exactly one bootstrap organization must be defined",
		}
	}

	for _, profile := range config.Profiles {
		consensusType := profile.Consensus.Type

		if consensusType != "" && consensusType != "etcdraft" && consensusType != "BFT" {
			return &ValidationError{
				Layer:  "Semantic Analysis",
				RuleID: "SM09",
				Rule:   "Invalid Consensus Type",
				Detail: fmt.Sprintf("invalid consensus type for the profile %s", profile.Name),
			}
		}

		for _, orgName := range profile.Organizations {
			if _, ok := organizationNames[orgName]; !ok {
				return &ValidationError{
					Layer:  "Semantic Analysis",
					RuleID: "SM10",
					Rule:   "Profile References Undefined Org",
					Detail: fmt.Sprintf("organization not defined: %s", orgName),
				}
			}
		}
	}

	for _, ch := range config.Channels {
		if !channelNameRegex.MatchString(ch.Name) {
			return &ValidationError{
				Layer:  "Semantic Analysis",
				RuleID: "SM12",
				Rule:   "Invalid Channel Name",
				Detail: fmt.Sprintf("invalid channel name: %s (must be lowercase alphanumeric, start with a letter, and contain only '.', '-', or alphanumeric characters, max 249 characters)", ch.Name),
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
