package main

import (
	"context"
	"os/signal"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/rodruizronald/ticos-in-tech/internal/database"
	"github.com/rodruizronald/ticos-in-tech/internal/technology"
	"github.com/rodruizronald/ticos-in-tech/internal/technology_alias"
)

func main() {
	// Configure logger
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Set up database connection
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Get database config
	dbConfig := database.DefaultConfig()

	log.Infof("Connecting to database %s at %s:%d", dbConfig.DBName, dbConfig.Host, dbConfig.Port)

	// Connect to the database
	dbpool, err := database.Connect(ctx, dbConfig)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbpool.Close()

	// Create technology repository
	techRepo := technology.NewRepository(dbpool)

	// Create technology alias repository
	aliasRepo := technology_alias.NewRepository(dbpool)

	// Create a map to store all technologies by name for lookup
	techMap := make(map[string]*technology.Technology)

	// Process and insert all technologies
	allTechnologies := [][]Technology{
		programmingLanguages,
		frontendTechnologies,
		backendTechnologies,
		databaseTechnologies,
		apiTechnologies,
		cloudTechnologies,
		devOpsTechnologies,
		observabilityTechnologies,
		testingTechnologies,
		osTechnologies,
		aiTechnologies,
		productivityTechnologies,
		dataScienceTechnologies,
		messagingTechnologies,
		otherTechnologies,
	}

	log.Info("Starting first pass: creating technologies without parent references")
	// First pass: insert all technologies without parent references
	for _, techGroup := range allTechnologies {
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

	log.Info("Starting second pass: updating technologies with parent references")
	// Second pass: update technologies with parent references
	for _, techGroup := range allTechnologies {
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

	log.Info("Technology import completed")
}

// addAliases adds aliases for a technology
func addAliases(ctx context.Context, log *logrus.Logger, aliasRepo *technology_alias.Repository, techID int, aliases []string) {
	for _, aliasName := range aliases {
		if aliasName == "" {
			continue
		}

		// Convert alias to lowercase
		lowerAlias := strings.ToLower(aliasName)

		// Create alias model
		newAlias := &technology_alias.TechnologyAlias{
			TechnologyID: techID,
			Alias:        lowerAlias,
		}

		// Insert into database
		err := aliasRepo.Create(ctx, newAlias)
		if err != nil {
			// Skip if it's a duplicate
			if technology_alias.IsDuplicate(err) {
				log.Infof("Alias already exists: %s", lowerAlias)
				continue
			}
			log.Warnf("Error creating alias %s for technology ID %d: %v", lowerAlias, techID, err)
			continue
		}

		log.Infof("Created alias: %s (ID: %d) for technology ID %d", lowerAlias, newAlias.ID, techID)
	}
}

type Technology struct {
	Name     string
	Category string
	Alias    []string
	Parent   string
}

var programmingLanguages = []Technology{
	{Name: "JavaScript", Category: "programming_language", Parent: "", Alias: []string{"JS", "ECMAScript"}},
	{Name: "Python", Category: "programming_language", Parent: "", Alias: []string{"Py", "CPython", "Python3"}},
	{Name: "Java", Category: "programming_language", Parent: "", Alias: []string{"JVM", "JDK"}},
	{Name: "TypeScript", Category: "programming_language", Parent: "JavaScript", Alias: []string{"TS", "Typed JS"}},
	{Name: "C", Category: "programming_language", Parent: "", Alias: []string{}},
	{Name: "C#", Category: "programming_language", Parent: "", Alias: []string{"CSharp", "C Sharp"}},
	{Name: "C++", Category: "programming_language", Parent: "", Alias: []string{"CPP", "ISO/IEC 14882"}},
	{Name: "Go", Category: "programming_language", Parent: "", Alias: []string{"Golang"}},
	{Name: "Ruby", Category: "programming_language", Parent: "", Alias: []string{"RoR", "Ruby Rails", "Ruby on Rails"}},
	{Name: "PHP", Category: "programming_language", Parent: "", Alias: []string{}},
	{Name: "Swift", Category: "programming_language", Parent: "", Alias: []string{"Swift Lang"}},
	{Name: "Kotlin", Category: "programming_language", Parent: "Java", Alias: []string{"KT"}},
	{Name: "Rust", Category: "programming_language", Parent: "", Alias: []string{"Rustlang", "RS"}},
	{Name: "Scala", Category: "programming_language", Parent: "Java", Alias: []string{"Scalable Language"}},
	{Name: "R", Category: "programming_language", Parent: "", Alias: []string{"R Language", "GNU R"}},
	{Name: "MATLAB", Category: "programming_language", Parent: "", Alias: []string{"Matrix Laboratory"}},
	{Name: "Perl", Category: "programming_language", Parent: "", Alias: []string{"Perl5", "Perl6"}},
	{Name: "Groovy", Category: "programming_language", Parent: "Java", Alias: []string{"Apache Groovy"}},
	{Name: "Bash", Category: "programming_language", Parent: "", Alias: []string{"Shell", "Shell Script"}},
	{Name: "PowerShell", Category: "programming_language", Parent: "", Alias: []string{"POSH", "PS"}},
	{Name: "Dart", Category: "programming_language", Parent: "", Alias: []string{"Flutter"}},
	{Name: "Clojure", Category: "programming_language", Parent: "Java", Alias: []string{"CLJ"}},
	{Name: "Elixir", Category: "programming_language", Parent: "Erlang", Alias: []string{"EX"}},
	{Name: "Erlang", Category: "programming_language", Parent: "", Alias: []string{"BEAM VM"}},
	{Name: "F#", Category: "programming_language", Parent: "", Alias: []string{"FSharp", "F Sharp"}},
	{Name: "Haskell", Category: "programming_language", Parent: "", Alias: []string{"GHC"}},
	{Name: "Julia", Category: "programming_language", Parent: "", Alias: []string{"JL"}},
	{Name: "Lua", Category: "programming_language", Parent: "", Alias: []string{"LuaJIT"}},
	{Name: "Objective-C", Category: "programming_language", Parent: "C", Alias: []string{"ObjC", "Obj-C"}},
	{Name: "OCaml", Category: "programming_language", Parent: "", Alias: []string{"Objective Caml"}},
	{Name: "SQL", Category: "programming_language", Parent: "", Alias: []string{"Structured Query Language", "T-SQL", "PL/SQL", "MySQL"}},
	{Name: "VBA", Category: "programming_language", Parent: "", Alias: []string{"Visual Basic for Applications", "Visual Basic", "Excel VBA"}},
}

var frontendTechnologies = []Technology{
	{Name: "React", Category: "frontend", Parent: "JavaScript", Alias: []string{"React.js", "ReactJS", "React DOM", "React Native"}},
	{Name: "Angular", Category: "frontend", Parent: "JavaScript", Alias: []string{"AngularJS", "Angular.js", "Angular 2+"}},
	{Name: "Vue", Category: "frontend", Parent: "JavaScript", Alias: []string{"Vue.js", "VueJS", "Vue 3"}},
	{Name: "Svelte", Category: "frontend", Parent: "JavaScript", Alias: []string{"SvelteKit", "Svelte 3"}},
	{Name: "jQuery", Category: "frontend", Parent: "JavaScript", Alias: []string{"JQ"}},
	{Name: "HTML", Category: "frontend", Parent: "", Alias: []string{"HTML5", "XHTML"}},
	{Name: "CSS", Category: "frontend", Parent: "", Alias: []string{"CSS3", "Stylesheets"}},
	{Name: "SASS", Category: "frontend", Parent: "CSS", Alias: []string{}},
	{Name: "SCSS", Category: "frontend", Parent: "CSS", Alias: []string{"Sassy CSS", "Sass CSS"}},
	{Name: "LESS", Category: "frontend", Parent: "CSS", Alias: []string{"Leaner CSS", "Less.js"}},
	{Name: "Bootstrap", Category: "frontend", Parent: "CSS", Alias: []string{"BS", "Bootstrap 5"}},
	{Name: "Tailwind CSS", Category: "frontend", Parent: "CSS", Alias: []string{"Tailwind", "TW CSS"}},
	{Name: "Material UI", Category: "frontend", Parent: "React", Alias: []string{"MUI", "Material-UI", "Material Design React"}},
	{Name: "Ant Design", Category: "frontend", Parent: "React", Alias: []string{"AntD", "Ant"}},
	{Name: "Redux", Category: "frontend", Parent: "JavaScript", Alias: []string{"React Redux", "Redux Toolkit", "RTK"}},
	{Name: "MobX", Category: "frontend", Parent: "JavaScript", Alias: []string{"MobX State Tree", "MST"}},
	{Name: "Next.js", Category: "frontend", Parent: "React", Alias: []string{"Next", "React Next", "Vercel Next"}},
	{Name: "Gatsby", Category: "frontend", Parent: "React", Alias: []string{"GatsbyJS", "Gatsby.js"}},
	{Name: "Webpack", Category: "frontend", Parent: "JavaScript", Alias: []string{"Webpack 5", "Web Pack"}},
	{Name: "Vite", Category: "frontend", Parent: "JavaScript", Alias: []string{"Vite.js", "ViteJS"}},
	{Name: "Ember", Category: "frontend", Parent: "JavaScript", Alias: []string{"Ember.js", "EmberJS", "Ember Data"}},
	{Name: "Backbone.js", Category: "frontend", Parent: "JavaScript", Alias: []string{"BackboneJS", "Backbone"}},
	{Name: "Alpine.js", Category: "frontend", Parent: "JavaScript", Alias: []string{"AlpineJS", "Alpine"}},
	{Name: "Storybook", Category: "frontend", Parent: "JavaScript", Alias: []string{"Storybook.js", "SB"}},
	{Name: "Stimulus", Category: "frontend", Parent: "JavaScript", Alias: []string{"Stimulus.js", "StimulusJS"}},
	{Name: "Preact", Category: "frontend", Parent: "JavaScript", Alias: []string{"Preact.js", "PreactJS", "Preact X"}},
	{Name: "Lit", Category: "frontend", Parent: "JavaScript", Alias: []string{"Lit Element", "Lit HTML", "Lit.js"}},
	{Name: "Web Components", Category: "frontend", Parent: "", Alias: []string{"HTML Templates"}},
}

var backendTechnologies = []Technology{
	{Name: "Django", Category: "backend", Parent: "Python", Alias: []string{"Django Framework", "Django REST Framework"}},
	{Name: "Flask", Category: "backend", Parent: "Python", Alias: []string{"Flask Framework", "Flask API"}},
	{Name: "FastAPI", Category: "backend", Parent: "Python", Alias: []string{"Fast API"}},
	{Name: "Spring", Category: "backend", Parent: "Java", Alias: []string{"Spring Framework", "Spring MVC", "Spring Cloud"}},
	{Name: "Spring Boot", Category: "backend", Parent: "Spring", Alias: []string{"SpringBoot", "Boot"}},
	{Name: "Express", Category: "backend", Parent: "Node.js", Alias: []string{"Express.js", "ExpressJS", "Express Framework"}},
	{Name: ".NET", Category: "backend", Parent: "", Alias: []string{".NET Framework", ".NET Core"}},
	{Name: "ASP.NET", Category: "backend", Parent: ".NET", Alias: []string{"ASP", "ASP.NET MVC"}},
	{Name: "ASP.NET Core", Category: "backend", Parent: "ASP.NET", Alias: []string{"ASP Core", "ASP.NET 6+"}},
	{Name: "Laravel", Category: "backend", Parent: "PHP", Alias: []string{"Laravel Framework", "Laravel API", "Laravel 10"}},
	{Name: "Nest.js", Category: "backend", Parent: "Node.js", Alias: []string{"NestJS", "Nest Framework", "Nest"}},
	{Name: "Play", Category: "backend", Parent: "Java", Alias: []string{"Play Framework", "Play Java", "Play Scala"}},
	{Name: "Phoenix", Category: "backend", Parent: "Elixir", Alias: []string{"Phoenix Framework"}},
	{Name: "Rocket", Category: "backend", Parent: "Rust", Alias: []string{"Rocket.rs", "Rocket Framework"}},
	{Name: "Gin", Category: "backend", Parent: "Go", Alias: []string{"Gin Gonic", "Gin Framework", "Gin Web"}},
	{Name: "Echo", Category: "backend", Parent: "Go", Alias: []string{"Echo Framework", "Echo Go", "Echo Web"}},
	{Name: "Symfony", Category: "backend", Parent: "PHP", Alias: []string{"Symfony Framework", "Symfony PHP", "Symfony 6"}},
	{Name: "Strapi", Category: "backend", Parent: "Node.js", Alias: []string{"Strapi CMS", "Strapi Headless", "Strapi API"}},
	{Name: "Node.js", Category: "backend", Parent: "JavaScript", Alias: []string{"NodeJS", "Node", "Node Runtime", "Node Server"}},
	{Name: "Deno", Category: "backend", Parent: "JavaScript", Alias: []string{"Deno Runtime", "Deno Land"}},
	{Name: "Bun", Category: "backend", Parent: "JavaScript", Alias: []string{"Bun Runtime", "Bun.js", "BunJS"}},
	{Name: "Micronaut", Category: "backend", Parent: "Java", Alias: []string{"Micronaut Framework", "Micronaut GraalVM"}},
	{Name: "Quarkus", Category: "backend", Parent: "Java", Alias: []string{"Quarkus.io", "Quarkus Framework"}},
	{Name: "Ktor", Category: "backend", Parent: "Kotlin", Alias: []string{"Ktor Framework", "Kotlin Server", "Ktor Server"}},
	{Name: "Actix", Category: "backend", Parent: "Rust", Alias: []string{"Actix Web", "Actix.rs", "Actix Framework"}},
	{Name: "Axum", Category: "backend", Parent: "Rust", Alias: []string{"Axum Framework", "Tokio Axum", "Axum Web"}},
	{Name: "Fiber", Category: "backend", Parent: "Go", Alias: []string{"Go Fiber", "Fiber Framework"}},
	{Name: "Buffalo", Category: "backend", Parent: "Go", Alias: []string{"Buffalo Framework", "Go Buffalo"}},
}

var databaseTechnologies = []Technology{
	{Name: "MySQL", Category: "databases", Parent: "", Alias: []string{"MySQL DB", "MySQL Server", "Oracle MySQL"}},
	{Name: "PostgreSQL", Category: "databases", Parent: "", Alias: []string{"Postgres", "PG", "PGSQL", "Postgre"}},
	{Name: "SQLite", Category: "databases", Parent: "", Alias: []string{"SQLite3", "Lite SQL", "SQL Lite"}},
	{Name: "Oracle", Category: "databases", Parent: "", Alias: []string{"Oracle DB", "Oracle Database", "OracleDB"}},
	{Name: "Microsoft SQL Server", Category: "databases", Parent: "", Alias: []string{"MSSQL", "SQL Server", "MS SQL"}},
	{Name: "MongoDB", Category: "databases", Parent: "", Alias: []string{"Mongo", "MongoDB Atlas", "Document DB", "NoSQL DB"}},
	{Name: "Redis", Category: "databases", Parent: "", Alias: []string{"Redis Cache", "Redis Stack"}},
	{Name: "Cassandra", Category: "databases", Parent: "", Alias: []string{"Apache Cassandra", "Cassandra DB"}},
	{Name: "DynamoDB", Category: "databases", Parent: "Amazon Web Services", Alias: []string{"AWS DynamoDB", "Dynamo", "Amazon DynamoDB"}},
	{Name: "Firestore", Category: "databases", Parent: "Firebase", Alias: []string{"Cloud Firestore", "Firebase DB", "Google Firestore"}},
	{Name: "Elasticsearch", Category: "databases", Parent: "", Alias: []string{"ELK Stack", "Elastic Search"}},
	{Name: "Neo4j", Category: "databases", Parent: "", Alias: []string{"Neo4j Graph", "Graph Database", "Neo"}},
	{Name: "CouchDB", Category: "databases", Parent: "", Alias: []string{"Apache CouchDB", "Couch", "Couch Database"}},
	{Name: "MariaDB", Category: "databases", Parent: "MySQL", Alias: []string{"Maria", "MySQL Fork", "MariaDB Server"}},
	{Name: "Snowflake", Category: "databases", Parent: "", Alias: []string{"Snowflake DB", "Snowflake Cloud"}},
	{Name: "BigQuery", Category: "databases", Parent: "Google Cloud Platform", Alias: []string{"Google BigQuery", "GCP BigQuery", "BQ"}},
	{Name: "Redshift", Category: "databases", Parent: "Amazon Web Services", Alias: []string{"Amazon Redshift", "AWS Redshift"}},
	{Name: "Supabase", Category: "databases", Parent: "PostgreSQL", Alias: []string{"Supabase DB", "Firebase Alternative", "Open Source backend"}},
	{Name: "InfluxDB", Category: "databases", Parent: "", Alias: []string{"Influx"}},
	{Name: "ArangoDB", Category: "databases", Parent: "", Alias: []string{"Arango", "Multi-model DB", "Graph DB"}},
	{Name: "RethinkDB", Category: "databases", Parent: "", Alias: []string{"Rethink", "Realtime DB", "Changefeeds DB"}},
	{Name: "H2", Category: "databases", Parent: "", Alias: []string{"H2 Database", "H2 Engine", "Embedded SQL"}},
	{Name: "Timescale", Category: "databases", Parent: "PostgreSQL", Alias: []string{"TimescaleDB", "Postgres Timeseries"}},

	// ORM and database access libraries
	{Name: "SQLAlchemy", Category: "databases", Parent: "Python", Alias: []string{"Python ORM", "SQL Toolkit", "Alembic"}},
	{Name: "Prisma", Category: "databases", Parent: "JavaScript", Alias: []string{"Prisma ORM", "Prisma Client", "Prisma Schema"}},
	{Name: "TypeORM", Category: "databases", Parent: "TypeScript", Alias: []string{"TS ORM", "TypeScript ORM", "Type ORM"}},
	{Name: "Mongoose", Category: "databases", Parent: "MongoDB", Alias: []string{"Mongoose ODM", "Mongo ORM", "MongoDB ODM"}},
	{Name: "Sequelize", Category: "databases", Parent: "JavaScript", Alias: []string{"Sequelize ORM", "Sequelize.js", "SQL ORM"}},
	{Name: "Knex", Category: "databases", Parent: "JavaScript", Alias: []string{"Knex.js", "SQL Query Builder", "JS Query Builder"}},
	{Name: "Drizzle", Category: "databases", Parent: "JavaScript", Alias: []string{"DrizzleORM", "Drizzle ORM"}},
}

var apiTechnologies = []Technology{
	{Name: "REST", Category: "api", Parent: "", Alias: []string{"RESTful API", "RESTful", "REST API"}},
	{Name: "GraphQL", Category: "api", Parent: "", Alias: []string{"GQL", "Graph Query Language", "Facebook API"}},
	{Name: "SOAP", Category: "api", Parent: "", Alias: []string{"XML Web Services", "SOAP API"}},
	{Name: "WebSockets", Category: "api", Parent: "", Alias: []string{"Web Sockets", "WebSocket Protocol", "Socket API"}},
	{Name: "Swagger", Category: "api", Parent: "OpenAPI", Alias: []string{"Swagger UI", "Swagger Docs", "API Documentation"}},
	{Name: "OpenAPI", Category: "api", Parent: "", Alias: []string{"OAS", "OpenAPI 3.0"}},
	{Name: "Apollo", Category: "api", Parent: "GraphQL", Alias: []string{"Apollo GraphQL"}},
	{Name: "Postman", Category: "api", Parent: "", Alias: []string{}},
	{Name: "OAuth", Category: "api", Parent: "", Alias: []string{"OAuth 2.0", "Open Authorization", "OAuth2", "API Authorization"}},
	{Name: "JWT", Category: "api", Parent: "", Alias: []string{"JSON Web Token", "JWS"}},
	{Name: "API Gateway", Category: "api", Parent: "", Alias: []string{"API Router", "Gateway"}},
	{Name: "Tyk", Category: "api", Parent: "API Gateway", Alias: []string{"Tyk Gateway"}},
	{Name: "Kong", Category: "api", Parent: "API Gateway", Alias: []string{"Kong Gateway", "Kong API", "Kong API Management"}},
	{Name: "Apigee", Category: "api", Parent: "API Gateway", Alias: []string{"Google Apigee", "Apigee Edge", "API Management Platform"}},
	{Name: "gRPC", Category: "api", Parent: "", Alias: []string{"Google RPC", "gRPC API"}},
	{Name: "Protocol Buffers", Category: "api", Parent: "gRPC", Alias: []string{"Protobuf", "Protobufs", "Proto"}},
	{Name: "tRPC", Category: "api", Parent: "TypeScript", Alias: []string{"TypeScript RPC", "TS RPC"}},
	{Name: "HATEOAS", Category: "api", Parent: "REST", Alias: []string{"Hypermedia API", "REST Hypermedia"}},
	{Name: "RPC", Category: "api", Parent: "", Alias: []string{"Remote Procedure Call", "Remote Call", "Procedure Call", "JSON-RPC"}},
	{Name: "GraphQL Federation", Category: "api", Parent: "GraphQL", Alias: []string{"Federated GraphQL", "Distributed GraphQL"}},
}

var cloudTechnologies = []Technology{
	{Name: "Amazon Web Services", Category: "cloud", Parent: "", Alias: []string{"AWS", "Amazon Cloud", "Amazon Services", "AWS Cloud"}},
	{Name: "Microsoft Azure", Category: "cloud", Parent: "", Alias: []string{"Azure", "MS Azure", "Azure Cloud", "MSFT Cloud"}},
	{Name: "Google Cloud Platform", Category: "cloud", Parent: "", Alias: []string{"GCP", "Google Cloud", "GC", "Cloud Platform"}},
	{Name: "IBM Cloud", Category: "cloud", Parent: "", Alias: []string{"IBM Cloud Services", "Bluemix"}},
	{Name: "Oracle Cloud", Category: "cloud", Parent: "", Alias: []string{"OCI", "Oracle Cloud Infrastructure", "Oracle Public Cloud"}},
	{Name: "DigitalOcean", Category: "cloud", Parent: "", Alias: []string{"DO", "Digital Ocean", "DO Cloud", "DO Droplets"}},
	{Name: "Heroku", Category: "cloud", Parent: "", Alias: []string{"Heroku Platform", "Heroku PaaS", "Heroku Cloud"}},
	{Name: "Netlify", Category: "cloud", Parent: "", Alias: []string{"Netlify Platform", "Netlify Hosting", "Jamstack Platform"}},
	{Name: "Vercel", Category: "cloud", Parent: "", Alias: []string{"Vercel Platform", "Zeit", "Next.js Platform", "Vercel Hosting"}},
	{Name: "Firebase", Category: "cloud", Parent: "Google Cloud Platform", Alias: []string{"FB", "Firebase Platform", "Google Firebase", "Firebase Cloud"}},
	{Name: "Cloudflare", Category: "cloud", Parent: "", Alias: []string{"CF", "Cloudflare Services", "Cloudflare Network", "Cloudflare Workers"}},
	{Name: "Fly.io", Category: "cloud", Parent: "", Alias: []string{"Fly", "Fly Platform", "Fly Edge"}},
	{Name: "Render", Category: "cloud", Parent: "", Alias: []string{"Render Cloud", "Render Platform", "Render Hosting"}},
	{Name: "Railway", Category: "cloud", Parent: "", Alias: []string{"Railway App", "Railway Platform", "Railway Hosting"}},
	{Name: "Linode", Category: "cloud", Parent: "", Alias: []string{"Akamai Cloud Computing", "Linode Cloud", "Linode Hosting"}},
	{Name: "Vultr", Category: "cloud", Parent: "", Alias: []string{"Vultr Cloud", "Vultr Compute", "Vultr Platform"}},
	{Name: "Scaleway", Category: "cloud", Parent: "", Alias: []string{"Scaleway Cloud", "Scaleway Platform", "Scaleway Elements"}},
	{Name: "OVHcloud", Category: "cloud", Parent: "", Alias: []string{"OVH", "OVH Cloud", "Online SAS", "OVH Hosting"}},
	{Name: "Backblaze", Category: "cloud", Parent: "", Alias: []string{"B2 Cloud Storage", "Backblaze B2", "B2", "Backblaze Storage"}},
	{Name: "S3", Category: "cloud", Parent: "Amazon Web Services", Alias: []string{"Simple Storage Service", "AWS S3", "Amazon S3", "Object Storage"}},
	{Name: "EC2", Category: "cloud", Parent: "Amazon Web Services", Alias: []string{"Elastic Compute Cloud", "AWS EC2", "Amazon EC2", "Virtual Servers"}},
	{Name: "Lambda", Category: "cloud", Parent: "Amazon Web Services", Alias: []string{"AWS Lambda", "Amazon Lambda", "Serverless Functions", "FaaS"}},
	{Name: "CloudFront", Category: "cloud", Parent: "Amazon Web Services", Alias: []string{"AWS CloudFront", "Amazon CloudFront", "CDN", "Content Delivery"}},
	{Name: "Route 53", Category: "cloud", Parent: "Amazon Web Services", Alias: []string{"AWS Route 53", "Amazon Route 53", "DNS Service", "AWS DNS"}},
	{Name: "IAM", Category: "cloud", Parent: "Amazon Web Services", Alias: []string{"Identity and Access Management", "AWS IAM", "Amazon IAM", "Access Control"}},
	{Name: "ECS", Category: "cloud", Parent: "Amazon Web Services", Alias: []string{"Elastic Container Service", "AWS ECS", "Amazon ECS"}},
	{Name: "EKS", Category: "cloud", Parent: "Amazon Web Services", Alias: []string{"Elastic Kubernetes Service", "AWS EKS", "Amazon EKS", "Managed Kubernetes"}},
	{Name: "Fargate", Category: "cloud", Parent: "Amazon Web Services", Alias: []string{"AWS Fargate", "Amazon Fargate", "Serverless Containers", "Serverless ECS"}},
}

var devOpsTechnologies = []Technology{
	{Name: "Docker", Category: "devops", Parent: "", Alias: []string{"Docker Engine", "Docker Containers", "Moby", "Docker Platform"}},
	{Name: "Kubernetes", Category: "devops", Parent: "", Alias: []string{"K8s", "Kube", "K-Eights"}},
	{Name: "Jenkins", Category: "devops", Parent: "", Alias: []string{"Jenkins CI", "Jenkins Pipeline", "Hudson", "Jenkins X"}},
	{Name: "GitHub Actions", Category: "devops", Parent: "GitHub", Alias: []string{"GH Actions", "Actions", "GitHub CI/CD", "GHA"}},
	{Name: "GitLab CI/CD", Category: "devops", Parent: "GitLab", Alias: []string{"GitLab Pipeline", "GL CI", "GitLab CI", "GitLab Pipelines"}},
	{Name: "CircleCI", Category: "devops", Parent: "", Alias: []string{"Circle", "Circle Pipeline", "CircleCI Pipeline", "Circle Build"}},
	{Name: "Travis CI", Category: "devops", Parent: "", Alias: []string{"Travis", "Travis Build", "Travis Pipeline", "Travis-CI"}},
	{Name: "Terraform", Category: "devops", Parent: "", Alias: []string{"HashiCorp Terraform", "IaC", "Terraform Cloud"}},
	{Name: "Ansible", Category: "devops", Parent: "", Alias: []string{"Red Hat Ansible", "Ansible Playbooks", "Ansible Tower"}},
	{Name: "Puppet", Category: "devops", Parent: "", Alias: []string{"Puppet Labs", "Puppet Enterprise", "Puppet Forge"}},
	{Name: "Chef", Category: "devops", Parent: "", Alias: []string{"Chef Infra", "Progress Chef", "Chef Cookbooks"}},
	{Name: "Vagrant", Category: "devops", Parent: "", Alias: []string{"HashiCorp Vagrant", "VM Management", "Development Environments", "Vagrant Boxes"}},
	{Name: "ArgoCD", Category: "devops", Parent: "Kubernetes", Alias: []string{"Argo CD", "Argo", "GitOps Controller", "Kubernetes CD"}},
	{Name: "Helm", Category: "devops", Parent: "Kubernetes", Alias: []string{"Helm Charts", "Kubernetes Package Manager", "Helm 3", "K8s Package Manager"}},
	{Name: "Harbor", Category: "devops", Parent: "", Alias: []string{"Harbor Registry", "Container Registry", "CNCF Harbor", "OCI Registry"}},
	{Name: "Rancher", Category: "devops", Parent: "Kubernetes", Alias: []string{"Rancher Labs", "Rancher K8s", "Container Management Platform", "SUSE Rancher"}},
	{Name: "Podman", Category: "devops", Parent: "", Alias: []string{"Pod Manager", "Red Hat Podman", "Docker Alternative", "Daemonless Containers"}},
	{Name: "containerd", Category: "devops", Parent: "", Alias: []string{"Container Runtime", "CNCF containerd", "OCI Runtime", "Container Daemon"}},
	{Name: "Buildkite", Category: "devops", Parent: "", Alias: []string{"Buildkite Pipeline", "BK", "Buildkite CI", "Buildkite Agent"}},
	{Name: "TeamCity", Category: "devops", Parent: "", Alias: []string{"JetBrains TeamCity", "TC", "TeamCity CI", "TeamCity Server"}},
	{Name: "Octopus Deploy", Category: "devops", Parent: "", Alias: []string{"Octopus", "Octopus CD", "Deployment Automation", "Release Management"}},
	{Name: "Spinnaker", Category: "devops", Parent: "", Alias: []string{"CD Platform", "Multi-Cloud CD", "Netflix Spinnaker", "Continuous Delivery"}},
	{Name: "FluxCD", Category: "devops", Parent: "Kubernetes", Alias: []string{"Flux", "GitOps Operator", "Flux v2", "CNCF Flux", "Kubernetes GitOps"}},
}

var observabilityTechnologies = []Technology{
	{Name: "New Relic", Category: "observability", Parent: "", Alias: []string{"NR", "New Relic One", "NRDB", "New Relic APM"}},
	{Name: "Datadog", Category: "observability", Parent: "", Alias: []string{"DD", "Datadog Platform", "Datadog APM", "Datadog Monitoring"}},
	{Name: "Splunk", Category: "observability", Parent: "", Alias: []string{"Splunk Enterprise", "Splunk Cloud", "Splunk Platform", "SPL"}},
	{Name: "ELK Stack", Category: "observability", Parent: "", Alias: []string{"Elastic Stack", "ELK", "Elasticsearch-Logstash-Kibana", "Elastic"}},
	{Name: "Logstash", Category: "observability", Parent: "ELK Stack", Alias: []string{"LS", "Log Pipeline", "Log Processor", "ETL Pipeline"}},
	{Name: "Kibana", Category: "observability", Parent: "ELK Stack", Alias: []string{"KB", "Elastic Dashboards", "Kibana Dashboard"}},
	{Name: "Grafana", Category: "observability", Parent: "", Alias: []string{"Grafana Labs", "Grafana Dashboards", "Grafana Cloud"}},
	{Name: "Prometheus", Category: "observability", Parent: "", Alias: []string{"Prom", "Prometheus Monitoring", "Time Series Database"}},
	{Name: "Sentry", Category: "observability", Parent: "", Alias: []string{"Sentry.io", "Error Tracking", "Exception Monitoring", "Crash Reporting"}},
	{Name: "PagerDuty", Category: "observability", Parent: "", Alias: []string{"PD", "On-Call Platform", "Alert Management"}},
	{Name: "AppDynamics", Category: "observability", Parent: "", Alias: []string{"AppD", "Cisco AppDynamics", "App Dynamics", "APM"}},
	{Name: "Dynatrace", Category: "observability", Parent: "", Alias: []string{"DT", "Dynatrace Platform", "DT Monitoring", "Full-Stack Monitoring"}},
	{Name: "Honeycomb", Category: "observability", Parent: "", Alias: []string{"Honeycomb.io", "Observability Platform", "BubbleUp", "High-Cardinality Analysis"}},
	{Name: "Lightstep", Category: "observability", Parent: "", Alias: []string{"Lightstep Observability", "ServiceNow Lightstep", "Observability Cloud"}},
	{Name: "OpenTelemetry", Category: "observability", Parent: "", Alias: []string{"OTel", "OTEL", "Open Telemetry", "CNCF OpenTelemetry", "OT"}},
	{Name: "Jaeger", Category: "observability", Parent: "OpenTelemetry", Alias: []string{"Jaeger Tracing", "CNCF Jaeger", "Uber Jaeger"}},
	{Name: "Zipkin", Category: "observability", Parent: "", Alias: []string{"Zipkin Tracing", "Twitter Zipkin", "Trace Collection"}},
	{Name: "Loki", Category: "observability", Parent: "Grafana", Alias: []string{"Grafana Loki", "Log Aggregation", "Promtail", "Log-like Prometheus"}},
	{Name: "Fluentd", Category: "observability", Parent: "", Alias: []string{"Fluent", "Log Collector", "CNCF Fluentd", "Log Forwarder"}},
	{Name: "Cloudwatch", Category: "observability", Parent: "Amazon Web Services", Alias: []string{"AWS CloudWatch", "Amazon CloudWatch", "AWS Logs", "AWS Metrics"}},
	{Name: "Nagios", Category: "observability", Parent: "", Alias: []string{"Nagios Core", "Nagios XI", "Infrastructure Monitoring", "Legacy Monitoring"}},
	{Name: "Zabbix", Category: "observability", Parent: "", Alias: []string{"Zabbix Monitoring", "Zabbix Platform", "Enterprise Monitoring", "Server Monitoring"}},
	{Name: "Instana", Category: "observability", Parent: "", Alias: []string{"IBM Instana", "Instana APM", "Instana Observability", "Auto-Discovery Platform"}},
	{Name: "Graphite", Category: "observability", Parent: "", Alias: []string{"Graphite Monitoring", "Carbon", "Whisper"}},
}

var testingTechnologies = []Technology{
	{Name: "Jest", Category: "testing", Parent: "JavaScript", Alias: []string{"Jest.js", "Facebook Jest", "React Testing"}},
	{Name: "Mocha", Category: "testing", Parent: "JavaScript", Alias: []string{"Mocha.js", "JS Test Framework", "Node Testing"}},
	{Name: "Cypress", Category: "testing", Parent: "JavaScript", Alias: []string{"Cypress.io", "E2E Testing", "Cypress Testing", "End-to-End Tests"}},
	{Name: "Selenium", Category: "testing", Parent: "", Alias: []string{"Selenium WebDriver", "Browser Automation", "Selenium Grid", "UI Testing"}},
	{Name: "Playwright", Category: "testing", Parent: "JavaScript", Alias: []string{"Microsoft Playwright", "Chromium Testing", "Multi-Browser Testing", "Headless Testing"}},
	{Name: "Puppeteer", Category: "testing", Parent: "JavaScript", Alias: []string{"Google Puppeteer", "Headless Chrome", "Chrome Automation", "Browser API"}},
	{Name: "JUnit", Category: "testing", Parent: "Java", Alias: []string{"JUnit 5", "Java Unit Testing", "JUnit Jupiter", "xUnit for Java"}},
	{Name: "TestNG", Category: "testing", Parent: "Java", Alias: []string{"Java Testing", "TestNG Framework", "Java Test Suite", "NextGen Testing"}},
	{Name: "pytest", Category: "testing", Parent: "Python", Alias: []string{"py.test", "Python Testing", "pytest Framework", "Python Test Framework"}},
	{Name: "unittest", Category: "testing", Parent: "Python", Alias: []string{"Python unittest", "PyUnit", "Python Test Module", "Standard Test Framework"}},
	{Name: "PHPUnit", Category: "testing", Parent: "PHP", Alias: []string{"PHP Testing", "PHP Test Framework", "xUnit for PHP", "PHP Unit Tests"}},
	{Name: "RSpec", Category: "testing", Parent: "Ruby", Alias: []string{"Ruby Spec", "BDD for Ruby", "Ruby Testing"}},
	{Name: "Jasmine", Category: "testing", Parent: "JavaScript", Alias: []string{"Jasmine.js", "Behavior-Driven JS"}},
	{Name: "Karma", Category: "testing", Parent: "JavaScript", Alias: []string{"Karma Runner", "JS Test Runner", "Test Environment"}},
	{Name: "Protractor", Category: "testing", Parent: "JavaScript", Alias: []string{"Protractor E2E", "E2E Framework", "Angular E2E"}},
	{Name: "SoapUI", Category: "testing", Parent: "", Alias: []string{"SmartBear SoapUI", "SOAP Testing", "REST Testing Tool"}},
	{Name: "RestAssured", Category: "testing", Parent: "Java", Alias: []string{"REST Assured", "Java REST Testing", "Java HTTP Testing"}},
	{Name: "Mockito", Category: "testing", Parent: "Java", Alias: []string{"Mock Framework", "Java Test Doubles", "Mockito Framework"}},
	{Name: "EasyMock", Category: "testing", Parent: "Java", Alias: []string{"Java EasyMock", "Mock Objects", "Test Doubles"}},
	{Name: "WireMock", Category: "testing", Parent: "Java", Alias: []string{"HTTP Mocking", "API Mocking", "Mock Server"}},
	{Name: "WebMock", Category: "testing", Parent: "Ruby", Alias: []string{"Ruby HTTP Mocking", "Ruby Network Mocks", "Request Mocking"}},
	{Name: "VCR", Category: "testing", Parent: "Ruby", Alias: []string{"VCR.rb", "HTTP Recording", "Request Playback", "Cassette Testing"}},
	{Name: "Artillery", Category: "testing", Parent: "JavaScript", Alias: []string{"Artillery.io", "Scalability Testing"}},
	{Name: "K6", Category: "testing", Parent: "JavaScript", Alias: []string{"k6.io", "Grafana k6", "JS Load Testing"}},
	{Name: "Gatling", Category: "testing", Parent: "Scala", Alias: []string{"Gatling.io", "Scala Testing", "Performance Framework"}},
	{Name: "Locust", Category: "testing", Parent: "Python", Alias: []string{"Locust.io", "Python Load Testing", "Distributed Load Testing"}},
	{Name: "Cucumber", Category: "testing", Parent: "", Alias: []string{"Gherkin", "BDD Framework", "Cucumber.io", "Specification by Example"}},
	{Name: "Robot Framework", Category: "testing", Parent: "", Alias: []string{"RF", "Robot", "Acceptance Testing", "ATDD Framework", "Test Automation"}},
}

var osTechnologies = []Technology{
	{Name: "Linux", Category: "os", Parent: "", Alias: []string{"GNU/Linux", "Linux Kernel", "Tux", "Linux OS"}},
	{Name: "Ubuntu", Category: "os", Parent: "Linux", Alias: []string{"Ubuntu Linux", "Canonical Ubuntu", "Debian-based", "Ubuntu LTS"}},
	{Name: "Debian", Category: "os", Parent: "Linux", Alias: []string{"Debian GNU/Linux", "Debian Linux", "The Universal OS", "Stable Distribution"}},
	{Name: "CentOS", Category: "os", Parent: "Linux", Alias: []string{"CentOS Linux", "Community Enterprise OS", "RHEL Clone", "CentOS Stream"}},
	{Name: "Red Hat", Category: "os", Parent: "Linux", Alias: []string{"RHEL", "Red Hat Enterprise Linux", "Red Hat Linux", "Enterprise Linux"}},
	{Name: "Windows", Category: "os", Parent: "", Alias: []string{"Microsoft Windows", "Windows OS", "Win", "Windows NT"}},
	{Name: "macOS", Category: "os", Parent: "", Alias: []string{"Mac OS X", "Mac", "Apple macOS", "OS X", "Mac Operating System"}},
	{Name: "iOS", Category: "os", Parent: "", Alias: []string{"iPhone OS", "Apple iOS", "iPadOS", "Mobile OS"}},
	{Name: "Android", Category: "os", Parent: "", Alias: []string{"Google Android", "Android OS", "AOSP", "Android Platform"}},
	{Name: "Alpine", Category: "os", Parent: "Linux", Alias: []string{"Alpine Linux", "musl libc", "BusyBox", "Container OS"}},
	{Name: "Fedora", Category: "os", Parent: "Linux", Alias: []string{"Fedora Linux", "Fedora Core", "Red Hat Fedora", "RHEL Upstream"}},
	{Name: "Arch Linux", Category: "os", Parent: "Linux", Alias: []string{"Arch", "Rolling Release", "Pacman", "Minimalist Linux"}},
	{Name: "FreeBSD", Category: "os", Parent: "", Alias: []string{"Free BSD", "Berkeley Software Distribution", "BSD OS", "FreeBSD Unix"}},
	{Name: "OpenBSD", Category: "os", Parent: "", Alias: []string{"Open BSD", "Secure by Default", "OpenBSD Unix"}},
	{Name: "NetBSD", Category: "os", Parent: "", Alias: []string{"Net BSD", "Portable BSD", "NetBSD Unix"}},
	{Name: "Solaris", Category: "os", Parent: "", Alias: []string{"Oracle Solaris", "Sun Solaris", "SunOS"}},
	{Name: "Unix", Category: "os", Parent: "", Alias: []string{"UNIX", "Unix System", "Unix-like", "Original OS"}},
	{Name: "ChromeOS", Category: "os", Parent: "Linux", Alias: []string{"Chrome OS", "Google ChromeOS", "Chromium OS", "Chromebook OS"}},
	{Name: "Windows Server", Category: "os", Parent: "Windows", Alias: []string{"Windows NT Server", "Microsoft Server", "Server OS", "WinServer"}},
	{Name: "AIX", Category: "os", Parent: "Unix", Alias: []string{"IBM AIX", "Power Systems"}},
	{Name: "HP-UX", Category: "os", Parent: "Unix", Alias: []string{"HP Unix", "Hewlett Packard Unix", "HP-UX Unix"}},
}

var aiTechnologies = []Technology{
	{Name: "TensorFlow", Category: "ai", Parent: "", Alias: []string{"Google TensorFlow", "TensorFlow ML", "TF Framework"}},
	{Name: "PyTorch", Category: "ai", Parent: "", Alias: []string{"Torch", "Facebook PyTorch", "Meta PyTorch", "PT"}},
	{Name: "scikit-learn", Category: "ai", Parent: "Python", Alias: []string{"sklearn", "Scikit", "SK Learn", "Python ML Library"}},
	{Name: "Keras", Category: "ai", Parent: "TensorFlow", Alias: []string{"TF Keras", "Keras API", "Deep Learning API", "High-Level API"}},
	{Name: "NLTK", Category: "ai", Parent: "Python", Alias: []string{"Natural Language Toolkit", "NLTK Library", "Text Processing"}},
	{Name: "spaCy", Category: "ai", Parent: "Python", Alias: []string{"Industrial NLP", "spaCy NLP", "Advanced NLP"}},
	{Name: "Hugging Face", Category: "ai", Parent: "", Alias: []string{"HF", "Transformers", "HuggingFace Hub", "Transformers Library"}},
	{Name: "OpenAI API", Category: "ai", Parent: "", Alias: []string{"GPT API", "OpenAI Platform", "GPT-4 API", "Completion API"}},
	{Name: "LangChain", Category: "ai", Parent: "", Alias: []string{"LC", "LLM Framework", "Chain of Thought", "Agents Framework"}},
	{Name: "OpenCV", Category: "ai", Parent: "", Alias: []string{"Open Computer Vision", "CV Library", "Image Processing", "Computer Vision"}},
	{Name: "CUDA", Category: "ai", Parent: "", Alias: []string{"NVIDIA CUDA", "GPU Computing", "Parallel Computing", "CUDA Toolkit"}},
	{Name: "JAX", Category: "ai", Parent: "", Alias: []string{"Google JAX", "JAX Framework", "Autograd", "XLA"}},
	{Name: "MLflow", Category: "ai", Parent: "", Alias: []string{"ML Pipeline", "ML Lifecycle", "Experiment Tracking", "Model Registry"}},
	{Name: "Kubeflow", Category: "ai", Parent: "Kubernetes", Alias: []string{"ML on Kubernetes", "KF", "K8s ML Platform", "ML Pipelines"}},
	{Name: "Weights & Biases", Category: "ai", Parent: "", Alias: []string{"W&B", "WandB", "ML Experiment Tracking", "ML Ops Platform"}},
	{Name: "XGBoost", Category: "ai", Parent: "", Alias: []string{"eXtreme Gradient Boosting", "XGBM", "XGB"}},
	{Name: "LightGBM", Category: "ai", Parent: "", Alias: []string{"Light Gradient Boosting", "LGBM", "Microsoft LightGBM", "Light GBM"}},
	{Name: "CatBoost", Category: "ai", Parent: "", Alias: []string{"Categorical Boosting", "Yandex CatBoost", "CB"}},
	{Name: "FastAI", Category: "ai", Parent: "PyTorch", Alias: []string{"Fast.ai", "FastAI Library", "High-level API", "Practical Deep Learning"}},
	{Name: "H2O", Category: "ai", Parent: "", Alias: []string{"H2O.ai", "H2O Platform", "AutoML", "Driverless AI"}},
	{Name: "Labelbox", Category: "ai", Parent: "", Alias: []string{"Data Labeling Platform", "Training Data Platform", "Data Annotation", "ML Data Platform"}},
	{Name: "Roboflow", Category: "ai", Parent: "", Alias: []string{"Computer Vision Platform", "CV Data Platform", "Vision AI", "Image Annotation"}},
	{Name: "TensorRT", Category: "ai", Parent: "", Alias: []string{"TRT", "NVIDIA Inference", "Model Optimization", "Inference Accelerator"}},
	{Name: "CoreML", Category: "ai", Parent: "", Alias: []string{"Core ML", "Apple ML", "iOS ML", "On-device ML"}},
	{Name: "TFLite", Category: "ai", Parent: "TensorFlow", Alias: []string{"TensorFlow Lite", "TF Mobile", "TF Edge", "Mobile ML"}},
	{Name: "ONNX", Category: "ai", Parent: "", Alias: []string{"Open Neural Network Exchange", "Model Exchange", "ML Interoperability", "Neural Network Format"}},
	{Name: "PyCaret", Category: "ai", Parent: "Python", Alias: []string{"Low-code ML", "AutoML Library", "ML Workflow", "Python AutoML"}},
}

var productivityTechnologies = []Technology{
	{Name: "Microsoft Office", Category: "productivity", Alias: []string{"Office", "MS Office", "Office Suite", "Microsoft 365 Apps"}, Parent: ""},
	{Name: "Google Workspace", Category: "productivity", Alias: []string{"G Suite", "Google Apps", "Google Suite", "Google for Business"}, Parent: ""},
	{Name: "Notion", Category: "productivity", Alias: []string{"Notion App", "Notion.so"}, Parent: ""},
	{Name: "Asana", Category: "productivity", Alias: []string{"Asana App", "Asana Project Management"}, Parent: ""},
	{Name: "Trello", Category: "productivity", Alias: []string{"Trello Boards", "Trello App"}, Parent: ""},
	{Name: "Slack", Category: "productivity", Alias: []string{"Slack App", "Slack Messenger", "Slack Chat", "Slack Workspace"}, Parent: ""},
	{Name: "Microsoft Teams", Category: "productivity", Alias: []string{"Teams", "MS Teams", "Teams App"}, Parent: ""},
	{Name: "Jira", Category: "productivity", Alias: []string{"Jira Software", "Atlassian Jira", "Jira Project Management", "Jira Agile"}, Parent: ""},
	{Name: "Confluence", Category: "productivity", Alias: []string{"Atlassian Confluence", "Confluence Wiki", "Confluence Docs"}, Parent: ""},
	{Name: "Airtable", Category: "productivity", Alias: []string{"Airtable Database", "Airtable Spreadsheets"}, Parent: ""},
	{Name: "ClickUp", Category: "productivity", Alias: []string{"ClickUp App", "ClickUp Project Management"}, Parent: ""},
	{Name: "Todoist", Category: "productivity", Alias: []string{"Todoist App", "Todoist Tasks"}, Parent: ""},
	{Name: "Basecamp", Category: "productivity", Alias: []string{"Basecamp App", "Basecamp Project Management", "Basecamp 3"}, Parent: ""},
	{Name: "Calendly", Category: "productivity", Alias: []string{"Calendly Scheduling", "Calendly App"}, Parent: ""},
	{Name: "Zoom", Category: "productivity", Alias: []string{"Zoom Meetings", "Zoom Video", "Zoom Calls", "Zoom App"}, Parent: ""},
	{Name: "Microsoft 365", Category: "productivity", Alias: []string{"M365", "Office 365", "O365", "Microsoft Subscription"}, Parent: ""},
	{Name: "OneNote", Category: "productivity", Alias: []string{"Microsoft OneNote", "MS OneNote", "OneNote App"}, Parent: ""},
	{Name: "Evernote", Category: "productivity", Alias: []string{"Evernote App", "Evernote Notes"}, Parent: ""},
	{Name: "SharePoint", Category: "productivity", Alias: []string{"Microsoft SharePoint", "SharePoint Online", "SharePoint Server"}, Parent: ""},
	{Name: "Miro", Category: "productivity", Alias: []string{"Miro Boards", "Miro App", "Miro Whiteboard"}, Parent: ""},
	{Name: "Figma", Category: "productivity", Alias: []string{"Figma Design", "Figma App", "Figma Platform"}, Parent: ""},
}

var dataScienceTechnologies = []Technology{
	{Name: "Pandas", Category: "data_science", Alias: []string{"pd", "pandas"}, Parent: "Python"},
	{Name: "NumPy", Category: "data_science", Alias: []string{"np", "numpy"}, Parent: "Python"},
	{Name: "Matplotlib", Category: "data_science", Alias: []string{"plt", "MPL"}, Parent: "Python"},
	{Name: "Tableau", Category: "data_science", Alias: []string{"Tableau Desktop", "Tableau Server", "Tableau Cloud"}, Parent: ""},
	{Name: "Power BI", Category: "data_science", Alias: []string{"Microsoft Power BI", "PowerBI", "Power BI Desktop", "Power BI Service"}, Parent: ""},
	{Name: "Excel", Category: "data_science", Alias: []string{"Microsoft Excel", "MS Excel", "Excel Spreadsheets", "Excel Analytics"}, Parent: ""},
	{Name: "R Studio", Category: "data_science", Alias: []string{"RStudio", "Posit", "RStudio IDE", "RStudio Desktop"}, Parent: "R"},
	{Name: "Jupyter", Category: "data_science", Alias: []string{"Jupyter Notebook", "Jupyter Lab", "IPython Notebook", "Jupyter Hub"}, Parent: ""},
	{Name: "SciPy", Category: "data_science", Alias: []string{"scipy", "Scientific Python", "python-scipy"}, Parent: "Python"},
	{Name: "Seaborn", Category: "data_science", Alias: []string{"sns", "python-seaborn", "seaborn-py"}, Parent: "Matplotlib"},
	{Name: "Plotly", Category: "data_science", Alias: []string{"Plotly Express", "px", "Plotly Dash", "Plotly.js"}, Parent: ""},
	{Name: "D3.js", Category: "data_science", Alias: []string{"D3", "Data-Driven Documents", "d3js", "D3 Visualization"}, Parent: "JavaScript"},
	{Name: "Databricks", Category: "data_science", Alias: []string{"Databricks Platform", "Databricks Workspace", "Databricks Notebooks", "Databricks Lakehouse"}, Parent: ""},
	{Name: "dbt", Category: "data_science", Alias: []string{"dbt Core", "data build tool", "dbt Cloud", "dbt Labs"}, Parent: ""},
	{Name: "Looker", Category: "data_science", Alias: []string{"Looker Studio", "Google Looker", "Looker Analytics", "Data Studio"}, Parent: "Google Cloud Platform"},
	{Name: "Metabase", Category: "data_science", Alias: []string{"Metabase BI", "Metabase Analytics", "Metabase Dashboard"}, Parent: ""},
	{Name: "Mode", Category: "data_science", Alias: []string{"Mode Analytics", "Mode BI", "Mode Dashboard"}, Parent: ""},
	{Name: "Redash", Category: "data_science", Alias: []string{"Redash Analytics", "Redash Dashboards", "Redash Query"}, Parent: ""},
	{Name: "Apache Spark", Category: "data_science", Alias: []string{"Spark", "PySpark", "Spark SQL", "Spark MLlib"}, Parent: ""},
	{Name: "KNIME", Category: "data_science", Alias: []string{"KNIME Analytics Platform", "KNIME Workbench", "KNIME Server"}, Parent: ""},
	{Name: "RapidMiner", Category: "data_science", Alias: []string{"RapidMiner Studio", "RapidMiner Server", "RapidMiner Platform"}, Parent: ""},
	{Name: "Orange", Category: "data_science", Alias: []string{"Orange Data Mining", "Orange Canvas", "Orange3"}, Parent: ""},
	{Name: "Alteryx", Category: "data_science", Alias: []string{"Alteryx Designer", "Alteryx Server", "Alteryx Analytics", "Alteryx Platform"}, Parent: ""},
	{Name: "QlikView", Category: "data_science", Alias: []string{"Qlik", "QlikSense", "Qlik Analytics", "QlikTech"}, Parent: ""},
}

var messagingTechnologies = []Technology{
	{Name: "Kafka", Category: "messaging", Alias: []string{"Apache Kafka", "Kafka Streams", "Kafka Connect", "Event Streaming Platform"}, Parent: ""},
	{Name: "RabbitMQ", Category: "messaging", Alias: []string{"Rabbit", "RMQ", "Rabbit Message Queue", "AMQP Broker"}, Parent: ""},
	{Name: "ActiveMQ", Category: "messaging", Alias: []string{"Apache ActiveMQ", "ActiveMQ Artemis", "AMQ"}, Parent: ""},
	{Name: "ZeroMQ", Category: "messaging", Alias: []string{"0MQ", "ØMQ", "ZMQ", "zeromq"}, Parent: ""},
	{Name: "Apache Pulsar", Category: "messaging", Alias: []string{"Pulsar", "Apache Pulsar Service", "Pulsar Functions"}, Parent: ""},
	{Name: "NATS", Category: "messaging", Alias: []string{"NATS Server", "NATS Streaming", "NATS.io", "NATS Jetstream"}, Parent: ""},
	{Name: "Redis Pub/Sub", Category: "messaging", Alias: []string{"Redis PubSub", "Redis Messaging", "Redis Streams"}, Parent: "Redis"},
	{Name: "AWS SQS", Category: "messaging", Alias: []string{"Simple Queue Service", "Amazon SQS", "SQS", "AWS Simple Queue Service"}, Parent: "Amazon Web Services"},
	{Name: "AWS SNS", Category: "messaging", Alias: []string{"Simple Notification Service", "Amazon SNS", "SNS", "AWS Simple Notification Service"}, Parent: "Amazon Web Services"},
	{Name: "Google Pub/Sub", Category: "messaging", Alias: []string{"Google PubSub", "GCP Pub/Sub", "Cloud Pub/Sub", "Google Cloud Pub/Sub"}, Parent: "Google Cloud Platform"},
	{Name: "Azure Service Bus", Category: "messaging", Alias: []string{"ASB", "Microsoft Service Bus", "Azure Messaging", "Azure Message Bus"}, Parent: "Microsoft Azure"},
	{Name: "IBM MQ", Category: "messaging", Alias: []string{"WebSphere MQ", "MQ Series", "WMQ", "IBM Message Queue"}, Parent: "IBM Cloud"},
	{Name: "RocketMQ", Category: "messaging", Alias: []string{"Apache RocketMQ", "Alibaba RocketMQ", "Rocket Message Queue"}, Parent: ""},
	{Name: "MQTT", Category: "messaging", Alias: []string{"Message Queuing Telemetry Transport", "MQ Telemetry Transport", "MQTT Protocol", "IoT Messaging Protocol"}, Parent: ""},
	{Name: "AMQP", Category: "messaging", Alias: []string{"Advanced Message Queuing Protocol", "AMQP Protocol", "AMQP 0-9-1", "AMQP 1.0"}, Parent: ""},
	{Name: "Apache Camel", Category: "messaging", Alias: []string{"Camel", "Camel Integration", "Apache Integration Framework", "Camel Routes"}, Parent: ""},
	{Name: "Apache NiFi", Category: "messaging", Alias: []string{"NiFi", "Dataflow System", "NiFi Registry", "Hortonworks DataFlow"}, Parent: ""},
	{Name: "Celery", Category: "messaging", Alias: []string{"Celery Task Queue", "Python Celery", "Celery Distributed Tasks", "Celery Workers"}, Parent: "Python"},
}

var otherTechnologies = []Technology{
	{Name: "Git", Category: "other", Alias: []string{"Git SCM", "Git VCS", "Git Version Control", "Git Source Control"}, Parent: ""},
	{Name: "GitHub", Category: "other", Alias: []string{"GH", "GitHub.com", "GitHub Pages", "GitHub Actions"}, Parent: "Git"},
	{Name: "GitLab", Category: "other", Alias: []string{"GL", "GitLab.com", "GitLab CI/CD", "GitLab Runner"}, Parent: "Git"},
	{Name: "Bitbucket", Category: "other", Alias: []string{"BB", "Atlassian Bitbucket", "Bitbucket Cloud", "Bitbucket Server"}, Parent: "Git"},
	{Name: "Maven", Category: "other", Alias: []string{"Apache Maven", "MVN", "Maven Project", "Maven Build Tool"}, Parent: "Java"},
	{Name: "Gradle", Category: "other", Alias: []string{"Gradle Build Tool", "Gradle Wrapper", "Gradle DSL", "Groovy DSL"}, Parent: "Java"},
	{Name: "npm", Category: "other", Alias: []string{"Node Package Manager", "npm CLI", "npm Registry", "npmjs"}, Parent: "Node.js"},
	{Name: "yarn", Category: "other", Alias: []string{"Yarn Package Manager", "Yarn Classic", "Yarn Berry", "Yarn PnP"}, Parent: "npm"},
	{Name: "pnpm", Category: "other", Alias: []string{"Performant npm", "Fast npm", "pnpm Package Manager"}, Parent: "npm"},
	{Name: "Babel", Category: "other", Alias: []string{"Babel.js", "Babel Compiler", "Babel Transpiler", "JavaScript Transpiler"}, Parent: "JavaScript"},
	{Name: "ESLint", Category: "other", Alias: []string{"JavaScript Linter", "JS Linter", "ECMAScript Linter"}, Parent: "JavaScript"},
	{Name: "Prettier", Category: "other", Alias: []string{"Prettier Formatter", "Code Formatter", "JavaScript Formatter"}, Parent: "JavaScript"},
	{Name: "Sketch", Category: "other", Alias: []string{"Sketch App", "Sketch Design", "Sketch for Mac", "Bohemian Coding"}, Parent: ""},
	{Name: "WordPress", Category: "other", Alias: []string{"WP", "WordPress CMS", "WordPress.org", "WordPress.com"}, Parent: "PHP"},
	{Name: "Drupal", Category: "other", Alias: []string{"Drupal CMS", "Drupal Core", "Drupal.org", "Acquia Drupal"}, Parent: "PHP"},
	{Name: "Magento", Category: "other", Alias: []string{"Adobe Commerce", "Magento Commerce", "Magento 2", "Magento Open Source"}, Parent: "PHP"},
	{Name: "Shopify", Category: "other", Alias: []string{"Shopify Platform", "Shopify Store", "Shopify Plus", "Shopify Liquid"}, Parent: ""},
	{Name: "Auth0", Category: "other", Alias: []string{"Auth0 Identity", "Auth0 IAM", "Okta Auth0", "Auth0 Authentication"}, Parent: ""},
	{Name: "Okta", Category: "other", Alias: []string{"Okta Identity", "Okta SSO", "Okta IAM", "Okta Platform"}, Parent: ""},
	{Name: "Nginx", Category: "other", Alias: []string{"nginx", "NGINX Web Server", "NGINX Proxy", "NGINX Plus"}, Parent: ""},
	{Name: "Apache", Category: "other", Alias: []string{"Apache HTTP Server", "httpd", "Apache Web Server", "Apache HTTPD"}, Parent: ""},
	{Name: "Tomcat", Category: "other", Alias: []string{"Apache Tomcat", "Tomcat Server", "Jakarta Tomcat", "Tomcat Servlet Container"}, Parent: "Java"},
	{Name: "WebAssembly", Category: "other", Alias: []string{"WASM", "Wasm", "Web Assembly", "WASI"}, Parent: ""},
	{Name: "Web3.js", Category: "other", Alias: []string{"Web3", "Ethereum JavaScript API", "Web3 Library", "ETH JavaScript Library"}, Parent: "JavaScript"},
	{Name: "Solidity", Category: "other", Alias: []string{"Solidity Language", "Ethereum Smart Contracts", "Sol", ".sol"}, Parent: ""},
	{Name: "Ethereum", Category: "other", Alias: []string{"ETH", "Ethereum Blockchain", "Ethereum Network", "Ether"}, Parent: ""},
}
