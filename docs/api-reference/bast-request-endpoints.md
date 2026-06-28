# BAST Request Endpoints

Endpoint transaksi inti: membuat permintaan BAST (otomatis generate nomor), mengubah status, dan melihat audit. Sumber kode: [`bast_request_handler.go`](../../internal/handlers/bast_request_handler.go), [`bast_request_service.go`](../../internal/services/bast_request_service.go), [`bast_request_repository.go`](../../internal/repositories/bast_request_repository.go).

> 🔑 **Terlindungi** — semua endpoint butuh token JWT.

---

## Skema: BastRequest

```json
{
  "bast_request_id": "breq-123-uuid",
  "customer_id": "cust-...",
  "customer": { ... },
  "project_id": "proj-...",
  "project": { ... },
  "format_id": "fmt-...",
  "format": { ... },
  "perihal": "Serah terima server produksi",
  "tipe_nomor": "Internal",
  "po_number": "",
  "bast_number": "BAST/INT/2026/06/0001",
  "status": "Active",
  "requested_by": "Mada",
  "requested_at": "2026-06-28T10:30:00Z",
  "created_at": "...",
  "updated_at": "...",
  "deleted_at": null
}
```

| Field | Tipe | Keterangan |
|---|---|---|
| `bast_request_id` | UUID | Primary key |
| `customer_id` | UUID | FK → Customer |
| `customer` | object | Nested data customer (hasil Preload, saat GET) |
| `project_id` | UUID | FK → Project |
| `project` | object | Nested data project (hasil Preload, saat GET) |
| `format_id` | UUID | FK → BastFormat |
| `format` | object | Nested data format (hasil Preload, saat GET) |
| `perihal` | string | Deskripsi permintaan |
| `tipe_nomor` | string | `Internal` (auto-generate) atau `PO` (pakai nomor PO) |
| `po_number` | string | Nomor PO dari customer (wajib jika `tipe_nomor: PO`) |
| `bast_number` | string | **Auto-generated** (Internal) atau dari `po_number` (PO). **Unique.** |
| `status` | string | `Active`, `Used`, `Void` |
| `requested_by` | string | Nama/identitas peminta |
| `requested_at` | datetime | Auto-filled saat create |

---

## Get All
<a id="get-all"></a>

### Request
```http
GET /api/bast-requests?customer_id=<uuid>&project_id=<uuid>&status=Active
Authorization: Bearer <token>
```

**Query param:**
| Param | Tipe | Wajib | Keterangan |
|---|---|---|---|
| `customer_id` | string | ❌ | Filter customer |
| `project_id` | string | ❌ | Filter project |
| `status` | string | ❌ | Filter status (`Active`/`Used`/`Void`) |

### Response — `200 OK`
Array of BastRequest (dengan nested customer, project, format).
```json
[
  {
    "bast_request_id": "...",
    "bast_number": "BAST/INT/2026/06/0001",
    "status": "Active",
    "perihal": "Serah terima server",
    "customer": { "customer_name": "PT. Maju Mundur", ... },
    "project": { "project_name": "ERP Terpadu", ... },
    "format": { "format_pattern": "BAST/INT/{YYYY}/{MM}/{SEQ}", ... },
    ...
  }
]
```

### Contoh curl
```bash
curl -H "Authorization: Bearer <TOKEN>" \
  "http://localhost:8080/api/bast-requests?status=Active"
```

---

## Get By ID
<a id="get-by-id"></a>

### Request
```http
GET /api/bast-requests/:id
Authorization: Bearer <token>
```

### Response — `200 OK`
Object BastRequest lengkap (nested).

### `404 Not Found`
```json
{ "error": "Request not found" }
```

---

## Create
<a id="create"></a>

Buat permintaan BAST baru. Jika `tipe_nomor: Internal`, nomor **di-generate otomatis**. Jika `tipe_nomor: PO`, nomor diambil dari `po_number`.

### Request
```http
POST /api/bast-requests
Authorization: Bearer <token>
Content-Type: application/json
```

**Body:**
```json
{
  "customer_id": "550e8400-customer-id",
  "project_id": "abc-123-project-id",
  "format_id": "fmt-internal-id",
  "perihal": "Serah terima server produksi",
  "tipe_nomor": "Internal",
  "requested_by": "Mada"
}
```

