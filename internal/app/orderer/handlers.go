package orderer

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Index(c *gin.Context, db *gorm.DB) {
	entity := Orderer{}
	data := entity.FindAll(db)
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func Show(c *gin.Context, db *gorm.DB) {
	entity := Orderer{}
	id, _ := strconv.Atoi(c.Param("id"))
	data := entity.FindById(db, id)
	c.JSON(http.StatusOK, data)
}

func CreateOrUpdate(c *gin.Context, db *gorm.DB) {
	entity := Orderer{}
	var data Orderer

	if err := c.BindJSON(&data); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	entity.CreateOrUpdate(db, &data)

	c.JSON(http.StatusCreated, data)
}

func Delete(c *gin.Context, db *gorm.DB) {
	entity := Orderer{}
	id, _ := strconv.Atoi(c.Param("id"))
	entity.Delete(db, id)
	c.Status(http.StatusOK)
}
