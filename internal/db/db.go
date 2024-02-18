package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"todo-list-service/internal/models"
)

//	func ConnectDB() *gorm.DB {
//		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
//			os.Getenv("DB_HOST"),
//			os.Getenv("DB_PORT"),
//			os.Getenv("DB_USER"),
//			os.Getenv("DB_PASSWORD"),
//			os.Getenv("DB_NAME"),
//			os.Getenv("DB_SSL_MODE"))
//
//		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
//		if err != nil {
//			log.Fatalf("Error connecting to database: %v", err)
//		}
//
//		if err := db.AutoMigrate(&models.User{}, &models.Todo{}); err != nil {
//			log.Fatalf("Error auto-migrating: %v", err)
//		}
//
//		return db
//	}

func ConnectDB() *gorm.DB {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Todo{}); err != nil {
		log.Fatalf("Error auto-migrating: %v", err)
	}

	return db
}
