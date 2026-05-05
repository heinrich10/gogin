-- +goose Up
DROP TABLE IF EXISTS person;
DROP TABLE IF EXISTS country;
DROP TABLE IF EXISTS continent;


CREATE TABLE continent
(
    code TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE country
(
    code           TEXT    NOT NULL PRIMARY KEY,
    name           TEXT    NOT NULL,
    phone          INTEGER NOT NULL,
    symbol         TEXT,
    capital        TEXT,
    currency       TEXT,
    continent_code TEXT,
    alpha_3        TEXT,
    updated_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (continent_code) REFERENCES continent (code) ON DELETE SET NULL
);

CREATE TABLE person
(
    id           INTEGER PRIMARY KEY,
    last_name    TEXT,
    first_name   TEXT NOT NULL,
    country_code TEXT,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (country_code) REFERENCES country (code) ON DELETE SET NULL
);

CREATE INDEX idx_person_country_code ON person (country_code);

-- +goose Down
DROP TABLE IF EXISTS person;
DROP TABLE IF EXISTS country;
DROP TABLE IF EXISTS continent;