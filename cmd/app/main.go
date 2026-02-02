package main

import (
	"log"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/command"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/fabric"
)

func main() {

	config, err := pkg.LoadConfigFromPath("./cmd/app/config.yml")

	fabric, err := fabric.NewFabric(config, &command.DefaultExecutor{})
	if err != nil {
		log.Fatal(err)
	}

	if err := fabric.DeployNetwork(); err != nil {
		log.Fatalf("Network deployment failed: %v", err)
	}
}
