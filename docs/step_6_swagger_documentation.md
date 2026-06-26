# Step 6: Dokumentasi API (Swagger)

Langkah terakhir dari pembangunan aplikasi ini adalah memastikan Frontend Developer atau Klien dapat mengkonsumsi API yang kita buat dengan mudah, tanpa perlu menebak-nebak *body JSON* apa yang harus dikirim. Di Golang, kita menggunakan `swaggo/swag`.

## 1. Menambahkan Anotasi ke Root Aplikasi
Pertama, kita menyisipkan deklarasi identitas aplikasi di atas fungsi `main()` pada `cmd/api/main.go`.

```go
// @title BAST Request API
// @version 1.0
// @description This is the API server for BAST Request System.
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
    // ...
}
```
Anotasi `securityDefinitions` di atas sangat penting agar antarmuka Swagger memunculkan fitur masuk/input *Bearer Token* secara universal (Icon gembok).

## 2. Menambahkan Anotasi di Tiap Handler
Setiap fungsi di dalam `internal/handlers/*.go` diberikan komentar rapi yang berfungsi sebagai dokumentasi terstruktur.

Contoh di `project_handler.go`:
```go
// CreateProject godoc
// @Summary Create a new project
// @Description Add a new project linked to a customer
// @Tags projects
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param project body models.Project true "Project Data"
// @Success 201 {object} models.Project
// @Router /projects [post]
func (h *ProjectHandler) CreateProject(c *gin.Context) {
    // ... isi kode handler
}
```
Penjelasan Tag:
- `@Security BearerAuth`: Menandakan bahwa fungsi ini terkunci dan butuh Token. (Sesuai nama yang didefinisikan di `main.go`).
- `@Param`: Apa yang harus diisi di *Request Body*? (Di sini adalah struktur `models.Project`).
- `@Success`: Apa tipe balasan jika *Request* sukses? (Mengembalikan JSON format dari struct `models.Project`).
- `@Router`: URL dan *Method HTTP*-nya.

## 3. Proses Generasi File JSON Swagger
Golang sendiri tidak mengerti komentar ini. Kita perlu pihak ketiga (Aplikasi *Swag CLI*) untuk membaca file `.go` kita, mengekstrak komentarnya, dan mengubahnya menjadi *file* spesifikasi OpenAPI yang diakui dunia internasional (`docs/swagger.json`).

Perintah yang dieksekusi di Terminal (Root):
```bash
swag init -g ./cmd/api/main.go --parseDependency --parseInternal
```
- `-g` menunjuk letak fungsi `main()` berada.
- `--parseDependency` memaksa *swag* untuk membaca tipe data lain (seperti *uuid* dan *json.RawMessage*) dari repositori di luar aplikasi ini jika diperlukan.

## 4. Meng-hosting Halaman Web Swagger (UI)
Untuk menampilkan halaman *web interaktif* agar API bisa langsung diklik, kita menghubungkan rute khusus di dalam Gin (`internal/routes/routes.go`).

```go
import (
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
    _ "bast-request/docs" // Memanggil folder docs hasil generate (Penting!)
)

func SetupRoutes(...) {
    // ...
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
```

## Ringkasan Step 6:
Dokumentasi *API Swagger* kini tertanam rapi ke dalam kode sumber, membuat dokumentasi API dan aplikasi selalu terhubung (*single source of truth*). Jika Anda mengubah struct tabel di Database, dokumentasi JSON API klien juga ikut berubah setelah Anda memanggil ulang fungsi *swag init*!

---
**Tamat.** Anda telah menyusuri pilar-pilar pondasi utama bagaimana **BAST Request API** dirakit menggunakan Golang!
