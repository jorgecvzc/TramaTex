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

	// Pricing Module Imports
	pricing_uc "github.com/joran-cortez/tramatex/internal/pricing/application"
	pricing_cache "github.com/joran-cortez/tramatex/internal/pricing/infrastructure/cache"
	pricing_repo "github.com/joran-cortez/tramatex/internal/pricing/infrastructure/persistence"
	pricing_productclient "github.com/joran-cortez/tramatex/internal/pricing/infrastructure/productclient"
	pricing_handler "github.com/joran-cortez/tramatex/internal/pricing/interfaces/http/handler"
	"github.com/redis/go-redis/v9"

	// MES Module Imports
	mes_uc "github.com/joran-cortez/tramatex/internal/mes/application"
	mes_repo "github.com/joran-cortez/tramatex/internal/mes/infrastructure/persistence"
	mes_handler "github.com/joran-cortez/tramatex/internal/mes/interfaces/http/handler"

	// Sales Module Imports
	sales_uc "github.com/joran-cortez/tramatex/internal/sales/application"
	sales_repo "github.com/joran-cortez/tramatex/internal/sales/infrastructure/persistence"
	sales_handler "github.com/joran-cortez/tramatex/internal/sales/interfaces/http/handler"

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
	partyRepo := party_repo.NewGORMPartyRepository(db)
	partyRelationshipRepo := party_repo.NewGORMPartyRelationshipRepository(db)
	partyAddressRepo := party_repo.NewGORMPartyAddressRepository(db)

	// 2. Use Cases
	createPartyHandler := party_uc.NewCreatePartyHandler(partyRepo)
	updatePartyHandler := party_uc.NewUpdatePartyHandler(partyRepo)
	changePartyStatusHandler := party_uc.NewChangePartyStatusHandler(partyRepo)
	deletePartyHandler := party_uc.NewDeletePartyHandler(partyRepo, partyRelationshipRepo)
	getPartyHandler := party_uc.NewGetPartyHandler(partyRepo)
	listPartiesHandler := party_uc.NewListPartiesHandler(partyRepo)
	getPartiesBatchHandler := party_uc.NewGetPartiesBatchHandler(partyRepo)
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
	updatePartyAddressHandler := party_uc.NewUpdatePartyAddressHandler(partyAddressRepo)
	removePartyAddressHandler := party_uc.NewRemovePartyAddressHandler(partyAddressRepo)
	listPartyAddressesHandler := party_uc.NewListPartyAddressesHandler(partyAddressRepo)

	// 3. HTTP Handlers
	partyHandler := party_handler.NewPartyHandler(
		createPartyHandler,
		updatePartyHandler,
		changePartyStatusHandler,
		deletePartyHandler,
		getPartyHandler,
		listPartiesHandler,
		getPartiesBatchHandler,
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
	partyAddressHandler := party_handler.NewPartyAddressHandler(
		addPartyAddressHandler,
		updatePartyAddressHandler,
		removePartyAddressHandler,
		listPartyAddressesHandler,
	)

	// --- Product Module Dependencies ---
	// 1. Repositories
	productRepository := product_repo.NewGORMProductRepository(db)
	attributeRepository := product_repo.NewGORMAttributeRepository(db)
	productVariantRepository := product_repo.NewGORMVariantRepository(db)
	brandRepository := product_repo.NewGORMBrandRepository(db)
	productGroupRepository := product_repo.NewGORMProductGroupRepository(db)
	partyServiceConfigurationRepository := product_repo.NewGORMPartyServiceConfigurationRepository(db)

	// 2. Use Cases
	productService := product_uc.NewProductService(productRepository, brandRepository, productGroupRepository, attributeRepository, productVariantRepository, partyServiceConfigurationRepository)

	// 3. HTTP Handlers
	productHandler := product_handler.NewProductHandler(productService)

	// --- Pricing Module Dependencies ---
	pricingRuleRepo := pricing_repo.NewGORMPricingRuleRepository(db)
	clientPricingRepo := pricing_repo.NewGORMClientPricingRepository(db)
	brandMarginRepo := pricing_repo.NewGORMBrandProfitMarginRepository(db)
	discountRuleRepo := pricing_repo.NewGORMSalesDiscountRuleRepository(db)
	calculationRepo := pricing_repo.NewGORMPriceCalculationRepository(db)
	productPricingClient := pricing_productclient.NewProductPricingClient(db)
	pricingService := pricing_uc.NewPricingService(
		pricingRuleRepo,
		clientPricingRepo,
		brandMarginRepo,
		discountRuleRepo,
		calculationRepo,
		productPricingClient,
	)
	pricingHandler := pricing_handler.NewPricingHandler(pricingService)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	basePriceCache := pricing_cache.NewRedisBasePriceCache(redisClient, cfg.Redis.TTL)
	baseRuleRepo := pricing_repo.NewGORMBaseSalesPriceRuleRepository(db)
	saleRuleRepo := pricing_repo.NewGORMSaleModificationRuleRepository(db)
	pricingEngineService := pricing_uc.NewPricingEngineService(
		baseRuleRepo,
		saleRuleRepo,
		productPricingClient,
		basePriceCache,
	)
	pricingEngineHandler := pricing_handler.NewPricingEngineHandler(pricingEngineService)

	// --- Sales Module Dependencies ---
	quoteRepo := sales_repo.NewGORMQuoteRepository(db)
	orderRepo := sales_repo.NewGORMSalesOrderRepository(db)
	deliveryRepo := sales_repo.NewGORMDeliveryNoteRepository(db)
	invoiceRepo := sales_repo.NewGORMInvoiceRepository(db)
	numberGenerator := sales_repo.NewTimeBasedNumberGenerator()
	partyLookup := sales_repo.NewPartyLookupAdapter(partyRepo)
	salesService := sales_uc.NewSalesService(
		quoteRepo,
		orderRepo,
		deliveryRepo,
		invoiceRepo,
		numberGenerator,
		pricingEngineService,
		partyLookup,
	)
	salesHandler := sales_handler.NewSalesHandler(salesService)

	// --- MES Module Dependencies ---
	taskRepo := mes_repo.NewGORMTaskRepository(db)
	positionRepo := mes_repo.NewGORMPositionRepository(db)
	serviceGroupRepo := mes_repo.NewGORMServiceGroupRepository(db)
	mesWorkRepo := mes_repo.NewGORMMESWorkRepository(db)
	mesService := mes_uc.NewMESService(taskRepo, positionRepo, serviceGroupRepo, mesWorkRepo)
	mesHandler := mes_handler.NewMESHandler(mesService)

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
				parties.DELETE("/:id", infra_middleware.RequireRole("admin", "commercial"), partyHandler.DeleteParty)

				parties.GET("", partyHandler.ListParties)
				parties.GET("/batch", partyHandler.GetPartiesBatch)
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
				parties.PUT("/:id/addresses/:addressId", infra_middleware.RequireRole("admin", "commercial"), partyAddressHandler.UpdateAddress)
				parties.DELETE("/:id/addresses/:addressId", infra_middleware.RequireRole("admin", "commercial"), partyAddressHandler.DeleteAddress)

				// New PartyServiceConfiguration routes
				parties.POST("/:id/service-configurations", infra_middleware.RequireRole("admin", "commercial"), productHandler.CreatePartyServiceConfiguration)
				parties.GET("/:id/service-configurations", productHandler.ListPartyServiceConfigurationsByPartyID)
				parties.GET("/:id/service-configurations/:configId", productHandler.GetPartyServiceConfigurationByID)
				parties.PUT("/:id/service-configurations/:configId", infra_middleware.RequireRole("admin", "commercial"), productHandler.UpdatePartyServiceConfiguration)
				parties.DELETE("/:id/service-configurations/:configId", infra_middleware.RequireRole("admin", "commercial"), productHandler.DeletePartyServiceConfiguration)
			}

			products := protected.Group("/products")
			{
				products.POST("", infra_middleware.RequireRole("admin", "commercial"), productHandler.CreateProduct)
				products.GET("", productHandler.ListProducts)       // New
				products.GET("/:id", productHandler.GetProductByID) // New
				products.PUT("/:id", infra_middleware.RequireRole("admin", "commercial"), productHandler.UpdateProduct)
				products.GET("/:id/calculated-option-sets", productHandler.GetCalculatedOptionSetsForProduct) // New
				products.POST("/:id/groups", infra_middleware.RequireRole("admin", "commercial"), productHandler.AddGroupToProduct)
				products.POST("/:id/attributes", infra_middleware.RequireRole("admin", "commercial"), productHandler.AddDirectAttributeToProduct)
				products.POST("/:id/direct-option-sets", infra_middleware.RequireRole("admin", "commercial"), productHandler.AddDirectAttributeToProduct)
				products.PATCH("/:id/sku", infra_middleware.RequireRole("admin", "commercial"), productHandler.UpdateProductSKU)

				// Product Variant routes nested under product
				products.POST("/:id/variants/generate", infra_middleware.RequireRole("admin", "commercial"), productHandler.GenerateProductVariants)          // New
				products.POST("/:id/variants/find-or-create", infra_middleware.RequireRole("admin", "commercial"), productHandler.FindOrCreateProductVariant) // New
				products.GET("/:id/variants", productHandler.ListProductVariantsByProductID)                                                                  // New
			}

			// New: Attributes (ProductOptionSet) routes
			attributes := protected.Group("/attributes")
			{
				attributes.POST("", infra_middleware.RequireRole("admin", "commercial"), productHandler.CreateAttribute)    // New
				attributes.GET("", productHandler.ListAttributes)                                                           // New
				attributes.GET("/:id", productHandler.GetAttributeByID)                                                     // New
				attributes.PUT("/:id", infra_middleware.RequireRole("admin", "commercial"), productHandler.UpdateAttribute) // New
				attributes.DELETE("/:id", infra_middleware.RequireRole("admin"), productHandler.DeleteAttribute)            // New
			}

			// API contract alias: ProductOptionSet
			productOptionSets := protected.Group("/product-option-sets")
			{
				productOptionSets.POST("", infra_middleware.RequireRole("admin", "commercial"), productHandler.CreateAttribute)
				productOptionSets.GET("", productHandler.ListAttributes)
				productOptionSets.GET("/:id", productHandler.GetAttributeByID)
				productOptionSets.PUT("/:id", infra_middleware.RequireRole("admin", "commercial"), productHandler.UpdateAttribute)
			}

			// New: Top-level Product Variant routes
			variants := protected.Group("/variants")
			{
				variants.GET("/:id", productHandler.GetProductVariantByID)                                                     // New
				variants.GET("", productHandler.GetProductVariantBySKU)                                                        // New (using query param sku)
				variants.PUT("/:id", infra_middleware.RequireRole("admin", "commercial"), productHandler.UpdateProductVariant) // New
			}

			// Brand routes
			brands := protected.Group("/brands")
			{
				brands.GET("", productHandler.ListBrands)
				brands.GET("/:id", productHandler.GetBrandByID)
				brands.POST("", infra_middleware.RequireRole("admin", "commercial"), productHandler.CreateBrand)
				brands.PUT("/:id", infra_middleware.RequireRole("admin", "commercial"), productHandler.UpdateBrand)
				brands.DELETE("/:id", infra_middleware.RequireRole("admin"), productHandler.DeleteBrand)
			}

			// Product Group routes
			productGroups := protected.Group("/product-groups")
			{
				productGroups.GET("", productHandler.ListProductGroups)
				productGroups.GET("/:id", productHandler.GetProductGroupByID)
				productGroups.POST("", infra_middleware.RequireRole("admin", "commercial"), productHandler.CreateProductGroup)
				productGroups.PUT("/:id", infra_middleware.RequireRole("admin", "commercial"), productHandler.UpdateProductGroup)
				productGroups.DELETE("/:id", infra_middleware.RequireRole("admin"), productHandler.DeleteProductGroup)
			}

			pricing := protected.Group("/pricing")
			{
				pricing.POST("/calculate", pricingHandler.CalculatePrice)
				pricing.GET("/rules", pricingHandler.ListPricingRules)
				pricing.POST("/rules", infra_middleware.RequireRole("admin", "commercial"), pricingHandler.CreatePricingRule)
				pricing.POST("/client-overrides", infra_middleware.RequireRole("admin", "commercial"), pricingHandler.CreateClientPricingOverride)
				pricing.GET("/history/:variantId", pricingHandler.GetPricingHistory)

				pricing.POST("/base-sales-rules", infra_middleware.RequireRole("admin", "commercial"), pricingEngineHandler.CreateBaseSalesPriceRule)
				pricing.PUT("/base-sales-rules/:id", infra_middleware.RequireRole("admin", "commercial"), pricingEngineHandler.UpdateBaseSalesPriceRule)
				pricing.POST("/sale-modification-rules", infra_middleware.RequireRole("admin", "commercial"), pricingEngineHandler.CreateSaleModificationRule)
				pricing.PUT("/sale-modification-rules/:id", infra_middleware.RequireRole("admin", "commercial"), pricingEngineHandler.UpdateSaleModificationRule)
				pricing.POST("/base-sales-price/calculate", pricingEngineHandler.CalculateBaseSalesPrice)
				pricing.POST("/final-sale-price/calculate", pricingEngineHandler.CalculateFinalSalePrice)
			}

			sales := protected.Group("/sales")
			{
				quotes := sales.Group("/quotes")
				{
					quotes.POST("", infra_middleware.RequireRole("admin", "commercial"), salesHandler.CreateQuote)
					quotes.GET("", salesHandler.ListQuotes)
					quotes.GET("/:id", salesHandler.GetQuote)
					quotes.PUT("/:id", infra_middleware.RequireRole("admin", "commercial"), salesHandler.UpdateQuote)
					quotes.PATCH("/:id/status", infra_middleware.RequireRole("admin", "commercial"), salesHandler.ChangeQuoteStatus)
					quotes.POST("/:id/convert", infra_middleware.RequireRole("admin", "commercial"), salesHandler.ConvertQuoteToOrder)
				}

				orders := sales.Group("/orders")
				{
					orders.POST("", infra_middleware.RequireRole("admin", "commercial"), salesHandler.CreateOrder)
					orders.GET("", salesHandler.ListOrders)
					orders.GET("/:id", salesHandler.GetOrder)
					orders.PUT("/:id", infra_middleware.RequireRole("admin", "commercial"), salesHandler.UpdateOrderDetails)
					orders.PATCH("/:id/status", infra_middleware.RequireRole("admin", "commercial"), salesHandler.ChangeOrderStatus)
					orders.POST("/:id/line-items", infra_middleware.RequireRole("admin", "commercial"), salesHandler.AddOrderLineItem)
					orders.PUT("/:id/line-items/:lineItemId", infra_middleware.RequireRole("admin", "commercial"), salesHandler.UpdateOrderLineItem)
					orders.DELETE("/:id/line-items/:lineItemId", infra_middleware.RequireRole("admin", "commercial"), salesHandler.RemoveOrderLineItem)
				}

				deliveryNotes := sales.Group("/delivery-notes")
				{
					deliveryNotes.POST("", infra_middleware.RequireRole("admin", "commercial"), salesHandler.CreateDeliveryNote)
					deliveryNotes.GET("", salesHandler.ListDeliveryNotes)
					deliveryNotes.GET("/:id", salesHandler.GetDeliveryNote)
				}

				invoices := sales.Group("/invoices")
				{
					invoices.POST("", infra_middleware.RequireRole("admin", "commercial"), salesHandler.CreateInvoice)
					invoices.POST("/simplified", infra_middleware.RequireRole("admin", "commercial", "cashier"), salesHandler.CreateSimplifiedInvoice)
					invoices.GET("", salesHandler.ListInvoices)
					invoices.GET("/:id", salesHandler.GetInvoice)
				}
			}

			mes := protected.Group("/mes")
			{
				tasks := mes.Group("/tasks")
				{
					tasks.POST("", infra_middleware.RequireRole("admin", "commercial"), mesHandler.CreateTask)
					tasks.GET("", mesHandler.ListTasks)
					tasks.GET("/:id", mesHandler.GetTask)
					tasks.PUT("/:id", infra_middleware.RequireRole("admin", "commercial"), mesHandler.UpdateTask)
					tasks.DELETE("/:id", infra_middleware.RequireRole("admin"), mesHandler.DeleteTask)
				}

				positions := mes.Group("/positions")
				{
					positions.POST("", infra_middleware.RequireRole("admin", "commercial"), mesHandler.CreatePosition)
					positions.GET("", mesHandler.ListPositions)
					positions.GET("/:id", mesHandler.GetPosition)
					positions.PUT("/:id", infra_middleware.RequireRole("admin", "commercial"), mesHandler.UpdatePosition)
					positions.DELETE("/:id", infra_middleware.RequireRole("admin"), mesHandler.DeletePosition)
				}

				serviceGroups := mes.Group("/service-groups")
				{
					serviceGroups.POST("", infra_middleware.RequireRole("admin", "commercial"), mesHandler.CreateServiceGroup)
					serviceGroups.GET("", mesHandler.ListServiceGroups)
					serviceGroups.GET("/:id", mesHandler.GetServiceGroup)
					serviceGroups.PUT("/:id", infra_middleware.RequireRole("admin", "commercial"), mesHandler.UpdateServiceGroup)
					serviceGroups.DELETE("/:id", infra_middleware.RequireRole("admin"), mesHandler.DeleteServiceGroup)
				}

				serviceTemplates := mes.Group("/service-templates")
				{
					serviceTemplates.POST("", infra_middleware.RequireRole("admin", "commercial"), mesHandler.CreateServiceTemplate)
					serviceTemplates.GET("", mesHandler.ListServiceTemplates)
					serviceTemplates.GET("/:id", mesHandler.GetServiceTemplate)
					serviceTemplates.PUT("/:id", infra_middleware.RequireRole("admin", "commercial"), mesHandler.UpdateServiceTemplate)
					serviceTemplates.DELETE("/:id", infra_middleware.RequireRole("admin"), mesHandler.DeleteServiceTemplate)
				}

				works := mes.Group("/works")
				{
					works.POST("", infra_middleware.RequireRole("admin", "commercial"), mesHandler.CreateMESWork)
					works.GET("", mesHandler.ListMESWorks)
					works.GET("/dashboard/stats", mesHandler.GetMESWorkDashboardStats)
					works.GET("/overdue", mesHandler.ListOverdueMESWorks)
					works.GET("/:id", mesHandler.GetMESWork)
					works.PUT("/:id", infra_middleware.RequireRole("admin", "commercial"), mesHandler.UpdateMESWork)
					works.PATCH("/:workId/tasks/:taskId/status", infra_middleware.RequireRole("admin", "commercial", "workshop"), mesHandler.UpdateMESWorkTaskStatus)
				}

				workDefinitions := mes.Group("/work-definitions")
				{
					workDefinitions.POST("", infra_middleware.RequireRole("admin", "commercial"), mesHandler.CreateWorkDefinition)
					workDefinitions.GET("", mesHandler.ListWorkDefinitions)
					workDefinitions.GET("/dashboard/stats", mesHandler.GetWorkDefinitionDashboardStats)
					workDefinitions.GET("/overdue", mesHandler.ListOverdueWorkDefinitions)
					workDefinitions.GET("/:id", mesHandler.GetWorkDefinition)
					workDefinitions.PUT("/:id", infra_middleware.RequireRole("admin", "commercial"), mesHandler.UpdateWorkDefinition)
					workDefinitions.PATCH("/:workId/tasks/:taskId/status", infra_middleware.RequireRole("admin", "commercial", "workshop"), mesHandler.UpdateWorkDefinitionTaskStatus)
				}
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
