"""
Alembic environment configuration.

This module configures the Alembic environment for database migrations.
"""

import os
import sys
from logging.config import fileConfig

from alembic import context
from sqlalchemy import engine_from_config, pool

# Add the project root directory to the Python path
sys.path.insert(
    0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../../../../../"))
)

# Import the SQLAlchemy metadata and database configuration
from src.infrastructure.database.base import metadata
from src.infrastructure.database.config import get_postgres_uri
from src.infrastructure.database.models import *  # Import all models to ensure they are registered with metadata

# This is the Alembic Config object, which provides access to the values within the .ini file
config = context.config

# Interpret the config file for Python logging.
# This line sets up loggers basically.
if config.config_file_name is not None:
    fileConfig(config.config_file_name)

# Set the SQLAlchemy URL in the Alembic configuration
config.set_main_option("sqlalchemy.url", get_postgres_uri())

# Add your model's MetaData object here for 'autogenerate' support
target_metadata = metadata


def run_migrations_offline() -> None:
    """
    Run migrations in 'offline' mode.

    This configures the context with just a URL and not an Engine,
    though an Engine is acceptable here as well. By skipping the Engine creation
    we don't even need a DBAPI to be available.

    Calls to context.execute() here emit the given string to the script output.
    """
    url = config.get_main_option("sqlalchemy.url")
    context.configure(
        url=url,
        target_metadata=target_metadata,
        literal_binds=True,
        dialect_opts={"paramstyle": "named"},
    )

    with context.begin_transaction():
        context.run_migrations()


def run_migrations_online() -> None:
    """
    Run migrations in 'online' mode.

    In this scenario we need to create an Engine and associate a connection with the context.
    """
    connectable = engine_from_config(
        config.get_section(config.config_ini_section, {}),
        prefix="sqlalchemy.",
        poolclass=pool.NullPool,
    )

    with connectable.connect() as connection:
        context.configure(
            connection=connection,
            target_metadata=target_metadata,
            # Compare types between models and database
            compare_type=True,
            # Compare server default values
            compare_server_default=True,
            # Include comments in the migration
            include_object=lambda obj, name, type_, reflected, compare_to: True,
            # Include schemas in the migration
            include_schemas=True,
            # Render as batch migrations for SQLite compatibility
            render_as_batch=True,
        )

        with context.begin_transaction():
            context.run_migrations()


if context.is_offline_mode():
    run_migrations_offline()
else:
    run_migrations_online()
