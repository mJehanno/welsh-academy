package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mjehanno/welsh-academy/pkg/error"
	"github.com/mjehanno/welsh-academy/pkg/ingredient"
)

// @Summary      Create an Ingredient
// @Description  Create an ingredient that you'll be able to use in a recipe.
// @Tags         ingredients
// @Accept       json
// @Produce      json
// @Param ingredient body ingredient.Ingredient true "ingredient to create"
// @Success      201  {integer}  id
// @Failure      400  {object}  error.ErrorResponse
// @Failure      500
// @Router       /ingredients [post]
func createIngredientEndpoint(c *gin.Context) {
	var json ingredient.Ingredient

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, error.ErrorResponse{ErrorMessage: err.Error()})
		return
	}

	id, err := ingredientService.CreateIngredient(json)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

// @Summary      Get the list of ingredients.
// @Description  Get the whole list of ingredients.
// @Tags         ingredients
// @Produce      json
// @Success      200  {array}  ingredient.Ingredient
// @Failure      500
// @Router       /ingredients [get]
func getIngredientEndpoint(c *gin.Context) {
	ingredients, err := ingredientService.GetAllIngredient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, ingredients)
}
