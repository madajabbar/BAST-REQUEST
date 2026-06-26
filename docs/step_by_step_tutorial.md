# Tutorial Langkah-demi-Langkah (Step-by-Step): Menambahkan Fitur Baru

Tutorial ini akan memandu Anda secara praktis bagaimana cara menambahkan satu modul/tabel baru (misalnya tabel `Division` / Divisi) ke dalam proyek **BAST Request API** dengan menerapkan pola **Clean Architecture** yang sama.

Kita akan berjalan dari level terendah (*Database/Model*) hingga ke level tertinggi (*Router/API*).

---

## Langkah 1: Membuat Model (Database Schema)
Pertama, kita harus mendefinisikan bentuk tabel `Division` di dalam folder `models`.

1. Buat file baru di `internal/models/division.go`.
2. Ketikkan kode berikut:
```go
package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Division struct {
	DivisionID   uuid.UUID      `gorm:"type:uuid;primary_key"`
	DivisionName string         `gorm:"type:varchar(100);not null"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// Hook sebelum Create untuk mengisi UUID otomatis
func (d *Division) BeforeCreate(tx *gorm.DB) (err error) {
	if d.DivisionID == uuid.Nil {
		d.DivisionID = uuid.New()
	}
	return
}
```

## Langkah 2: Menambahkan ke AutoMigrate
Agar GORM secara otomatis membuatkan tabel `divisions` di database SQLite Anda:
1. Buka file `internal/config/database.go`.
2. Cari fungsi `AutoMigrate()`.
3. Tambahkan `&models.Division{}` ke dalam daftarnya:
```go
func AutoMigrate() {
	DB.AutoMigrate(
		&models.Role{},
		&models.User{},
		// ... tabel lain ...
		&models.Division{}, // <--- Tambahkan baris ini
	)
}
```

## Langkah 3: Membuat Repository (Akses Database)
Repository adalah satu-satunya tempat yang diizinkan untuk melakukan Query SQL.
1. Buat file `internal/repositories/division_repository.go`.
2. Ketikkan fungsi untuk Create dan Read data:
```go
package repositories

import (
	"bast-request/internal/models"
	"gorm.io/gorm"
)

type DivisionRepository struct {
	db *gorm.DB
}

func NewDivisionRepository(db *gorm.DB) *DivisionRepository {
	return &DivisionRepository{db: db}
}

func (r *DivisionRepository) GetAll() ([]models.Division, error) {
	var divisions []models.Division
	err := r.db.Find(&divisions).Error
	return divisions, err
}

func (r *DivisionRepository) Create(division *models.Division) error {
	return r.db.Create(division).Error
}
```

## Langkah 4: Membuat Service (Logika Bisnis)
Service akan memanggil Repository. Jika Anda butuh validasi (misalnya: Nama Divisi tidak boleh kosong), letakkan di sini.
1. Buat file `internal/services/division_service.go`.
```go
package services

import (
	"errors"
	"bast-request/internal/models"
	"bast-request/internal/repositories"
)

type DivisionService struct {
	repo *repositories.DivisionRepository
}

func NewDivisionService(repo *repositories.DivisionRepository) *DivisionService {
	return &DivisionService{repo: repo}
}

func (s *DivisionService) GetAllDivisions() ([]models.Division, error) {
	return s.repo.GetAll()
}

func (s *DivisionService) CreateDivision(div *models.Division) error {
	if div.DivisionName == "" {
		return errors.New("nama divisi tidak boleh kosong") // Validasi Bisnis
	}
	return s.repo.Create(div)
}
```

## Langkah 5: Membuat Handler (Lalu Lintas HTTP & Gin)
Handler bertugas menerima *Request* JSON dari *User*, dan memanggil *Service*.
1. Buat file `internal/handlers/division_handler.go`.
2. Jangan lupa tambahkan komentar Swagger!
```go
package handlers

import (
	"net/http"
	"bast-request/internal/models"
	"bast-request/internal/services"
	"github.com/gin-gonic/gin"
)

type DivisionHandler struct {
	service *services.DivisionService
}

func NewDivisionHandler(service *services.DivisionService) *DivisionHandler {
	return &DivisionHandler{service: service}
}

// GetAllDivisions godoc
// @Summary Get all divisions
// @Description Retrieve a list of all divisions
// @Tags divisions
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Division
// @Router /divisions [get]
func (h *DivisionHandler) GetAllDivisions(c *gin.Context) {
	divisions, err := h.service.GetAllDivisions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, divisions)
}

// CreateDivision godoc
// @Summary Create a division
// @Description Add a new division
// @Tags divisions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param division body models.Division true "Division Data"
// @Success 201 {object} models.Division
// @Router /divisions [post]
func (h *DivisionHandler) CreateDivision(c *gin.Context) {
	var input models.Division
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateDivision(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, input)
}
```

## Langkah 6: Mendaftarkan Route
Semua modul yang sudah dibuat harus disambungkan ke dalam *router* utama.
1. Buka `internal/routes/routes.go`.
2. Di dalam fungsi `SetupRoutes`, inisialisasi ketiga *layer* tadi:
```go
	// Inisialisasi Dependencies (Tambahkan di bagian atas bersama yang lain)
	divisionRepo := repositories.NewDivisionRepository(db)
	divisionService := services.NewDivisionService(divisionRepo)
	divisionHandler := handlers.NewDivisionHandler(divisionService)
```
3. Gulir ke bawah ke blok `protected := api.Group("/")` dan tambahkan rute URL-nya:
```go
			// MASTER DIVISION
			protected.GET("/divisions", divisionHandler.GetAllDivisions)
			protected.POST("/divisions", divisionHandler.CreateDivision)
```

## Langkah 7: Build Swagger dan Jalankan!
Terakhir, Anda harus memberi tahu Swagger bahwa ada rute baru.
1. Buka terminal Anda.
2. Ketik perintah:
```bash
swag init -g ./cmd/api/main.go --parseDependency --parseInternal
```
3. Restart server Anda:
```bash
go run ./cmd/api/main.go
```
4. Buka Postman / Swagger UI di `http://localhost:8080/swagger/index.html`. Anda akan melihat *endpoint* `/divisions` sudah muncul dan siap dites (jangan lupa masukkan Bearer Token Anda!).

---

🎉 **Selamat!** Anda telah berhasil mengimplementasikan *Clean Architecture* dari Model (Database) hingga ke Endpoint API secara utuh! Anda dapat menerapkan 7 langkah ini berulang kali untuk tabel dan fitur apa pun di masa depan.
