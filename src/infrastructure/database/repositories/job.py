"""
Job repository implementation.

This module provides repository implementations for the Job model.
"""

from datetime import datetime, timedelta
from typing import Dict, List, Optional, Any, Tuple, cast, Union

from sqlalchemy import func, or_, and_, select, text
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.orm import Session, joinedload, selectinload

from src.infrastructure.database.models.company import Company
from src.infrastructure.database.models.job import Job
from src.infrastructure.database.models.job_technology import JobTechnology
from src.infrastructure.database.models.technology import Technology
from src.infrastructure.database.repositories.base import (
    BaseRepository,
    BaseAsyncRepository,
)


class JobRepository(BaseRepository[Job]):
    """
    Repository for Job model with synchronous operations.

    This class extends the BaseRepository to provide specific
    operations for the Job model.
    """

    def __init__(self, session: Session):
        """
        Initialize the repository with a session.

        Args:
            session: The SQLAlchemy session
        """
        super().__init__(Job, session)

    def get_with_company(self, id: int) -> Optional[Job]:
        """
        Get a job by ID with company data loaded.

        Args:
            id: The job ID

        Returns:
            The job with company data if found, None otherwise
        """
        query = select(Job).options(joinedload(Job.company)).where(Job.id == id)
        return self.session.scalar(query)

    def get_with_technologies(self, id: int) -> Optional[Job]:
        """
        Get a job by ID with technologies data loaded.

        Args:
            id: The job ID

        Returns:
            The job with technologies data if found, None otherwise
        """
        query = (
            select(Job)
            .options(
                selectinload(Job.job_technologies).joinedload(JobTechnology.technology)
            )
            .where(Job.id == id)
        )
        return self.session.scalar(query)

    def get_with_all_relations(self, id: int) -> Optional[Job]:
        """
        Get a job by ID with all related data loaded.

        Args:
            id: The job ID

        Returns:
            The job with all related data if found, None otherwise
        """
        query = (
            select(Job)
            .options(
                joinedload(Job.company),
                selectinload(Job.job_technologies).joinedload(JobTechnology.technology),
            )
            .where(Job.id == id)
        )
        return self.session.scalar(query)

    def get_by_signature(self, signature: str) -> Optional[Job]:
        """
        Get a job by its unique signature.

        Args:
            signature: The job signature

        Returns:
            The job if found, None otherwise
        """
        query = select(Job).where(Job.signature == signature)
        return self.session.scalar(query)

    def search(
        self,
        *,
        query: Optional[str] = None,
        company_id: Optional[int] = None,
        location: Optional[str] = None,
        work_mode: Optional[str] = None,
        experience_level: Optional[str] = None,
        employment_type: Optional[str] = None,
        job_function: Optional[str] = None,
        technology_ids: Optional[List[int]] = None,
        posted_after: Optional[datetime] = None,
        active_only: bool = True,
        skip: int = 0,
        limit: int = 100,
        sort_by: str = "posted_at",
        sort_desc: bool = True,
        cursor: Optional[Tuple[datetime, int]] = None,
    ) -> List[Job]:
        """
        Search jobs with various filters.

        Args:
            query: Search term for job title and description
            company_id: Filter by company ID
            location: Filter by location
            work_mode: Filter by work mode (e.g., "Remote", "Hybrid", "On-site")
            experience_level: Filter by experience level
            employment_type: Filter by employment type
            job_function: Filter by job function
            technology_ids: Filter by technology IDs
            posted_after: Filter by posted date
            active_only: Whether to include only active jobs
            skip: Number of records to skip (for pagination)
            limit: Maximum number of records to return
            sort_by: Field to sort by
            sort_desc: Whether to sort in descending order
            cursor: Cursor for pagination (tuple of posted_at and id)

        Returns:
            List of jobs matching the search criteria
        """
        # Start with a query that joins company and loads it
        stmt = (
            select(Job)
            .options(joinedload(Job.company))
            .join(Company, Job.company_id == Company.id)
        )

        # Apply filters
        if query:
            # Use to_tsvector for full-text search if using PostgreSQL
            # This is a simplified version that works with SQLite too
            stmt = stmt.where(
                or_(
                    Job.title.ilike(f"%{query}%"),
                    Job.description.ilike(f"%{query}%"),
                )
            )

        if company_id:
            stmt = stmt.where(Job.company_id == company_id)

        if location:
            stmt = stmt.where(Job.location.ilike(f"%{location}%"))

        if work_mode:
            stmt = stmt.where(Job.work_mode == work_mode)

        if experience_level:
            stmt = stmt.where(Job.experience_level == experience_level)

        if employment_type:
            stmt = stmt.where(Job.employment_type == employment_type)

        if job_function:
            stmt = stmt.where(Job.job_function == job_function)

        if posted_after:
            stmt = stmt.where(Job.posted_at >= posted_after)

        if active_only:
            stmt = stmt.where(Job.is_active == True)
            stmt = stmt.where(Company.active == True)

        # Filter by technologies
        if technology_ids and len(technology_ids) > 0:
            for tech_id in technology_ids:
                # Create a subquery for each technology
                tech_subquery = (
                    select(JobTechnology.job_id)
                    .where(JobTechnology.technology_id == tech_id)
                    .scalar_subquery()
                )
                stmt = stmt.where(Job.id.in_(tech_subquery))

        # Apply cursor-based pagination if cursor is provided
        if cursor:
            posted_at, job_id = cursor
            stmt = stmt.where(
                or_(
                    Job.posted_at < posted_at,
                    and_(Job.posted_at == posted_at, Job.id < job_id),
                )
            )

        # Apply sorting
        if sort_by and hasattr(Job, sort_by):
            if sort_desc:
                stmt = stmt.order_by(getattr(Job, sort_by).desc(), Job.id.desc())
            else:
                stmt = stmt.order_by(getattr(Job, sort_by).asc(), Job.id.asc())
        else:
            # Default sorting by posted_at desc
            stmt = stmt.order_by(Job.posted_at.desc(), Job.id.desc())

        # Apply offset/limit pagination
        if not cursor:
            stmt = stmt.offset(skip)
        stmt = stmt.limit(limit)

        return list(self.session.scalars(stmt).all())

    def count_search_results(
        self,
        *,
        query: Optional[str] = None,
        company_id: Optional[int] = None,
        location: Optional[str] = None,
        work_mode: Optional[str] = None,
        experience_level: Optional[str] = None,
        employment_type: Optional[str] = None,
        job_function: Optional[str] = None,
        technology_ids: Optional[List[int]] = None,
        posted_after: Optional[datetime] = None,
        active_only: bool = True,
    ) -> int:
        """
        Count the number of jobs matching the search criteria.

        Args:
            query: Search term for job title and description
            company_id: Filter by company ID
            location: Filter by location
            work_mode: Filter by work mode
            experience_level: Filter by experience level
            employment_type: Filter by employment type
            job_function: Filter by job function
            technology_ids: Filter by technology IDs
            posted_after: Filter by posted date
            active_only: Whether to include only active jobs

        Returns:
            Number of matching jobs
        """
        stmt = (
            select(func.count(Job.id))
            .select_from(Job)
            .join(Company, Job.company_id == Company.id)
        )

        # Apply filters
        if query:
            stmt = stmt.where(
                or_(
                    Job.title.ilike(f"%{query}%"),
                    Job.description.ilike(f"%{query}%"),
                )
            )

        if company_id:
            stmt = stmt.where(Job.company_id == company_id)

        if location:
            stmt = stmt.where(Job.location.ilike(f"%{location}%"))

        if work_mode:
            stmt = stmt.where(Job.work_mode == work_mode)

        if experience_level:
            stmt = stmt.where(Job.experience_level == experience_level)

        if employment_type:
            stmt = stmt.where(Job.employment_type == employment_type)

        if job_function:
            stmt = stmt.where(Job.job_function == job_function)

        if posted_after:
            stmt = stmt.where(Job.posted_at >= posted_after)

        if active_only:
            stmt = stmt.where(Job.is_active == True)
            stmt = stmt.where(Company.active == True)

        # Filter by technologies
        if technology_ids and len(technology_ids) > 0:
            for tech_id in technology_ids:
                # Create a subquery for each technology
                tech_subquery = (
                    select(JobTechnology.job_id)
                    .where(JobTechnology.technology_id == tech_id)
                    .scalar_subquery()
                )
                stmt = stmt.where(Job.id.in_(tech_subquery))

        return cast(int, self.session.scalar(stmt))

    def get_recent_jobs(
        self, days: int = 30, limit: int = 10, active_only: bool = True
    ) -> List[Job]:
        """
        Get recent jobs posted within the specified number of days.

        Args:
            days: Number of days to look back
            limit: Maximum number of jobs to return
            active_only: Whether to include only active jobs

        Returns:
            List of recent jobs
        """
        cutoff_date = datetime.utcnow() - timedelta(days=days)

        stmt = (
            select(Job)
            .options(joinedload(Job.company))
            .join(Company, Job.company_id == Company.id)
            .where(Job.posted_at >= cutoff_date)
        )

        if active_only:
            stmt = stmt.where(Job.is_active == True)
            stmt = stmt.where(Company.active == True)

        stmt = stmt.order_by(Job.posted_at.desc()).limit(limit)

        return list(self.session.scalars(stmt).all())

    def get_job_counts_by_field(
        self, field: str, active_only: bool = True
    ) -> Dict[str, int]:
        """
        Get job counts grouped by a specific field.

        Args:
            field: Field to group by (e.g., "work_mode", "experience_level")
            active_only: Whether to include only active jobs

        Returns:
            Dictionary mapping field values to job counts
        """
        if not hasattr(Job, field):
            raise ValueError(f"Invalid field: {field}")

        stmt = select(getattr(Job, field), func.count(Job.id)).group_by(
            getattr(Job, field)
        )

        if active_only:
            stmt = stmt.where(Job.is_active == True)

        result = self.session.execute(stmt).all()
        return {value: count for value, count in result if value is not None}

    def update_last_seen(self, id: int) -> Optional[Job]:
        """
        Update the last_seen_at timestamp of a job.

        Args:
            id: The job ID

        Returns:
            The updated job if found, None otherwise
        """
        job = self.get(id)
        if job is None:
            return None

        job.last_seen_at = datetime.utcnow()
        self.session.add(job)
        self.session.commit()
        self.session.refresh(job)
        return job

    def toggle_active(self, id: int) -> Optional[Job]:
        """
        Toggle the active status of a job.

        Args:
            id: The job ID

        Returns:
            The updated job if found, None otherwise
        """
        job = self.get(id)
        if job is None:
            return None

        job.is_active = not job.is_active
        self.session.add(job)
        self.session.commit()
        self.session.refresh(job)
        return job

    def add_technology(
        self, job_id: int, technology_id: int, is_primary: bool = False
    ) -> Optional[JobTechnology]:
        """
        Add a technology to a job.

        Args:
            job_id: The job ID
            technology_id: The technology ID
            is_primary: Whether this is a primary technology for the job

        Returns:
            The created job-technology relationship if successful, None otherwise
        """
        job = self.get(job_id)
        if job is None:
            return None

        # Check if the relationship already exists
        stmt = select(JobTechnology).where(
            JobTechnology.job_id == job_id, JobTechnology.technology_id == technology_id
        )
        existing = self.session.scalar(stmt)

        if existing:
            # Update is_primary if the relationship already exists
            existing.is_primary = is_primary
            self.session.add(existing)
            self.session.commit()
            self.session.refresh(existing)
            return existing

        # Create a new relationship
        job_tech = JobTechnology(
            job_id=job_id, technology_id=technology_id, is_primary=is_primary
        )
        self.session.add(job_tech)
        self.session.commit()
        self.session.refresh(job_tech)
        return job_tech

    def remove_technology(self, job_id: int, technology_id: int) -> bool:
        """
        Remove a technology from a job.

        Args:
            job_id: The job ID
            technology_id: The technology ID

        Returns:
            True if the technology was removed, False otherwise
        """
        stmt = select(JobTechnology).where(
            JobTechnology.job_id == job_id, JobTechnology.technology_id == technology_id
        )
        job_tech = self.session.scalar(stmt)

        if job_tech is None:
            return False

        self.session.delete(job_tech)
        self.session.commit()
        return True


