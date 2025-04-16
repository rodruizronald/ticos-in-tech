"""
Main script for running the job search automation agent.

This script provides a command-line interface for running the job search
automation agent.
"""

import argparse
import asyncio
import logging
import sys
from typing import Dict, Any

from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession
from sqlalchemy.orm import sessionmaker

from src.infrastructure.database.config import get_postgres_uri
from src.infrastructure.database.repositories.company import CompanyAsyncRepository
from src.infrastructure.database.repositories.job import JobAsyncRepository
from src.infrastructure.database.repositories.technology import (
    TechnologyAsyncRepository,
)
from src.scraper.agent.job_agent import JobAgent
from src.scraper.utils.logging import setup_logging

logger = logging.getLogger(__name__)


async def run_agent(args: argparse.Namespace) -> Dict[str, Any]:
    """
    Run the job search automation agent.

    Args:
        args: Command-line arguments

    Returns:
        Dictionary containing statistics about the run
    """
    # Set up database connection
    engine = create_async_engine(get_postgres_uri(async_=True))
    async_session = sessionmaker(engine, expire_on_commit=False, class_=AsyncSession)

    async with async_session() as session:
        # Create repositories
        company_repository = CompanyAsyncRepository(session)
        job_repository = JobAsyncRepository(session)
        technology_repository = TechnologyAsyncRepository(session)

        # Create and run the agent
        agent = JobAgent(
            company_repository=company_repository,
            job_repository=job_repository,
            technology_repository=technology_repository,
        )

        # Run the agent
        stats = await agent.run()

        return stats


def main() -> None:
    """
    Main entry point for the job search automation agent.
    """
    # Parse command-line arguments
    parser = argparse.ArgumentParser(description="Job Search Automation Agent")
    parser.add_argument(
        "--log-level",
        choices=["DEBUG", "INFO", "WARNING", "ERROR", "CRITICAL"],
        default=None,
        help="Set the logging level",
    )
    parser.add_argument(
        "--company",
        type=str,
        help="Process only the specified company (by name)",
    )
    args = parser.parse_args()

    # Set up logging
    setup_logging(args.log_level)

    # Run the agent
    try:
        logger.info("Starting job search automation agent")
        stats = asyncio.run(run_agent(args))
        logger.info(f"Job search automation agent completed: {stats}")
    except KeyboardInterrupt:
        logger.info("Job search automation agent interrupted")
        sys.exit(1)
    except Exception as e:
        logger.error(
            f"Error running job search automation agent: {str(e)}", exc_info=True
        )
        sys.exit(1)


if __name__ == "__main__":
    main()
