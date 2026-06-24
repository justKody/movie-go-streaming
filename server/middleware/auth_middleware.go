package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/justKody/movie-streaming-go/server/utils"
	"github.com/justKody/movie-streaming-go/server/utils/response"
)

func AuthMiddlWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := utils.GetAccessToken(c)

		if err != nil {
			response.ErrorResponse(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		if token == "" {
			response.ErrorResponse(c, http.StatusUnauthorized, "No token provided")
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(token)

		if err != nil {
			response.ErrorResponse(c, http.StatusUnauthorized, "Invalid Token")
			c.Abort()
			return
		}

		c.Set("userId", claims.UserId) // like req.userId
		c.Set("role", claims.Role)

		c.Next()
	}
}
