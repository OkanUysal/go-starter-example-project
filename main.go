package main

import (
	"github.com/OkanUysal/go-logger"
	"github.com/OkanUysal/go-metrics"
	"github.com/OkanUysal/go-starter-example-project/auth"
	"github.com/OkanUysal/go-starter-example-project/config"
	"github.com/OkanUysal/go-starter-example-project/handlers"
	"github.com/OkanUysal/go-starter-example-project/websocket"
	"github.com/OkanUysal/go-swagger"
	"github.com/gin-gonic/gin"

	docs "github.com/OkanUysal/go-starter-example-project/docs"

	_ "github.com/OkanUysal/go-starter-example-project/docs" // Import generated docs
)

// @title           Go Starter Example ProjectAPI
// @version         1.0.0
// @description     REST API for Go starter example project
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @schemes http https
func main() {
	// Load environment variables
	config.LoadEnv()

	// Initialize logger
	loggerConfig := &logger.Config{
		Level:      logger.LevelInfo,
		Format:     logger.FormatJSON,
		TimeFormat: "2006-01-02 15:04:05",
	}
	log := logger.New(loggerConfig)
	config.Logger = log // Set global logger for other packages

	// Connect to database
	if err := config.ConnectDatabase(); err != nil {
		log.Error("Failed to connect to database", logger.Err(err))
		return
	}
	log.Info("Database connected successfully")

	// Initialize cache
	if err := config.InitCache(); err != nil {
		log.Error("Failed to initialize cache", logger.Err(err))
		return
	}
	log.Info("Cache initialized successfully")

	// Initialize auth service
	if err := handlers.InitAuthService(); err != nil {
		log.Error("Failed to initialize auth service", logger.Err(err))
		return
	}
	log.Info("Auth service initialized successfully")

	// Initialize and start WebSocket room manager
	roomManager := websocket.GetRoomManager()
	roomManager.Start()
	log.Info("WebSocket room manager initialized")

	// Initialize metrics
	metricsConfig := &metrics.Config{
		ServiceName: config.GetEnv("SERVICE_NAME", "go-starter-example-project"),
	}
	metricsInstance := metrics.NewMetrics(metricsConfig)

	// Create Gin router
	r := gin.Default()

	// Serve static files
	r.Static("/static", "./static")

	// Setup metrics endpoints (/metrics and /health)
	metricsInstance.Setup(r)

	// Metrics middleware (automatic HTTP metrics collection)
	r.Use(metricsInstance.GinMiddleware())

	// Swagger documentation with auto host detection
	swagSpec, err := swagger.LoadSwagDocs(docs.SwaggerInfo.ReadDoc())
	if err != nil {
		logger.Error("Failed to load swagger docs", logger.Err(err))
	} else {
		swagger.SetupWithSwag(r, swagSpec, swagger.DefaultConfig())
		logger.Info("Swagger UI enabled", logger.String("path", "/swagger/index.html"))
	}

	// API routes group
	api := r.Group("/api")
	{
		api.GET("/hello", handlers.HelloHandler)

		// Auth routes
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/guest-login", handlers.GuestLogin)
			authGroup.POST("/refresh", handlers.RefreshToken)

			// Protected routes
			authGroup.GET("/me", auth.Middleware(), handlers.GetMe)
		}

		// Admin routes - requires authentication and admin role
		adminGroup := api.Group("/admin")
		adminGroup.Use(auth.Middleware())
		adminGroup.Use(auth.AdminMiddleware())
		{
			adminGroup.GET("/dashboard", handlers.AdminDashboard)
			adminGroup.GET("/users", handlers.ListUsers)
		}

		// WebSocket routes
		wsGroup := api.Group("/ws")
		wsGroup.Use(auth.Middleware()) // All WebSocket routes require authentication
		{
			// WebSocket connection endpoint
			wsGroup.GET("", websocket.WebSocketConnect)

			// Room management endpoints
			wsGroup.GET("/rooms", websocket.GetRooms)
			wsGroup.GET("/rooms/:room_id", websocket.GetRoomInfo)

			// Admin-only WebSocket endpoints
			wsAdminGroup := wsGroup.Group("")
			wsAdminGroup.Use(auth.AdminMiddleware())
			{
				wsAdminGroup.POST("/rooms", websocket.CreateRoom)
				wsAdminGroup.DELETE("/rooms/:room_id", websocket.CloseRoom)
				wsAdminGroup.POST("/invite", websocket.InviteToRoom)
			}
		}
	}

	// Start server
	log.Info("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Error("Failed to start server", logger.Err(err))
	}
}
