package application

type Policy struct {
	Type string `yaml:"Type"`
	Rule string `yaml:"Rule"`
}

type Policies struct {
	LifecycleEndorsement Policy `yaml:"LifecycleEndorsement"`
	Endorsement          Policy `yaml:"Endorsement"`
	Readers              Policy `yaml:"Readers"`
	Writers              Policy `yaml:"Writers"`
	Admins               Policy `yaml:"Admins"`
}

type Application struct {
	Policies      Policies `yaml:"Policies"`
	Organizations []string `yaml:"Organizations"`
	Capabilities  []string `yaml:"Capabilities"`
}

func NewApplication(organizations []string) Application {
	return Application{
		Policies: Policies{
			LifecycleEndorsement: Policy{Type: "ImplicitMeta", Rule: "ANY Endorsement"},
			Endorsement:          Policy{Type: "ImplicitMeta", Rule: "ANY Endorsement"},
			Readers:              Policy{Type: "ImplicitMeta", Rule: "ANY Readers"},
			Writers:              Policy{Type: "ImplicitMeta", Rule: "ANY Writers"},
			Admins:               Policy{Type: "ImplicitMeta", Rule: "ANY Admins"},
		},
		Organizations: organizations,
		Capabilities:  []string{"<<: *ApplicationCapabilities"},
	}
}
