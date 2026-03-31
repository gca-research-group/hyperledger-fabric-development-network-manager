package network

const (
	ORDERER_MSP_DIR        = "./%[1]s/ordererOrganizations/%[1]s/msp"
	PEER_MSP_DIR           = "./%[1]s/peerOrganizations/%[1]s/msp"
	ORDERER_TLS_SERVER_CRT = "./%[1]s/ordererOrganizations/%[1]s/orderers/%[2]s.%[1]s/tls/server.crt"
	ORDERER_MSP_SIGNCERT   = "./%[1]s/ordererOrganizations/%[1]s/orderers/%[2]s.%[1]s/msp/signcerts/cert.pem"

	CONTAINER_TLS_PEER_SERVER_CRT = "%[1]s/%[2]s/peerOrganizations/%[2]s/peers/%[3]s.%[2]s/tls/server.crt"
	CONTAINER_TLS_PEER_SERVER_KEY = "%[1]s/%[2]s/peerOrganizations/%[2]s/peers/%[3]s.%[2]s/tls/server.key"
	CONTAINER_TLS_PEER_SERVER_CA  = "%[1]s/%[2]s/peerOrganizations/%[2]s/peers/%[3]s.%[2]s/tls/ca.crt"
)
