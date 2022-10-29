package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	docs "github.com/mjehanno/welsh-academy/docs"
	"github.com/mjehanno/welsh-academy/pkg/ingredient"
	"github.com/mjehanno/welsh-academy/pkg/recipe"
	"github.com/mjehanno/welsh-academy/pkg/user"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var userService *user.UserService
var ingredientService *ingredient.IngredientService
var recipeService *recipe.RecipeService

func init() {
	var err error
	dsn := "host=" + os.Getenv("DB_HOST") + " user=" + os.Getenv("DB_USER") + " password=" + os.Getenv("DB_PASS") + " dbname=" + os.Getenv("DB_NAME") + " port=" + os.Getenv("DB_PORT") + " sslmode=disable "
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("couldn't connect to database : %s", err.Error())
	}
	db.AutoMigrate(&user.User{}, &ingredient.Ingredient{}, &recipe.Recipe{})

	userService = user.NewUserService(db)
	ingredientService = ingredient.NewIngredientService(db)
	recipeService = recipe.NewRecipeService(db)
}

// @title           Welsh-Academy OpenAPI Spec
// @version         1.0
// @description     This is a rest api made to handle some recipe so please have a sit and chees... chill !
// @termsOfService  http://swagger.io/terms/
// @contact.name   API Support
// @contact.email  mathob.jehanno@hotmail.fr
// @license.name  Beerware 42.0
// @host      localhost:9000
// @BasePath  /api/v1
func main() {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	docs.SwaggerInfo.BasePath = "/api/v1"
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
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.Run()
}
