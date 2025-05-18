package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/rodruizronald/ticos-in-tech/internal/company"
	"github.com/rodruizronald/ticos-in-tech/internal/database"
)

type Company struct {
	Name    string
	LogoURL string
}

var companies = []Company{
	{Name: "Growth Acceleration Partners", LogoURL: "https://example.com/logo1.png"},
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
	dbpool, err := database.Connect(ctx, dbConfig)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbpool.Close()

	// Create a company repository
	repo := company.NewRepository(dbpool)

	// Store each company in the database
	for _, c := range companies {
		cm := &company.Company{
			Name:     c.Name,
			LogoURL:  c.LogoURL,
			IsActive: true,
		}

		err := repo.Create(ctx, cm)
		if err != nil {
			if company.IsDuplicate(err) {
				log.Infof("Company already exists: %s", cm.Name)
				continue
			}
			log.Warnf("Error creating company %s: %v", c.Name, err)
			continue
		}

		log.Infof("Successfully added company: %s (ID: %d)", cm.Name, cm.ID)
	}

	log.Info("Company population completed")
}
