package main

import (
	"net/http"

	"github.com/gca-research-group/hyperledger-fabric-network-manager/config/database"
	"github.com/gca-research-group/hyperledger-fabric-network-manager/pkg"
	"github.com/gca-research-group/hyperledger-fabric-network-manager/repository"
	"github.com/gin-gonic/gin"
)

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

	database.Migrate()

	server := gin.Default()

	server.GET("/config", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "teste"})
	})

	server.POST("/config", func(ctx *gin.Context) {
		var body pkg.Config
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, repository.Create(body))
	})

	server.Run()
}
