package user

import (
	"database/sql"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mjehanno/welsh-academy/pkg/ingredient"
	"github.com/mjehanno/welsh-academy/pkg/recipe"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var mock sqlmock.Sqlmock
var db *sql.DB
var userService *UserService

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

	userService = NewUserService(gdb)

	return func(t *testing.T) {
		defer db.Close()

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}

func TestCreateUserSucceed(t *testing.T) {
	tearDown := Setup(t)
	defer tearDown(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("created_at","updated_at","deleted_at","username","password") VALUES ($1,$2,$3,$4,$5) RETURNING "id","username","password"`)).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "cam-amber", "5de4c437b552985b0fa4a9566a60d767ab89310343e4c5e3d7a373bc1b68747b").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	_, err := userService.CreateUser(User{Username: "cam-amber", Password: "mytopsecretpassword"})
	if err != nil {
		t.Errorf("error occured while it shouldn't have : %s", err.Error())
	}
}

func TestCreateUserFailOnEmptyName(t *testing.T) {
	tearDown := Setup(t)
	defer tearDown(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("created_at","updated_at","deleted_at","password") VALUES ($1,$2,$3,$4) RETURNING "id","username","password"`)).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "5de4c437b552985b0fa4a9566a60d767ab89310343e4c5e3d7a373bc1b68747b").WillReturnError(fmt.Errorf("can't create user with empty name"))
	mock.ExpectRollback()
	_, err := userService.CreateUser(User{Username: "", Password: "mytopsecretpassword"})
	if err == nil {
		t.Errorf("error did not occured while it should have")
	}
}

func TestCreateUserFailOnEmptyPassword(t *testing.T) {
	tearDown := Setup(t)
	defer tearDown(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("created_at","updated_at","deleted_at","username","password") VALUES ($1,$2,$3,$4,$5) RETURNING "id","username","password"`)).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "cam-amber", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855").WillReturnError(fmt.Errorf("can't create user with empty password"))
	mock.ExpectRollback()
	_, err := userService.CreateUser(User{Username: "cam-amber", Password: ""})
	if err == nil {
		t.Errorf("error did not occured while it should have")
	}
}

func TestLogUserSucceed(t *testing.T) {
	tearDown := Setup(t)
	defer tearDown(t)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users"."id","users"."username" FROM "users" WHERE username = $1 AND password = $2 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).WithArgs("cam-amber", "5de4c437b552985b0fa4a9566a60d767ab89310343e4c5e3d7a373bc1b68747b").WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "cam-amber"))

	user, err := userService.LogUser(User{Username: "cam-amber", Password: "mytopsecretpassword"})
	if err != nil {
		t.Errorf("error occured while it shouldn't have : %s", err.Error())
	}
	if user.Password != "" {
		t.Error("password should be empty here for security reason, we do not want to send password back in the frontend !")
	}

}

func TestLogUserFail(t *testing.T) {
	tearDown := Setup(t)
	defer tearDown(t)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users"."id","users"."username" FROM "users" WHERE username = $1 AND password = $2 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).WithArgs("cam-amber", "ed9909730fb6e9af1d563ce4a0019f0006141d0d818509ea4c1babf25821ecba").WillReturnError(fmt.Errorf("can't find matching user"))

	_, err := userService.LogUser(User{Username: "cam-amber", Password: "mynotsosecretpassword"})
	if err == nil {
		t.Errorf("error did not occured while it should have")
	}

}

func TestAddFavoriteRecipeSucceed(t *testing.T) {
	tearDown := Setup(t)
	defer tearDown(t)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id=$1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "cam-amber"))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "updated_at"=$1 WHERE "users"."deleted_at" IS NULL AND "id" = $2`)).WithArgs(sqlmock.AnyArg(), 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "recipes" ("created_at","updated_at","deleted_at","name") VALUES ($1,$2,$3,$4) ON CONFLICT DO NOTHING RETURNING "id","name"`)).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "welsh").WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "welsh"))
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "favorite_recipe" ("user_id","recipe_id") VALUES ($1,$2) ON CONFLICT DO NOTHING`)).WithArgs(1, 0).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := userService.AddFavoriteRecipe(recipe.Recipe{Name: "welsh", Ingredients: []*ingredient.Ingredient{{Name: "cheddar"}, {Name: "bière brune"}}}, 1)

	if err != nil {
		t.Errorf("error occured while it shouldn't have : %s", err.Error())
	}
}

