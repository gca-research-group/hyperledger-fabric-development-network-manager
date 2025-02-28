package orderer

import (
	"net/http"
	"strconv"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/errors"
	model "github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/models/orderer"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/models/sql"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Index(c *gin.Context, db *gorm.DB) {
	entity := model.Orderer{}

	var queryParams model.OrdererDto
	c.ShouldBindQuery(&queryParams)

	queryOptions := sql.QueryOptions{}
	queryOptions.UpdateFromContext(c)

	data, err := entity.FindAll(db, queryOptions, queryParams)

	if err != nil {
		c.Error(&errors.AppError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, data)
}

func Show(c *gin.Context, db *gorm.DB) {
	entity := model.Orderer{}
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
	var data model.Orderer

	if err := c.ShouldBindJSON(&data); err != nil {
		c.Error(&errors.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	entity := model.Orderer{}
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
	var data model.Orderer

	if err := c.ShouldBindJSON(&data); err != nil {
		c.Error(&errors.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	entity := model.Orderer{}
	_, err := entity.Update(db, data)

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

	entity := model.Orderer{}
	if err := entity.Delete(db, uint(id)); err != nil {
		c.Error(&errors.AppError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}
