# Step 1: Setup Proyek, Konfigurasi, & Database

Dokumen ini adalah bagian pertama dari seri bedah kode **BAST Request API**. Di sini kita akan mempelajari bagaimana pondasi aplikasi ini dibangun dari nol.

## 1. Inisialisasi Proyek Go
Proyek ini dimulai dengan membuat folder dan menjalankan perintah:
```bash
go mod init bast-request
```
Ini menghasilkan file `go.mod` yang bertugas mencatat semua *library* eksternal yang kita gunakan (seperti Gin, GORM, JWT, dll).

## 2. Struktur Folder
Aplikasi ini mengikuti pola **Clean Architecture** (atau *Layered Architecture*), sehingga struktur foldernya dibagi rapi:
*   `cmd/api/main.go` -> Pintu masuk utama aplikasi.
*   `internal/config/` -> Konfigurasi dan pengaturan (seperti koneksi database).
*   `internal/models/` -> Struktur tabel Database.
*   `internal/repositories/` -> Fungsi-fungsi yang melakukan *query* langsung ke Database.
*   `internal/services/` -> Logika Bisnis (*business rules*).
*   `internal/handlers/` -> Penerima *Request* HTTP (Controller).
*   `internal/routes/` -> Pengatur rute URL.
*   `internal/middlewares/` -> Filter (seperti pengecekan keamanan) di tengah jalan.
*   `internal/utils/` -> Fungsi pembantu (seperti *hashing password*).

## 3. Konfigurasi Database (`internal/config/database.go`)
Kita menggunakan SQLite sebagai database lokal yang disetir oleh ORM bernama **GORM**.

```go
package config

import (
	"log"
	"github.com/glebarez/sqlite" // Driver SQLite murni (tanpa butuh CGO)
	"gorm.io/gorm"
)

var DB *gorm.DB // Variabel global agar bisa dipanggil dari luar (dimulai dengan huruf Kapital)

func ConnectDB() {
	var err error
	// Membuka file database bernama "bast_request.db"
	DB, err = gorm.Open(sqlite.Open("bast_request.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal konek ke database:", err)
	}
	log.Println("Database connection successfully opened")
}
```
**Catatan Penting:** Penggunaan variabel global `DB` di sini mempermudah inisialisasi awal, namun nantinya pada layer *Repository*, koneksi ini akan disuntikkan (*dependency injection*) satu persatu.

## 4. Fungsi Utama (`cmd/api/main.go`)
Ini adalah titik mula (titik *start*) saat Anda menjalankan `go run main.go`.

```go
package main

import (
	"bast-request/internal/config"
	"bast-request/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Buka Koneksi ke Database
	config.ConnectDB()
	
	// 2. Lakukan Migrasi Tabel (Membuat tabel secara otomatis jika belum ada)
	config.AutoMigrate()

	// 3. Masukkan data-data awal (Seeds) seperti Role Admin
	config.SeedDB(config.DB)

	// 4. Siapkan server Web / API menggunakan Gin
	r := gin.Default()

	// 5. Daftarkan semua URL (Endpoints) aplikasi ke Gin
	routes.SetupRoutes(r, config.DB)

	// 6. Jalankan server di port 8080 (menunggu request dari luar)
	r.Run(":8080")
}
```

## Ringkasan Step 1:
Pada tahap ini, aplikasi sudah bisa berdiri. Koneksi ke database SQLite sudah terjalin, dan *web server* Gin sudah berjalan di *port* 8080. Langkah selanjutnya adalah menentukan apa saja isi tabel-tabel di database tersebut. 

Lanjut ke: **[Step 2: Models & Migration](step_2_models_and_migration.md)**
