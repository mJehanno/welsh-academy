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

	result := rs.db.Model(&Recipe{}).Preload("Ingredients").Find(&recipe)

	return recipe, result.Error
}

func (rs *RecipeService) CreateRecipe(recipe Recipe) (uint, error) {
	result := rs.db.Create(&recipe)

	return recipe.ID, result.Error
}

func (rs *RecipeService) GetRecipeById(recipeID uint) (Recipe, error) {
	var recipe Recipe
	result := rs.db.Where("id = ?", recipeID).Find(&recipe)

	return recipe, result.Error
}
