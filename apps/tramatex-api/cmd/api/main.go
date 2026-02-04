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

	// Product Module Imports
	product_uc "github.com/joran-cortez/tramatex/internal/product/application"
	product_repo "github.com/joran-cortez/tramatex/internal/product/infrastructure/persistence"
	product_handler "github.com/joran-cortez/tramatex/internal/product/interfaces/http/handler"

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
	partyRepo := party_repo.NewPostgreSQLPartyRepository(sqlDB)
	partyRelationshipRepo := party_repo.NewPostgreSQLPartyRelationshipRepository(sqlDB)
	partyAddressRepo := party_repo.NewPostgreSQLPartyAddressRepository(sqlDB)

	// 2. Use Cases
	createPartyHandler := party_uc.NewCreatePartyHandler(partyRepo)
	updatePartyHandler := party_uc.NewUpdatePartyHandler(partyRepo)
	changePartyStatusHandler := party_uc.NewChangePartyStatusHandler(partyRepo)
	getPartyHandler := party_uc.NewGetPartyHandler(partyRepo)
	listPartiesHandler := party_uc.NewListPartiesHandler(partyRepo)
	addPartyRoleHandler := party_uc.NewAddPartyRoleHandler(partyRepo)
	removePartyRoleHandler := party_uc.NewRemovePartyRoleHandler(partyRepo)
	addPartyRelationshipHandler := party_uc.NewAddPartyRelationshipHandler(partyRelationshipRepo)
	listPartyRelationshipsHandler := party_uc.NewListPartyRelationshipsHandler(partyRelationshipRepo)
	removePartyRelationshipHandler := party_uc.NewRemovePartyRelationshipHandler(partyRelationshipRepo)
	addContactDetailsHandler := party_uc.NewAddContactDetailsHandler(partyRepo)
	updateContactDetailsHandler := party_uc.NewUpdateContactDetailsHandler(partyRepo)
	listContactDetailsHandler := party_uc.NewListContactDetailsHandler(partyRepo)
	removeContactDetailsHandler := party_uc.NewRemoveContactDetailsHandler(partyRepo)
	addPartyAddressHandler := party_uc.NewAddPartyAddressHandler(partyAddressRepo)
	listPartyAddressesHandler := party_uc.NewListPartyAddressesHandler(partyAddressRepo)

	// 3. HTTP Handlers
	partyHandler := party_handler.NewPartyHandler(
		createPartyHandler,
		updatePartyHandler,
		changePartyStatusHandler,
		getPartyHandler,
		listPartiesHandler,
	)
	partyRoleHandler := party_handler.NewPartyRoleHandler(addPartyRoleHandler, removePartyRoleHandler)
	partyRelationshipHandler := party_handler.NewPartyRelationshipHandler(
		addPartyRelationshipHandler,
		listPartyRelationshipsHandler,
		removePartyRelationshipHandler,
	)
	contactDetailsHandler := party_handler.NewContactDetailsHandler(
		addContactDetailsHandler,
		updateContactDetailsHandler,
		listContactDetailsHandler,
		removeContactDetailsHandler,
	)
	partyAddressHandler := party_handler.NewPartyAddressHandler(addPartyAddressHandler, listPartyAddressesHandler)

	// --- Product Module Dependencies ---
	// 1. Repositories
	productRepository := product_repo.NewGORMProductRepository(db)
	attributeRepository := product_repo.NewGORMAttributeRepository(db)
	productVariantRepository := product_repo.NewGORMVariantRepository(db)
	brandRepository := product_repo.NewGORMBrandRepository(db)
	productGroupRepository := product_repo.NewGORMProductGroupRepository(db)

	// 2. Use Cases
	productService := product_uc.NewProductService(productRepository, brandRepository, productGroupRepository, attributeRepository, productVariantRepository)

	// 3. HTTP Handlers
	productHandler := product_handler.NewProductHandler(productService)

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
			parties := protected.Group("/parties")
			{
				parties.POST("", infra_middleware.RequireRole("admin", "commercial"), partyHandler.CreateParty)
				parties.PUT("/:id", infra_middleware.RequireRole("admin", "commercial"), partyHandler.UpdateParty)
				parties.PATCH("/:id/status", infra_middleware.RequireRole("admin", "commercial"), partyHandler.ChangePartyStatus)

				parties.GET("", partyHandler.ListParties)
				parties.GET("/:id", partyHandler.GetParty)

				parties.POST("/:id/roles", infra_middleware.RequireRole("admin", "commercial"), partyRoleHandler.AddRole)
				parties.DELETE("/:id/roles/:role", infra_middleware.RequireRole("admin", "commercial"), partyRoleHandler.RemoveRole)

				parties.POST("/:id/relationships", infra_middleware.RequireRole("admin", "commercial"), partyRelationshipHandler.AddRelationship)
				parties.GET("/:id/relationships", partyRelationshipHandler.ListRelationships)
				parties.DELETE("/:id/relationships/:relationship_id", infra_middleware.RequireRole("admin", "commercial"), partyRelationshipHandler.RemoveRelationship)

				parties.POST("/:id/contact-details", infra_middleware.RequireRole("admin", "commercial"), contactDetailsHandler.AddContactDetails)
				parties.GET("/:id/contact-details", contactDetailsHandler.ListContactDetails)
				parties.PUT("/:id/contact-details/:contact_id", infra_middleware.RequireRole("admin", "commercial"), contactDetailsHandler.UpdateContactDetails)
				parties.DELETE("/:id/contact-details/:contact_id", infra_middleware.RequireRole("admin", "commercial"), contactDetailsHandler.RemoveContactDetails)

				parties.POST("/:id/addresses", infra_middleware.RequireRole("admin", "commercial"), partyAddressHandler.AddAddress)
				parties.GET("/:id/addresses", partyAddressHandler.ListAddresses)

				// New PartyServiceConfiguration routes
				parties.POST("/:partyId/service-configurations", infra_middleware.RequireRole("admin", "commercial"), productHandler.CreatePartyServiceConfiguration)
				parties.GET("/:partyId/service-configurations", productHandler.ListPartyServiceConfigurationsByPartyID)
				parties.GET("/:partyId/service-configurations/:id", productHandler.GetPartyServiceConfigurationByID)
				parties.PUT("/:partyId/service-configurations/:id", infra_middleware.RequireRole("admin", "commercial"), productHandler.UpdatePartyServiceConfiguration)
				parties.DELETE("/:partyId/service-configurations/:id", infra_middleware.RequireRole("admin", "commercial"), productHandler.DeletePartyServiceConfiguration)
			}

			products := protected.Group("/products")
			{
				products.POST("", infra_middleware.RequireRole("admin", "commercial"), productHandler.CreateProduct)
				products.GET("", productHandler.ListProducts) // New
				products.GET("/:id", productHandler.GetProductByID) // New
				products.GET("/:id/calculated-option-sets", productHandler.GetCalculatedOptionSetsForProduct) // New
				products.POST("/:id/groups", infra_middleware.RequireRole("admin", "commercial"), productHandler.AddGroupToProduct)
				products.POST("/:id/attributes", infra_middleware.RequireRole("admin", "commercial"), productHandler.AddDirectAttributeToProduct)
				products.PATCH("/:id/sku", infra_middleware.RequireRole("admin", "commercial"), productHandler.UpdateProductSKU)

				// Product Variant routes nested under product
				products.POST("/:productId/variants/generate", infra_middleware.RequireRole("admin", "commercial"), productHandler.GenerateProductVariants) // New
				products.POST("/:productId/variants/find-or-create", infra_middleware.RequireRole("admin", "commercial"), productHandler.FindOrCreateProductVariant) // New
				products.GET("/:productId/variants", productHandler.ListProductVariantsByProductID) // New
			}

			// New: Attributes (ProductOptionSet) routes
			attributes := protected.Group("/attributes")
			{
				attributes.POST("", infra_middleware.RequireRole("admin", "commercial"), productHandler.CreateAttribute) // New
				attributes.GET("", productHandler.ListAttributes) // New
				attributes.GET("/:id", productHandler.GetAttributeByID) // New
				attributes.PUT("/:id", infra_middleware.RequireRole("admin", "commercial"), productHandler.UpdateAttribute) // New
			}

			// New: Top-level Product Variant routes
			variants := protected.Group("/variants")
			{
				variants.GET("/:id", productHandler.GetProductVariantByID) // New
				variants.GET("", productHandler.GetProductVariantBySKU) // New (using query param sku)
				variants.PUT("/:id", infra_middleware.RequireRole("admin", "commercial"), productHandler.UpdateProductVariant) // New
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
