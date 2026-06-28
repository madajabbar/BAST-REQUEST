# Step 2: Models, Migrasi & Seeding

> Seri Tutorial · **Step 2 dari 8**

Pada langkah ini kita membedah folder `internal/models/`: bagaimana `struct` Go "ajaib" berubah menjadi tabel SQL, apa arti tag GORM, hook `BeforeCreate`, soft delete, serta bagaimana data awal (seed) disisipkan.

---

## 1. Apa Itu Model?

**Model** adalah `struct` Go yang merepresentasikan **satu baris data** di tabel database. GORM membaca **tag** (tulisan dalam *backtick* di samping field) untuk merancang kolom tabel.

---

## 2. Anatomi Satu Model: `Customer`

File: [`internal/models/customer.go`](../../internal/models/customer.go)

```go
package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Customer struct {
	CustomerID   uuid.UUID      `gorm:"type:uuid;primary_key"`
	CustomerCode string         `gorm:"type:varchar(50);uniqueIndex;not null"`
	CustomerName string         `gorm:"type:varchar(255);not null"`
	Status       string         `gorm:"type:varchar(50);default:'active'"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// TableName menimpa nama tabel default
func (Customer) TableName() string {
	return "master_customer"
}

// Hook sebelum Create
func (c *Customer) BeforeCreate(tx *gorm.DB) (err error) {
	if c.CustomerID == uuid.Nil {
		c.CustomerID = uuid.New()
	}
	return
}
```

### Penjelasan bagian per bagian

#### A. Tag GORM
| Tag | Arti |
|---|---|
| `type:uuid` | Kolom bertipe `UUID` |
| `primary_key` | Primary key tabel |
| `varchar(50)` | String maksimal 50 karakter |
| `uniqueIndex` | Tidak boleh duplikat (DB menolak) |
| `not null` | Wajib diisi |
| `default:'active'` | Nilai default jika tidak diisi saat insert |
| `index` | Buat index (untuk soft delete & query cepat) |

#### B. `TableName()` — Override Nama Tabel
Tanpa method ini, GORM akan otomatis menamai tabel versi jamak dari struct: `customers`. Tapi di sini kita **paksa** nama `master_customer` agar konsisten dengan konvensi "master data".

```go
func (Customer) TableName() string {
	return "master_customer"
}
```

#### C. Hook `BeforeCreate`
Hook adalah fungsi yang **otomatis dipanggil GORM** pada momen tertentu. `BeforeCreate` dipanggil tepat sebelum `INSERT`:

```go
func (c *Customer) BeforeCreate(tx *gorm.DB) (err error) {
	if c.CustomerID == uuid.Nil {     // kalau ID masih kosong
		c.CustomerID = uuid.New()      // generate UUID baru
	}
	return
}
```

**Tujuan:** agar kita tidak perlu manual generate ID saat create. Cukup kirim data tanpa ID, GORM+hook yang isi otomatis.

#### D. Field Waktu
- `CreatedAt` / `UpdatedAt` → GORM mengisi otomatis saat `Create`/`Save`.
- `DeletedAt gorm.DeletedAt` → **mengaktifkan soft delete**.

### Soft Delete — Penjelasan Mendalam
Dengan `DeletedAt gorm.DeletedAt`, saat Anda "menghapus" record via `db.Delete(...)`, GORM **tidak benar-benar menghapus** baris. Ia hanya mengisi `deleted_at` dengan timestamp. Selanjutnya, query normal (`Find`, `First`) **otomatis menyembunyikan** record yang sudah di-soft-delete.

Manfaat: data tidak hilang permanen, bisa di-restore, aman dari kesalahan hapus.

> ⚠️ **Catatan:** Di proyek ini, beberapa operasi "delete" (Customer, Project) justru **tidak** memakai soft delete, melainkan mengubah `status` jadi `inactive`. Lihat [`customer_repository.go:45-48`](../../internal/repositories/customer_repository.go). Ini pilihan desain agar data tetap muncul tapi ditandai non-aktif.

---

## 3. Daftar Model Lengkap

| Model | File | Primary Key | Nama Tabel | Catatan |
|---|---|---|---|---|
| `Role` | [`role.go`](../../internal/models/role.go) | `RoleID` | `roles` | Tanpa `TableName()` → auto-plural |
| `User` | [`user.go`](../../internal/models/user.go) | `UserID` | `users` | Punya relasi `Role` |
| `Customer` | [`customer.go`](../../internal/models/customer.go) | `CustomerID` | `master_customer` | Soft delete |
| `Project` | [`project.go`](../../internal/models/project.go) | `ProjectID` | `master_project` | FK ke Customer |
| `BastFormat` | [`bast_format.go`](../../internal/models/bast_format.go) | `FormatID` | `master_bast_format` | Menyimpan pattern nomor |
| `BastSequence` | [`bast_sequence.go`](../../internal/models/bast_sequence.go) | `SequenceID` | `bast_sequence` | Unique index komposit |
| `BastRequest` | [`bast_request.go`](../../internal/models/bast_request.go) | `BastRequestID` | `bast_request` | Tabel transaksi inti |
| `AuditLog` | [`audit_log.go`](../../internal/models/audit_log.go) | `AuditLogID` | `audit_log` | Tipe JSON (`datatypes.JSON`) |

> Lihat diagram lengkapnya di [Skema Database & ERD](../architecture/database-schema-erd.md).

---

## 4. Unique Index Komposit di `BastSequence`

Contoh lanjutan — model [`internal/models/bast_sequence.go:10-19`](../../internal/models/bast_sequence.go):

```go
type BastSequence struct {
	SequenceID uuid.UUID  `gorm:"type:uuid;primary_key"`
	FormatID   uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:idx_format_year_month"`
	Format     BastFormat `gorm:"foreignKey:FormatID;references:FormatID"`
	Year       int        `gorm:"not null;uniqueIndex:idx_format_year_month"`
	Month      int        `gorm:"not null;uniqueIndex:idx_format_year_month"`
	LastNumber int        `gorm:"not null;default:0"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
```

Perhatikan tiga field (`FormatID`, `Year`, `Month`) memakai tag yang sama:
```
uniqueIndex:idx_format_year_month
```
Nama index yang sama (`idx_format_year_month`) menyatukan ketiganya menjadi **satu unique index komposit**. Artinya: kombinasi `(format_id, year, month)` harus unik — tidak boleh ada 2 baris dengan format+periode yang sama. Ini jantung anti-race-condition penomoran BAST (dibahas [Step 5](step-05-bast-numbering-engine.md)).

### Relasi (Foreign Key)
```go
Format BastFormat `gorm:"foreignKey:FormatID;references:FormatID"`
```
Field `Format` ini **bukan kolom database**, melainkan "anchor" relasi. Saat kita `Preload("Format")`, GORM otomatis JOIN data format terkait.

---

## 5. AutoMigrate (Rekap)

[`internal/config/database.go:26-42`](../../internal/config/database.go) mendaftarkan semua model:
```go
DB.AutoMigrate(
	&models.Customer{},
	&models.Project{},
	&models.BastFormat{},
	&models.BastSequence{},
	&models.BastRequest{},
	&models.AuditLog{},
	&models.User{},
	&models.Role{},
)
```

> **Urutan penting?** Sebenarnya tidak, karena GORM menangani dependency FK. Tapi biasanya tabel yang jadi "induk" (yang dirujuk) diletakkan lebih awal demi keterbacaan.

---

## 6. Seeding — Mengisi Data Awal

File: [`internal/config/seed.go`](../../internal/config/seed.go)

Seeding = menanamkan "bibit" data awal agar aplikasi langsung bisa dipakai tanpa input manual.

```go
func SeedDB(db *gorm.DB) {
	// (1) Cek apakah tabel Role kosong
	var roleCount int64
	db.Model(&models.Role{}).Count(&roleCount)
	if roleCount == 0 {
		log.Println("Seeding roles...")
		roles := []models.Role{
			{Name: "superadmin"},
			{Name: "admin"},
			{Name: "user"},
		}
		for _, role := range roles {
			db.Create(&role)
		}
	}

	// (2) Cek apakah customer sudah ada
	var count int64
	db.Model(&models.Customer{}).Count(&count)
	if count > 0 {
		log.Println("Database already contains data, skipping seed.")
		return   // ← JANGAN seed ulang kalau sudah ada data
	}

	// (3) Seed customer, project, format contoh
	// ... (lihat file lengkap)
}
```

### Pola penting
1. **Cek dulu sebelum insert** (`Count`) → mencegah duplikat data saat aplikasi di-restart.
2. **`return` lebih awal** kalau sudah ada data → efisien.
3. Data contoh: 2 customer, 2 project, 2 format (`BAST/PO/...` & `BAST/INT/...`).

### Cara re-seed dari nol
Hapus file `bast_request.db`, lalu jalankan ulang server. Semua tabel & seed akan dibuat fresh.

---

## 7. Uji Coba

Setelah server berjalan, cek database. Cara termudah: gunakan [DB Browser for SQLite](https://sqlitebrowser.org/) untuk membuka `bast_request.db`. Anda akan melihat:
- 8 tabel sesuai model.
- Tabel `roles` berisi 3 baris (superadmin/admin/user).
- Tabel `master_customer` berisi 2 customer contoh.

Atau via API (setelah login):
```bash
curl -H "Authorization: Bearer <token>" http://localhost:8080/api/customers
```

---

## ✅ Ringkasan Step 2
- Model = struct + tag GORM → menjadi tabel.
- `TableName()` meng-override nama tabel default.
- Hook `BeforeCreate` mengisi UUID otomatis.
- `DeletedAt` mengaktifkan soft delete.
- `uniqueIndex` dengan nama yang sama = index komposit.
- `AutoMigrate` membuat/memperbarui tabel otomatis.
- `SeedDB` mengisi data awal **hanya jika tabel masih kosong**.

Database sudah siap berisi tabel. Tapi siapa pun bisa melihat/mengubah data? Tentu tidak — kita butuh autentikasi.

---

⬅️ **[Step 1: Setup & Konfigurasi](step-01-setup-and-config.md)** · ➡️ **[Step 3: Autentikasi JWT](step-03-authentication-jwt.md)**
