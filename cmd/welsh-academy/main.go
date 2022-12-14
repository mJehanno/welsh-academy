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
var sharedKey = []byte("asupersecrettokenthatnooneshouldknow")

func init() {
	var err error
	dsn := "host=" + os.Getenv("DB_HOST") + " user=" + os.Getenv("DB_USER") + " password=" + os.Getenv("DB_PASS") + " dbname=" + os.Getenv("DB_NAME") + " port=" + os.Getenv("DB_PORT") + " sslmode=disable "
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("couldn't connect to database : %s", err.Error())
	}

	result := db.Exec(`DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'roles') THEN
        CREATE TYPE roles as ENUM
        (
            'basicuser',
            'cheddarexpert',
            'admin'
        );
    END IF;
END$$;`)
	if result.Error != nil {
		log.Fatalf("couldn'nt create role type in db : %s", err.Error())
	}

	err = db.AutoMigrate(&user.User{}, &ingredient.Ingredient{}, &recipe.Recipe{})
	if err != nil {
		log.Fatalf("couldn't not create the database via migration : %s", err.Error())
	}

	userService = user.NewUserService(db)
	ingredientService = ingredient.NewIngredientService(db)
	recipeService = recipe.NewRecipeService(db)

}

func createAdminUser() {
	var admin user.User
	result := db.Model(&user.User{}).First(admin)
	if result.Error != nil {
		admin.Username = os.Getenv("ADMIN_USERNAME")
		admin.Password = os.Getenv("ADMIN_PASSWORD")
		admin.Role = user.Admin
		_, err := userService.CreateUser(admin)
		if err != nil {
			log.Println("couldn't create admin user, you'll need to create it manually in the database")
		}
	}
}

// @title           Welsh-Academy OpenAPI Spec
// @version         1.2.3
// @description     This is a rest api made to handle some recipe so please have a sit and chees... chill !
// @termsOfService  http://swagger.io/terms/
// @contact.name   API Support
// @contact.email  mathob.jehanno@hotmail.fr
// @license.name  GPL-3.0
// @host      localhost:9000
// @BasePath  /api/v1
func main() {
	createAdminUser()
	r := gin.Default()
	err := r.SetTrustedProxies(nil)
	if err != nil {
		log.Printf("couldn't unset trusted proxies on http server : %s", err.Error())
	}

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

	err = r.Run()
	if err != nil {
		log.Fatalf("couldn't start http server : %s", err.Error())
	}
}
