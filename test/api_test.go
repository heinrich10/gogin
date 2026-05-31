package api_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gogin/internal/app"
	"gogin/internal/config"
	"gogin/internal/lib"
	"gogin/internal/model"

	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "test.db"))
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	require.NoError(t, db.Ping())

	require.NoError(t, goose.SetDialect("sqlite3"))
	migrationsDir, err := lib.MigrationsDir()
	require.NoError(t, err)
	require.NoError(t, goose.RunContext(t.Context(), "up", db, migrationsDir))

	return db
}

func TestAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupTestDB(t)
	defer db.Close()

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(t.Context())
	// Use t.Cleanup for more reliable cleanup
	t.Cleanup(func() {
		cancel()
		wg.Wait()
	})

	cfg := config.LoadConfig()
	router, updateChan := app.NewRouter(ctx, &wg, db, cfg)
	_ = updateChan // explicitly use or ignore to satisfy compiler
	// No defer close(updateChan) here, we will manage it per subtest if needed or let the worker drain on ctx.Done()

	t.Run("continents", func(t *testing.T) {
		t.Run("list default pagination", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/continents/", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			var result []model.Continent
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
			assert.Len(t, result, 7)
		})

		t.Run("list with limit", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/continents/?limit=3", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			var result []model.Continent
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
			assert.Len(t, result, 3)
			assert.Equal(t, "AF", result[0].Code)
			assert.Equal(t, "AN", result[1].Code)
			assert.Equal(t, "AS", result[2].Code)
		})

		t.Run("list with limit and page", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/continents/?limit=3&page=2", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			var result []model.Continent
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
			assert.Len(t, result, 3)
			assert.Equal(t, "EU", result[0].Code)
			assert.Equal(t, "NA", result[1].Code)
			assert.Equal(t, "OC", result[2].Code)
		})

		t.Run("list limit over max", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/continents/?limit=200", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			var result []model.Continent
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
			assert.Len(t, result, 7)
		})

		t.Run("get one success", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/continents/AF", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			var result model.Continent
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
			assert.Equal(t, "AF", result.Code)
			assert.Equal(t, "Africa", result.Name)
		})

		t.Run("get one not found", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/continents/XX", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNotFound, w.Code)
			var result map[string]string
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
			assert.Equal(t, "Continent not found", result["error"])
		})
	})

	t.Run("countries", func(t *testing.T) {
		t.Run("list default pagination", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/countries/", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			var result []model.Country
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
			assert.Len(t, result, 10) // default limit
		})

		t.Run("list with pagination", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/countries/?limit=2&page=1", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			var result []model.Country
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
			assert.Len(t, result, 2)
		})

		t.Run("get one success", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/countries/US", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			var result model.Country
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
			assert.Equal(t, "US", result.Code)
			assert.Equal(t, "United States", result.Name)
		})

		t.Run("get one not found", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/countries/XX", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNotFound, w.Code)
			var result map[string]string
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
			assert.Equal(t, "Country not found", result["error"])
		})
	})

	t.Run("persons read-only", func(t *testing.T) {
		t.Run("list default pagination", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/persons/", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			var result []model.Person
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
			assert.Len(t, result, 10) // default limit
		})

		t.Run("list with pagination", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/persons/?limit=1&page=2", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			var result []model.Person
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
			assert.Len(t, result, 1)
		})

		t.Run("get one success", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/persons/1", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			var result model.Person
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
			assert.Equal(t, int64(1), result.Id)
			assert.Equal(t, "John", result.FirstName)
			assert.Equal(t, "Doe", result.LastName)
			assert.Equal(t, "US", result.CountryCode)
		})

		t.Run("get one not found", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/persons/999", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNotFound, w.Code)
			var result map[string]string
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
			assert.Equal(t, "Person not found", result["error"])
		})
	})

	t.Run("person create", func(t *testing.T) {
		t.Run("valid person returns accepted", func(t *testing.T) {
			body, _ := json.Marshal(model.Person{
				FirstName:   "Integration",
				LastName:    "Test",
				CountryCode: "JP",
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/persons/", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusAccepted, w.Code)
			var result map[string]string
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
			assert.Equal(t, "queued", result["status"])
		})

		t.Run("async worker persists person", func(t *testing.T) {
			body, _ := json.Marshal(model.Person{
				FirstName:   "AsyncWorker",
				LastName:    "Isolation",
				CountryCode: "JP",
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/persons/", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)
			require.Equal(t, http.StatusAccepted, w.Code)

			// Drain the worker via context cancellation
			cancel()
			wg.Wait()

			// After worker finishes, check the DB
			w = httptest.NewRecorder()
			req, _ = http.NewRequest(http.MethodGet, "/persons/?limit=100", nil)
			router.ServeHTTP(w, req)

			var persons []model.Person
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &persons))
			var found bool
			for _, p := range persons {
				if p.FirstName == "AsyncWorker" && p.LastName == "Isolation" {
					found = true
					break
				}
			}
			assert.True(t, found, "async worker should insert the new person")
		})

		t.Run("invalid json returns bad request", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/persons/", bytes.NewReader([]byte("not-json")))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			var result map[string]string
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
			assert.Contains(t, result["error"], "Invalid request body")
		})

		t.Run("missing required field returns bad request", func(t *testing.T) {
			body, _ := json.Marshal(model.Person{
				LastName:    "MissingFirst",
				CountryCode: "JP",
			})
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/persons/", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			var result map[string]any
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
			assert.Contains(t, result["error"], "validation failed")
		})

		t.Run("invalid country code length returns bad request", func(t *testing.T) {
			body, _ := json.Marshal(model.Person{
				FirstName:   "Invalid",
				LastName:    "Country",
				CountryCode: "USA",
			})
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/persons/", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			var result map[string]any
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
			assert.Contains(t, result["error"], "validation failed")
		})
	})
}
