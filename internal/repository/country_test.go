package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCountryRepository_GetCountryByCode(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := CountryRepository{Db: db}
	code := "US"
	name := "United States"

	rows := sqlmock.NewRows([]string{"code", "name"}).
		AddRow(code, name)

	mock.ExpectQuery("SELECT code, name FROM country WHERE code = ?").
		WithArgs(code).
		WillReturnRows(rows)

	res, err := repo.GetCountryByCode(t.Context(), code)
	assert.NoError(t, err)
	assert.Equal(t, code, res.Code)
	assert.Equal(t, name, res.Name)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCountryRepository_GetMany(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := CountryRepository{Db: db}
	limit, offset := 10, 0

	rows := sqlmock.NewRows([]string{"code", "name"}).
		AddRow("US", "United States").
		AddRow("CA", "Canada")

	mock.ExpectQuery("SELECT code, name FROM country ORDER BY code LIMIT \\? OFFSET \\?").
		WithArgs(limit, offset).
		WillReturnRows(rows)

	res, err := repo.GetMany(t.Context(), limit, offset)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, "US", res[0].Code)
	assert.Equal(t, "CA", res[1].Code)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCountryRepository_Count(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := CountryRepository{Db: db}

	rows := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(250)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM country").WillReturnRows(rows)

	res, err := repo.Count(t.Context())
	assert.NoError(t, err)
	assert.Equal(t, 250, res)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
