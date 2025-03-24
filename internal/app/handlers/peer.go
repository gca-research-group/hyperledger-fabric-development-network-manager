package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/dtos"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/errors"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/models"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/models/sql"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/services"
	"github.com/gin-gonic/gin"
)

type PeerHandler struct {
	service *services.PeerService
}

func NewPeerHandler(service *services.PeerService) *PeerHandler {
	return &PeerHandler{service: service}
}

func (h *PeerHandler) Index(c *gin.Context) {

	var queryParams dtos.PeerDto
	c.ShouldBindQuery(&queryParams)

	queryOptions := sql.QueryOptions{}
	queryOptions.UpdateFromContext(c)

	data, err := h.service.FindAll(queryOptions, queryParams)

	if err != nil {
		c.Error(&errors.AppError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *PeerHandler) Show(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	data, err := h.service.FindById(uint(id))

	if err != nil {
		c.Error(&errors.AppError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *PeerHandler) Create(c *gin.Context) {
	var data models.Peer

	if err := c.ShouldBindJSON(&data); err != nil {
		slog.Error(fmt.Sprintf("[Peer -> Create]: %v\n", err))
		c.Error(&errors.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	_, err := h.service.Create(&data)

	if err != nil {
		c.Error(&errors.AppError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, data)
}

func (h *PeerHandler) Update(c *gin.Context) {
	var data models.Peer

	if err := c.ShouldBindJSON(&data); err != nil {
		slog.Error(fmt.Sprintf("[Peer -> Update]: %v\n", err))
		c.Error(&errors.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	_, err := h.service.Update(data)

	if err != nil {
		c.Error(&errors.AppError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *PeerHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.service.Delete(uint(id)); err != nil {
		c.Error(&errors.AppError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}
