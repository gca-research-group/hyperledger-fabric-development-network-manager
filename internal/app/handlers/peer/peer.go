package peer

import (
	"net/http"
	"strconv"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/errors"
	model "github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/models/peer"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Index(c *gin.Context, db *gorm.DB) {
	entity := model.Peer{}
	data, err := entity.FindAll(db)

	if err != nil {
		c.Error(&errors.AppError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

func Show(c *gin.Context, db *gorm.DB) {
	entity := model.Peer{}
	id, _ := strconv.Atoi(c.Param("id"))
	data, err := entity.FindById(db, uint(id))

	if err != nil {
		c.Error(&errors.AppError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, data)
}

func Create(c *gin.Context, db *gorm.DB) {
	var data model.Peer

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	entity := model.Peer{}
	_, err := entity.Create(db, &data)

	if err != nil {
		c.Error(&errors.AppError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, data)
}

func Update(c *gin.Context, db *gorm.DB) {
	var data model.Peer
	id, _ := strconv.Atoi(c.Param("id"))

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	entity := model.Peer{}
	_, err := entity.Update(db, uint(id), &data)

	if err != nil {
		c.Error(&errors.AppError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, data)
}

func Delete(c *gin.Context, db *gorm.DB) {
	id, _ := strconv.Atoi(c.Param("id"))

	entity := model.Peer{}
	if err := entity.Delete(db, uint(id)); err != nil {
		c.Error(&errors.AppError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}
