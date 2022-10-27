package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mjehanno/welsh-academy/pkg/recipe"
	"github.com/mjehanno/welsh-academy/pkg/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var userService *user.UserService
var ingredientService *recipe.IngredientService
var recipeService *recipe.RecipeService

func init() {
	var err error
	dsn := "host=localhost user=nimda password=nimda dbname=welsh port=5432 sslmode=disable "
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("couldn't connect to database : %w", err)
	}
	db.AutoMigrate(&user.User{}, &recipe.Ingredient{}, &recipe.Recipe{})
}

func main() {
	userService = user.NewUserService(db)
	ingredientService = recipe.NewIngredientService(db)
	recipeService = recipe.NewRecipeService(db)

	r := gin.Default()
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			user := v1.Group("/users")
			{
				user.POST("/", createUserEndpoint)
				user.POST("/login", loginUserEndpoint)

				favorites := user.Group("/favorites")
				{
					favorites.POST("/", createFavoriteRecipeEndpoint)
					favorites.GET("/", getFavoriteRecipeEndpoint)
					favorites.DELETE("/:recipeId", deleteFavoriteRecipeEndpoint)
				}
			}
			ingredient := v1.Group("/ingredients")
			{
				ingredient.POST("/", createIngredientEndpoint)
				ingredient.GET("/", getIngredientEndpoint)
			}
			recipe := v1.Group("/recipes")
			{
				recipe.GET("/", getRecipeEndoint)
				recipe.POST("/", createRecipeEndpoint)
			}
		}
	}

	r.Run()
}
