package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mjehanno/welsh-academy/pkg/recipe"
	"github.com/mjehanno/welsh-academy/pkg/user"
)

func createUserEndpoint(c *gin.Context) {
	var json user.User

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if json.Password == "" || json.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	id, err := userService.CreateUser(json)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

func loginUserEndpoint(c *gin.Context) {
	var json user.User

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := userService.LogUser(json)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	if user != nil {
		c.SetCookie("user", user.Username, 36000, "/", "localhost", false, true)
		c.SetCookie("userID", strconv.Itoa(int(user.ID)), 36000, "/", "localhost", false, true)
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{})
}

func createFavoriteRecipeEndpoint(c *gin.Context) {
	var json recipe.Recipe

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cookie, err := c.Cookie("userID")
	if err != nil {
		c.JSON(http.StatusUnauthorized, nil)
		return
	}

	userId, err := strconv.ParseUint(cookie, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	err = userService.AddFavoriteRecipe(json, uint(userId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusCreated, nil)
}

func getFavoriteRecipeEndpoint(c *gin.Context) {
	cookie, err := c.Cookie("userID")
	if err != nil {
		c.JSON(http.StatusUnauthorized, nil)
		return
	}

	userId, err := strconv.ParseUint(cookie, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	recipes, err := userService.GetFavoriteRecipe(uint(userId))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, recipes)
}

func deleteFavoriteRecipeEndpoint(c *gin.Context) {
	cookie, err := c.Cookie("userID")
	if err != nil {
		c.JSON(http.StatusUnauthorized, nil)
		return
	}

	userId, err := strconv.ParseUint(cookie, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	paramId := c.Param("recipeId")

	recipeID, err := strconv.ParseUint(paramId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	recipe, err := recipeService.GetRecipeById(uint(recipeID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	err = userService.DeleteFavoriteRecipe(recipe, uint(userId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusNoContent, nil)

}
