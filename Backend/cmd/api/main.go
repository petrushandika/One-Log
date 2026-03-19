package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/petrushandika/one-log/internal/handler"
	"github.com/petrushandika/one-log/internal/middleware"
	"github.com/petrushandika/one-log/internal/repository"
	"github.com/petrushandika/one-log/internal/service"
	"github.com/petrushandika/one-log/internal/worker"
	"github.com/petrushandika/one-log/pkg/ai"
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

	chatService := service.NewChatService(logRepo, ai.NewGroqClient())
	chatHandler := handler.NewChatHandler(chatService)

	// Phase 4: Incident Management
	incidentRepo := repository.NewIncidentRepository(db)
	incidentService := service.NewIncidentService(incidentRepo)
	incidentHandler := handler.NewIncidentHandler(incidentService)

	// Phase 2: Activity Analytics
	activityAnalyticsRepo := repository.NewActivityAnalyticsRepository(db)
	activityAnalyticsService := service.NewActivityAnalyticsService(activityAnalyticsRepo)
	activityAnalyticsHandler := handler.NewActivityAnalyticsHandler(activityAnalyticsService)

	// Phase 2 Extended: Activity Monitor (Feed, Top Users, Compliance Export)
	activityMonitorRepo := repository.NewActivityMonitorRepository(db)
	activityMonitorService := service.NewActivityMonitorService(activityMonitorRepo)
	activityMonitorHandler := handler.NewActivityMonitorHandler(activityMonitorService)

	// Phase 3: APM Thresholds
	apmThresholdRepo := repository.NewAPMThresholdRepository(db)
	apmThresholdService := service.NewAPMThresholdService(apmThresholdRepo, logRepo)
	apmThresholdHandler := handler.NewAPMThresholdHandler(apmThresholdService)

	// Phase 4: Status Page
	statusPageRepo := repository.NewStatusPageRepository(db)
	statusPageService := service.NewStatusPageService(statusPageRepo)
	statusPageHandler := handler.NewStatusPageHandler(statusPageService)

	// Phase 7: Export Handler
	exportHandler := handler.NewExportHandler(logRepo)

	// 5. Start Background Workers
	retentionWorker := worker.NewRetentionWorker(logRepo, 30) // 30 days retention
	retentionWorker.Start()

	uptimeWorker := worker.NewUptimeWorker(sourceRepo, incidentRepo, logService, notifySvc)
	uptimeWorker.Start()

	// Phase 5: Regression Detection Worker
	regressionDetector := worker.NewRegressionDetector(logRepo, notifySvc)
	regressionDetector.Start()

	// Phase 5: Daily Digest Worker
	aiClient := ai.NewGroqClient()
	dailyDigestWorker := worker.NewDailyDigestWorker(logRepo, aiClient, notifySvc)
	dailyDigestWorker.Start()

	// 5. Setup Router (Gin Framework)
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(middleware.RequestIDMiddleware()) // Request tracing
	r.Use(middleware.CORSMiddleware())      // Apply CORS middleware

	// 6. Register Routes
	r.GET("/health", func(c *gin.Context) {
		utils.Success(c, http.StatusOK, "System is healthy", gin.H{"app": "ULAM API"})
	})

	// Public Status Page routes (outside /api group)
	r.GET("/status/:slug", statusPageHandler.PublicStatusPage)
	r.GET("/embed/:token", statusPageHandler.EmbedWidget)

	api := r.Group("/api")
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

		// 6b. API Key Protected (For client applications) - Rate limited: 100 req/min per API key
		ingest := api.Group("/ingest")
		ingest.Use(middleware.RateLimitByAPIKey(100, time.Minute))
		ingest.Use(middleware.APIKeyAuth(sourceRepo))
		{
			ingest.POST("", logHandler.Ingest)
		}

		// 6c. JWT Protected (For admins in the UI dashboard) - Rate limited: 60 req/min per user
		admin := api.Group("")
		admin.Use(middleware.RateLimitByJWT(60, time.Minute))
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

			// Activity Analytics (Phase 2 Extended)
			admin.GET("/activity/analytics/methods", activityAnalyticsHandler.GetAuthMethodBreakdown)
			admin.GET("/activity/analytics/timeline", activityAnalyticsHandler.GetLoginTimeline)
			admin.GET("/activity/analytics/heatmap", activityAnalyticsHandler.GetFailedLoginHeatmap)
			admin.GET("/activity/sessions", activityAnalyticsHandler.GetRecentSessions)

			// Activity Monitor (Phase 2 Extended)
			admin.GET("/activity/feed", activityMonitorHandler.GetActivityFeed)
			admin.GET("/activity/top-users", activityMonitorHandler.GetTopActiveUsers)
			admin.GET("/activity/by-resource", activityMonitorHandler.GetActivityByResource)
			admin.GET("/activity/users/:user_id/profile", activityMonitorHandler.GetUserProfile)
			admin.POST("/activity/compliance-export", activityMonitorHandler.RequestComplianceExport)
			admin.GET("/activity/compliance-exports", activityMonitorHandler.GetComplianceExports)

			// APM (Phase 3)
			admin.GET("/apm/endpoints", apmHandler.EndpointStats)
			admin.GET("/apm/timeline", apmHandler.ResponseTimeTimeline)

			// APM Thresholds (Phase 3 Extended)
			admin.GET("/apm/thresholds", apmThresholdHandler.List)
			admin.POST("/apm/thresholds", apmThresholdHandler.Create)
			admin.GET("/apm/thresholds/:id", apmThresholdHandler.Get)
			admin.PATCH("/apm/thresholds/:id", apmThresholdHandler.Update)
			admin.DELETE("/apm/thresholds/:id", apmThresholdHandler.Delete)
			admin.GET("/apm/slow-queries", apmThresholdHandler.GetSlowQueries)
			admin.GET("/apm/slow-queries/trend", apmThresholdHandler.GetSlowQueryTrend)
			admin.GET("/apm/apdex", apmThresholdHandler.GetApdexScore)
			admin.GET("/apm/threshold-alerts", apmThresholdHandler.GetThresholdAlerts)

			// Issues (Phase 5)
			admin.GET("/issues", issueHandler.List)
			admin.GET("/issues/analytics/trend", issueHandler.ErrorRateTrend)
			admin.GET("/issues/analytics/heatmap", issueHandler.ErrorHeatmap)
			admin.GET("/issues/:fingerprint", issueHandler.Get)
			admin.PATCH("/issues/:fingerprint", issueHandler.UpdateStatus)
			admin.GET("/issues/:fingerprint/logs", issueHandler.Logs)

			// Phase 4: Incident Management
			admin.GET("/incidents", incidentHandler.List)
			admin.GET("/incidents/timeline", incidentHandler.GetTimeline)

			// Phase 4: Status Page Admin
			admin.GET("/admin/status-pages", statusPageHandler.List)
			admin.POST("/admin/status-pages", statusPageHandler.Create)
			admin.GET("/admin/status-pages/:source_id", statusPageHandler.Get)
			admin.PATCH("/admin/status-pages/:source_id", statusPageHandler.Update)
			admin.DELETE("/admin/status-pages/:source_id", statusPageHandler.Delete)
			admin.GET("/admin/status-pages/:source_id/uptime", statusPageHandler.GetUptime)
			admin.POST("/admin/status-pages/:source_id/embed", statusPageHandler.CreateEmbed)

			// Phase 7: Export
			admin.GET("/logs/export/excel", exportHandler.ExportLogsExcel)
			admin.GET("/logs/export/pdf", exportHandler.ExportAuditPDF)

			// Configs
			admin.POST("/sources/:id/configs", configHandler.Save)
			admin.GET("/sources/:id/configs", configHandler.GetBySource)
			admin.GET("/sources/:id/configs/history", configHandler.History)

			// Sources
			admin.POST("/sources", sourceHandler.Create)
			admin.GET("/sources", sourceHandler.GetAll)
			admin.GET("/sources/:id", sourceHandler.GetByID)
			admin.PATCH("/sources/:id", sourceHandler.Update)
			admin.DELETE("/sources/:id", sourceHandler.Delete)
			admin.POST("/sources/:id/rotate-key", sourceHandler.RotateKey)

			// AI Chat Copilot (Phase 5)
			admin.POST("/chat", chatHandler.Ask)
		}
	}

	// 7. Start the Server
	log.Printf("Server starting on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
