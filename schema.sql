CREATE TABLE IF NOT EXISTS "schema_migrations" (version varchar(128) primary key);
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
-- Dbmate schema migrations
INSERT INTO "schema_migrations" (version) VALUES
  ('20260115134728'),
  ('20260115134736');
