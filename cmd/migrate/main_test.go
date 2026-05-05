package main

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gogin/internal/lib"

	_ "modernc.org/sqlite"
)

func TestMigrations_UpAndDown(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := sql.Open("sqlite", dbPath)
	require.NoError(t, err)
	defer db.Close()

	require.NoError(t, goose.SetDialect("sqlite3"))

	dir, err := lib.MigrationsDir()
	require.NoError(t, err)
	ctx := context.Background()

	// Run all up migrations
	require.NoError(t, goose.RunContext(ctx, "up", db, dir))

	// Verify schema and seed data
	var continentCount int
	require.NoError(t, db.QueryRow("SELECT COUNT(*) FROM continent").Scan(&continentCount))
	assert.Equal(t, 7, continentCount)

	var countryCount int
	require.NoError(t, db.QueryRow("SELECT COUNT(*) FROM country").Scan(&countryCount))
	assert.GreaterOrEqual(t, countryCount, 200)

	var personCount int
	require.NoError(t, db.QueryRow("SELECT COUNT(*) FROM person").Scan(&personCount))
	assert.Equal(t, 13, personCount)

	// Rollback seed data migration
	require.NoError(t, goose.RunContext(ctx, "down", db, dir))

	// Verify tables exist but are empty
	require.NoError(t, db.QueryRow("SELECT COUNT(*) FROM continent").Scan(&continentCount))
	assert.Zero(t, continentCount)

	require.NoError(t, db.QueryRow("SELECT COUNT(*) FROM country").Scan(&countryCount))
	assert.Zero(t, countryCount)

	require.NoError(t, db.QueryRow("SELECT COUNT(*) FROM person").Scan(&personCount))
	assert.Zero(t, personCount)

	// Rollback schema migration
	require.NoError(t, goose.RunContext(ctx, "down", db, dir))

	// Verify tables are dropped
	var tableCount int
	require.NoError(t, db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name IN ('continent', 'country', 'person')").Scan(&tableCount))
	assert.Zero(t, tableCount)
}
