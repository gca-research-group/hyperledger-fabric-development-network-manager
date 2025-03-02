package channel

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/errors"
	model "github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/models/channel"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/models/sql"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Index(c *gin.Context, db *gorm.DB) {
	entity := model.Channel{}

	var queryParams model.ChannelDto
	c.ShouldBindQuery(&queryParams)

	queryOptions := sql.QueryOptions{}
	queryOptions.UpdateFromContext(c)

	data, err := entity.FindAll(db, queryOptions, queryParams)

	if err != nil {
		slog.Error("[Channel -> Index]", "err", err)
		c.Error(&errors.AppError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, data)
}

func Show(c *gin.Context, db *gorm.DB) {
	entity := model.Channel{}
	id, _ := strconv.Atoi(c.Param("id"))
	data, err := entity.FindById(db, uint(id))

	if err != nil {
		slog.Error("[Channel -> Show]", "err", err)
		c.Error(&errors.AppError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, data)
}

func Create(c *gin.Context, db *gorm.DB) {
	var data model.ChannelDto

	if err := c.ShouldBindJSON(&data); err != nil {
		slog.Error("[Channel -> Create ->  ShouldBindJSON]", "err", err)
		c.Error(&errors.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	entity := model.Channel{}
	channel := data.ToEntity()
	_, err := entity.Create(db, &channel)

	if err != nil {
		c.Error(&errors.AppError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, channel)
}

func Update(c *gin.Context, db *gorm.DB) {
	var data model.ChannelDto

	if err := c.ShouldBindJSON(&data); err != nil {
		slog.Error("[Channel -> Update ->  ShouldBindJSON]", "err", err)
		c.Error(&errors.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
		})
		return
	}

	channel := data.ToEntity()
	entity := model.Channel{}
	updatedChannel, err := entity.Update(db, &channel)

	if err != nil {
		slog.Error("[Channel -> Update]", "err", err)
		c.Error(&errors.AppError{
			Code:    http.StatusBadRequest,
			Message: "FAILED_TO_UPDATE",
		})
		return
	}

	c.JSON(http.StatusOK, updatedChannel)
}

func Delete(c *gin.Context, db *gorm.DB) {
	id, _ := strconv.Atoi(c.Param("id"))

	entity := model.Channel{}
	if err := entity.Delete(db, uint(id)); err != nil {
		slog.Error("[Channel -> Delete]", "err", err)
		c.Error(&errors.AppError{
			Code:    http.StatusBadRequest,
			Message: "FAILED_TO_DELETE",
		})
		return
	}

	c.Status(http.StatusOK)
}
