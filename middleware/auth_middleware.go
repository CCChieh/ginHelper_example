package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")
		if name=="admin"{
			c.Next()
		}else {
			c.JSON(http.StatusUnauthorized,gin.H{"message": "只能admin访问!"})
			c.Abort()
		}
	}
}

func UnAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")
		if name=="admin"{
			c.JSON(http.StatusUnauthorized,gin.H{"message": "admin不能访问!"})
			c.Abort()
		}else {
			c.Next()
		}
	}
}

