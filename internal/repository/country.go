package repository

import (
	"database/sql"
	"gogin/internal/model"
	"log/slog"
)

type CountryRepositoryInterface interface {
	GetCountryByCode(code string) (model.Country, error)
	GetMany(limit, offset int) ([]model.Country, error)
}

type CountryRepository struct {
	Db *sql.DB
}

func (r CountryRepository) GetCountryByCode(code string) (model.Country, error) {
	var country model.Country
	if err := r.Db.QueryRow("SELECT code, name FROM country WHERE code = ?", code).
		Scan(&country.Code, &country.Name); err != nil {
		return model.Country{}, err
	}
	return country, nil
}

func (r CountryRepository) GetMany(limit, offset int) ([]model.Country, error) {
	rows, err := r.Db.Query("SELECT code, name FROM country LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		if closeErr := rows.Close(); closeErr != nil {
			slog.Error("GetMany", "error", err)
		}
	}(rows)

	var countries []model.Country
	for rows.Next() {
		var country model.Country
		if err := rows.Scan(&country.Code, &country.Name); err != nil {
			return nil, err
		}
		countries = append(countries, country)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return countries, nil
}
