# Job Search Automation Agent

An AI-powered agent for scraping and processing job listings from company career websites, focusing on opportunities in Costa Rica and LATAM.

## System Architecture

The job search automation agent is built using a clean, feature-based architecture that separates concerns and promotes maintainability. The system consists of the following main components:

```
                                 ┌─────────────────┐
                                 │                 │
                                 │   Job Agent     │
                                 │                 │
                                 └────────┬────────┘
                                          │
                                          │ orchestrates
                                          ▼
           ┌───────────────┬──────────────────────────┬───────────────┐
           │               │                          │               │
┌──────────▼─────────┐ ┌───▼───────────────┐ ┌────────▼──────────┐   │
│                    │ │                   │ │                   │   │
│  Browser Manager   │ │  Job Processor    │ │ Technology        │   │
│                    │ │                   │ │ Extractor         │   │
└──────────┬─────────┘ └───┬───────────────┘ └────────┬──────────┘   │
           │               │                          │               │
           │               │                          │               │
┌──────────▼─────────┐    │                          │               │
│                    │    │                          │               │
│  Page Handler      │    │                          │               │
│                    │    │                          │               │
└────────────────────┘    │                          │               │
                          │                          │               │
                          │                          │               │
                          ▼                          ▼               ▼
                    ┌─────────────────────────────────────────────────┐
                    │                                                 │
                    │           Database Repositories                 │
                    │                                                 │
                    └─────────────────────────────────────────────────┘
```

### Components

1. **Job Agent**: Orchestrates the entire scraping and processing workflow, managing the interaction between different components.

2. **Browser Manager**: Handles browser initialization, navigation, and cleanup using Playwright.

3. **Page Handler**: Provides methods for interacting with web pages and extracting data.

4. **Job Processor**: Processes job data, including signature generation, slug creation, and database operations.

5. **Technology Extractor**: Uses AI to extract technology information from job descriptions and match them to the technology database.

6. **Database Repositories**: Interface with the database to store and retrieve data.

## Installation

### Prerequisites

- Python 3.9+
- PostgreSQL database
- Node.js (for Playwright)

### Setup

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd <repository-directory>
   ```

2. Create and activate a virtual environment:
   ```bash
   python -m venv venv
   source venv/bin/activate  # On Windows: venv\Scripts\activate
   ```

3. Install dependencies:
   ```bash
   pip install -r requirements.txt
   ```

4. Install Playwright browsers:
   ```bash
   playwright install chromium
   ```

5. Set up environment variables:
   ```bash
   # Database connection
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_USER=postgres
   export DB_PASSWORD=postgres
   export DB_NAME=jobboard
   
   # OpenAI API key (for LangChain)
   export OPENAI_API_KEY=your-api-key
   ```

## Usage

### Running the Agent

To run the job search automation agent:

```bash
python -m src.scraper.main
```

### Command-line Options

- `--log-level`: Set the logging level (DEBUG, INFO, WARNING, ERROR, CRITICAL)
- `--company`: Process only the specified company (by name)

Example:
```bash
python -m src.scraper.main --log-level=DEBUG --company="Acme Corporation"
```

## How It Works

1. The agent fetches all active companies from the database.
2. For each company, it:
   - Visits their careers page using Playwright
   - Extracts job listings
   - For each job listing:
     - Determines if the job is available in Costa Rica or LATAM
     - Extracts structured data using AI
     - Generates a signature to identify duplicate jobs
     - Stores new jobs in the database
     - Updates existing jobs' last_seen_at timestamp
     - Extracts and stores technologies mentioned in the job description
   - Marks jobs no longer listed as inactive

## Extending the Agent

### Adding Support for New Career Page Formats

To add support for new career page formats, you can modify the `_extract_job_links` method in the `JobAgent` class to include additional CSS selectors for job links.

### Customizing Location Filtering

Location filtering settings can be adjusted in the `config.py` file by modifying the `target_locations` and `excluded_locations` lists.

### Adding New AI Capabilities

The agent uses LangChain and OpenAI's models for AI tasks. You can extend the AI capabilities by:

1. Adding new prompt templates in the `JobAgent` or `TechnologyExtractor` classes
2. Creating new methods that use the language model for specific tasks
3. Adjusting the system prompts to improve extraction accuracy

## Troubleshooting

### Common Issues

1. **Browser Initialization Fails**:
   - Ensure Playwright is properly installed: `playwright install chromium`
   - Check if you have sufficient memory and disk space

2. **Database Connection Errors**:
   - Verify database connection settings in environment variables
   - Ensure the PostgreSQL server is running

3. **API Rate Limiting**:
   - If you encounter rate limiting from OpenAI, adjust the request delay in `config.py`

4. **Memory Issues**:
   - For large job descriptions, you might need to increase the available memory
   - Consider processing companies in smaller batches

### Logging

The agent uses Python's logging module to provide detailed logs. You can adjust the log level using the `--log-level` command-line option.

Log files are written to stdout by default. You can modify the logging configuration in `utils/logging.py` to write logs to a file if needed.

## License

[Specify the license here]
