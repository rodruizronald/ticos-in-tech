"""
Job data processors.

This module provides processors for extracting and processing job data
from career websites.
"""

from src.scraper.processors.job_processor import JobProcessor
from src.scraper.processors.technology_extractor import TechnologyExtractor

__all__ = ["JobProcessor", "TechnologyExtractor"]
