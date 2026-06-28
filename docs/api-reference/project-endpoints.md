# Project Endpoints

Operasi CRUD master data Project. Setiap project terikat ke satu customer. Sumber kode: [`project_handler.go`](../../internal/handlers/project_handler.go), [`project_service.go`](../../internal/services/project_service.go), [`project_repository.go`](../../internal/repositories/project_repository.go).

> 🔑 **Terlindungi.** Endpoint **POST (create)** hanya untuk `admin`/`superadmin`.

---

## Skema: Project

```json
{
  "project_id": "abc-123-project-id",
  "customer_id": "550e8400-customer-id",
  "customer": { ... },          // nested (hasil Preload)
  "project_code": "PRJ-MM-01",
  "project_name": "Implementasi Sistem ERP Terpadu",
  "status": "active",
  "created_at": "...",
  "updated_at": "...",
  "deleted_at": null
}
```

| Field | Tipe | Keterangan |
|---|---|---|
| `project_id` | UUID | Primary key |
| `customer_id` | UUID | **FK** ke Customer (wajib) |
| `customer` | object | Data customer terkait (otomatis di-Preload saat GET) |
| `project_code` | string | Unique |
| `project_name` | string | Nama project |
| `status` | string | `active`/`inactive` |

---

## Get All
<a id="get-all"></a>

### Request
```http
GET /api/projects?customer_id=<uuid>
Authorization: Bearer <token>
```

**Query param:**
| Param | Tipe | Wajib | Keterangan |
|---|---|---|---|
| `customer_id` | string | ❌ | Filter project milik customer tertentu |

### Response — `200 OK`
Array of Project (dengan nested `customer`).

### Contoh curl
```bash
curl -H "Authorization: Bearer <TOKEN>" \
  "http://localhost:8080/api/projects?customer_id=550e8400-..."
```

---

## Get By ID
<a id="get-by-id"></a>

### Request
```http
GET /api/projects/:id
Authorization: Bearer <token>
```

### Response — `200 OK`
Object project lengkap (dengan nested customer).

### `404 Not Found`
```json
{ "error": "Project not found" }
```

---

## Create
<a id="create"></a>

> 🛡️ **Admin only** (`admin`/`superadmin`).

### Request
```http
POST /api/projects
Authorization: Bearer <token-admin>
Content-Type: application/json
```

**Body:**
```json
{
  "customer_id": "550e8400-customer-id",
  "project_code": "PRJ-NEW-01",
  "project_name": "Pengembangan Aplikasi Mobile",
  "status": "active"
}
```
> `customer_id` wajib ada & valid (harus exist di tabel customer).

### Response — `201 Created`
Object project lengkap.

### Contoh curl
```bash
curl -X POST http://localhost:8080/api/projects \
  -H "Authorization: Bearer <TOKEN-ADMIN>" \
  -H "Content-Type: application/json" \
  -d '{"customer_id":"550e8400-...","project_code":"PRJ-NEW-01","project_name":"App Mobile","status":"active"}'
```

---

## Update
<a id="update"></a>

### Request
```http
PUT /api/projects/:id
Authorization: Bearer <token>
Content-Type: application/json
```

**Body:** field yang ingin diubah (`project_code`, `project_name`, `status`).
```json
{
  "project_code": "PRJ-NEW-01",
  "project_name": "App Mobile v2",
  "status": "active"
}
```

> ⚠️ `customer_id` **tidak** ikut di-update (lihat [`project_service.go:38`](../../internal/services/project_service.go) — sengaja, project tidak pindah customer).

### Response — `200 OK`
Object project setelah update.

---

## Delete
<a id="delete"></a>

Soft-deactivate: ubah `status` jadi `inactive`.

### Request
```http
DELETE /api/projects/:id
Authorization: Bearer <token>
```

### Response — `200 OK`
```json
{ "message": "Project deleted successfully" }
```

---

## Alur Umum Penggunaan

1. Pastikan customer target sudah ada (buat dulu via `POST /customers` kalau belum).
2. Ambil `customer_id` dari customer tersebut.
3. Buat project dengan `customer_id` itu.

---

← **[Daftar Referensi API](README.md)**
