package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kataras/jwt"
	"github.com/mjehanno/welsh-academy/pkg/error"
	"github.com/mjehanno/welsh-academy/pkg/ingredient"
	"github.com/mjehanno/welsh-academy/pkg/user"
)

// @Summary      Create an Ingredient
// @Description  Create an ingredient that you'll be able to use in a recipe.
// @Tags         ingredients
// @Accept       json
// @Produce      json
// @Param ingredient body ingredient.Ingredient true "ingredient to create"
// @Success      201  {integer}  id
// @Failure      400  {object}  error.ErrorResponse
// @Failure   	 401
// @Failure 	 	 403
// @Failure      500
// @Router       /ingredients [post]
func createIngredientEndpoint(c *gin.Context) {
	var json ingredient.Ingredient

	cookie, err := c.Cookie("jwt")
	if err != nil {
		c.JSON(http.StatusUnauthorized, nil)
		return
	}

	verifiedToken, err := jwt.Verify(jwt.HS256, sharedKey, []byte(cookie))
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	var currentUser user.User
	err = verifiedToken.Claims(&currentUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	if currentUser.Role != user.CheddarExpert {
		c.JSON(http.StatusForbidden, nil)
		return
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, error.ErrorResponse{ErrorMessage: err.Error()})
		return
	}

	if json.Name == "" {
		c.JSON(http.StatusBadRequest, error.ErrorResponse{ErrorMessage: "can't create ingredient with empty name"})
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
