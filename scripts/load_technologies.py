"""
Script to load initial technology data into the database.

This script loads a predefined set of technologies into the database,
including programming languages, frameworks, databases, etc.
"""

import asyncio
import logging
import os
from pathlib import Path
import sys

# Add the project root directory to the Python path
project_root = Path(__file__).resolve().parent.parent
sys.path.append(str(project_root))


from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession
from sqlalchemy.orm import sessionmaker

from src.infrastructure.database.config import get_postgres_uri
from src.infrastructure.database.models.technology import Technology
from src.infrastructure.database.repositories.technology import (
    TechnologyAsyncRepository,
)


# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
)
logger = logging.getLogger(__name__)


# Define technology data
TECHNOLOGIES = [
    # Programming Languages
    {"name": "Python", "category": "Programming Language"},
    {"name": "JavaScript", "category": "Programming Language"},
    {"name": "TypeScript", "category": "Programming Language"},
    {"name": "Java", "category": "Programming Language"},
    {"name": "C#", "category": "Programming Language"},
    {"name": "Go", "category": "Programming Language"},
    {"name": "Rust", "category": "Programming Language"},
    {"name": "PHP", "category": "Programming Language"},
    {"name": "Ruby", "category": "Programming Language"},
    {"name": "Swift", "category": "Programming Language"},
    {"name": "Kotlin", "category": "Programming Language"},
    {"name": "C++", "category": "Programming Language"},
    {"name": "C", "category": "Programming Language"},
    # Frontend Frameworks
    {"name": "React", "category": "Frontend Framework", "parent": "JavaScript"},
    {"name": "Angular", "category": "Frontend Framework", "parent": "TypeScript"},
    {"name": "Vue.js", "category": "Frontend Framework", "parent": "JavaScript"},
    {"name": "Svelte", "category": "Frontend Framework", "parent": "JavaScript"},
    {"name": "Next.js", "category": "Frontend Framework", "parent": "React"},
    {"name": "Nuxt.js", "category": "Frontend Framework", "parent": "Vue.js"},
    # Backend Frameworks
    {"name": "Django", "category": "Backend Framework", "parent": "Python"},
    {"name": "Flask", "category": "Backend Framework", "parent": "Python"},
    {"name": "FastAPI", "category": "Backend Framework", "parent": "Python"},
    {"name": "Express", "category": "Backend Framework", "parent": "JavaScript"},
    {"name": "NestJS", "category": "Backend Framework", "parent": "TypeScript"},
    {"name": "Spring Boot", "category": "Backend Framework", "parent": "Java"},
    {"name": "ASP.NET Core", "category": "Backend Framework", "parent": "C#"},
    {"name": "Ruby on Rails", "category": "Backend Framework", "parent": "Ruby"},
    {"name": "Laravel", "category": "Backend Framework", "parent": "PHP"},
    # Databases
    {"name": "PostgreSQL", "category": "Database"},
    {"name": "MySQL", "category": "Database"},
    {"name": "MongoDB", "category": "Database"},
    {"name": "Redis", "category": "Database"},
    {"name": "SQLite", "category": "Database"},
    {"name": "Elasticsearch", "category": "Database"},
    {"name": "Cassandra", "category": "Database"},
    {"name": "DynamoDB", "category": "Database"},
    {"name": "Firebase", "category": "Database"},
    # Cloud Platforms
    {"name": "AWS", "category": "Cloud Platform"},
    {"name": "Google Cloud", "category": "Cloud Platform"},
    {"name": "Azure", "category": "Cloud Platform"},
    {"name": "Heroku", "category": "Cloud Platform"},
    {"name": "DigitalOcean", "category": "Cloud Platform"},
    {"name": "Vercel", "category": "Cloud Platform"},
    {"name": "Netlify", "category": "Cloud Platform"},
    # DevOps
    {"name": "Docker", "category": "DevOps"},
    {"name": "Kubernetes", "category": "DevOps"},
    {"name": "Jenkins", "category": "DevOps"},
    {"name": "GitHub Actions", "category": "DevOps"},
    {"name": "CircleCI", "category": "DevOps"},
    {"name": "Travis CI", "category": "DevOps"},
    {"name": "Terraform", "category": "DevOps"},
    {"name": "Ansible", "category": "DevOps"},
    # Mobile
    {"name": "React Native", "category": "Mobile", "parent": "React"},
    {"name": "Flutter", "category": "Mobile"},
    {"name": "iOS", "category": "Mobile"},
    {"name": "Android", "category": "Mobile"},
    {"name": "Xamarin", "category": "Mobile", "parent": "C#"},
    # Testing
    {"name": "Jest", "category": "Testing", "parent": "JavaScript"},
    {"name": "Pytest", "category": "Testing", "parent": "Python"},
    {"name": "JUnit", "category": "Testing", "parent": "Java"},
    {"name": "Selenium", "category": "Testing"},
    {"name": "Cypress", "category": "Testing"},
    {"name": "Playwright", "category": "Testing"},
    # Other
    {"name": "GraphQL", "category": "API"},
    {"name": "REST", "category": "API"},
    {"name": "WebSockets", "category": "API"},
    {"name": "gRPC", "category": "API"},
    {"name": "Git", "category": "Tool"},
    {"name": "Linux", "category": "Operating System"},
    {"name": "Agile", "category": "Methodology"},
    {"name": "Scrum", "category": "Methodology"},
    {"name": "Kanban", "category": "Methodology"},
]


async def load_technologies() -> None:
    """
    Load technologies into the database.
    """
    # Set up database connection
    engine = create_async_engine(get_postgres_uri(async_=True))
    async_session = sessionmaker(engine, expire_on_commit=False, class_=AsyncSession)

    async with async_session() as session:
        # Create repository
        repository = TechnologyAsyncRepository(session)

        # Create a mapping of technology names to IDs
        tech_map = {}

        # First pass: Create all technologies without parent relationships
        for tech_data in TECHNOLOGIES:
            name = tech_data["name"]
            category = tech_data["category"]

            # Check if technology already exists
            existing_tech = await repository.get_by_name(name)
            if existing_tech:
                logger.info(f"Technology already exists: {name}")
                tech_map[name] = existing_tech.id
                continue

            # Create technology
            tech = await repository.create(
                {
                    "name": name,
                    "category": category,
                    "parent_id": None,  # Will be updated in second pass
                }
            )

            logger.info(f"Created technology: {name} (ID: {tech.id})")
            tech_map[name] = tech.id

        # Second pass: Update parent relationships
        for tech_data in TECHNOLOGIES:
            if "parent" in tech_data:
                name = tech_data["name"]
                parent_name = tech_data["parent"]

                if name in tech_map and parent_name in tech_map:
                    tech_id = tech_map[name]
                    parent_id = tech_map[parent_name]

                    # Update technology with parent ID
                    await repository.update(tech_id, {"parent_id": parent_id})
                    logger.info(f"Updated technology {name} with parent {parent_name}")

        logger.info(f"Loaded {len(tech_map)} technologies into the database")


if __name__ == "__main__":
    try:
        asyncio.run(load_technologies())
        logger.info("Technology loading completed successfully")
    except Exception as e:
        logger.error(f"Error loading technologies: {str(e)}", exc_info=True)
        sys.exit(1)
