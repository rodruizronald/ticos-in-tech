package jobs

import "time"

// Database entities and repository-level structs for job management.
// This file contains the core database models and search parameters used by the repository layer.
// These models map directly to database tables and are used for data persistence operations.

// Job represents the database entity
type Job struct {
	ID              int       `db:"id"`
	CompanyID       int       `db:"company_id"`
	Title           string    `db:"title"`
	Description     string    `db:"description"`
	ExperienceLevel string    `db:"experience_level"`
	EmploymentType  string    `db:"employment_type"`
	Location        string    `db:"location"`
	WorkMode        string    `db:"work_mode"`
	ApplicationURL  string    `db:"application_url"`
	IsActive        bool      `db:"is_active"`
	Signature       string    `db:"signature"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

// JobWithCompany represents a job with company details (for read operations only)
type JobWithCompany struct {
	Job                   // Embed the original Job struct
	CompanyName    string `db:"company_name"`
	CompanyLogoURL string `db:"company_logo_url"`
}

// SearchParams defines parameters for job search (repository layer)
type SearchParams struct {
	Query  string
	Limit  int
	Offset int
	// Optional filters
	ExperienceLevel *string
	EmploymentType  *string
	Location        *string
	WorkMode        *string
	DateFrom        *time.Time
	DateTo          *time.Time
}
