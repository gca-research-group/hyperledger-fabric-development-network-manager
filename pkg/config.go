package pkg

type Orderer struct {
	Name     string `yaml:"name"`
	Hostname string `yaml:"hostname"`
	Port     int
}

type AnchorPeer struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type CertificateAuthority struct {
	ExposePort int `yaml:"exposePort,omitempty"`
}

type Organization struct {
	Name                 string               `yaml:"name"`
	Domain               string               `yaml:"domain"`
	AnchorPeer           AnchorPeer           `yaml:"anchorPeer"`
	Peers                int                  `yaml:"peers"`
	Users                int                  `yaml:"users"`
	CertificateAuthority CertificateAuthority `yaml:"certificateAuthority"`
	Orderers             []Orderer            `yaml:"orderers"`
	Bootstrap            bool                 `yaml:"bootstrap,omitempty"`
}

type Profile struct {
	Name          string   `yaml:"name"`
	Organizations []string `yaml:"organizations"`
}

type Config struct {
	Output        string         `yaml:"output"`
	Network       string         `yaml:"network"`
	Organizations []Organization `yaml:"organizations"`
	Profiles      []Profile      `yaml:"profiles"`
}
