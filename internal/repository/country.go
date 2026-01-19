package repository

import (
	"database/sql"
	"gogin/internal/model"
	"log/slog"
)

type CountryRepository struct {
	Db *sql.DB
}

func (r CountryRepository) GetCountryByCode(code string) (model.Continent, error) {
	var continent model.Continent
	if err := r.Db.QueryRow("SELECT code, name FROM country WHERE code = ?", code).
		Scan(&continent.Code, &continent.Name); err != nil {
		return model.Continent{}, err
	}
	return continent, nil
}

func (r CountryRepository) GetMany(limit, offset int) ([]model.Continent, error) {
	rows, err := r.Db.Query("SELECT code, name FROM country LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			slog.Error("func", "GetMany", "close", err)
		}
	}(rows)

	var continents []model.Continent
	for rows.Next() {
		var continent model.Continent
		if err := rows.Scan(&continent.Code, &continent.Name); err != nil {
			return nil, err
		}
		continents = append(continents, continent)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return continents, nil
}
