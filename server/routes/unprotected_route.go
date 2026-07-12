package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/justKody/movie-streaming-go/server/controllers"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SetupUnprotectedRoutes(router *gin.Engine, Client *mongo.Client) {

	router.GET("/movies", controllers.GetMovies(Client))
	router.POST("/signup", controllers.SignUp(Client))
	router.POST("/login", controllers.Login(Client))

}
