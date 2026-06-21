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
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection("users")

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User

		if err := c.ShouldBindJSON(&user); err != nil {
			response.SuccessResponse(c, http.StatusBadRequest, "Invalid input data.")
		}
		if err := validate.Struct(user); err != nil {
			response.ValidationErrorResponse(c, err)
			return
		}

		// check if user already exist

		count, err := userCollection.CountDocuments(ctx, bson.M{
			"email": user.Email,
		})

		if err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to check existing user.")
			return
		}

		if count > 1 {
			response.ErrorResponse(c, http.StatusConflict, "User already exists.")
			return
		}

		// hash to password

		hashPass, err := hashPassword(user.Password)
		if err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Unabled to hash password.")
			return
		}

		user.Password = hashPass
		user.UserID = bson.NewObjectID().Hex()
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()

		// insert to db
		result, err := userCollection.InsertOne(ctx, user)

		if err != nil {
			response.SuccessResponse(c, http.StatusInternalServerError, err, "Failed to add movie.")
			return
		}

		response.SuccessResponse(c, http.StatusCreated, result)
	}
}

// func Login() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
// 		defer cancel()
// 	}
// }

func hashPassword(password string) (string, error) {
	hashpassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashpassword), nil
}