class JobAsyncRepository(BaseAsyncRepository[Job]):
    """
    Repository for Job model with asynchronous operations.

    This class extends the BaseAsyncRepository to provide specific
    operations for the Job model.
    """

    def __init__(self, session: AsyncSession):
        """
        Initialize the repository with a session.

        Args:
            session: The SQLAlchemy async session
        """
        super().__init__(Job, session)

    async def get_with_company(self, id: int) -> Optional[Job]:
        """
        Get a job by ID with company data loaded.

        Args:
            id: The job ID

        Returns:
            The job with company data if found, None otherwise
        """
        query = select(Job).options(joinedload(Job.company)).where(Job.id == id)
        result = await self.session.scalar(query)
        return result

    async def get_with_technologies(self, id: int) -> Optional[Job]:
        """
        Get a job by ID with technologies data loaded.

        Args:
            id: The job ID

        Returns:
            The job with technologies data if found, None otherwise
        """
        query = (
            select(Job)
            .options(
                selectinload(Job.job_technologies).joinedload(JobTechnology.technology)
            )
            .where(Job.id == id)
        )
        result = await self.session.scalar(query)
        return result

    async def get_with_all_relations(self, id: int) -> Optional[Job]:
        """
        Get a job by ID with all related data loaded.

        Args:
            id: The job ID

        Returns:
            The job with all related data if found, None otherwise
        """
        query = (
            select(Job)
            .options(
                joinedload(Job.company),
                selectinload(Job.job_technologies).joinedload(JobTechnology.technology),
            )
            .where(Job.id == id)
        )
        result = await self.session.scalar(query)
        return result

    async def get_by_signature(self, signature: str) -> Optional[Job]:
        """
        Get a job by its unique signature.

        Args:
            signature: The job signature

        Returns:
            The job if found, None otherwise
        """
        query = select(Job).where(Job.signature == signature)
        result = await self.session.scalar(query)
        return result

    async def search(
        self,
        *,
        query: Optional[str] = None,
        company_id: Optional[int] = None,
        location: Optional[str] = None,
        work_mode: Optional[str] = None,
        experience_level: Optional[str] = None,
        employment_type: Optional[str] = None,
        job_function: Optional[str] = None,
        technology_ids: Optional[List[int]] = None,
        posted_after: Optional[datetime] = None,
        active_only: bool = True,
        skip: int = 0,
        limit: int = 100,
        sort_by: str = "posted_at",
        sort_desc: bool = True,
        cursor: Optional[Tuple[datetime, int]] = None,
    ) -> List[Job]:
        """
        Search jobs with various filters.

        Args:
            query: Search term for job title and description
            company_id: Filter by company ID
            location: Filter by location
            work_mode: Filter by work mode (e.g., "Remote", "Hybrid", "On-site")
            experience_level: Filter by experience level
            employment_type: Filter by employment type
            job_function: Filter by job function
            technology_ids: Filter by technology IDs
            posted_after: Filter by posted date
            active_only: Whether to include only active jobs
            skip: Number of records to skip (for pagination)
            limit: Maximum number of records to return
            sort_by: Field to sort by
            sort_desc: Whether to sort in descending order
            cursor: Cursor for pagination (tuple of posted_at and id)

        Returns:
            List of jobs matching the search criteria
        """
        # Start with a query that joins company and loads it
        stmt = (
            select(Job)
            .options(joinedload(Job.company))
            .join(Company, Job.company_id == Company.id)
        )

        # Apply filters
        if query:
            # Use to_tsvector for full-text search if using PostgreSQL
            # This is a simplified version that works with SQLite too
            stmt = stmt.where(
                or_(
                    Job.title.ilike(f"%{query}%"),
                    Job.description.ilike(f"%{query}%"),
                )
            )

        if company_id:
            stmt = stmt.where(Job.company_id == company_id)

        if location:
            stmt = stmt.where(Job.location.ilike(f"%{location}%"))

        if work_mode:
            stmt = stmt.where(Job.work_mode == work_mode)

        if experience_level:
            stmt = stmt.where(Job.experience_level == experience_level)

        if employment_type:
            stmt = stmt.where(Job.employment_type == employment_type)

        if job_function:
            stmt = stmt.where(Job.job_function == job_function)

        if posted_after:
            stmt = stmt.where(Job.posted_at >= posted_after)

        if active_only:
            stmt = stmt.where(Job.is_active == True)
            stmt = stmt.where(Company.active == True)

        # Filter by technologies
        if technology_ids and len(technology_ids) > 0:
            for tech_id in technology_ids:
                # Create a subquery for each technology
                tech_subquery = (
                    select(JobTechnology.job_id)
                    .where(JobTechnology.technology_id == tech_id)
                    .scalar_subquery()
                )
                stmt = stmt.where(Job.id.in_(tech_subquery))

        # Apply cursor-based pagination if cursor is provided
        if cursor:
            posted_at, job_id = cursor
            stmt = stmt.where(
                or_(
                    Job.posted_at < posted_at,
                    and_(Job.posted_at == posted_at, Job.id < job_id),
                )
            )

        # Apply sorting
        if sort_by and hasattr(Job, sort_by):
            if sort_desc:
                stmt = stmt.order_by(getattr(Job, sort_by).desc(), Job.id.desc())
            else:
                stmt = stmt.order_by(getattr(Job, sort_by).asc(), Job.id.asc())
        else:
            # Default sorting by posted_at desc
            stmt = stmt.order_by(Job.posted_at.desc(), Job.id.desc())

        # Apply offset/limit pagination
        if not cursor:
            stmt = stmt.offset(skip)
        stmt = stmt.limit(limit)

        result = await self.session.scalars(stmt)
        return list(result.all())

    async def count_search_results(
        self,
        *,
        query: Optional[str] = None,
        company_id: Optional[int] = None,
        location: Optional[str] = None,
        work_mode: Optional[str] = None,
        experience_level: Optional[str] = None,
        employment_type: Optional[str] = None,
        job_function: Optional[str] = None,
        technology_ids: Optional[List[int]] = None,
        posted_after: Optional[datetime] = None,
        active_only: bool = True,
    ) -> int:
        """
        Count the number of jobs matching the search criteria.

        Args:
            query: Search term for job title and description
            company_id: Filter by company ID
            location: Filter by location
            work_mode: Filter by work mode
            experience_level: Filter by experience level
            employment_type: Filter by employment type
            job_function: Filter by job function
            technology_ids: Filter by technology IDs
            posted_after: Filter by posted date
            active_only: Whether to include only active jobs

        Returns:
            Number of matching jobs
        """
        stmt = (
            select(func.count(Job.id))
            .select_from(Job)
            .join(Company, Job.company_id == Company.id)
        )

        # Apply filters
        if query:
            stmt = stmt.where(
                or_(
                    Job.title.ilike(f"%{query}%"),
                    Job.description.ilike(f"%{query}%"),
                )
            )

        if company_id:
            stmt = stmt.where(Job.company_id == company_id)

        if location:
            stmt = stmt.where(Job.location.ilike(f"%{location}%"))

        if work_mode:
            stmt = stmt.where(Job.work_mode == work_mode)

        if experience_level:
            stmt = stmt.where(Job.experience_level == experience_level)

        if employment_type:
            stmt = stmt.where(Job.employment_type == employment_type)

        if job_function:
            stmt = stmt.where(Job.job_function == job_function)

        if posted_after:
            stmt = stmt.where(Job.posted_at >= posted_after)

        if active_only:
            stmt = stmt.where(Job.is_active == True)
            stmt = stmt.where(Company.active == True)

        # Filter by technologies
        if technology_ids and len(technology_ids) > 0:
            for tech_id in technology_ids:
                # Create a subquery for each technology
                tech_subquery = (
                    select(JobTechnology.job_id)
                    .where(JobTechnology.technology_id == tech_id)
                    .scalar_subquery()
                )
                stmt = stmt.where(Job.id.in_(tech_subquery))

        result = await self.session.scalar(stmt)
        return cast(int, result)

    async def get_recent_jobs(
        self, days: int = 30, limit: int = 10, active_only: bool = True
    ) -> List[Job]:
        """
        Get recent jobs posted within the specified number of days.

        Args:
            days: Number of days to look back
            limit: Maximum number of jobs to return
            active_only: Whether to include only active jobs

        Returns:
            List of recent jobs
        """
        cutoff_date = datetime.utcnow() - timedelta(days=days)

        stmt = (
            select(Job)
            .options(joinedload(Job.company))
            .join(Company, Job.company_id == Company.id)
            .where(Job.posted_at >= cutoff_date)
        )

        if active_only:
            stmt = stmt.where(Job.is_active == True)
            stmt = stmt.where(Company.active == True)

        stmt = stmt.order_by(Job.posted_at.desc()).limit(limit)

        result = await self.session.scalars(stmt)
        return list(result.all())

    async def get_job_counts_by_field(
        self, field: str, active_only: bool = True
    ) -> Dict[str, int]:
        """
        Get job counts grouped by a specific field.

        Args:
            field: Field to group by (e.g., "work_mode", "experience_level")
            active_only: Whether to include only active jobs

        Returns:
            Dictionary mapping field values to job counts
        """
        if not hasattr(Job, field):
            raise ValueError(f"Invalid field: {field}")

        stmt = select(getattr(Job, field), func.count(Job.id)).group_by(
            getattr(Job, field)
        )

        if active_only:
            stmt = stmt.where(Job.is_active == True)

        result = await self.session.execute(stmt)
        return {value: count for value, count in result.all() if value is not None}

    async def update_last_seen(self, id: int) -> Optional[Job]:
        """
        Update the last_seen_at timestamp of a job.

        Args:
            id: The job ID

        Returns:
            The updated job if found, None otherwise
        """
        job = await self.get(id)
        if job is None:
            return None

        job.last_seen_at = datetime.utcnow()
        self.session.add(job)
        await self.session.commit()
        await self.session.refresh(job)
        return job

    async def toggle_active(self, id: int) -> Optional[Job]:
        """
        Toggle the active status of a job.

        Args:
            id: The job ID

        Returns:
            The updated job if found, None otherwise
        """
        job = await self.get(id)
        if job is None:
            return None

        job.is_active = not job.is_active
        self.session.add(job)
        await self.session.commit()
        await self.session.refresh(job)
        return job

    async def add_technology(
        self, job_id: int, technology_id: int, is_primary: bool = False
    ) -> Optional[JobTechnology]:
        """
        Add a technology to a job.

        Args:
            job_id: The job ID
            technology_id: The technology ID
            is_primary: Whether this is a primary technology for the job

        Returns:
            The created job-technology relationship if successful, None otherwise
        """
        job = await self.get(job_id)
        if job is None:
            return None

        # Check if the relationship already exists
        stmt = select(JobTechnology).where(
            JobTechnology.job_id == job_id, JobTechnology.technology_id == technology_id
        )
        existing = await self.session.scalar(stmt)

        if existing:
            # Update is_primary if the relationship already exists
            existing.is_primary = is_primary
            self.session.add(existing)
            await self.session.commit()
            await self.session.refresh(existing)
            return existing

        # Create a new relationship
        job_tech = JobTechnology(
            job_id=job_id, technology_id=technology_id, is_primary=is_primary
        )
        self.session.add(job_tech)
        await self.session.commit()
        await self.session.refresh(job_tech)
        return job_tech

    async def remove_technology(self, job_id: int, technology_id: int) -> bool:
        """
        Remove a technology from a job.

        Args:
            job_id: The job ID
            technology_id: The technology ID

        Returns:
            True if the technology was removed, False otherwise
        """
        stmt = select(JobTechnology).where(
            JobTechnology.job_id == job_id, JobTechnology.technology_id == technology_id
        )
        job_tech = await self.session.scalar(stmt)

        if job_tech is None:
            return False

        await self.session.delete(job_tech)
        await self.session.commit()
        return True
