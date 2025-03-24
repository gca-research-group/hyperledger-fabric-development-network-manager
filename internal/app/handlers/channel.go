package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/dtos"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/errors"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/models/sql"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/services"
	"github.com/gin-gonic/gin"
)

type ChannelHandler struct {
	service *services.ChannelService
}

func NewChannelHandler(service *services.ChannelService) *ChannelHandler {
	return &ChannelHandler{service: service}
}

func (h *ChannelHandler) Index(c *gin.Context) {
	var queryParams dtos.ChannelDto
	c.ShouldBindQuery(&queryParams)

	queryOptions := sql.QueryOptions{}
	queryOptions.UpdateFromContext(c)

	data, err := h.service.FindAll(queryOptions, queryParams)

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

func (h *ChannelHandler) Show(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	data, err := h.service.FindById(uint(id))

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

func (h *ChannelHandler) Create(c *gin.Context) {
	var data dtos.ChannelDto

	if err := c.ShouldBindJSON(&data); err != nil {
		slog.Error("[Channel -> Create ->  ShouldBindJSON]", "err", err)
		c.Error(&errors.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	channel := data.ToEntity()
	_, err := h.service.Create(&channel)

	if err != nil {
		c.Error(&errors.AppError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, channel)
}

func (h *ChannelHandler) Update(c *gin.Context) {
	var data dtos.ChannelDto

	if err := c.ShouldBindJSON(&data); err != nil {
		slog.Error("[Channel -> Update ->  ShouldBindJSON]", "err", err)
		c.Error(&errors.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
		})
		return
	}

	channel := data.ToEntity()
	updatedChannel, err := h.service.Update(&channel)

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

func (h *ChannelHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.service.Delete(uint(id)); err != nil {
		slog.Error("[Channel -> Delete]", "err", err)
		c.Error(&errors.AppError{
			Code:    http.StatusBadRequest,
			Message: "FAILED_TO_DELETE",
		})
		return
	}

	c.Status(http.StatusOK)
}
