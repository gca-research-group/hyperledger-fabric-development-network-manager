package routes

import (
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/handlers"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/repositories"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupCryptoConfigRoutes(engine *gin.RouterGroup, db *gorm.DB) *gin.RouterGroup {

	channelRepository := repositories.NewChannelRepository(db)
	channelService := services.NewChannelService(channelRepository)
	service := services.NewCryptoConfigService(channelService)
	handler := handlers.NewCryptoConfigHandler(service)

	routes := engine.Group("/cryptoconfig")
	{
		routes.POST("/:channelId", handler.GenerateCryptoConfig)
	}

	return routes
}
