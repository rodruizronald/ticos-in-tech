package jobs

import (
	"fmt"
	"time"
)

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
	Company         string `form:"company" example:"Tech Corp"`
	DateFrom        string `form:"date_from" example:"2024-01-01"`
	DateTo          string `form:"date_to" example:"2024-12-31"`
}

// ToSearchParams converts a SearchRequest to SearchParams
func (req *SearchRequest) ToSearchParams() (*SearchParams, error) {
	searchParams := &SearchParams{
		Query:  req.Query,
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	// Set optional filters
	if req.ExperienceLevel != "" {
		searchParams.ExperienceLevel = &req.ExperienceLevel
	}
	if req.EmploymentType != "" {
		searchParams.EmploymentType = &req.EmploymentType
	}
	if req.Location != "" {
		searchParams.Location = &req.Location
	}
	if req.WorkMode != "" {
		searchParams.WorkMode = &req.WorkMode
	}
	if req.Company != "" {
		searchParams.Company = &req.Company
	}

	// Parse dates if provided
	if req.DateFrom != "" && req.DateTo != "" {
		dateFrom, err := time.Parse("2006-01-02", req.DateFrom)
		if err != nil {
			return nil, fmt.Errorf("invalid date_from format: %w", err)
		}
		dateTo, err := time.Parse("2006-01-02", req.DateTo)
		if err != nil {
			return nil, fmt.Errorf("invalid date_to format: %w", err)
		}
		searchParams.DateFrom = &dateFrom
		searchParams.DateTo = &dateTo
	}

	return searchParams, nil
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

// SearchResponse represents the search response with pagination
type SearchResponse struct {
	Data       []*JobResponse    `json:"data"`
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
