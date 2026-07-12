package controllers

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"errors"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/justKody/movie-streaming-go/server/database"
	"github.com/justKody/movie-streaming-go/server/models"
	"github.com/justKody/movie-streaming-go/server/utils"
	"github.com/justKody/movie-streaming-go/server/utils/response"
	"github.com/tmc/langchaingo/llms/openai"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var movieCollection *mongo.Collection = database.OpenCollection("movies")
var rankingCollection *mongo.Collection = database.OpenCollection("rankings")

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
			response.SuccessResponse(c, http.StatusInternalServerError, err, "Failed to add movie.")
			return
		}

		response.SuccessResponse(c, http.StatusCreated, result)
	}
}

func AdminReviewUpdate() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, err := utils.GetRoleFromContext(c)
		if role != "admin" {
			response.ErrorResponse(c, http.StatusForbidden, "You are not authorized to update admin review")
			return
		}

		movieId := c.Param("imdb_id")
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		cancel()

		if movieId == "" {
			response.ErrorResponse(c, http.StatusBadRequest, "Movie Id required")
			return
		}

		var req struct {
			AdminReview string `json:"admin_review"`
		}

		var resp struct {
			RankingName  string `json:"ranking_name"`
			RankingValue int    `json:"ranking_value"`
			AdminReview  string `json:"admin_review"`
		}

		if err := c.ShouldBind((&req)); err != nil {
			response.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
			return
		}

		sentiment, rankVal, err := GetReviewRanking(req.AdminReview)

		if err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		filter := bson.M{
			"imdb_id": movieId,
		}

		update := bson.M{
			"$set": bson.M{
				"adming_review": req.AdminReview,
				"ranking": bson.M{
					"ranking_value": rankVal,
					"ranking_name":  sentiment,
				},
			},
		}

		result, err := movieCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		if result.MatchedCount == 0 {
			response.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}

		resp.RankingName = sentiment
		resp.RankingValue = rankVal
		resp.AdminReview = req.AdminReview

		response.SuccessResponse(c, http.StatusOK, resp)

	}
}

func GetReviewRanking(admin_review string) (string, int, error) {
	rankings, err := GetRankings()

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	if err != nil {
		return "", 0, err
	}

	sentimentDelimited := ""

	for _, ranking := range rankings {
		if ranking.RankingValue != 999 {
			sentimentDelimited = sentimentDelimited + ranking.RankingName + ","
		}
	}

	sentimentDelimited = strings.Trim(sentimentDelimited, ",")

	err = godotenv.Load()

	if err != nil {
		log.Println("Warning: Unable to find .env file")
	}

	ApiKey := os.Getenv("LLM_API_KEY")
	BaseUrl := os.Getenv("LLM_API_KEY")
	Model := os.Getenv("LLM_API_KEY")

	if ApiKey == "" {
		return "", 0, errors.New("could not read OPENAI_API_KEY")
	}

	llm, err := openai.New(
		openai.WithToken(ApiKey),
		openai.WithBaseURL(BaseUrl),
		openai.WithModel(Model),
	)

	if err != nil {
		return "", 0, errors.New("Something went wrong with connection to LLM model")
	}

	base_prompt_template := os.Getenv("BASE_PROMPT_TEMPLATE")

	base_prompt := strings.Replace(base_prompt_template, "{rankings}", sentimentDelimited, 1)

	response, err := llm.Call(ctx, base_prompt+admin_review)

	if err != nil {
		return "", 0, err
	}

	rankVal := 0

	for _, ranking := range rankings {
		if ranking.RankingName == response {
			rankVal = ranking.RankingValue
			break
		}
	}

	return response, rankVal, nil
}

// excellence, good, ok, bad etc
func GetRankings() ([]models.Ranking, error) {
	var rankings []models.Ranking

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	cursor, err := rankingCollection.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &rankings); err != nil {
		return nil, err
	}

	return rankings, nil

}

func GetRecommendedMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId, err := utils.GetUserIdFromContext(c)

		if err != nil {
			response.ErrorResponse(c, http.StatusBadRequest, "User Id in context not found")
			return
		}

		favGenres, err := GetUsersFavouriteGenres(userId)

		if err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to get favourite genres")
			return
		}

		// get movieRecemondation limit from .env
		movieRecommendationLimit := os.Getenv("MOVIE_RECOMMENDATION_LIMIT")
		if movieRecommendationLimit == "" {
			movieRecommendationLimit = "10"
		}

		limit, err := strconv.ParseInt(movieRecommendationLimit, 10, 64)
		if err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to convert movie recommendation limit to integer")
			return
		}

		filter := bson.M{
			"genres": bson.M{"$in": favGenres},
		}

		opts := options.Find().SetSort(bson.M{"ranking.ranking_value": -1}).SetLimit(limit)

		cursor, err := movieCollection.Find(ctx, filter, opts)

		if err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch movies")
		}
		defer cursor.Close(ctx)

		var movies []models.Movie

		if err = cursor.All(ctx, &movies); err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to decode movies")
		}

		response.SuccessResponse(c, http.StatusOK, movies)
	}
}

func GetUsersFavouriteGenres(userId string) ([]string, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userId}

	projection := bson.M{
		"favourite_genres.genre_name": 1,
		"_id":                         0,
	}

	opts := options.FindOne().SetProjection(projection)
	var result bson.M

	err := userCollection.FindOne(ctx, filter, opts).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []string{}, nil
		}
	}

	favGenresArray, ok := result["favourite_genres"].(bson.A)

	if !ok {
		return []string{}, errors.New("favourite_genres is not an array")
	}

	favGenres := make([]string, len(favGenresArray))

	for i, genre := range favGenresArray {
		genreMap, ok := genre.(bson.M)
		if !ok {
			return nil, errors.New("invalid genre object")
		}

		name, ok := genreMap["genre_name"].(string)
		if !ok {
			return nil, errors.New("genre_name is not a string")
		}

		favGenres[i] = name
	}
	return favGenres, nil
}
