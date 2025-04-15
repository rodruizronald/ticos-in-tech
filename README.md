# Ticos in Tech - Job Search Automation

An AI-powered job search automation system that monitors company career websites for relevant opportunities in Costa Rica and LATAM.

## Overview

This project provides an intelligent agent that scrapes job listings from company career websites, processes them using AI, and stores them in a PostgreSQL database. The system focuses on identifying job opportunities that are available in Costa Rica or open to candidates in Latin America.

## Features

- **AI-Powered Scraping**: Uses LangChain and OpenAI to intelligently extract and process job data
- **Automated Monitoring**: Regularly checks company career pages for new job listings
- **Location Filtering**: Identifies jobs available in Costa Rica or LATAM
- **Technology Extraction**: Automatically extracts and categorizes technologies mentioned in job descriptions
- **Database Integration**: Stores job data in a PostgreSQL database with a clean repository pattern

## Project Structure

```
ticos-in-tech/
├── document/                # Documentation files
│   └── database/            # Database documentation
├── src/                     # Source code
│   ├── api/                 # API layer (future)
│   ├── core/                # Core domain logic (future)
│   ├── domain/              # Domain models (future)
│   ├── infrastructure/      # Infrastructure layer
│   │   └── database/        # Database models and repositories
│   └── scraper/             # Job search automation agent
│       ├── agent/           # AI agent for orchestration
│       ├── browser/         # Browser automation with Playwright
│       ├── processors/      # Job and technology data processors
│       └── utils/           # Utility functions
└── tests/                   # Test suite (future)
```

## Architecture

The project follows clean architecture principles with a feature-based organization:

1. **Domain Layer**: Contains the core business logic and domain models
2. **Infrastructure Layer**: Provides implementations for external services like databases
3. **Scraper Module**: Contains the job search automation agent
4. **API Layer**: (Future) Will provide REST API endpoints for accessing job data

## Getting Started

### Prerequisites

- Python 3.9+
- PostgreSQL database
- Node.js (for Playwright)

### Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd ticos-in-tech
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
   cp .env.example .env
   # Edit .env with your configuration
   ```

6. Set up the database:
   ```bash
   # Run database migrations
   cd src/infrastructure/database/migrations
   alembic upgrade head
   ```

### Running the Job Search Agent

To run the job search automation agent:

```bash
python -m src.scraper.main
```

For more options and details, see the [Scraper README](src/scraper/README.md).

## Database Schema

The database schema includes the following main tables:

- **companies**: Stores information about employers posting jobs
- **jobs**: Contains job listings with references to the posting company
- **technologies**: Maintains a catalog of technologies/skills with hierarchical relationships
- **job_technologies**: Links jobs to their required technologies

For more details on the database schema and repository layer, see the [Database README](src/infrastructure/database/README.md).

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature-name`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin feature/your-feature-name`
5. Submit a pull request

## License

[Specify the license here]

## Acknowledgments

- This project uses [LangChain](https://github.com/hwchase17/langchain) for AI capabilities
- Web scraping is powered by [Playwright](https://playwright.dev/)
- Database access is handled by [SQLAlchemy](https://www.sqlalchemy.org/)
