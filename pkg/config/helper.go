package config

import "fmt"

func ResolveOrganizationMSPID(organization Organization) string {
	return fmt.Sprintf("%sMSP", organization.Name)
}

func ResolveOrdererMSPID(organization Organization) string {
	return fmt.Sprintf("%sOrdererMSP", organization.Name)
}
