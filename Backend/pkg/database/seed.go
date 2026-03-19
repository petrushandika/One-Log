package database

import (
	"log"
	"os"
	"time"

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

		// Create sample activity feed data for testing
		seedActivityFeed(db, adminUser.ID)
	} else {
		log.Printf("Admin user already exists: %s, skipping.", adminEmail)
	}

	log.Println("Database seeding completed.")
}

// seedActivityFeed creates sample activity feed entries for testing
func seedActivityFeed(db *gorm.DB, userID uint) {
	log.Println("Seeding sample activity feed data...")

	// Check if activity_feeds table exists and has data
	var count int64
	db.Model(&domain.ActivityFeed{}).Count(&count)
	if count > 0 {
		log.Println("Activity feed already has data, skipping.")
		return
	}

	// Create sample activity feed entries
	activities := []domain.ActivityFeed{
		{
			UserID:       "admin",
			SourceID:     "system",
			Action:       "login",
			ResourceType: "auth",
			ResourceID:   "",
			Context:      map[string]interface{}{"method": "password", "success": true},
			IPAddress:    "127.0.0.1",
			CreatedAt:    time.Now().Add(-1 * time.Hour),
		},
		{
			UserID:       "admin",
			SourceID:     "system",
			Action:       "create",
			ResourceType: "config",
			ResourceID:   "app-settings",
			Context:      map[string]interface{}{"key": "debug_mode", "value": "false"},
			IPAddress:    "127.0.0.1",
			CreatedAt:    time.Now().Add(-30 * time.Minute),
		},
		{
			UserID:       "admin",
			SourceID:     "system",
			Action:       "view",
			ResourceType: "logs",
			ResourceID:   "",
			Context:      map[string]interface{}{"page": "overview", "filters": "none"},
			IPAddress:    "127.0.0.1",
			CreatedAt:    time.Now().Add(-15 * time.Minute),
		},
	}

	for _, activity := range activities {
		if err := db.Create(&activity).Error; err != nil {
			log.Printf("Failed to seed activity feed: %v", err)
		}
	}

	log.Printf("Created %d sample activity feed entries", len(activities))

	// Seed sample sessions
	seedSessions(db)
}

// seedSessions creates sample session data for testing
func seedSessions(db *gorm.DB) {
	log.Println("Seeding sample sessions data...")

	// Check if sessions table exists
	var tableExists bool
	err := db.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'sessions')").Scan(&tableExists).Error
	if err != nil || !tableExists {
		log.Println("Sessions table does not exist, skipping.")
		return
	}

	// Check if already has data
	var count int64
	db.Model(&domain.Session{}).Count(&count)
	if count > 0 {
		log.Println("Sessions already has data, skipping.")
		return
	}

	// Create sample sessions
	sessions := []domain.Session{
		{
			UserID:       "admin",
			SourceID:     "system",
			AuthMethod:   "password",
			IPAddress:    "127.0.0.1",
			Browser:      "Chrome 120.0",
			Device:       "macOS Desktop",
			IsActive:     true,
			LastActivity: time.Now().Add(-5 * time.Minute),
			CreatedAt:    time.Now().Add(-2 * time.Hour),
		},
		{
			UserID:       "admin",
			SourceID:     "system",
			AuthMethod:   "password",
			IPAddress:    "192.168.1.100",
			Browser:      "Firefox 121.0",
			Device:       "Windows Desktop",
			IsActive:     true,
			LastActivity: time.Now().Add(-30 * time.Minute),
			CreatedAt:    time.Now().Add(-5 * time.Hour),
		},
		{
			UserID:       "admin",
			SourceID:     "mobile-app",
			AuthMethod:   "oauth",
			IPAddress:    "10.0.0.50",
			Browser:      "Safari Mobile",
			Device:       "iPhone 14 Pro",
			IsActive:     true,
			LastActivity: time.Now().Add(-10 * time.Minute),
			CreatedAt:    time.Now().Add(-1 * time.Hour),
		},
	}

	for _, session := range sessions {
		if err := db.Create(&session).Error; err != nil {
			log.Printf("Failed to seed session: %v", err)
		}
	}

	log.Printf("Created %d sample sessions", len(sessions))
}
