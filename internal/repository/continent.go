package repository

import (
	"context"
	"database/sql"
	"gogin/internal/model"
	"log/slog"
)

type ContinentRepositoryInterface interface {
	GetContinentByCode(ctx context.Context, code string) (model.Continent, error)
	GetMany(ctx context.Context, limit, offset int) ([]model.Continent, error)
}

type ContinentRepository struct {
	Db *sql.DB
}

func (r ContinentRepository) GetContinentByCode(ctx context.Context, code string) (model.Continent, error) {
	var continent model.Continent
	if err := r.Db.QueryRowContext(ctx, "SELECT code, name FROM continent WHERE code = ?", code).
		Scan(&continent.Code, &continent.Name); err != nil {
		return model.Continent{}, err
	}
	return continent, nil
}

func (r ContinentRepository) GetMany(ctx context.Context, limit, offset int) ([]model.Continent, error) {
	rows, err := r.Db.QueryContext(ctx, "SELECT code, name FROM continent ORDER BY code LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		if closeErr := rows.Close(); closeErr != nil {
			slog.Error("GetMany", "error", closeErr)
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
