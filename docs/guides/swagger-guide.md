# Panduan Swagger — Regenerasi & Troubleshooting

Panduan praktis cara mengelola dokumentasi Swagger/OpenAPI di proyek ini: instalasi, regenerasi, penggunaan UI, dan penyelesaian masalah umum.

> 📚 Pendalaman dari [Tutorial Step 8](../tutorials/step-08-swagger-documentation.md).

---

## 1. Arsitektur Swagger di Proyek Ini

```text
Komentar anotasi (@Summary, @Router, ...)
              │
              ▼  dibaca oleh
        [ swag init CLI ]
              │
              ▼  menulis
   docs/docs.go + swagger.json + swagger.yaml
              │
              ▼  di-import (blank import) di routes.go
        [ Swagger UI ter-host ]
              │
              ▼  diakses di
   http://localhost:8080/swagger/index.html
```

**Filosofi:** dokumentasi = kode. Ubah kode → regenerasi → dokumen update otomatis.

---

## 2. Instalasi Swag CLI (sekali)

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Verifikasi:
```bash
swag --version
```
Jika `swag: command not found`, tambahkan `$GOPATH/bin` ke PATH:
- **Linux/Mac:** `export PATH=$PATH:$(go env GOPATH)/bin` (tambah ke `~/.bashrc`).
- **Windows:** tambah `%USERPROFILE%\go\bin` ke Environment Variable PATH.

---

## 3. Regenerasi Dokumentasi

Setiap kali Anda **mengubah anotasi Swagger** atau **menambah/mengubah handler**, jalankan dari root proyek:

```bash
swag init -g cmd/api/main.go --parseDependency --parseInternal
```

### Penjelasan flag
| Flag | Wajib? | Fungsi |
|---|---|---|
| `-g cmd/api/main.go` | ✅ | Lokasi file dengan anotasi global (`@title`, `@host`, dst) |
| `--parseDependency` | ✅ | Baca tipe dari package eksternal (`uuid.UUID`, `datatypes.JSON`) |
| `--parseInternal` | ✅ | Baca tipe dari package internal proyek (`models.*`) |
| `-o docs` | opsional | Folder output (default: `docs/`) |
| `--quiet` | opsional | Kurangi output log |

> ⚠️ **Tanpa** `--parseDependency --parseInternal`, banyak tipe (terutama `uuid.UUID`) tidak ter-resolve dan Swagger jadi error/kosong.

### Output
Swag menulis **3 file** ke `docs/`:
- `docs.go` — package Go berisi spesifikasi ter-embed.
- `swagger.json` — spesifikasi OpenAPI format JSON.
- `swagger.yaml` — versi YAML.

**Jangan edit manual** — akan ditimpa saat `swag init` berikutnya.

---

## 4. Anotasi yang Wajib & Opsional

### Di `main.go` (global) — wajib sekali
```go
// @title Nama API
// @version 1.0
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
```

### Per handler — minimal
| Tag | Wajib? | Contoh |
|---|---|---|
| `@Summary` | ✅ | `Create a new customer` |
| `@Router` | ✅ | `/customers [post]` |
| `@Tags` | ✅ (untuk grouping) | `customers` |
| `@Produce` | opsional | `json` |
| `@Accept` | opsional | `json` |
| `@Param` | jika ada param | `customer body models.Customer true "..."` |
| `@Success` | disarankan | `201 {object} models.Customer` |
| `@Failure` | disarankan | `400 {object} map[string]interface{}` |
| `@Security` | opsional | `BearerAuth` (tandai butuh token) |

---

## 5. Menggunakan Swagger UI

### A. Buka UI
Setelah server jalan, buka browser:
👉 **http://localhost:8080/swagger/index.html**

### B. Otorisasi (untuk endpoint terkunci)
1. Klik tombol **Authorize** 🔒 (kanan atas halaman).
2. Masukkan token dengan format:
   ```
   Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6...
   ```
   ⚠️ **Penting:** sertakan kata `Bearer ` diikuti spasi, baru token.
3. Klik **Authorize** → gembok tertutup.

Cara dapat token: jalankan `POST /api/auth/login` di Swagger itu sendiri, copy field `token` dari respons.

### C. Test endpoint
1. Klik endpoint → expand.
2. Klik **Try it out** (kanan).
3. Isi parameter / body JSON.
4. Klik **Execute**.
5. Lihat **Response body** & **Code** di bawah.

