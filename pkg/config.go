package pkg

type Orderer struct {
	Name     string `yaml:"name" json:"name" toml:"name"`
	Hostname string `yaml:"hostname" json:"hostname" toml:"hostname"`
	Port     int
}

type AnchorPeer struct {
	Host string `yaml:"host" json:"host" toml:"host"`
	Port int    `yaml:"port" json:"port" toml:"port"`
}

type CertificateAuthority struct {
	ExposePort int `yaml:"exposePort,omitempty" json:"exposePort,omitempty" toml:"bootstrap,omitempty"`
}

type Organization struct {
	Name                 string               `yaml:"name" json:"name" toml:"name"`
	Domain               string               `yaml:"domain" json:"domain" toml:"domain"`
	AnchorPeer           AnchorPeer           `yaml:"anchorPeer" json:"anchorPeer" toml:"anchorPeer"`
	Peers                int                  `yaml:"peers" json:"peers" toml:"peers"`
	Users                int                  `yaml:"users" json:"users" toml:"users"`
	CertificateAuthority CertificateAuthority `yaml:"certificateAuthority" json:"certificateAuthority" toml:"certificateAuthority"`
	Orderers             []Orderer            `yaml:"orderers" json:"orderers" toml:"orderers"`
	Bootstrap            bool                 `yaml:"bootstrap,omitempty" json:"bootstrap,omitempty" toml:"bootstrap,omitempty"`
}

type Profile struct {
	Name          string   `yaml:"name" json:"name" toml:"name"`
	Organizations []string `yaml:"organizations" json:"organizations" toml:"organizations"`
}

type Config struct {
	Output        string         `yaml:"output" json:"output" toml:"output"`
	Network       string         `yaml:"network" json:"network" toml:"network"`
	Organizations []Organization `yaml:"organizations" json:"organizations" toml:"organizations"`
	Profiles      []Profile      `yaml:"profiles" json:"profiles" toml:"profiles"`
}
