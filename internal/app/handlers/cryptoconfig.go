package handlers

import (
	"net/http"
	"strconv"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/errors"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/services"
	"github.com/gin-gonic/gin"
)

type CryptoConfigHandler struct {
	service *services.CryptoConfigService
}

func NewCryptoConfigHandler(service *services.CryptoConfigService) *CryptoConfigHandler {
	return &CryptoConfigHandler{service: service}
}

func (h *CryptoConfigHandler) GenerateCryptoConfig(c *gin.Context) {
	channelId, _ := strconv.Atoi(c.Param("channelId"))

	out, err := h.service.GenerateCryptoConfig(channelId)

	if err != nil {
		c.Error(&errors.AppError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	c.Data(200, "application/x-yaml", out)
}
