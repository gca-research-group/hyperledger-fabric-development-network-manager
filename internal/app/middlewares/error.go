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

		if c.Writer.Written() {
			return
		}

		if len(c.Errors) > 0 {
			var appErr *custom.AppError
			message := "UNEXPECTED_ERROR"

			if errors.As(c.Errors[0].Err, &appErr) && appErr.Code != http.StatusInternalServerError {
				c.JSON(appErr.Code, custom.AppError{
					Code:    appErr.Code,
					Message: appErr.Message,
				})
			} else {
				c.JSON(http.StatusInternalServerError, custom.AppError{
					Code:    http.StatusInternalServerError,
					Message: message,
				})
			}
		}
	}
}
