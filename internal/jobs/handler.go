package jobs

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/rodruizronald/ticos-in-tech/internal/jobtech"
)

// DataRepository interface to make database operations for the Job model.
type DataRepository interface {
	SearchJobsWithCount(ctx context.Context, params *SearchParams) ([]*JobWithCompany, int, error)
	GetJobTechnologiesBatch(ctx context.Context, jobIDs []int) (map[int][]*jobtech.JobTechnologyWithDetails, error)
}

// Repositories struct to hold repositories for job and jobtech models
type Repositories struct {
	jobRepo     *Repository
	jobtechRepo *jobtech.Repository
}

// SearchJobsWithCount delegates to the job repository's SearchJobsWithCount method
func (r *Repositories) SearchJobsWithCount(ctx context.Context, params *SearchParams) ([]*JobWithCompany, int, error) {
	return r.jobRepo.SearchJobsWithCount(ctx, params)
}

// GetJobTechnologiesBatch delegates to the jobtech repository's GetJobTechnologiesBatch method
func (r *Repositories) GetJobTechnologiesBatch(ctx context.Context, jobIDs []int) (
	map[int][]*jobtech.JobTechnologyWithDetails, error) {
	return r.jobtechRepo.GetJobTechnologiesBatch(ctx, jobIDs)
}

// Handler handles HTTP requests for job operations
type Handler struct {
	repos DataRepository
}

// NewRepositories creates a new job and jobtech repositories
func NewRepositories(jobRepo *Repository, jobtechRepo *jobtech.Repository) *Repositories {
	return &Repositories{jobRepo: jobRepo, jobtechRepo: jobtechRepo}
}

// NewHandler creates a new job handler
func NewHandler(repos DataRepository) *Handler {
	return &Handler{repos: repos}
}

// RegisterRoutes registers job routes with the given router group
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET(JobsRoute, h.SearchJobs)
}

// SearchJobs godoc
// @Summary Search for jobs
// @Description Search for jobs with optional filters and pagination
// @Tags jobs
// @Accept json
// @Produce json
// @Param q query string true "Search query" example("golang developer")
// @Param limit query int false "Number of results to return (max 100)" default(20) example(20)
// @Param offset query int false "Number of results to skip" default(0) example(0)
// @Param experience_level query string false "Experience level filter" \
// Enums(Entry-level,Junior,Mid-level,Senior,Lead,Principal,Executive) example("Senior")
// @Param employment_type query string false "Employment type filter" \
// Enums(Full-time,Part-time,Contract,Freelance,Temporary,Internship) example("Full-time")
// @Param location query string false "Location filter" Enums(Costa Rica,LATAM) example("Costa Rica")
// @Param work_mode query string false "Work mode filter" Enums(Remote,Hybrid,Onsite) example("Remote")
// @Param date_from query string false "Start date filter (YYYY-MM-DD)" example("2024-01-01")
// @Param date_to query string false "End date filter (YYYY-MM-DD)" example("2024-12-31")
// @Success 200 {object} SearchResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /jobs [get]
func (h *Handler) SearchJobs(c *gin.Context) {
	var req SearchRequest

	// Bind query parameters
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetails{
				Code:    ErrCodeInvalidRequest,
				Message: "Invalid request parameters",
				Details: []string{err.Error()},
			},
		})
		return
	}

	// Validate request
	if validationErrors := validateSearchRequest(&req); len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetails{
				Code:    ErrCodeValidationError,
				Message: "Invalid search parameters",
				Details: validationErrors,
			},
		})
		return
	}

	// Convert request to search parameters
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

	// Parse dates if provided
	if req.DateFrom != "" && req.DateTo != "" {
		dateFrom, _ := time.Parse("2006-01-02", req.DateFrom)
		dateTo, _ := time.Parse("2006-01-02", req.DateTo)
		searchParams.DateFrom = &dateFrom
		searchParams.DateTo = &dateTo
	}

	// Perform search with count in single query
	jobs, total, err := h.repos.SearchJobsWithCount(c.Request.Context(), searchParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: ErrorDetails{
				Code:    ErrCodeSearchError,
				Message: "Failed to search jobs",
				Details: []string{err.Error()},
			},
		})
		return
	}

	// Get job IDs for batch fetching technologies
	jobIDs := make([]int, len(jobs))
	for i, job := range jobs {
		jobIDs[i] = job.ID
	}

	// Batch fetch technologies for all jobs
	technologiesMap, err := h.repos.GetJobTechnologiesBatch(c.Request.Context(), jobIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: ErrorDetails{
				Code:    ErrCodeSearchError,
				Message: "Failed to fetch job technologies",
				Details: []string{err.Error()},
			},
		})
		return
	}

	// Convert jobs to response format with technologies
	companyJobsResponseList := GroupJobsByCompany(jobs, technologiesMap)

	// Build response with correct pagination
	hasMore := searchParams.Offset+len(jobs) < total

	response := SearchResponse{
		Data: companyJobsResponseList,
		Pagination: PaginationDetails{
			Total:   total,
			Limit:   searchParams.Limit,
			Offset:  searchParams.Offset,
			HasMore: hasMore,
		},
	}

	c.JSON(http.StatusOK, response)
}
