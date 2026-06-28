# Fondasi Golang untuk Pemula

Dokumen ini ditujukan bagi Anda yang **baru pertama kali** menyentuh Go (Golang). Bahasa ini punya beberapa konsep unik yang berbeda dari PHP, Python, atau Java. Pahami 4 pilar ini dulu, sisanya akan jauh lebih mudah.

> 💡 Bacalah pelan-pelan. Setiap konsep disertai contoh dari kode nyata proyek ini.

---

## 1. Tidak Ada "Class" — Yang Ada "Struct"

Jika Anda pernah belajar OOP, Anda mengenal `Class`. **Di Go, itu tidak ada.** Sebagai gantinya, Go memakai **`struct`**: wadah (cetakan) untuk mengelompokkan beberapa variabel jadi satu kesatuan.

### Contoh dari proyek ini
[`internal/models/customer.go:10-18`](../../internal/models/customer.go):
```go
type Customer struct {
    CustomerID   uuid.UUID      `gorm:"type:uuid;primary_key"`
    CustomerCode string         `gorm:"type:varchar(50);uniqueIndex;not null"`
    CustomerName string         `gorm:"type:varchar(255);not null"`
    Status       string         `gorm:"type:varchar(50);default:'active'"`
    CreatedAt    time.Time
    UpdatedAt    time.Time
    DeletedAt    gorm.DeletedAt `gorm:"index"`
}
```

Penjelasan:
- `type Customer struct { ... }` → mendefinisikan cetakan bernama `Customer`.
- Tiap baris di dalamnya adalah **field** (variabel anggota).
- Tulisan dalam *backtick* `` `gorm:"..."` `` disebut **struct tag** — metadata yang dibaca GORM untuk membentuk kolom database.

### Cara pakai struct
```go
cust := Customer{
    CustomerCode: "CUST-001",
    CustomerName: "PT. Maju Mundur",
    Status:       "active",
}
```

### Method: "Perilaku" untuk Struct
Lalu bagaimana kita menambah *behavior* ke struct? Pakai **method** — fungsi yang menempel pada struct tertentu:

```go
// Fungsi ini dimiliki oleh struct Customer (receiver: c *Customer)
func (c *Customer) BeforeCreate(tx *gorm.DB) (err error) {
    if c.CustomerID == uuid.Nil {
        c.CustomerID = uuid.New()
    }
    return
}
```
Di atas adalah **hook** GORM: otomatis dipanggil sebelum `Create` untuk mengisi UUID.

---

## 2. Pointer (`*` dan `&`) — Topik Paling Memusingkan Pemula

Bayangkan variabel sebagai **rumah**:
- **Variabel biasa** = isi rumah (TV, meja, kursi).
- **Pointer** (`*`) = **alamat** rumah (Jl. Merdeka No.1).

### Mengapa pointer penting?
Secara bawaan, saat Anda melempar variabel ke fungsi, Go melakukan **copy-paste** (salin). Jika "rumah" besar (struct kompleks), copy-paste boros memori. Lebih parah: jika salinan diubah, **aslinya tidak berubah**!

Dengan memberi **pointer** (alamatnya saja), fungsi bisa pergi ke "rumah asli" dan mengubahnya langsung — tanpa copy-paste mahal.

### Contoh
```go
// Parameter (k *Karyawan) memakai *, artinya minta ALAMATNYA
func NaikkanGaji(k *Karyawan) {
    k.Gaji = k.Gaji + 1000000   // ubah data asli
}

func main() {
    pegawai := Karyawan{Gaji: 5000000}
    NaikkanGaji(&pegawai)        // & = ambil alamat pegawai
    // pegawai.Gaji sekarang 6.000.000!
}
```

### Di proyek ini: pointer di mana-mana
Hampir seluruh Handler, Service, Repository saling berkomunikasi via pointer agar cepat & hemat. Contoh konstruktor:

```go
// internal/repositories/customer_repository.go:12-14
func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
    return &CustomerRepository{db: db}
}
```
- `*CustomerRepository` (nilai kembalian) = fungsi ini mengembalikan **alamat** object repository.
- `&CustomerRepository{...}` = ambil alamat dari struct baru yang dibuat.

**Singkatnya:** `*` = "ini alamat tipe X", `&` = "ambil alamat variabel ini".

---

## 3. Penanganan Error — Tidak Ada `try/catch`

Go **tidak** memakai blok `try { } catch (e) { }`. Sebaliknya, error diperlakukan sebagai **nilai biasa** yang dikembalikan fungsi. Ini memaksa programmer **sadar** dan menangani error secara langsung.

