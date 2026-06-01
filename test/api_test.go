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
	router, updateChan := app.NewRouter(ctx, context.Background(), &wg, db, cfg)
	_ = updateChan // explicitly use or ignore to satisfy compiler
	// No defer close(updateChan) here, we will manage it per subtest if needed or let the worker drain on ctx.Done()

	t.Run("continents", func(t *testing.T) {
		t.Run("list default pagination", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/continents/", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			var result struct {
				Data []model.Continent `json:"data"`
				Meta struct {
					Page       int `json:"page"`
					Limit      int `json:"limit"`
					Total      int `json:"total"`
					TotalPages int `json:"total_pages"`
				} `json:"meta"`
			}
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
			assert.Len(t, result.Data, 7)
			assert.Equal(t, 1, result.Meta.Page)
			assert.Equal(t, 10, result.Meta.Limit)
			assert.Equal(t, 7, result.Meta.Total)
			assert.Equal(t, 1, result.Meta.TotalPages)
		})

		t.Run("list with limit", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/continents/?limit=3", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			var result struct {
				Data []model.Continent `json:"data"`
				Meta struct {
					Page       int `json:"page"`
					Limit      int `json:"limit"`
					Total      int `json:"total"`
					TotalPages int `json:"total_pages"`
				} `json:"meta"`
			}
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
			assert.Len(t, result.Data, 3)
			assert.Equal(t, "AF", result.Data[0].Code)
			assert.Equal(t, "AN", result.Data[1].Code)
			assert.Equal(t, "AS", result.Data[2].Code)
			assert.Equal(t, 1, result.Meta.Page)
			assert.Equal(t, 3, result.Meta.Limit)
			assert.Equal(t, 7, result.Meta.Total)
			assert.Equal(t, 3, result.Meta.TotalPages)
		})

		t.Run("list with limit and page", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/continents/?limit=3&page=2", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			var result struct {
				Data []model.Continent `json:"data"`
				Meta struct {
					Page       int `json:"page"`
					Limit      int `json:"limit"`
					Total      int `json:"total"`
					TotalPages int `json:"total_pages"`
				} `json:"meta"`
			}
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
			assert.Len(t, result.Data, 3)
			assert.Equal(t, "EU", result.Data[0].Code)
			assert.Equal(t, "NA", result.Data[1].Code)
			assert.Equal(t, "OC", result.Data[2].Code)
			assert.Equal(t, 2, result.Meta.Page)
			assert.Equal(t, 3, result.Meta.TotalPages)
		})

		t.Run("list limit over max", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/continents/?limit=200", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			var result struct {
				Data []model.Continent `json:"data"`
				Meta struct {
					Page       int `json:"page"`
					Limit      int `json:"limit"`
					Total      int `json:"total"`
					TotalPages int `json:"total_pages"`
				} `json:"meta"`
			}
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
			assert.Len(t, result.Data, 7)
			assert.Equal(t, 100, result.Meta.Limit)
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
			var result struct {
				Error struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				} `json:"error"`
			}
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
			assert.Equal(t, "ERR_NOT_FOUND", result.Error.Code)
			assert.Equal(t, "Continent not found", result.Error.Message)
		})
	})

	t.Run("countries", func(t *testing.T) {
		t.Run("list default pagination", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/countries/", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			var result struct {
				Data []model.Country `json:"data"`
				Meta struct {
					Page       int `json:"page"`
					Limit      int `json:"limit"`
					Total      int `json:"total"`
					TotalPages int `json:"total_pages"`
				} `json:"meta"`
			}
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
			assert.Len(t, result.Data, 10) // default limit
			assert.Equal(t, 1, result.Meta.Page)
			assert.Equal(t, 10, result.Meta.Limit)
			assert.Greater(t, result.Meta.Total, 0)
		})

		t.Run("list with pagination", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/countries/?limit=2&page=1", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			var result struct {
				Data []model.Country `json:"data"`
				Meta struct {
					Page       int `json:"page"`
					Limit      int `json:"limit"`
					Total      int `json:"total"`
					TotalPages int `json:"total_pages"`
				} `json:"meta"`
			}
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
			assert.Len(t, result.Data, 2)
			assert.Equal(t, 1, result.Meta.Page)
			assert.Equal(t, 2, result.Meta.Limit)
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
			var result struct {
				Error struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				} `json:"error"`
			}
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
			assert.Equal(t, "ERR_NOT_FOUND", result.Error.Code)
			assert.Equal(t, "Country not found", result.Error.Message)
		})
	})

	t.Run("persons read-only", func(t *testing.T) {
		t.Run("list default pagination", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/persons/", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			var result struct {
				Data []model.Person `json:"data"`
				Meta struct {
					Page       int `json:"page"`
					Limit      int `json:"limit"`
					Total      int `json:"total"`
					TotalPages int `json:"total_pages"`
				} `json:"meta"`
			}
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
			assert.Len(t, result.Data, 10) // default limit
			assert.Equal(t, 1, result.Meta.Page)
			assert.Equal(t, 10, result.Meta.Limit)
		})

		t.Run("list with pagination", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/persons/?limit=1&page=2", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			var result struct {
				Data []model.Person `json:"data"`
				Meta struct {
					Page       int `json:"page"`
					Limit      int `json:"limit"`
					Total      int `json:"total"`
					TotalPages int `json:"total_pages"`
				} `json:"meta"`
			}
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
			assert.Len(t, result.Data, 1)
			assert.Equal(t, 2, result.Meta.Page)
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
			var result struct {
				Error struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				} `json:"error"`
			}
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
			assert.Equal(t, "ERR_NOT_FOUND", result.Error.Code)
			assert.Equal(t, "Person not found", result.Error.Message)
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

			var result struct {
				Data []model.Person `json:"data"`
				Meta struct {
					Page       int `json:"page"`
					Limit      int `json:"limit"`
					Total      int `json:"total"`
					TotalPages int `json:"total_pages"`
				} `json:"meta"`
			}
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
			var found bool
			for _, p := range result.Data {
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
			var result struct {
				Error struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				} `json:"error"`
			}
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
			assert.Equal(t, "ERR_BAD_REQUEST", result.Error.Code)
			assert.Contains(t, result.Error.Message, "Invalid request body")
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
			var result struct {
				Error struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				} `json:"error"`
			}
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
			assert.Equal(t, "ERR_BAD_REQUEST", result.Error.Code)
			assert.Contains(t, result.Error.Message, "validation failed")
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
			var result struct {
				Error struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				} `json:"error"`
			}
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
			assert.Equal(t, "ERR_BAD_REQUEST", result.Error.Code)
			assert.Contains(t, result.Error.Message, "validation failed")
		})
	})
}
