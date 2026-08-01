package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	visitors = make(map[string]*visitor)
	mu       sync.Mutex
)

func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	v, exists := visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(rate.Every(3*time.Minute), 5)
		visitors[ip] = &visitor{
			limiter:  limiter,
			lastSeen: time.Now(),
		}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {

		limiter := getVisitor(c.ClientIP())

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
			})
			return
		}

		c.Next()
	}
}

func StartCleanup() {
	go func() {
		for {
			time.Sleep(time.Hour)

			mu.Lock()

			for ip, v := range visitors {
				if time.Since(v.lastSeen) > time.Hour {
					delete(visitors, ip)
				}
			}

			mu.Unlock()
		}
	}()
}
