# Step 1: Setup Proyek, Konfigurasi & Database

> Seri Tutorial · **Step 1 dari 8**

Pada langkah pertama ini, kita membedah **pondasi aplikasi**: bagaimana proyek diinisialisasi, bagaimana koneksi database dibuka, dan apa yang terjadi saat Anda menjalankan `go run ./cmd/api/main.go`.

---

## 1. Inisialisasi Proyek Go

Semua proyek Go dimulai dengan satu perintah untuk membuat file `go.mod`:

```bash
go mod init bast-request
```

`go.mod` adalah "KTP" proyek: mencatat **nama modul** (`bast-request`) dan **daftar dependensi** eksternal (Gin, GORM, JWT, dll). Lihat [`go.mod`](../../go.mod) — nama modul `bast-request` inilah yang dipakai saat import antar package:

```go
import "bast-request/internal/config"   // modul/nama-folder
```

Setelah menambah dependensi, jalankan:
```bash
go mod tidy
```
untuk merapikan & mengunduh semua paket yang dibutuhkan.

---

## 2. Struktur Folder Clean Architecture

Aplikasi mengikuti pola **Clean Architecture** (detail: [di sini](../architecture/clean-architecture.md)). Struktur foldernya:

```text
BAST Request/
├── cmd/api/main.go              # Entry point
├── internal/
│   ├── config/                  # Koneksi DB, migrasi, seed
│   ├── models/                  # Struct tabel database
│   ├── repositories/            # Query SQL (GORM)
│   ├── services/                # Logika bisnis
│   ├── handlers/                # Penerima HTTP request
│   ├── middlewares/             # Filter keamanan (JWT, role)
│   ├── routes/                  # Pendaftaran endpoint
│   └── utils/                   # Fungsi bantu (hash, JWT)
└── docs/                        # Dokumentasi (file ini + Swagger)
```

> Kenapa folder utama bernama `internal/`? Di Go, folder `internal` punya aturan khusus: **hanya bisa di-import dari dalam modul yang sama**. Ini mencegah orang luar memakai kode internal Anda. Fitur keamanan bawaan Go.

---

## 3. Konfigurasi Database

File: [`internal/config/database.go`](../../internal/config/database.go)

```go
package config

import (
	"fmt"
	"log"

	"bast-request/internal/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB  // (1) variabel global

func ConnectDB() {
	// (2) Buka file SQLite lokal
	db, err := gorm.Open(sqlite.Open("bast_request.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
	}

	DB = db  // (3) simpan ke variabel global
	fmt.Println("Database connection successfully opened")
}
```

### Penjelasan baris penting
- **(1) `var DB *gorm.DB`** — variabel global yang menyimpan koneksi database. Huruf kapital `DB` artinya **exported** (bisa diakses dari package lain via `config.DB`).
- **(2) `sqlite.Open("bast_request.db")`** — membuka/membuat file database bernama `bast_request.db` di root proyek. Driver `glebarez/sqlite` dipilih karena **pure Go (tanpa CGO)** — tidak perlu compiler C.
- **(3) `DB = db`** — menyimpan koneksi agar bisa dipakai di seluruh aplikasi.
- **`log.Fatal(...)`** — jika koneksi gagal, aplikasi **langsung berhenti** (exit). Logis: tanpa DB, aplikasi tak berguna.

> 💡 **Catatan desain:** Pemakaian variabel global `DB` mempermudah inisialisasi awal. Namun di layer Repository, koneksi ini akan **disuntikkan** (dependency injection) secara eksplisit agar tetap testable.

---

## 4. AutoMigrate (Buat Tabel Otomatis)

Masih di [`internal/config/database.go:26-42`](../../internal/config/database.go):

```go
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
```

`AutoMigrate` menerima daftar **struct model**. Untuk tiap struct, GORM:
1. Mengecek apakah tabel sudah ada di DB.
2. Jika **belum** → buat tabel sesuai definisi struct.
3. Jika **sudah ada** → tambahkan kolom baru yang belum ada (tapi **tidak menghapus** kolom lama).

Hasilnya: Anda tidak perlu menulis SQL `CREATE TABLE` manual. Detail model akan dibahas di [Step 2](step-02-models-and-migration.md).

---

## 5. Fungsi Utama `main()`

File: [`cmd/api/main.go`](../../cmd/api/main.go)

```go
package main

import (
	"bast-request/internal/config"
	"bast-request/internal/routes"

	"github.com/gin-gonic/gin"
)

// @title BAST Request API
// @version 1.0
// @description This is the API server for BAST Request System.
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// (1) Buka koneksi database
	config.ConnectDB()

	// (2) Buat tabel jika belum ada
	config.AutoMigrate()

	// (3) Isi data awal (roles, customer contoh, dll)
	config.SeedDB(config.DB)

	// (4) Siapkan server web Gin
	r := gin.Default()

	// (5) Endpoint test publik
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// (6) Daftarkan semua URL endpoint
	routes.SetupRoutes(r, config.DB)

	// (7) Nyalakan server di port 8080
	r.Run(":8080")
}
```

### Urutan eksekusi
1. **ConnectDB** → buka `bast_request.db`.
2. **AutoMigrate** → buat tabel-tabel.
3. **SeedDB** → isi data contoh.
4. **gin.Default()** → instance server Gin (sudah termasuk logger & recovery middleware bawaan).
5. **`/ping`** → endpoint test sederhana tanpa autentikasi.
6. **SetupRoutes** → daftarkan endpoint API (dibahas [Step 7](step-07-routing-and-rbac.md)).
7. **`r.Run(":8080")`** → server mendengarkan request. Baris ini **memblokir** — aplikasi terus berjalan sampai dihentikan (Ctrl+C).

### Komentar aneh di atas `main()`?
Baris-baris `// @title`, `@version`, dst itu **bukan komentar biasa** — itu **anotasi Swagger**. Akan dibahas di [Step 8](step-08-swagger-documentation.md).

---

## 6. Uji Coba

Saat server jalan, coba endpoint `/ping`:
```bash
curl http://localhost:8080/ping
```
Respons:
```json
{ "message": "pong" }
```

Jika ini muncul, berarti **seluruh pondasi berfungsi**: Go terkompilasi, dependensi terpasang, DB terhubung, Gin merespons HTTP.

---

## ✅ Ringkasan Step 1
- `go mod init` + `go mod tidy` menyiapkan modul & dependensi.
- `internal/config/database.go` membuka koneksi SQLite & menjalankan AutoMigrate.
- `cmd/api/main.go` adalah urutan boot aplikasi: DB → Migrate → Seed → Gin → Routes → Run.
- Server mendengarkan di port 8080, dengan endpoint `/ping` sebagai tes awal.

Aplikasi sudah berdiri, tapi tabel-tabelnya masih kosong definisinya. Mari kita lihat bagaimana model Go menjadi tabel.

---

⬅️ **[Daftar Tutorial](README.md)** · ➡️ **[Step 2: Models & Migrasi](step-02-models-and-migration.md)**
