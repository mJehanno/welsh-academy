package ingredient

import (
	"database/sql"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var mock sqlmock.Sqlmock
var db *sql.DB
var ingredientService *IngredientService

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

	ingredientService = NewIngredientService(gdb)

	return func(t *testing.T) {
		defer db.Close()

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}

func TestGetAllIngredient(t *testing.T) {
	tearDown := Setup(t)
	defer tearDown(t)

	rows := mock.NewRows([]string{"id", "name"}).AddRow(1, "cheddar").AddRow(2, "comte").AddRow(3, "camembert")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ingredients" WHERE "ingredients"."deleted_at" IS NULL`)).WillReturnRows(rows)

	_, err := ingredientService.GetAllIngredient()
	if err != nil {
		t.Errorf("error occured while it shouldn't have : %s", err.Error())
	}

}

func TestCreateIngredientSucceed(t *testing.T) {
	tearDown := Setup(t)
	defer tearDown(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ingredients" ("created_at","updated_at","deleted_at","name") VALUES ($1,$2,$3,$4) RETURNING "id"`)).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "brie").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	_, err := ingredientService.CreateIngredient(Ingredient{Name: "brie"})
	if err != nil {
		t.Errorf("error occured while it shouldn't have : %s", err.Error())
	}

}

func TestCreateIngredientFail(t *testing.T) {
	tearDown := Setup(t)
	defer tearDown(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ingredients" ("created_at","updated_at","deleted_at") VALUES ($1,$2,$3) RETURNING "id"`)).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnError(fmt.Errorf("can't create ingredient without name"))
	mock.ExpectRollback()

	_, err := ingredientService.CreateIngredient(Ingredient{Name: ""})
	if err == nil {
		t.Errorf("error did not occured while it should have : %s", err.Error())
	}

}

func TestGetIngredientByNameSucceed(t *testing.T) {
	tearDown := Setup(t)
	defer tearDown(t)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ingredients" WHERE name = $1 AND "ingredients"."deleted_at" IS NULL ORDER BY "ingredients"."id" LIMIT 1`)).WithArgs("cheddar").WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "cheddar"))

	_, err := ingredientService.GetIngredientByName("cheddar")
	if err != nil {
		t.Errorf("error occured while it shouldn't have : %s", err.Error())
	}

}

func TestGetIngredientByNameFail(t *testing.T) {
	tearDown := Setup(t)
	defer tearDown(t)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ingredients" WHERE name = $1 AND "ingredients"."deleted_at" IS NULL ORDER BY "ingredients"."id" LIMIT 1`)).WithArgs("mimolette").WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

	_, err := ingredientService.GetIngredientByName("mimolette")
	if err == nil {
		t.Errorf("error did not occured while it should have : %s", err.Error())
	}

}
