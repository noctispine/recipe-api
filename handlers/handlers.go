package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/noctispine/recipe-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RecipesHandler struct {
	collectionRecipes *mongo.Collection
	ctx               context.Context
	redisClient       *redis.Client
}

func NewRecipesHandler(ctx context.Context, collection *mongo.
	Collection, redisClient *redis.Client) *RecipesHandler {
	return &RecipesHandler{
		collectionRecipes: collection,
		ctx:               ctx,
		redisClient:       redisClient,
	}
}

// swagger:operation POST /recipes recipes newRecipe
// Create new recipe
// ---
// produces:
// - application/json
// responses:
//    '201':
//      description: Successful
//    '400':
//      description: Bad Request
func (handler *RecipesHandler) NewRecipeHandler(c *gin.Context) {
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error:": err.Error()})

		return
	}

	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()
	_, err := handler.collectionRecipes.InsertOne(handler.ctx, recipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})

		return
	}

	handler.redisClient.Del(c, "recipes")
	c.JSON(http.StatusCreated, recipe)

}

// swagger:operation GET /recipes recipes listRecipes
// Returns list of recipes
// ---
// produces:
// - application/json
// responses:
// '200':
//         description: Successful operation
func (handler *RecipesHandler) ListRecipesHandler(c *gin.Context) {
	val, err := handler.redisClient.Get(c, "recipes").Result()
	if err == redis.Nil {
		log.Printf("Request to MongoDB")
		cur, err := handler.collectionRecipes.Find(handler.ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()})
		}
		defer cur.Close(handler.ctx)

		recipes := make([]models.Recipe, 0)

		for cur.Next(handler.ctx) {
			var recipe models.Recipe
			cur.Decode(&recipe)
			recipes = append(recipes, recipe)
		}

		data, _ := json.Marshal(recipes)
		handler.redisClient.Set(c, "recipes", string(data), 0)

		c.JSON(http.StatusOK, recipes)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
	} else {
		log.Printf("Request to Redis")
		recipes := make([]models.Recipe, 0)
		json.Unmarshal([]byte(val), &recipes)
		c.JSON(http.StatusOK, recipes)

	}

}

// swagger:operation PUT /recipes/{id} recipes updateRecipe
// Update an existing recipe
//
// ---
// produces:
// - application/json
// - application/xml
// - text/xml
// - text/html
// parameters:
// - name: id
//   in: path
//   description: Recipe ID
//   required: true
//   type: string
// responses:
//    '200':
//      description: Successful operation
//    '400':
//      description: Invalid input
//    '404':
//      description: Invalid recipe ID
func (handler *RecipesHandler) UpdateRecipeHandler(c *gin.Context) {
	id := c.Params.ByName("id")
	objectID, _ := primitive.ObjectIDFromHex(id)

	var incomingRecipe models.Recipe
	if err := c.ShouldBindJSON(&incomingRecipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": error.Error(err)})

		return
	}

	filter := bson.M{"_id": objectID}

	update := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "name", Value: incomingRecipe.Name},
		primitive.E{Key: "tags", Value: incomingRecipe.Tags},
		primitive.E{Key: "ingredients", Value: incomingRecipe.Ingredients},
		primitive.E{Key: "instructions", Value: incomingRecipe.Instructions}}}}

	result, err := handler.collectionRecipes.UpdateOne(handler.ctx, filter, update)

	if result.MatchedCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "this item not exists"})

		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})

		return
	}
	handler.redisClient.Del(c, "recipes")
	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe has been updated"})
}

// swagger:operation DELETE /recipes/{id} recipes deleteRecipe
// Delete an existing Recipe
// ---
// produces:
// - application/json
// parameters:
// - name: id
//   in: path
//   required: true
//   description: Recipe ID
//   type: string
// responses:
//    '200':
//      description: Successful
//    '404':
//      description: Recipe Not Found
func (handler *RecipesHandler) DeleteRecipeHandler(c *gin.Context) {
	id := c.Params.ByName("id")

	objectID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objectID}

	result, err := handler.collectionRecipes.DeleteOne(handler.ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})

		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "this item not exists"})

		return
	}

	handler.redisClient.Del(c, "recipes")
	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe has been deleted"})
}

// swagger:operation GET /recipes/search/ recipes searchRecipe
// return recipes according to given tag
// ---
// produces:
// - application/json
// parameters:
// - name: tag
//   in: query
//   required: false
//   description: Tag Name
//   type: string
// responses:
//    '200':
//      description: Successful
// func SearchRecipesHandler(c *gin.Context) {
// 	tag := c.Query("tag")
// 	newRecipes := make([]Recipe, 0)

// 	for _, recipe := range recipes {
// 		found := false
// 		for _, recipeTag := range recipe.Tags {
// 			if strings.EqualFold(tag, recipeTag) {
// 				found = true
// 				continue
// 			}
// 		}

// 		if found {
// 			newRecipes = append(newRecipes, recipe)
// 		}
// 	}

// 	c.JSON(http.StatusOK, newRecipes)
// }
