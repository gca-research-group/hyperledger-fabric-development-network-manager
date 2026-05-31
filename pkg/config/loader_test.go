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
				Name: "mychannel",
				Profile: Profile{
					Name: "TwoOrgsChannel",
				},
			},
		},
	}
}

func TestValidateSyntactic(t *testing.T) {
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
			expectError: "the output directory cannot be empty",
		},
		{
			name: "No organizations",
			modify: func(c *Config) {
				c.Organizations = []Organization{}
			},
			expectError: "at least one organization must be defined",
		},
		{
			name: "Organization name empty",
			modify: func(c *Config) {
				c.Organizations[0].Name = ""
			},
			expectError: "name of the organization index 0 is undefined",
		},
		{
			name: "Organization domain empty",
			modify: func(c *Config) {
				c.Organizations[0].Domain = ""
			},
			expectError: "domain of the organization index 0 is undefined",
		},
		{
			name: "CA expose port negative",
			modify: func(c *Config) {
				c.Organizations[0].CertificateAuthority.ExposePort = 0
			},
			expectError: "expose port of the certificate authority of the organization Org1 should be greater than zero",
		},
		{
			name: "Peer name empty",
			modify: func(c *Config) {
				c.Organizations[0].Peers[0].Name = ""
			},
			expectError: "name of the peer index 0 of the organization Org1 is undefined",
		},
		{
			name: "Peer subdomain empty",
			modify: func(c *Config) {
				c.Organizations[0].Peers[0].Subdomain = ""
			},
			expectError: "subdomain of the peer peer0 of the organization Org1 is undefined",
		},
		{
			name: "Peer expose port invalid",
			modify: func(c *Config) {
				c.Organizations[0].Peers[0].ExposePort = -10
			},
			expectError: "expose port of the peer peer0 of the organization Org1 should be greater than zero",
		},
		{
			name: "Orderer name empty",
			modify: func(c *Config) {
				c.Organizations[0].Orderers[0].Name = ""
			},
			expectError: "name of the orderer index 0 of the organization Org1 is undefined",
		},
		{
			name: "Orderer subdomain empty",
			modify: func(c *Config) {
				c.Organizations[0].Orderers[0].Subdomain = ""
			},
			expectError: "subdomain of the orderer orderer0 of the organization Org1 is undefined",
		},
		{
			name: "Orderer expose port invalid",
			modify: func(c *Config) {
				c.Organizations[0].Orderers[0].ExposePort = 0
			},
			expectError: "expose port of the orderer orderer0 of the organization Org1 should be greater than zero",
		},
		{
			name: "Chaincode name empty",
			modify: func(c *Config) {
				c.Chaincodes = []Chaincode{
					{Path: "some/path", Version: "1.0"},
				}
			},
			expectError: "name of the chaincode 0 is empty",
		},
		{
			name: "Chaincode path empty",
			modify: func(c *Config) {
				c.Chaincodes = []Chaincode{
					{Name: "mycc", Version: "1.0"},
				}
			},
			expectError: "path of the chaincode 0 is empty",
		},
		{
			name: "Chaincode version empty",
			modify: func(c *Config) {
				c.Chaincodes = []Chaincode{
					{Name: "mycc", Path: "some/path"},
				}
			},
			expectError: "path of the chaincode 0 is empty",
		},
		{
			name: "Profile organizations empty",
			modify: func(c *Config) {
				c.Profiles[0].Organizations = []string{}
			},
			expectError: "profile TwoOrgsChannel must include at least one organization",
		},
		{
			name: "Channel profile empty",
			modify: func(c *Config) {
				c.Channels[0].Profile.Name = ""
			},
			expectError: "channel mychannel must reference a profile",
		},
		{
			name: "Channel name empty",
			modify: func(c *Config) {
				c.Channels[0].Name = ""
			},
			expectError: "channel name cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := getValidBaseConfig()
			tt.modify(&conf)
			err := ValidateSyntactic(conf)
			if tt.expectError == "" {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.expectError)
				} else if err.Error() != tt.expectError {
					t.Errorf("expected error %q, got %q", tt.expectError, err.Error())
				}
			}
		})
	}
}

