package orderer

type Policy struct {
	Type string `yaml:"Type"`
	Rule string `yaml:"Rule"`
}

type Policies struct {
	Readers         Policy `yaml:"Readers"`
	Writers         Policy `yaml:"Writers"`
	Admins          Policy `yaml:"Admins"`
	BlockValidation Policy `yaml:"BlockValidation"`
}

type BatchSize struct {
	MaxMessageCount   int    `yaml:"MaxMessageCount"`
	AbsoluteMaxBytes  string `yaml:"AbsoluteMaxBytes"`
	PreferredMaxBytes string `yaml:"PreferredMaxBytes"`
}

type Orderer struct {
	OrdererType  string    `yaml:"OrdererType"`
	Addresses    []string  `yaml:"Addresses"`
	Capabilities []string  `yaml:"Capabilities"`
	Policies     Policies  `yaml:"Policies"`
	BatchTimeout string    `yaml:"BatchTimeout"`
	BatchSize    BatchSize `yaml:"BatchSize"`
}

func NewOrderer(addresses []string) Orderer {
	return Orderer{
		OrdererType:  "solo",
		Addresses:    addresses,
		Capabilities: []string{"<<: *OrdererCapabilities"},
		Policies: Policies{
			Readers: Policy{
				Type: "ImplicitMeta",
				Rule: "ANY Readers",
			},
			Writers: Policy{
				Type: "ImplicitMeta",
				Rule: "ANY Writers",
			},
			Admins: Policy{
				Type: "ImplicitMeta",
				Rule: "ANY Admins",
			},
			BlockValidation: Policy{
				Type: "ImplicitMeta",
				Rule: "ANY Writers",
			},
		},
		BatchTimeout: "2s",
		BatchSize: BatchSize{
			MaxMessageCount:   10,
			AbsoluteMaxBytes:  "98 MB",
			PreferredMaxBytes: "512 KB",
		},
	}
}
