-- Companies Table
CREATE TABLE companies (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    careers_page_url VARCHAR(255) NOT NULL,
    logo_url VARCHAR(255),
    description TEXT,
    industry VARCHAR(100),
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Jobs Table
CREATE TABLE jobs (
    id SERIAL PRIMARY KEY,
    company_id INT REFERENCES companies(id),
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    requirements TEXT NOT NULL,
    preferred_skills TEXT,
    experience_level VARCHAR(50) NOT NULL, 
    employment_type VARCHAR(50) NOT NULL,
    location VARCHAR(50),
    work_mode VARCHAR(20) NOT NULL,
    application_url VARCHAR(255),
    job_function VARCHAR(100), -- Business function
    first_seen_at TIMESTAMP NOT NULL DEFAULT NOW(), -- timestamp when system first discovered the job
    last_seen_at TIMESTAMP NOT NULL DEFAULT NOW(), -- timestamp when system most recently saw the job still active
    posted_at TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    signature VARCHAR(64) UNIQUE, -- Store hash
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
);

-- Technologies Table
CREATE TABLE technologies (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    category VARCHAR(50),
    parent_id INT REFERENCES technologies(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Job Technologies Junction Table
CREATE TABLE job_technologies (
    id SERIAL PRIMARY KEY,
    job_id INT REFERENCES jobs(id),
    technology_id INT REFERENCES technologies(id),
    is_primary BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create Indexes for Optimized Searching

-- Jobs Indexes 

-- Primary Search Indexes
CREATE INDEX idx_jobs_title_tsvector ON jobs USING GIN (to_tsvector('english', title));
CREATE INDEX idx_jobs_description_tsvector ON jobs USING GIN (to_tsvector('english', description));

-- Filter Indexes
CREATE INDEX idx_jobs_location ON jobs(location);
--  partial index for active jobs since they'll be queried most frequently
CREATE INDEX idx_active_jobs ON jobs(id) WHERE is_active = TRUE;
CREATE INDEX idx_jobs_work_mode ON jobs(work_mode);
CREATE INDEX idx_jobs_posted_at ON jobs(posted_at);
CREATE INDEX idx_jobs_company_id ON jobs(company_id);
CREATE INDEX idx_jobs_job_function ON jobs(job_function);
CREATE UNIQUE INDEX idx_jobs_signature ON jobs(signature);
CREATE INDEX idx_jobs_employment_type ON jobs(employment_type);
CREATE INDEX idx_jobs_experience_level ON jobs(experience_level);


-- Technology Indexes
CREATE INDEX idx_job_technologies_job_id ON job_technologies(job_id);
CREATE INDEX idx_job_technologies_technology_id ON job_technologies(technology_id);

-- Company Indexes
CREATE INDEX idx_companies_name ON companies(name);
CREATE INDEX idx_companies_industry ON companies(industry);