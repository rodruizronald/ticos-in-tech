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

// SearchRequest represents the search request parameters
type SearchRequest struct {
	Query           string `form:"q" binding:"required" example:"golang developer"`
	Limit           int    `form:"limit" example:"20"`
	Offset          int    `form:"offset" example:"0"`
	ExperienceLevel string `form:"experience_level" example:"Senior"`
	EmploymentType  string `form:"employment_type" example:"Full-time"`
	Location        string `form:"location" example:"Costa Rica"`
	WorkMode        string `form:"work_mode" example:"Remote"`
	DateFrom        string `form:"date_from" example:"2024-01-01"`
	DateTo          string `form:"date_to" example:"2024-12-31"`
}

// SearchResponse represents the search response with pagination
type SearchResponse struct {
	Data       []*Job            `json:"data"`
	Pagination PaginationDetails `json:"pagination"`
}

// PaginationDetails contains pagination metadata
type PaginationDetails struct {
	Total   int  `json:"total"`
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	HasMore bool `json:"has_more"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error ErrorDetails `json:"error"`
}

// ErrorDetails contains error information
type ErrorDetails struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}
