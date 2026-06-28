# Customer Endpoints

Operasi CRUD master data Customer. Sumber kode: [`internal/handlers/customer_handler.go`](../../internal/handlers/customer_handler.go), [`services/customer_service.go`](../../internal/services/customer_service.go), [`repositories/customer_repository.go`](../../internal/repositories/customer_repository.go).

> 🔑 **Terlindungi** — semua endpoint butuh token JWT. Endpoint **DELETE** hanya untuk role `admin`/`superadmin`.

---

## Skema: Customer

```json
{
  "customer_id": "550e8400-e29b-41d4-a716-446655440000",
  "customer_code": "CUST-001",
  "customer_name": "PT. Maju Mundur",
  "status": "active",
  "created_at": "2026-06-01T08:00:00Z",
  "updated_at": "2026-06-01T08:00:00Z",
  "deleted_at": null
}
```

| Field | Tipe | Keterangan |
|---|---|---|
| `customer_id` | UUID | Primary key (auto-generated) |
| `customer_code` | string | **Unique**, mis. `CUST-001` |
| `customer_name` | string | Nama lengkap customer |
| `status` | string | `active` atau `inactive` (default `active`) |
| `created_at`/`updated_at` | datetime | Otomatis |
| `deleted_at` | datetime/null | Soft delete |

---

## Get All
<a id="get-all"></a>

Ambil daftar customer, dengan filter opsional.

### Request
```http
GET /api/customers?status=active&nama=Maju
Authorization: Bearer <token>
```

**Query param:**
| Param | Tipe | Wajib | Keterangan |
|---|---|---|---|
| `status` | string | ❌ | Filter status (`active`/`inactive`) |
| `nama` | string | ❌ | Filter nama (fuzzy, case-insensitive) |

> ⚠️ Param pencarian namanya **`nama`**, bukan `name`.

### Response — `200 OK`
```json
[
  {
    "customer_id": "550e8400-...",
    "customer_code": "CUST-001",
    "customer_name": "PT. Maju Mundur",
    "status": "active",
    "created_at": "...",
    "updated_at": "...",
    "deleted_at": null
  }
]
```

### Contoh curl
```bash
curl -H "Authorization: Bearer <TOKEN>" \
  "http://localhost:8080/api/customers?status=active&nama=Maju"
```

---

## Get By ID
<a id="get-by-id"></a>

### Request
```http
GET /api/customers/:id
Authorization: Bearer <token>
```
| Param | Lokasi | Keterangan |
|---|---|---|
| `id` | path | UUID customer |

### Response — `200 OK`
Lihat skema di atas.

### `404 Not Found`
```json
{ "error": "Customer not found" }
```

---

## Create
<a id="create"></a>

### Request
```http
POST /api/customers
Authorization: Bearer <token>
Content-Type: application/json
```

**Body:**
```json
{
  "customer_code": "CUST-003",
  "customer_name": "PT. Sukses Jaya",
  "status": "active"
}
```
> `customer_id` auto-generated (jangan kirim). `status` opsional (default `active`).

### Response — `201 Created`
Object customer lengkap (termasuk `customer_id` yang baru di-generate).

### `400` / `500`
```json
{ "error": "pesan error" }
```

---

## Update
<a id="update"></a>

### Request
```http
PUT /api/customers/:id
Authorization: Bearer <token>
Content-Type: application/json
```

**Body:** field yang ingin diubah (`customer_code`, `customer_name`, `status`).
```json
{
  "customer_code": "CUST-003",
  "customer_name": "PT. Sukses Jaya Abadi",
  "status": "active"
}
```

### Response — `200 OK`
Object customer setelah update.

---

## Delete
<a id="delete"></a>

> 🛡️ **Admin only** — butuh role `admin`/`superadmin`.

**Catatan:** ini **soft-deactivate**, bukan hapus permanen. Mengubah `status` jadi `inactive`. Data tetap ada di DB.

### Request
```http
DELETE /api/customers/:id
Authorization: Bearer <token>
```

### Response — `200 OK`
```json
{ "message": "Customer deleted successfully" }
```

### `403 Forbidden` (jika role `user`)
```json
{ "error": "Anda tidak memiliki akses (Forbidden)" }
```

### Contoh curl
```bash
curl -X DELETE http://localhost:8080/api/customers/<ID> \
  -H "Authorization: Bearer <TOKEN-ADMIN>"
```

---

← **[Daftar Referensi API](README.md)**
