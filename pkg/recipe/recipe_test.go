package recipe

import (
	"database/sql"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mjehanno/welsh-academy/pkg/ingredient"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var mock sqlmock.Sqlmock
var db *sql.DB
var recipeService *RecipeService

func Setup(t *testing.T) func(t *testing.T) {
	var err error

	db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp)) // mock sql.DB
	if err != nil {
		t.Fatalf("error shouldn't have occured while mocking db")
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	gdb, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("error shouldn't have occured while opening the mocked db")
	}

	recipeService = NewRecipeService(gdb)

	return func(t *testing.T) {
		defer db.Close()

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}

func TestGetAllRecipe(t *testing.T) {
	tearDown := Setup(t)
	defer tearDown(t)

	rows := mock.NewRows([]string{"id", "name"}).AddRow(1, "welsh").AddRow(2, "raclette").AddRow(3, "tartiflette")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipes" WHERE "recipes"."deleted_at" IS NULL`)).WillReturnRows(rows)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipe_ingredient" WHERE "recipe_ingredient"."recipe_id" IN ($1,$2,$3)`)).WithArgs(1, 2, 3).WillReturnRows(sqlmock.NewRows([]string{"recipe_id", "ingredient_id"}).AddRow(1, 1).AddRow(1, 2).AddRow(1, 3).AddRow(2, 4).AddRow(3, 4))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ingredients" WHERE "ingredients"."id" IN ($1,$2,$3,$4) AND "ingredients"."deleted_at" IS NULL`)).WithArgs(1, 2, 3, 4).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "cheddar").AddRow(2, "bière brune").AddRow(3, "pain").AddRow(4, "reblochon"))

	_, err := recipeService.GetAllRecipe()
	if err != nil {
		t.Errorf("an error occured while it shouldn't have : %s", err.Error())
	}

}

func TestGetRecipeByIngredientSucceed(t *testing.T) {
	tearDown := Setup(t)
	defer tearDown(t)

	rows := mock.NewRows([]string{"id", "name"}).AddRow(1, "welsh")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "recipes"."id","recipes"."created_at","recipes"."updated_at","recipes"."deleted_at","recipes"."name" FROM "recipes" inner join recipe_ingredient ri on ri.recipe_id = recipes.id inner join ingredients i on ri.ingredient_id = i.id inner join recipe_ingredient rii on rii.recipe_id = recipes.id inner join ingredients ii on rii.ingredient_id = ii.id WHERE i.id=$1 AND ii.id=$2 AND "recipes"."deleted_at" IS NULL`)).WithArgs(0, 0).WillReturnRows(rows)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipe_ingredient" WHERE "recipe_ingredient"."recipe_id" = $1`)).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"recipe_id", "ingredient_id"}).AddRow(1, 1).AddRow(1, 2).AddRow(1, 3))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ingredients" WHERE "ingredients"."id" IN ($1,$2,$3) AND "ingredients"."deleted_at" IS NULL`)).WithArgs(1, 2, 3).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "cheddar").AddRow(2, "bière brune").AddRow(3, "pain"))

	_, err := recipeService.GetRecipeByIngredient([]ingredient.Ingredient{{Name: "cheddar"}, {Name: "bière brune"}})
	if err != nil {
		t.Errorf("an error occured while it shouldn't have : %s", err.Error())
	}
}

func TestGetRecipeByIngredientFail(t *testing.T) {
	tearDown := Setup(t)
	defer tearDown(t)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "recipes"."id","recipes"."created_at","recipes"."updated_at","recipes"."deleted_at","recipes"."name" FROM "recipes" inner join recipe_ingredient ri on ri.recipe_id = recipes.id inner join ingredients i on ri.ingredient_id = i.id inner join recipe_ingredient rii on rii.recipe_id = recipes.id inner join ingredients ii on rii.ingredient_id = ii.id WHERE i.id=$1 AND ii.id=$2 AND "recipes"."deleted_at" IS NULL`)).WithArgs(0, 0).WillReturnError(fmt.Errorf("no record found"))

	_, err := recipeService.GetRecipeByIngredient([]ingredient.Ingredient{{Name: "cheddar"}, {Name: "bière brune"}})
	if err == nil {
		t.Error("an error did not occured while it should have")
	}
}

func TestCreateRecipeSucceed(t *testing.T) {
	tearDown := Setup(t)
	defer tearDown(t)

	any := sqlmock.AnyArg()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "recipes" ("created_at","updated_at","deleted_at","name") VALUES ($1,$2,$3,$4) RETURNING "id","name"`)).WithArgs(any, any, any, "welsh").WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "welsh"))
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ingredients" ("created_at","updated_at","deleted_at","name") VALUES ($1,$2,$3,$4),($5,$6,$7,$8),($9,$10,$11,$12) ON CONFLICT DO NOTHING RETURNING "id","name"`)).WithArgs(any, any, any, "cheddar", any, any, any, "bière brune", any, any, any, "pain").WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "cheddar").AddRow(2, "bière brune").AddRow(3, "pain"))
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "recipe_ingredient" ("recipe_id","ingredient_id") VALUES ($1,$2),($3,$4),($5,$6) ON CONFLICT DO NOTHING`)).WithArgs(1, 0, 1, 0, 1, 0).WillReturnResult(sqlmock.NewResult(1, 3))
	mock.ExpectCommit()

	_, err := recipeService.CreateRecipe(Recipe{Name: "welsh", Ingredients: []*ingredient.Ingredient{
		{Name: "cheddar"},
		{Name: "bière brune"},
		{Name: "pain"},
	}})

	if err != nil {
		t.Errorf("an error occured while it shouldn't have : %s", err.Error())
	}
}

func TestCreateRecipeFail(t *testing.T) {
	tearDown := Setup(t)
	defer tearDown(t)

	any := sqlmock.AnyArg()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "recipes" ("created_at","updated_at","deleted_at","name") VALUES ($1,$2,$3,$4) RETURNING "id","name"`)).WithArgs(any, any, any, "welsh").WillReturnError(fmt.Errorf("recipe already exist"))
	mock.ExpectRollback()

	_, err := recipeService.CreateRecipe(Recipe{Name: "welsh", Ingredients: []*ingredient.Ingredient{
		{Name: "cheddar"},
		{Name: "bière brune"},
		{Name: "pain"},
	}})

	if err == nil {
		t.Error("an error did not occured while it should have")
	}
}

func TestGetRecipeByIdSucceed(t *testing.T) {
	tearDown := Setup(t)
	defer tearDown(t)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipes" WHERE id = $1 AND "recipes"."deleted_at" IS NULL`)).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "welsh"))

	_, err := recipeService.GetRecipeById(1)
	if err != nil {
		t.Errorf("an error occured while it shouldn't have : %s", err.Error())
	}
}

func TestGetRecipeByIdFail(t *testing.T) {
	tearDown := Setup(t)
	defer tearDown(t)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipes" WHERE id = $1 AND "recipes"."deleted_at" IS NULL`)).WithArgs(1).WillReturnError(fmt.Errorf("record not found"))

	_, err := recipeService.GetRecipeById(1)
	if err == nil {
		t.Error("an error did not occured while it should have")
	}
}
