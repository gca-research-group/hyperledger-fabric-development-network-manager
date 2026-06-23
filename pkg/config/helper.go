package config

import "fmt"

func ResolveOrganizationMSPID(organization Organization) string {
	return fmt.Sprintf("%sMSP", organization.Name)
}

func ResolveOrdererMSPID(organization Organization) string {
	return fmt.Sprintf("%sOrdererMSP", organization.Name)
}

func ResolveChaincodes(config Config) []Chaincode {
	items := ChaincodesMap(config)
	chaincodes := make([]Chaincode, 0, len(items))

	for _, chaincode := range items {
		chaincodes = append(chaincodes, chaincode)
	}

	return chaincodes
}

func ChaincodesMap(config Config) map[string]Chaincode {
	chaincodes := make(map[string]Chaincode)

	for _, channel := range config.Channels {
		for _, chaincode := range channel.Chaincodes {
			if _, exists := chaincodes[chaincode.Name]; !exists {
				chaincodes[chaincode.Name] = chaincode
			}
		}
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
