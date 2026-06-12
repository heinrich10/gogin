# AGENTS Guide

## Scope and sources
- This guide is derived from the current codebase and one convention-source scan: `README.md` was found; no existing `AGENT*.md`/Copilot/Claude/Cursor/Windsurf/Cline rule files were present.
- Prefer code as source of truth if docs ever conflict.

## Big picture architecture
- App is a single Gin HTTP service with layered structure: `cmd/server/main.go` calls `internal/app.NewRouter(db)` to wire `controller` -> `service` -> `repository` -> SQLite (`internal/lib/db.go`).
- Domain split is by resource: `continent`, `country`, `person` in `internal/{controller,service,repository,model}`.
- Controllers depend on service interfaces (`internal/controller/*.go`); services depend on repository interfaces (`internal/service/*.go`). This enables mock-based unit tests at each layer.
- Repository layer uses `database/sql` with SQLite driver `modernc.org/sqlite` and positional `?` placeholders.
- `person` writes are asynchronous: `POST /persons` enqueues to `UpdatePersonChan` via `PersonService`; `PersonService.StartWorker()` performs DB insert.

## Request/data flow patterns
- Router injects standard security headers (HSTS, CSP, X-Frame-Options, etc.) to all responses.
- List endpoints (`GET /continents`, `/countries`, `/persons`) always call `util.Paginate(c)` and pass `(limit, offset)` to `GetMany`.
- Single-resource endpoints map `sql.ErrNoRows` to `404`; other errors return `500` with `{"error":"Something went wrong"}`.
- `POST /persons` returns `202 {"status":"queued"}` after JSON bind; malformed JSON returns `400`.
- Logging uses `log/slog` in all layers with structured key/value pairs; follow existing key names like `"func"`, `"ip"`, `"err"`.
- **Privacy:** Never log PII (Personally Identifiable Information) such as names or phone numbers.

## Config, DB, and migrations
- Runtime config auto-loads from env via `godotenv/autoload` (`internal/config/config.go`); key vars: `DB_HOST`, `HOST`, `PORT`, `TRUSTED_PROXIES`.
- `.env.example` default DB path is `./gogin.db`; DB connection currently sets `SetMaxOpenConns(1)` for SQLite safety.
- Schema/migration source is `migrations/` (goose format: `-- +goose Up/Down`), with full snapshot in `schema.sql`.
- Migrations are run via the native Go CLI: `go run ./cmd/migrate/main.go up` (uses `pressly/goose/v3`, no Docker required).
- Seed migration (`migrations/20260115134736_seed_data.sql`) inserts large country dataset + sample persons; keep new seeds idempotency expectations in mind.

## Developer workflows (verified)
- Run tests: `go test ./...` (verified passing on 2026-06-12).
- Run service from repo root: `go run ./cmd/server/main.go`.
- Build binary: `go build -o main ./cmd/server/main.go`.
- Apply migrations: `go run ./cmd/migrate/main.go up` (native Go, no Docker required).
- Manual API testing: open `test/http/api.http` in a JetBrains IDE to execute requests against the running server.

## Project-specific coding/testing conventions
- Add repository and service interfaces first, then wire concrete structs in `internal/app/app.go`; controllers should depend on service interfaces, and services on repository interfaces.
- Controller tests use Gin test context + `testify/mock` service doubles (see `internal/controller/person_test.go`).
- Service tests use `testify/mock` repository doubles (see `internal/service/person_test.go`).
- Repository tests use `DATA-DOG/go-sqlmock` with explicit SQL expectation strings (see `internal/repository/*_test.go`).
- API integration tests use `httptest` with the real `app.NewRouter(db)` against a temp-file SQLite DB running goose migrations (see `test/api_test.go`).
- Pagination contract: default `limit=10`, `page=1`, max limit `100`; invalid/negative query values fall back to defaults (`internal/util/pagination.go`).
- Keep HTTP JSON field names aligned with model tags in `internal/model/*.go` when adding fields/endpoints.

## Change checklist for new endpoints/features
- Add/extend model in `internal/model`.
- Add repository interface method + SQL implementation + sqlmock tests.
- Add service interface method + implementation + repository-mock tests.
- Add controller handler using existing error/response conventions + service-mock controller tests.
- Wire route in `internal/app/app.go`; if async behavior is needed, follow the person channel/worker pattern.
- Update migrations and, when schema changes, refresh `schema.sql` snapshot.

