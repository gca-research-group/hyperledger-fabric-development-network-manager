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

func LoadConfigFromPath(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Error when loading the config file: %v\n", err)
	}

	var config Config

	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".json":
		err = json.Unmarshal(data, &config)
	case ".yml":
		err = yaml.Unmarshal(data, &config)
	case ".yaml":
		err = yaml.Unmarshal(data, &config)
	case ".toml":
		err = toml.Unmarshal(data, &config)
	default:
		return nil, errors.New("unsupported config format")
	}

	if err != nil {
		return nil, fmt.Errorf("Error when loading the config file: %v\n", err)
	}

	if len(config.Organizations) == 0 {
		return nil, fmt.Errorf("You must define at least one organization.")
	}

	supportedCapabilityVersions := []string{"V2_0", "V2_5", "V3_0"}

	supportedVersions := map[string][]string{
		"V2_0": {"2.0", "2.1", "2.2", "2.3", "2.4", "2.5"},
		"V2_5": {"2.5"},
		"V3_0": {"3.0"},
	}

	supportedApplicationVersions := supportedVersions[config.Capabilties.Application]
	supportedOrdererVersions := supportedVersions[config.Capabilties.Orderer]

	organizationNames := make(map[string]int)

	isApplicationCapabilityValid := false
	isChannelCapabilityValid := false
	isOrdererCapabilityValid := false

	for _, capabilityVersion := range supportedCapabilityVersions {
		if capabilityVersion == config.Capabilties.Application {
			isApplicationCapabilityValid = true
		}

		if capabilityVersion == config.Capabilties.Channel {
			isChannelCapabilityValid = true
		}

		if capabilityVersion == config.Capabilties.Orderer {
			isOrdererCapabilityValid = true
		}
	}

	if !isApplicationCapabilityValid {
		return nil, fmt.Errorf("Application capability version not suppported. Supported versions: %v.", supportedCapabilityVersions)
	}

	if !isChannelCapabilityValid {
		return nil, fmt.Errorf("Channel capability version not suppported. Supported versions: %v.", supportedCapabilityVersions)
	}

	if !isOrdererCapabilityValid {
		return nil, fmt.Errorf("Orderer capability version not suppported. Supported versions: %v.", supportedCapabilityVersions)
	}

	quantityOfBootstrapOrganizations := 0
	hasOrderers := false

	for _, organization := range config.Organizations {
		if !isVersionSupported(organization.Version.Peer, supportedApplicationVersions) {
			return nil, fmt.Errorf("The peer version %s of the organization %s is not supported since the application capability is %s. The supported versions are %v.", organization.Version.Peer, organization.Name, config.Capabilties.Application, supportedApplicationVersions)
		}

		if !isVersionSupported(organization.Version.Orderer, supportedOrdererVersions) {
			return nil, fmt.Errorf("The orderer version %s of the organization %s is not supported since the orderer capability is %s. The supported versions are %v.", organization.Version.Orderer, organization.Name, config.Capabilties.Orderer, supportedOrdererVersions)
		}

		if organization.Bootstrap {
			quantityOfBootstrapOrganizations += 1
		}

		if len(organization.Orderers) > 0 {
			hasOrderers = true
		}

		organizationNames[organization.Name] += 1
	}

	if !hasOrderers {
		return nil, fmt.Errorf("At least one orderer must be configured.")
	}

	if quantityOfBootstrapOrganizations == 0 {
		return nil, fmt.Errorf("The bootstrap organization must be configured.")
	}

	if quantityOfBootstrapOrganizations > 1 {
		return nil, fmt.Errorf("Just one booststrap organization must be configured.")
	}

	for key, value := range organizationNames {
		if value > 1 {
			return nil, fmt.Errorf("The organization's name should be unique. The organization %s was defined multiple times.", key)
		}
	}

	for _, profile := range config.Profiles {
		if len(profile.Organizations) == 0 {
			return nil, fmt.Errorf("Add at least one organization to the profile %s.", profile.Name)
		}

		for _, organizationName := range profile.Organizations {
			if !isOrganizationDefined(organizationName, config.Organizations) {
				return nil, fmt.Errorf("Organization is not defined: %s.", organizationName)
			}
		}
	}

	for o := range config.Organizations {
		for i := range config.Organizations[o].Orderers {
			if config.Organizations[o].Orderers[i].Port == 0 {
				config.Organizations[o].Orderers[i].Port = 7050
			}
		}
	}

	return &config, nil
}

func isOrganizationDefined(organizationName string, organizations []Organization) bool {
	for _, organization := range organizations {
		if organizationName == organization.Name {
			return true
		}
	}

	return false
}

func isVersionSupported(version string, suportedVersions []string) bool {
	if version == "" {
		return true
	}

	for _, supportedVersion := range suportedVersions {
		if strings.HasPrefix(version, supportedVersion) {
			return true
		}
	}

	return false
}
