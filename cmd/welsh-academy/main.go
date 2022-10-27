package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mjehanno/welsh-academy/pkg/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	var err error

	dsn := "host=localhost user=nimda password=nimda dbname=welsh port=5432 sslmode=disable "
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	db.AutoMigrate(&user.User{})

	if err != nil {
		log.Fatalf("couldn't connect to database : %w", err)
	}

	r := gin.Default()
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			user := v1.Group("/users")
			{
				user.POST("/", createUserEndpoint)
			}
		}
	}

	r.Run()
}

func createUserEndpoint(c *gin.Context) {
	userService := user.NewUserService(db)

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
