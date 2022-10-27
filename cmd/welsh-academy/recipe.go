package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mjehanno/welsh-academy/pkg/recipe"
)

func getRecipeEndoint(c *gin.Context) {
	if len(c.QueryArray("ingredient")) > 0 {
		ingredientsName := c.QueryArray("ingredient")
		ingredients := make([]recipe.Ingredient, len(ingredientsName))

		for i, name := range ingredientsName {
			ing, err := ingredientService.GetIngredientByName(name)
			if err != nil {
				c.JSON(http.StatusBadRequest, nil)
				return
			}

			ingredients[i] = ing
		}

		recipes, err := recipeService.GetRecipeByIngredient(ingredients)
		if err != nil {
			c.JSON(http.StatusInternalServerError, nil)
			return
		}

		c.JSON(http.StatusOK, recipes)
		return

	}

	recipes, err := recipeService.GetAllRecipe()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, recipes)
}

func createRecipeEndpoint(c *gin.Context) {
	var json recipe.Recipe

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := recipeService.CreateRecipe(json)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

func getRecipeByIngredientEndpoint(c *gin.Context) {
}
