# AGENTS Guide

## Scope and sources
- This guide is derived from the current codebase and one convention-source scan: `README.md` was found; no existing `AGENT*.md`/Copilot/Claude/Cursor/Windsurf/Cline rule files were present.
- Prefer code as source of truth when docs conflict (example: README starts `cmd/server/main.go`, but entrypoint is `cmd/main.go`).

## Big picture architecture
- App is a single Gin HTTP service with layered structure: `cmd/main.go` wires `controller` -> `repository` -> SQLite (`internal/lib/db.go`).
- Domain split is by resource: `continent`, `country`, `person` in `internal/{controller,repository,model}`.
- Controllers depend on repository interfaces (`internal/controller/*.go`), enabling mock-based unit tests.
- Repository layer uses `database/sql` with SQLite driver `modernc.org/sqlite` and positional `?` placeholders.
- `person` writes are asynchronous: `POST /persons` enqueues to `UpdatePersonChan`; `PersonController.StartWorker()` performs DB insert.

## Request/data flow patterns
- List endpoints (`GET /continents`, `/countries`, `/persons`) always call `util.Paginate(c)` and pass `(limit, offset)` to `GetMany`.
- Single-resource endpoints map `sql.ErrNoRows` to `404`; other errors return `500` with `{"error":"Something went wrong"}`.
- `POST /persons` returns `202 {"status":"queued"}` after JSON bind; malformed JSON returns `400`.
- Logging uses `log/slog` in all layers with structured key/value pairs; follow existing key names like `"func"`, `"ip"`, `"err"`.

## Config, DB, and migrations
- Runtime config auto-loads from env via `godotenv/autoload` (`internal/config/config.go`); key vars: `DB_HOST`, `HOST`, `PORT`, `TRUSTED_PROXIES`.
- `.env.example` default DB path is `./gogin.db`; DB connection currently sets `SetMaxOpenConns(1)` for SQLite safety.
- Schema/migration source is `migrations/` (dbmate format: `-- migrate:up/down`), with full snapshot in `schema.sql`.
- Seed migration (`migrations/20260115134736_seed_data.sql`) inserts large country dataset + sample persons; keep new seeds idempotency expectations in mind.

## Developer workflows (verified)
- Run tests: `go test ./...` (verified passing on 2026-05-03).
- Run service from repo root: `go run ./cmd/main.go`.
- Build binary: `go build -o main ./cmd/main.go`.
- Apply migrations (from README pattern): `dbmate up` against SQLite file in repo root.

## Project-specific coding/testing conventions
- Add repository interfaces first, then inject concrete structs in `cmd/main.go`; controllers should depend on interfaces, not concrete DB structs.
- Controller tests use Gin test context + `testify/mock` repository doubles (see `internal/controller/person_test.go`).
- Repository tests use `DATA-DOG/go-sqlmock` with explicit SQL expectation strings (see `internal/repository/*_test.go`).
- Pagination contract: default `limit=10`, `page=1`, max limit `100`; invalid/negative query values fall back to defaults (`internal/util/pagination.go`).
- Keep HTTP JSON field names aligned with model tags in `internal/model/*.go` when adding fields/endpoints.

## Change checklist for new endpoints/features
- Add/extend model in `internal/model`.
- Add repository interface method + SQL implementation + sqlmock tests.
- Add controller handler using existing error/response conventions + mock-based controller tests.
- Wire route in `cmd/main.go`; if async behavior is needed, follow the person channel/worker pattern.
- Update migrations and, when schema changes, refresh `schema.sql` snapshot.

