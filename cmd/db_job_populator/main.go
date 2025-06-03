// Package main provides a utility to populate the database with job information.
// It reads job data from JSON files and inserts them into the database along with
// their associated technologies.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"

	"github.com/rodruizronald/ticos-in-tech/internal/company"
	"github.com/rodruizronald/ticos-in-tech/internal/database"
	"github.com/rodruizronald/ticos-in-tech/internal/job"
	"github.com/rodruizronald/ticos-in-tech/internal/jobtech"
	"github.com/rodruizronald/ticos-in-tech/internal/techalias"
	"github.com/rodruizronald/ticos-in-tech/internal/technology"
)

// Job define a type to represent a single job
type jobData struct {
	Company         string `json:"company"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	ApplicationURL  string `json:"application_url"`
	Location        string `json:"location"`
	WorkMode        string `json:"work_mode"`
	ExperienceLevel string `json:"experience_level"`
	EmploymentType  string `json:"employment_type"`
	Technologies    []struct {
		Name     string `json:"name"`
		Category string `json:"category"`
		Required bool   `json:"required"`
	} `json:"technologies"`
	Signature string `json:"signature"`
}

// Update the jobs struct to use the Job type
type jobs struct {
	Jobs []jobData `json:"jobs"`
}

func main() {
	var err error
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		stop()
		if err != nil {
			os.Exit(1)
		}
	}()
	err = run(ctx)
}

func run(ctx context.Context) error {
	// Initialize logger
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Setup database and repositories
	dbpool, repos, err := setupDatabase(ctx, log)
	if err != nil {
		return err
	}
	defer dbpool.Close()

	// Get file paths
	today := time.Now().Format("20060102")
	inputDir := filepath.Join("data", today)
	inputFile := filepath.Join(inputDir, "jobs.json")
	missingTechFile := filepath.Join(inputDir, "missing_technologies.json")

	// Read and parse job data
	jobData, err := readJobData(inputFile, log)
	if err != nil {
		return err
	}

	// Process jobs and collect missing technologies
	missingTechnologies, err := processJobs(ctx, jobData, repos, log)
	if err != nil {
		return err
	}

	// Write missing technologies to file if any
	if err := writeMissingTechnologies(missingTechnologies, missingTechFile, log); err != nil {
		return err
	}

	log.Info("Job population completed")
	return nil
}

// setupDatabase initializes the database connection and repositories
func setupDatabase(ctx context.Context, log *logrus.Logger) (*pgxpool.Pool, *repositories, error) {
	// Get database config
	dbConfig := database.DefaultConfig()

	// Connect to the database
	dbpool, err := database.Connect(ctx, &dbConfig)
	if err != nil {
		log.Errorf("Unable to connect to database: %v", err)
		return nil, nil, err
	}

	// Create repositories
	repos := &repositories{
		job:     job.NewRepository(dbpool),
		company: company.NewRepository(dbpool),
		jobtech: jobtech.NewRepository(dbpool),
		tech:    technology.NewRepository(dbpool),
		alias:   techalias.NewRepository(dbpool),
	}

	return dbpool, repos, nil
}

// repositories holds all the database repositories needed
type repositories struct {
	job     *job.Repository
	company *company.Repository
	jobtech *jobtech.Repository
	tech    *technology.Repository
	alias   *techalias.Repository
}

// readJobData reads and parses the job data from the input file
func readJobData(inputFile string, log *logrus.Logger) (*jobs, error) {
	log.Infof("Reading job data from %s", inputFile)

	// Read job data from file
	data, err := os.ReadFile(inputFile)
	if err != nil {
		log.Errorf("Failed to read job data file: %v", err)
		return nil, err
	}

	// Parse job data
	var jobData jobs
	if err := json.Unmarshal(data, &jobData); err != nil {
		log.Errorf("Failed to parse job data: %v", err)
		return nil, err
	}

	log.Infof("Found %d jobs to process", len(jobData.Jobs))
	return &jobData, nil
}

// processJobs processes each job and returns a map of missing technologies
func processJobs(ctx context.Context, jobData *jobs, repos *repositories,
	log *logrus.Logger) (map[string][]string, error) {
	// Create a map to track missing technologies
	missingTechnologies := make(map[string][]string) // company -> list of missing tech names

	// Process each job
	for i := range jobData.Jobs {
		j := &jobData.Jobs[i] // Use a pointer to the job instead of copying it

		// Process job and its technologies
		jobMissingTechs, err := processJob(ctx, j, repos, log)
		if err != nil {
			// Log error but continue with next job
			log.Warnf("Error processing job %s: %v", j.Title, err)
			continue
		}

		// Add any missing technologies to the map
		if len(jobMissingTechs) > 0 {
			missingTechnologies[j.Company] = append(missingTechnologies[j.Company], jobMissingTechs...)
		}
	}

	return missingTechnologies, nil
}

// Update the processJob function signature
func processJob(ctx context.Context, j *jobData, repos *repositories, log *logrus.Logger) ([]string, error) {
	// Find company by name
	jobCompany, err := repos.company.GetByName(ctx, j.Company)
	if err != nil {
		log.Warnf("Error finding company %s: %v", j.Company, err)
		return nil, err
	}

	companyID := jobCompany.ID

	// Create job model
	jobModel := &job.Job{
		CompanyID:       companyID,
		Title:           j.Title,
		Description:     j.Description,
		ExperienceLevel: j.ExperienceLevel,
		EmploymentType:  j.EmploymentType,
		Location:        j.Location,
		WorkMode:        j.WorkMode,
		ApplicationURL:  j.ApplicationURL,
		IsActive:        true,
		Signature:       j.Signature,
	}
	fmt.Print("Processing job: ", jobModel.Title, " at ", j.Company, "\n")

	// Insert or retrieve job
	if err := createOrRetrieveJob(ctx, jobModel, j, repos.job, log); err != nil {
		return nil, err
	}

	log.Infof("Successfully added job: %s at %s (ID: %d)",
		jobModel.Title, j.Company, jobModel.ID)

	// Process technologies for this job
	return processTechnologies(ctx, j, jobModel, repos, log)
}

// createOrRetrieveJob creates a new job or retrieves an existing one
func createOrRetrieveJob(ctx context.Context, jobModel *job.Job, j *jobData, jobRepo *job.Repository,
	log *logrus.Logger) error {
	err := jobRepo.Create(ctx, jobModel)
	if err != nil {
		if job.IsDuplicate(err) {
			log.Infof("Job already exists: %s at %s", j.Title, j.Company)

			// Get the existing job by signature to retrieve its ID
			existingJob, findErr := jobRepo.GetBySignature(ctx, j.Signature)
			if findErr != nil {
				log.Warnf("Failed to retrieve existing job %s: %v", j.Title, findErr)
				return findErr
			}

			// Use the existing job's ID for technology associations
			jobModel.ID = existingJob.ID
			log.Infof("Using existing job ID: %d", jobModel.ID)
			return nil
		}
		log.Warnf("Failed to insert job %s: %v", j.Title, err)
		return err
	}
	return nil
}

// processTechnologies processes all technologies for a job
func processTechnologies(ctx context.Context, j *jobData, jobModel *job.Job, repos *repositories,
	log *logrus.Logger) ([]string, error) {
	var missingTechs []string

	for _, tech := range j.Technologies {
		techName := strings.ToLower(tech.Name)

		// Find technology by name or alias
		techModel, err := findTechnology(ctx, techName, repos, log)
		if err != nil {
			missingTechs = append(missingTechs, techName)
			continue
		}

		// Create job technology association
		if err := createJobTechnology(ctx, jobModel.ID, techModel.ID,
			tech.Required, techName, repos.jobtech, log); err != nil {
			continue
		}
	}

	return missingTechs, nil
}

// findTechnology tries to find a technology by name or alias
func findTechnology(ctx context.Context, techName string, repos *repositories,
	log *logrus.Logger) (*technology.Technology, error) {
	// Find technology by name
	techModel, err := repos.tech.GetByName(ctx, techName)
	if err == nil {
		return techModel, nil
	}

	// If not found by exact name, try to find by alias
	alias, aliasErr := repos.alias.GetByAlias(ctx, techName)
	if aliasErr != nil {
		log.Warnf("Technology not found by name or alias: %s: %v", techName, err)
		return nil, aliasErr
	}

	// Get the technology using the alias's technology ID
	techModel, err = repos.tech.GetByID(ctx, alias.TechnologyID)
	if err != nil {
		log.Warnf("Error finding technology by alias ID %d: %v", alias.TechnologyID, err)
		return nil, err
	}

	log.Infof("Found technology %s via alias %s", techModel.Name, techName)
	return techModel, nil
}

// createJobTechnology creates a job-technology association
func createJobTechnology(ctx context.Context, jobID, techID int, isRequired bool, techName string,
	jobtechRepo *jobtech.Repository, log *logrus.Logger) error {
	jobTechModel := &jobtech.JobTechnology{
		JobID:        jobID,
		TechnologyID: techID,
		IsRequired:   isRequired,
	}

	// Insert job technology into database
	err := jobtechRepo.Create(ctx, jobTechModel)
	if err != nil {
		if jobtech.IsDuplicate(err) {
			log.Debugf("Job technology association already exists: %s for job ID %d", techName, jobID)
			return nil
		}
		log.Warnf("Failed to insert job technology %s: %v", techName, err)
		return err
	}

	log.Infof("Added technology %s to job ID %d", techName, jobID)
	return nil
}

// writeMissingTechnologies writes missing technologies to a file
func writeMissingTechnologies(missingTechnologies map[string][]string,
	missingTechFile string, log *logrus.Logger) error {
	if len(missingTechnologies) == 0 {
		return nil
	}

	missingTechData, err := json.MarshalIndent(missingTechnologies, "", "  ")
	if err != nil {
		log.Errorf("Failed to marshal missing technologies: %v", err)
		return err
	}

	err = os.WriteFile(missingTechFile, missingTechData, 0o644)
	if err != nil {
		log.Errorf("Failed to write missing technologies file: %v", err)
		return err
	}

	log.Infof("Missing technologies saved to %s", missingTechFile)
	return nil
}
