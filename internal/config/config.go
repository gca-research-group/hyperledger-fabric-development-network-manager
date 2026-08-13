package config

import "github.com/gca-research-group/fabric-network-orchestrator/internal/spec"

type CapabilityLevel = spec.CapabilityLevel

const (
	V2_0 = spec.V2_0
	V2_5 = spec.V2_5
	V3_0 = spec.V3_0
)

var CapabilityMap = spec.CapabilityMap
var MinBinaryVersion = spec.MinBinaryVersion
var DefaultVersionByCapability = spec.DefaultVersionByCapability

type Orderer = spec.Orderer
type Peer = spec.Peer
type CertificateAuthority = spec.CertificateAuthority
type Organization = spec.Organization
type Channel = spec.Channel
type Consensus = spec.Consensus
type Profile = spec.Profile
type Capabilities = spec.Capabilities
type Language = spec.Language
type Chaincode = spec.Chaincode
type Config = spec.Config
