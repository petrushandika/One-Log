package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/petrushandika/one-log/internal/handler"
	"github.com/petrushandika/one-log/internal/middleware"
	"github.com/petrushandika/one-log/internal/repository"
	"github.com/petrushandika/one-log/internal/service"
	"github.com/petrushandika/one-log/internal/worker"
	"github.com/petrushandika/one-log/pkg/database"
	"github.com/petrushandika/one-log/pkg/utils"
)

func main() {
	// 1. Setup Arguments (Check if -migrate flag is provided during run)
	migrateFlag := flag.Bool("migrate", false, "Run database migrations and exit")
	seedFlag := flag.Bool("seed", false, "Run database seeders and exit")
	flag.Parse()

	// 2. Load Environment file (.env)
	// Look for .env in the root Backend folder first, fallback to current dir
	if err := godotenv.Load("../../.env"); err != nil {
		if err := godotenv.Load(".env"); err != nil {
			log.Println("No .env file found, using system environment variables")
		}
	}

	// 3. Verify Required Secure Credentials
	if os.Getenv("ADMIN_EMAIL") == "" || os.Getenv("ADMIN_PASSWORD") == "" || os.Getenv("JWT_SECRET") == "" {
		log.Fatal("CRITICAL SECURITY ERROR: ADMIN_EMAIL, ADMIN_PASSWORD, or JWT_SECRET is missing from environment variables.")
	}

	// 4. Initialize Database
	db := database.InitDB()

	// Execute migration if the -migrate flag is set
	if *migrateFlag {
		database.Migrate(db)
		return
	}

	if *seedFlag {
		database.Seed(db)
		return
	}

	// 4. Dependency Injection (Wire all layers)
	logRepo := repository.NewLogRepository(db)

	notifySvc := service.NewNotificationService()
	aiSvc := service.NewAIService(logRepo)

	logService := service.NewLogService(logRepo, notifySvc, aiSvc)
	logHandler := handler.NewLogHandler(logService)
	activityService := service.NewActivityService(logRepo)
	activityHandler := handler.NewActivityHandler(activityService)
	apmService := service.NewAPMService(logRepo)
	apmHandler := handler.NewAPMHandler(apmService)
	issueService := service.NewIssueService(logRepo)
	issueHandler := handler.NewIssueHandler(issueService)

	sourceRepo := repository.NewSourceRepository(db)
	sourceService := service.NewSourceService(sourceRepo)
	sourceHandler := handler.NewSourceHandler(sourceService)

	authHandler := handler.NewAuthHandler(db, logService)
	statusHandler := handler.NewStatusHandler(sourceRepo)

	configRepo := repository.NewConfigRepository(db)
	configService := service.NewConfigService(configRepo)
	configHandler := handler.NewConfigHandler(configService, sourceService)

	// 5. Start Background Workers
	retentionWorker := worker.NewRetentionWorker(logRepo, 30) // 30 days retention
	retentionWorker.Start()

	uptimeWorker := worker.NewUptimeWorker(sourceRepo, logService)
	uptimeWorker.Start()

	// 5. Setup Router (Gin Framework)
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(middleware.CORSMiddleware()) // Apply CORS middleware

	// 6. Register Routes
	r.GET("/health", func(c *gin.Context) {
		utils.Success(c, http.StatusOK, "System is healthy", gin.H{"app": "ULAM API"})
	})

	api := r.Group("/api/v1")
	{
		// Public (no auth) status page data
		api.GET("/status", statusHandler.PublicStatus)

		// 6a. Public Endpoint
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/logout", authHandler.Logout)
		}

		// 6b. API Key Protected (For client applications)
		ingest := api.Group("/ingest")
		ingest.Use(middleware.APIKeyAuth(sourceRepo))
		{
			ingest.POST("", logHandler.Ingest)
		}

		// 6c. JWT Protected (For admins in the UI dashboard)
		admin := api.Group("")
		admin.Use(middleware.JWTAuth())
		{
			// Logs
			admin.GET("/logs", logHandler.GetAll)
			admin.GET("/logs/:id", logHandler.GetByID)
			admin.POST("/logs/:id/analyze", logHandler.Analyze)
			admin.GET("/logs/export", logHandler.ExportCSV)

			// Stats
			admin.GET("/stats/overview", logHandler.GetStatsOverview)
			admin.GET("/stats/activity", logHandler.GetActivitySummary)

			// Activity (Phase 2)
			admin.GET("/activity", activityHandler.List)
			admin.GET("/activity/summary", activityHandler.Summary)
			admin.GET("/activity/users/:user_id", activityHandler.ByUser)
			admin.GET("/activity/suspicious", activityHandler.Suspicious)

			// APM (Phase 3)
			admin.GET("/apm/endpoints", apmHandler.EndpointStats)

			// Issues (Phase 5)
			admin.GET("/issues", issueHandler.List)
			admin.GET("/issues/:fingerprint", issueHandler.Get)
			admin.PATCH("/issues/:fingerprint", issueHandler.UpdateStatus)
			admin.GET("/issues/:fingerprint/logs", issueHandler.Logs)

			// Configs
			admin.POST("/sources/:id/configs", configHandler.Save)
			admin.GET("/sources/:id/configs", configHandler.GetBySource)
			admin.GET("/sources/:id/configs/history", configHandler.History)

			// Sources
			admin.POST("/sources", sourceHandler.Create)
			admin.GET("/sources", sourceHandler.GetAll)
			admin.GET("/sources/:id", sourceHandler.GetByID)
			admin.PATCH("/sources/:id", sourceHandler.Update)
			admin.POST("/sources/:id/rotate-key", sourceHandler.RotateKey)
		}
	}

	// 7. Start the Server
	log.Printf("Server starting on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
