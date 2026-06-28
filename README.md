# BAST Request API

Sistem API berbasis **Golang** untuk mengelola pembuatan dan permintaan **Berita Acara Serah Terima (BAST)** secara otomatis. Aplikasi ini dikembangkan dengan **Clean Architecture**, menjamin pemisahan tugas yang jelas (*Separation of Concerns*), skalabilitas, dan kemudahan pengujian.

## 🌟 Fitur Utama

- **Penomoran Otomatis yang Aman (*Atomic Sequence*):** Pembuatan nomor BAST secara berurutan dan kebal dari bentrok (Race Condition) berkat *Database Transactions*.
- **Role-Based Access Control (RBAC):** Sistem *Login* & Registrasi menggunakan **JWT (JSON Web Tokens)** dengan pengamanan hak akses level pengguna (*Superadmin, Admin, User*).
- **Master Data Terpusat:** Manajemen CRUD untuk `Customer`, `Project`, dan `Format Penomoran BAST`.
- **Audit Trail Log:** Semua aktivitas dan perubahan krusial di dalam sistem tercatat secara rapi dan otomatis ke dalam log database yang siap ditinjau kapanpun.
- **Auto-Generated Documentation:** Terintegrasi penuh dengan **Swagger UI** untuk memudahkan *frontend developer* membaca skema API.

---

## 🛠️ Teknologi yang Digunakan

- **Bahasa Pemrograman:** [Go (Golang)](https://go.dev/)
- **Web Framework:** [Gin Gonic](https://gin-gonic.com/)
- **ORM & Database:** [GORM](https://gorm.io/) dengan Driver [SQLite murni](https://github.com/glebarez/sqlite) (tanpa CGO).
- **Security:** JWT (`golang-jwt/jwt`) dan Bcrypt Password Hashing (`golang.org/x/crypto/bcrypt`).
- **Dokumentasi API:** [Swaggo](https://github.com/swaggo/swag)

---

## 🏗️ Struktur Arsitektur (Clean Architecture)

Proyek ini disusun dengan batas *layer* yang sangat ketat:

```text
.
├── cmd/
│   └── api/main.go          # Pintu masuk (Entry point) aplikasi.
├── docs/                    # Dokumentasi (Swagger + Panduan & Tutorial).
├── internal/
│   ├── config/              # Konfigurasi Database (GORM, Migrasi, Seed).
│   ├── handlers/            # Layer Terluar: Menerima HTTP Request, Parsing JSON (Gin).
│   ├── middlewares/         # Satpam API: Memverifikasi JWT dan Role Akses.
│   ├── models/              # Struktur Database (Entity/GORM Struct).
│   ├── repositories/        # Layer Terdalam: Akses kueri SQL Database murni.
│   ├── routes/              # Routing & Dependency Injection setup.
│   ├── services/            # Logika Bisnis Utama (Core Logic & Validation).
│   └── utils/               # Fungsi bantuan teknis (Hash, Token).
```
*Aturan Emas: Handler memanggil Service, Service memanggil Repository, Repository memanggil Database.*

---

## 🚀 Panduan Instalasi & Menjalankan Aplikasi

Aplikasi ini sudah dipaket dengan SQLite lokal dan sistem *Seeding* otomatis, sehingga Anda dapat langsung menjalankannya tanpa perlu menginstal database eksternal!

### 1. Prasyarat
Pastikan Anda sudah menginstal [Go](https://go.dev/dl/) versi `1.20` atau lebih tinggi di mesin Anda.

### 2. Kloning Repositori
```bash
git clone https://github.com/madajabbar/BAST-REQUEST.git
cd BAST-REQUEST
```

### 3. Unduh Dependensi
```bash
go mod tidy
```

### 4. Menjalankan Server
```bash
go run .\cmd\api\main.go
```
*Saat dijalankan pertama kali, aplikasi akan otomatis membuat file `bast_request.db`, membuat tabel, dan menyisipkan data contoh (termasuk Role).*

---

## 📚 Menjelajahi Dokumentasi API (Swagger)

Aplikasi memiliki antarmuka grafis (UI) untuk Anda menguji coba seluruh *Endpoint* secara langsung.
Saat server menyala, buka *browser* Anda dan kunjungi:

👉 **http://localhost:8080/swagger/index.html**

*Tips: Agar dapat mengakses Endpoint yang dikunci, Anda harus mendaftar (Register) atau masuk (Login) terlebih dahulu untuk mendapatkan token JWT, lalu memasukkan token tersebut ke dalam gembok "Authorize" (menggunakan format `Bearer <token>`).*

---

## 📖 Belajar & Panduan Lebih Dalam

Khusus untuk Anda yang ingin membedah kode ini atau baru belajar Golang, silakan buka **pusat dokumentasi** di folder [`docs/`](docs/README.md):

> 📂 **[`docs/README.md`](docs/README.md)** — Peta navigasi lengkap semua dokumen.

### Mulai dari sini (Pemula)
- 🚀 [**Gambaran Umum Aplikasi**](docs/getting-started/overview.md) — fitur, teknologi, alur kerja.
- 🚀 [**Panduan Instalasi**](docs/getting-started/installation.md) — setup, run, smoke test.

### Pahami Fondasi
- 🏗️ [**Clean Architecture**](docs/architecture/clean-architecture.md) — pola 4-layer + analogi restoran.
- 🗄️ [**Skema Database & ERD**](docs/architecture/database-schema-erd.md) — diagram Mermaid + penjelasan tabel.
- 🎓 [**Fondasi Golang**](docs/architecture/golang-fundamentals.md) — struct, pointer, error handling.

### Membedah Kode (Tutorial Step-by-Step)
1. [**Step 1: Setup & Konfigurasi**](docs/tutorials/step-01-setup-and-config.md)
2. [**Step 2: Models & Migrasi**](docs/tutorials/step-02-models-and-migration.md)
3. [**Step 3: Autentikasi JWT**](docs/tutorials/step-03-authentication-jwt.md)
4. [**Step 4: Master Data CRUD**](docs/tutorials/step-04-master-data-crud.md)
5. [**Step 5: Mesin Penomoran BAST**](docs/tutorials/step-05-bast-numbering-engine.md)
6. [**Step 6: Audit Log**](docs/tutorials/step-06-audit-log.md)
7. [**Step 7: Routing & RBAC**](docs/tutorials/step-07-routing-and-rbac.md)
8. [**Step 8: Dokumentasi Swagger**](docs/tutorials/step-08-swagger-documentation.md)

### Panduan Topik Spesifik
- 🔐 [**Autentikasi & RBAC End-to-End**](docs/guides/authentication-guide.md)
- 🤿 [**Deep Dive Penomoran BAST**](docs/guides/bast-numbering-deep-dive.md) — race condition, locking, reset.
- 🛠️ [**Menambahkan Fitur Baru**](docs/guides/add-new-feature-guide.md) — contoh modul Division.
- 📘 [**Panduan Swagger**](docs/guides/swagger-guide.md) — regenerasi & troubleshooting.

### Referensi API per Endpoint
- 🔌 [**Daftar Lengkap Endpoint**](docs/api-reference/README.md) — method, parameter, contoh request/response.

---

*Dikembangkan untuk efisiensi dan pencatatan yang solid.*
