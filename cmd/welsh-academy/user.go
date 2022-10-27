package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{})
}
