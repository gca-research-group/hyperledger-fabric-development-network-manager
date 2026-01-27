package pkg

type Orderer struct {
	Name     string
	Hostname string
	Domain   string
	Port     int
}

type AnchorPeer struct {
	Host string
	Port int
}

type CertificateAuthority struct {
	ExposePort int
}

type Organization struct {
	Name                 string
	Domain               string
	AnchorPeer           AnchorPeer
	Peers                int
	Users                int
	CertificateAuthority CertificateAuthority
}

type Profile struct {
	Name          string
	Organizations []string
}

type Docker struct {
	NetworkName string
}

type Config struct {
	Orderers      []Orderer
	Organizations []Organization
	Profiles      []Profile
	Docker        Docker
}
