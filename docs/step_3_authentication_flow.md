# Step 3: Alur Autentikasi (JWT & Role)

Dokumen ini menjelaskan bagaimana fitur registrasi dan login bekerja, mulai dari ujung belakang (Database) hingga diterjemahkan ke *Token* di ujung depan (Client).

## 1. Pembuatan Hashing Password (`internal/utils/hash.go`)
Menyimpan *password* di database dalam bentuk teks asli sangat berbahaya. Oleh karena itu, kita selalu menggunakan *algoritma hashing satu arah* bernama `bcrypt`.
```go
package utils
import "golang.org/x/crypto/bcrypt"

// Fungsi ini mengacak teks sandi menjadi Hash panjang
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Fungsi ini membandingkan Hash di database dengan teks sandi yang diketik user
func ValidatePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
```

## 2. Pembangkitan JWT Token (`internal/utils/jwt.go`)
JWT (*JSON Web Token*) adalah "karcis" agar user tidak perlu login setiap kali meminta data. JWT ini memuat ID User dan nama Role-nya.

```go
package utils
import "github.com/golang-jwt/jwt/v5"

func GenerateToken(userID string, roleName string) (string, error) {
	// Memasukkan klaim (data) yang akan disembunyikan di dalam token
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    roleName,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token mati setelah 24 jam
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Ditandatangani menggunakan "Kunci Rahasia" (Secret Key)
	return token.SignedString([]byte("rahasia-negara")) 
}
```

## 3. Register & Login (`internal/services/auth_service.go`)
Ini adalah otak dari proses otentikasi. Service ini akan memanggil `Repository` untuk mencari user, dan menggunakan `utils` untuk mencocokkan password dan membuat JWT.

```go
func (s *AuthService) Login(email, password string) (string, error) {
	// 1. Cek di Database, adakah email ini?
	user, err := s.authRepo.FindUserByEmail(email)
	if err != nil {
		return "", errors.New("email tidak ditemukan")
	}

	// 2. Cocokkan Hash Password-nya
	if !utils.ValidatePassword(password, user.PasswordHash) {
		return "", errors.New("password salah")
	}

	// 3. Password benar! Buatkan Karcis Token JWT
	token, err := utils.GenerateToken(user.UserID.String(), user.Role.Name)
	return token, err
}
```

## 4. Satpam API: JWT Middleware (`internal/middlewares/auth_middleware.go`)
Bahkan jika token berhasil diberikan kepada *User*, sistem harus memiliki cara untuk mengecek setiap request baru untuk melihat apakah karcis tersebut asli, dan peran apa yang dimiliki user tersebut. Inilah tugas **Middleware**.

Middleware pertama `RequireAuth()` mengecek apakah *Token* JWT asli atau palsu. Jika asli, ia menyimpan informasi (misal ID dan Role) tersebut di memori (Context) untuk *request* itu, sehingga bisa dibaca oleh *Handler*.

Middleware kedua `RequireRole()` bertindak seperti pengecekan level VIP.
```go
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Mengambil nama 'role' dari pengecekan Token di middleware sebelumnya
		userRole, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Role tidak ditemukan"})
			return
		}

		// Mengecek apakah Role user ada di dalam daftar VIP (allowedRoles)?
		isAllowed := false
		for _, allowed := range allowedRoles {
			if userRole == allowed {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Anda tidak berhak mengakses rute ini"})
			return
		}

		// Silakan Lewat!
		c.Next()
	}
}
```

## Ringkasan Step 3:
Otentikasi di *backend* Go melibatkan Hash Password, Pembangkitan JWT, dan verifikasi menggunakan Middleware sebelum sebuah Endpoint diakses. 

Lanjut ke: **[Step 4: Master Data CRUD](step_4_master_data_crud.md)**
