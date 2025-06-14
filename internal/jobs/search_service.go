package jobs

import (
	"context"

	"github.com/rodruizronald/ticos-in-tech/internal/httpservice"
)

type JobSearchService struct {
	repos DataRepository
}

func NewJobSearchService(repos DataRepository) httpservice.SearchService[*SearchParams, JobResponseList] {
	return &JobSearchService{repos: repos}
}

func (s *JobSearchService) ExecuteSearch(ctx context.Context, params *SearchParams) (JobResponseList, int, error) {
	// Your existing business logic
	jobs, total, err := s.repos.SearchJobsWithCount(ctx, params)
	if err != nil {
		return nil, 0, &httpservice.SearchError{Operation: "search jobs", Err: err}
	}

	// Get job IDs for batch fetching technologies
	jobIDs := make([]int, len(jobs))
	for i, job := range jobs {
		jobIDs[i] = job.ID
	}

	// Batch fetch technologies for all jobs
	technologiesMap, err := s.repos.GetJobTechnologiesBatch(ctx, jobIDs)
	if err != nil {
		return nil, 0, &httpservice.SearchError{Operation: "fetch job technologies", Err: err}
	}

	// Convert jobs to response format with technologies
	searchResult := MapJobsToResponse(jobs, technologiesMap)

	return searchResult, total, nil
}
