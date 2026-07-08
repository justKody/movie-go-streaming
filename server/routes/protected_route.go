package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/justKody/movie-streaming-go/server/controllers"
	"github.com/justKody/movie-streaming-go/server/middleware"
)

func SetupProtectedRoutes(router *gin.Engine) {
	router.Use(middleware.AuthMiddlWare())

	router.GET("/movie/:imdb_id", controllers.GetMovie())
	router.POST("/movie", controllers.AddMovie())
	router.POST("/movie/:imdb_id/admin-review", controllers.AdminReviewUpdate())
}
