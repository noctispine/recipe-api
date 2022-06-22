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

	"github.com/gin-contrib/sessions"
	redi "github.com/gin-contrib/sessions/redis"
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
var collectionUsers *mongo.Collection
var recipesHandler *handlers.RecipesHandler
var authHandler *handlers.AuthHandler

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

	collectionRecipes = client.Database(os.Getenv("MONGO_DB")).Collection("recipes")
	collectionUsers = client.Database(os.Getenv("MONGO_DB")).Collection("users")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URI"),
		Password: "",
		DB:       0,
	})

	status := redisClient.Ping(ctx)
	fmt.Println("redis status ", status)

	recipesHandler = handlers.NewRecipesHandler(ctx, collectionRecipes, redisClient)

	authHandler = handlers.NewAuthHandler(ctx, collectionUsers)
}

func main() {
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	router := gin.Default()
	store, _ := redi.NewStore(10, "tcp", os.Getenv("REDIS_URI"), os.Getenv("REDIS_PASSWORD"), []byte(os.Getenv("REDIS_SECRET")))
	router.Use(sessions.Sessions("recipes_api", store))
	router.Use(handlers.CustomLogger())
	router.POST("/signin", authHandler.SignInHandler)
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.POST("/refresh", authHandler.RefreshHandler)
	router.POST("/register", authHandler.RegisterHandler)

	authorized := router.Group("/")
	authorized.Use(authHandler.AuthMiddleware())
	{
		authorized.POST("/recipes", recipesHandler.NewRecipeHandler)
		authorized.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
		authorized.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
		authorized.GET("/signout", authHandler.SignOutHandler)
	}
	// router.GET("/recipes/search", SearchRecipesHandler)
	router.Run()

}
