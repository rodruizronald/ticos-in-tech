"""
JobTechnology model for the job board database.

This module defines the SQLAlchemy ORM model for the job_technologies junction table.
"""

from datetime import datetime
from typing import Optional, TYPE_CHECKING

from sqlalchemy import Boolean, ForeignKey
from sqlalchemy.orm import Mapped, mapped_column, relationship

from src.infrastructure.database.base import Base

if TYPE_CHECKING:
    from src.infrastructure.database.models.job import Job
    from src.infrastructure.database.models.technology import Technology


class JobTechnology(Base):
    """
    JobTechnology model representing the many-to-many relationship between jobs and technologies.

    Attributes:
        id: Unique identifier for the job-technology relationship
        job_id: Foreign key to the job
        technology_id: Foreign key to the technology
        is_primary: Whether this is a primary/key technology for the job
        created_at: Timestamp when the record was created
        job: Relationship to the job
        technology: Relationship to the technology
    """

    id: Mapped[int] = mapped_column(primary_key=True)
    job_id: Mapped[int] = mapped_column(ForeignKey("job.id"))
    technology_id: Mapped[int] = mapped_column(ForeignKey("technology.id"))
    is_primary: Mapped[bool] = mapped_column(Boolean, default=False, nullable=False)
    created_at: Mapped[datetime] = mapped_column(
        default=datetime.utcnow, nullable=False
    )

    # Relationships
    job: Mapped["Job"] = relationship("Job", back_populates="job_technologies")
    technology: Mapped["Technology"] = relationship(
        "Technology", back_populates="job_technologies"
    )

    def __repr__(self) -> str:
        """String representation of the JobTechnology model."""
        return f"<JobTechnology(job_id={self.job_id}, technology_id={self.technology_id}, is_primary={self.is_primary})>"
