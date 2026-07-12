package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/justKody/movie-streaming-go/server/database"
	"github.com/justKody/movie-streaming-go/server/routes"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func main() {

	router := gin.Default()

	var Client *mongo.Client = database.Connect()
	routes.SetupProtectedRoutes(router, Client)
	routes.SetupUnprotectedRoutes(router, Client)

	if Client == nil {
		log.Fatal("Failed to connect to MongoDB")
	}
	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server", err)
	}
}
