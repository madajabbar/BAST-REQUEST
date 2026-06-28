# Step 6: Audit Log (Pelacakan Jejak Perubahan)

> Seri Tutorial · **Step 6 dari 8**

Audit log adalah sistem yang mencatat **setiap perubahan penting** di database: apa yang berubah, kapan, oleh siapa, dan bagaimana kondisi data sebelum & sesudahnya. Ini krusial untuk transparansi, kepatuhan (*compliance*), dan investigasi masalah.

---

## 1. Mengapa Butuh Audit Log?

Tanpa audit log, jika ada data Customer mendadak hilang/berubah, kita **tidak tahu** siapa pelakunya. Audit log menjawab:
- 🕵️ **Siapa** yang melakukan aksi?
- 📅 **Kapan** dilakukan?
- 🗂️ Di **tabel & record mana**?
- 🔄 Data **sebelum** dan **sesudah** seperti apa?

---

## 2. Model AuditLog

File: [`internal/models/audit_log.go`](../../internal/models/audit_log.go)

```go
type AuditLog struct {
	AuditLogID  uuid.UUID      `gorm:"type:uuid;primary_key"`
	TargetTable string         `gorm:"type:varchar(100);not null;column:table_name"`
	RecordID    string         `gorm:"type:varchar(100);not null"`
	Action      string         `gorm:"type:varchar(20);not null"`  // POST/PUT/DELETE/PATCH
	OldData     datatypes.JSON `gorm:"type:jsonb"`
	NewData     datatypes.JSON `gorm:"type:jsonb"`
	PerformedBy string         `gorm:"type:varchar(100);not null"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
}

func (AuditLog) TableName() string {
	return "audit_log"
}
```

### Penjelasan field

| Field | Tipe | Isi |
|---|---|---|
| `TargetTable` | string | Nama tabel yang diubah. **Catatan:** nama field Go `TargetTable`, tapi nama kolom DB di-override jadi `table_name` via tag `column:table_name`. |
| `RecordID` | string | ID record (UUID) yang diubah |
| `Action` | string | Jenis operasi: `POST` (create), `PUT`/`PATCH` (update), `DELETE` |
| `OldData` | JSON | Snapshot data **sebelum** diubah |
| `NewData` | JSON | Snapshot data **sesudah** diubah |
| `PerformedBy` | string | Siapa pelaku (bisa userID/username) |

### Tipe `datatypes.JSON`
Pakai `gorm.io/datatypes`. Di SQLite tersimpan sebagai teks JSON; di PostgreSQL sebagai `jsonb` (bisa di-query). Ini fleksibel — bisa menyimpan struktur data apa pun tanpa perlu tabel terpisah.

---

## 3. Repository — Filter Fleksibel

File: [`internal/repositories/audit_log_repository.go`](../../internal/repositories/audit_log_repository.go)

```go
func (r *AuditLogRepository) FindAll(targetTable, recordID, performedBy, dateFrom, dateTo string) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	query := r.db.Model(&models.AuditLog{})

	if targetTable != "" {
		query = query.Where("table_name = ?", targetTable)
	}
	if recordID != "" {
		query = query.Where("record_id = ?", recordID)
	}
	if performedBy != "" {
		query = query.Where("performed_by = ?", performedBy)
	}
	if dateFrom != "" {
		query = query.Where("created_at >= ?", dateFrom)
	}
	if dateTo != "" {
		query = query.Where("created_at <= ?", dateTo)
	}

	err := query.Order("created_at DESC").Find(&logs).Error
	return logs, err
}
```

### Pola filter dinamis (lagi)
Sama seperti di Customer (Step 4): `Where` ditambahkan hanya jika parameter tidak kosong. Bisa kombinasi:
- Filter berdasarkan **satu tabel** saja: `?table_name=master_customer`
- Filter **rentang waktu**: `?date_from=2026-06-01&date_to=2026-06-30`
- Filter **pelaku**: `?performed_by=admin1`

### Urutan terbaru dulu
`Order("created_at DESC")` → log terbaru muncul di atas. Wajar untuk audit: kita ingin lihat aktivitas terkini.

---

## 4. Service & Handler

### Service
File: [`internal/services/audit_log_service.go`](../../internal/services/audit_log_service.go)
```go
func (s *AuditLogService) GetAllAuditLogs(tableName, recordID, performedBy, dateFrom, dateTo string) ([]models.AuditLog, error) {
	return s.repo.FindAll(tableName, recordID, performedBy, dateFrom, dateTo)
}

