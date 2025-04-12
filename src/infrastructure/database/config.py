"""
Database configuration module.

This module provides configuration for database connections,
including connection string generation and session management.
"""

import os
from typing import AsyncGenerator, Generator, Optional

from sqlalchemy import create_engine
from sqlalchemy.ext.asyncio import (
    AsyncEngine,
    AsyncSession,
    async_sessionmaker,
    create_async_engine,
)
from sqlalchemy.orm import Session, sessionmaker

# Default connection parameters
DEFAULT_DB_HOST = "localhost"
DEFAULT_DB_PORT = "5432"
DEFAULT_DB_USER = "postgres"
DEFAULT_DB_PASSWORD = "postgres"
DEFAULT_DB_NAME = "jobboard"


def get_postgres_uri() -> str:
    """
    Generate a PostgreSQL connection URI from environment variables.

    Returns:
        str: PostgreSQL connection URI
    """
    host = os.environ.get("DB_HOST", DEFAULT_DB_HOST)
    port = os.environ.get("DB_PORT", DEFAULT_DB_PORT)
    user = os.environ.get("DB_USER", DEFAULT_DB_USER)
    password = os.environ.get("DB_PASSWORD", DEFAULT_DB_PASSWORD)
    db_name = os.environ.get("DB_NAME", DEFAULT_DB_NAME)

    return f"postgresql://{user}:{password}@{host}:{port}/{db_name}"


def get_postgres_async_uri() -> str:
    """
    Generate an asynchronous PostgreSQL connection URI from environment variables.

    Returns:
        str: Asynchronous PostgreSQL connection URI
    """
    host = os.environ.get("DB_HOST", DEFAULT_DB_HOST)
    port = os.environ.get("DB_PORT", DEFAULT_DB_PORT)
    user = os.environ.get("DB_USER", DEFAULT_DB_USER)
    password = os.environ.get("DB_PASSWORD", DEFAULT_DB_PASSWORD)
    db_name = os.environ.get("DB_NAME", DEFAULT_DB_NAME)

    return f"postgresql+asyncpg://{user}:{password}@{host}:{port}/{db_name}"


# Synchronous engine and session factory
engine = create_engine(
    get_postgres_uri(),
    pool_pre_ping=True,
    pool_size=10,
    max_overflow=20,
    pool_recycle=3600,
)
SessionFactory = sessionmaker(autocommit=False, autoflush=False, bind=engine)


# Asynchronous engine and session factory
async_engine = create_async_engine(
    get_postgres_async_uri(),
    pool_pre_ping=True,
    pool_size=10,
    max_overflow=20,
    pool_recycle=3600,
)
AsyncSessionFactory = async_sessionmaker(
    async_engine, autocommit=False, autoflush=False, expire_on_commit=False
)


def get_session() -> Generator[Session, None, None]:
    """
    Get a synchronous database session.

    Yields:
        Session: SQLAlchemy session
    """
    session = SessionFactory()
    try:
        yield session
    finally:
        session.close()


async def get_async_session() -> AsyncGenerator[AsyncSession, None]:
    """
    Get an asynchronous database session.

    Yields:
        AsyncSession: SQLAlchemy asynchronous session
    """
    session = AsyncSessionFactory()
    try:
        yield session
    finally:
        await session.close()
