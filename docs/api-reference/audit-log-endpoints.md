# Audit Log Endpoints

Endpoint untuk membaca log audit (jejak perubahan data). Sumber kode: [`audit_log_handler.go`](../../internal/handlers/audit_log_handler.go), [`audit_log_service.go`](../../internal/services/audit_log_service.go), [`audit_log_repository.go`](../../internal/repositories/audit_log_repository.go).

> 🔑 **Terlindungi** — semua endpoint butuh token JWT.
>
> ⚠️ **Catatan:** Sistem audit log sudah memiliki struktur dan endpoint lengkap, namun **penulisan log belum otomatis** di setiap operasi CRUD. Endpoint ini akan mengembalikan array kosong `[]` sampai audit logging diintegrasikan ke service layer. Lihat [Tutorial Step 6](../tutorials/step-06-audit-log.md) untuk detail.

---

## Skema: AuditLog

```json
{
  "audit_log_id": "audit-123-uuid",
  "table_name": "master_customer",
  "record_id": "cust-456-uuid",
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

| Field | Tipe | Keterangan |
|---|---|---|
| `audit_log_id` | UUID | Primary key |
| `table_name` | string | Nama tabel yang diubah (mis. `master_customer`, `bast_request`) |
| `record_id` | string | ID record yang diubah |
| `action` | string | Jenis aksi: `POST`, `PUT`, `DELETE`, `PATCH` |
| `old_data` | JSON | Snapshot data **sebelum** diubah (`null` untuk `POST`) |
| `new_data` | JSON | Snapshot data **sesudah** diubah |
| `performed_by` | string | Siapa yang melakukan aksi |
| `created_at` | datetime | Waktu aksi dilakukan |

> 💡 `old_data` dan `new_data` bertipe JSON — bisa berisi struktur data apa pun sesuai tabel sumber.

---

## Get All Audit Logs
<a id="get-all"></a>

Ambil daftar audit log dengan filter opsional. Hasil diurutkan **terbaru dulu** (`created_at DESC`).

### Request
```http
GET /api/audit-logs?table_name=master_customer&date_from=2026-06-01&date_to=2026-06-30
Authorization: Bearer <token>
```

**Query param:**
| Param | Tipe | Wajib | Keterangan |
|---|---|---|---|
| `table_name` | string | ❌ | Filter nama tabel (mis. `bast_request`, `master_customer`) |
| `record_id` | string | ❌ | Filter record ID tertentu |
| `performed_by` | string | ❌ | Filter pelaku (userID/username) |
| `date_from` | string | ❌ | Awal rentang waktu (`YYYY-MM-DD`) |
| `date_to` | string | ❌ | Akhir rentang waktu (`YYYY-MM-DD`) |

### Contoh Kombinasi Filter
| URL | Arti |
|---|---|
| `/audit-logs` | Semua log |
| `/audit-logs?table_name=bast_request` | Semua perubahan BAST request |
| `/audit-logs?performed_by=admin1` | Semua aksi oleh user admin1 |
| `/audit-logs?date_from=2026-06-01&date_to=2026-06-30` | Aktivitas bulan Juni 2026 |
| `/audit-logs?record_id=abc-123` | Semua perubahan pada record tertentu |
| `/audit-logs?table_name=bast_request&performed_by=admin1&date_from=2026-06-15` | Kombinasi 3 filter |

### Response — `200 OK`
```json
[
  {
    "audit_log_id": "...",
    "table_name": "master_customer",
    "record_id": "cust-...",
    "action": "PUT",
    "old_data": { ... },
    "new_data": { ... },
    "performed_by": "admin1",
    "created_at": "2026-06-28T10:30:00Z"
  },
  {
    "audit_log_id": "...",
    "table_name": "bast_request",
    "record_id": "breq-...",
    "action": "PATCH",
    "old_data": { "status": "Active" },
    "new_data": { "status": "Used" },
    "performed_by": "user1",
    "created_at": "2026-06-27T14:00:00Z"
  }
]
```

### Contoh curl
```bash
# Semua log
curl -H "Authorization: Bearer <TOKEN>" http://localhost:8080/api/audit-logs

# Filter per tabel
curl -H "Authorization: Bearer <TOKEN>" "http://localhost:8080/api/audit-logs?table_name=bast_request"

# Filter rentang tanggal
curl -H "Authorization: Bearer <TOKEN>" \
  "http://localhost:8080/api/audit-logs?date_from=2026-06-01&date_to=2026-06-30"

# Kombinasi
curl -H "Authorization: Bearer <TOKEN>" \
  "http://localhost:8080/api/audit-logs?table_name=master_customer&performed_by=admin1"
```

---

## Get Audit Log By ID
<a id="get-by-id"></a>

### Request
```http
GET /api/audit-logs/:id
Authorization: Bearer <token>
```

| Param | Lokasi | Keterangan |
|---|---|---|
| `id` | path | UUID audit log |

### Response — `200 OK`
Object audit log tunggal.

### `404 Not Found`
```json
{ "error": "Audit log not found" }
```

### Contoh curl
```bash
curl -H "Authorization: Bearer <TOKEN>" \
  http://localhost:8080/api/audit-logs/<AUDIT-LOG-ID>
```

---

## Contoh Use Case Audit Log

### Investigasi: "Siapa mengubah customer X?"
```bash
curl -H "Authorization: Bearer <TOKEN>" \
  "http://localhost:8080/api/audit-logs?table_name=master_customer&record_id=<CUSTOMER_ID>"
```

### Monitoring: "Aktivitas hari ini"
```bash
curl -H "Authorization: Bearer <TOKEN>" \
  "http://localhost:8080/api/audit-logs?date_from=2026-06-28&date_to=2026-06-28"
```

### Forensik: "Semua aksi user Y"
```bash
curl -H "Authorization: Bearer <TOKEN>" \
  "http://localhost:8080/api/audit-logs?performed_by=user-y-id"
```

---

← **[Daftar Referensi API](README.md)** · 🎓 **[Tutorial Step 6 — Audit Log](../tutorials/step-06-audit-log.md)**