func TestValidateSemantic(t *testing.T) {
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
			name: "Unsupported channel capability",
			modify: func(c *Config) {
				c.Capabilities.Channel = "V1_4"
			},
			expectError: "unsupported channel capability: V1_4",
		},
		{
			name: "Unsupported application capability",
			modify: func(c *Config) {
				c.Capabilities.Application = "V1_4"
			},
			expectError: "unsupported application capability: V1_4",
		},
		{
			name: "Unsupported orderer capability",
			modify: func(c *Config) {
				c.Capabilities.Orderer = "V1_4"
			},
			expectError: "unsupported orderer capability: V1_4",
		},
		{
			name: "Duplicate organization name",
			modify: func(c *Config) {
				c.Organizations = append(c.Organizations, Organization{
					Name:   "Org1",
					Domain: "org1dup.example.com",
					CertificateAuthority: CertificateAuthority{
						ExposePort: 8054,
					},
				})
			},
			expectError: "duplicate organization name: Org1",
		},
		{
			name: "Invalid peer version (lower than required)",
			modify: func(c *Config) {
				c.Capabilities.Channel = "V2_5" // min version is 2.5.0
				c.Organizations[0].Peers[0].Version = "2.4.0"
			},
			expectError: "peer version of org Org1 invalid: version 2.4.0 is lower than required 2.5.0",
		},
		{
			name: "Invalid orderer version (lower than required)",
			modify: func(c *Config) {
				c.Capabilities.Channel = "V2_5" // min version is 2.5.0
				c.Organizations[0].Orderers[0].Version = "2.4.0"
			},
			expectError: "orderer version of org Org1 invalid: version 2.4.0 is lower than required 2.5.0",
		},
		{
			name: "No orderers configured",
			modify: func(c *Config) {
				c.Organizations[0].Orderers = []Orderer{}
			},
			expectError: "at least one orderer must be configured",
		},
		{
			name: "Multiple bootstrap organizations",
			modify: func(c *Config) {
				c.Organizations = append(c.Organizations, Organization{
					Name:   "Org2",
					Domain: "org2.example.com",
					CertificateAuthority: CertificateAuthority{
						ExposePort: 8054,
					},
					Orderers: []Orderer{
						{Name: "orderer1", Subdomain: "orderer1", ExposePort: 8050},
					},
					Bootstrap: true,
				})
			},
			expectError: "exactly one bootstrap organization must be defined",
		},
		{
			name: "Invalid consensus type in profile",
			modify: func(c *Config) {
				c.Profiles[0].Consensus.Type = "solo"
			},
			expectError: "invalid consensus type for the profile TwoOrgsChannel",
		},
		{
			name: "Profile references undefined organization",
			modify: func(c *Config) {
				c.Profiles[0].Organizations = []string{"UndefinedOrg"}
			},
			expectError: "organization not defined: UndefinedOrg",
		},
		{
			name: "Port conflict between CA and Peer",
			modify: func(c *Config) {
				c.Organizations[0].CertificateAuthority.ExposePort = 7051 // peer expose port is 7051
			},
			expectError: "port conflict: expose port 7051 is configured for Certificate Authority of org Org1 and peer peer0 of org Org1",
		},
		{
			name: "Port conflict between Peer and Orderer",
			modify: func(c *Config) {
				c.Organizations[0].Orderers[0].ExposePort = 7051 // peer expose port is 7051
			},
			expectError: "port conflict: expose port 7051 is configured for orderer orderer0 of org Org1 and peer peer0 of org Org1",
		},
		{
			name: "Invalid channel name format (uppercase)",
			modify: func(c *Config) {
				c.Channels[0].Name = "MYCHANNEL"
			},
			expectError: "invalid channel name: MYCHANNEL (must be lowercase alphanumeric, start with a letter, and contain only '.', '-', or alphanumeric characters, max 249 characters)",
		},
		{
			name: "Invalid channel name format (spaces)",
			modify: func(c *Config) {
				c.Channels[0].Name = "my channel"
			},
			expectError: "invalid channel name: my channel (must be lowercase alphanumeric, start with a letter, and contain only '.', '-', or alphanumeric characters, max 249 characters)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := getValidBaseConfig()
			tt.modify(&conf)
			err := ValidateSemantic(conf)
			if tt.expectError == "" {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.expectError)
				} else if err.Error() != tt.expectError {
					t.Errorf("expected error %q, got %q", tt.expectError, err.Error())
				}
			}
		})
	}
}

func TestValidateConfig(t *testing.T) {
	t.Run("syntactic failure", func(t *testing.T) {
		conf := getValidBaseConfig()
		conf.Output = ""
		err := validateConfig(conf)
		if err == nil || err.Error() != "the output directory cannot be empty" {
			t.Errorf("expected syntactic error, got: %v", err)
		}
	})

	t.Run("semantic failure", func(t *testing.T) {
		conf := getValidBaseConfig()
		conf.Capabilities.Channel = "V1_4"
		err := validateConfig(conf)
		if err == nil || err.Error() != "unsupported channel capability: V1_4" {
			t.Errorf("expected semantic error, got: %v", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		conf := getValidBaseConfig()
		err := validateConfig(conf)
		if err != nil {
			t.Errorf("expected success, got error: %v", err)
		}
	})
}
