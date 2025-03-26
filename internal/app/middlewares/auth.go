package middlewares

import (
	"net/http"
	"strings"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/errors"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/utils"
	"github.com/gin-gonic/gin"
)

func AuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		authorization := c.GetHeader("Authorization")

		parts := strings.Split(authorization, " ")

		if len(parts) != 2 {
			c.Error(&errors.AppError{
				Code:    http.StatusUnauthorized,
				Message: "INVALID_TOKEN",
			})
			c.Abort()
			return
		}

		token := parts[1]

		_, err := utils.VerifyToken(token)

		if err != nil {
			c.Error(&errors.AppError{
				Code:    http.StatusUnauthorized,
				Message: "TOKEN_EXPIRED",
			})
			c.Abort()
			return
		}

		c.Next()

	}
}
