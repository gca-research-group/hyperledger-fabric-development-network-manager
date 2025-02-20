package capabilities

type Capability struct {
	V2_0 bool `yaml:"V2_0"`
}

type Capabilities struct {
	Application Capability `yaml:"Application"`
	Orderer     Capability `yaml:"Orderer"`
	Channel     Capability `yaml:"Channel"`
}

func NewCapabilities() Capabilities {
	return Capabilities{
		Application: Capability{V2_0: true},
		Orderer:     Capability{V2_0: true},
		Channel:     Capability{V2_0: true},
	}
}
