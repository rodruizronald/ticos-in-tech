"""
Browser module for web scraping.

This module provides browser management and web scraping functionality
using Playwright.
"""

from src.scraper.browser.manager import BrowserManager
from src.scraper.browser.page_handler import PageHandler

__all__ = ["BrowserManager", "PageHandler"]
