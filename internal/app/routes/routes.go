package routes

import (
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/middlewares"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetUpPublicRoutes(engine *gin.RouterGroup, db *gorm.DB) {
	SetupAuthRoutes(engine, db)
}

func SetUpProtectedRoutes(engine *gin.RouterGroup, db *gorm.DB) {
	engine.Use(middlewares.AuthHandler())

	SetupChannelRoutes(engine, db)
	SetupOrdererRoutes(engine, db)
	SetupPeerRoutes(engine, db)
	SetupUserRoutes(engine, db)
}
