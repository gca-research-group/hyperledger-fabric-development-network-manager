<h1 align="center">
  <br>
    <img src="assets/logo.svg" height="256px" alt="Fabric Network Orchestrator">
  <br>
  Fabric Network Orchestrator
  <br>
</h1>

<p align="center">
    <img alt="Docker" src="https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white" />
    <img alt="Hyperledger Fabric" src="https://img.shields.io/badge/Hyperledger_Fabric-2.0-ff69b4?style=for-the-badge&logo=hyperledger&logoColor=white" />
    <img alt="Go" src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" />
</p>

<div align="center">

🚧 **This project is currently under development.** 🚧  
Expect frequent updates and changes. Your feedback is appreciated!

</div>

## Overview

The **Fabric Network Orchestrator** is a CLI-based orchestration tool designed to simplify the configuration, generation, and deployment of local Hyperledger Fabric networks.

It automates:

- Docker Compose generation  
- `configtx.yml` generation  
- MSP identity generation  
- Certificate authority setup  
- Network bootstrapping  
- Chaincode lifecycle management  

The tool is intended for:

- Researchers  
- Developers  
- Academic experimentation  
- Rapid prototyping of Fabric networks 

## Table of contents

- [Getting Started](#getting-started)
  - [Executing with Binaries](#executing-with-binaries)
  - [Executing From Source Code](#executing-from-source-code)
- [Configuration File](#configuration-file)
- [Network Samples](#network-samples)
- [Chaincode Samples](#chaincode-samples)
- [Project Repositories](#project-repositories)
- [Related Publications](#related-publications)
- [License](#license)

## Getting Started

### Executing with Binaries

Compiled binaries are available in the [./.bin](./.bin) directory for Windows, macOS, and Linux. Choose the one that matches your OS:

- **Windows:** `fno_windows_amd64.exe`
- **Linux:** `fno_linux_arm64`
- **macOS:** `fno_darwin_arm64`

> **Note:** Replace `<binary>` in the commands below with your specific file name (e.g., `./fno_linux_arm64`). Use a configuration file from the [samples](./samples/) folder.

#### Commands Workflow

1. **Generate Artifacts** Generates Docker Compose and `configtx` files.

```bash
  <binary> artifacts generate --config=samples/minimal-network.yml
```

2. **Deploy Network** Starts CAs, generates identities, and initializes orderers/peers.

```bash
  <binary> deploy --config=samples/minimal-network.yml
```

3. **Start/Stop Containers** Manage the runtime state of the network.

```bash
  # Start network
  <binary> network up --config=samples/minimal-network.yml

  # Stop containers
  <binary> network down --config=samples/minimal-network.yml
```

4. **Clean Everything** Removes all generated files, ledger state, and identities. Use with caution.

```bash
  <binary> artifacts clean --config=samples/minimal-network.yml
```

### Executing From Source Code

#### Prerequisites

- **Docker** & **Docker Compose**
- **Go** (version 1.26 or higher)

#### Clone the repository

```bash
git clone https://github.com/gca-research-group/hyperledger-fabric-development-network-manager
```

#### Install the dependencies

```bash
go mod tidy
```

#### Commands Workflow

Run the same commands from the "Executing with Binaries" section by replacing the <binary> placeholder with `go run cmd/cli/main.go` (for development/testing).

## Configuration File

```yml
output: output/minimal-network
network: minimal-network
chaincodes:
  - samples/chaincodes

capabilities:
  channel: V2_0
  orderer: V2_0
  application: V2_5

organizations:
  - name: Org1
    domain: org1.minimal-network.com
    orderers:
      - name: Orderer
        subdomain: orderer
    bootstrap: true
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
  - &DefaultProfile
    name: DefaultProfile
    organizations:
      - Org1
      - Org2

channels:
  - name: defaultchannel
    profile: *DefaultProfile
```

## Network Samples

You will find network samples in the directory [samples](./samples/)


## Chaincode Samples

You will find chaincode samples in the directory [samples](./samples/chaincodes)

## Project repositories

- [Smart Contract Execution Monitoring System](https://github.com/gca-research-group/smart-contract-execution-monitoring-system)
- [Fabric Network Orchestrator](https://github.com/gca-research-group/hyperledger-fabric-development-network-manager)
- [Transformation Engine](https://github.com/gca-research-group/jabuti-ce-transformation-engine)
- [Jabuti CE (VSCode Plug-in)](https://github.com/gca-research-group/jabuti-ce-vscode-plugin)
- [Jabuti DSL Grammar](https://github.com/gca-research-group/jabuti-ce-jabuti-dsl-grammar)
- [Jabuti XText/Xtend implementation](https://github.com/gca-research-group/dsl-smart-contract-eai)

## Related Publications

- 2025
  - [Proposing a Tool to Monitor Smart Contract Execution in Integration Processes](https://sol.sbc.org.br/index.php/sbsi_estendido/article/view/34617)
  - [Towards a Smart Contract Toolkit for Application Integration](#)
 
- 2024
  - [Jabuti CE: A Tool for Specifying Smart Contracts in the Domain of Enterprise Application Integration](https://www.scitepress.org/Link.aspx?doi=10.5220/0012413300003645)

- 2022
  - [Advances in a DSL to Specify Smart Contracts for Application Integration Processes](https://sol.sbc.org.br/index.php/cibse/article/view/20962)
  - [On the Need to Use Smart Contracts in Enterprise Application Integration](https://idus.us.es/handle/11441/140199)

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contact

For any questions or issues, please open an issue on GitHub or contact the maintainers.
