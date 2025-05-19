package job

import (
	"time"

	"github.com/rodruizronald/ticos-in-tech/internal/jobtech"
)

// Job represents a job posting on the platform.
type Job struct {
	ID              int       `json:"id" db:"id"`
	CompanyID       int       `json:"company_id" db:"company_id"`
	Title           string    `json:"title" db:"title"`
	Description     string    `json:"description" db:"description"`
	ExperienceLevel string    `json:"experience_level" db:"experience_level"`
	EmploymentType  string    `json:"employment_type" db:"employment_type"`
	Location        string    `json:"location" db:"location"`
	WorkMode        string    `json:"work_mode" db:"work_mode"`
	ApplicationURL  string    `json:"application_url" db:"application_url"`
	IsActive        bool      `json:"is_active" db:"is_active"`
	Signature       string    `json:"signature" db:"signature"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`

	// Relationships (not stored in database)
	Technologies []jobtech.JobTechnology `json:"technologies,omitempty" db:"-"`
}
