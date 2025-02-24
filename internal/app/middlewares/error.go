package middlewares

import (
	"errors"
	"net/http"

	custom "github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/errors"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			var appErr *custom.AppError
			if errors.As(c.Errors[0].Err, &appErr) {
				c.JSON(appErr.Code, appErr)
			} else {
				c.JSON(http.StatusInternalServerError, custom.AppError{
					Code:    http.StatusInternalServerError,
					Message: "An unexpected error occurred",
				})
			}
		}
	}
}
