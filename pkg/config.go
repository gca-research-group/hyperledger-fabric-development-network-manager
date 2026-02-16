package pkg

type Orderer struct {
	Name       string `yaml:"name" json:"name" toml:"name"`
	Subdomain  string `yaml:"subdomain" json:"subdomain" toml:"subdomain"`
	ExposePort int    `yaml:"exposePort" json:"exposePort" toml:"exposePort"`
	Port       int    `yaml:"port" json:"port" toml:"port"`
}

type Peer struct {
	Name       string `yaml:"name" json:"name" toml:"name"`
	Subdomain  string `yaml:"subdomain" json:"subdomain" toml:"subdomain"`
	Port       int    `yaml:"port" json:"port" toml:"port"`
	ExposePort int    `yaml:"exposePort" json:"exposePort" toml:"exposePort"`
	IsAnchor   bool   `yaml:"isAnchor" json:"isAnchor" toml:"isAnchor"`
}

type CertificateAuthority struct {
	ExposePort int `yaml:"exposePort,omitempty" json:"exposePort,omitempty" toml:"bootstrap,omitempty"`
}

type Organization struct {
	Name                 string               `yaml:"name" json:"name" toml:"name"`
	Domain               string               `yaml:"domain" json:"domain" toml:"domain"`
	Peers                []Peer               `yaml:"peers" json:"peers" toml:"peers"`
	Orderers             []Orderer            `yaml:"orderers" json:"orderers" toml:"orderers"`
	Users                int                  `yaml:"users" json:"users" toml:"users"`
	CertificateAuthority CertificateAuthority `yaml:"certificateAuthority" json:"certificateAuthority" toml:"certificateAuthority"`
	Bootstrap            bool                 `yaml:"bootstrap,omitempty" json:"bootstrap,omitempty" toml:"bootstrap,omitempty"`
}

type Channel struct {
	Name    string  `yaml:"name" json:"name" toml:"name"`
	Profile Profile `yaml:"profile" json:"profile" toml:"profile"`
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
	Channels      []Channel      `yaml:"channels" json:"channels" toml:"channels"`
}
