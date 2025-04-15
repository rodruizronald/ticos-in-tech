"""
Job processor for job data.

This module provides a job processor class that processes job data
and interacts with the database.
"""

import hashlib
import logging
import re
from datetime import datetime
from typing import Dict, List, Optional, Tuple, Any

from src.infrastructure.database.models.job import Job
from src.infrastructure.database.models.job_technology import JobTechnology
from src.infrastructure.database.repositories.job import JobRepository
from src.infrastructure.database.repositories.technology import TechnologyRepository
from src.scraper.processors.technology_extractor import TechnologyExtractor

logger = logging.getLogger(__name__)


class JobProcessor:
    """
    Processes job data and interacts with the database.

    This class handles job data processing, including location filtering,
    signature generation, and database operations.
    """

    def __init__(
        self,
        job_repository: JobRepository,
        technology_repository: TechnologyRepository,
        technology_extractor: Optional[TechnologyExtractor] = None,
    ):
        """
        Initialize the job processor.

        Args:
            job_repository: Repository for job data
            technology_repository: Repository for technology data
            technology_extractor: Technology extractor, or None to create a new one
        """
        self.job_repository = job_repository
        self.technology_repository = technology_repository
        self.technology_extractor = technology_extractor or TechnologyExtractor(
            technology_repository
        )

    def generate_signature(self, company_id: int, job_title: str) -> str:
        """
        Generate a unique signature for a job.

        Args:
            company_id: The company ID
            job_title: The job title

        Returns:
            A unique signature for the job
        """
        # Normalize the job title by removing extra whitespace and converting to lowercase
        normalized_title = re.sub(r"\s+", " ", job_title.strip().lower())

        # Create a string to hash
        hash_input = f"{company_id}:{normalized_title}"

        # Generate a SHA-256 hash
        signature = hashlib.sha256(hash_input.encode()).hexdigest()

        return signature

    def generate_slug(self, job_title: str) -> str:
        """
        Generate a URL-friendly slug from a job title.

        Args:
            job_title: The job title

        Returns:
            A URL-friendly slug
        """
        # Convert to lowercase
        slug = job_title.lower()

        # Replace non-alphanumeric characters with hyphens
        slug = re.sub(r"[^a-z0-9]+", "-", slug)

        # Remove leading and trailing hyphens
        slug = slug.strip("-")

        # Limit length
        if len(slug) > 100:
            slug = slug[:100]

        return slug

    async def process_job(
        self, company_id: int, job_data: Dict[str, Any]
    ) -> Optional[int]:
        """
        Process a job and store it in the database.

        Args:
            company_id: The company ID
            job_data: The job data

        Returns:
            The job ID if the job was stored, None otherwise
        """
        try:
            # Extract required fields
            title = job_data.get("title")
            description = job_data.get("description")

            if not title or not description:
                logger.warning("Job missing required fields: title or description")
                return None

            # Generate signature and check if job already exists
            signature = self.generate_signature(company_id, title)
            existing_job = await self.job_repository.get_by_signature(signature)

            if existing_job:
                # Update last_seen_at timestamp for existing job
                logger.info(f"Job already exists: {title} (ID: {existing_job.id})")
                updated_job = await self.job_repository.update_last_seen(
                    existing_job.id
                )
                return existing_job.id

            # Generate slug
            slug = self.generate_slug(title)

            # Prepare job data for database
            job_dict = {
                "company_id": company_id,
                "title": title,
                "slug": slug,
                "description": description,
                "requirements": job_data.get("requirements", ""),
                "preferred_skills": job_data.get("preferred_skills"),
                "experience_level": job_data.get("experience_level", "Not specified"),
                "employment_type": job_data.get("employment_type", "Full-time"),
                "location": job_data.get("location"),
                "work_mode": job_data.get("work_mode", "Not specified"),
                "application_url": job_data.get("application_url"),
                "job_function": job_data.get("job_function"),
                "posted_at": job_data.get("posted_at", datetime.utcnow()),
                "signature": signature,
                "is_active": True,
            }

            # Create job in database
            job = await self.job_repository.create(job_dict)
            logger.info(f"Created new job: {title} (ID: {job.id})")

            # Extract and store technologies
            await self._process_technologies(job.id, title, description)

            return job.id
        except Exception as e:
            logger.error(f"Error processing job: {str(e)}")
            return None

    async def _process_technologies(
        self, job_id: int, job_title: str, job_description: str
    ) -> None:
        """
        Extract technologies from a job description and store them in the database.

        Args:
            job_id: The job ID
            job_title: The job title
            job_description: The job description
        """
        try:
            # Extract technologies
            technologies = await self.technology_extractor.extract_technologies(
                job_title, job_description
            )

            if not technologies:
                logger.warning(f"No technologies found for job ID {job_id}")
                return

            # Store technologies
            for tech_id, is_primary in technologies:
                await self.job_repository.add_technology(job_id, tech_id, is_primary)

            logger.info(f"Added {len(technologies)} technologies to job ID {job_id}")
        except Exception as e:
            logger.error(f"Error processing technologies for job ID {job_id}: {str(e)}")

    async def mark_inactive_jobs(
        self, company_id: int, active_job_ids: List[int]
    ) -> int:
        """
        Mark jobs that are no longer active as inactive.

        Args:
            company_id: The company ID
            active_job_ids: List of active job IDs

        Returns:
            Number of jobs marked as inactive
        """
        try:
            # Get all active jobs for the company
            all_active_jobs = await self.job_repository.get_multi(
                filters={"company_id": company_id, "is_active": True}
            )

            # Find jobs that are no longer active
            inactive_job_ids = []
            for job in all_active_jobs:
                if job.id not in active_job_ids:
                    inactive_job_ids.append(job.id)

            # Mark jobs as inactive
            count = 0
            for job_id in inactive_job_ids:
                job_data = {"is_active": False}
                await self.job_repository.update(job_id, job_data)
                count += 1

            if count > 0:
                logger.info(
                    f"Marked {count} jobs as inactive for company ID {company_id}"
                )

            return count
        except Exception as e:
            logger.error(f"Error marking inactive jobs: {str(e)}")
            return 0
