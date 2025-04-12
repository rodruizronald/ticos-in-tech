"""
Technology repository implementation.

This module provides repository implementations for the Technology model.
"""

from typing import Dict, List, Optional, Any, Tuple, cast

from sqlalchemy import func, or_, select
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.orm import Session, joinedload, selectinload

from src.infrastructure.database.models.technology import Technology
from src.infrastructure.database.models.job_technology import JobTechnology
from src.infrastructure.database.repositories.base import (
    BaseRepository,
    BaseAsyncRepository,
)


class TechnologyRepository(BaseRepository[Technology]):
    """
    Repository for Technology model with synchronous operations.

    This class extends the BaseRepository to provide specific
    operations for the Technology model.
    """

    def __init__(self, session: Session):
        """
        Initialize the repository with a session.

        Args:
            session: The SQLAlchemy session
        """
        super().__init__(Technology, session)

    def get_by_name(self, name: str) -> Optional[Technology]:
        """
        Get a technology by name.

        Args:
            name: The technology name

        Returns:
            The technology if found, None otherwise
        """
        query = select(Technology).where(Technology.name == name)
        return self.session.scalar(query)

    def get_with_parent(self, id: int) -> Optional[Technology]:
        """
        Get a technology by ID with its parent technology loaded.

        Args:
            id: The technology ID

        Returns:
            The technology with parent data if found, None otherwise
        """
        query = (
            select(Technology)
            .options(joinedload(Technology.parent))
            .where(Technology.id == id)
        )
        return self.session.scalar(query)

    def get_with_children(self, id: int) -> Optional[Technology]:
        """
        Get a technology by ID with its child technologies loaded.

        Args:
            id: The technology ID

        Returns:
            The technology with children data if found, None otherwise
        """
        query = (
            select(Technology)
            .options(selectinload(Technology.children))
            .where(Technology.id == id)
        )
        return self.session.scalar(query)

    def get_with_all_relations(self, id: int) -> Optional[Technology]:
        """
        Get a technology by ID with all related data loaded.

        Args:
            id: The technology ID

        Returns:
            The technology with all related data if found, None otherwise
        """
        query = (
            select(Technology)
            .options(
                joinedload(Technology.parent),
                selectinload(Technology.children),
                selectinload(Technology.job_technologies),
            )
            .where(Technology.id == id)
        )
        return self.session.scalar(query)

    def search(
        self,
        *,
        query: Optional[str] = None,
        category: Optional[str] = None,
        parent_id: Optional[int] = None,
        skip: int = 0,
        limit: int = 100,
        sort_by: str = "name",
        sort_desc: bool = False,
    ) -> List[Technology]:
        """
        Search technologies with various filters.

        Args:
            query: Search term for technology name
            category: Filter by category
            parent_id: Filter by parent technology ID
            skip: Number of records to skip (for pagination)
            limit: Maximum number of records to return
            sort_by: Field to sort by
            sort_desc: Whether to sort in descending order

        Returns:
            List of technologies matching the search criteria
        """
        stmt = select(Technology)

        # Apply filters
        if query:
            stmt = stmt.where(Technology.name.ilike(f"%{query}%"))

        if category:
            stmt = stmt.where(Technology.category == category)

        if parent_id is not None:
            stmt = stmt.where(Technology.parent_id == parent_id)

        # Apply sorting
        if sort_by and hasattr(Technology, sort_by):
            if sort_desc:
                stmt = stmt.order_by(getattr(Technology, sort_by).desc())
            else:
                stmt = stmt.order_by(getattr(Technology, sort_by).asc())

        # Apply pagination
        stmt = stmt.offset(skip).limit(limit)

        return list(self.session.scalars(stmt).all())

    def get_top_level_technologies(self) -> List[Technology]:
        """
        Get all top-level technologies (those without a parent).

        Returns:
            List of top-level technologies
        """
        stmt = select(Technology).where(Technology.parent_id.is_(None))
        return list(self.session.scalars(stmt).all())

    def get_technology_hierarchy(
        self, root_id: Optional[int] = None
    ) -> List[Dict[str, Any]]:
        """
        Get the technology hierarchy as a nested structure.

        Args:
            root_id: The ID of the root technology to start from (None for all top-level)

        Returns:
            List of dictionaries representing the technology hierarchy
        """
        # First, get all technologies
        stmt = select(Technology).order_by(Technology.name)
        all_technologies = list(self.session.scalars(stmt).all())

        # Create a dictionary for quick lookup
        tech_dict = {tech.id: tech for tech in all_technologies}

        # Build the hierarchy
        result = []

        # Helper function to build a tree
        def build_tree(tech_id: Optional[int]) -> List[Dict[str, Any]]:
            tree = []
            for tech in all_technologies:
                if tech.parent_id == tech_id:
                    children = build_tree(tech.id)
                    tree.append(
                        {
                            "id": tech.id,
                            "name": tech.name,
                            "category": tech.category,
                            "children": children,
                        }
                    )
            return tree

        # If root_id is provided, build tree from that node
        if root_id is not None:
            root_tech = tech_dict.get(root_id)
            if root_tech:
                children = build_tree(root_id)
                return [
                    {
                        "id": root_tech.id,
                        "name": root_tech.name,
                        "category": root_tech.category,
                        "children": children,
                    }
                ]
            return []

        # Otherwise, build tree from all top-level technologies
        return build_tree(None)

    def get_categories(self) -> List[str]:
        """
        Get a list of all unique categories.

        Returns:
            List of unique category names
        """
        stmt = (
            select(Technology.category)
            .distinct()
            .where(Technology.category.is_not(None))
        )
        result = self.session.scalars(stmt).all()
        return [category for category in result if category]  # Filter out None values

    def get_technologies_for_job(self, job_id: int) -> List[Tuple[Technology, bool]]:
        """
        Get all technologies associated with a job, along with their primary status.

        Args:
            job_id: The job ID

        Returns:
            List of tuples containing (technology, is_primary)
        """
        stmt = (
            select(Technology, JobTechnology.is_primary)
            .join(JobTechnology, Technology.id == JobTechnology.technology_id)
            .where(JobTechnology.job_id == job_id)
            .order_by(JobTechnology.is_primary.desc(), Technology.name)
        )
        result = self.session.execute(stmt).all()
        return [(tech, is_primary) for tech, is_primary in result]

    def get_popular_technologies(self, limit: int = 10) -> List[Tuple[Technology, int]]:
        """
        Get the most popular technologies based on job count.

        Args:
            limit: Maximum number of technologies to return

        Returns:
            List of tuples containing (technology, job_count)
        """
        stmt = (
            select(Technology, func.count(JobTechnology.job_id).label("job_count"))
            .join(JobTechnology, Technology.id == JobTechnology.technology_id)
            .group_by(Technology.id)
            .order_by(func.count(JobTechnology.job_id).desc())
            .limit(limit)
        )
        result = self.session.execute(stmt).all()
        return [(tech, count) for tech, count in result]


