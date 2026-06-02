package repository

import (
	"context"
	"database/sql"
	"gogin/internal/model"
	"log/slog"
)

type CountryRepositoryInterface interface {
	GetCountryByCode(ctx context.Context, code string) (model.Country, error)
	GetMany(ctx context.Context, limit, offset int) ([]model.Country, error)
	Count(ctx context.Context) (int, error)
}

type CountryRepository struct {
	Db *sql.DB
}

func (r CountryRepository) GetCountryByCode(ctx context.Context, code string) (model.Country, error) {
	var country model.Country
	if err := r.Db.QueryRowContext(ctx, "SELECT code, name FROM country WHERE code = ?", code).
		Scan(&country.Code, &country.Name); err != nil {
		return model.Country{}, err
	}
	return country, nil
}

func (r CountryRepository) GetMany(ctx context.Context, limit, offset int) ([]model.Country, error) {
	rows, err := r.Db.QueryContext(ctx, "SELECT code, name FROM country ORDER BY code LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		if closeErr := rows.Close(); closeErr != nil {
			slog.Error("GetMany", "error", closeErr)
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

func (r CountryRepository) Count(ctx context.Context) (int, error) {
	var count int
	if err := r.Db.QueryRowContext(ctx, "SELECT COUNT(*) FROM country").Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
