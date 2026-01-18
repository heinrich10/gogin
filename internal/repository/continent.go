package repository

import (
	"database/sql"
	"gogin/internal/model"
)

type ContinentRepository struct {
	Db *sql.DB
}

func (r ContinentRepository) GetContinentByCode(code string) (model.Continent, error) {
	var continent model.Continent
	err := r.Db.QueryRow("SELECT code, name FROM continent WHERE code = ?", code).Scan(&continent.Code, &continent.Name)
	if err != nil {
		return model.Continent{}, err
	}
	return continent, nil
}

func (r ContinentRepository) GetMany(limit, offset int) ([]model.Continent, error) {
	rows, err := r.Db.Query("SELECT code, name FROM continent LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
