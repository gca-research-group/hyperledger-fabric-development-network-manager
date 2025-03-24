package routes

import (
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/handlers"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/repositories"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupOrdererRoutes(engine *gin.RouterGroup, db *gorm.DB) *gin.RouterGroup {

	repository := repositories.NewOrdererRepository(db)
	service := services.NewOrdererService(repository)
	handler := handlers.NewOrdererHandler(service)

	routes := engine.Group("/orderer")
	{
		routes.POST("/", handler.Create)
		routes.PUT("/", handler.Update)
		routes.GET("/", handler.Index)
		routes.GET("/:id", handler.Show)
		routes.DELETE("/:id", handler.Delete)
	}

	return routes
}
