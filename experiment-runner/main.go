package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

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

func main() {
	seed, err := yaml.FromBytes([]byte(seedYAML))

	if err != nil {
		log.Fatalf("error %v", err)
	}

	quantityOfMutations := 3

	scenarios := GenerateCombinations(operators, quantityOfMutations)

	scenarioRules := make([]ScenarioRules, 0)

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
		err = clone.ToFile(fmt.Sprintf("output/config/%s.yaml", scenarioName))
	}

	data, err := json.MarshalIndent(
		scenarioRules,
		"",
		"  ",
	)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(
		"output/scenarios.json",
		data,
		0644,
	)
	if err != nil {
		log.Fatal(err)
	}
}
