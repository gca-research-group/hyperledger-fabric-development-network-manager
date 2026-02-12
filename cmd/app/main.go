package main

import (
	"flag"
	"log"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/command"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/fabric"
)

func main() {

	configFlag := flag.String("config", "", "Path to config file")

	flag.Parse()

	if *configFlag == "" {
		log.Fatal("Config is required")
	}

	config, err := pkg.LoadConfigFromPath(*configFlag)

	if err != nil {
		log.Fatal(err)
	}

	fabric, err := fabric.NewFabric(*config, &command.DefaultExecutor{})

	if err != nil {
		log.Fatal(err)
	}

	if err := fabric.DeployNetwork(); err != nil {
		log.Fatalf("Network deployment failed: %v", err)
	}
}
