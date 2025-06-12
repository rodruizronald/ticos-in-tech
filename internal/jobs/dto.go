package jobs

import "time"

// Data Transfer Objects (DTOs) for the job API layer.
// This file contains request/response structures used for HTTP API communication.
// These models define the external API contract and handle JSON serialization/deserialization.
// They are decoupled from database models to allow independent evolution of API and database schemas.

// SearchRequest represents the search request parameters (API layer)
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

// JobResponse represents the API response for a single job
type JobResponse struct {
	ID              int                  `json:"job_id"`
	CompanyID       int                  `json:"company_id"`
	Title           string               `json:"title"`
	Description     string               `json:"description"`
	ExperienceLevel string               `json:"experience_level"`
	EmploymentType  string               `json:"employment_type"`
	Location        string               `json:"location"`
	WorkMode        string               `json:"work_mode"`
	ApplicationURL  string               `json:"application_url"`
	Technologies    []TechnologyResponse `json:"technologies"`
	PostedAt        time.Time            `json:"posted_at"`
}

// TechnologyResponse represents the API response for job technologies
type TechnologyResponse struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	Required bool   `json:"required"`
}

// CompanyJobsResponse represents grouped jobs by company
type CompanyJobsResponse struct {
	CompanyName    string         `json:"company_name"`
	CompanyLogoURL string         `json:"company_logo_url"`
	Jobs           []*JobResponse `json:"jobs"`
}

// SearchResponse represents the search response with pagination
type SearchResponse struct {
	Data       []*CompanyJobsResponse `json:"data"`
	Pagination PaginationDetails      `json:"pagination"`
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
