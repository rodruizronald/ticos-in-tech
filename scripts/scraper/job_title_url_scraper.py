import asyncio
import json
import os
import sys
from datetime import datetime
from pathlib import Path

from dotenv import load_dotenv
import openai
from loguru import logger
from playwright.async_api import async_playwright

# Get the root directory
root_dir = Path(__file__).parent.parent.parent

# Load environment variables from .env file
load_dotenv(root_dir / ".env")

# Configuration
OPENAI_API_KEY = os.environ.get("OPENAI_API_KEY")

OUTPUT_DIR = "jobs"  # Directory to store results
INPUT_FILE = "companies.json"  # JSON file with company career URLs
MODEL = "o4-mini"  # OpenAI model to use
PROMPT_FILE = "prompts/basic_parser.md"  # File containing the prompt template
LOG_LEVEL = "INFO"  # Log level (DEBUG, INFO, WARNING, ERROR, CRITICAL)

# Configure logger
logger.remove()  # Remove default handler
logger.add(sys.stderr, level=LOG_LEVEL)  # Add stderr handler with desired log level
logger.add(
    f"{OUTPUT_DIR}/job_url_scraper.log", rotation="10 MB", level=LOG_LEVEL
)  # Add file handler

# Initialize OpenAI client
client = openai.OpenAI(api_key=OPENAI_API_KEY)


async def extract_html_content(url: str, selectors: list[str] = None) -> str:
    """
    Use Playwright to fetch HTML content from multiple selectors on a webpage.

    Args:
        url: The URL to fetch content from
        selectors: List of CSS selectors to extract content from (optional)

    Returns:
        String containing concatenated HTML content from all selectors,
        or full page HTML if no selectors provided
    """
    try:
        async with async_playwright() as p:
            browser = await p.chromium.launch()
            page = await browser.new_page()

            # Navigate to the URL
            logger.info(f"Navigating to {url}")
            await page.goto(url, wait_until="networkidle", timeout=60000)

            # Wait a bit for any dynamic content to load
            await page.wait_for_timeout(3000)

            if selectors:
                contents = []
                for selector in selectors:
                    try:
                        # Wait for the specific element to be available
                        element = await page.wait_for_selector(selector, timeout=5000)
                        if element:
                            # Get the HTML content of this element
                            content = await element.inner_html()
                            contents.append(content)
                            logger.info(
                                f"Successfully extracted content from selector: {selector}"
                            )
                        else:
                            logger.warning(f"Selector not found: {selector}")
                    except Exception as e:
                        logger.error(
                            f"Error extracting content from selector {selector}: {str(e)}"
                        )

                # Concatenate all contents with a newline between them
                content = "\n".join(contents) if contents else None
            else:
                # Get the full page content if no selectors specified
                content = await page.content()

            await browser.close()
            return content
    except Exception as e:
        logger.error(f"Error fetching {url}: {str(e)}")
        return None


def read_prompt_template():
    """Read the prompt template from a file."""
    try:
        with open(PROMPT_FILE, "r") as f:
            return f.read()
    except FileNotFoundError:
        logger.error(f"Error: Prompt template file '{PROMPT_FILE}' not found.")
        exit(1)
    except Exception as e:
        logger.error(f"Error reading prompt template: {str(e)}")
        exit(1)


async def process_company(
    company_name: str, career_url: str, selectors: list[str] = None
):
    """Process a single company's career page."""
    logger.info(f"Processing {company_name}...")

    # Extract HTML content
    html_content = await extract_html_content(career_url, selectors)
    if not html_content:
        logger.warning(f"Could not fetch content for {company_name}")
        return {
            "jobs": [],
            "error": "Failed to fetch HTML content",
        }

    # Read prompt template and fill it with HTML content
    prompt_template = read_prompt_template()
    # Escape any curly braces in the HTML content
    escaped_html = html_content.replace("{", "{{").replace("}", "}}")
    filled_prompt = prompt_template.replace("{html_content}", escaped_html)

    # Send to OpenAI
    try:
        logger.info(f"Sending content to OpenAI for {company_name}...")
        response = client.chat.completions.create(
            model=MODEL,
            messages=[
                {
                    "role": "system",
                    "content": "You extract job href links from HTML content.",
                },
                {"role": "user", "content": filled_prompt},
            ],
            response_format={"type": "json_object"},
        )

        # Parse response
        response_text = response.choices[0].message.content
        job_data = json.loads(response_text)

        # Add company metadata
        result = {
            "jobs": job_data.get("jobs", []),
            "timestamp": datetime.now().isoformat(),
        }

        logger.success(f"Found {len(result['jobs'])} job links for {company_name}")
        return result

    except Exception as e:
        logger.error(f"Error processing {company_name} with OpenAI: {str(e)}")
        return {
            "jobs": [],
            "error": str(e),
        }


async def main():
    """Main function to process all companies."""
    # Create output directory if it doesn't exist
    Path(OUTPUT_DIR).mkdir(exist_ok=True)

    # Check if prompt template file exists
    if not Path(PROMPT_FILE).exists():
        logger.error(f"Prompt template file '{PROMPT_FILE}' not found")
        return

    # Read input file with company data
    try:
        with open(INPUT_FILE, "r") as f:
            companies = json.load(f)
    except Exception as e:
        logger.error(f"Error reading input file: {str(e)}")
        return

    # Process each company
    companies_jobs = {"jobs": []}  # Initialize the structure for all jobs

    # Process each company
    for company in companies:
        company_name = company.get("name")
        career_url = company.get("career_url")
        job_board_selector = company.get("html_selectors", {}).get(
            "job_board_selector", []
        )
        job_description_selector = company.get("html_selectors", {}).get(
            "job_description_selector", []
        )

        if not company_name or not career_url:
            logger.warning("Skipping entry with missing name or URL")
            continue

        result = await process_company(company_name, career_url, job_board_selector)

        # Add to the all_jobs structure
        companies_jobs["companies"].append(
            {
                "company": company_name,
                "job_description_selector": job_description_selector,
                "jobs": result["jobs"],
            }
        )

        # Delay to avoid rate limiting
        await asyncio.sleep(1)

    # Create dated directory and save results
    timestamp = datetime.now().strftime("%Y%m%d")
    output_dir = Path(OUTPUT_DIR) / timestamp
    output_dir.mkdir(exist_ok=True, parents=True)

    output_file = output_dir / "companies_jobs.json"
    with open(output_file, "w") as f:
        json.dump(companies_jobs, f, indent=2)

    logger.info(f"Processing complete. Results saved to {output_file}")


if __name__ == "__main__":
    # Check for API key
    if not OPENAI_API_KEY:
        logger.error("OPENAI_API_KEY environment variable is not set")
        exit(1)

    logger.info("Starting job link extraction process")
    # Run the async main function
    asyncio.run(main())
    logger.info("Process completed")
