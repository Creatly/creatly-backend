package limiter

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/creatly-backend/pkg/logger"
	"golang.org/x/time/rate"
)

// visitor holds limiter and lastSeen for specific user.
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// rateLimiter used to rate limit an incoming requests.
type rateLimiter struct {
	sync.RWMutex

	visitors map[string]*visitor
	limit    rate.Limit
	burst    int
	ttl      time.Duration
}

// newRateLimiter creates an instance of the rateLimiter.
func newRateLimiter(rps, burst int, ttl time.Duration) *rateLimiter {
	return &rateLimiter{
		visitors: make(map[string]*visitor),
		limit:    rate.Limit(rps),
		burst:    burst,
		ttl:      ttl,
	}
}

// getVisitor returns limiter for the specific visitor by its IP,
// looking up within the visitors map.
func (l *rateLimiter) getVisitor(ip string) *rate.Limiter {
	l.RLock()
	v, exists := l.visitors[ip]
	l.RUnlock()

	if !exists {
		limiter := rate.NewLimiter(l.limit, l.burst)
		l.Lock()
		l.visitors[ip] = &visitor{limiter, time.Now()}
		l.Unlock()

		return limiter
	}

	v.lastSeen = time.Now()

	return v.limiter
}

// cleanupVisitors removes old entries from the visitors map.
func (l *rateLimiter) cleanupVisitors() {
	for {
		time.Sleep(time.Minute)

		l.Lock()
		for ip, v := range l.visitors {
			if time.Since(v.lastSeen) > l.ttl {
				delete(l.visitors, ip)
			}
		}
		l.Unlock()
	}
}

// Limit creates a new rate limiter middleware handler.
func Limit(rps int, burst int, ttl time.Duration) gin.HandlerFunc {
	l := newRateLimiter(rps, burst, ttl)

	// run a background worker to clean up old entries
	go l.cleanupVisitors()

	return func(c *gin.Context) {
		ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
		if err != nil {
			logger.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)

			return
		}

		if !l.getVisitor(ip).Allow() {
			c.AbortWithStatus(http.StatusTooManyRequests)

			return
		}

		c.Next()
	}
}
