package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := os.Getenv("ADMIN_SECRET")
		//TODO: ambil header Authorization
		auth := c.Request.Header.Get("Authorization")

		// TODO: cek apakah sesuai dengan sandi 
		if auth == "" {
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		if auth != key {
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		// TODO: jika sesuai, lanjutkan ke handler
		c.Next()
	}
}