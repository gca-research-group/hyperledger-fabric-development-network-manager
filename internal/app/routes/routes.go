package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetUp(engine *gin.RouterGroup, db *gorm.DB) {
	SetupOrdererRoutes(engine, db)
	SetupPeerRoutes(engine, db)
	SetupChannelRoutes(engine, db)
}
