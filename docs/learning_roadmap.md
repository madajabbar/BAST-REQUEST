# Peta Belajar (Learning Roadmap) Golang Backend Secara Mendalam

Membangun aplikasi *Backend* menggunakan Golang pada awalnya mungkin terasa membingungkan, terutama karena Golang memiliki gaya pemrograman yang berbeda dibandingkan bahasa seperti PHP, Python, atau Java. 

Dokumen ini ditulis **khusus untuk Anda** sebagai panduan komprehensif. Bacalah perlahan-lahan dari atas ke bawah.

---

## 1. Fondasi Utama: Memahami Konsep Unik di Golang

Sebelum kita berbicara tentang *framework* atau database, Anda wajib memahami aturan main dasar di Golang.

### A. Tidak Ada Class, Adanya "Struct"
Jika Anda pernah belajar OOP (Object-Oriented Programming), Anda pasti tahu tentang `Class`. Di Golang, konsep tersebut **tidak ada**. Sebagai gantinya, Golang menggunakan `Struct`. `Struct` hanyalah sebuah wadah (cetakan) untuk mengelompokkan beberapa variabel menjadi satu kesatuan.

**Contoh:**
```go
// Mendefinisikan cetakan bernama 'Karyawan'
type Karyawan struct {
    NamaLengkap string
    Gaji        int
    Aktif       bool
}

func main() {
    // Membuat data karyawan baru menggunakan cetakan di atas
    pegawai1 := Karyawan{
        NamaLengkap: "Budi Santoso",
        Gaji:        5000000,
        Aktif:       true,
    }
}
```

Lalu bagaimana cara kita menambahkan *behavior* atau aksi ke dalam *Struct* tersebut? Kita menggunakan **Method**.
```go
// Fungsi ini hanya menempel (dimiliki) oleh struct Karyawan
func (k Karyawan) HitungPajak() int {
    return k.Gaji * 10 / 100
}
```

### B. Konsep "Pointer" (Simbol `*` dan `&`)
Ini adalah topik yang paling sering membuat pemula bingung. Bayangkan variabel sebagai sebuah **rumah**.
- Variabel biasa (misal: `karyawan1`) adalah isi dari rumah tersebut (TV, Meja, Kursi).
- **Pointer** (ditandai dengan `*`) adalah **Alamat** dari rumah tersebut (Jl. Merdeka No.1).

**Kenapa Pointer ini sangat penting?**
Secara bawaan, saat Anda melempar sebuah variabel ke dalam fungsi lain, Golang akan melakukan **"Copy-Paste"** (Salin). Jika rumahnya besar, *copy-paste* akan memakan banyak sekali RAM memori, dan jika rumah salinan itu diubah, rumah aslinya tidak akan berubah!

Dengan memberikan **Pointer** (alamatnya saja), fungsi tersebut dapat pergi ke "rumah asli" dan mengubah nilainya tanpa harus menguras RAM untuk melakukan *copy-paste*.

```go
// Parameter (k *Karyawan) menggunakan tanda bintang, artinya kita meminta ALAMATNYA
func NaikkanGaji(k *Karyawan) {
    k.Gaji = k.Gaji + 1000000 // Mengubah data asli
}

func main() {
    pegawai := Karyawan{Gaji: 5000000}
    
    // Tanda & (ampersand) digunakan untuk MENDAPATKAN alamat dari variabel pegawai
    NaikkanGaji(&pegawai) 
    
    // Gaji sekarang menjadi 6 Juta!
}
```
*Dalam kode proyek ini, hampir seluruh `Handler`, `Service`, dan `Repository` kita saling berkomunikasi menggunakan Pointer agar aplikasi sangat cepat dan hemat memori.*

### C. Penanganan Error (No Try-Catch)
Golang tidak menggunakan blok `try { ... } catch (e) { ... }`. Golang memperlakukan error sebagai nilai biasa (seperti string atau int) yang dikembalikan oleh fungsi. Ini memaksa pemrogram untuk sadar dan menangani error secara langsung.

