package recipe

import (
	"gorm.io/gorm"
)

type Recipe struct {
	gorm.Model
	Name        string
	Ingredients []*Ingredient `gorm:"many2many:recipe_ingredient;"`
}

type RecipeService struct {
	db *gorm.DB
}

func NewRecipeService(db *gorm.DB) *RecipeService {
	return &RecipeService{
		db: db,
	}
}

func (rs *RecipeService) GetAllRecipe() ([]Recipe, error) {
	var recipe []Recipe

	result := rs.db.Find(&recipe)

	return recipe, result.Error
}
