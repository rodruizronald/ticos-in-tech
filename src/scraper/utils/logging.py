"""
Logging utilities for the job search automation agent.

This module provides logging utilities for the job search automation agent.
"""

import logging
import sys
from typing import Optional

from src.scraper.config import settings


def setup_logging(log_level: Optional[str] = None) -> None:
    """
    Set up logging for the job search automation agent.

    Args:
        log_level: The log level to use, or None to use the default from settings
    """
    # Get the log level from settings if not provided
    level = log_level or settings.log_level

    # Convert string log level to numeric value
    numeric_level = getattr(logging, level.upper(), None)
    if not isinstance(numeric_level, int):
        raise ValueError(f"Invalid log level: {level}")

    # Configure the root logger
    logging.basicConfig(
        level=numeric_level,
        format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
        handlers=[
            logging.StreamHandler(sys.stdout),
        ],
    )

    # Set up loggers for external libraries
    logging.getLogger("playwright").setLevel(logging.WARNING)
    logging.getLogger("urllib3").setLevel(logging.WARNING)
    logging.getLogger("asyncio").setLevel(logging.WARNING)

    # Create a logger for the scraper
    logger = logging.getLogger("scraper")
    logger.setLevel(numeric_level)

    logger.info(f"Logging initialized with level: {level}")
