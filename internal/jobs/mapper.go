package jobs

import "github.com/rodruizronald/ticos-in-tech/internal/jobtech"

// Mapping functions to convert between database and API models.
// This file contains transformation logic that bridges the repository layer (database models)
// and the API layer (DTOs). Mappers handle the conversion of internal data structures
// to external API representations, including data aggregation and formatting.

// MapJobToResponse converts a single job with company data to API response format.
// It transforms a database model into a DTO suitable for API responses.
func MapJobToResponse(job *JobWithCompany, technologies []TechnologyResponse) *JobResponse {
	return &JobResponse{
		ID:              job.ID,
		CompanyID:       job.CompanyID,
		Title:           job.Title,
		Description:     job.Description,
		ExperienceLevel: job.ExperienceLevel,
		EmploymentType:  job.EmploymentType,
		Location:        job.Location,
		WorkMode:        job.WorkMode,
		ApplicationURL:  job.ApplicationURL,
		Technologies:    technologies,
		PostedAt:        job.CreatedAt,
	}
}

// GroupJobsByCompany groups jobs by company and converts them to API response format.
// It takes jobs with company data and technologies, organizing them by company
// and transforming database models into CompanyJobsResponse DTOs.
func GroupJobsByCompany(jobs []*JobWithCompany, techMap map[int][]*jobtech.JobTechnologyWithDetails) []*CompanyJobsResponse {
	// Group jobs by company
	companyJobsMap := make(map[int]*CompanyJobsResponse)

	for _, job := range jobs {
		if _, exists := companyJobsMap[job.CompanyID]; !exists {
			companyJobsMap[job.CompanyID] = &CompanyJobsResponse{
				CompanyName:    job.CompanyName,
				CompanyLogoURL: job.CompanyLogoURL,
				Jobs:           []*JobResponse{},
			}
		}

		// Convert technologies for this job
		jobTechnologies := techMap[job.ID]
		technologies := make([]TechnologyResponse, len(jobTechnologies))
		for i, tech := range jobTechnologies {
			technologies[i] = TechnologyResponse{
				Name:     tech.TechName,
				Category: tech.TechCategory,
				Required: tech.IsRequired,
			}
		}

		// Use the single job mapper
		jobResponse := MapJobToResponse(job, technologies)
		companyJobsMap[job.CompanyID].Jobs = append(companyJobsMap[job.CompanyID].Jobs, jobResponse)
	}

	// Convert map to slice
	companyJobsResponseList := make([]*CompanyJobsResponse, 0, len(companyJobsMap))
	for _, companyJobsResp := range companyJobsMap {
		companyJobsResponseList = append(companyJobsResponseList, companyJobsResp)
	}

	return companyJobsResponseList
}
