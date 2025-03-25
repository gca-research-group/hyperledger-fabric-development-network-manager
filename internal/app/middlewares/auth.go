package middlewares

import (
	"net/http"
	"strings"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/utils"
	"github.com/gin-gonic/gin"
)

func AuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		authorization := c.GetHeader("Authorization")

		parts := strings.Split(authorization, " ")

		if len(parts) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "TOKEN_IS_REQUIRED",
			})
			return
		}

		token := parts[1]

		_, err := utils.VerifyToken(token)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": err,
			})
			return
		}

		c.Next()

	}
}
