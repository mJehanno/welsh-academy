package recipe

import (
	"github.com/mjehanno/welsh-academy/pkg/ingredient"
	"gorm.io/gorm"
)

// swagger:model Recipe
// Recipe define a meal made with ingredients.
type Recipe struct {
	gorm.Model
	// The name of the Recipe
	Name string `example:"welsh"`
	// The list of ingredients in the recipe.
	Ingredients []*ingredient.Ingredient `gorm:"many2many:recipe_ingredient;"`
}

// RecipeService define a service made to handle recipes.
type RecipeService struct {
	db *gorm.DB
}

// NewRecipeService is the RecipeService constructor.
func NewRecipeService(db *gorm.DB) *RecipeService {
	return &RecipeService{
		db: db,
	}
}

// GetAllRecipe returns all recipe.
func (rs *RecipeService) GetAllRecipe() ([]Recipe, error) {
	var recipe []Recipe

	result := rs.db.Model(&Recipe{}).Preload("Ingredients").Find(&recipe)

	return recipe, result.Error
}

// GetRecipeByIngredient takes a list of ingredients and returns only the recipe that contains ALL the listed ingredients or an error.
func (rs *RecipeService) GetRecipeByIngredient(ingredients []ingredient.Ingredient) ([]Recipe, error) {
	var recipes []Recipe

	query := rs.db.Table("recipes")
	for i := range ingredients {
		tableAlias := ""
		for j := 0; j < i+1; j++ {
			tableAlias += "i"
		}
		query.Joins("inner join recipe_ingredient r" + tableAlias + " on r" + tableAlias + ".recipe_id = recipes.id")
		query.Joins("inner join ingredients " + tableAlias + " on r" + tableAlias + ".ingredient_id = " + tableAlias + ".id")
	}

	for i, ing := range ingredients {
		tableAlias := ""
		for j := 0; j < i+1; j++ {
			tableAlias += "i"
		}
		query.Where(tableAlias+".id=?", ing.ID)
	}

	result := query.Preload("Ingredients").Find(&recipes)
	return recipes, result.Error
}

// CreateRecipe takes a recipe object and insert it to DB, returning it's new ID or an error.
func (rs *RecipeService) CreateRecipe(recipe Recipe) (uint, error) {
	result := rs.db.Create(&recipe)

	return recipe.ID, result.Error
}

// GetRecipeById takes a recipe ID and returns the corresponding recipe or an error.
func (rs *RecipeService) GetRecipeById(recipeID uint) (Recipe, error) {
	var recipe Recipe
	result := rs.db.Where("id = ?", recipeID).Find(&recipe)

	return recipe, result.Error
}
