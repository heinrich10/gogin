# Comprehensive API Test Plan

## Evaluation: Is `httptest` Fit for This Project?

**Verdict: Yes. `httptest` is the right tool.**

### Why `httptest` Works Here
- It exercises the full HTTP stack: routing, middleware, JSON binding, status codes, and response bodies.
- It avoids the fragility of subprocess-based E2E tests: no port conflicts, compilation delays, `SIGKILL` orphans, or process lifecycle races.
- The async `POST /persons` worker runs in-process, so `httptest` can verify the end-to-end flow naturally.
- It integrates with `go test` without custom scripts or external tools.

### Caveat
`cmd/main.go` encapsulates all wiring inside `main()`. An `httptest` integration test must recreate that wiring. This is a small, acceptable duplication for a service with only three resource groups.

### Why the Subprocess E2E Approach Failed
- Starting `go run ./cmd/main.go` per test suite is slow and unreliable.
- The test runner killed the server mid-shutdown, leaving orphaned processes.
- SQLite file locking caused race conditions between the test assertions and the worker goroutine.
- Tests were not isolated: `TestPersonCreateSuccess` mutated DB state, causing `TestPersonListDefaultPagination` to fail when it ran later alphabetically.

---

## Implementation Plan

### 1. Clean Up Broken Artifacts
- Delete `test/e2e/api_test.go` (the subprocess approach).
- Remove the empty `test/integration/` directory if it exists.

### 2. Test Strategy
Create **one** file: `test/api_test.go`.

**Database setup**
- Use a **temporary file DB** (`os.CreateTemp`) instead of `:memory:`.
- Call `db.SetMaxOpenConns(1)` immediately after opening to match production behavior and eliminate SQLite connection-pool races.
- Programmatically create schema and seed data using the same DDL inferred from the repository layer.

**Router setup**
- Recreate the Gin router exactly as `cmd/main.go` does:
  - Instantiate real repositories with the test DB.
  - Instantiate real controllers.
  - Start `personController.StartWorker()` in a goroutine.
  - Register routes with `gin.Default()` + `cors.Default()`.
- Return the `updatePersonChan` so the test can close it during cleanup.

**Test structure**
- Use a single top-level `TestAPI` function with `t.Run` subtests.
- Run subtests **sequentially** to control mutation order:
  1. Read-only tests for continents, countries, and persons.
  2. `POST /persons` success (returns `202 Accepted`).
  3. Poll `GET /persons` until the async worker persists the new record.
  4. `POST /persons` with invalid JSON (returns `400 Bad Request`).
- This avoids the isolation bug where a create test pollutes a later list test.

### 3. Coverage Checklist
| Endpoint | Scenario | Assertion |
|---|---|---|
| `GET /continents` | Default pagination | `200 OK`, 7 items |
| `GET /continents?limit=3` | Custom limit | `200 OK`, 3 items |
| `GET /continents?limit=3&page=2` | Pagination page 2 | `200 OK`, 3 items |
| `GET /continents?limit=200` | Limit capped at max | `200 OK`, 7 items |
| `GET /continents/AF` | Existing resource | `200 OK`, code=`AF` |
| `GET /continents/XX` | Missing resource | `404 Not Found` |
| `GET /countries` | Default pagination | `200 OK`, 5 items |
| `GET /countries?limit=2&page=1` | Custom pagination | `200 OK`, 2 items |
| `GET /countries/US` | Existing resource | `200 OK`, code=`US` |
| `GET /countries/XX` | Missing resource | `404 Not Found` |
| `GET /persons` | Default pagination | `200 OK`, 3 items |
| `GET /persons?limit=1&page=2` | Custom pagination | `200 OK`, 1 item |
| `GET /persons/1` | Existing resource | `200 OK`, id=`1` |
| `GET /persons/999` | Missing resource | `404 Not Found` |
| `POST /persons` | Valid body | `202 Accepted`, status=`queued` |
| *(poll)* | Async worker insert | New person appears in list |
| `POST /persons` | Invalid JSON | `400 Bad Request` |

### 4. Files to Change
- **Delete:** `test/e2e/api_test.go`
- **Create:** `test/api_test.go`

### 5. Verification
Run:
```bash
go test ./test/... -v -count=1
```
All subtests should pass in a single run.
