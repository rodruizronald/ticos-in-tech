"""
Browser manager for web scraping.

This module provides a browser manager class that handles browser initialization,
navigation, and cleanup for web scraping operations.
"""

import asyncio
import logging
from typing import Optional

from playwright.async_api import (
    async_playwright,
    Browser,
    BrowserContext,
    Page,
    Playwright,
)

from src.scraper.config import settings

logger = logging.getLogger(__name__)


class BrowserManager:
    """
    Manages browser instances for web scraping operations.

    This class handles browser initialization, context creation, and cleanup.
    It provides methods for creating and managing browser pages.
    """

    def __init__(self):
        """Initialize the browser manager."""
        self._playwright: Optional[Playwright] = None
        self._browser: Optional[Browser] = None
        self._context: Optional[BrowserContext] = None

    async def initialize(self) -> None:
        """
        Initialize the browser and create a browser context.

        This method launches a new browser instance and creates a browser context
        with the configured settings.

        Raises:
            Exception: If browser initialization fails
        """
        try:
            logger.info("Initializing browser")
            self._playwright = await async_playwright().start()
            self._browser = await self._playwright.chromium.launch(
                headless=settings.headless
            )
            self._context = await self._browser.new_context(
                user_agent=settings.user_agent,
                viewport={"width": 1920, "height": 1080},
                ignore_https_errors=True,
            )

            # Set default timeout
            self._context.set_default_timeout(settings.timeout)

            logger.info("Browser initialized successfully")
        except Exception as e:
            logger.error(f"Failed to initialize browser: {str(e)}")
            await self.cleanup()
            raise

    async def new_page(self) -> Page:
        """
        Create a new browser page.

        Returns:
            A new browser page

        Raises:
            ValueError: If the browser is not initialized
        """
        if not self._context:
            raise ValueError("Browser not initialized. Call initialize() first.")

        return await self._context.new_page()

    async def goto(self, page: Page, url: str) -> bool:
        """
        Navigate to a URL with error handling.

        Args:
            page: The browser page
            url: The URL to navigate to

        Returns:
            True if navigation was successful, False otherwise
        """
        try:
            logger.info(f"Navigating to {url}")
            response = await page.goto(url, wait_until="domcontentloaded")

            # Add a small delay to ensure page is fully loaded
            await asyncio.sleep(settings.request_delay)

            if not response:
                logger.warning(f"No response received when navigating to {url}")
                return False

            if not response.ok:
                logger.warning(
                    f"Received status {response.status} when navigating to {url}"
                )
                return False

            return True
        except Exception as e:
            logger.error(f"Failed to navigate to {url}: {str(e)}")
            return False

    async def cleanup(self) -> None:
        """
        Clean up browser resources.

        This method closes the browser context and browser, and stops the playwright
        instance.
        """
        logger.info("Cleaning up browser resources")

        if self._context:
            try:
                await self._context.close()
            except Exception as e:
                logger.warning(f"Error closing browser context: {str(e)}")
            self._context = None

        if self._browser:
            try:
                await self._browser.close()
            except Exception as e:
                logger.warning(f"Error closing browser: {str(e)}")
            self._browser = None

        if self._playwright:
            try:
                await self._playwright.stop()
            except Exception as e:
                logger.warning(f"Error stopping playwright: {str(e)}")
            self._playwright = None

    async def __aenter__(self) -> "BrowserManager":
        """
        Enter the async context manager.

        Returns:
            The browser manager instance
        """
        await self.initialize()
        return self

    async def __aexit__(self, exc_type, exc_val, exc_tb) -> None:
        """
        Exit the async context manager.

        Args:
            exc_type: Exception type
            exc_val: Exception value
            exc_tb: Exception traceback
        """
        await self.cleanup()
