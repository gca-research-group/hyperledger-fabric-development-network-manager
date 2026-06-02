package config

import (
	"testing"
)

func getValidBaseConfig() Config {
	return Config{
		Output: "./output",
		Capabilities: Capabilities{
			Channel:     "V2_5",
			Application: "V2_5",
			Orderer:     "V2_5",
		},
		Organizations: []Organization{
			{
				Name:   "Org1",
				Domain: "org1.example.com",
				CertificateAuthority: CertificateAuthority{
					ExposePort: 7054,
					Version:    "latest",
				},
				Peers: []Peer{
					{
						Name:       "peer0",
						Subdomain:  "peer0",
						Port:       7051,
						ExposePort: 7051,
						Version:    "2.5.15",
						IsAnchor:   true,
					},
				},
				Orderers: []Orderer{
					{
						Name:       "orderer0",
						Subdomain:  "orderer0",
						Port:       7050,
						ExposePort: 7050,
						Version:    "2.5.15",
					},
				},
				Bootstrap: true,
			},
		},
		Profiles: []Profile{
			{
				Name:          "TwoOrgsChannel",
				Organizations: []string{"Org1"},
				Consensus: Consensus{
					Type: "etcdraft",
				},
			},
		},
		Channels: []Channel{
			{
				Name:    "mychannel",
				Profile: "TwoOrgsChannel",
			},
		},
	}
}

func formatError(ruleID, rule, detail string) string {
	return (&ValidationError{
		RuleID: ruleID,
		Rule:   rule,
		Detail: detail,
	}).Error()
}

