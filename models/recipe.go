package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// swagger:parameters recipes NewRecipe
type Recipe struct {
	// swagger:ignore
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Name         string             `json:"name" bson:"name"`
	Tags         []string           `json:"tags" bson:"tags"`
	Ingredients  []string           `json:"ingredients" bson:"ingredients"`
	Instructions []string           `json:"instructions" bson:"instructions"`
	PublishedAt  time.Time          `json:"publishedAt" bson:"publishedAt"`
}

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Username  string             `json:"username"`
	Password  string             `json:"password"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}
