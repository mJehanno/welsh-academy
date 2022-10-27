package recipe

import "gorm.io/gorm"

type Ingredient struct {
	gorm.Model
	Name    string
	Recipes []*Recipe `gorm:"many2many:recipe_ingredient;"`
}

func NewIngredientService(db *gorm.DB) *IngredientService {
	return &IngredientService{
		db: db,
	}
}

type IngredientService struct {
	db *gorm.DB
}

func (is *IngredientService) CreateIngredient(ingredient Ingredient) (uint, error) {

	result := is.db.Create(&ingredient)

	return ingredient.ID, result.Error
}

func (is *IngredientService) GetAllIngredient() ([]Ingredient, error) {
	var ingredients []Ingredient

	result := is.db.Find(&ingredients)

	return ingredients, result.Error
}
