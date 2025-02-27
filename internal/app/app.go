package app

import (
	"time"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/database"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/middlewares"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Run() {

	godotenv.Load(".env")

	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	server.Use(gin.Recovery())
	server.Use(middlewares.ErrorHandler())

	db := database.Connection()
	api := server.Group("/api/v1")
	routes.SetUp(api, db)

	server.Run()
}
