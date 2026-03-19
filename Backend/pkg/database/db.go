package database

import (
	"fmt"
	"log"
	"os"

	"github.com/petrushandika/one-log/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB initialization databse connection
func InitDB() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Fallback for local development
		host := os.Getenv("DB_HOST")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		port := os.Getenv("DB_PORT")
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", host, user, password, dbname, port)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed connect to database: %v", err)
	}

	log.Println("Database connection established")
	return db
}

// Migrate do syncronize database schema
func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&domain.Source{},
		&domain.LogEntry{},
		&domain.Issue{},
		&domain.SourceConfig{},
		&domain.SourceConfigHistory{},
		&domain.User{},
		&domain.Incident{},
		&domain.APMThreshold{},
		&domain.Session{},
	)
	if err != nil {
		log.Fatalf("Migration failed %v", err)
	}

	log.Println("Database migration completed")
}
