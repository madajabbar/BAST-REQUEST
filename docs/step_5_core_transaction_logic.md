# Step 5: Core Transaction Logic (Penomoran Otomatis BAST)

Sistem pembuatan nomor BAST secara berurutan (*Running Number*) adalah jantung dari sistem ini. Jika di sebuah rumah sakit mesin tiket antrian dipencet secara bersamaan oleh 2 orang berbeda, mereka tidak boleh mendapatkan angka "001" berbarengan. Begitupun dengan BAST.

Bagian ini membahas isi dari `internal/services/bast_request_service.go` dan hubungannya dengan `BastSequence` dan `BastFormat`.

## 1. Pengecekan Format (`BastFormat`)
Setiap meminta nomor, sistem membutuhkan ID Format.
Contoh sebuah format: `BAST/PO/{YYYY}/{MM}/{SEQ}`.
*Service* bertugas mengambil format ini dari Database, dan mengubah kata pengganti (*placeholder*):
- `{YYYY}` menjadi `2026`
- `{MM}` menjadi `06`
- `{SEQ}` menjadi nomor urut (*running number*) terbaru.

## 2. Pengambilan Sequence (`BastSequence`) secara Aman (Atomic)
Untuk mendapatkan nomor urut, kita tidak bisa sekedar melakukan pencarian jumlah dokumen BAST di database, karena rawan angka kembar (*Race Condition*).

**Logika Bisnis di `bast_sequence_service.go`**:
1. Mulai *Database Transaction* (Kunci tabel sementara).
2. Cari record `BastSequence` berdasarkan ID Format, Tahun Ini, dan Bulan Ini.
3. Jika TIDAK KETEMU: Artinya ini adalah BAST pertama di bulan ini. Buat row sequence baru dengan `last_number = 1`.
4. Jika KETEMU: Tambahkan `last_number = last_number + 1`.
5. Akhiri Transaksi dan kembalikan angkanya.

Proses ini wajib menggunakan transaksi `gorm` (`tx.Begin()`, `tx.Commit()`, `tx.Rollback()`) sehingga dua orang yang meng-klik tombol "Buat BAST" secara bersamaan dijamin 100% akan mendapat angka yang berbeda.

## 3. Pembuatan String Dokumen BAST
Setelah *Service* mendapatkan format dan angka urut (contoh: angka = 15). Service menggunakan library bawaan Go `strings.Replace` dan `fmt.Sprintf`:

```go
// bast_request_service.go (Ilustrasi Logika)

// 1. Ambil format string
pattern := "BAST/PO/{YYYY}/{MM}/{SEQ}"

// 2. Ganti YYYY dan MM
pattern = strings.Replace(pattern, "{YYYY}", "2026", -1)
pattern = strings.Replace(pattern, "{MM}", "06", -1)

// 3. Ubah angka 15 menjadi 3 digit "015" menggunakan fmt.Sprintf("%03d")
seqString := fmt.Sprintf("%03d", 15)

// 4. Ganti SEQ
hasilAkhir := strings.Replace(pattern, "{SEQ}", seqString, -1)

// Hasilnya: "BAST/PO/2026/06/015"
```

## 4. Pencatatan Jejak (Audit Log)
Setelah `BastRequest` tersimpan, *Handler/Service* tidak berhenti di situ. GORM memiliki fitur Trigger (Hooks).
Di dalam `internal/models/bast_request.go` (atau di Service), setiap perubahan harus memicu penyimpanan log. Saat terjadi UPDATE (misal status diubah dari `Active` menjadi `Void`), kode akan menangkap format JSON sebelum dan sesudah diganti, lalu menaruhnya ke dalam tabel `AuditLog`. Ini memastikan sistem memiliki sejarah (history) yang transparan.

## Ringkasan Step 5:
Tabel transaksional utama (`bast_request`) sangat bergantung pada kerja sama yang solid dengan tabel pendukung (`bast_format` dan `bast_sequence`) dalam satu payung **Transaction Database** agar urutan nomor bisa terjaga konsistensinya di setiap bulannya.

Lanjut ke: **[Step 6: Swagger Documentation](step_6_swagger_documentation.md)**
