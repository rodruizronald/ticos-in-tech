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
MODEL = "o4-mini"  # OpenAI model to use
PROMPT_FILE = "prompts/job_technologies.md"  # File containing the prompt template
LOG_LEVEL = "DEBUG"  # Log level (DEBUG, INFO, WARNING, ERROR, CRITICAL)

# Configure logger
logger.remove()  # Remove default handler
logger.add(sys.stderr, level=LOG_LEVEL)  # Add stderr handler with desired log level
logger.add(
    f"{OUTPUT_DIR}/job_technologies_scraper.log", rotation="10 MB", level=LOG_LEVEL
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


async def extract_job_technologies(job_url: str, selectors: list[str]):
    """Extract technologies from a job posting URL."""
    logger.info(f"Extracting technologies for job at {job_url}")

    # Extract HTML content
    html_content = await extract_html_content(job_url, selectors)
    if not html_content:
        logger.warning(f"Could not fetch content for job {job_url}")
        return None

    # Read prompt template and fill it with HTML content
    prompt_template = read_prompt_template()
    # Escape any curly braces in the HTML content
    escaped_html = html_content.replace("{", "{{").replace("}", "}}")
    filled_prompt = prompt_template.replace("{html_content}", escaped_html)

    # Send to OpenAI
    try:
        logger.info(f"Sending content to OpenAI for job {job_url}...")
        response = client.chat.completions.create(
            model=MODEL,
            messages=[
                {
                    "role": "system",
                    "content": "You extract technologies from job postings.",
                },
                {"role": "user", "content": filled_prompt},
            ],
            response_format={"type": "json_object"},
        )

        # Parse response
        response_text = response.choices[0].message.content
        tech_data = json.loads(response_text)

        logger.success(f"Successfully extracted technologies for job {job_url}")
        return tech_data.get("technologies", [])

    except Exception as e:
        logger.error(f"Error extracting technologies for job {job_url}: {str(e)}")
        return None


async def main():
    """Main function to process all jobs."""
    # Get today's date for the input directory
    today = datetime.now().strftime("%Y%m%d")
    input_dir = Path(OUTPUT_DIR) / today

    if not input_dir.exists():
        logger.error(f"Input directory {input_dir} does not exist")
        return

    input_file = input_dir / "jobs_with_descriptions.json"
    if not input_file.exists():
        logger.error(f"Input file {input_file} does not exist")
        return

    # Read input file with jobs data
    try:
        with open(input_file, "r") as f:
            data = json.load(f)
    except Exception as e:
        logger.error(f"Error reading input file: {str(e)}")
        return

    # Process each job
    final_jobs = []

    for job in data.get("jobs", []):
        job_url = job.get("application_url")
        is_job_new = job.get("new", True)

        if not job_url:
            logger.warning(f"Job missing application_url, skipping: {job.get('title')}")
            continue

        # Skip technology extraction for jobs that aren't new
        if not is_job_new:
            logger.debug(
                f"Job is not new, skipping technology extraction: {job.get('title')}"
            )
            # Create a clean job object without the excluded fields
            clean_job = {
                k: v
                for k, v in job.items()
                if k not in ["job_description_selector", "eligible", "new"]
            }

            # If the job already has technologies, keep them
            if "technologies" not in clean_job:
                clean_job["technologies"] = []

            final_jobs.append(clean_job)
            continue

        # Get job description selectors directly from the job object
        selectors = job.get("job_description_selector", [])

        # Extract job technologies
        technologies = await extract_job_technologies(job_url, selectors)

        if technologies:
            # Add technologies to job data
            job["technologies"] = technologies
            logger.info(
                f"Added {len(technologies)} technologies to job: {job.get('title')}"
            )
        else:
            # Add empty technologies array if extraction failed
            job["technologies"] = []
            logger.warning(
                f"Failed to extract technologies for job: {job.get('title')}"
            )

        # Create a clean job object without the excluded fields
        clean_job = {
            k: v
            for k, v in job.items()
            if k not in ["job_description_selector", "eligible", "new"]
        }

        # Add job to final jobs list
        final_jobs.append(clean_job)

        # Delay to avoid rate limiting
        await asyncio.sleep(1)

    # Save final jobs data
    output_file = input_dir / "final_jobs.json"
    with open(output_file, "w", encoding="utf-8") as f:
        json.dump(
            {"jobs": final_jobs},
            f,
            indent=2,
            ensure_ascii=False,  # This will prevent Unicode escaping
        )

    logger.info(f"Processing complete. Results saved to {output_file}")


if __name__ == "__main__":
    # Check for API key
    if not OPENAI_API_KEY:
        logger.error("OPENAI_API_KEY environment variable is not set")
        exit(1)

    logger.info("Starting job technologies extraction process")
    # Run the async main function
    asyncio.run(main())
    logger.info("Process completed")
