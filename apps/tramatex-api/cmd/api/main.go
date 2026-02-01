package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joran-cortez/tramatex/internal/shared/infrastructure/config"
	"github.com/joran-cortez/tramatex/internal/shared/infrastructure/logging"
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
	infra_middleware "github.com/joran-cortez/tramatex/internal/shared/infrastructure/middleware"
	"github.com/joran-cortez/tramatex/internal/shared/infrastructure/security"
	"github.com/joran-cortez/tramatex/internal/shared/interfaces/http/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logging.InitLogger(cfg.Server.Environment)
	logging.Logger.Info("🚀 TramaTex API starting...")
	logging.Logger.WithFields(map[string]interface{}{
		"host":        cfg.Server.Host,
		"port":        cfg.Server.Port,
		"environment": cfg.Server.Environment,
	}).Info("Server configuration loaded")

	fmt.Println("🚀 TramaTex API starting...")
	fmt.Printf("Server: %s:%s\n", cfg.Server.Host, cfg.Server.Port)
	fmt.Printf("Database: %s\n", cfg.DB.Host)

	// Connect to database
	db, err := cfg.DB.Connect()
	if err != nil {
		logging.Logger.WithError(err).Fatal("Failed to connect to database")
	}
	logging.Logger.Info("✓ Database connected")
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

	// Global middlewares (order matters!)
	router.Use(infra_middleware.SecurityHeadersMiddleware())                     // 1. Add security headers
	router.Use(infra_middleware.RequestIDMiddleware())                           // 2. Generate request ID
	router.Use(infra_middleware.SecurityLoggerMiddleware())                      // 3. Log all requests with security context
	router.Use(infra_middleware.CORSMiddleware(cfg.Security.CORSAllowedOrigins)) // 4. CORS with whitelist
	router.Use(infra_middleware.ErrorHandlerMiddleware(cfg.Server.Environment))  // 5. Error sanitization

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
	// 2.1 Token Blacklist (logout invalidation)
	tokenBlacklist := security.NewPostgresTokenBlacklist(db)
	// 3. Use Cases
	loginUseCase := iam_usecase.NewLoginUseCase(userRepository, jwtService)
	registerUseCase := iam_usecase.NewRegisterUserUseCase(userRepository)
	createUserUseCase := iam_usecase.NewCreateUserUseCase(userRepository)
	refreshUseCase := iam_usecase.NewRefreshTokenUseCase(userRepository, jwtService)
	logoutUseCase := iam_usecase.NewLogoutUserUseCase(jwtService, tokenBlacklist)
	assignRoleUseCase := iam_usecase.NewAssignRoleUseCase(userRepository)
	checkAuthUseCase := iam_usecase.NewCheckAuthorizationUseCase(userRepository)
	listUsersUseCase := iam_usecase.NewListUsersUseCase(userRepository)
	deleteUserUseCase := iam_usecase.NewDeleteUserUseCase(userRepository)
	// 4. HTTP Handler
	iamHandler := iam_handler.NewIAMHandler(
		loginUseCase,
		registerUseCase,
		createUserUseCase,
		refreshUseCase,
		logoutUseCase,
		assignRoleUseCase,
		checkAuthUseCase,
		listUsersUseCase,
		deleteUserUseCase,
	)

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
	authMiddleware := middleware.AuthMiddleware(jwtService, tokenBlacklist)

	// --- API Routes ---
	auth := router.Group("/auth")
	{
		// Public auth routes
		auth.POST("/register", iamHandler.Register)
		auth.POST("/login", infra_middleware.RateLimitMiddleware(10, time.Minute), iamHandler.Login)
		auth.POST("/refresh", iamHandler.Refresh)

		// Protected auth routes
		protectedAuth := auth.Group("/")
		protectedAuth.Use(authMiddleware)
		{
			// Logout requires Authorization header
			protectedAuth.POST("/logout", iamHandler.Logout)
			// Admin only: create user
			protectedAuth.POST("/users", infra_middleware.RequireRole("admin"), iamHandler.CreateUser)
			// Admin only: assign role
			protectedAuth.POST("/assign-role", infra_middleware.RequireRole("admin"), iamHandler.AssignRole)
			// Admin only: list users
			protectedAuth.GET("/users", infra_middleware.RequireRole("admin"), iamHandler.ListUsers)
			// Admin only: delete user
			protectedAuth.DELETE("/users/:id", infra_middleware.RequireRole("admin"), iamHandler.DeleteUser)
			// Authorization checks (authenticated)
			protectedAuth.POST("/authorize", iamHandler.CheckAuthorization)
		}
	}

	api := router.Group("/api")
	{
		// Public routes
		api.GET("/health", iamHandler.Health)

		// Protected routes
		protected := api.Group("/")
		protected.Use(authMiddleware)
		{
			organizations := protected.Group("/organizations")
			{
				// Write operations: admin and commercial only
				organizations.POST("", infra_middleware.RequireRole("admin", "commercial"), orgHandler.CreateOrganization)
				organizations.PUT("/:id", infra_middleware.RequireRole("admin", "commercial"), orgHandler.UpdateOrganization)
				organizations.PATCH("/:id/status", infra_middleware.RequireRole("admin", "commercial"), orgHandler.ChangeStatus)

				// Read operations: all authenticated users
				organizations.GET("", orgHandler.ListOrganizations)
				organizations.GET("/:id", orgHandler.GetOrganization)

				// Person operations
				organizations.POST("/:id/persons", infra_middleware.RequireRole("admin", "commercial"), personHandler.AddPerson)
				organizations.GET("/:id/persons", personHandler.ListPersons)
				organizations.GET("/:id/primary-contact", personHandler.GetPrimaryContact)

				// Address operations
				organizations.POST("/:id/addresses", infra_middleware.RequireRole("admin", "commercial"), addressHandler.AddAddress)
				organizations.GET("/:id/addresses", addressHandler.ListAddresses)
			}

			persons := protected.Group("/persons")
			{
				persons.GET("/:id", personHandler.GetPerson)
			}
		}
	}
	// =========================================================================

	// Start server
	logging.Logger.WithFields(map[string]interface{}{
		"host": cfg.Server.Host,
		"port": cfg.Server.Port,
	}).Info("Server starting")

	if err := router.Run(cfg.Server.Host + ":" + cfg.Server.Port); err != nil {
		logging.Logger.WithError(err).Fatal("Failed to start server")
	}
}
