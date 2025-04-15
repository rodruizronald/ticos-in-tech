"""
Technology extractor for job descriptions.

This module provides a technology extractor class that extracts technology
information from job descriptions using AI.
"""

import logging
import re
from typing import Dict, List, Optional, Set, Tuple

from langchain.llms import BaseLLM
from langchain.chat_models import ChatOpenAI
from langchain.schema import HumanMessage, SystemMessage
from langchain.prompts import PromptTemplate

from src.infrastructure.database.models.technology import Technology
from src.infrastructure.database.repositories.technology import TechnologyRepository

logger = logging.getLogger(__name__)


class TechnologyExtractor:
    """
    Extracts technology information from job descriptions using AI.

    This class uses LangChain and a language model to extract technology
    information from job descriptions and match them to the technology
    database.
    """

    def __init__(
        self, technology_repository: TechnologyRepository, llm: Optional[BaseLLM] = None
    ):
        """
        Initialize the technology extractor.

        Args:
            technology_repository: Repository for technology data
            llm: Language model to use for extraction, or None to use the default
        """
        self.technology_repository = technology_repository
        self.llm = llm or ChatOpenAI(temperature=0, model_name="gpt-3.5-turbo")

        # Cache of technology names to IDs for faster lookups
        self._technology_cache: Dict[str, int] = {}

        # System prompt for the language model
        self._system_prompt = """
        You are a technology extraction assistant. Your task is to extract technology names, 
        programming languages, frameworks, and tools mentioned in job descriptions.
        
        Focus on extracting specific technologies, not general concepts. For example:
        - Extract "Python", "JavaScript", "React", "Docker", "AWS", "PostgreSQL"
        - Do not extract general terms like "programming", "database", "cloud", "web development"
        
        Return the technologies as a comma-separated list, with no additional text or explanation.
        If no technologies are found, return an empty string.
        """

        # Prompt template for technology extraction
        self._extract_prompt_template = PromptTemplate(
            input_variables=["job_description"],
            template="""
            Extract all technology names, programming languages, frameworks, and tools from the following job description.
            Return them as a comma-separated list, with no additional text or explanation.
            
            Job Description:
            {job_description}
            
            Technologies:
            """,
        )

        # Prompt template for primary technology identification
        self._primary_prompt_template = PromptTemplate(
            input_variables=["technologies", "job_title", "job_description"],
            template="""
            From the following list of technologies mentioned in a job description, identify the PRIMARY technology 
            that is most central to the job role. This should be the main technology the job requires.
            
            Job Title: {job_title}
            
            Job Description Summary: {job_description}
            
            Technologies mentioned: {technologies}
            
            Primary Technology:
            """,
        )

    async def load_technology_cache(self) -> None:
        """
        Load the technology cache from the database.

        This method loads all technologies from the database into a cache
        for faster lookups.
        """
        logger.info("Loading technology cache")
        technologies = await self.technology_repository.get_multi(limit=10000)

        # Create a case-insensitive mapping of technology names to IDs
        self._technology_cache = {tech.name.lower(): tech.id for tech in technologies}
        logger.info(f"Loaded {len(self._technology_cache)} technologies into cache")

    async def extract_technologies(
        self, job_title: str, job_description: str
    ) -> List[Tuple[int, bool]]:
        """
        Extract technologies from a job description and match them to the database.

        Args:
            job_title: The job title
            job_description: The job description

        Returns:
            List of tuples containing (technology_id, is_primary)
        """
        # Ensure the technology cache is loaded
        if not self._technology_cache:
            await self.load_technology_cache()

        # Extract technology names from the job description
        tech_names = await self._extract_technology_names(job_description)
        if not tech_names:
            logger.warning(f"No technologies found in job: {job_title}")
            return []

        # Match technology names to the database
        matched_techs = self._match_technologies(tech_names)
        if not matched_techs:
            logger.warning(f"No technologies matched in job: {job_title}")
            return []

        # Identify the primary technology
        primary_tech_id = await self._identify_primary_technology(
            matched_techs, job_title, job_description
        )

        # Create the result list with is_primary flag
        result = [(tech_id, tech_id == primary_tech_id) for tech_id in matched_techs]

        logger.info(f"Extracted {len(result)} technologies from job: {job_title}")
        return result

    async def _extract_technology_names(self, job_description: str) -> List[str]:
        """
        Extract technology names from a job description using AI.

        Args:
            job_description: The job description

        Returns:
            List of technology names
        """
        try:
            # Prepare the prompt
            prompt = self._extract_prompt_template.format(
                job_description=job_description
            )

            # Get the response from the language model
            messages = [
                SystemMessage(content=self._system_prompt),
                HumanMessage(content=prompt),
            ]
            response = self.llm.predict_messages(messages)

            # Parse the response
            tech_text = response.content.strip()
            if not tech_text:
                return []

            # Split the comma-separated list and clean up each item
            tech_names = [name.strip() for name in tech_text.split(",") if name.strip()]

            return tech_names
        except Exception as e:
            logger.error(f"Error extracting technologies: {str(e)}")
            return []

    def _match_technologies(self, tech_names: List[str]) -> List[int]:
        """
        Match technology names to technology IDs in the database.

        Args:
            tech_names: List of technology names

        Returns:
            List of technology IDs
        """
        matched_ids = set()

        for name in tech_names:
            # Try exact match first (case-insensitive)
            tech_id = self._technology_cache.get(name.lower())
            if tech_id:
                matched_ids.add(tech_id)
                continue

            # Try partial matches
            for db_name, db_id in self._technology_cache.items():
                # Check if the extracted name is a substring of a database name
                # or if a database name is a substring of the extracted name
                if name.lower() in db_name or db_name in name.lower():
                    matched_ids.add(db_id)
                    break

        return list(matched_ids)

    async def _identify_primary_technology(
        self, tech_ids: List[int], job_title: str, job_description: str
    ) -> Optional[int]:
        """
        Identify the primary technology for a job.

        Args:
            tech_ids: List of technology IDs
            job_title: The job title
            job_description: The job description

        Returns:
            The primary technology ID, or None if no primary technology is identified
        """
        if not tech_ids:
            return None

        # If there's only one technology, it's the primary one
        if len(tech_ids) == 1:
            return tech_ids[0]

        try:
            # Get the technology names
            tech_names = []
            for tech_id in tech_ids:
                # Reverse lookup from ID to name
                for name, id_ in self._technology_cache.items():
                    if id_ == tech_id:
                        tech_names.append(name.capitalize())
                        break

            # Prepare the prompt
            prompt = self._primary_prompt_template.format(
                technologies=", ".join(tech_names),
                job_title=job_title,
                job_description=job_description[:500],  # Limit description length
            )

            # Get the response from the language model
            messages = [
                SystemMessage(content=self._system_prompt),
                HumanMessage(content=prompt),
            ]
            response = self.llm.predict_messages(messages)

            # Parse the response
            primary_tech = response.content.strip()

            # Match the primary technology to an ID
            for name, id_ in self._technology_cache.items():
                if (
                    name.lower() == primary_tech.lower()
                    or name.lower() in primary_tech.lower()
                ):
                    return id_

            # If no match is found, return the first technology ID
            return tech_ids[0]
        except Exception as e:
            logger.error(f"Error identifying primary technology: {str(e)}")
            return tech_ids[0]  # Default to the first technology
