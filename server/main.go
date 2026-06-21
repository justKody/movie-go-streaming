package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/justKody/movie-streaming-go/server/controllers"
)

func main() {

	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.String(200, "Hello world")
	})

	router.GET("/movies", controllers.GetMovies())
	router.GET("/movie/:imdb_id", controllers.GetMovie())
	router.POST("/movie", controllers.AddMovie())

	if err := router.Run(":8080"); err !=nil {
		fmt.Println("Failed to start server", err)
	}
}