package config

import (
	"log"

	"bast-request/internal/models"
	"gorm.io/gorm"
)

func SeedDB(db *gorm.DB) {
	// Check and Seed Roles
	var roleCount int64
	db.Model(&models.Role{}).Count(&roleCount)
	if roleCount == 0 {
		log.Println("Seeding roles...")
		roles := []models.Role{
			{Name: "superadmin"},
			{Name: "admin"},
			{Name: "user"},
		}
		for _, role := range roles {
			if err := db.Create(&role).Error; err != nil {
				log.Printf("Failed to seed role %s: %v\n", role.Name, err)
			}
		}
	}

	// Check if data already exists
	var count int64
	db.Model(&models.Customer{}).Count(&count)
	if count > 0 {
		log.Println("Database already contains data, skipping seed.")
		return
	}

	log.Println("Seeding database with dummy data...")

	// Seed Customers
	customer1 := models.Customer{
		CustomerCode: "CUST-001",
		CustomerName: "PT. Maju Mundur",
		Status:       "active",
	}
	customer2 := models.Customer{
		CustomerCode: "CUST-002",
		CustomerName: "CV. Sukses Selalu",
		Status:       "active",
	}
	if err := db.Create(&customer1).Error; err != nil {
		log.Printf("Failed to seed customer1: %v\n", err)
	}
	if err := db.Create(&customer2).Error; err != nil {
		log.Printf("Failed to seed customer2: %v\n", err)
	}

	// Seed Projects
	project1 := models.Project{
		CustomerID:  customer1.CustomerID,
		ProjectCode: "PRJ-MM-01",
		ProjectName: "Implementasi Sistem ERP Terpadu",
		Status:      "active",
	}
	project2 := models.Project{
		CustomerID:  customer2.CustomerID,
		ProjectCode: "PRJ-SS-01",
		ProjectName: "Migrasi Infrastruktur Cloud",
		Status:      "active",
	}
	if err := db.Create(&project1).Error; err != nil {
		log.Printf("Failed to seed project1: %v\n", err)
	}
	if err := db.Create(&project2).Error; err != nil {
		log.Printf("Failed to seed project2: %v\n", err)
	}

	// Seed BAST Formats
	format1 := models.BastFormat{
		FormatName:    "Format PO Standar",
		FormatType:    "PO",
		FormatPattern: "BAST/PO/{YYYY}/{MM}/{SEQ}",
		IsActive:      true,
	}
	format2 := models.BastFormat{
		FormatName:    "Format Internal Perusahaan",
		FormatType:    "Internal",
		FormatPattern: "BAST/INT/{YYYY}/{MM}/{SEQ}",
		IsActive:      true,
	}
	if err := db.Create(&format1).Error; err != nil {
		log.Printf("Failed to seed format1: %v\n", err)
	}
	if err := db.Create(&format2).Error; err != nil {
		log.Printf("Failed to seed format2: %v\n", err)
	}

	log.Println("Database successfully seeded.")
}
