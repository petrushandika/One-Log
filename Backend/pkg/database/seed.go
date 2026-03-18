package database

import (
	"log"
	"os"

	"github.com/petrushandika/one-log/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Seed populates the database with initial data for a functional first run.
func Seed(db *gorm.DB) {
	log.Println("Database seeding started...")

	// 0. Create Default Admin — credentials come from env vars only
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	adminName := os.Getenv("ADMIN_NAME")

	if adminEmail == "" {
		adminEmail = "admin@onelog.com"
	}
	if adminPassword == "" {
		adminPassword = "123456"
	}
	if adminName == "" {
		adminName = "System Admin"
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash admin password: %v", err)
	}

	adminUser := domain.User{
		Email:    adminEmail,
		Password: string(hashedPassword),
		Name:     adminName,
	}

	var existingUser domain.User
	result := db.Where("email = ?", adminEmail).First(&existingUser)
	if result.Error != nil {
		// User not found — create it
		if err := db.Create(&adminUser).Error; err != nil {
			log.Printf("Failed to seed admin user: %v", err)
			return
		}
		log.Printf("Admin user created: %s", adminEmail)
		existingUser = adminUser
	} else {
		log.Printf("Admin user already exists: %s, skipping creation.", adminEmail)
	}

	// 1. Create Dummy Sources (associated with the admin user)
	sources := []domain.Source{
		{
			UserID:    existingUser.ID,
			Name:      "Authentication Service",
			APIKey:    "dev_auth_svc_001",
			HealthURL: "https://auth.sample.com/health",
			Status:    "ONLINE",
		},
		{
			UserID:    existingUser.ID,
			Name:      "API Gateway",
			APIKey:    "dev_api_gw_001",
			HealthURL: "https://gateway.sample.com/health",
			Status:    "ONLINE",
		},
		{
			UserID:    existingUser.ID,
			Name:      "Database Analytics",
			APIKey:    "dev_db_analytics_001",
			HealthURL: "https://db.sample.com/health",
			Status:    "OFFLINE",
		},
	}

	for _, src := range sources {
		var count int64
		db.Model(&domain.Source{}).Where("name = ?", src.Name).Count(&count)
		if count == 0 {
			if err := db.Create(&src).Error; err != nil {
				log.Printf("Failed to create source %s: %v", src.Name, err)
				continue
			}
			log.Printf("Created source: %s (ID: %s)", src.Name, src.ID)

			// 2. Create Dummy Log Entries for this source
			dummyLogs := []domain.LogEntry{
				{
					SourceID:   src.ID,
					Category:   "SYSTEM_ERROR",
					Level:      "ERROR",
					Message:    "Database connection pool exhausted",
					StackTrace: "main.connectDB() at db.go:45\ngoroutine 1 [running]...",
				},
				{
					SourceID: src.ID,
					Category: "SECURITY",
					Level:    "CRITICAL",
					Message:  "Brute force attack detected from IP 192.168.1.1",
				},
				{
					SourceID: src.ID,
					Category: "USER_ACTIVITY",
					Level:    "INFO",
					Message:  "User (id: usr_123) logged in successfully",
				},
			}

			for _, l := range dummyLogs {
				if err := db.Create(&l).Error; err != nil {
					log.Printf("Failed to create log for source %s: %v", src.Name, err)
				}
			}
		} else {
			log.Printf("Source already exists: %s, skipping.", src.Name)
		}
	}

	log.Println("Database seeding completed successfully!")
}
