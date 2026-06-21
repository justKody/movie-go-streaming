package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID             bson.ObjectID `bson:"_id,omitempty" json:"_id,omitempty" `
	UserID         string        `json:"user_id" bson:"user_id"`
	FirstName      string        `json:"first_name" bson:"first_name" validate:"required,min=2,max=100"`
	LastName       string        `json:"last_name" bson:"last_name" validate:"required,min=2,max=100"`
	Email          string        `json:"email" bson:"email" validate:"required,email"`
	Password       string        `json:"password" bson:"password" validate:"required,min=8"`
	Role           string        `json:"role" bson:"role" validate:"required,oneof=USER"`
	CreatedAt      time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at" bson:"updated_at"`
	Token          string        `json:"token" bson:"token"`
	RefreshToken   string        `json:"refresh_token" bson:"refresh_token"`
	FavouriteGenre []Genre       `json:"favourite_genres" bson:"favourite_genre" validate:"required,dive"`
}
