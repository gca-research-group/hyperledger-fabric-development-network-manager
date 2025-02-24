package app

import (
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/database"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/middlewares"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/routes"
	"github.com/gin-gonic/gin"
)

func Run() {

	server := gin.Default()

	server.Use(gin.Recovery())
	server.Use(middlewares.ErrorHandler())

	db := database.Connection()

	api := server.Group("/api/v1")
	routes.SetupOrdererRoutes(api, db)

	server.Run()
}
