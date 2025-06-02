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

	"github.com/rodruizronald/ticos-in-tech/internal/company"
	"github.com/rodruizronald/ticos-in-tech/internal/database"
	"github.com/rodruizronald/ticos-in-tech/internal/job"
	"github.com/rodruizronald/ticos-in-tech/internal/jobtech"
	"github.com/rodruizronald/ticos-in-tech/internal/techalias"
	"github.com/rodruizronald/ticos-in-tech/internal/technology"
	"github.com/sirupsen/logrus"
)

type JobData struct {
	Jobs []struct {
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
	} `json:"jobs"`
}

func main() {
	// Initialize logger
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Set up database connection
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Get database config
	dbConfig := database.DefaultConfig()

	// Connect to the database
	dbpool, err := database.Connect(ctx, &dbConfig)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbpool.Close()

	// Create repositories
	jobRepo := job.NewRepository(dbpool)
	companyRepo := company.NewRepository(dbpool)
	jobtechRepo := jobtech.NewRepository(dbpool)
	techRepo := technology.NewRepository(dbpool)
	aliasRepo := techalias.NewRepository(dbpool)

	// Get today's date for the input directory
	today := time.Now().Format("20060102")
	inputDir := filepath.Join("data", today)
	inputFile := filepath.Join(inputDir, "jobs.json")
	missingTechFile := filepath.Join(inputDir, "missing_technologies.json")

	// Create a map to track missing technologies
	missingTechnologies := make(map[string][]string) // company -> list of missing tech names

	log.Infof("Reading job data from %s", inputFile)

	// Read job data from file
	data, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Failed to read job data file: %v", err)
	}

	// Parse job data
	var jobData JobData
	if err := json.Unmarshal(data, &jobData); err != nil {
		log.Fatalf("Failed to parse job data: %v", err)
	}

	log.Infof("Found %d jobs to process", len(jobData.Jobs))

	// Process each job
	for _, j := range jobData.Jobs {
		// Find company by name
		company, err := companyRepo.GetByName(ctx, j.Company)
		if err != nil {
			log.Warnf("Error finding company %s: %v", j.Company, err)
			continue
		}

		companyID := company.ID

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

		// Insert job into database
		err = jobRepo.Create(ctx, jobModel)
		if err != nil {
			if job.IsDuplicate(err) {
				log.Infof("Job already exists: %s at %s", j.Title, j.Company)
				
				// Get the existing job by signature to retrieve its ID
				existingJob, findErr := jobRepo.GetBySignature(ctx, j.Signature)
				if findErr != nil {
					log.Warnf("Failed to retrieve existing job %s: %v", j.Title, findErr)
					continue
				}
				
				// Use the existing job's ID for technology associations
				jobModel.ID = existingJob.ID
				log.Infof("Using existing job ID: %d", jobModel.ID)
			} else {
				log.Warnf("Failed to insert job %s: %v", j.Title, err)
				continue
			}
		}

		log.Infof("Successfully added job: %s at %s (ID: %d)",
			jobModel.Title, j.Company, jobModel.ID)

		for _, tech := range j.Technologies {
			techName := strings.ToLower(tech.Name)

			// Find technology by name
			techModel, err := techRepo.GetByName(ctx, techName)
			if err != nil {
				// If not found by exact name, try to find by alias
				alias, aliasErr := aliasRepo.GetByAlias(ctx, techName)

				if aliasErr != nil {
					log.Warnf("Technology not found by name or alias: %s: %v", techName, err)
					missingTechnologies[j.Company] = append(missingTechnologies[j.Company], techName)
					continue
				}

				// Get the technology using the alias's technology ID
				techModel, err = techRepo.GetByID(ctx, alias.TechnologyID)
				if err != nil {
					log.Warnf("Error finding technology by alias ID %d: %v", alias.TechnologyID, err)
					continue
				}

				log.Infof("Found technology %s via alias %s", techModel.Name, techName)
			}

			// Create job technology association
			jobTechModel := &jobtech.JobTechnology{
				JobID:        jobModel.ID,
				TechnologyID: techModel.ID,
				IsRequired:   tech.Required,
			}

			// Insert job technology into database
			err = jobtechRepo.Create(ctx, jobTechModel)
			if err != nil {
				if jobtech.IsDuplicate(err) {
					log.Debugf("Job technology association already exists: %s for job ID %d", techName, jobModel.ID)
					continue
				}
				log.Warnf("Failed to insert job technology %s: %v", techName, err)
				continue
			}

			log.Infof("Added technology %s to job ID %d", techName, jobModel.ID)
		}
	}

	// Write missing technologies to file
	if len(missingTechnologies) > 0 {
		missingTechData, err := json.MarshalIndent(missingTechnologies, "", "  ")
		if err != nil {
			log.Fatalf("Failed to marshal missing technologies: %v", err)
		}

		err = os.WriteFile(missingTechFile, missingTechData, 0644)
		if err != nil {
			log.Fatalf("Failed to write missing technologies file: %v", err)
		}

		log.Infof("Missing technologies saved to %s", missingTechFile)
	}

	log.Info("Job population completed")
}
