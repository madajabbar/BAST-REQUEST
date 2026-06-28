# Panduan Instalasi — BAST Request API

Panduan ini mengantarkan Anda dari nol hingga server berjalan dan bisa diuji coba dalam waktu kurang dari 5 menit. Aplikasi ini sengaja dirancang *self-contained*: **tidak perlu instal database eksternal** apa pun.

---

## 1. Prasyarat (Wajib)

Satu-satunya hal yang **wajib** Anda miliki adalah **Go (Golang)**.

| Prasyarat | Versi | Cara Cek | Link Unduh |
|---|---|---|---|
| **Go** | `1.20` atau lebih tinggi | `go version` | [go.dev/dl](https://go.dev/dl/) |

**Opsional (untuk regenerasi Swagger):**
- [Git](https://git-scm.com/) — untuk meng-clone repositori.
- `swag` CLI — lihat [Panduan Swagger](../guides/swagger-guide.md) jika ingin regenerasi dokumentasi.

### Verifikasi instalasi Go
Buka terminal (Command Prompt / PowerShell / Git Bash), lalu jalankan:
```bash
go version
```
Jika muncul `go version go1.2x.x ...`, berarti Go sudah terinstal dengan benar.

---

## 2. Mengkloning Repositori

```bash
git clone https://github.com/madajabbar/BAST-REQUEST.git
cd BAST-REQUEST
```

> **Tanpa Git?** Anda juga bisa mengunduh arsip ZIP dari halaman GitHub, lalu mengekstraknya dan membuka folder hasil ekstrak di terminal.

---

## 3. Mengunduh Dependensi

Aplikasi memakai beberapa *library* eksternal (Gin, GORM, JWT, bcrypt, dll). Semua tercatat di `go.mod`. Unduh semuanya dengan:

```bash
go mod tidy
```

Perintah ini akan:
- Membaca `go.mod` & `go.sum`,
- Mengunduh semua paket yang dibutuhkan,
- Menyimpannya di *module cache* Go (biasanya `$GOPATH/pkg/mod`).

> **Tips:** Jika koneksi lambat, Anda bisa set proxy module Go bawaan: `go env -w GOPROXY=https://goproxy.io,direct`.

---

## 4. Menjalankan Server

Dari root proyek, jalankan:

```bash
go run ./cmd/api/main.go
```

### Apa yang terjadi saat pertama kali dijalankan?

Urutan eksekusi mengikuti fungsi `main()` di [`cmd/api/main.go:18-42`](../../cmd/api/main.go):

1. **`config.ConnectDB()`** → membuat / membuka file `bast_request.db` di root proyek.
2. **`config.AutoMigrate()`** → membuat semua tabel otomatis (kalau belum ada).
3. **`config.SeedDB(config.DB)`** → menyisipkan data awal (Role, Customer, Project, Format contoh).
4. **`gin.Default()`** → menyiapkan web server Gin.
5. **`routes.SetupRoutes(...)`** → mendaftarkan semua endpoint API.
6. **`r.Run(":8080")`** → server mendengarkan permintaan di **port 8080**.

### Output yang diharapkan di terminal:
```text
Database connection successfully opened
Database Migration Completed
Seeding roles...
Seeding database with dummy data...
Database successfully seeded.
[GIN-debug] listening and serving HTTP on :8080
```

Sekarang server hidup di **http://localhost:8080** 🎉

---

## 5. Uji Coba Pertama (Smoke Test)

### A. Tes endpoint publik `/ping`
Buka terminal baru (biarkan server tetap jalan), lalu:
```bash
curl http://localhost:8080/ping
```
**Respons yang diharapkan:**
```json
{ "message": "pong" }
```
Endpoint ini didefinisikan langsung di [`cmd/api/main.go:31-35`](../../cmd/api/main.go) dan **tidak butuh token**.

### B. Buka Swagger UI
Buka browser ke:
👉 **http://localhost:8080/swagger/index.html**

Anda akan melihat daftar lengkap endpoint yang dikelompokkan per tag (`Auth`, `customers`, `projects`, dll).

### C. Mencoba API yang butuh autentikasi
Sebagian besar endpoint terkunci. Untuk mencobanya, Anda perlu token JWT. Alurnya:

1. **Register** akun baru → `POST /api/auth/register`
2. **Login** → `POST /api/auth/login` → dapat `token`
3. Klik ikon 🔒 **Authorize** di Swagger, masukkan `Bearer <token-anda>`
4. Sekarang endpoint terkunci bisa dipakai.

Detail lengkap dengan contoh JSON ada di [Referensi API Auth](../api-reference/auth-endpoints.md).

---

## 6. Menghentikan Server

Tekan `Ctrl + C` di terminal tempat server berjalan.

---

## 7. Troubleshooting Umum

| Masalah | Penyebab | Solusi |
|---|---|---|
| `go: command not found` | Go belum terinstal / tidak di PATH | Instal dari go.dev/dl, lalu buka terminal baru |
| `bind: address already in use :8080` | Port 8080 dipakai aplikasi lain | Tutup aplikasi itu, atau ubah port di `cmd/api/main.go:41` (`r.Run(":8081")`) |
| `go mod tidy` gagal / lambat | Proxy module bermasalah | `go env -w GOPROXY=https://goproxy.io,direct` lalu ulangi |
| Database error / data rusak | File `bast_request.db` korup | Hapus file `bast_request.db`, aplikasi akan membuat ulang saat dijalankan |
| Swagger UI kosong / tidak update | Belum regenerasi setelah ubah kode | Jalankan `swag init` (lihat [Panduan Swagger](../guides/swagger-guide.md)) |

> 💡 **Mau mulai dari bersih?** Hapus file `bast_request.db`. Saat server dijalankan ulang, seluruh tabel & seed data akan dibuat kembali.

---

## 8. File Penting yang Dibuat Saat Run

| File | Dibuat oleh | Fungsi |
|---|---|---|
| `bast_request.db` | `config.ConnectDB()` | Database SQLite lokal (semua data tersimpan di sini) |
| `main.exe` *(opsional)* | `go build` | Binary jika Anda melakukan *build* manual |

> ⚠️ Jangan commit `bast_request.db` ke Git — itu adalah data lokal. Idealnya tambahkan ke `.gitignore`.

---

## 9. Langkah Selanjutnya

Server sudah jalan, sekarang waktunya **memahami kode** yang membuatnya bekerja:

➡️ **Memahami konsep dulu:** [Clean Architecture](../architecture/clean-architecture.md)
➡️ **Atau langsung bedah kode:** [Tutorial Step 1 — Setup & Konfigurasi](../tutorials/step-01-setup-and-config.md)
