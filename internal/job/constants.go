package job

// SQL query constants
const (
	// Base query for selecting job fields
	selectJobBaseQuery = `
        SELECT id, company_id, title, description, experience_level, employment_type,
               location, work_mode, application_url, is_active, signature, created_at, updated_at
        FROM jobs
    `

	createJobQuery = `
        INSERT INTO jobs (
            company_id, title, description, experience_level, employment_type,
            location, work_mode, application_url, is_active, signature
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING id, created_at, updated_at
    `

	getJobByIDQuery = selectJobBaseQuery + `
        WHERE id = $1
    `

	getJobBySignatureQuery = selectJobBaseQuery + `
        WHERE signature = $1
    `

	updateJobQuery = `
        UPDATE jobs
        SET company_id = $1, title = $2, description = $3, experience_level = $4,
            employment_type = $5, location = $6, work_mode = $7, application_url = $8,
            is_active = $9, signature = $10, updated_at = NOW()
        WHERE id = $11
        RETURNING updated_at
    `

	deleteJobQuery = `DELETE FROM jobs WHERE id = $1`

	// Text search query with weighted ranking (title has more weight than description)
	searchJobsBaseQuery = `
        WITH search_query AS (
            SELECT plainto_tsquery('english', $1) AS query
        )
        SELECT 
            j.id, j.company_id, j.title, j.description, j.experience_level, j.employment_type,
            j.location, j.work_mode, j.application_url, j.is_active, j.signature, j.created_at, j.updated_at
        FROM jobs j, search_query sq
        WHERE j.is_active = true AND j.search_vector @@ sq.query
	`
)

// Constants for job attributes and values
const (
	// Experience levels
	experienceLevelEntry     = "Entry-level"
	experienceLevelJunior    = "Junior"
	experienceLevelMid       = "Mid-level"
	experienceLevelSenior    = "Senior"
	experienceLevelLead      = "Lead"
	experienceLevelPrincipal = "Principal"
	experienceLevelExecutive = "Executive"

	// Employment types
	employmentTypeFullTime   = "Full-time"
	employmentTypePartTime   = "Part-time"
	employmentTypeContract   = "Contract"
	employmentTypeFreelance  = "Freelance"
	employmentTypeTemporary  = "Temporary"
	employmentTypeInternship = "Internship"

	// Locations
	locationCostaRica = "Costa Rica"
	locationLATAM     = "LATAM"

	// Work modes
	workModeRemote = "Remote"
	workModeHybrid = "Hybrid"
	workModeOnsite = "Onsite"
)

// Validation collections for job attributes and values
var (
	validExperienceLevels = []string{
		experienceLevelEntry,
		experienceLevelJunior,
		experienceLevelMid,
		experienceLevelSenior,
		experienceLevelLead,
		experienceLevelPrincipal,
		experienceLevelExecutive,
	}
	validEmploymentTypes = []string{
		employmentTypeFullTime,
		employmentTypePartTime,
		employmentTypeContract,
		employmentTypeFreelance,
		employmentTypeTemporary,
		employmentTypeInternship,
	}
	validLocations = []string{
		locationCostaRica,
		locationLATAM,
	}
	validWorkModes = []string{
		workModeRemote,
		workModeHybrid,
		workModeOnsite,
	}
)

// Constants for job routes and endpoints
const (
	JobsRoute = "/jobs"
)

// Constants for error codes and messages
const (
	ErrCodeInvalidRequest  = "INVALID_REQUEST"
	ErrCodeValidationError = "VALIDATION_ERROR"
	ErrCodeSearchError     = "SEARCH_ERROR"
)
