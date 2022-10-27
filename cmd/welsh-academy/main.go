package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mjehanno/welsh-academy/pkg/ingredient"
	"github.com/mjehanno/welsh-academy/pkg/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var userService *user.UserService
var ingredientService *ingredient.IngredientService

func init() {
	var err error
	dsn := "host=localhost user=nimda password=nimda dbname=welsh port=5432 sslmode=disable "
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("couldn't connect to database : %w", err)
	}
	db.AutoMigrate(&user.User{}, &ingredient.Ingredient{})
}

func main() {
	userService = user.NewUserService(db)
	ingredientService = ingredient.NewIngredientService(db)

	r := gin.Default()
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			user := v1.Group("/users")
			{
				user.POST("/", createUserEndpoint)
			}
			ingredient := v1.Group("/ingredients")
			{
				ingredient.POST("/", createIngredientEndpoint)
				ingredient.GET("/", getIngredientEndpoint)
			}
		}
	}

	r.Run()
}

func createUserEndpoint(c *gin.Context) {
	var json user.User

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := userService.CreateUser(json)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

func createIngredientEndpoint(c *gin.Context) {
	var json ingredient.Ingredient

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := ingredientService.CreateIngredient(json)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

func getIngredientEndpoint(c *gin.Context) {
	ingredients, err := ingredientService.GetAllIngredient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
	}

	c.JSON(http.StatusAccepted, ingredients)
}
