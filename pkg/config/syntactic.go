package config

import "fmt"

func ValidateSyntactic(config Config) error {
	if config.Output == "" {
		return &ValidationError{
			Layer:  "Syntactic Analysis",
			RuleID: "SYN01",
			Rule:   "Empty Output",
			Detail: "the output directory cannot be empty",
		}
	}

	if len(config.Organizations) == 0 {
		return &ValidationError{
			Layer:  "Syntactic Analysis",
			RuleID: "SYN02",
			Rule:   "No Organizations",
			Detail: "at least one organization must be defined",
		}
	}

	for i := range config.Organizations {
		organization := &config.Organizations[i]

		if organization.Name == "" {
			return &ValidationError{
				Layer:  "Syntactic Analysis",
				RuleID: "SYN03",
				Rule:   "Empty Org Name",
				Detail: fmt.Sprintf("name of the organization index %d is undefined", i),
			}
		}

		if organization.Domain == "" {
			return &ValidationError{
				Layer:  "Syntactic Analysis",
				RuleID: "SYN04",
				Rule:   "Empty Org Domain",
				Detail: fmt.Sprintf("domain of the organization index %d is undefined", i),
			}
		}

		if organization.CertificateAuthority.ExposePort < 0 {
			return &ValidationError{
				Layer:  "Syntactic Analysis",
				RuleID: "SYN05",
				Rule:   "Invalid CA Port",
				Detail: fmt.Sprintf("expose port of the certificate authority of the organization %s should be greater than zero", organization.Name),
			}
		}

		for j, peer := range organization.Peers {
			if peer.Name == "" {
				return &ValidationError{
					Layer:  "Syntactic Analysis",
					RuleID: "SYN06",
					Rule:   "Empty Peer Name",
					Detail: fmt.Sprintf("name of the peer index %d of the organization %s is undefined", j, organization.Name),
				}
			}

			if peer.Subdomain == "" {
				return &ValidationError{
					Layer:  "Syntactic Analysis",
					RuleID: "SYN07",
					Rule:   "Empty Peer Subdomain",
					Detail: fmt.Sprintf("subdomain of the peer %s of the organization %s is undefined", peer.Name, organization.Name),
				}
			}

			if peer.ExposePort < 0 {
				return &ValidationError{
					Layer:  "Syntactic Analysis",
					RuleID: "SYN08",
					Rule:   "Invalid Peer Port",
					Detail: fmt.Sprintf("expose port of the peer %s of the organization %s should be greater than zero", peer.Name, organization.Name),
				}
			}
		}

		for j, orderer := range organization.Orderers {
			if orderer.Name == "" {
				return &ValidationError{
					Layer:  "Syntactic Analysis",
					RuleID: "SYN09",
					Rule:   "Empty Orderer Name",
					Detail: fmt.Sprintf("name of the orderer index %d of the organization %s is undefined", j, organization.Name),
				}
			}

			if orderer.Subdomain == "" {
				return &ValidationError{
					Layer:  "Syntactic Analysis",
					RuleID: "SYN10",
					Rule:   "Empty Orderer Subdomain",
					Detail: fmt.Sprintf("subdomain of the orderer %s of the organization %s is undefined", orderer.Name, organization.Name),
				}
			}

			if orderer.ExposePort < 0 {
				return &ValidationError{
					Layer:  "Syntactic Analysis",
					RuleID: "SYN11",
					Rule:   "Invalid Orderer Port",
					Detail: fmt.Sprintf("expose port of the orderer %s of the organization %s should be greater than zero", orderer.Name, organization.Name),
				}
			}
		}
	}

	for _, profile := range config.Profiles {
		if len(profile.Organizations) == 0 {
			return &ValidationError{
				Layer:  "Syntactic Analysis",
				RuleID: "SYN15",
				Rule:   "Empty Profile Orgs",
				Detail: fmt.Sprintf("profile %s must include at least one organization", profile.Name),
			}
		}
	}

	for _, ch := range config.Channels {
		if ch.Name == "" {
			return &ValidationError{
				Layer:  "Syntactic Analysis",
				RuleID: "SYN17",
				Rule:   "Empty Channel Name",
				Detail: "channel name cannot be empty",
			}
		}

		if ch.Profile == "" {
			return &ValidationError{
				Layer:  "Syntactic Analysis",
				RuleID: "SYN16",
				Rule:   "Empty Channel Profile",
				Detail: fmt.Sprintf("channel %s must reference a profile", ch.Name),
			}
		}

		for i := range ch.Chaincodes {
			chaincode := &ch.Chaincodes[i]

			if chaincode.Name == "" {
				return &ValidationError{
					Layer:  "Syntactic Analysis",
					RuleID: "SYN12",
					Rule:   "Empty Chaincode Name",
					Detail: fmt.Sprintf("name of the chaincode %d is empty", i),
				}
			}

			if chaincode.Path == "" {
				return &ValidationError{
					Layer:  "Syntactic Analysis",
					RuleID: "SYN13",
					Rule:   "Empty Chaincode Path",
					Detail: fmt.Sprintf("path of the chaincode %d is empty", i),
				}
			}

			if chaincode.Version == "" {
				return &ValidationError{
					Layer:  "Syntactic Analysis",
					RuleID: "SYN14",
					Rule:   "Empty Chaincode Version",
					Detail: fmt.Sprintf("version of the chaincode %d is empty", i),
				}
			}
		}
	}

	return nil
}