func (s *AuditLogService) LogAction(log *models.AuditLog) error {
	return s.repo.Create(log)
}
```

Method `LogAction` adalah **API publik** untuk layer lain mencatat audit. Misalnya, di masa depan service BAST bisa panggil ini setelah `UpdateStatus`.

### Handler
File: [`internal/handlers/audit_log_handler.go`](../../internal/handlers/audit_log_handler.go)
```go
func (h *AuditLogHandler) GetAllAuditLogs(c *gin.Context) {
	tableName := c.Query("table_name")
	recordID := c.Query("record_id")
	performedBy := c.Query("performed_by")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	logs, err := h.service.GetAllAuditLogs(tableName, recordID, performedBy, dateFrom, dateTo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, logs)
}
```
Membaca 5 query param → teruskan ke service → balas JSON.

---

## 5. Catatan Status Implementasi Saat Ini

> ⚠️ **Transparansi penting:** Saat kode ini ditulis, pencatatan audit **belum otomatis** di-hook ke setiap operasi CRUD. Method `LogAction` sudah ada, tapi belum dipanggil dari service lain secara sistematis.

Artinya:
- Tabel `audit_log` sudah **siap** & ter-migrate.
- Endpoint GET `/api/audit-logs` berfungsi untuk **membaca** (kalau ada isinya).
- Tapi belum ada kode yang **menulis** audit log otomatis setiap create/update/delete.

### Bagaimana mengaktifkannya (ide pengembangan)

**Opsi 1: Panggil manual di Service**
```go
func (s *CustomerService) UpdateCustomer(id string, input *models.Customer) (models.Customer, error) {
	old, _ := s.repo.FindByID(id)
	// ... lakukan update ...

	// Catat audit
	s.auditService.LogAction(&models.AuditLog{
		TargetTable: "master_customer",
		RecordID:    id,
		Action:      "PUT",
		OldData:     toJSON(old),
		NewData:     toJSON(customer),
		PerformedBy: "user-xxx",
	})
	return customer, nil
}
```

**Opsi 2: GORM Callbacks/Hooks (lebih elegan)**
Buat callback global GORM yang otomatis trigger setiap BeforeUpdate/AfterCreate, lalu inject audit log. Ini lebih "set & forget" tapi lebih kompleks.

**Opsi 3: Middleware-level logging**
Catat di middleware berdasarkan method HTTP (POST/PUT/DELETE).

> 📝 Untuk sekarang, fokus kita memahami **struktur** audit log. Implementasi otomatis bisa jadi tugas pengembangan berikutnya.

---

## 6. Contoh Data Audit Log (Konseptual)

Jika sudah aktif, satu baris audit log akan terlihat seperti:
```json
{
  "audit_log_id": "550e8400-e29b-41d4-a716-446655440000",
  "table_name": "master_customer",
  "record_id": "abc-123-customer-id",
  "action": "PUT",
  "old_data": {
    "customer_name": "PT. Maju Mundur",
    "status": "active"
  },
  "new_data": {
    "customer_name": "PT. Maju Mundur Jaya",
    "status": "active"
  },
  "performed_by": "user-xyz",
  "created_at": "2026-06-28T10:30:00Z"
}
```

Dari sini kita tahu: user-xyz mengubah nama customer pada 28 Juni 2026 pukul 10:30.

---

## 7. Uji Coba

### Ambil semua audit log
```bash
curl -H "Authorization: Bearer <TOKEN>" http://localhost:8080/api/audit-logs
```

### Filter by tabel
```bash
curl -H "Authorization: Bearer <TOKEN>" "http://localhost:8080/api/audit-logs?table_name=bast_request"
```

### Filter rentang tanggal
```bash
curl -H "Authorization: Bearer <TOKEN>" "http://localhost:8080/api/audit-logs?date_from=2026-06-01&date_to=2026-06-30"
```

> ⚠️ Karena audit belum ditulis otomatis (lihat bagian 5), response kemungkinan array kosong `[]` kecuali Anda isi manual untuk eksperimen.

---

## ✅ Ringkasan Step 6
- Audit log mencatat: tabel, record, aksi, data lama, data baru, pelaku, waktu.
- Field `OldData`/`NewData` bertipe `datatypes.JSON` — fleksibel & queryable (di PostgreSQL).
- Repository mendukung filter dinamis (table, record, performer, rentang tanggal).
- **Status saat ini:** struktur siap, tapi penulisan audit belum otomatis — butuh integrasi ke service atau GORM callbacks.

Audit log sudah dipahami. Sekarang mari lihat bagaimana **semua endpoint disusun & diamankan** dalam satu file routing.

---

⬅️ **[Step 5: Mesin Penomoran BAST](step-05-bast-numbering-engine.md)** · ➡️ **[Step 7: Routing & RBAC](step-07-routing-and-rbac.md)**