func TestAddFavoriteRecipeFail(t *testing.T) {
	tearDown := Setup(t)
	defer tearDown(t)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id=$1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "cam-amber"))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "updated_at"=$1 WHERE "users"."deleted_at" IS NULL AND "id" = $2`)).WithArgs(sqlmock.AnyArg(), 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "recipes" ("created_at","updated_at","deleted_at","name") VALUES ($1,$2,$3,$4) ON CONFLICT DO NOTHING RETURNING "id","name"`)).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "welsh").WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "welsh"))
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "favorite_recipe" ("user_id","recipe_id") VALUES ($1,$2) ON CONFLICT DO NOTHING`)).WithArgs(1, 0).WillReturnError(fmt.Errorf("can't add a non existing recipe to favorites"))
	mock.ExpectRollback()

	err := userService.AddFavoriteRecipe(recipe.Recipe{Name: "welsh", Ingredients: []*ingredient.Ingredient{{Name: "cheddar"}, {Name: "bière brune"}}}, 1)

	if err == nil {
		t.Error("error occured while it shouldn't have")
	}
}

func TestGetFavoriteRecipeSucceed(t *testing.T) {
	tearDown := Setup(t)
	defer tearDown(t)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id=$1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "cam-amber"))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "recipes"."id","recipes"."created_at","recipes"."updated_at","recipes"."deleted_at","recipes"."name" FROM "recipes" JOIN "favorite_recipe" ON "favorite_recipe"."recipe_id" = "recipes"."id" AND "favorite_recipe"."user_id" = $1 WHERE "recipes"."deleted_at" IS NULL`)).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "welsh"))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "recipe_ingredient"."recipe_id","recipe_ingredient"."ingredient_id" FROM "recipe_ingredient" WHERE "recipe_ingredient"."recipe_id" = $1`)).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"recipe_id", "ingredient_id"}).AddRow(1, 1).AddRow(1, 2).AddRow(1, 3))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "ingredients"."id","ingredients"."created_at","ingredients"."updated_at","ingredients"."deleted_at","ingredients"."name" FROM "ingredients" WHERE "ingredients"."id" IN ($1,$2,$3) AND "ingredients"."deleted_at" IS NULL`)).WithArgs(1, 2, 3).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "cheddar").AddRow(2, "bière brune").AddRow(3, "pain"))

	_, err := userService.GetFavoriteRecipe(1)
	if err != nil {
		t.Errorf("error occured while it shouldn't have : %s", err.Error())
	}
}

func TestGetFavoriteRecipeFail(t *testing.T) {
	tearDown := Setup(t)
	defer tearDown(t)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id=$1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "cam-amber"))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "recipes"."id","recipes"."created_at","recipes"."updated_at","recipes"."deleted_at","recipes"."name" FROM "recipes" JOIN "favorite_recipe" ON "favorite_recipe"."recipe_id" = "recipes"."id" AND "favorite_recipe"."user_id" = $1 WHERE "recipes"."deleted_at" IS NULL`)).WithArgs(1).WillReturnError(fmt.Errorf("record not found"))

	_, err := userService.GetFavoriteRecipe(1)
	if err == nil {
		t.Error("error occured while it shouldn't have")
	}
}

func TestDeleteFavoriteRecipeSucceed(t *testing.T) {
	tearDown := Setup(t)
	defer tearDown(t)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id=$1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "cam-amber"))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "favorite_recipe" WHERE "favorite_recipe"."user_id" = $1 AND "favorite_recipe"."recipe_id" IN (NULL)`)).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := userService.DeleteFavoriteRecipe(recipe.Recipe{Name: "welsh", Ingredients: []*ingredient.Ingredient{{Name: "cheddar"}, {Name: "bière brune"}}}, 1)
	if err != nil {
		t.Errorf("error occured while it shouldn't have : %s", err.Error())
	}
}

func TestDeleteFavoriteRecipeFail(t *testing.T) {
	tearDown := Setup(t)
	defer tearDown(t)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id=$1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "cam-amber"))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "favorite_recipe" WHERE "favorite_recipe"."user_id" = $1 AND "favorite_recipe"."recipe_id" IN (NULL)`)).WithArgs(1).WillReturnError(fmt.Errorf("can't delete inexistent record"))
	mock.ExpectRollback()

	err := userService.DeleteFavoriteRecipe(recipe.Recipe{Name: "welsh", Ingredients: []*ingredient.Ingredient{{Name: "cheddar"}, {Name: "bière brune"}}}, 1)
	if err == nil {
		t.Error("error occured while it shouldn't have")
	}
}
