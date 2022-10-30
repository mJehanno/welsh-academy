package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mjehanno/welsh-academy/pkg/error"
	"github.com/mjehanno/welsh-academy/pkg/recipe"
	"github.com/mjehanno/welsh-academy/pkg/user"
	"gorm.io/gorm"
)

// @Summary Create user
// @Schemes
// @Description Create a user with a username and a password (hashed)
// @Tags users
// @Accept json
// @Produce json
// @Param user body user.User true "user that need to be created"
// @Success 201 {number} id
// @Failure 400 {object} error.ErrorResponse
// @Failure 500
// @Router /users [post]
func createUserEndpoint(c *gin.Context) {
	var json user.User

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, error.ErrorResponse{ErrorMessage: err.Error()})
		return
	}

	if json.Password == "" || json.Username == "" {
		c.JSON(http.StatusBadRequest, error.ErrorResponse{ErrorMessage: "can't create user with empty username/password"})
		return
	}

	id, err := userService.CreateUser(json)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

// @Summary Log a user
// @Schemes
// @Description Log a user with his username and password, returns some cookie (but not made with cheese).
// @Tags users
// @Accept json
// @Produce json
// @Param user body user.User true "user information in order to log in"
// @Success 200
// @Failure 400 {object} error.ErrorResponse
// @Failure 500
// @Router /users/login [post]
func loginUserEndpoint(c *gin.Context) {
	var json user.User

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, error.ErrorResponse{ErrorMessage: err.Error()})
		return
	}

	user, err := userService.LogUser(json)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, error.ErrorResponse{ErrorMessage: "wrong data for user/password"})
			return
		}

		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	if user != nil {
		c.SetCookie("user", user.Username, 36000, "/", "localhost", false, true)
		c.SetCookie("userID", strconv.Itoa(int(user.ID)), 36000, "/", "localhost", false, true)
		c.JSON(http.StatusOK, nil)
		return
	}

	c.JSON(http.StatusBadRequest, error.ErrorResponse{ErrorMessage: "wrong data for username/password"})
}

// @Summary      Flag a favorite recipe
// @Description  Flag a favorite recipe by adding an entry in db
// @Tags         favorites
// @Accept       json
// @Produce      json
// @Param recipe body recipe.Recipe true "flagged favorite recipe"
// @Success      201
// @Failure      400  {object}  error.ErrorResponse
// @Failure      401
// @Failure      500
// @Router       /users/favorites [post]
func createFavoriteRecipeEndpoint(c *gin.Context) {
	var json recipe.Recipe

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, error.ErrorResponse{ErrorMessage: err.Error()})
		return
	}

	cookie, err := c.Cookie("userID")
	if err != nil {
		c.JSON(http.StatusUnauthorized, nil)
		return
	}

	userId, err := strconv.ParseUint(cookie, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, error.ErrorResponse{ErrorMessage: err.Error()})
		return
	}

	err = userService.AddFavoriteRecipe(json, uint(userId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusCreated, nil)
}

// @Summary      Get favorites recipe
// @Description  Get a user's favorites recipe
// @Tags         favorites
// @Produce      json
// @Success      200  {array}  recipe.Recipe
// @Failure      400  {object}  error.ErrorResponse
// @Failure      401
// @Failure      500
// @Router       /users/favorites [get]
func getFavoriteRecipeEndpoint(c *gin.Context) {
	cookie, err := c.Cookie("userID")
	if err != nil {
		c.JSON(http.StatusUnauthorized, nil)
		return
	}

	userId, err := strconv.ParseUint(cookie, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, error.ErrorResponse{ErrorMessage: err.Error()})
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

// @Summary      Unflag favorite recipe
// @Description  Unflag a favorite recipe by deleting it
// @Tags         favorites
// @Produce      json
// @Param        id   path      int  true  "Recipe ID"
// @Success      204
// @Failure      400  {object}  error.ErrorResponse
// @Failure      401
// @Failure      500
// @Router       /users/favorites/{id} [delete]
func deleteFavoriteRecipeEndpoint(c *gin.Context) {
	cookie, err := c.Cookie("userID")
	if err != nil {
		c.JSON(http.StatusUnauthorized, nil)
		return
	}

	userId, err := strconv.ParseUint(cookie, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, error.ErrorResponse{ErrorMessage: err.Error()})
		return
	}

	paramId := c.Param("recipeId")

	recipeID, err := strconv.ParseUint(paramId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, error.ErrorResponse{ErrorMessage: err.Error()})
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