func TestInvalidConfigurations(t *testing.T) {
	tests := []struct {
		name        string
		modify      func(*Config)
		expectError string
	}{
		{
			name:        "Valid base config",
			modify:      func(c *Config) {},
			expectError: "",
		},
		{
			name: "Empty output",
			modify: func(c *Config) {
				c.Output = ""
			},
			expectError: formatError("R01", "Empty Output", "the output directory cannot be empty"),
		},
		{
			name: "No organizations",
			modify: func(c *Config) {
				c.Organizations = []Organization{}
			},
			expectError: formatError("R02", "No Organizations", "at least one organization must be defined"),
		},
		{
			name: "Organization name empty",
			modify: func(c *Config) {
				c.Organizations[0].Name = ""
			},
			expectError: formatError("R03", "Empty Org Name", "name of the organization index 0 is undefined"),
		},
		{
			name: "Organization domain empty",
			modify: func(c *Config) {
				c.Organizations[0].Domain = ""
			},
			expectError: formatError("R04", "Empty Org Domain", "domain of the organization index 0 is undefined"),
		},
		{
			name: "CA expose port negative",
			modify: func(c *Config) {
				c.Organizations[0].CertificateAuthority.ExposePort = -1
			},
			expectError: formatError("R05", "Invalid CA Port", "expose port of the certificate authority of the organization Org1 should be greater than zero"),
		},
		{
			name: "Peer name empty",
			modify: func(c *Config) {
				c.Organizations[0].Peers[0].Name = ""
			},
			expectError: formatError("R06", "Empty Peer Name", "name of the peer index 0 of the organization Org1 is undefined"),
		},
		{
			name: "Peer subdomain empty",
			modify: func(c *Config) {
				c.Organizations[0].Peers[0].Subdomain = ""
			},
			expectError: formatError("R07", "Empty Peer Subdomain", "subdomain of the peer peer0 of the organization Org1 is undefined"),
		},
		{
			name: "Peer expose port invalid",
			modify: func(c *Config) {
				c.Organizations[0].Peers[0].ExposePort = -10
			},
			expectError: formatError("R08", "Invalid Peer Port", "expose port of the peer peer0 of the organization Org1 should be greater than zero"),
		},
		{
			name: "Orderer name empty",
			modify: func(c *Config) {
				c.Organizations[0].Orderers[0].Name = ""
			},
			expectError: formatError("R09", "Empty Orderer Name", "name of the orderer index 0 of the organization Org1 is undefined"),
		},
		{
			name: "Orderer subdomain empty",
			modify: func(c *Config) {
				c.Organizations[0].Orderers[0].Subdomain = ""
			},
			expectError: formatError("R10", "Empty Orderer Subdomain", "subdomain of the orderer orderer0 of the organization Org1 is undefined"),
		},
		{
			name: "Orderer expose port negative",
			modify: func(c *Config) {
				c.Organizations[0].Orderers[0].ExposePort = -1
			},
			expectError: formatError("R11", "Invalid Orderer Port", "expose port of the orderer orderer0 of the organization Org1 should be greater than zero"),
		},
		{
			name: "Chaincode name empty",
			modify: func(c *Config) {
				c.Channels[0].Chaincodes = []Chaincode{
					{Path: "some/path", Version: "1.0"},
				}
			},
			expectError: formatError("R12", "Empty Chaincode Name", "name of the chaincode 0 is empty"),
		},
		{
			name: "Chaincode path empty",
			modify: func(c *Config) {
				c.Channels[0].Chaincodes = []Chaincode{
					{Name: "mycc", Version: "1.0"},
				}
			},
			expectError: formatError("R13", "Empty Chaincode Path", "path of the chaincode 0 is empty"),
		},
		{
			name: "Chaincode version empty",
			modify: func(c *Config) {
				c.Channels[0].Chaincodes = []Chaincode{
					{Name: "mycc", Path: "some/path"},
				}
			},
			expectError: formatError("R14", "Empty Chaincode Version", "version of the chaincode 0 is empty"),
		},
		{
			name: "Profile organizations empty",
			modify: func(c *Config) {
				c.Profiles[0].Organizations = []string{}
			},
			expectError: formatError("R15", "Empty Profile Orgs", "profile TwoOrgsChannel must include at least one organization"),
		},
		{
			name: "Channel name empty",
			modify: func(c *Config) {
				c.Channels[0].Name = ""
			},
			expectError: formatError("R17", "Empty Channel Name", "channel name cannot be empty"),
		},
		{
			name: "Channel profile empty",
			modify: func(c *Config) {
				c.Channels[0].Profile = ""
			},
			expectError: formatError("R16", "Empty Channel Profile", "channel mychannel must reference a profile"),
		},
		{
			name:        "Valid base config",
			modify:      func(c *Config) {},
			expectError: "",
		},
		{
			name: "Invalid channel capability",
			modify: func(c *Config) {
				c.Capabilities.Channel = "V1_4"
			},
			expectError: formatError("R18", "Unsupported Channel Capability", "unsupported channel capability: V1_4"),
		},
		{
			name: "Invalid application capability",
			modify: func(c *Config) {
				c.Capabilities.Application = "V1_4"
			},
			expectError: formatError("R19", "Unsupported Application Capability", "unsupported application capability: V1_4"),
		},
		{
			name: "Invalid orderer capability",
			modify: func(c *Config) {
				c.Capabilities.Orderer = "V1_4"
			},
			expectError: formatError("R20", "Unsupported Orderer Capability", "unsupported orderer capability: V1_4"),
		},
		{
			name: "Duplicate organization name",
			modify: func(c *Config) {
				c.Organizations = append(c.Organizations, Organization{
					Name:   "Org1",
					Domain: "org1.sample.com",
				})
			},
			expectError: formatError("R21", "Duplicate Org Name", "duplicate organization name: Org1"),
		},
		{
			name: "Peer version below capability",
			modify: func(c *Config) {
				c.Capabilities.Channel = "V2_5"
				c.Organizations[0].Peers[0].Version = "2.4.0"
			},
			expectError: formatError("R22", "Invalid Peer Version", "peer version of org Org1 invalid: version 2.4.0 is lower than required 2.5.0"),
		},
		{
			name: "Orderer version below capability",
			modify: func(c *Config) {
				c.Capabilities.Channel = "V2_5"
				c.Organizations[0].Orderers[0].Version = "2.4.0"
			},
			expectError: formatError("R23", "Invalid Orderer Version", "orderer version of org Org1 invalid: version 2.4.0 is lower than required 2.5.0"),
		},
		{
			name: "No orderers configured",
			modify: func(c *Config) {
				c.Organizations[0].Orderers = []Orderer{}
			},
			expectError: formatError("R24", "No Orderer Topology", "at least one orderer must be configured"),
		},
		{
			name: "Multiple bootstrap organizations",
			modify: func(c *Config) {
				c.Organizations = []Organization{
					{Name: "Org1", Domain: "org1.sample.com", Bootstrap: true, Orderers: []Orderer{{Name: "ord1", Subdomain: "ord1"}}},
					{Name: "Org2", Domain: "org2.sample.com", Bootstrap: true},
				}
			},
			expectError: formatError("R25", "Multiple Bootstrap Orgs", "exactly one bootstrap organization must be defined"),
		},
		{
			name: "Invalid consensus type",
			modify: func(c *Config) {
				c.Profiles[0].Consensus.Type = "raft"
			},
			expectError: formatError("R26", "Invalid Consensus Type", "invalid consensus type for the profile TwoOrgsChannel"),
		},
		{
			name: "Profile references undefined organization",
			modify: func(c *Config) {
				c.Profiles[0].Organizations = []string{"UndefinedOrg"}
			},
			expectError: formatError("R27", "Profile References Undefined Org", "organization not defined: UndefinedOrg"),
		},
		{
			name: "Conflict port CA and peer",
			modify: func(c *Config) {
				c.Organizations[0].CertificateAuthority.ExposePort = 7051
				c.Organizations[0].Peers[0].ExposePort = 7051
			},
			expectError: formatError("R28", "Exposed Port Conflict", "Port 7051 is assigned to both Certificate Authority of org 'Org1' and peer 'peer0'."),
		},
		{
			name: "Conflict port peer and orderer",
			modify: func(c *Config) {
				c.Organizations[0].Peers[0].ExposePort = 7051
				c.Organizations[0].Orderers[0].ExposePort = 7051
			},
			expectError: formatError("R28", "Exposed Port Conflict", "Port 7051 is assigned to both peer 'peer0' and orderer 'orderer0'."),
		},
		{
			name: "Invalid channel name with uppercase",
			modify: func(c *Config) {
				c.Channels[0].Name = "MYCHANNEL"
			},
			expectError: formatError("R29", "Invalid Channel Name", "invalid channel name: MYCHANNEL (must be lowercase alphanumeric, start with a letter, and contain only '.', '-', or alphanumeric characters, max 249 characters)"),
		},
		{
			name: "Invalid channel name with spaces",
			modify: func(c *Config) {
				c.Channels[0].Name = "my channel"
			},
			expectError: formatError("R29", "Invalid Channel Name", "invalid channel name: my channel (must be lowercase alphanumeric, start with a letter, and contain only '.', '-', or alphanumeric characters, max 249 characters)"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := getValidBaseConfig()
			tt.modify(&conf)
			err := ValidateConfig(conf)
			if tt.expectError == "" {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if err.Error() != tt.expectError {
					t.Errorf("expected error:\n%s\n\ngot:\n%s", tt.expectError, err.Error())
				}
			}
		})
	}
}

func TestValidConfigurations(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		conf := getValidBaseConfig()
		err := ValidateConfig(conf)
		if err != nil {
			t.Errorf("expected success, got error: %v", err)
		}
	})
}
