package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/dtos"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/errors"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/services"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var data dtos.LoginDto
	if err := c.ShouldBindJSON(&data); err != nil {
		slog.Error("[Channel -> Create ->  ShouldBindJSON]", "err", err)
		c.Error(&errors.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	response, err := h.service.Login(data.Email, data.Password)

	if err != nil {
		c.Error(&errors.AppError{
			Code:    http.StatusUnauthorized,
			Message: err.Error(),
		})
		return
	}

	utils.SetRefreshTokenAsCookie(c, response.RefreshToken)

	c.JSON(http.StatusOK, map[string]interface{}{
		"accessToken": response.AccessToken,
		"user":        response.User,
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("jrt")

	if err != nil {
		c.Error(&errors.AppError{
			Code:    http.StatusUnauthorized,
			Message: err.Error(),
		})
		return
	}

	response, err := h.service.Refresh(refreshToken)

	if err != nil {
		c.Error(&errors.AppError{
			Code:    http.StatusUnauthorized,
			Message: "SESSION_EXPIRED",
		})
		return
	}

	utils.SetRefreshTokenAsCookie(c, response.RefreshToken)

	c.JSON(http.StatusOK, map[string]interface{}{
		"accessToken": response.AccessToken,
		"user":        response.User,
	})
}
