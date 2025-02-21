package pkg

type Orderer struct {
	Name   string
	Domain string
	Port   int
}

type Peer struct {
	Name   string
	Domain string
	Port   int
	Peers  int
	Users  int
}

type Channel struct {
	Name          string
	Organizations []string
}

type Config struct {
	Orderers []Orderer
	Peers    []Peer
	Channels []Channel
}