### Pola standar: dua nilai kembalian
```go
func Bagi(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("tidak bisa dibagi nol")
    }
    return a / b, nil     // nil = tidak ada error
}

func main() {
    hasil, err := Bagi(10, 0)
    if err != nil {
        fmt.Println("Ada masalah:", err)
        return
    }
    fmt.Println("Hasil:", hasil)
}
```

### Di proyek ini
Hampir semua fungsi Repository mengembalikan `error`. Contoh handler yang mengeceknya:

```go
// internal/handlers/customer_handler.go:71-83
func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
    var input models.Customer
    if err := c.ShouldBindJSON(&input); err != nil {           // cek error binding
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := h.service.CreateCustomer(&input); err != nil {   // cek error service
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, input)
}
```

**Filosofi:** selalu cek error. Jangan pernah mengabaikannya. `nil` berarti "aman, lanjut".

---

## 4. Aturan Huruf Besar/Kecil (Exported vs Unexported)

Go mengatur visibilitas (public/private) **hanya lewat huruf pertama**:

| Huruf awal | Sifat | Bisa diakses dari? |
|---|---|---|
| **Kapital** (`CreateCustomer`, `Customer`) | **Exported** (public) | Package mana pun |
| **Kecil** (`createCustomer`, `customer`) | **Unexported** (private) | Hanya package yang sama |

### Contoh
```go
package models

type Customer struct {        // Exported: bisa dipakai di package handlers/services
    CustomerName string       // Exported: bisa diakses
    internalCode  string      // Unexported: HANYA bisa diakses dalam package models
}
```

### Implikasi penting
- Field struct yang ingin disimpan/ dibaca GORM **harus** diawali huruf kapital.
- Fungsi konstruktor `New...()` kapital agar bisa dipanggil dari package lain.
- Konstanta rahasia seperti `jwtKey` boleh kecil agar tersembunyi:

```go
// internal/utils/jwt.go:9
var jwtKey = []byte("bast-request-secret-key-05")   // unexported, private!
```

> ⚠️ Dalam produksi, secret sebaiknya dari **environment variable**, bukan hardcode. Ini ide pengembangan.

---

## 5. Package & Import

Setiap folder di bawah proyek adalah satu **package**. Nama package biasanya sama dengan nama folder.

```go
package handlers                          // file internal/handlers/*.go

import (
    "bast-request/internal/models"        // import package lain
    "bast-request/internal/services"
    "github.com/gin-gonic/gin"
)
```

Aturan: **semua file `.go` dalam satu folder HARUS punya nama package yang sama.**

---

## 6. Slice, Map, dan `make()`

Go punya beberapa tipe kumpulan:
- **Slice** (`[]Customer`) → mirip array dinamis.
- **Map** (`map[string]bool`) → mirip associative array/dictionary.

Contoh di proyek ini — validasi status pakai map:
```go
// internal/services/bast_request_service.go:76
validStatuses := map[string]bool{"Active": true, "Used": true, "Void": true}
if !validStatuses[status] {
    return errors.New("invalid status")
}
```

---

## 7. Goroutine & Defer (Sekilas)

Dua fitur andalan Go (meski di API ini dipakai minimal):

- **`go func()`** → menjalankan fungsi secara **asynchronous** (concurrent).
- **`defer`** → menunda eksekusi sampai fungsi selesai (biasanya untuk cleanup). Contoh: `defer db.Close()`.

---

## 8. Ringkasan Saran Latihan

Agar teori tidak menguap, praktekkan di komputer Anda:

1. Buat folder kosong: `mkdir belajargo && cd belajargo`
2. Init modul: `go mod init belajargo`
3. Buat `main.go` sederhana yang print "Hello".
4. Buat `struct` `Barang`, isi field, lalu print.
5. Coba buat fungsi dengan parameter pointer vs non-pointer, amati bedanya.
6. Cobalah pola error handling (dua nilai kembalian).
7. Setelah nyaman, lanjut baca [Clean Architecture](clean-architecture.md) & [Tutorial Step 1](../tutorials/step-01-setup-and-config.md).

---

## 9. Bacaan Lanjutan
- 🌐 [Tour of Go (resmi)](https://go.dev/tour/) — interaktif, sangat direkomendasikan.
- 📘 [Effective Go](https://go.dev/doc/effective_go)
- 🏗️ Terapkan konsep di kode nyata: [Clean Architecture](clean-architecture.md)
