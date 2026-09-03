# OpenMat Backend

OpenMat is a Go backend application using PostgreSQL and Redis, with Docker Compose providing the development environment.

## Requirements

- Go
- Docker
- Docker Compose
- [Task](https://taskfile.dev/)
- Tern

## Getting Started

Clone the repository and enter the backend directory:

```bash
cd apps/backend
```

Create your local environment file:

```bash
cp .env.example .env
```

Fill in the required database and application configuration in `.env`.

### Start the development environment

Start OpenMat and its dependencies with:

```bash
task dev
```

This runs:

```bash
docker compose up --build
```

The development environment includes:

- OpenMat
- PostgreSQL
- Redis

### Run the application locally

To run the Go application directly:

```bash
task run
```

This runs:

```bash
go run ./cmd/openmat/main.go
```

## Docker Commands

### Stop the development environment

```bash
task down
```

### Stop the environment and remove volumes

```bash
task down:volumes
```

> **Warning:** `task down:volumes` removes Docker volumes, including the PostgreSQL data volume. Any local database data stored in that volume will be deleted.

### View application logs

```bash
task logs
```

## Database Migrations

OpenMat uses [Tern](https://github.com/jackc/tern) for database migrations.

### Create a migration

```bash
task migrations:new name=create_users
```

### Apply migrations

```bash
task migrations:up
```

### Roll back the latest migration

```bash
task migrations:down
```

This requires confirmation before running.

### Roll back to a specific version

```bash
task migrations:down:to version=1
```

This also requires confirmation.

## Development Utilities

### Format and verify dependencies

```bash
task tidy
```

This will:

1. Format the Go source code
2. Run `go mod tidy`
3. Verify Go module dependencies

### View available tasks

```bash
task help
```

or:

```bash
task --list-all
```

## Project Structure

```text
apps/backend/
├── cmd/
│   └── openmat/
│       └── main.go
├── internal/
│   └── database/
│       └── migrations/
├── .env.example
├── .env
├── .dockerignore
├── compose.yml
├── Dockerfile
├── Taskfile.yml
├── go.mod
└── go.sum
```

## Configuration

OpenMat uses environment variables for configuration.

Nested configuration uses a double underscore (`__`) separator:

```env
OPENMAT_DATABASE__HOST=localhost
OPENMAT_DATABASE__PORT=5431
OPENMAT_REDIS__ADDRESS=redis://localhost:6380
```

These are converted into nested Koanf configuration keys:

```text
database.host
database.port
redis.address
```

Docker uses container-specific values such as:

```env
OPENMAT_DATABASE__HOST=postgres
OPENMAT_DATABASE__PORT=5432
OPENMAT_REDIS__ADDRESS=redis://redis:6379
```

Local `.env` files containing secrets should not be committed to the repository. Use `.env.example` to document the required variables without exposing credentials.

## Development Workflow

A typical development workflow is:

```bash
# Start the development environment
task dev

# Run database migrations
task migrations:up

# View application logs
task logs

# Stop the environment
task down
```

For a clean local database reset:

```bash
task down:volumes
task dev
task migrations:up
```

## Taskfile

Common commands are exposed through `Taskfile.yml` so development commands remain consistent and easy to discover.

Run:

```bash
task help
```

to see all available tasks.
