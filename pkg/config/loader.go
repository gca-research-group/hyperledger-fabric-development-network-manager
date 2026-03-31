package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

type capabilityLevel int

const (
	V2_0 capabilityLevel = iota + 1
	V2_5
	V3_0
)

var capabilityMap = map[string]capabilityLevel{
	"V2_0": V2_0,
	"V2_5": V2_5,
	"V3_0": V3_0,
}

var minBinaryVersion = map[string]string{
	"V2_0": "2.0.0",
	"V2_5": "2.5.0",
	"V3_0": "3.0.0",
}

var defaultVersionByCapability = map[string]string{
	"V2_0": "2.5.0",
	"V2_5": "2.5.0",
	"V3_0": "3.1.4",
}

func LoadConfigFromPath(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("loading config file: %w", err)
	}

	var config Config

	switch strings.ToLower(filepath.Ext(path)) {
	case ".json":
		err = json.Unmarshal(data, &config)
	case ".yml", ".yaml":
		err = yaml.Unmarshal(data, &config)
	case ".toml":
		err = toml.Unmarshal(data, &config)
	default:
		return nil, errors.New("unsupported config format")
	}

	if err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	setUpDefaultValues(&config)

	return &config, nil
}

func setUpDefaultValues(config *Config) {
	for i := range config.Channels {
		channel := &config.Channels[i]
		if channel.Profile.Consensus.Type == "" {
			channel.Profile.Consensus.Type = "etcdraft"
		}
	}

	for i := range config.Profiles {
		profile := &config.Profiles[i]
		if profile.Consensus.Type == "" {
			profile.Consensus.Type = "etcdraft"
		}
	}

	channelCapability, _ := capabilityMap[config.Capabilities.Channel]
	applicationCapability, _ := capabilityMap[config.Capabilities.Application]
	ordererCapability, _ := capabilityMap[config.Capabilities.Orderer]

	hasBootstrapOrganization := false

	for i := range config.Organizations {
		organization := &config.Organizations[i]

		for j := range organization.Orderers {
			orderer := &organization.Orderers[j]
			if orderer.Port == 0 {
				orderer.Port = 7050
			}
		}

		if organization.Version.Peer == "" {
			if channelCapability > applicationCapability {
				organization.Version.Peer = defaultVersionByCapability[config.Capabilities.Channel]
			} else {
				organization.Version.Peer = defaultVersionByCapability[config.Capabilities.Application]
			}
		}

		if organization.Version.Orderer == "" {
			if channelCapability > ordererCapability {
				organization.Version.Orderer = defaultVersionByCapability[config.Capabilities.Channel]
			} else {
				organization.Version.Orderer = defaultVersionByCapability[config.Capabilities.Orderer]
			}
		}

		if organization.Bootstrap {
			hasBootstrapOrganization = true
		}

		hasAnchorPeer := false
		for i := range organization.Peers {
			peer := &organization.Peers[i]
			if peer.IsAnchor {
				hasAnchorPeer = true
			}
		}

		if !hasAnchorPeer && len(organization.Peers) > 0 {
			organization.Peers[0].IsAnchor = true
		}
	}

	if !hasBootstrapOrganization {
		for i, organization := range config.Organizations {
			if len(organization.Orderers) != 0 {
				config.Organizations[i].Bootstrap = true
				break
			}
		}
	}
}

func validateConfig(config Config) error {
	if len(config.Organizations) == 0 {
		return fmt.Errorf("at least one organization must be defined")
	}

	if _, ok := capabilityMap[config.Capabilities.Channel]; !ok {
		return fmt.Errorf("unsupported channel capability: %s", config.Capabilities.Channel)
	}

	if _, ok := capabilityMap[config.Capabilities.Application]; !ok {
		return fmt.Errorf("unsupported application capability: %s", config.Capabilities.Application)
	}

	if _, ok := capabilityMap[config.Capabilities.Orderer]; !ok {
		return fmt.Errorf("unsupported orderer capability: %s", config.Capabilities.Orderer)
	}

	organizationNames := make(map[string]struct{})
	bootstrapCount := 0
	hasOrderer := false

	for i := range config.Organizations {
		organization := &config.Organizations[i]

		if _, exists := organizationNames[organization.Name]; exists {
			return fmt.Errorf("duplicate organization name: %s", organization.Name)
		}

		organizationNames[organization.Name] = struct{}{}

		if organization.Bootstrap {
			bootstrapCount++
		}

		if len(organization.Orderers) > 0 {
			hasOrderer = true
		}

		if err := validateBinary(organization.Version.Peer, minBinaryVersion[config.Capabilities.Channel]); err != nil {
			return fmt.Errorf("peer version of org %s invalid: %w", organization.Name, err)
		}

		if err := validateBinary(organization.Version.Orderer, minBinaryVersion[config.Capabilities.Channel]); err != nil {
			return fmt.Errorf("orderer version of org %s invalid: %w", organization.Name, err)
		}
	}

	if !hasOrderer {
		return fmt.Errorf("at least one orderer must be configured")
	}

	if bootstrapCount > 1 {
		return fmt.Errorf("exactly one bootstrap organization must be defined")
	}

	for _, profile := range config.Profiles {
		if len(profile.Organizations) == 0 {
			return fmt.Errorf("profile %s must include at least one organization", profile.Name)
		}

		for _, orgName := range profile.Organizations {
			if _, ok := organizationNames[orgName]; !ok {
				return fmt.Errorf("organization not defined: %s", orgName)
			}
		}
	}

	for _, ch := range config.Channels {
		if ch.Profile.Name == "" {
			return fmt.Errorf("channel %s must reference a profile", ch.Name)
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
