# Ticos in Tech - Job Board Application

A job board application that connects companies with job seekers, focusing on technology skills.

## Project Structure

```
.
├── cmd/                  # Application entry points
├── internal/             # Private application code
│   ├── models/           # Data models
│   └── repository/       # Database access layer
├── migrations/           # Database migration files
└── schema/               # Database schema definition
```

## Database Models

The application uses the following data models:

- **Company**: Represents companies that post jobs
- **Job**: Represents job postings
- **Technology**: Represents technology skills (programming languages, frameworks, etc.)
- **TechnologyAlias**: Represents alternative names for technologies
- **JobTechnology**: Represents the association between jobs and technologies

## Database Setup

### Prerequisites

- PostgreSQL 12+
- [golang-migrate](https://github.com/golang-migrate/migrate) for running migrations

### Installing golang-migrate

```bash
# macOS
brew install golang-migrate

# Linux
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/migrate

# Windows (using scoop)
scoop install migrate
```

### Running Migrations

1. Create a PostgreSQL database:

```bash
createdb ticos_in_tech
```

2. Run the migrations:

```bash
# Apply migrations
migrate -path migrations -database "postgres://username:password@localhost:5432/ticos_in_tech?sslmode=disable" up

# Rollback migrations
migrate -database "postgres://username:password@localhost:5432/ticos_in_tech?sslmode=disable" -path migrations down
```

Replace `username` and `password` with your PostgreSQL credentials.

## Development

### Dependencies

This project uses:

- Go 1.24+
- pgx v5 for PostgreSQL interaction
- golang-migrate for database migrations

### Adding New Migrations

To create a new migration:

```bash
migrate create -ext sql -dir migrations -seq migration_name
```

This will create two files:
- `migrations/{timestamp}_migration_name.up.sql` - For applying the migration
- `migrations/{timestamp}_migration_name.down.sql` - For rolling back the migration

## Usage

The repository layer provides interfaces for interacting with the database:

- `CompanyRepository`: CRUD operations for companies
- `JobRepository`: CRUD operations for jobs
- `TechnologyRepository`: CRUD operations for technologies
- `JobTechnologyRepository`: Operations for job-technology associations

Example usage:

```go
import (
    "context"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/rodruizronald/ticos-in-tech/internal/models"
    "github.com/rodruizronald/ticos-in-tech/internal/repository"
)

func main() {
    // Connect to database
    connStr := "postgres://username:password@localhost:5432/ticos_in_tech"
    pool, err := pgxpool.New(context.Background(), connStr)
    if err != nil {
        panic(err)
    }
    defer pool.Close()

    // Create repositories
    companyRepo := repository.NewCompanyRepository(pool)
    jobRepo := repository.NewJobRepository(pool)
    techRepo := repository.NewTechnologyRepository(pool)
    jobTechRepo := repository.NewJobTechnologyRepository(pool)

    // Use repositories
    company := &models.Company{
        Name: "Example Company",
        Active: true,
    }
    err = companyRepo.Create(context.Background(), company)
    // ...
}
