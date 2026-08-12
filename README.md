# Fabric Network Orchestrator (FNO)

## Overview

Fabric Network Orchestrator (FNO) is a command-line tool and supporting domain-specific language (DSL) for specifying, validating, and deploying local Hyperledger Fabric networks.

FNO enables developers and researchers to describe Hyperledger Fabric topologies using a declarative specification rather than manually managing configuration files, certificates, deployment scripts, and chaincode lifecycle operations. Before deployment, specifications are validated against a set of formal validation constraints to detect configuration errors early.

The tool supports:

- Network specification through YAML, JSON, and TOML files
- Configuration validation
- Hyperledger Fabric artefact generation
- Certificate and MSP generation
- Network deployment and management
- Chaincode lifecycle management

FNO is primarily intended for:

- Research and experimentation
- Rapid prototyping
- Educational environments
- Local Hyperledger Fabric deployments

## Key Features

- Declarative network specification
- Validation of network configurations before deployment
- Support for channels, organisations, peers, orderers, and chaincodes
- Automated generation of Hyperledger Fabric artefacts
- Automated network provisioning
- Cross-platform execution
- No dependency on external Hyperledger Fabric CLI scripts

## Installation

### Using Prebuilt Binaries

Download the latest release from:

https://github.com/gca-research-group/fabric-network-orchestrator/releases

### Building from Source

Requirements:

- Docker
- Go 1.26 or later

```bash
git clone https://github.com/gca-research-group/fabric-network-orchestrator
cd fabric-network-orchestrator

go mod tidy
```

## Command Workflow

### Validate a Configuration

```bash
fno configuration validate --config samples/minimal-network.yml
```

### Generate Artefacts

```bash
fno artifacts generate --config samples/minimal-network.yml
```

Generates:

- Cryptographic material
- MSP structures
- TLS certificates
- Channel configuration artefacts
- Docker Compose descriptors

### Deploy a Network

```bash
fno network deploy --config samples/minimal-network.yml
```

This command performs:

1. Configuration validation
2. Identity generation
3. Artefact generation
4. Network provisioning
5. Channel creation and joining
6. Chaincode deployment (when configured)

### Manage Network Lifecycle

```bash
fno network up --config samples/minimal-network.yml
fno network down --config samples/minimal-network.yml
```

### Deploy Chaincodes

```bash
fno chaincode deploy --config samples/minimal-network.yml
```

### Clean Generated Artefacts

```bash
fno artifacts clean --config samples/minimal-network.yml
```

## Configuration Example

```yaml
output: output/minimal-network
network: minimal-network

capabilities:
  channel: V2_0
  orderer: V2_0
  application: V2_5

organizations:
  - name: Org1
    domain: org1.minimal-network.com
    bootstrap: true
    orderers:
      - name: Orderer
        subdomain: orderer
    peers:
      - name: Peer0
        subdomain: peer0
        isAnchor: true

  - name: Org2
    domain: org2.minimal-network.com
    peers:
      - name: Peer0
        subdomain: peer0
        isAnchor: true

profiles:
  - name: DefaultProfile
    organizations:
      - Org1
      - Org2

channels:
  - name: defaultchannel
    profile: DefaultProfile
```

## Validation Constraints

FNO validates specifications before deployment using the following 29 rules. The rule ID is included in each validation error and is also used by the experiment runner to identify the mutation applied to a scenario.

| Area | Rule ID | Constraint |
| --- | --- | --- |
| Output | `output.directory-name.invalid` | The output path must contain valid directory names. |
| Organizations | `organizations.required` | At least one organization must be defined. |
| Organizations | `organization.name.required` | Every organization must have a name. |
| Organizations | `organization.domain.required` | Every organization must have a domain. |
| Organizations | `organization.name.duplicate` | Organization names must be unique. |
| Organizations | `organization.bootstrap.multiple` | At most one organization may be marked as bootstrap. |
| Certificate authorities | `certificate-authority.port.invalid` | CA exposed ports cannot be negative. |
| Peers | `peer.name.required` | Every peer must have a name. |
| Peers | `peer.subdomain.required` | Every peer must have a subdomain. |
| Peers | `peer.port.invalid` | Peer exposed ports cannot be negative. |
| Peers | `peer.version.invalid` | A configured peer version must satisfy the channel capability's minimum binary version. |
| Orderers | `orderer.name.required` | Every orderer must have a name. |
| Orderers | `orderer.subdomain.required` | Every orderer must have a subdomain. |
| Orderers | `orderer.port.invalid` | Orderer exposed ports cannot be negative. |
| Orderers | `orderer.version.invalid` | A configured orderer version must satisfy the channel capability's minimum binary version. |
| Orderers | `orderer.topology.required` | The topology must contain at least one orderer. |
| Chaincodes | `chaincode.name.required` | Every chaincode must have a name. |
| Chaincodes | `chaincode.path.required` | Every chaincode must have a source path. |
| Chaincodes | `chaincode.version.required` | Every chaincode must have a version. |
| Profiles | `profile.organizations.required` | Every profile must reference at least one organization. |
| Profiles | `profile.consensus-type.invalid` | Consensus type must be empty, `etcdraft`, or `BFT`. |
| Profiles | `profile.organization.undefined` | Every organization referenced by a profile must be defined. |
| Channels | `channel.name.required` | Every channel must have a name. |
| Channels | `channel.name.invalid` | Channel names must follow Hyperledger Fabric naming restrictions. |
| Channels | `channel.profile.required` | Every channel must reference a profile. |
| Capabilities | `capability.channel.unsupported` | Channel capability must be `V2_0`, `V2_5`, or `V3_0`. |
| Capabilities | `capability.application.unsupported` | Application capability must be `V2_0`, `V2_5`, or `V3_0`. |
| Capabilities | `capability.orderer.unsupported` | Orderer capability must be `V2_0`, `V2_5`, or `V3_0`. |
| Networking | `exposed-port.conflict` | Positive exposed ports must be unique across CAs, peers, and orderers. |

Invalid configurations are rejected before artefacts are generated or infrastructure is provisioned.

## Experiment Runner

The [experiment runner](./experiment-runner/) contains two reusable packages:

- `generator` contains mutation operators for all 29 validation rules. Each operator starts from a valid seed configuration, introduces a targeted fault, and records the expected rule ID. Scenarios combine three different rule groups, with at most one mutation selected from each group.
- `runner` loads the generated manifest, validates every scenario configuration, and checks that the reported validation rule is one of the mutations expected for that scenario.

Run its verification tests with:

```bash
go test ./experiment-runner/...
```

Generate the scenario corpus with:

```bash
go run ./experiment-runner
```

The command generates and runs the complete scenario corpus. It writes mutated configurations to `output/config/`, the scenario-to-rule mapping to `output/scenarios.json`, and execution results to `output/results.json`. It exits with a non-zero status if any scenario does not report an expected rule.

## Samples

Network and chaincode samples are available in:

- [samples](./samples/)
- [chaincodes](./samples/chaincodes)

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contact

For any questions or issues, please open an issue on GitHub or contact the maintainers.
