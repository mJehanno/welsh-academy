package user

import (
	"crypto/sha256"
	"fmt"

	"github.com/mjehanno/welsh-academy/pkg/recipe"
	"gorm.io/gorm"
)

// User represent a basic
type User struct {
	gorm.Model
	Username         string          `gorm:"type:varchar(40);unique`
	Password         string          `gorm:"size:255", json:",omitempty"`
	FavoritesRecipes []recipe.Recipe `gorm:"many2many:favorite_recipe"`
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

func (us *UserService) LogUser(user User) (*User, error) {
	h := sha256.New()
	h.Write([]byte(user.Password))
	hash := h.Sum(nil)
	var dbUser User

	err := us.db.Where("username = ?", user.Username).First(&dbUser).Error
	if err != nil {
		return nil, err
	}

	if dbUser.Password == fmt.Sprintf("%x", string(hash)) {
		dbUser.Password = ""
		return &dbUser, nil
	}

	return nil, nil
}

func (us *UserService) GetFavorites() ([]recipe.Recipe, error) {
	return nil, nil
}
