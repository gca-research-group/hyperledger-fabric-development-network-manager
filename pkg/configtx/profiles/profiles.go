package profiles

type Orderer struct {
	Orderer       string   `yaml:"Orderer"`
	Organizations []string `yaml:"Organizations"`
}

type Organizations struct {
	Organizations []string `yaml:"Organizations"`
}

type Application struct {
	Application   string   `yaml:"Application"`
	Organizations []string `yaml:"Organizations"`
}

type OrdererProfile struct {
	Channel     string                   `yaml:"Channel"`
	Orderer     Orderer                  `yaml:"Orderer"`
	Consortiums map[string]Organizations `yaml:"Consortiums"`
	Application Application              `yaml:"Application"`
}

type ChannelProfile struct {
	Channel     string      `yaml:"Channel"`
	Consortium  string      `yaml:"Consortium"`
	Application Application `yaml:"Application"`
}
