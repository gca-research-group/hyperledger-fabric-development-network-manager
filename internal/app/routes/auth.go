package routes

import (
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/handlers"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/repositories"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupAuthRoutes(engine *gin.RouterGroup, db *gorm.DB) *gin.RouterGroup {

	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	service := services.NewAuthService(userService)
	handler := handlers.NewAuthHandler(service)

	routes := engine.Group("/auth")
	{
		routes.POST("/login", handler.Login)
		routes.POST("/refresh", handler.Refresh)
	}

	return routes
}
