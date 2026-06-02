package config

import "fmt"

func ResolveOrganizationMSPID(organization Organization) string {
	return fmt.Sprintf("%sMSP", organization.Name)
}

func ResolveOrdererMSPID(organization Organization) string {
	return fmt.Sprintf("%sOrdererMSP", organization.Name)
}

func ResolveAllChaincodes(config Config) []Chaincode {
	chaincodesMap := make(map[string]Chaincode)

	for _, channel := range config.Channels {
		for _, chaincode := range channel.Chaincodes {
			if _, exists := chaincodesMap[chaincode.Name]; !exists {
				chaincodesMap[chaincode.Name] = chaincode
			}
		}
	}

	chaincodes := make([]Chaincode, 0, len(chaincodesMap))
	for _, chaincode := range chaincodesMap {
		chaincodes = append(chaincodes, chaincode)
	}

	return chaincodes
}

func ProfilesMap(config Config) map[string]Profile {
	profiles := make(map[string]Profile)

	for _, profile := range config.Profiles {
		profiles[profile.Name] = profile
	}

	return profiles
}
