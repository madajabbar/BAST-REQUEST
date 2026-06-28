# BAST Format Endpoints

CRUD master format penomoran BAST. Format menentukan pola (pattern) nomor yang akan dipakai saat membuat BAST Request. Sumber kode: [`bast_format_handler.go`](../../internal/handlers/bast_format_handler.go), [`bast_format_service.go`](../../internal/services/bast_format_service.go), [`bast_format_repository.go`](../../internal/repositories/bast_format_repository.go).

> 🔑 **Terlindungi** — semua endpoint butuh token.

---

## Skema: BastFormat

```json
{
  "format_id": "fmt-123-uuid",
  "format_name": "Format Internal Perusahaan",
  "format_type": "Internal",
  "format_pattern": "BAST/INT/{YYYY}/{MM}/{SEQ}",
  "is_active": true,
  "created_at": "...",
  "updated_at": "...",
  "deleted_at": null
}
```

| Field | Tipe | Keterangan |
|---|---|---|
| `format_id` | UUID | Primary key |
| `format_name` | string | Nama deskriptif |
| `format_type` | string | `PO` atau `Internal` |
| `format_pattern` | string | Pola nomor dengan placeholder `{YYYY}`/`{MM}`/`{SEQ}` |
| `is_active` | bool | Status aktif (default `true`) |

### Placeholder Pattern
| Placeholder | Diganti jadi |
|---|---|
| `{YYYY}` | Tahun 4 digit (`2026`) |
| `{MM}` | Bulan 2 digit (`06`) |
| `{SEQ}` | Nomor urut 4 digit (`0001`) |

Detail mekanisme: [Deep Dive Penomoran](../guides/bast-numbering-deep-dive.md).

---

## Get All
<a id="get-all"></a>

### Request
```http
GET /api/bast-formats
Authorization: Bearer <token>
```

### Response — `200 OK`
```json
[
  {
    "format_id": "fmt-...",
    "format_name": "Format PO Standar",
    "format_type": "PO",
    "format_pattern": "BAST/PO/{YYYY}/{MM}/{SEQ}",
    "is_active": true,
    ...
  },
  {
    "format_id": "fmt-...",
    "format_name": "Format Internal Perusahaan",
    "format_type": "Internal",
    "format_pattern": "BAST/INT/{YYYY}/{MM}/{SEQ}",
    "is_active": true,
    ...
  }
]
```

### Contoh curl
```bash
curl -H "Authorization: Bearer <TOKEN>" http://localhost:8080/api/bast-formats
```

---

## Get By ID
<a id="get-by-id"></a>

### Request
```http
GET /api/bast-formats/:id
Authorization: Bearer <token>
```

### Response — `200 OK`
Object format lengkap.

### `404 Not Found`
```json
{ "error": "Format not found" }
```

---

## Create
<a id="create"></a>

### Request
```http
POST /api/bast-formats
Authorization: Bearer <token>
Content-Type: application/json
```

**Body:**
```json
{
  "format_name": "Format Project Khusus",
  "format_type": "Internal",
  "format_pattern": "BAST/KHUSUS/{YYYY}/{SEQ}",
  "is_active": true
}
```

### Response — `201 Created`
Object format lengkap (dengan `format_id` baru).

### Contoh curl
```bash
curl -X POST http://localhost:8080/api/bast-formats \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"format_name":"Format Khusus","format_type":"Internal","format_pattern":"BAST/X/{YYYY}/{SEQ}","is_active":true}'
```

---

## Update
<a id="update"></a>

### Request
```http
PUT /api/bast-formats/:id
Authorization: Bearer <token>
Content-Type: application/json
```

**Body:** field yang ingin diubah.
```json
{
  "format_name": "Format Project Khusus v2",
  "format_type": "Internal",
  "format_pattern": "BAST/X/{YYYY}/{MM}/{SEQ}",
  "is_active": true
}
```

### Response — `200 OK`
Object format setelah update.

---

## Delete
<a id="delete"></a>

Soft-deactivate: ubah `is_active` jadi `false`. Format tidak benar-benar dihapus (karena mungkin sudah dipakai BAST Request sebelumnya).

### Request
```http
DELETE /api/bast-formats/:id
Authorization: Bearer <token>
```

### Response — `200 OK`
```json
{ "message": "Format deleted successfully" }
```

---

## Contoh Format Pattern Populer

| Pola | Contoh Hasil |
|---|---|
| `BAST/INT/{YYYY}/{MM}/{SEQ}` | `BAST/INT/2026/06/0001` |
| `BAST/PO/{YYYY}/{SEQ}` | `BAST/PO/2026/0001` |
| `{YYYY}-{MM}-INT-{SEQ}` | `2026-06-INT-0001` |
| `BAST/HRD/{YYYY}/{SEQ}` | `BAST/HRD/2026/0001` |

---

## Tips

- Setelah buat format, catat `format_id` — dibutuhkan saat membuat `POST /bast-requests`.
- Format `is_active: false` masih bisa dipakai (saat ini tidak ada filter otomatis). Pertimbangkan validasi di service jika ingin hanya format aktif yang dipakai.

---

← **[Daftar Referensi API](README.md)** · 🤿 **[Deep Dive Penomoran](../guides/bast-numbering-deep-dive.md)**
