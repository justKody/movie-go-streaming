package response

import "github.com/gin-gonic/gin"

func SuccessResponse(c *gin.Context, code int, data any, message ...string) {
	response := gin.H{
		"data": data,
	}

	if len(message) > 0 {
		response["message"] = message[0]
	}

	c.JSON(code, response)
}