package app

import (
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/database"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/orderer"
	"github.com/gin-gonic/gin"
)

func Run() {

	server := gin.Default()

	db := database.Connection()

	api := server.Group("/api/v1")
	orderer.SetupRoutes(api, db)

	server.Run()
}
