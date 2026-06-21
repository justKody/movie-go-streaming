package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ErrorResponse(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{
		"error": msg,
	})
}

func ValidationErrorResponse(c *gin.Context, err error) {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		details := make([]gin.H, 0, len(validationErrors))
		for _, e := range validationErrors {
			details = append(details, gin.H{
				"field": e.Field(),
				"tag":   e.Tag(),
				"error": e.Error(),
			})
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed.",
			"details": details,
		})
		return
	}

	ErrorResponse(c, http.StatusBadRequest, err.Error())
}