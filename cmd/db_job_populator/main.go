package main

// type JobData struct {
//     Jobs []struct {
//         Company        string   `json:"company"`
//         Title          string   `json:"title"`
//         Description    string   `json:"description"`
//         ApplicationURL string   `json:"application_url"`
//         Location       string   `json:"location"`
//         WorkMode       string   `json:"work_mode"`
//         ExperienceLevel string  `json:"experience_level"`
//         EmploymentType string   `json:"employment_type"`
//         Technologies   []string `json:"technologies"`
//         Signature      string   `json:"signature"`
//     } `json:"jobs"`
// }

// func main() {
//     // Initialize logger
//     log := logrus.New()
//     log.SetFormatter(&logrus.TextFormatter{
//         FullTimestamp: true,
//     })

//     // Set up database connection
//     ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
//     defer stop()

//     // Get database config
//     dbConfig := database.DefaultConfig()

//     // Connect to the database
//     dbpool, err := database.Connect(ctx, dbConfig)
//     if err != nil {
//         log.Fatalf("Unable to connect to database: %v", err)
//     }
//     defer dbpool.Close()

//     // Create repositories
//     jobRepo := job.NewRepository(dbpool)
//     companyRepo := company.NewRepository(dbpool)

//     // Get today's date for the input directory
//     today := time.Now().Format("20060102")
//     inputDir := filepath.Join("scripts", "scraper", "output", today)
//     inputFile := filepath.Join(inputDir, "final_jobs.json")

//     log.Infof("Reading job data from %s", inputFile)

//     // Read job data from file
//     data, err := os.ReadFile(inputFile)
//     if err != nil {
//         log.Fatalf("Failed to read job data file: %v", err)
//     }

//     // Parse job data
//     var jobData JobData
//     if err := json.Unmarshal(data, &jobData); err != nil {
//         log.Fatalf("Failed to parse job data: %v", err)
//     }

//     log.Infof("Found %d jobs to process", len(jobData.Jobs))

//     // Process each job
//     for _, j := range jobData.Jobs {
//         // Find company by name
//         companies, err := companyRepo.List(ctx, company.CompanyFilter{Name: &j.Company})
//         if err != nil {
//             log.Warnf("Error finding company %s: %v", j.Company, err)
//             continue
//         }

//         if len(companies) == 0 {
//             log.Warnf("Company not found: %s", j.Company)
//             continue
//         }

//         companyID := companies[0].ID

//         // Create job model
//         jobModel := &job.Job{
//             CompanyID:       companyID,
//             Title:           j.Title,
//             Description:     j.Description,
//             ExperienceLevel: j.ExperienceLevel,
//             EmploymentType:  j.EmploymentType,
//             Location:        j.Location,
//             WorkMode:        j.WorkMode,
//             ApplicationURL:  j.ApplicationURL,
//             IsActive:        true,
//             Signature:       j.Signature,
//         }

//         // Insert job into database
//         err = jobRepo.Create(ctx, jobModel)
//         if err != nil {
//             if job.IsDuplicate(err) {
//                 log.Infof("Job already exists: %s at %s", j.Title, j.Company)
//                 continue
//             }
//             log.Warnf("Failed to insert job %s: %v", j.Title, err)
//             continue
//         }

//         log.Infof("Successfully added job: %s at %s (ID: %d)", jobModel.Title, j.Company, jobModel.ID)

//         // TODO: Add job technologies using job_technology repository
//     }

//     log.Info("Job population completed")
// }

func main() {}
