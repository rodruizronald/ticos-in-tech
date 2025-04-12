"""
Company repository implementation.

This module provides repository implementations for the Company model.
"""

from typing import Dict, List, Optional, Any, cast

from sqlalchemy import func, or_, select
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.orm import Session

from src.infrastructure.database.models.company import Company
from src.infrastructure.database.repositories.base import (
    BaseRepository,
    BaseAsyncRepository,
)


class CompanyRepository(BaseRepository[Company]):
    """
    Repository for Company model with synchronous operations.

    This class extends the BaseRepository to provide specific
    operations for the Company model.
    """

    def __init__(self, session: Session):
        """
        Initialize the repository with a session.

        Args:
            session: The SQLAlchemy session
        """
        super().__init__(Company, session)

    def get_by_name(self, name: str) -> Optional[Company]:
        """
        Get a company by name.

        Args:
            name: The company name

        Returns:
            The company if found, None otherwise
        """
        query = select(Company).where(Company.name == name)
        return self.session.scalar(query)

    def search(
        self,
        *,
        query: Optional[str] = None,
        industry: Optional[str] = None,
        active_only: bool = True,
        skip: int = 0,
        limit: int = 100,
        sort_by: str = "name",
        sort_desc: bool = False,
    ) -> List[Company]:
        """
        Search companies with various filters.

        Args:
            query: Search term for company name or description
            industry: Filter by industry
            active_only: Whether to include only active companies
            skip: Number of records to skip (for pagination)
            limit: Maximum number of records to return
            sort_by: Field to sort by
            sort_desc: Whether to sort in descending order

        Returns:
            List of companies matching the search criteria
        """
        stmt = select(Company)

        # Apply filters
        if query:
            stmt = stmt.where(
                or_(
                    Company.name.ilike(f"%{query}%"),
                    Company.description.ilike(f"%{query}%"),
                )
            )

        if industry:
            stmt = stmt.where(Company.industry == industry)

        if active_only:
            stmt = stmt.where(Company.active == True)

        # Apply sorting
        if sort_by and hasattr(Company, sort_by):
            if sort_desc:
                stmt = stmt.order_by(getattr(Company, sort_by).desc())
            else:
                stmt = stmt.order_by(getattr(Company, sort_by).asc())

        # Apply pagination
        stmt = stmt.offset(skip).limit(limit)

        return list(self.session.scalars(stmt).all())

    def count_search_results(
        self,
        *,
        query: Optional[str] = None,
        industry: Optional[str] = None,
        active_only: bool = True,
    ) -> int:
        """
        Count the number of companies matching the search criteria.

        Args:
            query: Search term for company name or description
            industry: Filter by industry
            active_only: Whether to include only active companies

        Returns:
            Number of matching companies
        """
        stmt = select(func.count()).select_from(Company)

        # Apply filters
        if query:
            stmt = stmt.where(
                or_(
                    Company.name.ilike(f"%{query}%"),
                    Company.description.ilike(f"%{query}%"),
                )
            )

        if industry:
            stmt = stmt.where(Company.industry == industry)

        if active_only:
            stmt = stmt.where(Company.active == True)

        return cast(int, self.session.scalar(stmt))

    def get_industries(self) -> List[str]:
        """
        Get a list of all unique industries.

        Returns:
            List of unique industry names
        """
        stmt = select(Company.industry).distinct().where(Company.industry.is_not(None))
        result = self.session.scalars(stmt).all()
        return [industry for industry in result if industry]  # Filter out None values

    def toggle_active(self, id: int) -> Optional[Company]:
        """
        Toggle the active status of a company.

        Args:
            id: The company ID

        Returns:
            The updated company if found, None otherwise
        """
        company = self.get(id)
        if company is None:
            return None

        company.active = not company.active
        self.session.add(company)
        self.session.commit()
        self.session.refresh(company)
        return company


class CompanyAsyncRepository(BaseAsyncRepository[Company]):
    """
    Repository for Company model with asynchronous operations.

    This class extends the BaseAsyncRepository to provide specific
    operations for the Company model.
    """

    def __init__(self, session: AsyncSession):
        """
        Initialize the repository with a session.

        Args:
            session: The SQLAlchemy async session
        """
        super().__init__(Company, session)

    async def get_by_name(self, name: str) -> Optional[Company]:
        """
        Get a company by name.

        Args:
            name: The company name

        Returns:
            The company if found, None otherwise
        """
        query = select(Company).where(Company.name == name)
        result = await self.session.scalar(query)
        return result

    async def search(
        self,
        *,
        query: Optional[str] = None,
        industry: Optional[str] = None,
        active_only: bool = True,
        skip: int = 0,
        limit: int = 100,
        sort_by: str = "name",
        sort_desc: bool = False,
    ) -> List[Company]:
        """
        Search companies with various filters.

        Args:
            query: Search term for company name or description
            industry: Filter by industry
            active_only: Whether to include only active companies
            skip: Number of records to skip (for pagination)
            limit: Maximum number of records to return
            sort_by: Field to sort by
            sort_desc: Whether to sort in descending order

        Returns:
            List of companies matching the search criteria
        """
        stmt = select(Company)

        # Apply filters
        if query:
            stmt = stmt.where(
                or_(
                    Company.name.ilike(f"%{query}%"),
                    Company.description.ilike(f"%{query}%"),
                )
            )

        if industry:
            stmt = stmt.where(Company.industry == industry)

        if active_only:
            stmt = stmt.where(Company.active == True)

        # Apply sorting
        if sort_by and hasattr(Company, sort_by):
            if sort_desc:
                stmt = stmt.order_by(getattr(Company, sort_by).desc())
            else:
                stmt = stmt.order_by(getattr(Company, sort_by).asc())

        # Apply pagination
        stmt = stmt.offset(skip).limit(limit)

        result = await self.session.scalars(stmt)
        return list(result.all())

    async def count_search_results(
        self,
        *,
        query: Optional[str] = None,
        industry: Optional[str] = None,
        active_only: bool = True,
    ) -> int:
        """
        Count the number of companies matching the search criteria.

        Args:
            query: Search term for company name or description
            industry: Filter by industry
            active_only: Whether to include only active companies

        Returns:
            Number of matching companies
        """
        stmt = select(func.count()).select_from(Company)

        # Apply filters
        if query:
            stmt = stmt.where(
                or_(
                    Company.name.ilike(f"%{query}%"),
                    Company.description.ilike(f"%{query}%"),
                )
            )

        if industry:
            stmt = stmt.where(Company.industry == industry)

        if active_only:
            stmt = stmt.where(Company.active == True)

        result = await self.session.scalar(stmt)
        return cast(int, result)

    async def get_industries(self) -> List[str]:
        """
        Get a list of all unique industries.

        Returns:
            List of unique industry names
        """
        stmt = select(Company.industry).distinct().where(Company.industry.is_not(None))
        result = await self.session.scalars(stmt)
        return [
            industry for industry in result.all() if industry
        ]  # Filter out None values

    async def toggle_active(self, id: int) -> Optional[Company]:
        """
        Toggle the active status of a company.

        Args:
            id: The company ID

        Returns:
            The updated company if found, None otherwise
        """
        company = await self.get(id)
        if company is None:
            return None

        company.active = not company.active
        self.session.add(company)
        await self.session.commit()
        await self.session.refresh(company)
        return company
