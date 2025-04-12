"""
Base repository classes for database access.

This module provides base repository classes that define common CRUD operations
for database access. Both synchronous and asynchronous versions are provided.
"""

from abc import ABC, abstractmethod
from typing import (
    Any,
    Dict,
    Generic,
    List,
    Optional,
    Sequence,
    Type,
    TypeVar,
    Union,
    cast,
)

from sqlalchemy import Select, asc, desc, func, select
from sqlalchemy.exc import IntegrityError, SQLAlchemyError
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.orm import Session

from src.infrastructure.database.base import Base, ModelType

# Type variable for model instances
T = TypeVar("T", bound=Base)


class BaseRepository(Generic[T], ABC):
    """
    Base repository class for synchronous database operations.

    This class provides common CRUD operations for database access.
    It should be extended by concrete repository classes.

    Attributes:
        model: The SQLAlchemy model class
        session: The SQLAlchemy session
    """

    def __init__(self, model: Type[T], session: Session):
        """
        Initialize the repository with a model class and session.

        Args:
            model: The SQLAlchemy model class
            session: The SQLAlchemy session
        """
        self.model = model
        self.session = session

    def create(self, obj_in: Dict[str, Any]) -> T:
        """
        Create a new record in the database.

        Args:
            obj_in: Dictionary containing the data for the new record

        Returns:
            The created model instance

        Raises:
            IntegrityError: If there's a constraint violation
            SQLAlchemyError: If there's a database error
        """
        try:
            db_obj = self.model(**obj_in)
            self.session.add(db_obj)
            self.session.commit()
            self.session.refresh(db_obj)
            return db_obj
        except IntegrityError as e:
            self.session.rollback()
            raise IntegrityError(
                f"Failed to create {self.model.__name__}: {str(e)}", e.params, e.orig
            )
        except SQLAlchemyError as e:
            self.session.rollback()
            raise SQLAlchemyError(f"Failed to create {self.model.__name__}: {str(e)}")

    def get(self, id: Any) -> Optional[T]:
        """
        Get a record by ID.

        Args:
            id: The ID of the record to get

        Returns:
            The model instance if found, None otherwise
        """
        return self.session.get(self.model, id)

    def get_multi(
        self,
        *,
        skip: int = 0,
        limit: int = 100,
        filters: Optional[Dict[str, Any]] = None,
        sort_by: Optional[str] = None,
        sort_desc: bool = False,
    ) -> List[T]:
        """
        Get multiple records with pagination, filtering, and sorting.

        Args:
            skip: Number of records to skip (for pagination)
            limit: Maximum number of records to return
            filters: Dictionary of field-value pairs to filter by
            sort_by: Field to sort by
            sort_desc: Whether to sort in descending order

        Returns:
            List of model instances
        """
        query = select(self.model)

        # Apply filters
        if filters:
            for field, value in filters.items():
                if hasattr(self.model, field):
                    query = query.where(getattr(self.model, field) == value)

        # Apply sorting
        if sort_by and hasattr(self.model, sort_by):
            if sort_desc:
                query = query.order_by(desc(getattr(self.model, sort_by)))
            else:
                query = query.order_by(asc(getattr(self.model, sort_by)))

        # Apply pagination
        query = query.offset(skip).limit(limit)

        return list(self.session.scalars(query).all())

    def count(self, filters: Optional[Dict[str, Any]] = None) -> int:
        """
        Count the number of records matching the given filters.

        Args:
            filters: Dictionary of field-value pairs to filter by

        Returns:
            Number of matching records
        """
        query = select(func.count()).select_from(self.model)

        # Apply filters
        if filters:
            for field, value in filters.items():
                if hasattr(self.model, field):
                    query = query.where(getattr(self.model, field) == value)

        return cast(int, self.session.scalar(query))

    def update(self, id: Any, obj_in: Dict[str, Any]) -> Optional[T]:
        """
        Update a record by ID.

        Args:
            id: The ID of the record to update
            obj_in: Dictionary containing the data to update

        Returns:
            The updated model instance if found, None otherwise

        Raises:
            IntegrityError: If there's a constraint violation
            SQLAlchemyError: If there's a database error
        """
        try:
            db_obj = self.get(id)
            if db_obj is None:
                return None

            for field, value in obj_in.items():
                if hasattr(db_obj, field):
                    setattr(db_obj, field, value)

            self.session.add(db_obj)
            self.session.commit()
            self.session.refresh(db_obj)
            return db_obj
        except IntegrityError as e:
            self.session.rollback()
            raise IntegrityError(
                f"Failed to update {self.model.__name__}: {str(e)}", e.params, e.orig
            )
        except SQLAlchemyError as e:
            self.session.rollback()
            raise SQLAlchemyError(f"Failed to update {self.model.__name__}: {str(e)}")

    def delete(self, id: Any) -> bool:
        """
        Delete a record by ID.

        Args:
            id: The ID of the record to delete

        Returns:
            True if the record was deleted, False otherwise

        Raises:
            SQLAlchemyError: If there's a database error
        """
        try:
            db_obj = self.get(id)
            if db_obj is None:
                return False

            self.session.delete(db_obj)
            self.session.commit()
            return True
        except SQLAlchemyError as e:
            self.session.rollback()
            raise SQLAlchemyError(f"Failed to delete {self.model.__name__}: {str(e)}")

    def execute_query(self, query: Select) -> Sequence[T]:
        """
        Execute a custom query.

        Args:
            query: The SQLAlchemy query to execute

        Returns:
            Sequence of model instances
        """
        return self.session.scalars(query).all()


