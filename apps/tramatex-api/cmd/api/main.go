package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joran-cortez/tramatex/internal/shared/infrastructure/config"
	"github.com/joran-cortez/tramatex/internal/shared/infrastructure/migrations"

	// IAM Module Imports
	iam_usecase "github.com/joran-cortez/tramatex/internal/iam/application/usecase"
	iam_repo "github.com/joran-cortez/tramatex/internal/iam/infrastructure/persistence"
	iam_handler "github.com/joran-cortez/tramatex/internal/iam/interfaces/http/handler"

	// Party Module Imports
	party_uc "github.com/joran-cortez/tramatex/internal/party/application"
	party_handler "github.com/joran-cortez/tramatex/internal/party/interfaces"
	party_repo "github.com/joran-cortez/tramatex/internal/party/persistence"

	// Security Service & Middleware Import
	"github.com/joran-cortez/tramatex/internal/shared/infrastructure/security"
	"github.com/joran-cortez/tramatex/internal/shared/interfaces/http/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	fmt.Println("🚀 TramaTex API starting...")
	fmt.Printf("Server: %s:%s\n", cfg.Server.Host, cfg.Server.Port)
	fmt.Printf("Database: %s\n", cfg.DB.Host)

	// Connect to database
	db, err := cfg.DB.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	fmt.Println("✓ Database connected")

	// Get underlying SQL DB from GORM for repositories that need *sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get SQL database instance: %v", err)
	}

	// Run migrations
	fmt.Println("🔄 Starting migrations...")
	if err := migrations.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	fmt.Println("✓ Migrations completed")

	// Setup Gin router
	router := gin.Default()

	// CORS middleware
	router.Use(corsMiddleware())

	// =========================================================================
	// DEPENDENCY INJECTION & ROUTE SETUP
	// =========================================================================

	// --- IAM Module Dependencies ---
	// 1. Repository
	userRepository := iam_repo.NewPostgresUserRepository(db)
	// 2. Security Service (JWT)
	jwtService, err := security.NewJWTService(cfg.Security.JWTSecret, cfg.Security.JWTAccessTTL, cfg.Security.JWTRefreshTTL)
	if err != nil {
		log.Fatalf("Failed to initialize JWT service: %v", err)
	}
	// 3. Use Cases
	loginUseCase := iam_usecase.NewLoginUseCase(userRepository, jwtService)
	// 4. HTTP Handler
	iamHandler := iam_handler.NewIAMHandler(loginUseCase)

	// --- Party Module Dependencies ---
	// 1. Repositories
	organizationRepo := party_repo.NewPostgreSQLOrganizationRepository(sqlDB)
	personRepo := party_repo.NewPostgreSQLPersonRepository(sqlDB)
	addressRepo := party_repo.NewPostgreSQLAddressRepository(sqlDB)

	// 2. Use Cases
	createOrgHandler := party_uc.NewCreateOrganizationHandler(organizationRepo)
	updateOrgHandler := party_uc.NewUpdateOrganizationHandler(organizationRepo)
	changeOrgStatusHandler := party_uc.NewChangeOrganizationStatusHandler(organizationRepo)
	getOrgHandler := party_uc.NewGetOrganizationHandler(organizationRepo)
	listOrgsHandler := party_uc.NewListOrganizationsHandler(organizationRepo)
	listOrgsByRoleHandler := party_uc.NewListOrganizationsByRoleHandler(organizationRepo)

	addPersonHandler := party_uc.NewAddPersonHandler(organizationRepo, personRepo)
	getPersonHandler := party_uc.NewGetPersonHandler(personRepo)
	listPersonsHandler := party_uc.NewListPersonsByOrganizationHandler(personRepo)
	getPersonByEmailHandler := party_uc.NewGetPersonByEmailHandler(personRepo)
	getPrimaryContactHandler := party_uc.NewGetPrimaryContactHandler(personRepo)

	addAddressHandler := party_uc.NewAddAddressHandler(organizationRepo, addressRepo)
	listAddressesHandler := party_uc.NewListAddressesByOrganizationHandler(addressRepo)
	getPrimaryAddressHandler := party_uc.NewGetPrimaryAddressHandler(addressRepo)

	// 3. HTTP Handlers
	orgHandler := party_handler.NewOrganizationHandler(
		createOrgHandler,
		updateOrgHandler,
		changeOrgStatusHandler,
		getOrgHandler,
		listOrgsHandler,
		listOrgsByRoleHandler,
	)
	personHandler := party_handler.NewPersonHandler(
		addPersonHandler,
		getPersonHandler,
		listPersonsHandler,
		getPersonByEmailHandler,
		getPrimaryContactHandler,
	)
	addressHandler := party_handler.NewAddressHandler(
		addAddressHandler,
		listAddressesHandler,
		getPrimaryAddressHandler,
	)

	// --- Middleware ---
	authMiddleware := middleware.AuthMiddleware(jwtService)

	// --- API Routes ---
	api := router.Group("/api")
	{
		// Public routes
		api.POST("/iam/login", iamHandler.Login)
		api.GET("/health", iamHandler.Health)

		// Protected routes
		protected := api.Group("/")
		protected.Use(authMiddleware)
		{
			organizations := protected.Group("/organizations")
			{
				organizations.POST("", orgHandler.CreateOrganization)
				organizations.GET("", orgHandler.ListOrganizations)
				organizations.GET("/:id", orgHandler.GetOrganization)
				organizations.PUT("/:id", orgHandler.UpdateOrganization)
				organizations.PATCH("/:id/status", orgHandler.ChangeStatus)

				organizations.POST("/:org_id/persons", personHandler.AddPerson)
				organizations.GET("/:org_id/persons", personHandler.ListPersons)
				organizations.GET("/:org_id/primary-contact", personHandler.GetPrimaryContact)

				organizations.POST("/:org_id/addresses", addressHandler.AddAddress)
				organizations.GET("/:org_id/addresses", addressHandler.ListAddresses)
			}

			persons := protected.Group("/persons")
			{
				persons.GET("/:id", personHandler.GetPerson)
			}
		}
	}
	// =========================================================================

	// Start server
	if err := router.Run(cfg.Server.Host + ":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// corsMiddleware adds CORS headers
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-User-ID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
