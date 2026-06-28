# 🔌 Referensi API — BAST Request API

Dokumentasi teknis tiap endpoint: method, path, parameter, contoh request/response, dan aturan akses.

> 💡 Untuk UI interaktif (bisa langsung klik & test), gunakan **Swagger UI** di `http://localhost:8080/swagger/index.html`. Dokumen ini adalah referensi tertulis untuk dibaca offline / dokumentasi resmi.

---

## Base URL

```
http://localhost:8080/api
```

Semua endpoint (kecuali `/ping` dan `/swagger/*`) berada di bawah prefix `/api`.

---

## Autentikasi

Mayoritas endpoint butuh token JWT. Cara mendapatkannya:

1. **Register** akun: `POST /api/auth/register`
2. **Login**: `POST /api/auth/login` → dapat `token` di response.
3. Sertakan token di header setiap request berikutnya:
   ```
   Authorization: Bearer <token-anda>
   ```

> ⚠️ Format header **wajib** `Bearer ` (dengan spasi) + token. Tanpa ini → 401 Unauthorized.

---

## Tingkat Akses (RBAC)

| Tingkat | Syarat | Endpoint |
|---|---|---|
| 🔓 **Publik** | Tanpa token | `/auth/register`, `/auth/login`, `/ping`, `/swagger/*` |
| 🔑 **Terlindungi** | Token valid (role apa pun) | Mayoritas endpoint (customers, projects, dll) |
| 🛡️ **Admin** | Token + role `admin`/`superadmin` | `POST /projects`, `DELETE /customers/:id` |

Detail tiap endpoint di file masing-masing modul di bawah.

---

## Status Code Umum

| Code | Arti |
|---|---|
| `200` | OK (GET, PUT, PATCH, DELETE sukses) |
| `201` | Created (POST sukses) |
| `400` | Bad Request (body JSON salah / validasi gagal) |
| `401` | Unauthorized (token tidak ada / invalid) |
| `403` | Forbidden (role tidak diizinkan) |
| `404` | Not Found (resource tidak ditemukan) |
| `500` | Internal Server Error (error server / DB) |

---

## Daftar Endpoint per Modul

### 🔐 Auth (Publik)
| Method | Path | Deskripsi |
|---|---|---|
| POST | [`/auth/register`](auth-endpoints.md#register) | Daftar user baru |
| POST | [`/auth/login`](auth-endpoints.md#login) | Login → dapat token |

### 👥 Customer (Terlindungi, delete = Admin)
| Method | Path | Deskripsi |
|---|---|---|
| GET | [`/customers`](customer-endpoints.md#get-all) | List customer (filter) |
| GET | [`/customers/:id`](customer-endpoints.md#get-by-id) | Detail customer |
| POST | [`/customers`](customer-endpoints.md#create) | Buat customer |
| PUT | [`/customers/:id`](customer-endpoints.md#update) | Update customer |
| DELETE | [`/customers/:id`](customer-endpoints.md#delete) | Hapus (admin only) |

### 📁 Project
| Method | Path | Deskripsi |
|---|---|---|
| GET | [`/projects`](project-endpoints.md#get-all) | List project (filter) |
| GET | [`/projects/:id`](project-endpoints.md#get-by-id) | Detail project |
| POST | [`/projects`](project-endpoints.md#create) | Buat project (admin) |
| PUT | [`/projects/:id`](project-endpoints.md#update) | Update project |
| DELETE | [`/projects/:id`](project-endpoints.md#delete) | Hapus project |

### 📋 BAST Format
| Method | Path | Deskripsi |
|---|---|---|
| GET | [`/bast-formats`](bast-format-endpoints.md#get-all) | List format |
| GET | [`/bast-formats/:id`](bast-format-endpoints.md#get-by-id) | Detail format |
| POST | [`/bast-formats`](bast-format-endpoints.md#create) | Buat format |
| PUT | [`/bast-formats/:id`](bast-format-endpoints.md#update) | Update format |
| DELETE | [`/bast-formats/:id`](bast-format-endpoints.md#delete) | Nonaktifkan format |

### 🔢 BAST Sequence
| Method | Path | Deskripsi |
|---|---|---|
| GET | [`/bast-sequences`](bast-sequence-endpoints.md#get) | Cek nomor urut |
| POST | [`/bast-sequences/reset`](bast-sequence-endpoints.md#reset) | Reset nomor urut |

### 📄 BAST Request
| Method | Path | Deskripsi |
|---|---|---|
| GET | [`/bast-requests`](bast-request-endpoints.md#get-all) | List request |
| GET | [`/bast-requests/:id`](bast-request-endpoints.md#get-by-id) | Detail request |
| POST | [`/bast-requests`](bast-request-endpoints.md#create) | Buat request (auto nomor) |
| PATCH | [`/bast-requests/:id/status`](bast-request-endpoints.md#update-status) | Ubah status |
| GET | [`/bast-requests/:id/audit`](bast-request-endpoints.md#audit) | Audit request |

### 📜 Audit Log
| Method | Path | Deskripsi |
|---|---|---|
| GET | [`/audit-logs`](audit-log-endpoints.md#get-all) | List audit (filter) |
| GET | [`/audit-logs/:id`](audit-log-endpoints.md#get-by-id) | Detail audit |

---

## Format Request/Response

- **Content-Type:** `application/json`
- **Body:** JSON object/array.
- **Field penamaan:** `snake_case` (mis. `customer_id`, `po_number`) — sesuai tag `json:"..."` di model.
- **Tanggal:** ISO 8601 (mis. `2026-06-28T10:30:00Z`).
- **UUID:** format standar (mis. `550e8400-e29b-41d4-a716-446655440000`).

---

## Contoh Error Response

Semua error mengikuti format:
```json
{
  "error": "pesan error dalam bahasa Inggris/Indonesia"
}
```

---

## Tips Menguji

- **Postman / Insomnia:** import file `docs/swagger.json` untuk auto-generate collection.
- **curl:** lihat contoh di tiap file endpoint.
- **Swagger UI:** paling cepat untuk eksplorasi (ada tombol "Try it out").

---

## Bacaan Lanjutan
- 🛠️ [Panduan Swagger](../guides/swagger-guide.md)
- 🔐 [Panduan Autentikasi](../guides/authentication-guide.md)
- 🏗️ [Clean Architecture](../architecture/clean-architecture.md)
