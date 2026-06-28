# Step 3: Alur Autentikasi (JWT & Role)

> Seri Tutorial · **Step 3 dari 8**

Database sudah ada tabelnya. Sekarang kita pasang **pengaman**: bagaimana user register, login, dapat token, dan bagaimana sistem memverifikasi setiap request berikutnya. Ini melibatkan **bcrypt** (hashing password), **JWT** (token), **Service Auth**, dan **Middleware**.

---

## 1. Mengapa Perlu Autentikasi?

REST API bersifat **stateless** (pelupa) — server tidak mengingat siapa yang barusan datang. Maka:
- Login tidak bisa pakai *session* seperti web PHP biasa.
- Solusinya: setelah login, server kasih **token (karcis)** ke klien. Klien wajib bawa token tsb di setiap request berikutnya.

Kita pakai **JWT (JSON Web Token)** sebagai "karcis" itu.

---

## 2. Hashing Password — `internal/utils/hash.go`

File: [`internal/utils/hash.go`](../../internal/utils/hash.go)

```go
package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)   // (1) cost = 14
	return string(bytes), err
}

func ValidatePassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))  // (2)
	return err == nil
}
```

### Penjelasan
- **(1) `GenerateFromPassword`** → mengubah password plain-text jadi **hash** acak. Angka `14` adalah **cost factor** (semakin tinggi, semakin aman & lambat). Hash hasilnya berbeda setiap kali, walau input sama.
- **(2) `CompareHashAndPassword`** → membandingkan hash di DB dengan password yang diketik user. Return `nil` kalau cocok.

> ⚠️ **Aturan emas:** JANGAN PERNAH simpan password plain-text di database. Selalu hash. Bcrypt bersifat **satu arah** — tidak bisa di-decode, hanya bisa dibandingkan.

---

## 3. Membuat & Memvalidasi JWT — `internal/utils/jwt.go`

File: [`internal/utils/jwt.go`](../../internal/utils/jwt.go)

```go
package utils

import (
	"fmt"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("bast-request-secret-key-05")   // (1) secret key

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims                              // (2) klaim bawaan JWT
}

func GenerateToken(userID string, roleName string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)  // (3) token berlaku 24 jam
	claims := &Claims{
		UserID: userID,
		Role:   roleName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "bast-request",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)   // (4) HS256
	return token.SignedString(jwtKey)                            // (5) tanda tangani
}

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
```

### Penjelasan
- **(1) `jwtKey`** — "kunci rahasia" untuk menandatangani token. Hanya server yang tahu. **Siapa pun yang pegang kunci ini bisa memalsukan token** → wajib dijaga rahasia (di produksi: simpan di env variable).
- **(2) `Claims`** — struct berisi data yang "ditelan" ke dalam token. Di sini: `UserID` & `Role`.
- **(3) `ExpiresAt`** — token **kedaluwarsa** 24 jam. Setelah itu wajib login ulang.
- **(4) `SigningMethodHS256`** — algoritma HMAC-SHA256.
- **(5) `SignedString`** — token final berbentuk string panjang `xxx.yyy.zzz`.

Token ini lalu dikirim klien sebagai `Authorization: Bearer <token>`.

> 💡 **Anatomi JWT:** terdiri dari 3 bagian dipisah titik: **header.signature.payload**. Payload ter-encode base64 (bukan terenkripsi!) — jadi JANGAN taruh data sensitif di JWT.

---

## 4. Model Permintaan Auth — `internal/models/auth.go`

File: [`internal/models/auth.go`](../../internal/models/auth.go)

