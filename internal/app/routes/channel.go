package routes

import (
	handler "github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/handlers/channel"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupChannelRoutes(engine *gin.RouterGroup, db *gorm.DB) *gin.RouterGroup {

	routes := engine.Group("/channel")
	{
		routes.POST("/", func(c *gin.Context) {
			handler.Create(c, db)
		})
		routes.PUT("/", func(c *gin.Context) {
			handler.Update(c, db)
		})
		routes.GET("/", func(c *gin.Context) {
			handler.Index(c, db)
		})
		routes.GET("/:id", func(c *gin.Context) {
			handler.Show(c, db)
		})
		routes.DELETE("/:id", func(c *gin.Context) {
			handler.Delete(c, db)
		})
	}

	return routes
}
