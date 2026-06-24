package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/justKody/movie-streaming-go/server/routes"
)

func main() {

	router := gin.Default()
	routes.SetupProtectedRoutes(router)
	routes.SetupUnprotectedRoutes(router)

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server", err)
	}
}
