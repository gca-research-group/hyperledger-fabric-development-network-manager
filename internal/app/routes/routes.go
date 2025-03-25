package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetUp(engine *gin.RouterGroup, db *gorm.DB) {
	SetupAuthRoutes(engine, db)
	SetupChannelRoutes(engine, db)
	SetupOrdererRoutes(engine, db)
	SetupPeerRoutes(engine, db)
	SetupUserRoutes(engine, db)
}
