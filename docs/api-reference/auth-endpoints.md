# Auth Endpoints

Endpoint autentikasi: registrasi user baru dan login untuk mendapatkan token JWT.

> 🔓 **Publik** — tidak butuh token. Sumber kode: [`internal/handlers/auth_handler.go`](../../internal/handlers/auth_handler.go)

---

## Register

Daftarkan user baru. Password di-hash dengan bcrypt, lalu user disimpan dengan role yang ditentukan.

### Request

```http
POST /api/auth/register
Content-Type: application/json
```

**Body:**
| Field | Tipe | Wajib | Aturan |
|---|---|---|---|
| `username` | string | ✅ | - |
| `email` | string | ✅ | format email valid |
| `password` | string | ✅ | minimal 6 karakter |
| `role` | string | ✅ | salah satu: `superadmin`, `admin`, `user` |

**Contoh:**
```json
{
  "username": "admin1",
  "email": "admin1@test.com",
  "password": "rahasia123",
  "role": "admin"
}
```

### Response

**`200 OK`** — registrasi sukses.
```json
{
  "token": "User registered successfully"
}
```

> ⚠️ **Catatan implementasi:** field key-nya `"token"`, tapi **isinya pesan sukses**, bukan JWT. Untuk dapat JWT, harus login setelah register. Ini sedikit *inconsistency* di kode — bisa diperbaiki ke `{ "message": "..." }` atau langsung return JWT.

**`400 Bad Request`** — validasi gagal (field kosong / format email salah / password <6).
```json
{
  "error": "Key: 'RegisterRequest.Password' Error:Field validation for 'Password' failed on the 'min' tag"
}
```

**`500 Internal Server Error`** — email sudah dipakai, role tidak ditemukan, atau error DB.
```json
{ "error": "user already exists" }
```

### Contoh curl
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"admin1","email":"admin1@test.com","password":"rahasia123","role":"admin"}'
```

> 🚨 **Peringatan keamanan:** saat ini klien bisa memilih `role` apa pun (termasuk `superadmin`) saat register publik. Lihat [Panduan Autentikasi bagian Celah](../guides/authentication-guide.md#8-celah-keamanan-yang-perlu-diperbaiki).

---

## Login

Login dengan email & password → dapat token JWT (berlaku 24 jam).

### Request

```http
POST /api/auth/login
Content-Type: application/json
```

**Body:**
| Field | Tipe | Wajib | Aturan |
|---|---|---|---|
| `email` | string | ✅ | format email valid |
| `password` | string | ✅ | - |

**Contoh:**
```json
{
  "email": "admin1@test.com",
  "password": "rahasia123"
}
```

### Response

**`200 OK`** — login sukses, dapat token.
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoi..."
}
```

> 💡 Token ini (string panjang) yang dipakai di header `Authorization: Bearer <token>` untuk endpoint lain.

**`400 Bad Request`** — body salah / field kosong.
```json
{ "error": "Key: 'LoginRequest.Email' Error:..." }
```

**`500 Internal Server Error`** — email/password salah.
```json
{ "error": "invalid email or password" }
```

> 💡 **Catatan keamanan:** pesan error identik untuk "email tidak ditemukan" dan "password salah". Ini sengaja — mencegah *user enumeration* (penyerang menebak email mana yang terdaftar).

### Contoh curl
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin1@test.com","password":"rahasia123"}'
```

---

## Alur Setelah Login

1. Copy nilai `token` dari response.
2. Untuk semua endpoint terlindungi, tambah header:
   ```
   Authorization: Bearer <token>
   ```
3. Token berlaku **24 jam**. Setelah itu, login ulang.

---

## Decode Token (Opsional)

Payload JWT bisa di-decode (base64) untuk melihat isinya tanpa secret:
```bash
echo "<token>" | cut -d. -f2 | base64 -d 2>/dev/null
```
Output (contoh):
```json
{
  "user_id": "abc-123-def",
  "role": "admin",
  "exp": 1780000000,
  "iat": 1779900000,
  "iss": "bast-request"
}
```

> ⚠️ Payload **tidak terenkripsi** — hanya ter-encode. Jangan taruh data sensitif di JWT.

---

← **[Daftar Referensi API](README.md)** · 🔐 **[Panduan Autentikasi](../guides/authentication-guide.md)**
