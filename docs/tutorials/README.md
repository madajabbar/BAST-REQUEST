# 🎓 Seri Tutorial Step-by-Step

Seri ini **membedah kode aplikasi BAST Request API** langkah demi langkah, dari pondasi paling bawah (database) sampai paling atas (dokumentasi Swagger).

> **Cara terbaik:** baca **berurutan** dari Step 1. Tiap langkah merujuk ke **file & nomor baris kode nyata**, jadi Anda bisa membuka kode sambil membaca.

---

## Peta Seri

| # | Tutorial | Yang Dipelajari | File Kode Utama |
|---|---|---|---|
| 1 | [Setup & Konfigurasi](step-01-setup-and-config.md) | `go mod init`, koneksi DB SQLite, fungsi `main()` | `cmd/api/main.go`, `internal/config/database.go` |
| 2 | [Models & Migrasi](step-02-models-and-migration.md) | Struct GORM, `AutoMigrate`, seeding data | `internal/models/*`, `internal/config/seed.go` |
| 3 | [Autentikasi JWT](step-03-authentication-jwt.md) | Bcrypt, JWT, Register/Login, Middleware | `internal/utils/*`, `internal/services/auth_service.go`, `internal/middlewares/*` |
| 4 | [Master Data CRUD](step-04-master-data-crud.md) | Pola Repository → Service → Handler | `internal/{repositories,services,handlers}/customer_*`, `project_*`, `bast_format_*` |
| 5 | [Mesin Penomoran BAST](step-05-bast-numbering-engine.md) | Format pattern, sequence atomik, transaksi | `internal/services/bast_request_service.go`, `bast_sequence_*` |
| 6 | [Audit Log](step-06-audit-log.md) | Audit trail, JSON old/new data, filter query | `internal/models/audit_log.go`, `audit_log_*` |
| 7 | [Routing & RBAC](step-07-routing-and-rbac.md) | `routes.go`, dependency injection, `RequireAuth`/`RequireRole` | `internal/routes/routes.go`, `internal/middlewares/*` |
| 8 | [Dokumentasi Swagger](step-08-swagger-documentation.md) | Anotasi `@Summary`, `swag init`, hosting UI | `cmd/api/main.go`, semua handler |

---

## Prasyarat

Sebelum mulai, pastikan Anda sudah:
- ✅ Menginstal Go 1.20+ (cek: `go version`)
- ✅ Membaca [Gambaran Umum](../getting-started/overview.md)
- ✅ Membaca [Panduan Instalasi](../getting-started/installation.md) & berhasil menjalankan server
- ✅ (Disarankan) Membaca [Fondasi Golang](../architecture/golang-fundamentals.md) & [Clean Architecture](../architecture/clean-architecture.md)

---

## Filosofi Seri Ini

Setiap tutorial mengikuti pola:
1. **Konsep** — apa & mengapa.
2. **Kode Nyata** — kutipan dari repo, lengkap dengan path + nomor baris (klik untuk langsung loncat).
3. **Penjelasan baris penting** — kalimat per kalimat.
4. **Uji coba** — cara melihat hasilnya bekerja.
5. **Navigasi** — link ke step sebelumnya & selanjutnya.

➡️ **Mulai dari: [Step 1 — Setup & Konfigurasi](step-01-setup-and-config.md)**
