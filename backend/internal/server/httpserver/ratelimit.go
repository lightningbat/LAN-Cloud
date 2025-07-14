package httpserver

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var ipLimiter = cache.New(1*time.Minute, 5*time.Minute)

func getRateLimiter(ip string, limit int, duration time.Duration) bool {
	countRaw, found := ipLimiter.Get(ip)
	if !found {
		ipLimiter.Set(ip, 1, duration)
		return true
	}

	count := countRaw.(int)
	if count >= limit {
		return false
	}
	ipLimiter.Set(ip, count+1, duration)
	return true
}
