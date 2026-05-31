package repository

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestContinentRepository_GetContinentByCode(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := ContinentRepository{Db: db}
	code := "AF"
	name := "Africa"

	rows := sqlmock.NewRows([]string{"code", "name"}).
		AddRow(code, name)

	mock.ExpectQuery("SELECT code, name FROM continent WHERE code = ?").
		WithArgs(code).
		WillReturnRows(rows)

	res, err := repo.GetContinentByCode(t.Context(), code)
	assert.NoError(t, err)
	assert.Equal(t, code, res.Code)
	assert.Equal(t, name, res.Name)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestContinentRepository_GetContinentByCode_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := ContinentRepository{Db: db}
	code := "XX"

	mock.ExpectQuery("SELECT code, name FROM continent WHERE code = ?").
		WithArgs(code).
		WillReturnError(sql.ErrNoRows)

	_, err = repo.GetContinentByCode(t.Context(), code)
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
}

func TestContinentRepository_GetMany(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := ContinentRepository{Db: db}
	limit, offset := 10, 0

	rows := sqlmock.NewRows([]string{"code", "name"}).
		AddRow("AF", "Africa").
		AddRow("AN", "Antarctica")

	mock.ExpectQuery("SELECT code, name FROM continent ORDER BY code LIMIT \\? OFFSET \\?").
		WithArgs(limit, offset).
		WillReturnRows(rows)

	res, err := repo.GetMany(t.Context(), limit, offset)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, "AF", res[0].Code)
	assert.Equal(t, "AN", res[1].Code)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
