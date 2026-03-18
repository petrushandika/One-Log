package database

import (
	"log"

	"github.com/petrushandika/one-log/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Seed populates the database with initial dummy data layout sets setups.
func Seed(db *gorm.DB) {
	log.Println("Database seeding started...")

	// 0. Create Default Admin
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	adminUser := domain.User{
		Email:    "admin@example.com",
		Password: string(hashedPassword),
		Name:     "System Admin",
	}

	var userCount int64
	db.Model(&domain.User{}).Where("email = ?", adminUser.Email).Count(&userCount)
	if userCount == 0 {
		if err := db.Create(&adminUser).Error; err != nil {
			log.Printf("Failed to seed admin: %v", err)
		} else {
			log.Println("Admin user seeded: admin@example.com / admin123")
		}
	} else {
		log.Println("Admin user already exists, skipping...")
	}

	// 1. Create Dummy Sources
	sources := []domain.Source{
		{
			Name:      "Authentication Service",
			APIKey:    "REDACTED_KEY",
			HealthURL: "https://auth.sample.com/health",
			Status:    "ONLINE",
		},
		{
			Name:      "API Gateway",
			APIKey:    "REDACTED_KEY",
			HealthURL: "https://gateway.sample.com/health",
			Status:    "ONLINE",
		},
		{
			Name:      "Database Analytics",
			APIKey:    "REDACTED_KEY",
			HealthURL: "https://db.sample.com/health",
			Status:    "OFFLINE",
		},
	}

	for _, src := range sources {
		// Check if source already exists to avoid duplicates node layout triggers setups configurations triggers sets
		var count int64
		db.Model(&domain.Source{}).Where("name = ?", src.Name).Count(&count)
		if count == 0 {
			if err := db.Create(&src).Error; err != nil {
				log.Printf("Failed to create source %s: %v", src.Name, err)
			} else {
				log.Printf("Created source %s with ID %s", src.Name, src.ID)

				// 2. Create Dummy Log Entries for this Source triggers sets setups
				dummyLogs := []domain.LogEntry{
					{
						SourceID:   src.ID,
						Category:   "SYSTEM_ERROR",
						Level:      "ERROR",
						Message:    "Database connection pool exhausted",
						StackTrace: "main.connectDB() at db.go:45\ngoroutine 1 [running]...",
					},
					{
						SourceID:   src.ID,
						Category:   "SECURITY",
						Level:      "CRITICAL",
						Message:    "Brute force attack detected from IP 192.168.1.1",
						StackTrace: "",
					},
					{
						SourceID:   src.ID,
						Category:   "USER_ACTIVITY",
						Level:      "INFO",
						Message:    "User (id: usr_123) logged in successfully",
						StackTrace: "",
					},
				}

				for _, l := range dummyLogs {
					if err := db.Create(&l).Error; err != nil {
						log.Printf("Failed to create log for source %s: %v", src.Name, err)
					}
				}
			}
		} else {
			log.Printf("Source %s already exists, skipping...", src.Name)
		}
	}

	log.Println("Database seeding completed successfully!")
}
