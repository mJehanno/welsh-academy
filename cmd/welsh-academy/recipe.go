package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mjehanno/welsh-academy/pkg/error"
	"github.com/mjehanno/welsh-academy/pkg/ingredient"
	"github.com/mjehanno/welsh-academy/pkg/recipe"
)

// @Summary      Get All Recipe
// @Description  Get the list of every created recipe with their ingredient.
// @Tags         recipes
// @Produce      json
// @Param	ingredient query []string false "filter by ingredient"
// @Success      200  {array}  recipe.Recipe
// @Failure      400
// @Failure      500
// @Router       /recipes [get]
func getRecipeEndoint(c *gin.Context) {
	ingredientQuery := c.QueryArray("ingredient")
	if len(ingredientQuery) > 0 {
		ingredientsName := ingredientQuery
		if len(ingredientQuery) == 1 && strings.Contains(ingredientQuery[0], ",") {
			ingredientsName = strings.Split(ingredientQuery[0], ",")
		}
		ingredients := make([]ingredient.Ingredient, len(ingredientsName))

		for i, name := range ingredientsName {
			ing, err := ingredientService.GetIngredientByName(name)
			if err != nil {
				c.JSON(http.StatusBadRequest, error.ErrorResponse{ErrorMessage: err.Error()})
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
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, recipes)
}

// @Summary      Create a Recipe
// @Description  Create a Recipe with one or multiple ingredient.
// @Tags         recipes
// @Accept       json
// @Produce      json
// @Param recipe body recipe.Recipe true "recipe to create"
// @Success      201  {integer} id
// @Failure      400  {object}  error.ErrorResponse
// @Failure      500
// @Router       /recipes [post]
func createRecipeEndpoint(c *gin.Context) {
	var json recipe.Recipe

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, error.ErrorResponse{ErrorMessage: err.Error()})
		return
	}

	if json.Name == "" || len(json.Ingredients) == 0 {
		c.JSON(http.StatusBadRequest, error.ErrorResponse{ErrorMessage: "can't create a recipe without a name or without ingredients"})
		return
	}

	id, err := recipeService.CreateRecipe(json)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}
