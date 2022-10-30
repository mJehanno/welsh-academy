package user

import (
	"crypto/sha256"
	"database/sql/driver"
	"fmt"

	"github.com/mjehanno/welsh-academy/pkg/recipe"
	"gorm.io/gorm"
)

type Role string

const (
	BasicUser     Role = "basicuser"
	CheddarExpert Role = "cheddarexpert"
	Admin         Role = "admin"
)

func (r *Role) Scan(value interface{}) error {
	*r = Role(value.([]byte))
	return nil
}

func (r Role) Value() (driver.Value, error) {
	return string(r), nil
}

// User represent user.
type User struct {
	gorm.Model
	// The name of the user.
	Username string `gorm:"type:varchar(40);unique;not null;default:null" example:"admin"`
	// The password of the user
	Password string `gorm:"size:255;not null;default:null" json:",omitempty" example:"admin"`
	// The user's favorites recipes
	FavoritesRecipes []recipe.Recipe `gorm:"many2many:favorite_recipe" swaggerignore:"true"`
	Role             Role            `gorm:"type:roles"`
}

// NewUserService is the constructor for a UserService.
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

// UserService is a service made to manage user related queries.
type UserService struct {
	db *gorm.DB
}

// CreateUser takes a user and insert it in database, it returns the id of inserted user or an error.
func (us *UserService) CreateUser(user User) (uint, error) {
	h := sha256.New()
	h.Write([]byte(user.Password))
	hash := h.Sum(nil)
	user.Password = fmt.Sprintf("%x", string(hash))

	result := us.db.Create(&user)

	return user.ID, result.Error
}

// LogUser verifies user credential to log him or not.
func (us *UserService) LogUser(user User) (*User, error) {
	h := sha256.New()
	h.Write([]byte(user.Password))
	hash := h.Sum(nil)
	var dbUser User

	err := us.db.Where("username = ?", user.Username).Where("password = ?", fmt.Sprintf("%x", string(hash))).Omit("password", "created_at", "deleted_at", "updated_at").First(&dbUser).Error
	if err != nil {
		return nil, err
	}

	return &dbUser, nil
}

// AddFavoriteRecipe takes a recipe and a userID (uint) and add the recipe to the corresponding user.
func (us *UserService) AddFavoriteRecipe(recipe recipe.Recipe, userID uint) error {
	var user User
	us.db.Where("id=?", userID).First(&user)

	return us.db.Model(&user).Association("FavoritesRecipes").Append(&recipe)
}

// GetFavoriteRecipe takes a userID(uint) and return his favorites recipes.
func (us *UserService) GetFavoriteRecipe(userID uint) ([]recipe.Recipe, error) {
	var user User
	var recipes []recipe.Recipe
	us.db.Where("id=?", userID).First(&user)

	err := us.db.Preload("Ingredients").Model(&user).Association("FavoritesRecipes").Find(&recipes)

	return recipes, err
}

// DeleteFavoriteRecipe takes a recipe and a userID and remove this recipe from user's favorite recipe.
func (us *UserService) DeleteFavoriteRecipe(recipe recipe.Recipe, userID uint) error {
	var user User
	us.db.Where("id=?", userID).First(&user)

	return us.db.Model(&user).Association("FavoritesRecipes").Delete(&recipe)
}
