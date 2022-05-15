package main

import (
	"log"
	"net/http"
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

func UpdateRecipe(c *gin.Context) {
	id := c.Params.ByName("id")
	log.Println(id)
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

func init() {
	recipes = make([]Recipe, 0)
}

func main() {
	router := gin.Default()
	router.GET("/recipes", ListRecipesHandler)
	router.POST("/recipes", NewRecipeHandler)
	router.PUT("/recipes/:id", UpdateRecipe)
	router.Run()
}