| Field | Tipe | Wajib | Keterangan |
|---|---|---|---|
| `customer_id` | UUID | ✅ | Customer yang meminta |
| `project_id` | UUID | ✅ | Project terkait |
| `format_id` | UUID | ✅ | Format penomoran yang dipakai |
| `perihal` | string | ✅ | Deskripsi |
| `tipe_nomor` | string | ✅ | `Internal` atau `PO` |
| `po_number` | string | kondisional | **Wajib** jika `tipe_nomor: PO` |
| `requested_by` | string | ✅ | Nama peminta |
| `bast_number` | string | ❌ | **Jangan kirim** — di-auto-generate |
| `status` | string | ❌ | **Jangan kirim** — default `Active` |

### Contoh: BAST Tipe PO
```json
{
  "customer_id": "550e8400-...",
  "project_id": "abc-123-...",
  "format_id": "fmt-po-id",
  "perihal": "Serah terima berdasarkan PO",
  "tipe_nomor": "PO",
  "po_number": "PO/2026/001",
  "requested_by": "Mada"
}
```

### Response — `201 Created`
Object BastRequest lengkap dengan `bast_number` terisi (baik auto-generate atau dari PO).

### `500 Internal Server Error`
```json
{ "error": "invalid format ID" }
{ "error": "po_number is required for TipeNomor PO" }
{ "error": "invalid tipe_nomor, must be Internal or PO" }
{ "error": "failed to generate running number: ..." }
```

### Contoh curl (Internal)
```bash
curl -X POST http://localhost:8080/api/bast-requests \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id":"550e8400-...",
    "project_id":"abc-123-...",
    "format_id":"fmt-...",
    "perihal":"Serah terima server",
    "tipe_nomor":"Internal",
    "requested_by":"Mada"
  }'
```

### Contoh curl (PO)
```bash
curl -X POST http://localhost:8080/api/bast-requests \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id":"550e8400-...",
    "project_id":"abc-123-...",
    "format_id":"fmt-po...",
    "perihal":"Serah terima PO",
    "tipe_nomor":"PO",
    "po_number":"PO/2026/001",
    "requested_by":"Mada"
  }'
```

---

## Update Status
<a id="update-status"></a>

Ubah status BAST: `Active` → `Used` (sudah dipakai) atau `Void` (dibatalkan).

### Request
```http
PATCH /api/bast-requests/:id/status
Authorization: Bearer <token>
Content-Type: application/json
```

**Body:**
```json
{
  "status": "Used"
}
```

**Nilai status valid:** `Active`, `Used`, `Void`.

### Response — `200 OK`
```json
{ "message": "Status updated successfully" }
```

### `500 Internal Server Error`
```json
{ "error": "invalid status" }
```

### Contoh curl
```bash
curl -X PATCH http://localhost:8080/api/bast-requests/<ID>/status \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"status":"Used"}'
```

---

## Get Request Audit
<a id="audit"></a>

Lihat riwayat perubahan untuk satu request BAST tertentu.

> ⚠️ **Status saat ini:** endpoint ini mengembalikan pesan placeholder. Audit per-request belum terintegrasi penuh (lihat `// TODO` di kode).

### Request
```http
GET /api/bast-requests/:id/audit
Authorization: Bearer <token>
```

### Response — `200 OK`
```json
{ "message": "Audit logs for request <id>" }
```

> 💡 Saat ini mengembalikan pesan string. Di masa depan, seharusnya mengembalikan array audit log yang terkait request ID tersebut.

---

## Alur Tipikal Penggunaan

1. Pastikan sudah ada **customer** dan **project**.
2. Buat atau pilih **format** yang sesuai.
3. `POST /bast-requests` dengan data lengkap.
4. Sistem generate nomor otomatis → cek `bast_number` di respons.
5. Jika BAST sudah terpakai: `PATCH /bast-requests/:id/status` → `Used`.

---

← **[Daftar Referensi API](README.md)** · 🤿 **[Deep Dive Penomoran](../guides/bast-numbering-deep-dive.md)**
