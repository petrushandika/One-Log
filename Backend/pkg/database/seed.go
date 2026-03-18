package database

import (
	"log"
	"os"

	"github.com/petrushandika/one-log/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Seed creates the initial admin user from environment variables.
// No sample sources or logs are created — register sources via the dashboard.
func Seed(db *gorm.DB) {
	log.Println("Database seeding started...")

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

	var existingUser domain.User
	if err := db.Where("email = ?", adminEmail).First(&existingUser).Error; err != nil {
		adminUser := domain.User{
			Email:    adminEmail,
			Password: string(hashedPassword),
			Name:     adminName,
		}
		if err := db.Create(&adminUser).Error; err != nil {
			log.Fatalf("Failed to seed admin user: %v", err)
		}
		log.Printf("Admin user created: %s", adminEmail)
	} else {
		log.Printf("Admin user already exists: %s, skipping.", adminEmail)
	}

	log.Println("Database seeding completed.")
}
