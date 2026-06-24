package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/justKody/movie-streaming-go/server/controllers"
)

func SetupUnprotectedRoutes(router *gin.Engine) {

	router.GET("/movies", controllers.GetMovies())
	router.POST("/signup", controllers.SignUp())
	router.POST("/login", controllers.Login())

}