```go
// Fungsi mengembalikan DUA nilai: Hasil (int) dan Error (error)
func Bagi(a int, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("tidak bisa dibagi dengan nol!") // Kirim error
    }
    return a / b, nil // nil berarti "kosong" (tidak ada error)
}

func main() {
    hasil, err := Bagi(10, 0)
    if err != nil {
        fmt.Println("Waduh ada masalah:", err)
        return // Hentikan proses
    }
}
```

---

## 2. GORM: Jembatan Antara Golang dan Database SQL

Bayangkan Anda memiliki database SQLite atau MySQL. Biasanya, Anda harus menulis perintah SQL murni seperti:
`SELECT * FROM users WHERE age > 20;`

Hal ini cukup merepotkan, rentan typo, dan sulit dirawat. Maka dari itu kita menggunakan **ORM (Object Relational Mapping)** bernama **GORM**. GORM mengubah `Struct` Go Anda secara otomatis menjadi tabel SQL, dan mempermudah query.

### A. AutoMigrate (Membuat Tabel Otomatis)
Anda hanya perlu membuat Struct, lalu memanggil fungsi AutoMigrate. GORM akan mengecek: *"Apakah tabel ini sudah ada? Kalau belum, buatkan!"*
```go
type Buku struct {
    ID     uint   `gorm:"primaryKey"` // Memberitahu DB bahwa ini adalah Primary Key
    Judul  string `gorm:"type:varchar(100);not null"`
    Penulis string
}

// Di file main.go:
db.AutoMigrate(&Buku{}) // Otomatis membuat tabel 'bukus' di SQLite!
```

### B. Operasi CRUD (Create, Read, Update, Delete)
Dengan GORM, operasi database menjadi fungsi-fungsi Golang biasa:
```go
// CREATE (INSERT)
bukuBaru := Buku{Judul: "Belajar Golang", Penulis: "Mada"}
db.Create(&bukuBaru) // Menyimpan ke DB

// READ (SELECT)
var bukuCari Buku
db.First(&bukuCari, 1) // Mencari buku dengan ID 1
db.Where("penulis = ?", "Mada").Find(&kumpulanBuku) // Mencari banyak buku

// UPDATE
db.Model(&bukuCari).Update("Judul", "Mastering Golang")

// DELETE
db.Delete(&bukuCari, 1)
```

---

## 3. Gin Web Framework (Lalu Lintas API Anda)

Golang memiliki server web bawaan, tapi terlalu dasar. **Gin** adalah sebuah kerangka kerja (*framework*) yang membuat pembuatan REST API menjadi sangat mudah dan sangat cepat.

Fungsi utama Gin ada tiga:
1. **Routing:** Menentukan URL mana (`/api/login`, `/api/users`) yang memicu fungsi mana.
2. **Binding Data:** Menerjemahkan data JSON yang dikirim dari klien (Postman/Web) menjadi `Struct` Golang.
3. **Response:** Membungkus hasil kembali menjadi format JSON.

```go
func main() {
    r := gin.Default() // Membuat server Gin
    
    // Jika ada yang mengakses URL GET /ping, jalankan fungsi ini
    r.GET("/ping", func(c *gin.Context) {
        
        // c (Context) adalah alat serbaguna. Di sini dipakai untuk mengirim JSON balasan (Response)
        c.JSON(200, gin.H{
            "pesan": "pong!",
        })
        
    })
    
    r.Run(":8080") // Nyalakan server di port 8080
}
```

---

## 4. Filosofi "Clean Architecture" (Arsitektur Bersih)

Saat proyek API semakin besar, meletakkan semua kode di satu file `main.go` akan menjadi bencana. Jika ada error, Anda akan pusing mencarinya. 

Oleh karena itu, proyek BAST Request ini memecah kode ke dalam **4 Lapisan (Layer)** yang sangat ketat tugasnya. Mari kita gunakan **Analogi Restoran Cepat Saji**:

1. **Handler (Kasir / Pelayan)** 
   Tugasnya hanya satu: Menerima pesanan (Request) dari pelanggan (Client), mengecek apakah pesanannya masuk akal, lalu menyerahkan catatan pesanan tersebut ke Dapur. Kasir tidak memasak makanan.
   *Di kode:* Handler menerima URL dan JSON, lalu memanggil fungsi di *Service*.

