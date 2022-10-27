package recipe

import "gorm.io/gorm"

// Ingredient defines a product in cooking.
type Ingredient struct {
	gorm.Model
	Name    string
	Recipes []*Recipe `gorm:"many2many:recipe_ingredient;"`
}

// NewIngredientService is the IngredientService constructor.
func NewIngredientService(db *gorm.DB) *IngredientService {
	return &IngredientService{
		db: db,
	}
}

// IngredientService is a service made to manage ingredients.
type IngredientService struct {
	db *gorm.DB
}

// CreateIngredient insert an ingredient in the database and return it's ID.
func (is *IngredientService) CreateIngredient(ingredient Ingredient) (uint, error) {

	result := is.db.Create(&ingredient)

	return ingredient.ID, result.Error
}

// GetIngredientByName takes the name of the ingredient and check if it's in the database returning the existing ingredient or an error.
func (is *IngredientService) GetIngredientByName(name string) (Ingredient, error) {
	var ingredient Ingredient

	result := is.db.Where("name = ?", name).First(&ingredient)

	return ingredient, result.Error
}

// GetAllIngredient returns a list containing all created ingredient.
func (is *IngredientService) GetAllIngredient() ([]Ingredient, error) {
	var ingredients []Ingredient

	result := is.db.Find(&ingredients)

	return ingredients, result.Error
}
