package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/justKody/movie-streaming-go/server/database"
	"github.com/justKody/movie-streaming-go/server/models"
	"github.com/justKody/movie-streaming-go/server/utils/response"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"github.com/go-playground/validator/v10"
)

var movieCollection *mongo.Collection = database.OpenCollection("movies")

var validate = validator.New()

func GetMovies() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var movies []models.Movie

		cursor, err := movieCollection.Find(ctx, bson.M{})

		if err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch movies.")
		}
		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &movies); err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to decode movies.")
		}

		response.SuccessResponse(c, http.StatusOK, movies)
	}
}

func GetMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		movieId := c.Param("imdb_id")

		if movieId == "" {
			response.ErrorResponse(c, http.StatusBadRequest, "Movie ID is required.")
		}

		var movie models.Movie

		err := movieCollection.FindOne(ctx, bson.M{
			"imdb_id": movieId,
		}).Decode(&movie)

		if err != nil {
			response.ErrorResponse(c, http.StatusNotFound, "Movie not found.")
			return
		}

		response.SuccessResponse(c, http.StatusOK, movie)

	}

}

func AddMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var movie models.Movie
		if err := c.ShouldBindJSON(&movie); err != nil {
			response.ErrorResponse(c, http.StatusBadRequest, "Invalid inputs.")
			return
		}
		if err := validate.Struct(movie); err != nil {
			response.ValidationErrorResponse(c, err)
			return
		}

		result, err := movieCollection.InsertOne(ctx, movie)

		if err != nil {
			response.SuccessResponse(c, http.StatusInternalServerError,err, "Failed to add movie.")
			return
		}

		response.SuccessResponse(c, http.StatusCreated, result)
	}
}
