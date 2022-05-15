package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

var recipes []Recipe

type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
}

func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error:": err.Error()})

		return
	}

	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)
	c.JSON(http.StatusCreated, recipe)

}

func ListRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}

func UpdateRecipeHandler(c *gin.Context) {
	id := c.Params.ByName("id")
	var incomingRecipe Recipe
	if err := c.ShouldBindJSON(&incomingRecipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": error.Error(err)})

		return
	}

	idx := -1

	for i, recipe := range recipes {
		if recipe.ID == id {
			idx = i
		}
	}

	if idx == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found"})
		return
	}

	incomingRecipe.ID = id
	incomingRecipe.PublishedAt = recipes[idx].PublishedAt
	recipes[idx] = incomingRecipe

	c.JSON(http.StatusCreated, recipes[idx])
}

func DeleteRecipeHandler(c *gin.Context) {
	id := c.Params.ByName("id")
	idx := -1

	for i, recipe := range recipes {
		if recipe.ID == id {
			idx = i
		}
	}

	if idx == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found"})
		return
	}

	recipes = append(recipes[:idx], recipes[idx+1:]...)

	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe has been deleted"})
}

func SearchRecipesHandler(c *gin.Context) {
	tag := c.Query("tag")
	newRecipes := make([]Recipe, 0)

	for _, recipe := range recipes {
		found := false
		for _, recipeTag := range recipe.Tags {
			if strings.EqualFold(tag, recipeTag) {
				found = true
				continue
			}
		}

		if found {
			newRecipes = append(newRecipes, recipe)
		}
	}

	c.JSON(http.StatusOK, newRecipes)
}

func init() {
	recipes = make([]Recipe, 0)
}

func main() {
	router := gin.Default()
	router.GET("/recipes", ListRecipesHandler)
	router.POST("/recipes", NewRecipeHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	router.GET("/recipes/search", SearchRecipesHandler)
	router.Run()
}
