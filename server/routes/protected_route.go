package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/justKody/movie-streaming-go/server/controllers"
	"github.com/justKody/movie-streaming-go/server/middleware"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SetupProtectedRoutes(router *gin.Engine, Client *mongo.Client) {
	router.Use(middleware.AuthMiddlWare())

	router.GET("/movie/:imdb_id", controllers.GetMovie(Client))
	router.POST("/movie", controllers.AddMovie(Client))
	router.POST("/movie/:imdb_id/admin-review", controllers.AdminReviewUpdate(Client))
	router.GET("/movie/recommendations", controllers.GetRecommendedMovies(Client))
}
