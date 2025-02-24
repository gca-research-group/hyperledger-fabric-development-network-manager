package orderer

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(engine *gin.RouterGroup, db *gorm.DB) *gin.RouterGroup {

	routes := engine.Group("/orderer")
	{
		routes.POST("/", func(c *gin.Context) {
			CreateOrUpdate(c, db)
		})
		routes.PUT("/", func(c *gin.Context) {
			CreateOrUpdate(c, db)
		})
		routes.GET("/", func(c *gin.Context) {
			Index(c, db)
		})
		routes.GET("/:id", func(c *gin.Context) {
			Show(c, db)
		})
		routes.DELETE("/:id", func(c *gin.Context) {
			Delete(c, db)
		})
	}

	return routes
}
