package repository

import (
	"gogin/internal/model"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestPersonRepository_GetPersonById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := PersonRepository{Db: db}
	id := "1"

	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "country_code", "updated_at", "created_at"}).
		AddRow(1, "John", "Doe", "US", "2023-01-01", "2023-01-01")

	mock.ExpectQuery("SELECT id, first_name, last_name, country_code, updated_at, created_at FROM person WHERE id = ?").
		WithArgs(id).
		WillReturnRows(rows)

	res, err := repo.GetPersonById(id)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), res.Id)
	assert.Equal(t, "John", res.FirstName)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPersonRepository_GetMany(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := PersonRepository{Db: db}
	limit, offset := 10, 0

	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "country_code"}).
		AddRow(1, "John", "Doe", "US").
		AddRow(2, "Jane", "Smith", "CA")

	mock.ExpectQuery("SELECT id, first_name, last_name, country_code FROM person LIMIT \\? OFFSET \\?").
		WithArgs(limit, offset).
		WillReturnRows(rows)

	res, err := repo.GetMany(limit, offset)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(1), res[0].Id)
	assert.Equal(t, "John", res[0].FirstName)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPersonRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := PersonRepository{Db: db}
	person := model.Person{
		FirstName:   "John",
		LastName:    "Doe",
		CountryCode: "US",
	}

	mock.ExpectExec("INSERT INTO person").
		WithArgs(person.FirstName, person.LastName, person.CountryCode).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(person)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
