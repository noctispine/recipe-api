// Recipes API
//
// Very simple Recipes API
//
//     Schemes: http
//     Host: localhost:8080
//     BasePath: /
//     Version: 1.0.0
//     Contact: Eren Cam<erencam.dev@gmail.com> https://erencam.dev
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
// 	   swagger:meta
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/noctispine/recipe-api/handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var ctx context.Context
var err error
var client *mongo.Client
var collectionRecipes *mongo.Collection
var recipesHandler *handlers.RecipesHandler

func init() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading env file")
	}

	ctx = context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal("Error while connecting to Database")
	}

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")

	collectionRecipes = client.Database(os.Getenv("COLL_RECIPES")).Collection("recipes")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_RECIPES_ADDR"),
		Password: "",
		DB:       0,
	})

	status := redisClient.Ping(ctx)
	fmt.Println("redis status ", status)

	recipesHandler = handlers.NewRecipesHandler(ctx, collectionRecipes, redisClient)
}

func main() {
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	router := gin.Default()
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.POST("/recipes", recipesHandler.NewRecipeHandler)
	router.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
	router.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	// router.GET("/recipes/search", SearchRecipesHandler)
	router.Run()
}
