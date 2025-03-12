<h1 align="center">
  <br>
    <img src="assets/logo.svg" height="256px" alt="Hyperledger Fabric Development Network Manager">
  <br>
  Hyperledger Fabric Development Network Manager
  <br>
</h1>

<p align="center">
    <img alt="Docker" src="https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white" />
    <img alt="Hyperledger Fabric" src="https://img.shields.io/badge/Hyperledger_Fabric-2.0-ff69b4?style=for-the-badge&logo=hyperledger&logoColor=white" />
    <img alt="Go" src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" />
    <img alt="Angular" src="https://img.shields.io/badge/Angular-20232f?style=for-the-badge&logo=angular&logoColor=red" />
    <img alt="PostgreSQL" src="https://img.shields.io/badge/PostgreSQL-336791?style=for-the-badge&logo=postgresql&logoColor=white" />
</p>

<div align="center">

🚧 **This project is currently under development.** 🚧  
Expect frequent updates and changes. Your feedback is appreciated!

</div>

## Overview

Setting up a Hyperledger Fabric network can be a challenging task, even for experienced professionals. Therefore, this project aims to offer a straightforward way to configure and run a local Hyperledger Fabric network.

## Table of contents

- [Project papers](#project-papers)
- [Project repositories](#project-repositories)
- [Features](#features)
- [How to execute](#how-to-execute)

## Project papers

- [Advances in a DSL to Specify Smart Contracts for Application Integration Processes](https://sol.sbc.org.br/index.php/cibse/article/view/20962)
- [On the Need to Use Smart Contracts in Enterprise Application Integration](https://idus.us.es/handle/11441/140199)
- [Jabuti CE: A Tool for Specifying Smart Contracts in the Domain of Enterprise Application Integration](https://www.scitepress.org/Link.aspx?doi=10.5220/0012413300003645)

## Project repositories

- [Smart Contract Monitoring System](https://github.com/gca-research-group/smart-contract-execution-monitoring-system)
- [Hyperledger Fabric Network Manager](https://github.com/gca-research-group/hyperledger-fabric-development-network-manager)
- [Transformation Engine](https://github.com/gca-research-group/jabuti-ce-transformation-engine)
- [Jabuti CE (VSCode Plug-in)](https://github.com/gca-research-group/jabuti-ce-vscode-plugin)
- [Jabuti DSL Grammar](https://github.com/gca-research-group/jabuti-ce-jabuti-dsl-grammar)
- [Jabuti XText/Xtend implementation](https://github.com/gca-research-group/dsl-smart-contract-eai)

## Features

- **Set up Orderers**
- **Set up Peers**
- **Set up Channels**
- **Manage the Chaincode Lifecycle**
- **Start/Stop the Network**

## How to execute

> Currently, you can only execute this project by cloning it. However, we are working on developing a Docker image. Therefore, in the coming weeks, you will be able to run it with a single, fast command.

### Prerequisites

- Docker
- NodeJs +22.0
- Golang +1.24
- [Air](https://github.com/air-verse/air)

### Executing

- Clone this repository

```sh
git clone https://github.com/gca-research-group/hyperledger-fabric-development-network-manager.git
```

- Running the api

```sh
./.scripts/app/api/up.sh
```

- Running the frontend

```sh
./.scripts/app/web/up.sh
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contact

For any questions or issues, please open an issue on GitHub or contact the maintainers.
