package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mjehanno/welsh-academy/pkg/recipe"
)

func createIngredientEndpoint(c *gin.Context) {
	var json recipe.Ingredient

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := ingredientService.CreateIngredient(json)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

func getIngredientEndpoint(c *gin.Context) {
	ingredients, err := ingredientService.GetAllIngredient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, ingredients)
}
