package user

import (
	"gorm.io/gorm"
)

// User represent a basic
type User struct {
	gorm.Model
	Name string
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
	result := us.db.Create(&user)

	return user.ID, result.Error
}
