package channel

type Policy struct {
	Type string `yaml:"Type"`
	Rule string `yaml:"Rule"`
}

type Policies struct {
	Readers Policy `yaml:"Readers"`
	Writers Policy `yaml:"Writers"`
	Admins  Policy `yaml:"Admins"`
}

type Channel struct {
	Policies     Policies `yaml:"Policies"`
	Capabilities []string `yaml:"Capabilities"`
}

func NewChannel() Channel {
	return Channel{
		Policies: Policies{
			Readers: Policy{Type: "ImplicitMeta", Rule: "ANY Readers"},
			Writers: Policy{Type: "ImplicitMeta", Rule: "ANY Writers"},
			Admins:  Policy{Type: "ImplicitMeta", Rule: "ANY Admins"},
		},
		Capabilities: []string{"<<: *ChannelCapabilities"},
	}
}
