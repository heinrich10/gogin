package repository

import (
	"context"
	"database/sql"
	"gogin/internal/model"
	"log/slog"
)

type PersonRepositoryInterface interface {
	GetPersonById(ctx context.Context, id string) (model.Person, error)
	GetMany(ctx context.Context, limit, offset int) ([]model.Person, error)
	Create(ctx context.Context, body model.Person) error
}

type PersonRepository struct {
	Db *sql.DB
}

func (r PersonRepository) GetPersonById(ctx context.Context, id string) (model.Person, error) {
	var person model.Person
	if err := r.Db.QueryRowContext(ctx,
		"SELECT id, first_name, last_name, country_code, updated_at, created_at "+
			"FROM person WHERE id = ?", id,
	).Scan(
		&person.Id, &person.FirstName, &person.LastName, &person.CountryCode, &person.UpdatedAt, &person.CreatedAt,
	); err != nil {
		return model.Person{}, err
	}
	return person, nil
}

func (r PersonRepository) GetMany(ctx context.Context, limit, offset int) ([]model.Person, error) {
	rows, err := r.Db.QueryContext(ctx,
		"SELECT id, first_name, last_name, country_code "+
			"FROM person LIMIT ? OFFSET ?", limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		if closeErr := rows.Close(); closeErr != nil {

			slog.Error("GetMany", "error", err)
		}
	}(rows)

	var persons []model.Person
	for rows.Next() {
		var person model.Person
		if err := rows.Scan(&person.Id, &person.FirstName, &person.LastName, &person.CountryCode); err != nil {
			return nil, err
		}
		persons = append(persons, person)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return persons, nil
}

func (r PersonRepository) Create(ctx context.Context, body model.Person) error {
	_, err := r.Db.ExecContext(ctx,
		"INSERT INTO "+
			"person (first_name, last_name, country_code, updated_at, created_at) "+
			"VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)",
		body.FirstName, body.LastName, body.CountryCode,
	)
	return err
}
