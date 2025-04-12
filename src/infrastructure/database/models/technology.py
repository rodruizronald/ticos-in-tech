"""
Technology model for the job board database.

This module defines the SQLAlchemy ORM model for the technologies table.
"""

from datetime import datetime
from typing import List, Optional, TYPE_CHECKING

from sqlalchemy import String, ForeignKey
from sqlalchemy.orm import Mapped, mapped_column, relationship

from src.infrastructure.database.base import Base

if TYPE_CHECKING:
    from src.infrastructure.database.models.job_technology import JobTechnology


class Technology(Base):
    """
    Technology model representing skills and technologies used in jobs.

    Attributes:
        id: Unique identifier for the technology
        name: Technology name (unique)
        category: Category of the technology (e.g., "Programming Language", "Framework")
        parent_id: ID of the parent technology (for hierarchical relationships)
        parent: Relationship to the parent technology
        children: Relationship to child technologies
        created_at: Timestamp when the record was created
        job_technologies: Relationship to job-technology associations
    """

    id: Mapped[int] = mapped_column(primary_key=True)
    name: Mapped[str] = mapped_column(String(100), nullable=False, unique=True)
    category: Mapped[Optional[str]] = mapped_column(String(50), nullable=True)
    parent_id: Mapped[Optional[int]] = mapped_column(
        ForeignKey("technology.id"), nullable=True
    )
    created_at: Mapped[datetime] = mapped_column(
        default=datetime.utcnow, nullable=False
    )

    # Self-referential relationship for technology hierarchy
    parent: Mapped[Optional["Technology"]] = relationship(
        "Technology", remote_side=[id], back_populates="children"
    )
    children: Mapped[List["Technology"]] = relationship(
        "Technology", back_populates="parent"
    )

    # Relationship to job technologies
    job_technologies: Mapped[List["JobTechnology"]] = relationship(
        "JobTechnology", back_populates="technology", cascade="all, delete-orphan"
    )

    def __repr__(self) -> str:
        """String representation of the Technology model."""
        return f"<Technology(id={self.id}, name='{self.name}', category='{self.category}')>"
