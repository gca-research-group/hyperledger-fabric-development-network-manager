package organization

type Policy struct {
	Type string `yaml:"Type"`
	Rule string `yaml:"Rule"`
}

type OrdererPolicies struct {
	Readers Policy `yaml:"Readers"`
	Writers Policy `yaml:"Writers"`
	Admins  Policy `yaml:"Admins"`
}

type PeerPolicies struct {
	Readers     Policy `yaml:"Readers"`
	Writers     Policy `yaml:"Writers"`
	Admins      Policy `yaml:"Admins"`
	Endorsement Policy `yaml:"Endorsement"`
}

type AnchorPeer struct {
	Host string `yaml:"Host"`
	Port int    `yaml:"Port"`
}

type Orderer struct {
	Name     string          `yaml:"Name"`
	ID       string          `yaml:"ID"`
	MSPDir   string          `yaml:"MSPDir"`
	Policies OrdererPolicies `yaml:"Policies"`
}

type Peer struct {
	Name        string       `yaml:"Name"`
	ID          string       `yaml:"ID"`
	MSPDir      string       `yaml:"MSPDir"`
	Policies    PeerPolicies `yaml:"Policies"`
	AnchorPeers []AnchorPeer `yaml:"AnchorPeers"`
}

type Organization interface{}
