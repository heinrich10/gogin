# GoGin

Sample API server built with Go and Gin framework.

## Technologies
- Go (1.25+)
- Gin (HTTP Framework)
- SQLite (Database)
- Goose (Migrations)
- Log/slog (Structured Logging)

## Architecture
The project follows a clean, layered architecture to separate concerns and ensure testability:
- **Controllers**: Handle HTTP requests, input binding, and validation.
- **Services**: Contain business logic and orchestrate operations between repositories.
- **Repositories**: Encapsulate data access logic and SQL queries.
- **Models**: Plain Go structs representing the domain entities.

## Prerequisites
- Go 1.25+ installed
- An SQLite database

## Quickstart

1. Clone the repository
2. Download dependencies: `go mod download`
3. Copy the example environment file and adjust settings:
   ```bash
   cp .env.example .env
   ```
4. Set up your database and run migrations
    ```bash
    go run cmd/migrate/main.go up
    ```
5. Start the server:
    ```bash
   go run cmd/main.go
   ```

## Building
To build the application, run:
```bash
go build -o main cmd/main.go
```

To build docker image, run:
```bash
docker build -f build/Dockerfile -t gogin .
``` 

## Testing
To run tests, use:
```bash
go test ./...
``` 

## Project layout
- `cmd/` - application entry points
- `internal/` - internal packages
    - `app/` - application wiring and router setup
    - `config/` - configuration management
    - `controller/` - HTTP handlers
    - `service/` - business logic (Service/Usecase layer)
    - `repository/` - data access layer
    - `model/` - domain models
    - `lib/` - shared libraries (DB connection, migrations)
    - `util/` - utility functions
- `migrations/` - SQL migration files (goose format)
- `test/` - integration tests
- `AGENTS.md` - detailed agent-specific architecture and conventions guide

## Data Model

```mermaid
---
config:
  layout: dagre
  look: handDrawn
title: ERD
---
erDiagram
	person {
		int id PK ""  
		String last_name  ""  
		String first_name  ""  
		String country_code FK ""  
		timestamp updated_at  ""  
		timestamp created_at  ""  
	}
	country {
		String code PK ""  
		String name  ""  
		int phone  ""  
		String symbol  ""  
		String capital  ""  
		String currency  ""  
		String continent_code FK ""  
		String alpha_3  ""  
		timestamp updated_at  ""  
		timestamp created_at  ""  
	}
	continent {
		String code PK ""  
		String name  ""  
	}
	person}|--||country:"residesIn"
	country}|--||continent:"belongsTo"
```
