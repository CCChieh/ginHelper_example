package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		fmt.Printf("%s %s %s %d\n",
			c.Request.Method,
			c.Request.RequestURI,
			time.Since(start),
			c.Writer.Status(),
		)
	}
}

