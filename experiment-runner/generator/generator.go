package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

type MutationOperator struct {
	RuleID validate.RuleID
	Apply  func(node *yaml.Node)
}

type ScenarioRules struct {
	Scenario string            `json:"scenario"`
	Rules    []validate.RuleID `json:"rules"`
}

const seedYAML = `
output: output/example
network: example

capabilities:
  channel: V2_0
  orderer: V2_0
  application: V2_5

organizations:
  - name: Org1
    bootstrap: true
    domain: org1.example.com
    orderers:
      - name: Orderer
        subdomain: orderer
    peers:
      - name: Peer0
        subdomain: peer0
        exposePort: 7051
    certificateAuthority:
      exposePort: 7054

  - name: Org2
    bootstrap: false
    domain: org2.example.com
    peers:
      - name: Peer0
        subdomain: peer0
        exposePort: 8051

  - name: Org3
    domain: org3.example.com
    peers:
      - name: Peer0
        subdomain: peer0
        exposePort: 9051

profiles:
  - name: DefaultProfile
    organizations:
      - Org1
      - Org2
      - Org3

channels:
  - name: defaultchannel
    profile: DefaultProfile
    chaincodes:
      - name: Asset
        version: "1.0"
        path: samples/chaincodes/asset
        language:
          name: golang
          version: "1.26"

      - name: Product
        version: "1.0"
        path: samples/chaincodes/product

      - name: PrivateAgreement
        version: "1.0"
        path: samples/chaincodes/private-agreement
`

var operators = [][]MutationOperator{
	outputDirectoryOperators,
	organizationDomainRequiredOperators,
	certificateAuthorityPortInvalidOperators,
	peerSubdomainRequiredOperators,
	peerPortInvalidOperators,
	ordererSubdomainRequiredOperators,
	ordererPortInvalidOperators,
	chaincodePathRequiredOperators,
	chaincodeVersionRequiredOperators,
	channelProfileRequiredOperators,
	channelNameInvalidOperators,
	channelCapabilityUnsupportedOperators,
	applicationCapabilityUnsupportedOperators,
	ordererCapabilityUnsupportedOperators,
	peerVersionInvalidOperators,
	ordererVersionInvalidOperators,
	bootstrapOrganizationsMultipleOperators,
	consensusTypeInvalidOperators,
	profileOrganizationUndefinedOperators,
	exposedPortConflictOperators,
	organizationNameDuplicateOperators,
	peerNameRequiredOperators,
	ordererNameRequiredOperators,
	chaincodeNameRequiredOperators,
	channelNameRequiredOperators,
	ordererTopologyRequiredOperators,
	organizationNameRequiredOperators,
	profileOrganizationsRequiredOperators,
	organizationsRequiredOperators,
}

const DefaultMutationCount = 3

func Generate(outputDirectory string, mutationCount int) ([]ScenarioRules, error) {
	seed, err := yaml.FromBytes([]byte(seedYAML))
	if err != nil {
		return nil, fmt.Errorf("parse seed configuration: %w", err)
	}

	scenarios := GenerateCombinations(operators, mutationCount)
	configDirectory := filepath.Join(outputDirectory, "config")
	if err := os.MkdirAll(configDirectory, 0755); err != nil {
		return nil, fmt.Errorf("create config directory: %w", err)
	}

	scenarioRules := make([]ScenarioRules, 0, len(scenarios))

	for i := 0; i < len(scenarios); i++ {
		clone := seed.Clone()
		doc := clone.Document()

		rules := make([]validate.RuleID, 0)

		for j := range scenarios[i] {
			operator := scenarios[i][j]
			rules = append(rules, operator.RuleID)
			operator.Apply(doc)
		}

		scenarioName := fmt.Sprintf(
			"%06d",
			i+1,
		)

		scenarioRules = append(
			scenarioRules,
			ScenarioRules{
				Scenario: scenarioName,
				Rules:    rules,
			},
		)
		if err := clone.ToFile(filepath.Join(configDirectory, scenarioName+".yaml")); err != nil {
			return nil, fmt.Errorf("write scenario %s: %w", scenarioName, err)
		}
	}

	data, err := json.MarshalIndent(
		scenarioRules,
		"",
		"  ",
	)
	if err != nil {
		return nil, fmt.Errorf("encode scenario manifest: %w", err)
	}

	err = os.WriteFile(
		filepath.Join(outputDirectory, "scenarios.json"),
		data,
		0644,
	)
	if err != nil {
		return nil, fmt.Errorf("write scenario manifest: %w", err)
	}

	return scenarioRules, nil
}
