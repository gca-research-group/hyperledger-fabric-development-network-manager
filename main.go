package main

import "github.com/gca-research-group/hyperledger-fabric-development-network-manager/api"

func main() {
	// config := pkg.Config{
	// 	Orderers: []pkg.Orderer{
	// 		{
	// 			Name:   "Orderer",
	// 			Domain: "example.com",
	// 			Port:   7050,
	// 		},
	// 	},
	// 	Peers: []pkg.Peer{
	// 		{
	// 			Name:   "Org1",
	// 			Domain: "org1.example.com",
	// 			Port:   7051,
	// 		},
	// 		{
	// 			Name:   "Org2",
	// 			Domain: "org2.example.com",
	// 			Port:   7051,
	// 		},
	// 		{
	// 			Name:   "Org3",
	// 			Domain: "org3.example.com",
	// 		},
	// 	},
	// 	Networks: []pkg.Network{
	// 		{Name: "MultiChannel1", Organizations: []string{"Org1", "Org2", "Org3"}},
	// 	},
	// }

	// _cryptoconfig := yaml.Parse(cryptoconfig.Build(config))
	// yaml.Write(_cryptoconfig, "./crypto-config.yaml")

	// _configtx := configtx.UpdateAnchors(yaml.Parse(configtx.Build(config)), config.Networks[0].Organizations)
	// yaml.Write(_configtx, "./configtx.yaml")
	api.Run()
}
