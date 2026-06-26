# Step 2: Models, Migrasi, & Seeding Database

Bagian ini membahas isi dari folder `internal/models/` dan bagaimana kode Go Anda secara ajaib bisa menjadi tabel-tabel di dalam database.

## 1. Definisi Models (Database Schema)
Model adalah struktur data (menggunakan `struct`) yang mewakili satu baris data di tabel SQL. GORM membaca tag-tag khusus di samping tipe datanya untuk merancang tabel.

Mari kita ambil satu contoh file, yaitu `internal/models/customer.go`:

```go
package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Customer struct {
	// gorm:"type:uuid;primary_key" memberitahu GORM bahwa kolom ini adalah Primary Key bertipe UUID
	CustomerID   uuid.UUID      `gorm:"type:uuid;primary_key"`
	CustomerCode string         `gorm:"type:varchar(50);uniqueIndex;not null"`
	CustomerName string         `gorm:"type:varchar(255);not null"`
	Status       string         `gorm:"type:varchar(50);default:'active'"`
	
	// Kolom waktu yang diisi otomatis oleh GORM
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"` // Kolom ini membuat fitur "Soft Delete" aktif
}

// Hook BeforeCreate (Trigger sebelum data ditambahkan ke DB)
func (c *Customer) BeforeCreate(tx *gorm.DB) (err error) {
	if c.CustomerID == uuid.Nil {
		c.CustomerID = uuid.New() // Menghasilkan string acak sepanjang 36 karakter untuk ID
	}
	return
}
```

> **Catatan Soft Delete:**
> Dengan adanya `DeletedAt gorm.DeletedAt`, jika Anda menghapus *customer*, GORM tidak akan benar-benar menghapusnya dari tabel SQL. GORM hanya mengisi tanggal hapus di kolom `DeletedAt`. *Customer* ini otomatis akan tersembunyi jika Anda mencoba me-Read nya dari GORM. Ini sangat aman untuk menghindari kehilangan data tak terduga!

Proyek ini memiliki banyak model lain seperti `Project`, `User`, `BastFormat`, `BastRequest`, yang strukturnya mirip-mirip.

## 2. Auto-Migration (`internal/config/database.go`)
GORM memiliki fitur luar biasa bernama `AutoMigrate`. Fitur ini mencocokkan struktur `struct` di atas dengan database secara langsung saat aplikasi dijalankan. Jika tabel belum ada, tabel akan dibuat. Jika ada kolom baru, kolom akan ditambahkan.

```go
func AutoMigrate() {
	err := DB.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.Customer{},
		&models.Project{},
		&models.BastFormat{},
		&models.BastSequence{},
		&models.BastRequest{},
		&models.AuditLog{},
	)
	if err != nil {
		log.Fatal("Gagal menjalankan migrasi:", err)
	}
	log.Println("Database Migration Completed")
}
```

## 3. Database Seeding (`internal/config/seed.go`)
*Seeding* berarti menanamkan data "bibit" awal ke dalam database yang masih kosong, agar kita bisa langsung menggunakan aplikasi tanpa pusing memasukkan data secara manual satu-persatu.

```go
func SeedDB(db *gorm.DB) {
	// Mengecek apakah tabel Role sudah ada isinya?
	var roleCount int64
	db.Model(&models.Role{}).Count(&roleCount)
	
	// Jika masih kosong (0), buatkan role-role dasar
	if roleCount == 0 {
		roles := []models.Role{
			{Name: "superadmin"},
			{Name: "admin"},
			{Name: "user"},
		}
		for _, role := range roles {
			db.Create(&role)
		}
	}
	
	// ... (Sama halnya dengan seeding untuk Customer, Project, dll)
}
```

## Ringkasan Step 2:
Berkat model dan auto-migrate, developer tidak perlu mengotak-atik perintah SQL seperti `CREATE TABLE`. Semuanya diatur penuh oleh kode Go di dalam *layer* `Models`.

Lanjut ke: **[Step 3: Authentication Flow](step_3_authentication_flow.md)**