---

## 6. Troubleshooting

### Masalah: Swagger UI 404 / blank
**Kemungkinan & solusi:**
- ❌ Belum `swag init` → jalankan.
- ❌ Lupa blank import di routes.go → cek ada `_ "bast-request/docs"`.
- ❌ Server belum direstart setelah `swag init` → restart.

### Masalah: Endpoint baru tidak muncul
- ❌ Lupa `@Router` di handler → tambahkan.
- ❌ Lupa `swag init` setelah tambah handler → jalankan.
- ❌ Route tidak didaftarkan di `routes.go` → cek.

### Masalah: Token di Authorize tidak terkirim
- ❌ Format salah → harus `Bearer <token>` (dengan spasi, B huruf besar).
- ❌ Endpoint tidak diberi `@Security BearerAuth` → Swagger tidak tahu butuh token (tapi middleware tetap cek).

### Masalah: Tipe `uuid.UUID` / `datatypes.JSON` unknown
- ❌ Lupa flag `--parseDependency --parseInternal` → tambahkan.

### Masalah: `swag: command not found`
- ❌ `$GOPATH/bin` tidak di PATH → tambahkan.

### Masalah: Error parsing anotasi
Swag akan print error spesifik, mis. `Failed to find param type`. Cek syntax `@Param`:
```
@Param <nama> <body|query|path> <tipe> <required> "<deskripsi>"
```
Contoh benar: `@Param id path string true "Customer ID"`.

### Masalah: Respons schema kosong
- ❌ `@Success` salah format → `@Success 200 {object} models.Customer` (perhatikan `{object}`).

---

## 7. Best Practices

✅ **Lakukan:**
- Regenerasi `swag init` **setiap** ubah handler (masukkan ke workflow/git hook).
- Tulis `@Description` yang jelas — bantu pengguna API.
- Pakai tag `@Tags` yang konsisten (mis. semua endpoint customer → tag `customers`).
- Commit file hasil generate (`docs.go`, `swagger.json`, `swagger.yaml`) agar repo self-contained.

❌ **Hindari:**
- Edit `swagger.json`/`docs.go` manual.
- Lupa flag `--parseDependency --parseInternal`.
- Pakai `map[string]interface{}` untuk semua respons (tidak deskriptif) — sebaiknya buat DTO struct khusus.

---

## 8. Tips Lanjutan

### A. Buat script regenerasi
Buat file `regen-docs.sh` (atau `.bat` di Windows):
```bash
#!/bin/bash
swag init -g cmd/api/main.go --parseDependency --parseInternal
echo "✅ Swagger regenerated"
```
Jalankan: `./regen-docs.sh`.

### B. Pisah ke file swagger
Jika ingin dokumentasi terpisah dari kode, bisa generate ke folder lain:
```bash
swag init -g cmd/api/main.go -o docs/swagger --parseDependency --parseInternal
```
Tapi konsekuensinya: ubah import di routes.go jadi `_ "bast-request/docs/swagger"`. Di proyek ini kita **biarkan di root** `docs/`.

### C. Ekspor untuk klien eksternal
File `swagger.json` bisa diimpor ke:
- [Postman](https://www.postman.com/) — Import → pilih file.
- [Insomnia](https://insomnia.rest/).
- [Stoplight](https://stoplight.io/).
- Generator client code ([openapi-generator](https://openapi-generator.tech/)) → buat SDK otomatis.

---

## 9. Workflow Harian yang Disarankan

```text
1. Ubah/tambah handler (dengan anotasi Swagger)
2. Ubah/tambah route
3. Jalankan: swag init -g cmd/api/main.go --parseDependency --parseInternal
4. Restart server: go run ./cmd/api/main.go
5. Cek di /swagger/index.html
6. Uji coba endpoint
7. Commit (termasuk docs.go, swagger.json, swagger.yaml)
```

---

## Bacaan Lanjutan
- 📘 [swaggo/swag docs](https://github.com/swaggo/swag#how-to-use-it-with-gin)
- 📘 [OpenAPI Specification](https://swagger.io/specification/)
- 🛠️ [Tutorial Step 8 — Swagger](../tutorials/step-08-swagger-documentation.md)
- 📖 [Referensi API](../api-reference/README.md)
