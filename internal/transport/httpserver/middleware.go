package httpserver

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger middleware
func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		d := time.Since(start)
		log.Printf("%s %s -> %d (%s)", c.Request.Method, c.Request.URL.Path, c.Writer.Status(), d)
		c.Writer.Header().Set("X-Response-Time", d.String())
	}
}

// Auth middleware (Bearer <token>)
func (h *Handler) authMW(c *gin.Context) {
	hdr := c.GetHeader("Authorization")
	if !strings.HasPrefix(strings.ToLower(hdr), "bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
		return
	}
	token := strings.TrimSpace(hdr[7:])
	phone, err := h.jwt.Parse(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}
	c.Set("phone", phone)
	c.Next()
}
