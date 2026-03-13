package database

import (
"log"
"os"

"github.com/petrushandika/one-log/internal/domain"
"gorm.io/driver/postgres"
"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return db
}

func Migrate(db *gorm.DB) {
	log.Println("Running database migrations...")
	err := db.AutoMigrate(
&domain.Source{},
		&domain.LogEntry{},
	)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Migration completed successfully")
}
