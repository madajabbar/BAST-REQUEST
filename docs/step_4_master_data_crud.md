# Step 4: Operasi Master Data (CRUD)

Kita sudah memiliki koneksi database (Step 1), pembuatan tabel otomatis (Step 2), dan pengamanan pintu depan via Token JWT (Step 3). Sekarang, mari kita lihat bagaimana aplikasi memproses pembuatan, pembacaan, pengubahan, dan penghapusan (CRUD) terhadap tabel Master seperti **Customer**, **Project**, dan **Bast Format**.

Konsep pengerjaan di sini murni mendemonstrasikan pola kerja berantai: `Handler -> Service -> Repository`.

## 1. Repository: Mengobrol dengan Database
Misalkan kita ingin melihat file `internal/repositories/customer_repository.go`. Di sinilah GORM digunakan.

```go
func (r *CustomerRepository) Create(customer *models.Customer) error {
	// Menyimpan data Customer baru
	return r.db.Create(customer).Error
}

func (r *CustomerRepository) FindAll() ([]models.Customer, error) {
	var customers []models.Customer
	// Melakukan SELECT * FROM customers
	err := r.db.Find(&customers).Error
	return customers, err
}

func (r *CustomerRepository) Update(customer *models.Customer) error {
	// Meng-update data berdasarkan Primary Key dari variabel 'customer'
	return r.db.Save(customer).Error
}
```
*Tugas Repositori hanyalah menerjemahkan bahasa Go menjadi perintah SQL.*

## 2. Service: Mandor Logika Bisnis
Lihat file `internal/services/customer_service.go`. Ia menerima data dari Handler, memvalidasinya, lalu menyuruh Repository untuk mengeksekusi SQL.

```go
func (s *CustomerService) CreateCustomer(customer *models.Customer) error {
	// 1. Logika Bisnis (Contoh: Code Customer tidak boleh kosong)
	if customer.CustomerCode == "" {
		return errors.New("kode customer wajib diisi")
	}

	// 2. Suruh repository mengeksekusinya
	return s.repo.Create(customer)
}
```

## 3. Handler: Membaca Internet JSON (Gin)
Lihat `internal/handlers/customer_handler.go`. Ia tidak peduli aturan bisnis atau tipe SQL yang dipakai, ia hanya peduli untuk menerjemahkan HTTP JSON.

```go
// CreateCustomer godoc
// @Summary Create a new customer
func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	// 1. Buat 'cetakan' / struct kosong
	var input models.Customer
	
	// 2. Baca body JSON dan "Tuangkan" nilainya ke cetakan 'input'
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data JSON salah"})
		return
	}

	// 3. Lemparkan ke Service (Mandor)
	if err := h.service.CreateCustomer(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 4. Sukses! Balas dengan status 201 Created
	c.JSON(http.StatusCreated, input)
}
```

## Ringkasan Step 4:
Pembuatan CRUD di Clean Architecture Golang selalu berulang melalui 3 siklus file tersebut. Proses ini dijamin membuat kode lebih bersih, karena jika tiba-tiba kita diharuskan merubah validasi pembuatan Customer, kita tahu persis kita HANYA perlu masuk ke file `Service` tanpa mengganggu bagian Database atau HTTP JSON.

Lanjut ke: **[Step 5: Core Transaction Logic (Penomoran BAST)](step_5_core_transaction_logic.md)**
