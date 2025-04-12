# Job Board Database - CRUD Operations Documentation

## Overview
This document provides comprehensive CRUD (Create, Read, Update, Delete) operations for the job board database. The database contains four main tables: `companies`, `jobs`, `technologies`, and `job_technologies`. This documentation is intended for developers implementing application code that interacts with this database.

## Table of Contents
1. [Schema Overview](#schema-overview)
2. [Companies Table](#companies-table)
3. [Jobs Table](#jobs-table)
4. [Technologies Table](#technologies-table)
5. [Job Technologies Table](#job-technologies-table)

## Schema Overview

The job board database consists of the following tables:

1. **Companies Table**: Stores information about employers posting jobs.
2. **Jobs Table**: Contains job listings with references to the posting company.
3. **Technologies Table**: Maintains a catalog of technologies/skills with hierarchical relationships.
4. **Job Technologies Junction Table**: Links jobs to their required technologies.

![Database Schema Relationships]

The schema is optimized for:
- Efficient job searching by various criteria
- Technology/skill-based filtering
- Company-based filtering
- Full-text search capabilities

## Companies Table

### Table Overview
The `companies` table stores information about employers that post jobs on the platform. It includes company details such as name, industry, logo, and contact information.

#### Table Structure
```sql
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
```

### CREATE Operations

#### Required Fields vs. Optional Fields
- **Required Fields**: `name`, `careers_page_url`
- **Optional Fields**: `logo_url`, `description`, `industry`

#### Default Values and Auto-generated Values
- `id`: Auto-incremented unique identifier
- `active`: Defaults to `TRUE`
- `created_at`, `updated_at`: Default to current timestamp

#### Validation Rules and Constraints
- `name` must not be null
- `careers_page_url` must not be null

#### Example SQL and Description

```sql
INSERT INTO companies (
  name,
  careers_page_url,
  logo_url,
  description,
  industry
) VALUES (
  'Acme Corporation',
  'https://acme.example.com/careers',
  'https://acme.example.com/logo.png',
  'A leading provider of innovative solutions',
  'Technology'
);
```

**Plain English Description**: This operation creates a new company record with the company name, careers page URL, logo URL, description, and industry. The system automatically assigns an ID and timestamps for creation and last update. By default, the company is marked as active.

### READ Operations

#### Common Query Patterns

1. **Get Company by ID**
```sql
SELECT * FROM companies WHERE id = 123;
```

2. **Search Companies by Name**
```sql
SELECT * FROM companies 
WHERE name ILIKE '%acme%' 
ORDER BY name ASC;
```

3. **Filter Companies by Industry**
```sql
SELECT * FROM companies 
WHERE industry = 'Technology' AND active = TRUE 
ORDER BY name ASC;
```

4. **List All Active Companies**
```sql
SELECT * FROM companies 
WHERE active = TRUE 
ORDER BY name ASC;
```

#### Recommended Indexes for Performance
- `idx_companies_name`: Already created for name-based searches
- `idx_companies_industry`: Already created for industry filtering

#### Pagination Strategy
For paginated results, use LIMIT and OFFSET:

```sql
SELECT * FROM companies 
WHERE active = TRUE 
ORDER BY name ASC 
LIMIT 20 OFFSET 40; -- Page 3 with 20 items per page
```

**Plain English Description**: These operations retrieve company data based on different criteria such as ID, name, or industry. Results can be paginated to manage large result sets efficiently, typically 20-50 companies per page.

### UPDATE Operations

#### Fields that can/cannot be updated
All fields except `id` and `created_at` can be updated.

#### Concurrency Considerations
- Update the `updated_at` timestamp on each modification to track the most recent change.
- Consider using PostgreSQL's optimistic concurrency control with the `updated_at` field.

#### Example SQL and Description

```sql
UPDATE companies 
SET 
  name = 'Acme Corporation International',
  description = 'A global leader in innovative solutions',
  industry = 'Information Technology',
  logo_url = 'https://acme.example.com/new_logo.png',
  updated_at = NOW()
WHERE id = 123;
```

**Plain English Description**: This operation updates a company's information including name, description, industry, and logo URL. The `updated_at` timestamp is refreshed to the current time to track when the record was last modified.

To toggle a company's active status:

```sql
UPDATE companies 
SET active = NOT active, updated_at = NOW() 
WHERE id = 123;
```

**Plain English Description**: This toggles a company's active status (from active to inactive or vice versa), which affects whether the company and its jobs appear in normal search results.

### DELETE Operations

#### Hard Delete vs. Soft Delete Recommendations
It's recommended to use soft deletes by setting the `active` field to `FALSE` rather than physically removing records. This preserves historical data and references.

#### Cascading Delete Implications
Deleting a company may affect:
- Jobs associated with the company
- Job-technology relationships for those jobs

#### Referential Integrity Considerations
Due to foreign key constraints, a hard delete would require deleting or updating all related jobs first.

#### Example SQL and Description

**Soft Delete (Recommended)**:
```sql
UPDATE companies 
SET 
  active = FALSE,
  updated_at = NOW() 
WHERE id = 123;
```

**Plain English Description**: This marks a company as inactive rather than removing it from the database. This approach preserves historical data while effectively removing the company from active searches and listings.

**Hard Delete (Use with Caution)**:
```sql
DELETE FROM companies WHERE id = 123;
```

**Plain English Description**: This permanently removes the company from the database. This operation should be used with extreme caution as it will fail if there are any jobs still associated with this company due to foreign key constraints.

### Special Considerations
- Consider implementing audit logging for company changes, especially for status changes.
- The `active` field provides a way to temporarily or permanently remove companies from search results without losing data.
- Industry values should ideally come from a standardized list for consistency in filtering and reporting.

## Jobs Table

### Table Overview
The `jobs` table stores job postings with detailed information about positions. Each job is associated with a company and can be linked to multiple technologies through the `job_technologies` junction table.

#### Table Structure
```sql
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
  job_function VARCHAR(100),
  first_seen_at TIMESTAMP NOT NULL DEFAULT NOW(),
  last_seen_at TIMESTAMP NOT NULL DEFAULT NOW(),
  posted_at TIMESTAMP,
  is_active BOOLEAN DEFAULT TRUE,
  signature VARCHAR(64) UNIQUE,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### CREATE Operations

#### Required Fields vs. Optional Fields
- **Required Fields**: `company_id`, `title`, `slug`, `description`, `requirements`, `experience_level`, `employment_type`, `work_mode`
- **Optional Fields**: `preferred_skills`, `location`, `application_url`, `job_function`, `posted_at`, `signature`

#### Default Values and Auto-generated Values
- `id`: Auto-incremented unique identifier
- `is_active`: Defaults to `TRUE`
- `first_seen_at`, `last_seen_at`: Default to current timestamp
- `created_at`, `updated_at`: Default to current timestamp

#### Validation Rules and Constraints
- `title`, `slug`, `description`, `requirements` must not be null
- `experience_level`, `employment_type`, `work_mode` must not be null
- `signature` must be unique (used to prevent duplicate job listings)
- `company_id` must reference a valid company ID

#### Foreign Key Considerations
- `company_id` references `companies(id)` - the company must exist before creating a job

#### Example SQL and Description

```sql
INSERT INTO jobs (
  company_id,
  title,
  slug,
  description,
  requirements,
  preferred_skills,
  experience_level,
  employment_type,
  location,
  work_mode,
  application_url,
  job_function,
  posted_at,
  signature
) VALUES (
  42,
  'Senior Software Engineer',
  'senior-software-engineer-42',
  'We are looking for an experienced software engineer to join our team...',
  'Minimum 5 years experience with JavaScript. Strong problem-solving skills.',
  'Experience with React, Node.js, and AWS preferred.',
  'Senior',
  'Full-time',
  'San Francisco, CA',
  'Remote',
  'https://example.com/apply/job123',
  'Engineering',
  '2024-04-03 12:00:00',
  'a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6'
);
```

**Plain English Description**: This operation creates a new job posting with all the necessary details including title, description, requirements, and association with a company (company_id). The system automatically assigns an ID and timestamps. The unique signature helps prevent duplicate job listings.

### READ Operations

#### Common Query Patterns

1. **Get Job by ID**
```sql
SELECT * FROM jobs WHERE id = 456;
```

2. **Search Jobs by Title/Keywords**
```sql
SELECT j.*, c.name as company_name 
FROM jobs j
JOIN companies c ON j.company_id = c.id
WHERE 
  j.is_active = TRUE AND
  to_tsvector('english', j.title || ' ' || j.description) @@ to_tsquery('english', 'software & engineer');
```

3. **Filter Jobs by Multiple Criteria**
```sql
SELECT j.*, c.name as company_name
FROM jobs j
JOIN companies c ON j.company_id = c.id
WHERE 
  j.is_active = TRUE AND
  j.work_mode = 'Remote' AND
  j.experience_level = 'Senior' AND
  j.employment_type = 'Full-time' AND
  j.posted_at >= NOW() - INTERVAL '7 days';
```

4. **Filter Jobs by Company**
```sql
SELECT * FROM jobs 
WHERE company_id = 42 AND is_active = TRUE
ORDER BY posted_at DESC;
```

5. **Filter Jobs by Technology/Skill**
```sql
SELECT j.*, c.name as company_name
FROM jobs j
JOIN companies c ON j.company_id = c.id
JOIN job_technologies jt ON j.id = jt.job_id
JOIN technologies t ON jt.technology_id = t.id
WHERE 
  j.is_active = TRUE AND
  t.name = 'React';
```

6. **Filter Jobs by Multiple Technologies**
```sql
SELECT j.*, c.name as company_name
FROM jobs j
JOIN companies c ON j.company_id = c.id
WHERE j.id IN (
  SELECT jt1.job_id
  FROM job_technologies jt1
  JOIN technologies t1 ON jt1.technology_id = t1.id
  WHERE t1.name = 'React'
) AND j.id IN (
  SELECT jt2.job_id
  FROM job_technologies jt2
  JOIN technologies t2 ON jt2.technology_id = t2.id
  WHERE t2.name = 'Node.js'
) AND j.is_active = TRUE;
```

7. **Get Recent Jobs with Date Filter**
```sql
SELECT j.*, c.name as company_name
FROM jobs j
JOIN companies c ON j.company_id = c.id
WHERE 
  j.is_active = TRUE AND
  j.posted_at >= CURRENT_DATE - INTERVAL '30 days'
ORDER BY j.posted_at DESC;
```

#### Recommended Indexes for Performance
The schema already includes these important indexes:
- Text search indexes: `idx_jobs_title_tsvector`, `idx_jobs_description_tsvector`
- Filter indexes: `idx_jobs_location`, `idx_active_jobs`, `idx_jobs_work_mode`, etc.
- Unique index on `signature` to prevent duplicates
- Foreign key index on `company_id`

#### Pagination Strategy
For job search results, implement cursor-based pagination for better performance:

```sql
SELECT j.*, c.name as company_name
FROM jobs j
JOIN companies c ON j.company_id = c.id
WHERE 
  j.is_active = TRUE AND
  j.posted_at < '2024-03-15 12:34:56' -- Last item from previous page
ORDER BY j.posted_at DESC
LIMIT 20;
```

Or use standard offset pagination:

```sql
SELECT j.*, c.name as company_name
FROM jobs j
JOIN companies c ON j.company_id = c.id
WHERE j.is_active = TRUE
ORDER BY j.posted_at DESC
LIMIT 20 OFFSET 40; -- Page 3 with 20 items per page
```

**Plain English Description**: These operations retrieve job listings based on various criteria such as keywords in the title/description, specific filters (work mode, experience level), company, or required technologies. Results are typically sorted by posting date (newest first) and paginated to manage large result sets.

### UPDATE Operations

#### Fields that can/cannot be updated
Most fields can be updated except `id`, `created_at`, and typically `signature` (since it represents the job's unique identity).

#### Concurrency Considerations
- Update the `updated_at` timestamp on each modification
- Update `last_seen_at` whenever the job is verified as still active
- Be cautious when updating the `signature` as it's used to identify duplicate jobs

#### Cascading Update Effects
- Updating a job doesn't automatically affect linked technologies
- If a job's core requirements change significantly, technologies may need to be updated separately

#### Example SQL and Description

```sql
UPDATE jobs 
SET 
  title = 'Senior Full Stack Engineer',
  description = 'Updated description with new team information...',
  requirements = 'Updated requirements with additional qualifications...',
  preferred_skills = 'React, Node.js, PostgreSQL, AWS',
  location = 'New York, NY or Remote',
  last_seen_at = NOW(),
  updated_at = NOW()
WHERE id = 456;
```

**Plain English Description**: This operation updates various details of a job posting including title, description, requirements, and location. The `last_seen_at` and `updated_at` timestamps are refreshed to indicate when the job was last verified as active and when its data was last modified.

To mark a job as inactive (no longer available):

```sql
UPDATE jobs 
SET 
  is_active = FALSE,
  updated_at = NOW() 
WHERE id = 456;
```

**Plain English Description**: This marks a job as inactive, which will remove it from active search results while preserving its data in the system.

### DELETE Operations

#### Hard Delete vs. Soft Delete Recommendations
Soft deletes are strongly recommended for jobs by setting `is_active = FALSE`. This preserves historical data and prevents breaking relationships with job technology records.

#### Cascading Delete Implications
Hard deleting a job would require:
- Deleting all related job_technologies entries first
- Potentially impacting historical data and analytics

#### Example SQL and Description

**Soft Delete (Recommended)**:
```sql
UPDATE jobs 
SET 
  is_active = FALSE,
  updated_at = NOW() 
WHERE id = 456;
```

**Plain English Description**: This marks a job as inactive rather than removing it from the database. This approach preserves historical data while effectively removing the job from active searches.

**Hard Delete (Use with Caution)**:
```sql
-- First delete related job_technologies records
DELETE FROM job_technologies WHERE job_id = 456;

-- Then delete the job
DELETE FROM jobs WHERE id = 456;
```

**Plain English Description**: This permanently removes the job and its technology relationships from the database. This operation should be used with caution as it permanently removes data that might be valuable for historical analysis.

### Special Considerations
- The `signature` field is crucial for preventing duplicate jobs. It should be generated consistently, possibly as a hash of key job attributes.
- The `first_seen_at` and `last_seen_at` fields track when a job was first discovered and last confirmed active, useful for job freshness metrics.
- Consider implementing a scheduled task to automatically mark old jobs as inactive if they haven't been seen recently.
- Full-text search indexes on title and description enable powerful keyword searching.

## Technologies Table

### Table Overview
The `technologies` table stores a catalog of technology skills that can be associated with jobs. It supports hierarchical relationships (parent-child) between technologies and categorization.

#### Table Structure
```sql
CREATE TABLE technologies (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL UNIQUE,
  category VARCHAR(50),
  parent_id INT REFERENCES technologies(id),
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### CREATE Operations

#### Required Fields vs. Optional Fields
- **Required Fields**: `name`
- **Optional Fields**: `category`, `parent_id`

#### Default Values and Auto-generated Values
- `id`: Auto-incremented unique identifier
- `created_at`: Default to current timestamp

#### Validation Rules and Constraints
- `name` must not be null and must be unique
- `parent_id` must reference a valid technology ID if provided

#### Foreign Key Considerations
- `parent_id` is a self-reference to `technologies(id)` - the parent technology must exist before creating a child technology

#### Example SQL and Description

```sql
INSERT INTO technologies (
  name,
  category,
  parent_id
) VALUES (
  'React',
  'Frontend Framework',
  (SELECT id FROM technologies WHERE name = 'JavaScript')
);
```

**Plain English Description**: This operation creates a new technology entry for "React" in the "Frontend Framework" category, establishing it as a child of "JavaScript". The system automatically assigns an ID and creation timestamp.

To create a top-level technology without a parent:

```sql
INSERT INTO technologies (
  name,
  category
) VALUES (
  'JavaScript',
  'Programming Language'
);
```

**Plain English Description**: This creates a top-level technology entry for "JavaScript" in the "Programming Language" category with no parent technology.

### READ Operations

#### Common Query Patterns

1. **Get Technology by ID**
```sql
SELECT * FROM technologies WHERE id = 25;
```

2. **Get Technology by Name**
```sql
SELECT * FROM technologies WHERE name = 'React';
```

3. **List Technologies by Category**
```sql
SELECT * FROM technologies 
WHERE category = 'Programming Language' 
ORDER BY name ASC;
```

4. **Get Child Technologies**
```sql
SELECT * FROM technologies 
WHERE parent_id = 1 -- JavaScript ID
ORDER BY name ASC;
```

5. **Get Technology with Its Parent**
```sql
SELECT t.*, p.name as parent_name, p.category as parent_category
FROM technologies t
LEFT JOIN technologies p ON t.parent_id = p.id
WHERE t.name = 'React';
```

6. **Get Technology Hierarchy (Recursive Query)**
```sql
WITH RECURSIVE tech_tree AS (
  SELECT id, name, category, parent_id, 1 as level
  FROM technologies
  WHERE name = 'JavaScript' -- Start with this technology
  
  UNION ALL
  
  SELECT t.id, t.name, t.category, t.parent_id, tt.level + 1
  FROM technologies t
  JOIN tech_tree tt ON t.parent_id = tt.id
)
SELECT * FROM tech_tree ORDER BY level, name;
```

#### Example SQL and Description

**Plain English Description**: These operations retrieve technology data based on different criteria such as ID, name, or category. More complex queries can retrieve the parent-child relationships between technologies or the entire technology hierarchy starting from a specific technology.

### UPDATE Operations

#### Fields that can/cannot be updated
All fields except `id` and `created_at` can be updated.

#### Validation Requirements
- If updating `name`, ensure it remains unique
- If updating `parent_id`, ensure it references a valid technology and avoid circular references

#### Example SQL and Description

```sql
UPDATE technologies 
SET 
  name = 'ReactJS',
  category = 'JavaScript Framework'
WHERE id = 25;
```

**Plain English Description**: This operation updates a technology's name and category. The unique constraint on the name field ensures there won't be duplicate technology names.

To update a technology's parent:

```sql
UPDATE technologies 
SET parent_id = (SELECT id FROM technologies WHERE name = 'Web Technologies')
WHERE id = 25;
```

**Plain English Description**: This changes the parent technology of an existing technology, which effectively moves it to a different position in the technology hierarchy.

### DELETE Operations

#### Cascading Delete Implications
Deleting a technology may affect:
- Child technologies that reference it as a parent
- Job-technology relationships that use this technology

#### Referential Integrity Considerations
Due to foreign key constraints, deleting a technology with child technologies or associated jobs will fail.

#### Example SQL and Description

**Delete with No Dependencies**:
```sql
DELETE FROM technologies WHERE id = 25;
```

**Plain English Description**: This permanently removes a technology from the database. This will only succeed if no child technologies reference it and no jobs are associated with it.

**Safe Delete with Reassignment**:
```sql
-- First update any child technologies to point to a different parent
UPDATE technologies 
SET parent_id = (SELECT parent_id FROM technologies WHERE id = 25)
WHERE parent_id = 25;

-- Then remove any job_technologies references
DELETE FROM job_technologies WHERE technology_id = 25;

-- Finally delete the technology
DELETE FROM technologies WHERE id = 25;
```

**Plain English Description**: This safely removes a technology by first reassigning its children to its own parent (preserving hierarchy) and removing any job associations. This approach maintains referential integrity while removing the technology.

### Special Considerations
- The hierarchical structure allows for creating technology trees (e.g., JavaScript → React → React Native)
- The `category` field enables grouping similar technologies for easier browsing and filtering
- Consider implementing a check to prevent circular references in the parent-child relationships
- When displaying technologies to users, consider showing both the technology name and its parent/category for context

## Job Technologies Table

### Table Overview
The `job_technologies` junction table creates many-to-many relationships between jobs and technologies, indicating which technologies are required for each job. It also has an `is_primary` flag to highlight the most important technologies for a job.

#### Table Structure
```sql
CREATE TABLE job_technologies (
  id SERIAL PRIMARY KEY,
  job_id INT REFERENCES jobs(id),
  technology_id INT REFERENCES technologies(id),
  is_primary BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### CREATE Operations

#### Required Fields vs. Optional Fields
- **Required Fields**: `job_id`, `technology_id`
- **Optional Fields**: `is_primary`

#### Default Values and Auto-generated Values
- `id`: Auto-incremented unique identifier
- `is_primary`: Defaults to `FALSE`
- `created_at`: Default to current timestamp

#### Validation Rules and Constraints
- `job_id` must reference a valid job ID
- `technology_id` must reference a valid technology ID

#### Foreign Key Considerations
- `job_id` references `jobs(id)` - the job must exist before creating the relationship
- `technology_id` references `technologies(id)` - the technology must exist before creating the relationship

#### Example SQL and Description

```sql
INSERT INTO job_technologies (
  job_id,
  technology_id,
  is_primary
) VALUES (
  456, -- Job ID for "Senior Software Engineer"
  25,  -- Technology ID for "React"
  TRUE -- This is a primary/key technology for the job
);
```

**Plain English Description**: This operation creates a relationship between a job and a technology, indicating that the job requires that technology skill. Setting `is_primary` to TRUE marks this as a primary or key technology for the job, which may be highlighted in job listings or used for primary filtering.

To add multiple technologies to a job at once:

```sql
INSERT INTO job_technologies (job_id, technology_id, is_primary)
VALUES
  (456, 25, TRUE),  -- React (primary)
  (456, 17, FALSE), -- AWS
  (456, 42, FALSE); -- Docker
```

**Plain English Description**: This adds multiple technology requirements to a single job in one operation, with one marked as the primary technology.

### READ Operations

#### Common Query Patterns

1. **Get All Technologies for a Job**
```sql
SELECT t.name, t.category, jt.is_primary
FROM job_technologies jt
JOIN technologies t ON jt.technology_id = t.id
WHERE jt.job_id = 456
ORDER BY jt.is_primary DESC, t.name ASC;
```

2. **Get Primary Technologies for a Job**
```sql
SELECT t.name, t.category
FROM job_technologies jt
JOIN technologies t ON jt.technology_id = t.id
WHERE jt.job_id = 456 AND jt.is_primary = TRUE;
```

3. **Get All Jobs Requiring a Specific Technology**
```sql
SELECT j.*, c.name as company_name
FROM jobs j
JOIN companies c ON j.company_id = c.id
JOIN job_technologies jt ON j.id = jt.job_id
WHERE jt.technology_id = 25 AND j.is_active = TRUE
ORDER BY j.posted_at DESC;
```

#### Recommended Indexes for Performance
The schema already includes these important indexes:
- `idx_job_technologies_job_id` for looking up technologies by job
- `idx_job_technologies_technology_id` for looking up jobs by technology

#### Example SQL and Description

**Plain English Description**: These operations retrieve the relationships between jobs and technologies, allowing you to find all technologies required for a specific job or all jobs requiring a specific technology.

### UPDATE Operations

#### Fields that can/cannot be updated
Only the `is_primary` field should typically be updated. The relationship itself (`job_id` and `technology_id`) should generally be removed and recreated if it needs to change.

#### Example SQL and Description

```sql
UPDATE job_technologies 
SET is_primary = TRUE
WHERE job_id = 456 AND technology_id = 17;
```

**Plain English Description**: This operation marks a specific technology as primary for a job. You might use this when you want to highlight a different technology as the main skill for the job.

### DELETE Operations

#### Example SQL and Description

**Remove a Single Technology from a Job**:
```sql
DELETE FROM job_technologies 
WHERE job_id = 456 AND technology_id = 42;
```

**Plain English Description**: This removes a specific technology requirement from a job. This might be done when updating a job posting to reflect changed requirements.

**Remove All Technologies from a Job**:
```sql
DELETE FROM job_technologies WHERE job_id = 456;
```

**Plain English Description**: This removes all technology requirements from a job. This might be done before adding a completely new set of technologies to the job.

### Special Considerations
- When adding technologies to a job, consider limiting the number of primary technologies (typically 1-3) to maintain focus
- Ensure consistency between the technologies in this table and the skills mentioned in the job description
- This junction table is critical for technology-based job searches, which are one of the most important features of a technical job board
- Consider implementing a check to prevent duplicate job-technology combinations (although the database won't enforce this)

## Performance Optimization Considerations

### Recommended Query Patterns

1. **Job Search with Multiple Filters**
```sql
SELECT DISTINCT j.id, j.title, j.description, j.location, j.work_mode, 
       j.experience_level, j.employment_type, j.posted_at,
       c.name as company_name, c.industry
FROM jobs j
JOIN companies c ON j.company_id = c.id
LEFT JOIN job_technologies jt ON j.id = jt.job_id
LEFT JOIN technologies t ON jt.technology_id = t.id
WHERE 
  j.is_active = TRUE AND
  c.active = TRUE AND
  (
    to_tsvector('english', j.title || ' ' || j.description) @@ 
    to_tsquery('english', 'senior & (developer | engineer)')
  ) AND
  j.posted_at >= NOW() - INTERVAL '30 days' AND
  j.work_mode = 'Remote' AND
  (t.name = 'React' OR t.name = 'JavaScript')
ORDER BY j.posted_at DESC
LIMIT 20;
```

2. **Advanced Technology Search (With Hierarchy)**
```sql
WITH RECURSIVE tech_tree AS (
  SELECT id FROM technologies WHERE name = 'JavaScript'
  
  UNION
  
  SELECT t.id
  FROM technologies t
  JOIN tech_tree tt ON t.parent_id = tt.id
)
SELECT DISTINCT j.id, j.title, j.description
FROM jobs j
JOIN job_technologies jt ON j.id = jt.job_id
WHERE 
  j.is_active = TRUE AND
  jt.technology_id IN (SELECT id FROM tech_tree)
ORDER BY j.posted_at DESC
LIMIT 20;
```

### Indexing Strategies
The provided schema includes well-designed indexes for common query patterns. Additional considerations:

1. **Consider adding composite indexes** for common filter combinations:
```sql
CREATE INDEX idx_jobs_composite_filters ON jobs(is_active, work_mode, experience_level, employment_type);
```

2. **Consider adding function-based indexes** for case-insensitive searches:
```sql
CREATE INDEX idx_companies_name_lower ON companies(LOWER(name));
```

### Pagination Best Practices
For better performance with large datasets:

1. **Use keyset pagination** (recommended for chronological data):
```sql
-- First page
SELECT * FROM jobs 
WHERE is_active = TRUE
ORDER BY posted_at DESC, id DESC
LIMIT 20;

-- Next page (assuming last row had posted_at='2024-03-15 12:34:56' and id=789)
SELECT * FROM jobs 
WHERE 
  is_active = TRUE AND
  (posted_at < '2024-03-15 12:34:56' OR 
   (posted_at = '2024-03-15 12:34:56' AND id < 789))
ORDER BY posted_at DESC, id DESC
LIMIT 20;
```

2. **Use count estimates** rather than exact counts for pagination display:
```sql
-- Faster approximate count
SELECT count(*) FROM jobs 
WHERE is_active = TRUE AND work_mode = 'Remote'
LIMIT 1000;  -- Cap count for very large tables
```

### Query Optimization Tips
1. Always use the proper indexes when filtering
2. Use `EXPLAIN ANALYZE` to verify query performance
3. Consider materialized views for complex, frequently-run reports
4. Implement connection pooling for production environments
5. Consider partitioning the `jobs` table by date for very large job boards
