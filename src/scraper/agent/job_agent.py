"""
Job agent for orchestrating job search automation.

This module provides a job agent class that orchestrates the scraping and
processing of job listings from company career websites.
"""

import logging
from typing import Dict, List, Optional, Tuple, Any

from langchain.llms import BaseLLM
from langchain_community.chat_models import ChatOpenAI
from langchain.schema import HumanMessage, SystemMessage
from langchain.prompts import PromptTemplate

from src.infrastructure.database.models.company import Company
from src.infrastructure.database.repositories.company import CompanyRepository
from src.infrastructure.database.repositories.job import JobRepository
from src.infrastructure.database.repositories.technology import TechnologyRepository
from src.scraper.browser.manager import BrowserManager
from src.scraper.browser.page_handler import PageHandler
from src.scraper.processors.job_processor import JobProcessor
from src.scraper.processors.technology_extractor import TechnologyExtractor
from src.scraper.config import settings

logger = logging.getLogger(__name__)


class JobAgent:
    """
    Agent for orchestrating job search automation.

    This class orchestrates the scraping and processing of job listings
    from company career websites.
    """

    def __init__(
        self,
        company_repository: CompanyRepository,
        job_repository: JobRepository,
        technology_repository: TechnologyRepository,
        llm: Optional[BaseLLM] = None,
    ):
        """
        Initialize the job agent.

        Args:
            company_repository: Repository for company data
            job_repository: Repository for job data
            technology_repository: Repository for technology data
            llm: Language model to use for AI tasks, or None to use the default
        """
        self.company_repository = company_repository
        self.job_repository = job_repository
        self.technology_repository = technology_repository
        self.llm = llm or ChatOpenAI(temperature=0, model_name="o4-mini")

        # Create processor instances
        self.technology_extractor = TechnologyExtractor(
            technology_repository, llm=self.llm
        )
        self.job_processor = JobProcessor(
            job_repository, technology_repository, self.technology_extractor
        )

        # System prompt for the language model
        self._system_prompt = """
        You are a job data extraction assistant. Your task is to extract structured job data
        from job listings on company career websites.
        
        Focus on extracting the following information:
        - Job title
        - Job description
        - Job requirements
        - Preferred skills
        - Experience level
        - Employment type
        - Location
        - Work mode (Remote, Hybrid, On-site)
        - Application URL
        - Job function/department
        
        Return the data in a structured format that can be parsed as JSON.
        """

        # Prompt template for job data extraction
        self._extract_prompt_template = PromptTemplate(
            input_variables=["job_page_content"],
            template="""
            Extract structured job data from the following job listing page content.
            Return the data in a JSON format with the following fields:
            - title: The job title
            - description: The full job description
            - requirements: The job requirements
            - preferred_skills: Preferred skills (if specified)
            - experience_level: Experience level (e.g., Entry, Mid, Senior)
            - employment_type: Employment type (e.g., Full-time, Part-time, Contract)
            - location: Job location
            - work_mode: Work mode (Remote, Hybrid, On-site)
            - application_url: URL to apply for the job
            - job_function: Job function or department
            
            Job Page Content:
            {job_page_content}
            
            Structured Job Data (JSON):
            """,
        )

    async def run(self) -> Dict[str, Any]:
        """
        Run the job agent to scrape and process job listings.

        Returns:
            Dictionary containing statistics about the run
        """
        stats = {
            "companies_processed": 0,
            "companies_failed": 0,
            "jobs_found": 0,
            "jobs_processed": 0,
            "jobs_updated": 0,
            "jobs_failed": 0,
            "jobs_marked_inactive": 0,
        }

        try:
            # Get active companies
            companies = await self.company_repository.search(active_only=True)
            logger.info(f"Found {len(companies)} active companies")

            # Process each company
            for company in companies:
                try:
                    company_stats = await self._process_company(company)

                    # Update statistics
                    stats["companies_processed"] += 1
                    stats["jobs_found"] += company_stats["jobs_found"]
                    stats["jobs_processed"] += company_stats["jobs_processed"]
                    stats["jobs_updated"] += company_stats["jobs_updated"]
                    stats["jobs_failed"] += company_stats["jobs_failed"]
                    stats["jobs_marked_inactive"] += company_stats[
                        "jobs_marked_inactive"
                    ]
                except Exception as e:
                    logger.error(f"Error processing company {company.name}: {str(e)}")
                    stats["companies_failed"] += 1

            logger.info(f"Job agent run completed: {stats}")
            return stats
        except Exception as e:
            logger.error(f"Error running job agent: {str(e)}")
            return stats

    async def _process_company(self, company: Company) -> Dict[str, int]:
        """
        Process a company's career page to find and process job listings.

        Args:
            company: The company to process

        Returns:
            Dictionary containing statistics about the company processing
        """
        stats = {
            "jobs_found": 0,
            "jobs_processed": 0,
            "jobs_updated": 0,
            "jobs_failed": 0,
            "jobs_marked_inactive": 0,
        }

        logger.info(f"Processing company: {company.name}")

        # Initialize browser
        async with BrowserManager() as browser:
            try:
                # Create a new page
                page = await browser.new_page()
                page_handler = PageHandler(page)

                # Navigate to the company's career page
                success = await browser.goto(page, company.careers_page_url)
                if not success:
                    logger.error(f"Failed to navigate to {company.careers_page_url}")
                    return stats

                # Extract job links
                job_links = await self._extract_job_links(page_handler)
                stats["jobs_found"] = len(job_links)
                logger.info(f"Found {len(job_links)} job links for {company.name}")

                # Process each job
                active_job_ids = []
                for job_link in job_links:
                    try:
                        # Navigate to the job page
                        success = await browser.goto(page, job_link["url"])
                        if not success:
                            logger.warning(
                                f"Failed to navigate to job: {job_link['url']}"
                            )
                            stats["jobs_failed"] += 1
                            continue

                        # Check if the job is relevant (Costa Rica or LATAM)
                        is_relevant, reason = await self._check_job_relevance(
                            page_handler
                        )
                        if not is_relevant:
                            logger.info(f"Skipping job '{job_link['title']}': {reason}")
                            continue

                        # Extract job data
                        job_data = await self._extract_job_data(page_handler, job_link)
                        if not job_data:
                            logger.warning(
                                f"Failed to extract data for job: {job_link['title']}"
                            )
                            stats["jobs_failed"] += 1
                            continue

                        # Process the job
                        job_id = await self.job_processor.process_job(
                            company.id, job_data
                        )
                        if job_id:
                            active_job_ids.append(job_id)
                            stats["jobs_processed"] += 1
                        else:
                            stats["jobs_failed"] += 1
                    except Exception as e:
                        logger.error(
                            f"Error processing job {job_link['title']}: {str(e)}"
                        )
                        stats["jobs_failed"] += 1

                # Mark jobs that are no longer active as inactive
                if active_job_ids:
                    inactive_count = await self.job_processor.mark_inactive_jobs(
                        company.id, active_job_ids
                    )
                    stats["jobs_marked_inactive"] = inactive_count
            except Exception as e:
                logger.error(f"Error processing company {company.name}: {str(e)}")

        return stats

    async def _extract_job_links(
        self, page_handler: PageHandler
    ) -> List[Dict[str, str]]:
        """
        Extract job links from a company's career page.

        Args:
            page_handler: The page handler

        Returns:
            List of dictionaries containing job URLs and titles
        """
        # Common selectors for job links
        selectors = [
            "a[href*='job']",
            "a[href*='career']",
            "a[href*='position']",
            "a[href*='opening']",
            ".job-listing a",
            ".careers-list a",
            ".job-card a",
            ".position-card a",
        ]

        all_links = []
        for selector in selectors:
            links = await page_handler.extract_job_links(selector)
            all_links.extend(links)

        # Remove duplicates
        unique_links = []
        seen_urls = set()
        for link in all_links:
            if link["url"] not in seen_urls:
                seen_urls.add(link["url"])
                unique_links.append(link)

        return unique_links

    async def _check_job_relevance(self, page_handler: PageHandler) -> Tuple[bool, str]:
        """
        Check if a job is relevant (available in Costa Rica or LATAM).

        Args:
            page_handler: The page handler

        Returns:
            Tuple of (is_relevant, reason)
        """
        # Extract page content
        content = await page_handler.extract_page_content()
        main_text = content.get("main_text", "")

        # Check for excluded locations
        for excluded in settings.excluded_locations:
            if excluded.lower() in main_text.lower():
                return False, f"Excluded location: {excluded}"

        # Check for target locations
        for target in settings.target_locations:
            if target.lower() in main_text.lower():
                return True, f"Target location: {target}"

        # If no target location is found, check if the job might be remote
        if "remote" in main_text.lower():
            return True, "Potentially remote job"

        # If no clear indication, assume it might be relevant
        return True, "No location restrictions found"

    async def _extract_job_data(
        self, page_handler: PageHandler, job_link: Dict[str, str]
    ) -> Optional[Dict[str, Any]]:
        """
        Extract job data from a job page.

        Args:
            page_handler: The page handler
            job_link: Dictionary containing job URL and title

        Returns:
            Dictionary containing job data, or None if extraction failed
        """
        try:
            # Extract page content
            content = await page_handler.extract_page_content()

            # Use AI to extract structured job data
            job_data = await self._extract_structured_job_data(content["main_text"])

            # Add the job title from the link if not extracted
            if not job_data.get("title"):
                job_data["title"] = job_link["title"]

            # Add the application URL
            job_data["application_url"] = job_link["url"]

            return job_data
        except Exception as e:
            logger.error(f"Error extracting job data: {str(e)}")
            return None

    async def _extract_structured_job_data(self, text: str) -> Dict[str, Any]:
        """
        Extract structured job data from text using AI.

        Args:
            text: The text to extract data from

        Returns:
            Dictionary containing structured job data
        """
        try:
            # Prepare the prompt
            prompt = self._extract_prompt_template.format(
                job_page_content=text[:4000]
            )  # Limit text length

            # Get the response from the language model
            messages = [
                SystemMessage(content=self._system_prompt),
                HumanMessage(content=prompt),
            ]
            response = self.llm.predict_messages(messages)

            # Parse the response as JSON
            import json

            job_data = json.loads(response.content)

            return job_data
        except Exception as e:
            logger.error(f"Error extracting structured job data: {str(e)}")
            # Return a minimal job data dictionary
            return {
                "title": "",
                "description": text[
                    :1000
                ],  # Use the first 1000 characters as description
                "requirements": "",
            }