class TechnologyAsyncRepository(BaseAsyncRepository[Technology]):
    """
    Repository for Technology model with asynchronous operations.

    This class extends the BaseAsyncRepository to provide specific
    operations for the Technology model.
    """

    def __init__(self, session: AsyncSession):
        """
        Initialize the repository with a session.

        Args:
            session: The SQLAlchemy async session
        """
        super().__init__(Technology, session)

    async def get_by_name(self, name: str) -> Optional[Technology]:
        """
        Get a technology by name.

        Args:
            name: The technology name

        Returns:
            The technology if found, None otherwise
        """
        query = select(Technology).where(Technology.name == name)
        result = await self.session.scalar(query)
        return result

    async def get_with_parent(self, id: int) -> Optional[Technology]:
        """
        Get a technology by ID with its parent technology loaded.

        Args:
            id: The technology ID

        Returns:
            The technology with parent data if found, None otherwise
        """
        query = (
            select(Technology)
            .options(joinedload(Technology.parent))
            .where(Technology.id == id)
        )
        result = await self.session.scalar(query)
        return result

    async def get_with_children(self, id: int) -> Optional[Technology]:
        """
        Get a technology by ID with its child technologies loaded.

        Args:
            id: The technology ID

        Returns:
            The technology with children data if found, None otherwise
        """
        query = (
            select(Technology)
            .options(selectinload(Technology.children))
            .where(Technology.id == id)
        )
        result = await self.session.scalar(query)
        return result

    async def get_with_all_relations(self, id: int) -> Optional[Technology]:
        """
        Get a technology by ID with all related data loaded.

        Args:
            id: The technology ID

        Returns:
            The technology with all related data if found, None otherwise
        """
        query = (
            select(Technology)
            .options(
                joinedload(Technology.parent),
                selectinload(Technology.children),
                selectinload(Technology.job_technologies),
            )
            .where(Technology.id == id)
        )
        result = await self.session.scalar(query)
        return result

    async def search(
        self,
        *,
        query: Optional[str] = None,
        category: Optional[str] = None,
        parent_id: Optional[int] = None,
        skip: int = 0,
        limit: int = 100,
        sort_by: str = "name",
        sort_desc: bool = False,
    ) -> List[Technology]:
        """
        Search technologies with various filters.

        Args:
            query: Search term for technology name
            category: Filter by category
            parent_id: Filter by parent technology ID
            skip: Number of records to skip (for pagination)
            limit: Maximum number of records to return
            sort_by: Field to sort by
            sort_desc: Whether to sort in descending order

        Returns:
            List of technologies matching the search criteria
        """
        stmt = select(Technology)

        # Apply filters
        if query:
            stmt = stmt.where(Technology.name.ilike(f"%{query}%"))

        if category:
            stmt = stmt.where(Technology.category == category)

        if parent_id is not None:
            stmt = stmt.where(Technology.parent_id == parent_id)

        # Apply sorting
        if sort_by and hasattr(Technology, sort_by):
            if sort_desc:
                stmt = stmt.order_by(getattr(Technology, sort_by).desc())
            else:
                stmt = stmt.order_by(getattr(Technology, sort_by).asc())

        # Apply pagination
        stmt = stmt.offset(skip).limit(limit)

        result = await self.session.scalars(stmt)
        return list(result.all())

    async def get_top_level_technologies(self) -> List[Technology]:
        """
        Get all top-level technologies (those without a parent).

        Returns:
            List of top-level technologies
        """
        stmt = select(Technology).where(Technology.parent_id.is_(None))
        result = await self.session.scalars(stmt)
        return list(result.all())

    async def get_technology_hierarchy(
        self, root_id: Optional[int] = None
    ) -> List[Dict[str, Any]]:
        """
        Get the technology hierarchy as a nested structure.

        Args:
            root_id: The ID of the root technology to start from (None for all top-level)

        Returns:
            List of dictionaries representing the technology hierarchy
        """
        # First, get all technologies
        stmt = select(Technology).order_by(Technology.name)
        result = await self.session.scalars(stmt)
        all_technologies = list(result.all())

        # Create a dictionary for quick lookup
        tech_dict = {tech.id: tech for tech in all_technologies}

        # Build the hierarchy
        result = []

        # Helper function to build a tree
        def build_tree(tech_id: Optional[int]) -> List[Dict[str, Any]]:
            tree = []
            for tech in all_technologies:
                if tech.parent_id == tech_id:
                    children = build_tree(tech.id)
                    tree.append(
                        {
                            "id": tech.id,
                            "name": tech.name,
                            "category": tech.category,
                            "children": children,
                        }
                    )
            return tree

        # If root_id is provided, build tree from that node
        if root_id is not None:
            root_tech = tech_dict.get(root_id)
            if root_tech:
                children = build_tree(root_id)
                return [
                    {
                        "id": root_tech.id,
                        "name": root_tech.name,
                        "category": root_tech.category,
                        "children": children,
                    }
                ]
            return []

        # Otherwise, build tree from all top-level technologies
        return build_tree(None)

    async def get_categories(self) -> List[str]:
        """
        Get a list of all unique categories.

        Returns:
            List of unique category names
        """
        stmt = (
            select(Technology.category)
            .distinct()
            .where(Technology.category.is_not(None))
        )
        result = await self.session.scalars(stmt)
        return [
            category for category in result.all() if category
        ]  # Filter out None values

    async def get_technologies_for_job(
        self, job_id: int
    ) -> List[Tuple[Technology, bool]]:
        """
        Get all technologies associated with a job, along with their primary status.

        Args:
            job_id: The job ID

        Returns:
            List of tuples containing (technology, is_primary)
        """
        stmt = (
            select(Technology, JobTechnology.is_primary)
            .join(JobTechnology, Technology.id == JobTechnology.technology_id)
            .where(JobTechnology.job_id == job_id)
            .order_by(JobTechnology.is_primary.desc(), Technology.name)
        )
        result = await self.session.execute(stmt)
        return [(tech, is_primary) for tech, is_primary in result.all()]

    async def get_popular_technologies(
        self, limit: int = 10
    ) -> List[Tuple[Technology, int]]:
        """
        Get the most popular technologies based on job count.

        Args:
            limit: Maximum number of technologies to return

        Returns:
            List of tuples containing (technology, job_count)
        """
        stmt = (
            select(Technology, func.count(JobTechnology.job_id).label("job_count"))
            .join(JobTechnology, Technology.id == JobTechnology.technology_id)
            .group_by(Technology.id)
            .order_by(func.count(JobTechnology.job_id).desc())
            .limit(limit)
        )
        result = await self.session.execute(stmt)
        return [(tech, count) for tech, count in result.all()]
