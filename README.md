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

FNO validates specifications before deployment. Validation includes:

- Mandatory attributes
- Valid parameter values
- Version compatibility checks
- Reference consistency checks
- Port conflict detection
- Consensus configuration validation

Invalid configurations are rejected before artefacts are generated or infrastructure is provisioned.

## Samples

Network and chaincode samples are available in:

- [samples](./samples/)
- [chaincodes](./samples/chaincodes)

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contact

For any questions or issues, please open an issue on GitHub or contact the maintainers.