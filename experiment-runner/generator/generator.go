package generator

import (
	"bufio"
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

type Summary struct {
	Total int
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
	networkNameInvalidOperators,
	organizationDomainDuplicateOperators,
	domainInvalidOperators,
	organizationUsersInvalidOperators,
	peerNameDuplicateOperators,
	peerSubdomainDuplicateOperators,
	peerInternalPortInvalidOperators,
	ordererNameDuplicateOperators,
	ordererInternalPortInvalidOperators,
	profileNameDuplicateOperators,
	channelProfileUndefinedOperators,
	channelNameDuplicateOperators,
	chaincodeNameDuplicateOperators,
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
	profileNameRequiredOperators,
	networkNameRequiredOperators,
	organizationsRequiredOperators,
}

var incompatibilities = IncompatibilityPolicy{
	validate.RuleOrganizationsRequired: {
		validate.RuleOrganizationDomainDuplicate,
		validate.RuleDomainInvalid,
		validate.RuleOrganizationUsersInvalid,
		validate.RuleOrganizationNameRequired,
		validate.RuleOrganizationDomainRequired,
		validate.RuleCertificateAuthorityPortInvalid,
		validate.RulePeerNameRequired,
		validate.RulePeerSubdomainRequired,
		validate.RulePeerPortInvalid,
		validate.RuleOrdererNameRequired,
		validate.RuleOrdererSubdomainRequired,
		validate.RuleOrdererPortInvalid,
		validate.RuleOrganizationNameDuplicate,
		validate.RulePeerVersionInvalid,
		validate.RuleOrdererVersionInvalid,
		validate.RuleBootstrapOrganizationsMultiple,
		validate.RuleExposedPortConflict,
		validate.RulePeerNameDuplicate,
		validate.RulePeerSubdomainDuplicate,
		validate.RulePeerInternalPortInvalid,
		validate.RuleOrdererNameDuplicate,
		validate.RuleOrdererInternalPortInvalid,
	},
	validate.RuleOrdererTopologyRequired: {
		validate.RuleOrdererNameRequired,
		validate.RuleOrdererNameDuplicate,
		validate.RuleOrdererSubdomainRequired,
		validate.RuleOrdererPortInvalid,
		validate.RuleOrdererInternalPortInvalid,
		validate.RuleOrdererVersionInvalid,
	},
	validate.RuleChannelCapabilityUnsupported: {
		validate.RulePeerVersionInvalid,
		validate.RuleOrdererVersionInvalid,
	},
	validate.RuleProfileOrganizationsRequired: {
		validate.RuleProfileOrganizationUndefined,
	},
	validate.RulePeerPortInvalid: {
		validate.RuleExposedPortConflict,
	},
	validate.RuleOrganizationNameRequired: {
		validate.RuleOrganizationNameDuplicate,
	},
	validate.RuleProfileNameRequired: {
		validate.RuleProfileNameDuplicate,
		validate.RuleChannelProfileUndefined,
	},
	validate.RuleChannelNameRequired: {
		validate.RuleChannelNameDuplicate,
		validate.RuleChannelNameInvalid,
	},
	validate.RuleChannelProfileUndefined: {
		validate.RuleChannelProfileRequired,
	},
	validate.RuleNetworkNameRequired: {
		validate.RuleNetworkNameInvalid,
	},
	validate.RuleOrganizationDomainRequired: {
		validate.RuleDomainInvalid,
		validate.RuleOrganizationDomainDuplicate,
	},
	validate.RuleOrganizationDomainDuplicate: {
		validate.RuleDomainInvalid,
	},
	validate.RulePeerNameDuplicate: {
		validate.RulePeerNameRequired,
	},
	validate.RulePeerSubdomainDuplicate: {
		validate.RulePeerSubdomainRequired,
	},
	validate.RuleOrdererNameDuplicate: {
		validate.RuleOrdererNameRequired,
	},
	validate.RuleChaincodeNameDuplicate: {
		validate.RuleChaincodeNameRequired,
	},
}

const DefaultMutationCount = 3

func Generate(outputDirectory string, mutationCount int) (Summary, error) {
	seed, err := yaml.FromBytes([]byte(seedYAML))
	if err != nil {
		return Summary{}, fmt.Errorf("parse seed configuration: %w", err)
	}

	configDirectory := filepath.Join(outputDirectory, "config")

	if err := os.MkdirAll(configDirectory, 0755); err != nil {
		return Summary{}, fmt.Errorf("create config directory: %w", err)
	}

	manifest, err := os.Create(filepath.Join(outputDirectory, "scenarios.json"))

	if err != nil {
		return Summary{}, fmt.Errorf("create scenario manifest: %w", err)
	}

	writer := bufio.NewWriter(manifest)

	if _, err := writer.WriteString("[\n"); err != nil {
		manifest.Close()
		return Summary{}, fmt.Errorf("write scenario manifest: %w", err)
	}

	summary := Summary{}

	err = WalkCombinations(operators, mutationCount, incompatibilities, func(scenario []MutationOperator) error {
		clone := seed.Clone()
		doc := clone.Document()
		rules := make([]validate.RuleID, 0, len(scenario))

		for _, operator := range scenario {
			rules = append(rules, operator.RuleID)
			operator.Apply(doc)
		}

		scenarioName := fmt.Sprintf("%06d", summary.Total+1)

		if err := clone.ToFile(filepath.Join(configDirectory, scenarioName+".yaml")); err != nil {
			return fmt.Errorf("write scenario %s: %w", scenarioName, err)
		}

		entry, err := json.Marshal(ScenarioRules{Scenario: scenarioName, Rules: rules})

		if err != nil {
			return fmt.Errorf("encode scenario %s: %w", scenarioName, err)
		}

		if summary.Total > 0 {
			if _, err := writer.WriteString(",\n"); err != nil {
				return fmt.Errorf("write scenario manifest: %w", err)
			}
		}

		if _, err := writer.WriteString("  " + string(entry)); err != nil {
			return fmt.Errorf("write scenario manifest: %w", err)
		}

		summary.Total++

		return nil
	})

	if err != nil {
		manifest.Close()
		return Summary{}, err
	}

	if _, err := writer.WriteString("\n]\n"); err != nil {
		manifest.Close()
		return Summary{}, fmt.Errorf("write scenario manifest: %w", err)
	}

	if err := writer.Flush(); err != nil {
		manifest.Close()
		return Summary{}, fmt.Errorf("flush scenario manifest: %w", err)
	}

	if err := manifest.Close(); err != nil {
		return Summary{}, fmt.Errorf("close scenario manifest: %w", err)
	}

	return summary, nil
}
