"""
Page handler for web scraping.

This module provides a page handler class that handles page interactions
and data extraction for web scraping operations.
"""

import logging
import re
from typing import Any, Dict, List, Optional, Tuple, Union

from playwright.async_api import Page, Locator, ElementHandle

from src.scraper.config import settings

logger = logging.getLogger(__name__)


class PageHandler:
    """
    Handles page interactions and data extraction for web scraping.

    This class provides methods for interacting with web pages and extracting
    data from them.
    """

    def __init__(self, page: Page):
        """
        Initialize the page handler.

        Args:
            page: The Playwright page to handle
        """
        self.page = page

    async def extract_text(self, selector: str) -> Optional[str]:
        """
        Extract text from an element.

        Args:
            selector: CSS selector for the element

        Returns:
            The text content of the element, or None if not found
        """
        try:
            element = await self.page.query_selector(selector)
            if not element:
                return None

            return await element.text_content()
        except Exception as e:
            logger.warning(f"Failed to extract text from {selector}: {str(e)}")
            return None

    async def extract_attribute(self, selector: str, attribute: str) -> Optional[str]:
        """
        Extract an attribute from an element.

        Args:
            selector: CSS selector for the element
            attribute: The attribute to extract

        Returns:
            The attribute value, or None if not found
        """
        try:
            element = await self.page.query_selector(selector)
            if not element:
                return None

            return await element.get_attribute(attribute)
        except Exception as e:
            logger.warning(
                f"Failed to extract attribute {attribute} from {selector}: {str(e)}"
            )
            return None

    async def extract_elements(self, selector: str) -> List[ElementHandle]:
        """
        Extract multiple elements matching a selector.

        Args:
            selector: CSS selector for the elements

        Returns:
            List of element handles
        """
        try:
            elements = await self.page.query_selector_all(selector)
            return elements
        except Exception as e:
            logger.warning(
                f"Failed to extract elements with selector {selector}: {str(e)}"
            )
            return []

    async def click(self, selector: str) -> bool:
        """
        Click on an element.

        Args:
            selector: CSS selector for the element

        Returns:
            True if the click was successful, False otherwise
        """
        try:
            await self.page.click(selector)
            return True
        except Exception as e:
            logger.warning(f"Failed to click on {selector}: {str(e)}")
            return False

    async def wait_for_selector(
        self, selector: str, timeout: Optional[int] = None
    ) -> bool:
        """
        Wait for an element to be visible.

        Args:
            selector: CSS selector for the element
            timeout: Timeout in milliseconds, or None to use the default timeout

        Returns:
            True if the element became visible, False otherwise
        """
        try:
            await self.page.wait_for_selector(
                selector, timeout=timeout or settings.timeout
            )
            return True
        except Exception as e:
            logger.warning(f"Timeout waiting for selector {selector}: {str(e)}")
            return False

    async def is_text_visible(self, text: str) -> bool:
        """
        Check if text is visible on the page.

        Args:
            text: The text to check for

        Returns:
            True if the text is visible, False otherwise
        """
        try:
            # Use a more lenient text search that ignores case and whitespace
            content = await self.page.content()
            # Remove HTML tags and normalize whitespace
            content = re.sub(r"<[^>]*>", " ", content)
            content = re.sub(r"\s+", " ", content).lower()

            return text.lower() in content
        except Exception as e:
            logger.warning(f"Failed to check if text '{text}' is visible: {str(e)}")
            return False

    async def extract_job_links(self, link_selector: str) -> List[Dict[str, str]]:
        """
        Extract job links from a careers page.

        Args:
            link_selector: CSS selector for job links

        Returns:
            List of dictionaries containing job URLs and titles
        """
        job_links = []

        try:
            elements = await self.extract_elements(link_selector)

            for element in elements:
                try:
                    url = await element.get_attribute("href")
                    title = await element.text_content()

                    if url and title:
                        # Clean up the title and URL
                        title = title.strip()

                        # Make URL absolute if it's relative
                        if url.startswith("/"):
                            url = f"{self.page.url.rstrip('/')}{url}"

                        job_links.append({"url": url, "title": title})
                except Exception as e:
                    logger.warning(f"Failed to extract job link: {str(e)}")
                    continue
        except Exception as e:
            logger.error(f"Failed to extract job links: {str(e)}")

        return job_links

    async def check_location_relevance(
        self, location_text: Optional[str]
    ) -> Tuple[bool, str]:
        """
        Check if a job location is relevant (Costa Rica or LATAM).

        Args:
            location_text: The location text to check

        Returns:
            Tuple of (is_relevant, reason)
        """
        if not location_text:
            # If no location is specified, assume it might be relevant
            return True, "No location specified"

        location_text = location_text.lower()

        # Check for excluded locations first
        for excluded in settings.excluded_locations:
            if excluded.lower() in location_text:
                return False, f"Excluded location: {excluded}"

        # Check for target locations
        for target in settings.target_locations:
            if target.lower() in location_text:
                return True, f"Target location: {target}"

        # If no target location is found, assume it's not relevant
        return False, "No target location found"

    async def extract_page_content(self) -> Dict[str, Any]:
        """
        Extract the main content from the current page.

        Returns:
            Dictionary containing the page content
        """
        content = {
            "title": await self.page.title(),
            "url": self.page.url,
            "full_text": await self.page.content(),
        }

        # Extract the main text content (excluding scripts, styles, etc.)
        try:
            main_text = await self.page.evaluate(
                """
                () => {
                    return document.body.innerText;
                }
            """
            )
            content["main_text"] = main_text
        except Exception as e:
            logger.warning(f"Failed to extract main text: {str(e)}")
            content["main_text"] = ""

        return content
