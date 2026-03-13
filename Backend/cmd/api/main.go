package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/petrushandika/one-log/internal/domain"
	"github.com/petrushandika/one-log/internal/middleware"
	"github.com/petrushandika/one-log/pkg/database"
	"github.com/petrushandika/one-log/pkg/utils"
)

func main() {
	// Flags
	migrateFlag := flag.Bool("migrate", false, "Run database migrations and exit")
	flag.Parse()

	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize Database
	db := database.InitDB()

	// Handle Migration Flag
	if *migrateFlag {
		database.Migrate(db)
		return
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	// Set Gin Mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Apply CORS Middleware (Allow All)
	r.Use(middleware.CORSMiddleware())

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		utils.Success(c, http.StatusOK, "System is healthy", gin.H{
			"app": "ULAM API",
		})
	})

	// Example Ingest with Structured Validation
	r.POST("/api/ingest", func(c *gin.Context) {
		var req domain.IngestLogRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			// Extract specific validation errors would be better in a separate helper,
			// for now we send the raw bind error formatted correctly.
			utils.Error(c, http.StatusUnprocessableEntity, "Validation failed", []utils.ErrorDetail{
				{Field: "request_body", Message: err.Error()},
			})
			return
		}

		// Logic to save to DB would go here
		utils.Success(c, http.StatusAccepted, "Log received successfully", gin.H{
			"log_id": 12345, // Placeholder
		})
	})

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
