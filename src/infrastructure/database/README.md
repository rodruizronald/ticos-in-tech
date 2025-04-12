# Job Board Database Layer

This directory contains the database layer implementation for the job board application. It follows clean architecture principles and provides both synchronous and asynchronous database access.

## Directory Structure

```
database/
├── base.py                # Base classes and utilities for SQLAlchemy models
├── config.py              # Database connection configuration
├── migrations/            # Alembic migration scripts
│   ├── alembic/
│   │   ├── env.py         # Alembic environment configuration
│   │   ├── script.py.mako # Template for migration scripts
│   │   └── versions/      # Migration script versions
│   │       └── 001_initial_schema.py # Initial schema migration
│   └── alembic.ini        # Alembic configuration file
├── models/                # SQLAlchemy ORM models
│   ├── __init__.py        # Model imports
│   ├── company.py         # Company model
│   ├── job.py             # Job model
│   ├── job_technology.py  # JobTechnology junction model
│   └── technology.py      # Technology model
└── repositories/          # Repository pattern implementations
    ├── __init__.py        # Repository imports
    ├── base.py            # Base repository classes
    ├── company.py         # Company repository
    ├── job.py             # Job repository
    └── technology.py      # Technology repository
```

## Components

### Models

The database models are defined using SQLAlchemy ORM and represent the database tables:

- **Company**: Represents employers posting jobs
- **Job**: Represents job listings
- **Technology**: Represents skills and technologies
- **JobTechnology**: Junction table linking jobs to technologies

### Repositories

The repository pattern is implemented to provide a clean interface for database access:

- **BaseRepository**: Base class for synchronous repositories
- **BaseAsyncRepository**: Base class for asynchronous repositories
- **CompanyRepository**: Repository for Company model
- **JobRepository**: Repository for Job model
- **TechnologyRepository**: Repository for Technology model

Each repository provides methods for CRUD operations and additional query methods specific to the entity.

### Migrations

Database migrations are managed using Alembic:

- **001_initial_schema.py**: Initial migration that creates all tables and indexes

## Usage Examples

### Synchronous Usage

```python
from src.infrastructure.database.config import get_session
from src.infrastructure.database.repositories import JobRepository

# Get a database session
with get_session() as session:
    # Create a repository instance
    job_repo = JobRepository(session)
    
    # Use repository methods
    jobs = job_repo.search(
        query="python",
        work_mode="Remote",
        limit=10
    )
    
    # Process results
    for job in jobs:
        print(f"{job.title} at {job.company.name}")
```

### Asynchronous Usage

```python
import asyncio
from src.infrastructure.database.config import get_async_session
from src.infrastructure.database.repositories import JobAsyncRepository

async def get_jobs():
    # Get an async database session
    async with get_async_session() as session:
        # Create a repository instance
        job_repo = JobAsyncRepository(session)
        
        # Use repository methods
        jobs = await job_repo.search(
            query="python",
            work_mode="Remote",
            limit=10
        )
        
        # Process results
        for job in jobs:
            print(f"{job.title} at {job.company.name}")

# Run the async function
asyncio.run(get_jobs())
```

## Database Configuration

Database connection parameters are configured through environment variables:

- `DB_HOST`: Database host (default: "localhost")
- `DB_PORT`: Database port (default: "5432")
- `DB_USER`: Database user (default: "postgres")
- `DB_PASSWORD`: Database password (default: "postgres")
- `DB_NAME`: Database name (default: "jobboard")

## Migrations

To run migrations:

1. Set up the database connection parameters as environment variables
2. Navigate to the migrations directory: `cd src/infrastructure/database/migrations`
3. Run Alembic commands:
   - Initialize the database: `alembic upgrade head`
   - Create a new migration: `alembic revision --autogenerate -m "description"`
   - Apply migrations: `alembic upgrade head`
   - Rollback migrations: `alembic downgrade -1`
