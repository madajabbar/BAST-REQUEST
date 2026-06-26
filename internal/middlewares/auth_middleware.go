package middlewares

import (
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
	"bast-request/internal/utils"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}
		
		tokenString := strings.Replace(authHeader, "Bearer ","",1)
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// simpan data user dan role kedalam context agar bisa dibaca oleh handler nanti
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)
		c.Next()
	}
}

// RequireRole memblokir user yang rolenya tidak sesuai
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context){
		userRole, _ := c.Get("userRole")

		isAllowed := false
		for _, role := range allowedRoles {
			if role == userRole {
				isAllowed = true
				break
			}
		}

		if isAllowed {
			c.Next()
		}else{
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses (Forbidden)"})
			return
		}
	}
}