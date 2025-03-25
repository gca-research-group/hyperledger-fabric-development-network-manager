package handlers

import (
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

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Index(c *gin.Context) {
	var queryParams dtos.UserDto
	c.ShouldBindQuery(&queryParams)

	queryOptions := sql.QueryOptions{}
	queryOptions.UpdateFromContext(c)

	data, err := h.service.FindAll(queryOptions, queryParams)

	if err != nil {
		slog.Error("[User -> Index]", "err", err)
		c.Error(&errors.AppError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *UserHandler) Show(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))
	data, err := h.service.FindById(uint(id))

	if err != nil {
		slog.Error("[User -> Show]", "err", err)
		c.Error(&errors.AppError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *UserHandler) Create(c *gin.Context) {
	var data models.User

	if err := c.ShouldBindJSON(&data); err != nil {
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

func (h *UserHandler) Update(c *gin.Context) {
	var data models.User

	if err := c.ShouldBindJSON(&data); err != nil {
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

func (h *UserHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.service.Delete(uint(id)); err != nil {
		slog.Error("[User -> Delete]", "err", err)
		c.Error(&errors.AppError{
			Code:    http.StatusBadRequest,
			Message: "FAILED_TO_DELETE",
		})
		return
	}

	c.Status(http.StatusOK)
}
