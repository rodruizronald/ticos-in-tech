"""
Database configuration module.

This module provides configuration for database connections,
including connection string generation and session management.
"""

import os
from typing import AsyncGenerator, Generator, Optional

from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

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


def get_postgres_uri(async_: bool = False) -> str:
    """
    Generate a PostgreSQL connection URI from environment variables.

    Args:
        async_: Whether to generate an asynchronous URI

    Returns:
        str: PostgreSQL connection URI
    """
    # Load environment variables from .env file if it exists
    host = os.getenv("DB_HOST", DEFAULT_DB_HOST)
    port = os.getenv("DB_PORT", DEFAULT_DB_PORT)
    user = os.getenv("DB_USER", DEFAULT_DB_USER)
    password = os.getenv("DB_PASSWORD", DEFAULT_DB_PASSWORD)
    db_name = os.getenv("DB_NAME", DEFAULT_DB_NAME)

    # Use asyncpg driver for async connections
    driver = "postgresql+asyncpg" if async_ else "postgresql"

    return f"{driver}://{user}:{password}@{host}:{port}/{db_name}"


# Synchronous engine and session factory
engine = create_engine(
    get_postgres_uri(async_=False),
    pool_pre_ping=True,
    pool_size=10,
    max_overflow=20,
    pool_recycle=3600,
)
SessionFactory = sessionmaker(autocommit=False, autoflush=False, bind=engine)


# Asynchronous engine and session factory
async_engine = create_async_engine(
    get_postgres_uri(async_=True),
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
