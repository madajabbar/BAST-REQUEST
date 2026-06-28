# 📚 Dokumentasi BAST Request API

Selamat datang di pusat dokumentasi **BAST Request API**. Dokumentasi ini disusun untuk membantu Anda — baik sebagai pemula Golang maupun developer berpengalaman — memahami, menjalankan, hingga membedah kode aplikasi ini dari nol sampai mahir.

> **Bahasa:** Seluruh dokumentasi ditulis dalam Bahasa Indonesia agar mudah dipelajari. Contoh kode selalu diambil dari kode sumber nyata proyek ini (lengkap dengan path & nomor baris).

---

## 🗺️ Peta Navigasi

Pilih jalur belajar sesuai kebutuhan Anda:

### 🚀 Baru pertama kali di sini?
Mulai dari sini untuk mengenal aplikasi & menjalankannya dalam 5 menit.

| Dokumen | Isi |
|---|---|
| [Gambaran Umum](getting-started/overview.md) | Apa itu BAST Request API, fitur utama, teknologi, dan alur kerja aplikasi. |
| [Panduan Instalasi](getting-started/installation.md) | Prasyarat, clone, instalasi dependensi, menjalankan server, dan uji coba pertama. |

### 🏗️ Memahami Arsitektur & Konsep
Fondasi penting sebelum menyentuh kode.

| Dokumen | Isi |
|---|---|
| [Clean Architecture](architecture/clean-architecture.md) | Pola 4-layer (Handler→Service→Repository→DB), analogi restoran, dan alur data. |
| [Skema Database & ERD](architecture/database-schema-erd.md) | Diagram ERD Mermaid, penjelasan tiap tabel, relasi, soft delete, dan unique index. |
| [Fondasi Golang](architecture/golang-fundamentals.md) | Struct, Pointer, Error Handling, Exported vs Unexported — khusus pemula. |

### 🎓 Tutorial Step-by-Step (Seri Utama)
Membedah kode aplikasi langkah-demi-langkah, dari database sampai Swagger. **Baca berurutan dari Step 1.**

| Step | Dokumen | Fokus |
|---|---|---|
| 1 | [Setup & Konfigurasi](tutorials/step-01-setup-and-config.md) | `go mod init`, koneksi DB, `main.go` |
| 2 | [Models & Migrasi](tutorials/step-02-models-and-migration.md) | Struct GORM, `AutoMigrate`, seeding |
| 3 | [Autentikasi JWT](tutorials/step-03-authentication-jwt.md) | Bcrypt, JWT, Register/Login, Middleware |
| 4 | [Master Data CRUD](tutorials/step-04-master-data-crud.md) | Customer/Project/Format: Repo→Service→Handler |
| 5 | [Mesin Penomoran BAST](tutorials/step-05-bast-numbering-engine.md) | Pattern, sequence atomik, transaksi DB |
| 6 | [Audit Log](tutorials/step-06-audit-log.md) | Pelacakan perubahan data (old_data/new_data) |
| 7 | [Routing & RBAC](tutorials/step-07-routing-and-rbac.md) | `routes.go`, RequireAuth, RequireRole |
| 8 | [Dokumentasi Swagger](tutorials/step-08-swagger-documentation.md) | Anotasi `@Summary`, `swag init`, hosting UI |

### 📖 Panduan Topik Spesifik
Pendalaman untuk fitur tertentu.

| Dokumen | Isi |
|---|---|
| [Panduan Autentikasi & RBAC](guides/authentication-guide.md) | Alur lengkap Register→Login→Token→Akses terlindungi. |
| [Deep Dive Penomoran BAST](guides/bast-numbering-deep-dive.md) | Race condition, `clause.Locking`, reset sequence, format pattern. |
| [Menambahkan Fitur Baru](guides/add-new-feature-guide.md) | Tutorial praktis menambah modul `Division` dari Model → Route. |
| [Panduan Swagger](guides/swagger-guide.md) | Cara regenerasi dokumentasi & troubleshooting. |

### 🔌 Referensi API (per Endpoint)
Detail teknis tiap endpoint: method, path, parameter, contoh request & response JSON.

| Dokumen | Endpoint |
|---|---|
| [Cara Pakai Referensi API](api-reference/README.md) | Daftar lengkap & cara autentikasi |
| [Auth](api-reference/auth-endpoints.md) | `POST /auth/register`, `POST /auth/login` |
| [Customer](api-reference/customer-endpoints.md) | `/customers` (5 endpoint) |
| [Project](api-reference/project-endpoints.md) | `/projects` (5 endpoint) |
| [BAST Format](api-reference/bast-format-endpoints.md) | `/bast-formats` (5 endpoint) |
| [BAST Sequence](api-reference/bast-sequence-endpoints.md) | `/bast-sequences`, `/reset` |
| [BAST Request](api-reference/bast-request-endpoints.md) | `/bast-requests` (5 endpoint) |
| [Audit Log](api-reference/audit-log-endpoints.md) | `/audit-logs` (2 endpoint) |

---

## 🧭 Rekomendasi Jalur Belajar

**Saya pemula di Golang:**
> `getting-started/overview` → `getting-started/installation` → `architecture/golang-fundamentals` → `architecture/clean-architecture` → `tutorials/step-01` s.d. `step-08`

**Saya sudah jago Go, hanya ingin paham kode ini:**
> `getting-started/installation` → `architecture/clean-architecture` → `architecture/database-schema-erd` → `api-reference/` (semua)

**Saya ingin menambah fitur baru:**
> `guides/add-new-feature-guide` → lihat pola di `tutorials/step-04`

---

## 📂 Catatan Struktur Folder

```text
docs/
├── README.md              ← Anda di sini (index ini)
├── docs.go                ← (jangan diutak-atik) package Swagger, di-import routes.go
├── swagger.json           ← (jangan diutak-atik) hasil generate swag
├── swagger.yaml           ← (jangan diutak-atik) hasil generate swag
├── getting-started/       ← pengenalan & instalasi
├── architecture/          ← konsep & desain
├── tutorials/             ← seri membedah kode (8 step)
├── guides/                ← panduan topik spesifik
└── api-reference/         ← referensi endpoint
```

> ⚠️ **Penting:** File `docs.go`, `swagger.json`, dan `swagger.yaml` adalah hasil generate dari [swaggo/swag](https://github.com/swaggo/swag). Ketiganya **tetap berada di root** `docs/` karena merupakan package Go yang di-import oleh `internal/routes/routes.go` (baris `_ "bast-request/docs"`). Jangan memindahkannya ke subfolder, kecuali Anda siap mengubah import path dan flag `-o` saat `swag init`.

---

## 🌐 Swagger UI Interaktif

Selain dokumentasi Markdown ini, aplikasi juga menyediakan UI Swagger yang bisa langsung diklik untuk mencoba API:

👉 **http://localhost:8080/swagger/index.html** (setelah server berjalan)

Pelajari cara menggunakannya di [Panduan Swagger](guides/swagger-guide.md).

---

*Dokumentasi ini akan terus berkembang bersama kode. Jika menemui kesalahan, silakan perbaiki langsung atau laporkan.*