class BaseAsyncRepository(Generic[T], ABC):
    """
    Base repository class for asynchronous database operations.

    This class provides common CRUD operations for database access.
    It should be extended by concrete repository classes.

    Attributes:
        model: The SQLAlchemy model class
        session: The SQLAlchemy async session
    """

    def __init__(self, model: Type[T], session: AsyncSession):
        """
        Initialize the repository with a model class and session.

        Args:
            model: The SQLAlchemy model class
            session: The SQLAlchemy async session
        """
        self.model = model
        self.session = session

    async def create(self, obj_in: Dict[str, Any]) -> T:
        """
        Create a new record in the database.

        Args:
            obj_in: Dictionary containing the data for the new record

        Returns:
            The created model instance

        Raises:
            IntegrityError: If there's a constraint violation
            SQLAlchemyError: If there's a database error
        """
        try:
            db_obj = self.model(**obj_in)
            self.session.add(db_obj)
            await self.session.commit()
            await self.session.refresh(db_obj)
            return db_obj
        except IntegrityError as e:
            await self.session.rollback()
            raise IntegrityError(
                f"Failed to create {self.model.__name__}: {str(e)}", e.params, e.orig
            )
        except SQLAlchemyError as e:
            await self.session.rollback()
            raise SQLAlchemyError(f"Failed to create {self.model.__name__}: {str(e)}")

    async def get(self, id: Any) -> Optional[T]:
        """
        Get a record by ID.

        Args:
            id: The ID of the record to get

        Returns:
            The model instance if found, None otherwise
        """
        return await self.session.get(self.model, id)

    async def get_multi(
        self,
        *,
        skip: int = 0,
        limit: int = 100,
        filters: Optional[Dict[str, Any]] = None,
        sort_by: Optional[str] = None,
        sort_desc: bool = False,
    ) -> List[T]:
        """
        Get multiple records with pagination, filtering, and sorting.

        Args:
            skip: Number of records to skip (for pagination)
            limit: Maximum number of records to return
            filters: Dictionary of field-value pairs to filter by
            sort_by: Field to sort by
            sort_desc: Whether to sort in descending order

        Returns:
            List of model instances
        """
        query = select(self.model)

        # Apply filters
        if filters:
            for field, value in filters.items():
                if hasattr(self.model, field):
                    query = query.where(getattr(self.model, field) == value)

        # Apply sorting
        if sort_by and hasattr(self.model, sort_by):
            if sort_desc:
                query = query.order_by(desc(getattr(self.model, sort_by)))
            else:
                query = query.order_by(asc(getattr(self.model, sort_by)))

        # Apply pagination
        query = query.offset(skip).limit(limit)

        result = await self.session.scalars(query)
        return list(result.all())

    async def count(self, filters: Optional[Dict[str, Any]] = None) -> int:
        """
        Count the number of records matching the given filters.

        Args:
            filters: Dictionary of field-value pairs to filter by

        Returns:
            Number of matching records
        """
        query = select(func.count()).select_from(self.model)

        # Apply filters
        if filters:
            for field, value in filters.items():
                if hasattr(self.model, field):
                    query = query.where(getattr(self.model, field) == value)

        result = await self.session.scalar(query)
        return cast(int, result)

    async def update(self, id: Any, obj_in: Dict[str, Any]) -> Optional[T]:
        """
        Update a record by ID.

        Args:
            id: The ID of the record to update
            obj_in: Dictionary containing the data to update

        Returns:
            The updated model instance if found, None otherwise

        Raises:
            IntegrityError: If there's a constraint violation
            SQLAlchemyError: If there's a database error
        """
        try:
            db_obj = await self.get(id)
            if db_obj is None:
                return None

            for field, value in obj_in.items():
                if hasattr(db_obj, field):
                    setattr(db_obj, field, value)

            self.session.add(db_obj)
            await self.session.commit()
            await self.session.refresh(db_obj)
            return db_obj
        except IntegrityError as e:
            await self.session.rollback()
            raise IntegrityError(
                f"Failed to update {self.model.__name__}: {str(e)}", e.params, e.orig
            )
        except SQLAlchemyError as e:
            await self.session.rollback()
            raise SQLAlchemyError(f"Failed to update {self.model.__name__}: {str(e)}")

    async def delete(self, id: Any) -> bool:
        """
        Delete a record by ID.

        Args:
            id: The ID of the record to delete

        Returns:
            True if the record was deleted, False otherwise

        Raises:
            SQLAlchemyError: If there's a database error
        """
        try:
            db_obj = await self.get(id)
            if db_obj is None:
                return False

            await self.session.delete(db_obj)
            await self.session.commit()
            return True
        except SQLAlchemyError as e:
            await self.session.rollback()
            raise SQLAlchemyError(f"Failed to delete {self.model.__name__}: {str(e)}")

    async def execute_query(self, query: Select) -> Sequence[T]:
        """
        Execute a custom query.

        Args:
            query: The SQLAlchemy query to execute

        Returns:
            Sequence of model instances
        """
        result = await self.session.scalars(query)
        return result.all()
