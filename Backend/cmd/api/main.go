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

	r := gin.Default()

	// Apply CORS Middleware (Allow All)
	r.Use(middleware.CORSMiddleware())

	// Health Check
	r.GET("/health", func(c *gin.Context) {
c.JSON(http.StatusOK, gin.H{
"status": "healthy",
"app":    "ULAM API",
})
})

	// Example Ingest with Validation
	r.POST("/api/ingest", func(c *gin.Context) {
var req domain.IngestLogRequest
if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
"error":   "Validation failed",
"details": err.Error(),
			})
			return
		}

		// Logic to save to DB would go here
		c.JSON(http.StatusAccepted, gin.H{
"status":  "success",
"message": "Log ingested successfully",
})
	})

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