2. **Service (Koki / Dapur - Logika Bisnis Utama)** 
   Tugasnya memasak pesanan sesuai resep. Ia memastikan urutannya benar, bumbunya pas (Validasi Bisnis). Tapi Koki tidak mengambil bahan baku sendiri ke pasar. Ia menyuruh asistennya.
   *Di kode:* Service berisi `if-else` kompleks, perhitungan angka `BAST Sequence`, dan logika keamanan. Jika butuh data, ia memanggil *Repository*.

3. **Repository (Asisten Gudang)**
   Tugasnya murni hanya masuk ke ruang penyimpanan (Database), mengambil bahan baku (Data SQL), atau meletakkan barang baru, lalu memberikannya kembali ke Koki. Asisten ini sama sekali tidak tahu masakan apa yang sedang dibuat.
   *Di kode:* Di sinilah semua sintaks `gorm` (seperti `db.Create`, `db.Find`) harus diletakkan. Dilarang keras meletakkan `gorm` di Handler atau Service.

4. **Models (Bahan Baku)**
   Adalah bentuk fisik dari apa yang diproses. Wortel, Daging, Kentang.
   *Di kode:* Inilah definisi struktur tabel database (`Struct`).

**Kenapa harus serepot ini?**
Jika esok hari Anda ingin mengganti SQLite menjadi PostgreSQL, Anda **HANYA** perlu mengotak-atik folder `Repository`. Folder `Service` dan `Handler` sama sekali tidak perlu disentuh karena mereka tidak peduli database apa yang Anda gunakan, asalkan Asisten Gudang (Repository) memberikan datanya! Ini membuat aplikasi Anda sangat kebal terhadap perubahan.

---

## 5. Autentikasi Menggunakan JWT (JSON Web Tokens)

Aplikasi BAST ini memiliki sistem *Login*. Mengapa kita butuh JWT? 
Berbeda dengan web konvensional (PHP) yang menyimpan "Sesi Seseorang" (Session) di dalam RAM Server, REST API bersifat *Stateless* (Pelupa). Server tidak pernah mengingat siapa klien yang datang barusan.

**Alur Kerjanya:**
1. **Login:** Klien (Postman/Frontend) mengirimkan Email & Password.
2. **Validasi:** Server mengecek ke Database. Jika benar, server membuat sebuah **Karcis Masuk (Token)** bernama JWT.
3. **Penyimpanan:** JWT ini berisi data terenkripsi (misal: "Orang ini bernama Mada, jabatannya Admin"). Server memberikan karcis ini ke klien, dan klien menyimpannya sendiri. Server tidak menyimpan karcisnya!
4. **Akses Data:** Besoknya, saat klien ingin melihat daftar *Project*, klien harus membawa dan menempelkan karcis JWT tersebut di bagian *Headers* request (`Authorization: Bearer <Karcis>`).
5. **Verifikasi Tengah Jalan (Middleware):** Sebelum pesanan klien sampai ke *Handler*, Gin akan menugaskan satpam (Middleware) di pintu depan. Middleware akan membaca Karcis JWT. Jika palsu atau kadaluarsa, Satpam menendang klien (Error 401 Unauthorized). Jika asli, klien dipersilakan masuk ke *Handler*.

---

## Ringkasan Saran Langkah Praktek Anda:
Agar dokumentasi ini tidak sekadar menjadi teori, cobalah langkah berikut di komputer Anda:

1. Buat folder kosong baru. Jalankan `go mod init cobago`.
2. Install Gin: `go get -u github.com/gin-gonic/gin`.
3. Buat file `main.go` dan cobalah membuat satu *Endpoint* sederhana yang mem-print "Hello World".
4. Jika sudah lancar, buatlah `Struct` sederhana (misalnya Barang), dan coba hubungkan dengan GORM.
5. Cobalah membuat 3 file terpisah untuk mempraktikkan konsep *Repository*, *Service*, dan *Handler* yang telah saya jelaskan di atas.
