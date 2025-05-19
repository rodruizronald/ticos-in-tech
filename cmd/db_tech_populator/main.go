// Package main provides a utility to populate the database with technology information.
// It reads technology data and inserts them into the database.
package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/rodruizronald/ticos-in-tech/internal/database"
	"github.com/rodruizronald/ticos-in-tech/internal/techalias"
	"github.com/rodruizronald/ticos-in-tech/internal/technology"
)

// Technology represents a technology entity as stored in the configuration.
// It contains the basic information needed to create a technology record in the database
type Technology struct {
	Name     string   `json:"name"`
	Category string   `json:"category"`
	Alias    []string `json:"alias"`
	Parent   string   `json:"parent"`
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
	// Configure logger
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Get database config
	dbConfig := database.DefaultConfig()

	log.Infof("Connecting to database %s at %s:%d", dbConfig.DBName, dbConfig.Host, dbConfig.Port)

	// Connect to the database
	dbpool, err := database.Connect(ctx, &dbConfig)
	if err != nil {
		log.Errorf("Unable to connect to database: %v", err)
	}
	defer dbpool.Close()

	// Create repositories
	techRepo := technology.NewRepository(dbpool)
	aliasRepo := techalias.NewRepository(dbpool)

	// Process technologies
	processTechnologies(ctx, log, techRepo, aliasRepo)

	log.Info("Technology import completed")
	return nil
}

// processTechnologies handles the two-pass technology import process
func processTechnologies(ctx context.Context, log *logrus.Logger, techRepo *technology.Repository, aliasRepo *techalias.Repository) {
	// Create a map to store all technologies by name for lookup
	techMap := make(map[string]*technology.Technology)

	// Process and insert all technologies
	technologies := readTechnologiesFromJSON()

	// First pass: create technologies without parent references
	log.Info("Starting first pass: creating technologies without parent references")
	createTechnologies(ctx, log, techRepo, aliasRepo, technologies, techMap)

	// Second pass: update technologies with parent references
	log.Info("Starting second pass: updating technologies with parent references")
	updateTechnologyParents(ctx, log, techRepo, technologies, techMap)
}

// createTechnologies handles the first pass of creating technologies
func createTechnologies(ctx context.Context, log *logrus.Logger, techRepo *technology.Repository,
	aliasRepo *techalias.Repository, technologies [][]Technology, techMap map[string]*technology.Technology) {

	for _, techGroup := range technologies {
		for _, tech := range techGroup {
			// Convert name to lowercase
			techName := strings.ToLower(tech.Name)

			// Create the technology model
			newTech := &technology.Technology{
				Name:     techName,
				Category: tech.Category,
				// Parent ID will be set in the second pass
			}

			// Insert into database
			err := techRepo.Create(ctx, newTech)
			if err != nil {
				// Skip if it's a duplicate
				if technology.IsDuplicate(err) {
					log.Infof("Technology already exists: %s", techName)

					// Fetch the existing technology to use for parent mapping
					existingTech, err := techRepo.GetByName(ctx, techName)
					if err != nil {
						log.Warnf("Error fetching existing technology %s: %v", techName, err)
						continue
					}
					techMap[techName] = existingTech

					// Add aliases for existing technology
					addAliases(ctx, log, aliasRepo, existingTech.ID, tech.Alias)
					continue
				}
				log.Warnf("Error creating technology %s: %v", techName, err)
				continue
			}

			log.Infof("Created technology: %s (ID: %d)", techName, newTech.ID)
			techMap[techName] = newTech

			// Add aliases for new technology
			addAliases(ctx, log, aliasRepo, newTech.ID, tech.Alias)
		}
	}
}

// updateTechnologyParents handles the second pass of updating parent references
func updateTechnologyParents(ctx context.Context, log *logrus.Logger, techRepo *technology.Repository,
	technologies [][]Technology, techMap map[string]*technology.Technology) {

	for _, techGroup := range technologies {
		for _, tech := range techGroup {
			if tech.Parent == "" {
				continue // Skip technologies without parents
			}
			techName := strings.ToLower(tech.Name)
			parentName := strings.ToLower(tech.Parent)

			// Look up the current technology
			currentTech, exists := techMap[techName]
			if !exists {
				log.Warnf("Cannot find technology: %s", techName)
				continue
			}

			// Look up the parent technology
			parentTech, exists := techMap[parentName]
			if !exists {
				log.Warnf("Cannot find parent technology: %s for %s", parentName, techName)
				continue
			}

			// Update the parent ID
			currentTech.ParentID = &parentTech.ID
			err := techRepo.Update(ctx, currentTech)
			if err != nil {
				log.Warnf("Error updating parent for %s: %v", currentTech.Name, err)
				continue
			}

			log.Infof("Updated technology %s with parent %s (ID: %d)",
				currentTech.Name, parentTech.Name, parentTech.ID)
		}
	}
}

// addAliases adds aliases for a technology
func addAliases(ctx context.Context, log *logrus.Logger, aliasRepo *techalias.Repository, techID int, aliases []string) {
	for _, aliasName := range aliases {
		if aliasName == "" {
			continue
		}

		// Convert alias to lowercase
		lowerAlias := strings.ToLower(aliasName)

		// Create alias model
		newAlias := &techalias.TechnologyAlias{
			TechnologyID: techID,
			Alias:        lowerAlias,
		}

		// Insert into database
		err := aliasRepo.Create(ctx, newAlias)
		if err != nil {
			// Skip if it's a duplicate
			if techalias.IsDuplicate(err) {
				log.Infof("Alias already exists: %s", lowerAlias)
				continue
			}
			log.Warnf("Error creating alias %s for technology ID %d: %v", lowerAlias, techID, err)
			continue
		}

		log.Infof("Created alias: %s (ID: %d) for technology ID %d", lowerAlias, newAlias.ID, techID)
	}
}

// readTechnologiesFromJSON reads technology data from a JSON file
func readTechnologiesFromJSON() [][]Technology {
	// Get the directory of the current executable
	execDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logrus.Errorf("Failed to get executable directory: %v", err)
		return [][]Technology{}
	}

	// Path to the JSON file
	jsonPath := filepath.Join(execDir, "technologies.json")

	// For development, if the file doesn't exist in the executable directory,
	// try looking in the current directory
	if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
		jsonPath = "technologies.json"
	}

	// Read the JSON file
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		logrus.Errorf("Failed to read technologies file: %v", err)
		return [][]Technology{}
	}

	// Parse the JSON data
	var technologies [][]Technology
	if err := json.Unmarshal(data, &technologies); err != nil {
		logrus.Errorf("Failed to parse technologies JSON: %v", err)
		return [][]Technology{}
	}

	return technologies
}
