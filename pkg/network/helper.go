package network

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/constants"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/compose"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/config"
)

func ResolveOrdererTLSConnection(organizations []config.Organization) (string, string) {
	var orderer config.Orderer
	var ordererDomain string
	var ordererPort int

	for _, organization := range organizations {
		if len(organization.Orderers) > 0 {
			orderer = organization.Orderers[0]
			ordererDomain = organization.Domain
			ordererPort = orderer.Port

			if ordererPort == 0 {
				ordererPort = 7050
			}
			break
		}
	}

	ordererAddress := fmt.Sprintf("%s.%s:%d", orderer.Subdomain, ordererDomain, ordererPort)
	caFile := fmt.Sprintf("%[1]s/%[2]s/ordererOrganizations/%[2]s/orderers/%[3]s.%[2]s/tls/ca.crt", constants.DEFAULT_FABRIC_DIRECTORY, ordererDomain, orderer.Subdomain)

	return ordererAddress, caFile
}

func ResolvePeersTLSConnection(organizations []config.Organization) [][2]string {

	var data [][2]string

	for _, organization := range organizations {
		for _, peer := range organization.Peers {
			address := fmt.Sprintf("%s.%s:%d", peer.Subdomain, organization.Domain, compose.ResolvePeerPort(peer.Port))
			caFile := fmt.Sprintf("%[1]s/%[2]s/peerOrganizations/%[2]s/peers/%[3]s.%[2]s/tls/ca.crt", constants.DEFAULT_FABRIC_DIRECTORY, organization.Domain, peer.Subdomain)
			data = append(data, [2]string{address, caFile})
		}
	}

	return data
}
