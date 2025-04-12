"""
Company model for the job board database.

This module defines the SQLAlchemy ORM model for the companies table.
"""

from datetime import datetime
from typing import List, Optional, TYPE_CHECKING

from sqlalchemy import String, Text, Boolean
from sqlalchemy.orm import Mapped, mapped_column, relationship

from src.infrastructure.database.base import Base, TimestampMixin

if TYPE_CHECKING:
    from src.infrastructure.database.models.job import Job


class Company(Base, TimestampMixin):
    """
    Company model representing employers posting jobs.

    Attributes:
        id: Unique identifier for the company
        name: Company name
        careers_page_url: URL to the company's careers page
        logo_url: URL to the company's logo image
        description: Detailed description of the company
        industry: Industry category the company belongs to
        active: Whether the company is active in the system
        created_at: Timestamp when the record was created
        updated_at: Timestamp when the record was last updated
        jobs: Relationship to jobs posted by this company
    """

    id: Mapped[int] = mapped_column(primary_key=True)
    name: Mapped[str] = mapped_column(String(255), nullable=False)
    careers_page_url: Mapped[str] = mapped_column(String(255), nullable=False)
    logo_url: Mapped[Optional[str]] = mapped_column(String(255), nullable=True)
    description: Mapped[Optional[str]] = mapped_column(Text, nullable=True)
    industry: Mapped[Optional[str]] = mapped_column(String(100), nullable=True)
    active: Mapped[bool] = mapped_column(Boolean, default=True, nullable=False)

    # Relationships
    jobs: Mapped[List["Job"]] = relationship(
        "Job", back_populates="company", cascade="all, delete-orphan"
    )

    def __repr__(self) -> str:
        """String representation of the Company model."""
        return (
            f"<Company(id={self.id}, name='{self.name}', industry='{self.industry}')>"
        )
