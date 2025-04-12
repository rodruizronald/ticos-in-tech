"""
Job model for the job board database.

This module defines the SQLAlchemy ORM model for the jobs table.
"""

from datetime import datetime
from typing import List, Optional, TYPE_CHECKING

from sqlalchemy import String, Text, Boolean, ForeignKey
from sqlalchemy.orm import Mapped, mapped_column, relationship

from src.infrastructure.database.base import Base, TimestampMixin

if TYPE_CHECKING:
    from src.infrastructure.database.models.company import Company
    from src.infrastructure.database.models.job_technology import JobTechnology


class Job(Base, TimestampMixin):
    """
    Job model representing job listings.

    Attributes:
        id: Unique identifier for the job
        company_id: Foreign key to the company posting the job
        title: Job title
        slug: URL-friendly version of the job title
        description: Detailed job description
        requirements: Job requirements
        preferred_skills: Preferred skills for the job (optional)
        experience_level: Required experience level (e.g., "Entry", "Mid", "Senior")
        employment_type: Type of employment (e.g., "Full-time", "Part-time", "Contract")
        location: Job location (optional)
        work_mode: Work mode (e.g., "Remote", "Hybrid", "On-site")
        application_url: URL to apply for the job (optional)
        job_function: Business function of the job (optional)
        first_seen_at: Timestamp when the job was first discovered
        last_seen_at: Timestamp when the job was last seen active
        posted_at: Timestamp when the job was posted (optional)
        is_active: Whether the job is active
        signature: Unique hash to identify duplicate jobs
        created_at: Timestamp when the record was created
        updated_at: Timestamp when the record was last updated
        company: Relationship to the company posting the job
        job_technologies: Relationship to job-technology associations
    """

    id: Mapped[int] = mapped_column(primary_key=True)
    company_id: Mapped[int] = mapped_column(ForeignKey("company.id"))
    title: Mapped[str] = mapped_column(String(255), nullable=False)
    slug: Mapped[str] = mapped_column(String(255), nullable=False)
    description: Mapped[str] = mapped_column(Text, nullable=False)
    requirements: Mapped[str] = mapped_column(Text, nullable=False)
    preferred_skills: Mapped[Optional[str]] = mapped_column(Text, nullable=True)
    experience_level: Mapped[str] = mapped_column(String(50), nullable=False)
    employment_type: Mapped[str] = mapped_column(String(50), nullable=False)
    location: Mapped[Optional[str]] = mapped_column(String(50), nullable=True)
    work_mode: Mapped[str] = mapped_column(String(20), nullable=False)
    application_url: Mapped[Optional[str]] = mapped_column(String(255), nullable=True)
    job_function: Mapped[Optional[str]] = mapped_column(String(100), nullable=True)
    first_seen_at: Mapped[datetime] = mapped_column(
        default=datetime.utcnow, nullable=False
    )
    last_seen_at: Mapped[datetime] = mapped_column(
        default=datetime.utcnow, nullable=False
    )
    posted_at: Mapped[Optional[datetime]] = mapped_column(nullable=True)
    is_active: Mapped[bool] = mapped_column(Boolean, default=True, nullable=False)
    signature: Mapped[str] = mapped_column(String(64), unique=True, nullable=False)

    # Relationships
    company: Mapped["Company"] = relationship("Company", back_populates="jobs")
    job_technologies: Mapped[List["JobTechnology"]] = relationship(
        "JobTechnology", back_populates="job", cascade="all, delete-orphan"
    )

    def __repr__(self) -> str:
        """String representation of the Job model."""
        return (
            f"<Job(id={self.id}, title='{self.title}', company_id={self.company_id})>"
        )
