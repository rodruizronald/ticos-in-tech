"""
Configuration settings for the job search automation agent.

This module provides configuration settings for the scraper, browser,
and AI components of the job search automation agent.
"""

from typing import List
from pydantic import Field
from pydantic_settings import BaseSettings


class ScraperSettings(BaseSettings):
    """Settings for the scraper component."""

    # Rate limiting settings
    request_delay: float = Field(
        default=2.0,
        description="Delay between requests in seconds to avoid overloading target websites",
    )

    # Location filtering settings
    target_locations: List[str] = Field(
        default=["Costa Rica", "LATAM", "Latin America", "Remote"],
        description="List of location keywords to look for in job listings",
    )

    excluded_locations: List[str] = Field(
        default=["US Only", "United States Only", "North America Only"],
        description="List of location keywords that would exclude a job",
    )

    # Browser settings
    headless: bool = Field(
        default=True, description="Whether to run the browser in headless mode"
    )

    timeout: int = Field(
        default=30000, description="Timeout for browser operations in milliseconds"
    )

    user_agent: str = Field(
        default="Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
        description="User agent string to use for browser requests",
    )

    # Logging settings
    log_level: str = Field(default="INFO", description="Log level for the scraper")

    # Override settings from environment variables
    class Config:
        env_prefix = "SCRAPER_"
        case_sensitive = False


# Create a singleton instance of the settings
settings = ScraperSettings()