```go
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`   // min 6 karakter
	Role     string `json:"role" binding:"required"`
}
```

Tag `binding:"..."` adalah **validator** bawaan Gin. `required,email` → wajib diisi & format email valid. Jika gagal, Gin otomatis kirim error 400.

---

## 5. Auth Service — Otak Autentikasi

File: [`internal/services/auth_service.go`](../../internal/services/auth_service.go)

### Fungsi `Register`
```go
func (s *AuthService) Register(req *models.RegisterRequest) (string, error) {
	// (1) Cek apakah email sudah dipakai
	_, err := s.authRepo.FindUserByEmail(req.Email)
	if err == nil {
		return "", fmt.Errorf("user already exists")
	}

	// (2) Cari Role ID berdasarkan nama role
	role, err := s.authRepo.FindRoleByName(req.Role)
	if err != nil {
		return "", fmt.Errorf("invalid role: %v", err)
	}

	// (3) Hash password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return "", err
	}

	// (4) Buat user baru
	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
		RoleID:       role.RoleID,
	}
	err = s.authRepo.CreateUser(user)
	if err != nil {
		return "", err
	}
	return "User registered successfully", nil
}
```

**Alur:** cek duplikat → cari role → hash password → simpan user.

### Fungsi `Login`
```go
func (s *AuthService) Login(email, password string) (string, error) {
	// (1) Cari user by email (Preload Role agar nama role ikut)
	user, err := s.authRepo.FindUserByEmail(email)
	if err != nil {
		return "", fmt.Errorf("invalid email or password")   // ← jangan bocorin mana yang salah
	}

	// (2) Bandingkan password
	if !utils.ValidatePassword(user.PasswordHash, password) {
		return "", fmt.Errorf("invalid email or password")   // ← pesan sama, anti brute-force
	}

	// (3) Generate token
	token, err := utils.GenerateToken(user.UserID.String(), user.Role.Name)
	if err != nil {
		return "", err
	}
	return token, nil
}
```

### Tips Keamanan di kode ini
Perhatikan **pesan error identik** ("invalid email or password") baik untuk email salah maupun password salah. Ini disengaja: **jangan kasih tahu penyerang** bagian mana yang salah. Mencegah *user enumeration*.

### Repository — `internal/repositories/auth_repository.go`
```go
func (r *AuthRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Role").Where("email = ?", email).First(&user).Error
	// ...
}
```
`Preload("Role")` penting agar `user.Role.Name` terisi — dibutuhkan saat generate token.

---

## 6. Auth Handler — Penerima HTTP

File: [`internal/handlers/auth_handler.go`](../../internal/handlers/auth_handler.go)

```go
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {        // (1) parse + validasi
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := h.authService.Login(req.Email, req.Password)  // (2) panggil service
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})         // (3) kirim token
}
```

Handler hanya terjemahkan HTTP ↔ Service. Tidak ada logika bisnis di sini.

---

## 7. Satpam API: Middleware

File: [`internal/middlewares/auth_middleware.go`](../../internal/middlewares/auth_middleware.go)

Middleware berjalan **sebelum** Handler. Tugasnya: cek token di setiap request.

### `RequireAuth()` — Cek Token
```go
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header is required"})
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}

		// Simpan info user ke context agar handler bisa baca
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)
		c.Next()   // ← lanjut ke handler
	}
}
```

**Pola penting:**
- `c.AbortWithStatusJSON(...)` → hentikan request, kirim error, **jangan lanjut**.
- `c.Set("key", value)` → simpan data untuk dipakai handler berikutnya.
- `c.Next()` → teruskan request.

### `RequireRole(...)` — Cek Level Akses
```go
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, _ := c.Get("userRole")   // ambil dari RequireAuth tadi

		isAllowed := false
		for _, role := range allowedRoles {
			if role == userRole {
				isAllowed = true
				break
			}
		}

		if isAllowed {
			c.Next()
		} else {
			c.AbortWithStatusJSON(403, gin.H{"error": "Anda tidak memiliki akses (Forbidden)"})
			return
		}
	}
}
```
`...string` = **variadic** (bisa terima banyak argumen, mis. `RequireRole("admin", "superadmin")`).

---

## 8. Alur Lengkap Autentikasi

```mermaid
sequenceDiagram
    participant K as Klien
    participant A as /auth/login
    participant S as AuthService
    participant DB as Database
    participant M as Middleware
    participant H as Handler Terlindungi

    Note over K,A: 1. LOGIN
    K->>A: POST {email, password}
    A->>S: Login(email, password)
    S->>DB: FindUserByEmail(+Preload Role)
    DB-->>S: user + hash
    S->>S: ValidatePassword(hash, password)
    S->>S: GenerateToken(userID, role)
    S-->>A: token JWT
    A-->>K: 200 {token}

    Note over K,H: 2. AKSES ENDPOINT TERLINDUNGI
    K->>M: GET /api/customers (Header: Bearer <token>)
    M->>M: ValidateToken(token)
    M->>M: c.Set("userRole", role)
    M->>H: c.Next()
    H-->>K: 200 [data]
```

---

## 9. Uji Coba dengan curl

### Register
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"admin1","email":"admin1@test.com","password":"rahasia123","role":"admin"}'
```
**Respons:** `{"token":"User registered successfully"}` *(catatan: field key-nya "token" tapi isinya pesan sukses — detail di [Referensi Auth](../api-reference/auth-endpoints.md).)*

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin1@test.com","password":"rahasia123"}'
```
**Respons:** `{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6..."}` → simpan token ini.

### Akses endpoint terlindungi
```bash
curl http://localhost:8080/api/customers \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6..."
```

---

## ✅ Ringkasan Step 3
- Password **dihash** dengan bcrypt (tidak disimpan plain-text).
- Login menghasilkan **JWT** berisi `UserID` & `Role`, berlaku 24 jam.
- Middleware `RequireAuth` mengecek token setiap request; `RequireRole` membatasi akses berdasarkan role.
- Pesan error identik ("invalid email or password") mencegah *user enumeration*.
- Info user disimpan di Gin Context (`c.Set`) agar bisa dibaca handler.

Autentikasi siap. Sekarang, bagaimana operasi data (CRUD) bekerja setelah user lolos satpam?

---

⬅️ **[Step 2: Models & Migrasi](step-02-models-and-migration.md)** · ➡️ **[Step 4: Master Data CRUD](step-04-master-data-crud.md)**

> 📖 Butuh panduan menyeluruh tentang auth? Lihat [Panduan Autentikasi & RBAC](../guides/authentication-guide.md).
