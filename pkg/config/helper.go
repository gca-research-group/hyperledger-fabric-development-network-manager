package config

import "fmt"

func ResolveOrganizationMSPID(organization Organization) string {
	return fmt.Sprintf("%sMSP", organization.Name)
}

func ResolveOrdererMSPID(orderer Orderer) string {
	return fmt.Sprintf("%sMSP", orderer.Name)
}
