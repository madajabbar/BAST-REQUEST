package config

import (
	"fmt"
	"log"

	"bast-request/internal/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// Using SQLite for local development (pure go, no CGO required)
	db, err := gorm.Open(sqlite.Open("bast_request.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
	}

	DB = db

	fmt.Println("Database connection successfully opened")
}

func AutoMigrate() {
	err := DB.AutoMigrate(
		&models.Customer{},
		&models.Project{},
		&models.BastFormat{},
		&models.BastSequence{},
		&models.BastRequest{},
		&models.AuditLog{},
		&models.User{},
		&models.Role{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database. \n", err)
	}

	fmt.Println("Database Migration Completed")
}
